package timewheel

import (
	"context"
	"errors"
	"github.com/meow-pad/persian/utils/gopool"
	"github.com/meow-pad/persian/utils/loggers"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	typeTimer taskType = iota
	typeTicker

	modeIsCircle  = true
	modeNotCircle = false

	modeIsAsync  = true
	modeNotAsync = false
)

type taskType int64
type taskID int64

type Task struct {
	delay    time.Duration
	id       taskID
	round    int
	callback func()

	async  bool
	stop   bool
	circle bool
	// circleNum int
}

// Reset for sync.Pool
func (t *Task) Reset() {
	t.delay = 0
	t.id = 0
	t.round = 0
	t.callback = nil

	t.async = false
	t.stop = false
	t.circle = false
}

type OptionCall func(*TimeWheel) error

// SetGoPool
//
//	@Description: 设置执行池
//	@param goPool *gopool.GoroutinePool
//	@return OptionCall
func SetGoPool(goPool *gopool.GoroutinePool) OptionCall {
	return func(o *TimeWheel) error {
		o.goPool = goPool
		return nil
	}
}

// NewTimeWheel create new time wheel
func NewTimeWheel(tick time.Duration, bucketsNum int, options ...OptionCall) (*TimeWheel, error) {
	tw := &TimeWheel{}
	if err := tw.init(tick, bucketsNum, options...); err != nil {
		return nil, err
	}
	return tw, nil
}

// TimeWheel
//
//		@Description: 时间轮 的实现
//	 	代码参考自 https://github.com/rfyiamcool/go-timewheel
type TimeWheel struct {
	randomID atomic.Int64

	tick   time.Duration
	ticker *time.Ticker

	bucketsNum    int
	buckets       []map[taskID]*Task // key: added item, value: *Task
	bucketIndexes map[taskID]int     // key: added item, value: bucket position

	goPool *gopool.GoroutinePool

	currentIndex int

	onceStart sync.Once

	stopC chan struct{}

	exited atomic.Bool

	sync.RWMutex
}

func (tw *TimeWheel) init(tick time.Duration, bucketsNum int, options ...OptionCall) error {
	if tick.Milliseconds() < 1 {
		return errors.New("invalid params, must tick >= 1 ms")
	}
	if bucketsNum <= 0 {
		return errors.New("invalid params, must bucketsNum > 0")
	}
	// init properties
	{
		// tick
		tw.tick = tick

		// store
		tw.bucketsNum = bucketsNum
		tw.bucketIndexes = make(map[taskID]int, 1024*100)
		tw.buckets = make([]map[taskID]*Task, bucketsNum)
		tw.currentIndex = 0

		// signal
		tw.stopC = make(chan struct{})
	}

	for i := 0; i < bucketsNum; i++ {
		tw.buckets[i] = make(map[taskID]*Task, 16)
	}

	for _, op := range options {
		if err := op(tw); err != nil {
			loggers.Error("select option error:", zap.Error(err))
		}
	}
	return nil
}

// Start to start the time wheel
func (tw *TimeWheel) Start() {
	// only once start
	tw.onceStart.Do(
		func() {
			tw.ticker = time.NewTicker(tw.tick)
			go tw.scheduler()
		},
	)
}

func (tw *TimeWheel) scheduler() {
	queue := tw.ticker.C

	for {
		select {
		case <-queue:
			tw.handleTick()

		case <-tw.stopC:
			tw.ticker.Stop()
			return
		}
	}
}

// Stop to stop the time wheel
func (tw *TimeWheel) Stop() {
	if !tw.exited.CompareAndSwap(false, true) {
		// 已经关闭
		return
	}
	tw.stopC <- struct{}{}
}

func (tw *TimeWheel) collectTask(task *Task) {
	index := tw.bucketIndexes[task.id]
	delete(tw.bucketIndexes, task.id)
	delete(tw.buckets[index], task.id)
}

func (tw *TimeWheel) handleTick() {
	tw.Lock()
	defer tw.Unlock()

	bucket := tw.buckets[tw.currentIndex]
	for k, task := range bucket {
		if task.stop {
			tw.collectTask(task)
			continue
		}

		if bucket[k].round > 0 {
			bucket[k].round--
			continue
		}

		if task.async {
			if tw.goPool != nil {
				if err := tw.goPool.Submit(task.callback); err != nil {
					loggers.Error("submit TimeWheel task error:", zap.Error(err))
				}
			} else {
				go task.callback()
			}
		} else {
			// optimize gopool
			task.callback()
		}

		// circle
		if task.circle == true {
			tw.collectTask(task)
			tw.putCircle(task, modeIsCircle)
			continue
		}

		// gc
		tw.collectTask(task)
	}

	if tw.currentIndex == tw.bucketsNum-1 {
		tw.currentIndex = 0
		return
	}

	tw.currentIndex++
}

// Add to add a task
func (tw *TimeWheel) Add(delay time.Duration, callback func()) *Task {
	return tw.addAny(delay, callback, modeNotCircle, modeIsAsync)
}

// AddCron add interval task
func (tw *TimeWheel) AddCron(delay time.Duration, callback func()) *Task {
	return tw.addAny(delay, callback, modeIsCircle, modeIsAsync)
}

func (tw *TimeWheel) addAny(delay time.Duration, callback func(), circle, async bool) *Task {
	if delay <= 0 {
		delay = tw.tick
	}

	id := tw.genUniqueID()
	task := new(Task)

	task.delay = delay
	task.id = id
	task.callback = callback
	task.circle = circle
	task.async = async // refer to src/runtime/time.go

	tw.put(task)
	return task
}

func (tw *TimeWheel) put(task *Task) {
	tw.Lock()
	defer tw.Unlock()

	tw.store(task, false)
}

func (tw *TimeWheel) putCircle(task *Task, circleMode bool) {
	tw.store(task, circleMode)
}

func (tw *TimeWheel) store(task *Task, circleMode bool) {
	round := tw.calculateRound(task.delay)
	index := tw.calculateIndex(task.delay)

	if round > 0 && circleMode {
		task.round = round - 1
	} else {
		task.round = round
	}

	tw.bucketIndexes[task.id] = index
	tw.buckets[index][task.id] = task
}

func (tw *TimeWheel) calculateRound(delay time.Duration) (round int) {
	delaySeconds := delay.Seconds()
	tickSeconds := tw.tick.Seconds()
	round = int(delaySeconds / tickSeconds / float64(tw.bucketsNum))
	return
}

func (tw *TimeWheel) calculateIndex(delay time.Duration) (index int) {
	delaySeconds := delay.Seconds()
	tickSeconds := tw.tick.Seconds()
	index = (int(float64(tw.currentIndex) + delaySeconds/tickSeconds)) % tw.bucketsNum
	return
}

func (tw *TimeWheel) Remove(task *Task) error {
	// tw.removeC <- task
	tw.remove(task)
	return nil
}

func (tw *TimeWheel) remove(task *Task) {
	tw.Lock()
	defer tw.Unlock()

	tw.collectTask(task)
}

func (tw *TimeWheel) NewTimer(delay time.Duration) *Timer {
	queue := make(chan bool, 1) // buf = 1, refer to src/time/sleep.go
	task := tw.addAny(delay,
		func() {
			notifyChannel(queue)
		},
		modeNotCircle,
		modeNotAsync,
	)

	// init timer
	ctx, cancel := context.WithCancel(context.Background())
	timer := &Timer{
		tw:     tw,
		C:      queue, // faster
		task:   task,
		Ctx:    ctx,
		cancel: cancel,
	}

	return timer
}

func (tw *TimeWheel) AfterFunc(delay time.Duration, callback func()) *Timer {
	queue := make(chan bool, 1)
	task := tw.addAny(delay,
		func() {
			callback()
			notifyChannel(queue)
		},
		modeNotCircle, modeIsAsync,
	)

	// init timer
	ctx, cancel := context.WithCancel(context.Background())
	timer := &Timer{
		tw:     tw,
		C:      queue, // faster
		task:   task,
		Ctx:    ctx,
		cancel: cancel,
		fn:     callback,
	}

	return timer
}

func (tw *TimeWheel) NewTicker(delay time.Duration) *Ticker {
	queue := make(chan bool, 1)
	task := tw.addAny(delay,
		func() {
			notifyChannel(queue)
		},
		modeIsCircle,
		modeNotAsync,
	)

	// init ticker
	ctx, cancel := context.WithCancel(context.Background())
	ticker := &Ticker{
		task:   task,
		tw:     tw,
		C:      queue,
		Ctx:    ctx,
		cancel: cancel,
	}

	return ticker
}

func (tw *TimeWheel) After(delay time.Duration) <-chan time.Time {
	queue := make(chan time.Time, 1)
	tw.addAny(delay,
		func() {
			queue <- time.Now()
		},
		modeNotCircle, modeNotAsync,
	)
	return queue
}

func (tw *TimeWheel) Sleep(delay time.Duration) {
	queue := make(chan bool, 1)
	tw.addAny(delay,
		func() {
			queue <- true
		},
		modeNotCircle, modeNotAsync,
	)
	<-queue
}

// Timer similar to golang std timer
type Timer struct {
	task   *Task
	tw     *TimeWheel
	fn     func() // external custom func
	stopFn func() // call function when timer stop

	C chan bool

	cancel context.CancelFunc
	Ctx    context.Context
}

func (t *Timer) Reset(delay time.Duration) {
	// first stop old task
	t.task.stop = true

	// make new task
	var task *Task
	if t.fn != nil { // use AfterFunc
		task = t.tw.addAny(delay,
			func() {
				t.fn()
				notifyChannel(t.C)
			},
			modeNotCircle, modeIsAsync, // must async mode
		)
	} else {
		task = t.tw.addAny(delay,
			func() {
				notifyChannel(t.C)
			},
			modeNotCircle, modeNotAsync)
	}

	t.task = task
}

func (t *Timer) Stop() {
	if t.stopFn != nil {
		t.stopFn()
	}

	t.task.stop = true
	t.cancel()
	if err := t.tw.Remove(t.task); err != nil {
		loggers.Error("remove timer task error:", zap.Error(err))
	}
}

func (t *Timer) AddStopFunc(callback func()) {
	t.stopFn = callback
}

type Ticker struct {
	tw     *TimeWheel
	task   *Task
	cancel context.CancelFunc

	C   chan bool
	Ctx context.Context
}

func (t *Ticker) Stop() {
	t.task.stop = true
	t.cancel()
	if err := t.tw.Remove(t.task); err != nil {
		loggers.Error("remove ticker task error:", zap.Error(err))
	}
}

func notifyChannel(q chan bool) {
	select {
	case q <- true:
	default:
	}
}

func (tw *TimeWheel) genUniqueID() taskID {
	id := tw.randomID.Inc()
	return taskID(id)
}
