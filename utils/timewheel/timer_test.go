package timewheel

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestTimer_TimeWheel_Round(t *testing.T) {
	should := require.New(t)
	tw, err := NewTimeWheel(1*time.Second, 4)
	should.Nil(err)
	tw.Start()
	count := 0
	wg := sync.WaitGroup{}
	wg.Add(1)
	startTime := time.Now()
	task := tw.AddCron(3*time.Second, func() {
		t.Logf("do func :time=%v\n", time.Now().Sub(startTime))
		count++
		if count >= 2 {
			wg.Done()
		}
	})
	t.Logf("task = %v\n", task)
	wg.Wait()
	err = tw.Remove(task)
	should.Nil(err)
	tw.Stop()
}
