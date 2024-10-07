package cache

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/frame/predis"
	"github.com/meow-pad/persian/utils/coding"
	"github.com/meow-pad/persian/utils/timewheel"
	"strings"
	"sync/atomic"
	"time"
)

var (
	ErrCacheExist = errors.New("cache exist")
)

const (
	luaAddKeyCount = 1
	luaAddSource   = `
local cacheKey = KEYS[1]
local cacheValue = ARGV[1]
local signature = ARGV[2]
local lenSignature = ARGV[3]
local expireSec = ARGV[4]
local existValue = redis.call("get", cacheKey)
if (existValue) then
    if (string.sub(existValue, 1, lenSignature) == signature) then
        if (redis.call("set", cacheKey, cacheValue, "PX", expireSec)) then
            return "1:update"
        else
            return "0:errorSet"
        end
    else
        --签名未匹配，被其他服务占用
        return "0:cacheExist:"..string.sub(existValue, lenSignature+1)
    end
else
    --不存在则直接设置
    if (redis.call("set", cacheKey, cacheValue, "NX", "PX", expireSec)) then
        return "1:insert"
    else
        return "0:errorSetNX"
    end
end
`
	luaDelKeyCount = 1
	luaDelSource   = `
local cacheKey = KEYS[1]
local signature = ARGV[1]
local lenSignature = ARGV[2]
local existValue = redis.call("get", cacheKey)
if (existValue) then
    if (string.sub(existValue, 1, lenSignature) == signature) then
        if (redis.call("del", cacheKey)) then
            return "1:delete"
        else
            return "0:errorDel"
        end
    else
        --签名未匹配，被其他服务占用
        return "0:invalidSign"
    end
else
    --不存在则直接返回
    return "1:notExist"
end
`

	replyErrorStart      = "0:"
	replyErrorStartLen   = len(replyErrorStart)
	replyErrorCacheExist = "0:cacheExist"

	signatureLen = 9
)

func newKeepaliveTask(
	cache *Cache,
	key string, value string,
	signature string, expireSec int64,
	interval time.Duration,
	retryNum int, retryDelay time.Duration,
	callback KeepaliveCallback) *KeepaliveTask {
	return &KeepaliveTask{
		cache:      cache,
		key:        key,
		value:      value,
		signature:  signature,
		expireSec:  expireSec,
		interval:   interval,
		retryNum:   retryNum,
		retryDelay: retryDelay,
		callback:   callback,
	}
}

type KeepaliveTask struct {
	cache      *Cache
	tPointer   atomic.Pointer[timewheel.Task]
	retryCount int

	key        string
	value      string
	signature  string
	expireSec  int64
	interval   time.Duration
	retryNum   int // 可重试次数
	retryDelay time.Duration
	callback   KeepaliveCallback
}

func (task *KeepaliveTask) startTask() {
	task.addTask(task.interval)
}

func (task *KeepaliveTask) addTask(delay time.Duration) {
	t := task.cache.secTimer.Add(delay, func() {
		if task.tPointer.Load() == nil || task.cache == nil {
			// 已经取消
			return
		}
		err := task.cache.insertOrUpdate(task.key, task.value, task.signature, task.expireSec)
		// 先执行回调
		if sErr := coding.SafeRunSimple(func() {
			task.callback(err, task.retryCount)
		}); sErr != nil {
			plog.Error("callback runtime error:", pfield.Error(sErr))
		}
		// 之后的调度
		if err != nil {
			//	失败了则重试
			if task.retryNum > task.retryCount {
				task.retryCount++
				task.addTask(task.retryDelay)
			}
		} else {
			// 成功则以正常延迟继续
			task.addTask(task.interval)
		}
	})
	task.tPointer.Store(t)
	task.retryCount = 0
}

func (task *KeepaliveTask) Cancel() {
	cache := task.cache
	task.cache = nil
	t := task.tPointer.Swap(nil)
	if t != nil {
		if err := cache.secTimer.Remove(t); err != nil {
			plog.Error("remove time wheel task error:", pfield.Error(err))
		}
	}
}

func buildSignature(signature [8]byte) string {
	return string(signature[:]) + " "
}

func NewCache(pool *predis.Pool, secTimer *timewheel.TimeWheel) (*Cache, error) {
	cache := &Cache{}
	if err := cache.init(pool, secTimer); err != nil {
		return nil, err
	}
	return cache, nil
}

type Cache struct {
	pool      *predis.Pool
	secTimer  *timewheel.TimeWheel
	anuScript *redis.Script
	delScript *redis.Script
}

func (cache *Cache) init(pool *predis.Pool, secTimer *timewheel.TimeWheel) (err error) {
	if pool == nil || secTimer == nil {
		err = errdef.ErrInvalidParams
		return
	}
	cache.pool = pool
	cache.secTimer = secTimer
	cache.anuScript, err = cache.pool.NewScript(luaAddSource, luaAddKeyCount, true)
	cache.delScript, err = cache.pool.NewScript(luaDelSource, luaDelKeyCount, true)
	return
}

type KeepaliveCallback func(result error, retryCount int)

// AddOrUpdate
//
//	@Description: 添加或更新
//	@receiver cache
//	@param key
//	@param value
//	@param signature 签名
//	@param expireSec 过期时间
//	@param keepaliveSec 保活间隔时间
//	@param keepaliveRetry 保活失败重试次数
//	@param retryDelaySec 保活失败重试延迟
//	@param callback 保活结果回调
//	@return *KeepaliveTask 保活任务对象，可以通过该对象取消
//	@return error
func (cache *Cache) AddOrUpdate(key, value string, signature [8]byte, expireSec int64,
	keepaliveSec int64, keepaliveRetry int, retryDelaySec int64, callback KeepaliveCallback) (*KeepaliveTask, error) {
	if len(key) <= 0 || expireSec <= 0 {
		return nil, errdef.ErrInvalidParams
	}
	if keepaliveSec > 0 && callback == nil {
		return nil, fmt.Errorf("callback cant be nil")
	}
	if keepaliveRetry > 0 && retryDelaySec <= 0 {
		return nil, fmt.Errorf("retryDelaySec should be greater than 0")
	}
	signStr := buildSignature(signature)
	value = signStr + value
	if err := cache.insertOrUpdate(key, value, signStr, expireSec); err != nil {
		return nil, err
	}
	if keepaliveSec > 0 {
		kpTask := newKeepaliveTask(cache, key, value, signStr, expireSec, time.Duration(keepaliveSec)*time.Second,
			keepaliveRetry, time.Duration(retryDelaySec)*time.Second, callback)
		kpTask.startTask()
		return kpTask, nil
	}
	return nil, nil
}

func (cache *Cache) insertOrUpdate(key, value string, signature string, expireSec int64) error {
	conn := cache.pool.Get()
	if conn == nil {
		return predis.ErrLessConn
	}
	defer cache.pool.Back(conn)
	replyStr, err := redis.String(cache.anuScript.Do(conn, key, value, signature, len(signature), expireSec*1000))
	if err != nil {
		return err
	}
	//plog.Debug("insert/update distribution cache:",
	//	pfield.String("key", key), pfield.String("value", value),
	//	pfield.String("reply", replyStr),
	//	pfield.String("signature", signature),
	//)
	if strings.HasPrefix(replyStr, replyErrorStart) {
		if strings.HasPrefix(replyStr, replyErrorCacheExist) {
			return ErrCacheExist
		}
		return errors.New("insert/update error:" + replyStr[replyErrorStartLen:])
	}
	return nil
}

// Get
//
//	@Description: 获取当前值
//	@receiver cache
//	@param key
//	@return string 值
//	@return bool 缓存是否存在
//	@return error 错误
func (cache *Cache) Get(key string) (string, bool, error) {
	replyStr, err := redis.String(cache.pool.Do("get", key))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return "", false, nil
		}
		return "", false, err
	}
	if len(replyStr) >= signatureLen && replyStr[signatureLen-1] == ' ' {
		return replyStr[signatureLen:], true, nil
	} else {
		return replyStr, true, nil
	}
}

// Delete
//
//	@Description: 删除指定key值
//	@receiver cache
//	@param key
//	@param signature 签名，不能删除非自己维护的缓存
//	@return error
func (cache *Cache) Delete(key string, signature [8]byte) error {
	conn := cache.pool.Get()
	if conn == nil {
		return predis.ErrLessConn
	}
	defer cache.pool.Back(conn)
	signStr := buildSignature(signature)
	replyStr, err := redis.String(cache.delScript.Do(conn, key, signStr, len(signStr)))
	if err != nil {
		return err
	}
	if strings.Index(replyStr, replyErrorStart) == 0 {
		return errors.New("delete error:" + replyStr[replyErrorStartLen:])
	}
	return nil
}
