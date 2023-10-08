package session

// Context
//
//	@Description: 会话上下文
type Context interface {
	Id() uint64      // 编号
	Deadline() int64 // 存活时间
}

type BaseContext struct {
	id       uint64
	deadline int64
}

func (ctx *BaseContext) Init(id uint64) {
	ctx.id = id
}

func (ctx *BaseContext) Id() uint64 {
	return ctx.id
}

func (ctx *BaseContext) Deadline() int64 {
	return ctx.deadline
}

func (ctx *BaseContext) SetDeadline(deadline int64) {
	ctx.deadline = deadline
}
