package worker

import (
	"context"
	"github.com/meow-pad/persian/utils/coding"
	"github.com/pkg/errors"
	"sync/atomic"
)

var (
	ErrWorkerPoolQueueIsFull = errors.New("queue of worker pool is full")
	ErrWorkerPoolClosed      = errors.New("worker pool is closed")
)

//// IdentifiedTask
////	@Description: 有标识任务
////
//type IdentifiedTask interface {
//	Id() int // 标识
//	run()    // 可执行任务
//}

// GoroutineLocal
//
//	@Description: Goroutine 本地对象
type GoroutineLocal struct {
	localMap map[any]any
}

func (local *GoroutineLocal) Get(key any) (any, bool) {
	if local.localMap == nil {
		return nil, false
	}
	value, ok := local.localMap[key]
	return value, ok
}

func (local *GoroutineLocal) Set(key any, value any) {
	if local.localMap == nil {
		local.localMap = make(map[any]any)
	}
	local.localMap[key] = value
}

func (local *GoroutineLocal) Remove(key any) {
	if local.localMap == nil {
		return
	}
	delete(local.localMap, key)
}

func (local *GoroutineLocal) Range(op func(key, val any) bool) {
	if local.localMap == nil {
		return
	}
	for key, val := range local.localMap {
		if bBreak := op(key, val); bBreak {
			return
		}
	}
}

// Pool
//
//	@Description: 执行工作池
type Pool interface {

	// Submit
	//	@Description: 提交任务到工作池
	//	@param taskCategory 任务分类
	//	@param task 任务
	//	@return error 提交失败时返回错误
	//		如果工作队列已满,将返回 ErrWorkerPoolQueueIsFull 错误;
	//		如果工作池已经被关闭,则将返回 ErrWorkerPoolClosed 错误;
	//	@throwable 提交任务时可能会抛出异常,在必要时自行捕捉
	//
	Submit(category int, task func(*GoroutineLocal)) error

	// Shutdown
	//	@Description: 关闭工作池
	//	@param ctx
	//	@return error 关闭失败时返回错误
	//		如果工作池已经关闭或在关闭中,将返回 ErrWorkerPoolClosed
	//
	Shutdown(ctx context.Context) error
}

func NewSimpleWorkerPool() *SimpleWorkerPool {
	return &SimpleWorkerPool{}
}

// SimpleWorkerPool
//
//	@Description: 提交即执行的工作池
type SimpleWorkerPool struct {
	closed atomic.Bool
	local  GoroutineLocal
}

func (pool *SimpleWorkerPool) Submit(category int, task func(*GoroutineLocal)) error {
	if pool.closed.Load() {
		return ErrWorkerPoolClosed
	}
	defer coding.CachePanicError("SimpleWorkerPool run task error:", nil)
	task(&pool.local)
	return nil
}

func (pool *SimpleWorkerPool) Shutdown(ctx context.Context) error {
	if !pool.closed.CompareAndSwap(false, true) {
		return ErrWorkerPoolClosed
	}
	return nil
}
