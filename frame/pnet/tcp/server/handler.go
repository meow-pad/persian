package server

import (
	"github.com/panjf2000/gnet/v2"
	"persian/frame/plog"
	"persian/frame/plog/cfield"
	"reflect"
	"time"
)

type eventHandler struct {
	server *Server
}

func (handler *eventHandler) OnBoot(engine gnet.Engine) (action gnet.Action) {
	handler.server.startChan <- engine
	return
}

func (handler *eventHandler) OnShutdown(_ gnet.Engine) {
}

func (handler *eventHandler) OnOpen(gConn gnet.Conn) (out []byte, action gnet.Action) {
	conn, err := newConn(gConn)
	if err != nil {
		plog.Error("create connection error:", cfield.Error(err))
		action = gnet.Close
		return
	}
	sess, err := newSession(handler.server, conn)
	if err != nil {
		plog.Error("create session error:", cfield.Error(err))
		action = gnet.Close
		return
	}
	plog.Debug("open connecting:",
		cfield.String("server", handler.server.name),
		cfield.Uint64("conn", conn.Hash()))
	if err = handler.server.AddSession(sess); err != nil {
		plog.Error("add session error:", cfield.Error(err))
		action = gnet.Close
		return
	}
	handler.server.listener.OnOpened(sess)
	return
}

func (handler *eventHandler) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	connCtx := conn.Context()
	if connCtx == nil {
		plog.Warn("connection less context", cfield.NamedError("closeReason", err))
		return
	}
	sess, ok := connCtx.(*svrSession)
	if !ok {
		plog.Error("connection context is not session",
			cfield.String("contextType", reflect.TypeOf(connCtx).String()),
			cfield.NamedError("closeReason", err))
		return
	}
	// 转换到关闭状态
	sess.conn.ToClosed(err)
	// 移除会话
	handler.server.RemoveSession(sess)
	plog.Debug("close connecting:",
		cfield.String("server", handler.server.name),
		cfield.Uint64("conn", sess.conn.Hash()))
	// 触发关闭监听
	handler.server.listener.OnClosed(sess)
	return
}

func (handler *eventHandler) OnTraffic(conn gnet.Conn) (action gnet.Action) {
	connCtx := conn.Context()
	if connCtx == nil {
		plog.Warn("connection less context")
		return
	}
	sess, ok := connCtx.(*svrSession)
	if !ok {
		plog.Error("connection context is not session",
			cfield.String("contextType", reflect.TypeOf(connCtx).String()))
		action = gnet.Close
		return
	}
	msgArr, totalLen, err := handler.server.codec.Decode(conn)
	if err != nil {
		plog.Error("decode error", cfield.Error(err))
		action = gnet.Close
		return
	}
	msgNum := len(msgArr)
	if msgNum > 1 {
		_ = handler.server.listener.OnReceiveMulti(sess, msgArr, totalLen)
	} else if msgNum == 1 {
		_ = handler.server.listener.OnReceive(sess, msgArr[0], totalLen)
	}
	return
}

func (handler *eventHandler) OnTick() (delay time.Duration, action gnet.Action) {
	handler.server.CheckSessions()
	delay = handler.server.options.CheckSessionInterval
	return
}
