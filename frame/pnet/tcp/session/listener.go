package session

// Listener
//
//	@Description: 会话监听器
type Listener interface {

	// OnOpened
	//	@Description: 当会话开启
	//	@param session
	//
	OnOpened(session Session)

	// OnClosed
	//	@Description: 当会话关闭
	//	@param session
	//
	OnClosed(session Session)

	// OnReceive
	//	@Description: 成功接收消息
	//	@param session
	//	@param msg 接收的消息体
	//	@param msgLen 消息长度
	//	@return err
	//
	OnReceive(session Session, msg any, msgLen int) (err error)

	// OnReceiveMulti
	//	@Description: 成功接收消息
	//	@param session
	//	@param msg 接收的消息体
	//	@param totalLen 消息长度
	//	@return err
	//
	OnReceiveMulti(session Session, msg []any, totalLen int) (err error)

	// OnSend
	//	@Description: 成功发送消息
	//	@param session
	//	@param msg 要发送的消息体
	//	@param msgLen 消息长度
	//	@return err
	//
	OnSend(session Session, msg any, msgLen int) (err error)

	// OnSendMulti
	//	@Description: 成功发送消息
	//	@param session
	//	@param msg 要发送的消息体
	//	@param totalLen 消息长度
	//	@return out 经过处理传出的消息体
	//	@return err
	//
	OnSendMulti(session Session, msg []any, totalLen int) (err error)
}

// EmptyListener 空处理的会话监听器
type EmptyListener struct {
}

func (listener *EmptyListener) OnOpened(session Session) {
}

func (listener *EmptyListener) OnClosed(session Session) {
}

func (listener *EmptyListener) OnReceive(session Session, msg any, msgLen int) (err error) {
	return
}

func (listener *EmptyListener) OnSend(session Session, msg any, msgLen int) (err error) {
	return
}

func (listener *EmptyListener) OnReceiveMulti(session Session, msg []any, totalLen int) (err error) {
	return
}

func (listener *EmptyListener) OnSendMulti(session Session, msg []any, totalLen int) (err error) {
	return
}
