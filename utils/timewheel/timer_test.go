package timewheel

import (
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimer_TimeWheel_AddCron(t *testing.T) {
	should := require.New(t)
	tw, err := NewTimeWheel(1*time.Second, 4)
	should.Nil(err)
	tw.Start()
	count := 0
	wg := sync.WaitGroup{}
	wg.Add(1)
	startTime := time.Now()
	task := tw.AddCron(1*time.Second, func() {
		t.Logf("do func :time=%v\n", time.Now().Sub(startTime))
		count++
		if count >= 4 {
			wg.Done()
		}
	})
	t.Logf("task = %v\n", task)
	wg.Wait()
	err = tw.Remove(task)
	should.Nil(err)
	tw.Stop()
}

func TestTimer_TimeWheel_Add(t *testing.T) {
	should := require.New(t)
	tw, err := NewTimeWheel(1*time.Second, 4)
	should.Nil(err)
	tw.Start()
	count := atomic.Uint32{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	startTime := time.Now()
	var tFunc func()
	tFunc = func() {
		t.Logf("do func :time=%v\n", time.Now().Sub(startTime))
		count.Add(1)
		if count.Load() >= 4 {
			wg.Done()
		} else {
			tw.Add(1*time.Second, tFunc)
		}
	}
	task := tw.Add(1*time.Second, tFunc)
	t.Logf("task = %v\n", task)
	wg.Wait()
	err = tw.Remove(task)
	should.Nil(err)
	tw.Stop()
}
