package client

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"persian/errdef"
	"persian/frame/plog"
	"persian/frame/plog/cfield"
	"persian/frame/pnet"
	"persian/utils/coding"
)

// 写事件
type writeEvent struct {
	buf      []byte
	bufS     [][]byte
	callback func(c *Conn, err error) error
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
	readChan   chan []byte
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
	defer coding.CachePanicError("run error:", func() {
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
		case buf := <-loop.readChan:
			act := loop.handler.OnTraffic(loop.conn, buf)
			if act == actionClose {
				loop._stop(nil)
				return
			}
			loop.readDone <- struct{}{}
		case event := <-loop.writeChan:
			var err error
			var write bool
			if len(event.buf) > 0 {
				write = true
				err = loop.conn.WriteMessage(websocket.BinaryMessage, event.buf)
			}
			if err == nil && len(event.bufS) > 0 {
				for _, buf := range event.bufS {
					if len(event.buf) > 0 {
						write = true
						err = loop.conn.WriteMessage(websocket.BinaryMessage, buf)
						if err != nil {
							break
						}
					}
				}
			}
			if write {
				if event.callback != nil {
					if cErr := event.callback(loop.conn, err); cErr != nil {
						plog.Error("write callback error:", cfield.Error(cErr))
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
	defer coding.CachePanicError("read conn error:", func() {
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
		opType, payload, err := loop.conn.ReadMessage()
		if err == nil {
			switch opType {
			case websocket.TextMessage, websocket.BinaryMessage:
				// 提交已读数据
				if len(payload) > 0 {
					// 等待读处理
					select {
					case loop.readChan <- payload:
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
			case websocket.CloseMessage:
			case websocket.PingMessage:
			case websocket.PongMessage:
			default:
				plog.Warn("unknown ws opCode:", cfield.Int("opCode", opType))
			}
		}
		if err != nil {
			if sErr := loop.stop(context.Background(), err); sErr != nil && !errors.Is(sErr, pnet.ErrClosedClient) {
				plog.Error("close conn error:", cfield.Error(sErr))
			}
			return
		}
	} // end of for
}

func (loop *eventLoop) asyncWrite(b []byte, callback func(c *Conn, err error) error) error {
	if loop.client.IsClosed() {
		return pnet.ErrClosedClient
	}
	if len(b) <= 0 {
		return nil
	}
	select {
	case loop.writeChan <- writeEvent{buf: b, callback: callback}:
	default:
		return pnet.ErrWriteQueueFull
	}
	return nil
}

func (loop *eventLoop) asyncWritev(bs [][]byte, callback func(c *Conn, err error) error) error {
	if loop.client.IsClosed() {
		return pnet.ErrClosedClient
	}
	if len(bs) <= 0 {
		return nil
	}
	select {
	case loop.writeChan <- writeEvent{bufS: bs, callback: callback}:
	default:
		return pnet.ErrWriteQueueFull
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
	loop.readChan = make(chan []byte)
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
	if err := loop.conn.Close(); err != nil {
		plog.Error("close connection error:", cfield.Error(err))
	}
	// 即便连接关闭出错，也继续走关闭逻辑
	loop.handler.OnClose(loop.conn, reason)
	loop.cancelFunc()
}
