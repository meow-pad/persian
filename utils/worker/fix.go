package worker

import (
	"context"
	"errors"
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/utils/coding"
	"github.com/meow-pad/persian/utils/runtime"
	"go.uber.org/atomic"
	"math"
	"sync"
)

func NewFixedIoWorkerPool(ioEventTimeRatio float64, ioEventQueueSize int, blockingOnFull bool) (*FixedWorkerPool, error) {
	ioWorkerNum := (1 + int(math.Round(ioEventTimeRatio))) * runtime.MaxProcess()
	if ioWorker, err := NewFixedWorkerPool(ioWorkerNum, ioEventQueueSize, blockingOnFull); err != nil {
		return nil, err
	} else {
		return ioWorker, nil
	}
}

func NewFixedWorkerPool(slotNum, queueSize int, blockingOnFull bool) (*FixedWorkerPool, error) {
	pool := &FixedWorkerPool{}
	if err := pool.init(slotNum, queueSize, blockingOnFull); err != nil {
		return nil, err
	}
	return pool, nil
}

func newFixedPoolWorker(pool *FixedWorkerPool, allocateQueueSize int) *fixedPoolWorker {
	return &fixedPoolWorker{
		pool:  pool,
		queue: make(chan func(*GoroutineLocal), allocateQueueSize),
	}
}

// fixedPoolWorker
//
//	@Description: 固定工作器池的工作器
type fixedPoolWorker struct {
	pool  *FixedWorkerPool
	queue chan func(*GoroutineLocal)
	local GoroutineLocal
}

// run
//
//	@Description: 自修复的函数
//	@receiver pool *fixedPoolWorker
func (worker *fixedPoolWorker) run() {
	defer coding.CachePanicError("fixedPoolWorker run task error:", func() {
		if worker.pool.closed.Load() {
			worker.pool.waitGroup.Done()
			return
		}
		go worker.run()
	})
	for {
		select {
		case task := <-worker.queue:
			task(&worker.local)
		case <-worker.pool.closeCtx.Done():
			worker.pool.waitGroup.Done()
			return
		}
	} // end of for
}

// FixedWorkerPool
//
//	@Description: 固定工作器池
type FixedWorkerPool struct {
	taskWorkers []*fixedPoolWorker
	closeCtx    context.Context
	closeFunc   context.CancelFunc
	closed      atomic.Bool
	waitGroup   sync.WaitGroup
	// 配置
	slotNum        int
	queueSize      int
	blockingOnFull bool
}

func (pool *FixedWorkerPool) init(slotNum, queueSize int, blockingOnFull bool) error {
	if pool.taskWorkers != nil {
		return errors.New("cant repeatedly init FixedWorkerPool")
	}
	if slotNum <= 0 || queueSize <= 0 {
		return errdef.ErrInvalidParams
	}
	pool.closed.Store(true)
	pool.slotNum = slotNum
	pool.blockingOnFull = blockingOnFull
	pool.queueSize = queueSize
	allocateQueueSize := queueSize
	pool.taskWorkers = make([]*fixedPoolWorker, slotNum)
	pool.closeCtx, pool.closeFunc = context.WithCancel(context.Background())
	pool.waitGroup.Add(pool.slotNum)
	for i := 0; i < pool.slotNum; i++ {
		worker := newFixedPoolWorker(pool, allocateQueueSize)
		pool.taskWorkers[i] = worker
		go worker.run()
	}
	pool.closed.Store(false)
	return nil
}

func (pool *FixedWorkerPool) Submit(group int, task func(*GoroutineLocal)) error {
	if pool.closed.Load() {
		return ErrWorkerPoolClosed
	}
	index := 0
	if pool.slotNum > 1 {
		index = group % pool.slotNum
	}
	worker := pool.taskWorkers[index]
	if pool.blockingOnFull {
		select {
		case worker.queue <- task:
			return nil
		case <-pool.closeCtx.Done():
			return ErrWorkerPoolClosed
		}
	} else {
		select {
		case worker.queue <- task:
			return nil
		default:
			return ErrWorkerPoolQueueIsFull
		}
	} // end of else
}

func (pool *FixedWorkerPool) SubmitToAll(task func(*GoroutineLocal), blockingOnFull bool) error {
	if pool.closed.Load() {
		return ErrWorkerPoolClosed
	}
	for _, worker := range pool.taskWorkers {
		if blockingOnFull {
			select {
			case worker.queue <- task:
			case <-pool.closeCtx.Done():
				return ErrWorkerPoolClosed
			}
		} else {
			select {
			case worker.queue <- task:
			default:
			}
		} // end of else
	}
	return nil
}

func (pool *FixedWorkerPool) Shutdown(ctx context.Context) error {
	if !pool.closed.CompareAndSwap(false, true) {
		return ErrWorkerPoolClosed
	}
	// 关闭
	pool.closeFunc()
	done := make(chan struct{})
	go func() {
		pool.waitGroup.Wait()
		close(done)
	}()
	// 等待所有worker正确退出
	select {
	case <-ctx.Done():
		plog.Warn("abort FixedWorkerPool by context")
	case <-done:
	}
	return nil
}

func (pool *FixedWorkerPool) SlotNum() int {
	return pool.slotNum
}

func (pool *FixedWorkerPool) QueueSize() int {
	return pool.queueSize
}
