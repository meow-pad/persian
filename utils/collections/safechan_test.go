package collections

import (
	"errors"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func printSafeChanGResult[T any](t *testing.T, should *require.Assertions, channel *SafeChannel[T]) {
	result, err := channel.Get()
	if errors.Is(err, ErrEmptySafeChan) {
		t.Logf("less result")
	} else {
		should.Nil(err)
		t.Logf("result:%v", result)
	}
}

func printSafeChanBGResult[T any](t *testing.T, should *require.Assertions, channel *SafeChannel[T]) {
	result, err := channel.BlockingGet(nil)
	should.Nil(err)
	t.Logf("result:%v", result)
}

func TestSafeChan_Base1(t *testing.T) {
	should := require.New(t)
	channel := NewSafeChan[string](3)
	defer channel.Close()
	err := error(nil)
	err = channel.Put("tom")
	should.Nil(err)
	err = channel.Put("jimmy")
	should.Nil(err)
	err = channel.Put("cate")
	should.Nil(err)

	printSafeChanBGResult(t, should, channel)
	printSafeChanBGResult(t, should, channel)
	printSafeChanBGResult(t, should, channel)
}

func TestSafeChan_Get(t *testing.T) {
	should := require.New(t)
	channel := NewSafeChan[string](3)
	err := error(nil)
	err = channel.Put("tom")
	should.Nil(err)
	printSafeChanGResult(t, should, channel)
	printSafeChanGResult(t, should, channel)
	go func() {
		err = channel.Put("jimmy")
		should.Nil(err)
		time.Sleep(1 * time.Second)
		err = channel.Put("cate")
		should.Nil(err)

		time.Sleep(1 * time.Second)
		channel.Close()
	}()
	t.Logf("time:%v", time.Now().UnixMilli())
	printSafeChanBGResult(t, should, channel)
	t.Logf("time:%v", time.Now().UnixMilli())
	printSafeChanBGResult(t, should, channel)
	t.Logf("time:%v", time.Now().UnixMilli())
	_, err = channel.BlockingGet(nil)
	should.True(errors.Is(err, ErrClosedSafeChan))
	t.Logf("time:%v", time.Now().UnixMilli())
}

func TestSafeChan_Listen(t *testing.T) {
	should := require.New(t)
	channel := NewSafeChan[string](3)
	go func() {
		err := channel.Listen(func(val string) bool {
			t.Logf("result:%s", val)
			return true
		})
		should.True(errors.Is(err, ErrClosedSafeChan))
	}()
	for i := 0; i < 10; i++ {
		err := channel.BlockingPut(nil, strconv.Itoa(i))
		should.Nil(err)
	}
	time.Sleep(1 * time.Second)
	channel.Close()
}
