package client

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/frame/pnet"
	"github.com/meow-pad/persian/frame/pnet/message"
	"github.com/meow-pad/persian/frame/pnet/tcp/session"
	"github.com/meow-pad/persian/frame/pnet/utils"
	"net/url"
	"sync/atomic"
)

const (
	StatusInitial = iota
	StatusConnecting
	StatusConnected
	StatusClosed
)

var (
	ErrInvalidStatus = errors.New("invalid client status")
)

func NewClient(codec message.Codec, listener session.Listener, opts ...Option) (*Client, error) {
	options := newOptions(opts...)
	client := &Client{}
	if err := client.init(codec, listener, options); err != nil {
		return nil, err
	}
	return client, nil
}

type Client struct {
	*Options
	session.BaseSession

	// 状态
	status atomic.Uint32
	// 编解码器
	codec message.Codec
	// 会话监听器
	listener session.Listener
	// 连接，仅连接成功后有值
	conn   *Conn
	loop   *eventLoop
	connPT atomic.Pointer[Conn]
}

func (client *Client) init(codec message.Codec, listener session.Listener, options *Options) error {
	if codec == nil || listener == nil {
		return errdef.ErrInvalidParams
	}
	client.Options = options
	client.status.Store(StatusInitial)
	client.codec = codec
	client.listener = listener
	client.conn = nil
	client.connPT.Store(nil)
	client.loop = newEventLoop(client)
	return nil
}

// Dial
//
//	@Description: 连接
//	@receiver client
//	@param ctx
//	@param address 如：127.0.0.1:9999
//	@return error
func (client *Client) Dial(ctx context.Context, wsUrl *url.URL) error {
	status := client.status.Load()
	if status != StatusInitial && status != StatusClosed {
		return ErrInvalidStatus
	}
	if !client.status.CompareAndSwap(status, StatusConnecting) {
		return ErrInvalidStatus
	}
	client.status.Store(StatusConnecting)
	//wsUrl := url.URL{Scheme: utils.ProtoWebsocket, Host: address, Path: "/"}
	if wsUrl == nil || wsUrl.Scheme != utils.ProtoWebsocket {
		return errors.New("invalid websocket url")
	}
	var dialer websocket.Dialer
	conn, _, err := dialer.DialContext(ctx, wsUrl.String(), nil)
	if err != nil {
		return err
	}
	client.conn, err = NewConn(client, conn)
	if err != nil {
		client.conn = nil
		client.status.Store(StatusInitial)
		if tErr := conn.Close(); tErr != nil {
			plog.Error("", pfield.Error(tErr))
		}
		return err
	}
	client.connPT.Store(client.conn)
	client.status.Store(StatusConnected)
	if err = client.loop.start(client.conn); err != nil {
		client.conn = nil
		client.connPT.Store(nil)
		client.status.Store(StatusInitial)
		if tErr := conn.Close(); tErr != nil {
			plog.Error("", pfield.Error(tErr))
		}
		return err
	}
	return nil
}

func (client *Client) Status() uint32 {
	return client.status.Load()
}

func (client *Client) Connection() session.Conn {
	return client.connPT.Load()
}

func (client *Client) Close() error {
	if client.status.Load() == StatusClosed {
		return pnet.ErrClosedClient
	}
	return client.loop.stop(context.Background(), nil)
}

func (client *Client) CloseWithContext(ctx context.Context) error {
	if client.status.Load() == StatusClosed {
		return pnet.ErrClosedClient
	}
	return client.loop.stop(ctx, nil)
}

func (client *Client) toClosed(reason error) bool {
	if client.status.CompareAndSwap(StatusConnected, StatusClosed) {
		if reason == nil {
			plog.Debug("close client", pfield.String("name", client.Name))
		} else {
			plog.Error("close client", pfield.String("name", client.Name), pfield.NamedError("reason", reason))
		}
		return true
	}
	return false
}

func (client *Client) IsClosed() bool {
	return client.status.Load() != StatusConnected
}

// SendMessage
//
//	@Description: 发送消息
//	@receiver client
//	@param message
func (client *Client) SendMessage(message any) {
	if client.status.Load() != StatusConnected {
		plog.Error("connect first")
		return
	}
	buf, err := client.codec.Encode(message)
	if err != nil {
		plog.Error("encode message error:", pfield.Error(err))
		return
	}
	bufLen := len(buf)
	err = client.loop.asyncWrite(buf, func(c *Conn, err error) error {
		if err != nil {
			client.onSendingError("write message error:", err)
			return nil
		}
		if err = client.listener.OnSend(client, message, bufLen); err != nil {
			plog.Error("on send error:", pfield.Error(err))
		}
		return nil
	})
	if err != nil {
		plog.Error("async write error:", pfield.Error(err))
	}
}

// SendMessages
//
//	@Description: 发送多条消息
//	@receiver client
//	@param messages
func (client *Client) SendMessages(messages ...any) {
	if client.status.Load() != StatusConnected {
		plog.Error("connect first")
		return
	}
	totalLen := 0
	dataArr := make([][]byte, 0, len(messages))
	for _, msg := range messages {
		data, err := client.codec.Encode(msg)
		if err != nil {
			client.onSendingError("encode msg error:", err)
			return
		}
		dataArr = append(dataArr, data)
		totalLen += len(data)
	}
	err := client.loop.asyncWritev(dataArr, func(c *Conn, err error) error {
		if err != nil {
			client.onSendingError("write msg error:", err)
			return nil
		}
		if err = client.listener.OnSendMulti(client, messages, totalLen); err != nil {
			plog.Error("on send multi error:", pfield.Error(err))
		}
		return nil
	})
	if err != nil {
		plog.Error("async writev error:", pfield.Error(err))
	}
}

// onSendingError
//
//	@Description: 发送消息时错误处理
//	@receiver sess
//	@param tip 日志消息
//	@param err 错误
func (client *Client) onSendingError(tip string, err error) {
	plog.Error(tip, pfield.Error(err))
	// 无法处理的状态，关闭连接
	cErr := client.Close()
	if cErr != nil {
		plog.Error("close conn error", pfield.Error(cErr))
	}
}
