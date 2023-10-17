package server

import (
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/frame/pnet"
	"github.com/meow-pad/persian/frame/pnet/tcp/session"
)

// newSession
//
//	@Description: 构造session
//	@param server
//	@param conn
//	@return session.Session
//	@return error
func newSession(server *Server, conn *Conn) (*svrSession, error) {
	if server == nil || conn == nil {
		return nil, errdef.ErrInvalidParams
	}
	if closed, _ := conn.IsClosed(); closed {
		return nil, pnet.ErrClosedConn
	}
	svrSess := &svrSession{
		server: server,
		conn:   conn,
	}
	conn.SetContext(svrSess)
	return svrSess, nil
}

type svrSession struct {
	session.BaseSession

	// 关联的服务
	server *Server
	// 关联的连接
	conn *Conn
}

func (sess *svrSession) Connection() session.Conn {
	return sess.conn
}

func (sess *svrSession) Register(context session.Context) error {
	return sess.server.RegisterSession(sess, context)
}

func (sess *svrSession) Close() error {
	return sess.conn.Close()
}

func (sess *svrSession) IsClosed() bool {
	closed, _ := sess.conn.IsClosed()
	return closed
}

func (sess *svrSession) SendMessage(message any) {
	if closed, _ := sess.conn.IsClosed(); closed {
		plog.Debug("cant send to closed conn")
		return
	}
	data, err := sess.server.codec.Encode(message)
	if err != nil {
		sess.onSendingError("encode message error:", err)
		return
	}
	dataLen := len(data)
	err = sess.conn.AsyncWrite(data, func(c session.Conn, err error) error {
		if err != nil {
			sess.onSendingError("write message error:", err)
			return nil
		}
		if err = sess.server.listener.OnSend(sess, message, dataLen); err != nil {
			plog.Error("on send error:", pfield.Error(err))
		}
		return nil
	})
	if err != nil {
		sess.onSendingError("async write error:", err)
	}
}

func (sess *svrSession) SendMessages(messages ...any) {
	if closed, _ := sess.conn.IsClosed(); closed {
		plog.Debug("cant send to closed conn")
		return
	}
	totalLen := 0
	dataArr := make([][]byte, 0, len(messages))
	for _, message := range messages {
		data, err := sess.server.codec.Encode(message)
		if err != nil {
			sess.onSendingError("encode message error:", err)
			return
		}
		dataArr = append(dataArr, data)
		totalLen += len(data)
	}
	err := sess.conn.AsyncWritev(dataArr, func(c session.Conn, err error) error {
		if err != nil {
			sess.onSendingError("write messages error:", err)
			return nil
		}
		if err = sess.server.listener.OnSendMulti(sess, messages, totalLen); err != nil {
			plog.Error("on send error:", pfield.Error(err))
		}
		return nil
	})
	if err != nil {
		sess.onSendingError("async writev error:", err)
	}
}

// onSendingError
//
//	@Description: 发送消息时错误处理
//	@receiver sess
//	@param tip 日志消息
//	@param err 错误
func (sess *svrSession) onSendingError(tip string, err error) {
	plog.Error(tip, pfield.Error(err))
	// 无法处理的状态，关闭连接
	cErr := sess.conn.Close()
	if cErr != nil {
		plog.Error("close conn error", pfield.Error(cErr))
	}
}
