package worker

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testTask struct {
	t         *testing.T
	id        int
	sleepTime time.Duration
}

func (task *testTask) Id() int {
	return task.id
}

func (task *testTask) Run(*GoroutineLocal) {
	if task.sleepTime > 0 {
		time.Sleep(task.sleepTime)
	}
	task.t.Logf("run task %d", task.id)
}

func TestNonBlockingFixedPoolWorker(t *testing.T) {
	should := require.New(t)
	slotNum := 5
	queueSize := 2
	pool, err := NewFixedWorkerPool(slotNum, queueSize, false)
	should.Nil(err)
	for i := 0; i < slotNum; i++ {
		for j := 0; j < queueSize; j++ {
			task := &testTask{t: t, id: i, sleepTime: 3 * time.Second}
			sErr := pool.Submit(i, task.Run)
			should.Nil(sErr)
		}
	}
	time.Sleep(1 * time.Second)
	task1 := &testTask{t: t, id: 10086, sleepTime: 0}
	err = pool.Submit(task1.id, task1.Run)
	should.Nil(err)
	task2 := testTask{t: t, id: 10086, sleepTime: 0}
	err = pool.Submit(task2.id, task2.Run)
	should.True(errors.Is(err, ErrWorkerPoolQueueIsFull))
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	err = pool.Shutdown(ctx)
	should.Nil(err)
}
