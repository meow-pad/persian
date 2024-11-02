package client

import (
	"bytes"
	"context"
	"errors"
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/frame/pnet"
	"github.com/meow-pad/persian/frame/pnet/tcp/session"
	"github.com/meow-pad/persian/utils/coding"
)

const (
	loopBufferSize = 512
)

// 写事件
type writeEvent struct {
	buf      []byte
	bufS     [][]byte
	callback func(c session.Conn, err error) error
}

// newEventLoop
//
//	@Description: 构建事件循环处理
//	@param client 客户端
//	@return *eventLoop
func newEventLoop(client *Client) *eventLoop {
	return &eventLoop{
		client:  client,
		handler: newEventHandler(),
	}
}

type eventLoop struct {
	client  *Client
	handler *eventHandler

	conn       *Conn
	buffer     []byte       // 网络读缓存
	cache      bytes.Buffer // 读写临时缓存
	readChan   chan int
	readDone   chan struct{}
	writeChan  chan writeEvent
	closeChan  chan error
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
}

// run
//
//	@Description:
//	@receiver loop
func (loop *eventLoop) run(first bool) {
	defer coding.CatchPanicError("run error:", func() {
		if loop.client.IsClosed() {
			return
		}
		go loop.run(false)
	})
	// 启动监听
	if first {
		act := loop.handler.OnOpen(loop.conn)
		if act == actionClose {
			loop._stop(nil)
			return
		}
	}
	// 事件处理
	for {
		select {
		case n := <-loop.readChan:
			if loop.conn.inbound.Buffered() > loop.client.ReadBufferCap {
				// 缓存数据过多且未处理
				loop._stop(pnet.ErrOutOfReadCap)
				return
			}
			_, _ = loop.conn.inbound.Write(loop.buffer[:n]) // 目前的实现上不会返回err
			act := loop.handler.OnTraffic(loop.conn)
			if act == actionClose {
				loop._stop(nil)
				return
			}
			loop.readDone <- struct{}{}
		case event := <-loop.writeChan:
			if _, err := loop.conn.Write(event.buf); err != nil {
				loop._stop(err)
				return
			}
			if _, err := loop.conn.Writev(event.bufS); err != nil {
				loop._stop(err)
				return
			}
			// 有数据可写则发送
			if loop.conn.outbound.Len() > 0 {
				_, err := loop.conn.TCPConn.ReadFrom(&loop.conn.outbound)
				if event.callback != nil {
					if cErr := event.callback(loop.conn, err); cErr != nil {
						plog.Error("write callback error:", pfield.Error(cErr))
					}
				} else if err != nil {
					loop._stop(err)
					return
				}
			}
		case closeReason := <-loop.closeChan:
			loop._stop(closeReason)
			return
		case <-loop.cancelCtx.Done():
			// 正常也走不到这
			plog.Warn("cancel event loop running")
			return
		}
	}
}

// readConn
//
//	@Description: 持续读网络数据
//	@receiver loop
func (loop *eventLoop) readConn() {
	defer coding.CatchPanicError("read conn error:", func() {
		if loop.client.IsClosed() {
			return
		}
		go loop.readConn()
	})
	for {
		// 判定是否已关闭
		if loop.client.IsClosed() {
			return
		}
		// 读网络
		n, err := loop.conn.TCPConn.Read(loop.buffer)
		if err != nil {
			if sErr := loop.stop(context.Background(), err); sErr != nil && !errors.Is(sErr, pnet.ErrClosedClient) {
				plog.Error("", pfield.Error(sErr))
			}
			return
		}
		// 提交已读数据
		if n > 0 {
			// 等待读处理
			select {
			case loop.readChan <- n:
			case <-loop.cancelCtx.Done():
				return
			}
			// 等待读处理结束
			select {
			case <-loop.readDone:
			case <-loop.cancelCtx.Done():
				return
			}
		}
	} // end of for
}

func (loop *eventLoop) asyncWrite(b []byte, callback func(c session.Conn, err error) error) error {
	if loop.client.IsClosed() {
		return pnet.ErrClosedClient
	}
	select {
	case loop.writeChan <- writeEvent{buf: b, callback: callback}:
	default:
		return pnet.ErrWriteQueueFull
	}
	return nil
}

func (loop *eventLoop) asyncWritev(bs [][]byte, callback func(c session.Conn, err error) error) error {
	if loop.client.IsClosed() {
		return pnet.ErrClosedClient
	}
	select {
	case loop.writeChan <- writeEvent{bufS: bs, callback: callback}:
	default:
		return pnet.ErrWriteQueueFull
	}
	return nil
}

func (loop *eventLoop) flush() error {
	if loop.client.IsClosed() {
		return pnet.ErrClosedClient
	}
	select {
	case loop.writeChan <- writeEvent{}:
	default:
	}
	return nil
}

// start
//
//	@Description: 开启事件处理循环
//	@receiver loop
//	@param conn 关联的连接
//	@return error
func (loop *eventLoop) start(conn *Conn) error {
	if conn == nil {
		return errdef.ErrInvalidParams
	}
	loop.conn = conn
	loop.buffer = make([]byte, loopBufferSize)
	loop.readChan = make(chan int)
	loop.readDone = make(chan struct{})
	loop.writeChan = make(chan writeEvent, loop.client.WriteQueueCap)
	loop.closeChan = make(chan error)
	loop.cancelCtx, loop.cancelFunc = context.WithCancel(context.Background())
	go loop.run(true)
	go loop.readConn()
	return nil
}

// stop
//
//	@Description: 停止事件处理并关闭client
//	@receiver loop
//	@param ctx
//	@param reason 关闭原因
//	@return error
func (loop *eventLoop) stop(ctx context.Context, reason error) error {
	if loop.client.IsClosed() {
		return pnet.ErrClosedClient
	}
	select {
	case loop.closeChan <- reason:
	case <-ctx.Done():
	}
	return nil
}

// _stop
//
//	@Description: 在执行go
//	@receiver loop
//	@param reason
func (loop *eventLoop) _stop(reason error) {
	if !loop.client.toClosed(reason) {
		return
	}
	if err := loop.conn.TCPConn.Close(); err != nil {
		plog.Error("close connection error:", pfield.Error(err))
	}
	// 即便连接关闭出错，也继续走关闭逻辑
	loop.handler.OnClose(loop.conn, reason)
	loop.cancelFunc()
	loop.cache.Reset()
	loop.buffer = nil
}
