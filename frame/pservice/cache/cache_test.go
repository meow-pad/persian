package cache

import (
	"github.com/gomodule/redigo/redis"
	"github.com/meow-pad/persian/frame/predis"
	"github.com/meow-pad/persian/utils/timewheel"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

const (
	testRedisHost = "redis://192.168.91.130:6379"
	testRedisPass = "123456"
)

func TestCache_AddAndDel(t *testing.T) {
	should := require.New(t)
	pool := predis.NewPool(testRedisHost, predis.WithDialOptions(redis.DialPassword(testRedisPass)))
	err := pool.Start()
	should.Nil(err)
	defer func() {
		err = pool.Stop()
		should.Nil(err)
	}()
	tw, err := timewheel.NewTimeWheel(200*time.Millisecond, 40)
	should.Nil(err)
	tw.Start()
	defer tw.Stop()
	cache, err := NewCache(pool, tw)
	should.Nil(err)
	signature := [8]byte{6, 7, 8, 2, 1, 7, 4, 1}
	count := atomic.Int32{}
	key := "123"
	task, err := cache.AddOrUpdate(key, "val123", signature, 4,
		2, 1, 1, func(result error, retryCount int) {
			count.Add(1)
			t.Logf("keepalive result=%v,retryCount=%v", result, retryCount)
		})
	should.Nil(err)
	time.Sleep(5 * time.Second)
	// 此时应该会触发两次更新
	should.True(count.Load() == 2, count.Load())
	// 获取缓存值
	{
		cValue, ok, gErr := cache.Get(key)
		should.Nil(gErr)
		should.True(ok)
		t.Logf("cValue=%s", cValue)
	}
	{
		_, ok, gErr := cache.Get("invalid_key")
		should.Nil(gErr)
		should.False(ok)
	}
	task.Cancel()
	time.Sleep(3 * time.Second)
	// 停止了应该没有变化
	should.True(count.Load() == 2)
	err = cache.Delete(key, signature)
	should.Nil(err)
}
