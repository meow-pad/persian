package collections

import (
	"context"
	stdErrors "errors"
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/utils/loggers"
	"sync/atomic"
)

var (
	ErrClosedSafeChan          = stdErrors.New("SafeChannel has been closed")
	ErrEmptySafeChan           = stdErrors.New("SafeChannel is empty")
	ErrFullSafeChan            = stdErrors.New("SafeChannel is full")
	ErrDoneSafeChanContext     = stdErrors.New("SafeChannel context is done")
	ErrBrokenSafeChanListening = stdErrors.New("SafeChannel listening is broken")
)

// NewSafeChan
//
//	@Description: 构造新的安全队列
//	@param capacity int chan的容量
//	@return *SafeChannel[T]
func NewSafeChan[T any](capacity int) *SafeChannel[T] {
	channel := SafeChannel[T]{
		inner: make(chan T, capacity),
	}
	channel.cancelCtx, channel.cancelFunc = context.WithCancel(context.Background())
	return &channel
}

// SafeChannel
//
//	@Description: 一个安全 chan 的实现
type SafeChannel[T any] struct {
	closed     atomic.Bool
	inner      chan T
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
}

// Listen
//
//	@Description: 阻塞监听消息
//	@receiver channel *SafeChannel[T]
//	@param handler func(value T) bool 如果期望遇到异常也能继续处理后续消息，需要自行捕捉异常
//	@return error
func (channel *SafeChannel[T]) Listen(handler func(value T) bool) error {
	if handler == nil {
		return errdef.ErrInvalidParams
	}
	if channel.closed.Load() {
		return ErrClosedSafeChan
	}
	for {
		select {
		case value, ok := <-channel.inner:
			if !ok {
				loggers.Warn("SafeChannel inner chan is closed")
				return ErrClosedSafeChan
			}
			if !handler(value) {
				//return ErrBrokenSafeChanListening
				return nil
			}
		case <-channel.cancelCtx.Done():
			return ErrClosedSafeChan
		}
	} // end of for
}

// DirectGet
//
//	@Description: 获取队列头的数据，无论关闭与否
//	@receiver channel *SafeChannel[T]
//	@return value T
//	@return err error
func (channel *SafeChannel[T]) DirectGet() (value T, err error) {
	ok := false
	select {
	case value, ok = <-channel.inner:
		if !ok {
			err = ErrClosedSafeChan
			loggers.Warn("SafeChannel inner chan is closed")
		}
		return
	default:
		err = ErrEmptySafeChan
		return
	}
}

// Get
//
//	@Description: 获取队列头的数据
//	@receiver channel *SafeChannel[T]
//	@return value T
//	@return err error 如果channel已关闭，则返回 ErrClosedSafeChan；如果队列没有数据则即刻返回 ErrEmptySafeChan
func (channel *SafeChannel[T]) Get() (value T, err error) {
	if channel.closed.Load() {
		err = ErrClosedSafeChan
		return
	}
	return channel.DirectGet()
}

// BlockingGet
//
//	@Description: 获取队列头的的数据，如果队列为空则阻塞
//	@receiver channel *SafeChannel[T]
//	@param ctx context.Context
//	@return value T
//	@return err error 如果channel已关闭，则返回 ErrClosedSafeChan；上下文超时则返回 ErrDoneSafeChanContext
func (channel *SafeChannel[T]) BlockingGet(ctx context.Context) (value T, err error) {
	if channel.closed.Load() {
		err = ErrClosedSafeChan
		return
	}
	ok := false
	if ctx == nil {
		ctx = context.TODO()
	}
	select {
	case value, ok = <-channel.inner:
		if !ok {
			err = ErrClosedSafeChan
		}
	case <-ctx.Done():
		err = ErrDoneSafeChanContext
	case <-channel.cancelCtx.Done():
		err = ErrClosedSafeChan
	}
	return
}

// Put
//
//	@Description: 插入数据到队列
//	@receiver channel *SafeChannel[T]
//	@param value T
//	@return error 如果当前队列已满，则返回 ErrFullSafeChan
func (channel *SafeChannel[T]) Put(value T) error {
	if channel.closed.Load() {
		return ErrClosedSafeChan
	}
	select {
	case channel.inner <- value:
		return nil
	default:
		return ErrFullSafeChan
	}
}

// BlockingPut
//
//	@Description: 阻塞插入数据到队列
//	@receiver channel *SafeChannel[T]
//	@param ctx context.Context
//	@param value T
//	@return error 如果channel已关闭，则返回 ErrClosedSafeChan；上下文超时则返回 ErrDoneSafeChanContext
func (channel *SafeChannel[T]) BlockingPut(ctx context.Context, value T) error {
	if channel.closed.Load() {
		return ErrClosedSafeChan
	}
	if ctx == nil {
		ctx = context.TODO()
	}
	select {
	case channel.inner <- value:
		return nil
	case <-ctx.Done():
		return ErrDoneSafeChanContext
	case <-channel.cancelCtx.Done():
		return ErrClosedSafeChan
	}
}

// Close
//
//	@Description: 关闭队列
//	@receiver channel *SafeChannel[T]
func (channel *SafeChannel[T]) Close() {
	if !channel.closed.CompareAndSwap(false, true) {
		return
	}
	// 不关闭inner chan
	channel.cancelFunc()
}

// Capacity
//
//	@Description: 队列容量
//	@receiver channel *SafeChannel[T]
//	@return int
func (channel *SafeChannel[T]) Capacity() int {
	return cap(channel.inner)
}

// Length
//
//	@Description: 队列中数据量
//	@receiver channel *SafeChannel[T]
//	@return int
func (channel *SafeChannel[T]) Length() int {
	return len(channel.inner)
}

// IsClosed
//
//	@Description: 队列关闭状态
//	@receiver channel *SafeChannel[T]
//	@return bool
func (channel *SafeChannel[T]) IsClosed() bool {
	return channel.closed.Load()
}
