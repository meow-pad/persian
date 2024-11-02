package gopool

import (
	"context"
	"errors"
	"github.com/meow-pad/persian/utils/coding"
	"github.com/meow-pad/persian/utils/collections"
	"github.com/meow-pad/persian/utils/loggers"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"time"
)

const (
	workerExpiryDuration = 15 * time.Second
)

var (
	ErrInvalidMaxWorkerNum     = errors.New("invalid max worker num")
	ErrInvalidTaskQueueSize    = errors.New("invalid task queue size")
	ErrStartGoroutinePoolFirst = errors.New("please start GoroutinePool first")
	ErrGoroutinePoolClosed     = errors.New("GoroutinePool has been closed")
	ErrFullTaskQueue           = errors.New("task queue is full")
)

// NewGoroutinePool
//
//	@Description: 构建协程池
//		协程池最大同时可承受任务数为 maxWorkerNum + taskQueueSize
//	@param name string 命名
//	@param maxWorkerNum int 最大支持同时运行协程数（小于等于0时无限制）
//	@param taskQueueSize int 最大任务队列数，也是最大可并发提交任务数（非负数）
//	@param blockingOrNot bool 队列已满时是否进行阻塞，或直接返回错误
//	@return error
func NewGoroutinePool(name string, maxWorkerNum int, taskQueueSize int, blockingOrNot bool) (*GoroutinePool, error) {
	goPool := &GoroutinePool{}
	if err := goPool.init(name, maxWorkerNum, taskQueueSize, blockingOrNot); err != nil {
		return nil, err
	}
	return goPool, nil
}

// GoroutinePool
//
//	@Description: 协程池
//		为了使运行协程数与并发提交数可控而设计的协程池
//		该协程池能安全关闭，在停止时优先尝试处理完队列中的任务
type GoroutinePool struct {
	core *ants.Pool

	// 命名
	name string
	// 最大worker数量
	maxWorkerNum int
	// 缓冲队列长度
	taskQueueSize int
	// 队列满时阻塞标志
	blockingOrNot bool
	// 任务队列
	taskQueue *collections.SafeChannel[func()]
}

func (pool *GoroutinePool) init(name string, maxWorkerNum int, taskQueueSize int, blockingOrNot bool) error {
	if maxWorkerNum < 0 {
		return ErrInvalidMaxWorkerNum
	}
	if taskQueueSize < 0 {
		return ErrInvalidTaskQueueSize
	}
	pool.name = name
	pool.maxWorkerNum = maxWorkerNum
	pool.taskQueueSize = taskQueueSize
	pool.blockingOrNot = blockingOrNot
	return nil
}

func (pool *GoroutinePool) Start() error {
	opts := []ants.Option{
		ants.WithExpiryDuration(workerExpiryDuration),
		ants.WithPreAlloc(false),
		ants.WithNonblocking(false),
		// 这里实际只有1就可以了，毕竟只有一个提交协程
		ants.WithMaxBlockingTasks(10_000),
		ants.WithPanicHandler(pool.onPanic),
		ants.WithLogger(loggers.GetStdLogger()),
	}
	if poolCore, err := ants.NewPool(pool.maxWorkerNum, opts...); err != nil {
		return err
	} else {
		pool.taskQueue = collections.NewSafeChan[func()](pool.taskQueueSize)
		pool.core = poolCore
		go pool.consumeQueue()
		return nil
	}
}

func (pool *GoroutinePool) Stop(ctx context.Context) error {
	if pool.IsStopped() {
		return ErrGoroutinePoolClosed
	}

	if !pool.core.IsClosed() {
		released := make(chan struct{})
		go func() {
			// 关闭协程池
			err := pool.core.ReleaseTimeout(30 * time.Minute) // 等待足够长的时间
			if err != nil {
				loggers.Error("release pool error:", zap.Error(err))
			}
			released <- struct{}{}
		}()
		select {
		case <-ctx.Done():
			loggers.Warn("cancel GoroutinePool by context", zap.String("poolName", pool.name))
			break
		case <-released:
			// 释放完毕
			break
		} // end of select
	}
	// 关闭任务队列
	pool.taskQueue.Close()
	if pool.core.IsClosed() {
		return nil
	} else {
		return errors.New("pool release failure")
	}
}

// consumeQueue
//
//	@Description: 消耗任务队列
//	@receiver pool *GoroutinePool
func (pool *GoroutinePool) consumeQueue() {
	defer coding.CatchPanicError("consume GoroutinePool's queue error:", func() {
		if pool.IsStopped() {
			return
		}
		go pool.consumeQueue()
	})
	handler := func(task func()) bool {
		if err := pool.core.Submit(task); err != nil {
			switch {
			case errors.Is(err, ants.ErrPoolClosed):
				// 池子已关闭则退出
				return false
			case errors.Is(err, ants.ErrPoolOverload):
				// 池子已满？以目前配置的阻塞方式，不太可能触发
				loggers.Error("GoroutinePool is overloaded", zap.String("poolName", pool.name))
			default:
				loggers.Error("submitting task to GoroutinePool's core error:",
					zap.String("poolName", pool.name), zap.Error(err))
			}
		} // end of if
		return true
	}
	err := pool.taskQueue.Listen(handler)
	if !errors.Is(err, collections.ErrClosedSafeChan) {
		// 被关闭则尝试消费完剩余的任务
		for {
			task, dErr := pool.taskQueue.DirectGet()
			if dErr != nil {
				if errors.Is(dErr, collections.ErrClosedSafeChan) || errors.Is(dErr, collections.ErrEmptySafeChan) {
					return
				}
				loggers.Error("get task error:", zap.Error(dErr))
			} else {
				handler(task)
			}
		} // end of for
	} // end of if
}

// onPanic 未处理异常捕捉
func (pool *GoroutinePool) onPanic(obj any) {
	if err := coding.Cast[error](obj); err != nil {
		loggers.Error("panic error in GoroutinePool:", zap.String("poolName", pool.name), zap.Error(err), zap.StackSkip("stack", 2))
	} else {
		loggers.Error("panic event in GoroutinePool:", zap.String("poolName", pool.name), zap.Any("panicObj", obj), zap.StackSkip("stack", 2))
	}
}

// Submit
//
//	@Description: 提交任务到协程池
//		当任务队列中积压的消息(处理不过来)超过队列长度,该方法将阻塞直到队列有空位
//	@receiver pool *GoroutinePool
//	@param task func()
//	@return error 如果池子已经关闭,则会返回 ErrGoroutinePoolClosed
func (pool *GoroutinePool) Submit(task func()) error {
	if pool.core == nil || pool.taskQueue == nil {
		return ErrStartGoroutinePoolFirst
	}
	if pool.IsStopped() {
		return ErrGoroutinePoolClosed
	}
	if pool.blockingOrNot {
		err := pool.taskQueue.BlockingPut(nil, task)
		if errors.Is(err, collections.ErrClosedSafeChan) {
			return ErrGoroutinePoolClosed
		}
		return err
	} else {
		err := pool.taskQueue.Put(task)
		switch {
		case errors.Is(err, collections.ErrClosedSafeChan):
			return ErrGoroutinePoolClosed
		case errors.Is(err, collections.ErrFullSafeChan):
			return ErrFullTaskQueue
		default:
			return err
		}
	}
}

// AvailableInQueue
//
//	@Description: 队列中可用空间数
//	@receiver pool *GoroutinePool
//	@return int
func (pool *GoroutinePool) AvailableInQueue() int {
	return pool.taskQueue.Capacity() - pool.taskQueue.Length()
}

// WaitingTask
//
//	@Description: 队列中等待的任务数
//	@receiver pool *GoroutinePool
//	@return int
func (pool *GoroutinePool) WaitingTask() int {
	return pool.taskQueue.Length()
}

// QueueCapacity
//
//	@Description: 队列容量
//	@receiver pool *GoroutinePool
//	@return int
func (pool *GoroutinePool) QueueCapacity() int {
	return pool.taskQueue.Capacity()
}

// WorkerNum
//
//	@Description: 运行的工作协程数量(而非工作中协程数量)
//	@receiver pool *GoroutinePool
func (pool *GoroutinePool) WorkerNum() int {
	return pool.core.Running()
}

// MaxWorkerNum
//
//	@Description: 最大工作协程数
//	@receiver pool *GoroutinePool
func (pool *GoroutinePool) MaxWorkerNum() int {
	return pool.core.Cap()
}

// IsStopped
//
//	@Description: 判定池是否以关闭
//	@receiver pool *GoroutinePool
//	@return bool
func (pool *GoroutinePool) IsStopped() bool {
	return pool.taskQueue.IsClosed() || pool.core.IsClosed()
}
