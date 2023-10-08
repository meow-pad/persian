package session

import "persian/frame/plog"

const (
	InvalidSessionId = 0
)

// Session
//
//	@Description: 连接会话
type Session interface {

	// Id
	//	@Description: 会话编号,仅关联后的会话有合法的id值，否则都返回 InvalidSessionId
	//	@return ID 非nil值
	//
	Id() uint64

	// Context
	//	@Description: 返回当前关联的业务对象
	//	@return Context
	//
	Context() Context

	// Register
	//  @Description: 注册会话
	//  @param context
	//	@return error
	//
	Register(context Context) error

	// setContext
	//	@Description: 设置关联的业务对象
	//	@param context
	//
	setContext(context Context)

	// Connection
	//	@Description: 连接
	//	@return Conn
	//
	Connection() Conn

	// Close
	//  @Description: 关闭会话
	//  @return error
	//
	Close() error

	// IsClosed
	//  @Description: 是否已关闭
	//  @return bool
	//
	IsClosed() bool

	// SendMessage
	//	@Description: 发送消息
	//	@param message
	//
	SendMessage(message any)

	// SendMessages
	//	@Description: 发送消息
	//	@param messages
	//
	SendMessages(messages ...any)
}

type BaseSession struct {
	// 关联的业务对象
	context Context
}

func (sess *BaseSession) Id() uint64 {
	if sess.context == nil {
		return InvalidSessionId
	}
	return sess.context.Id()
}

func (sess *BaseSession) Context() Context {
	return sess.context
}

func (sess *BaseSession) setContext(context Context) {
	if context != nil && sess.context != nil {
		plog.Error("set session context again?")
		return
	}
	sess.context = context
}

func (sess *BaseSession) Register(context Context) error {
	sess.setContext(context)
	return nil
}
