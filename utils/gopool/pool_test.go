package gopool

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"sync"
	"testing"
	"time"
)

func TestGoroutinePool_Stop(t *testing.T) {
	should := require.New(t)
	pool, err := NewGoroutinePool("testPool", 5, 100, true)
	should.Nil(err)
	err = pool.Start()
	should.Nil(err)
	count := atomic.Int32{}
	wg := sync.WaitGroup{}
	num := 10
	sleepTime := 1 * time.Second
	time1 := time.Now().UnixMilli()
	wg.Add(num)
	for i := 0; i < num; i++ {
		err = pool.Submit(func() {
			time.Sleep(sleepTime)
			count.Inc()
			wg.Done()
		})
		should.Nil(err)
	}
	wg.Wait()
	time2 := time.Now().UnixMilli()
	should.True((time2 - time1) >= int64(sleepTime/time.Millisecond))
	t.Logf("time diff:%vs", float32(time2-time1)/1000)
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	err = pool.Stop(ctx)
	should.Nil(err)
	should.True(count.Load() == int32(num))
}

func TestGoroutinePool_Submit(t *testing.T) {
	should := require.New(t)
	pool, err := NewGoroutinePool("testPool", 3000, 6000, true)
	should.Nil(err)
	err = pool.Start()
	should.Nil(err)
	// 参数
	loopNum := 10000
	count := atomic.Int32{}
	wg := sync.WaitGroup{}
	wg.Add(loopNum)
	sleepTime := 3 * time.Second
	//
	for i := 0; i < loopNum; i++ {
		err = pool.Submit(func() {
			time.Sleep(sleepTime)
			count.Inc()
			wg.Done()
		})
		should.Nil(err)
	}
	wg.Wait()
	err = pool.Stop(context.TODO())
	should.Nil(err)
	should.True(count.Load() == int32(loopNum))
}

func TestGoroutinePool_Panic(t *testing.T) {
	should := require.New(t)
	pool := &GoroutinePool{}
	pool, err := NewGoroutinePool("testPool", 5, 100, true)
	should.Nil(err)
	err = pool.Start()
	should.Nil(err)
	// 参数
	wg := sync.WaitGroup{}
	wg.Add(1)
	sleepTime := 1 * time.Second
	//
	err = pool.Submit(func() {
		defer wg.Done()
		time.Sleep(sleepTime)
		if err == nil {
			panic("panic func")
		}
	})
	should.Nil(err)
	wg.Wait()
	err = pool.Stop(context.TODO())
	should.Nil(err)
	time.Sleep(3 * time.Second)
}
