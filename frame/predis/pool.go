package predis

import (
	"context"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
)

var (
	ErrLessConn = errors.New("less redis conn")
)

func NewPool(dialUrl string, opts ...PoolOption) *Pool {
	pool := &Pool{}
	pool.init(dialUrl, opts...)
	return pool
}

type Pool struct {
	options *Options
	dialUrl string
	inner   *redis.Pool
}

func (pool *Pool) init(dialUrl string, opts ...PoolOption) {
	pool.dialUrl = dialUrl
	pool.options = newOptions()
	for _, opt := range opts {
		opt(pool.options)
	}
}

// Start
//
//	@Description: 开启连接池
//	@receiver pool
//	@return error
func (pool *Pool) Start() error {
	pool.inner = &redis.Pool{
		MaxIdle:         pool.options.MaxIdle,
		MaxActive:       pool.options.MaxActive,
		IdleTimeout:     pool.options.IdleTimeout,
		Wait:            pool.options.Wait,
		MaxConnLifetime: pool.options.MaxConnLifetime,
		Dial:            func() (redis.Conn, error) { return redis.DialURL(pool.dialUrl, pool.options.DialOptions...) },
	}
	return nil
}

// Stop
//
//	@Description: 关闭连接池
//	@receiver pool
//	@return error
func (pool *Pool) Stop() error {
	return pool.inner.Close()
}

// Get
//
//	@Description: 获取一个连接
//	@receiver pool
//	@return redis.Conn
func (pool *Pool) Get() redis.Conn {
	return pool.inner.Get()
}

// Take
//
//	@Description: 获取一个连接
//	@receiver pool
//	@param ctx
//	@return redis.Conn
//	@return error
func (pool *Pool) Take(ctx context.Context) (redis.Conn, error) {
	return pool.inner.GetContext(ctx)
}

// Back
//
//	@Description: 返还一个连接
//	@receiver pool
//	@param conn
func (pool *Pool) Back(conn redis.Conn) {
	if err := conn.Close(); err != nil {
		plog.Error("close redis conn error", pfield.Error(err))
	}
}

// NewScript
//
//	@Description: 构建脚本
//	@receiver pool
//	@param source
//	@param keyCount 这里的keyCount用于决定执行时,前多少参数是KEYS（如KEYS[1],KEYS[2]...），那么之后参数的都是ARGV（如ARGV[1],ARGV[2]...）
//	@param autoLoad 是否自动加载（除非这个脚本大概率不会被执行，否则这里都应该传true）
//	@return *redis.Script
//	@return error
func (pool *Pool) NewScript(source string, keyCount int, autoLoad bool) (*redis.Script, error) {
	script := redis.NewScript(keyCount, source)
	if autoLoad {
		conn := pool.inner.Get()
		if conn == nil {
			return nil, ErrLessConn
		}
		defer pool.Back(conn)
		if err := script.Load(conn); err != nil {
			return nil, err
		}
	}
	return script, nil
}

// Do
//
//	@Description: 执行redis命令行操作
//	@receiver pool
//	@param cmd
//	@param args
//	@return interface{}
//	@return error
func (pool *Pool) Do(cmd string, args ...any) (any, error) {
	if len(cmd) < 0 {
		return nil, errdef.ErrInvalidParams
	}
	conn := pool.inner.Get()
	if conn == nil {
		return nil, ErrLessConn
	}
	defer pool.Back(conn)
	return conn.Do(cmd, args...)
}

// DoScript
//
// @Description: 执行lua脚本
// @receiver pool
// @param script
// @param args
// @return any
// @return error
func (pool *Pool) DoScript(script *redis.Script, args ...any) (any, error) {
	if script == nil {
		return nil, errdef.ErrInvalidParams
	}
	conn := pool.inner.Get()
	if conn == nil {
		return nil, ErrLessConn
	}
	defer pool.Back(conn)
	return script.Do(conn, args...)
}

// DoCommand
//
//	@Description: 以 Command 形式执行
//	@receiver pool
//	@param cmd
//	@return error
func (pool *Pool) DoCommand(cmd *Command) error {
	if len(cmd.Name) < 0 || cmd.Script == nil || cmd.ReplyHandler == nil {
		return errdef.ErrInvalidParams
	}
	conn := pool.inner.Get()
	if conn == nil {
		return ErrLessConn
	}
	defer pool.Back(conn)
	if cmd.Script != nil {
		reply, err := cmd.Script.Do(conn, cmd.Args...)
		cmd.ReplyHandler(err, reply, cmd.ReplyArgs...)
	} else {
		reply, err := conn.Do(cmd.Name, cmd.Args...)
		cmd.ReplyHandler(err, reply, cmd.ReplyArgs...)
	}
	return nil
}

// BatchDo
//
//	@Description: 批量执行 Command
//	@receiver pool
//	@param cmdArr
//	@return error
func (pool *Pool) BatchDo(cmdArr ...*Command) error {
	if len(cmdArr) <= 0 {
		return nil
	}
	conn := pool.inner.Get()
	if conn == nil {
		return ErrLessConn
	}
	defer pool.Back(conn)
	for _, cmd := range cmdArr {
		if cmd == nil || (len(cmd.Name) <= 0 && cmd.Script == nil) || cmd.ReplyHandler == nil {
			return errdef.ErrInvalidParams
		}
		if cmd.Script != nil {
			if err := cmd.Script.SendHash(conn, cmd.Args...); err != nil {
				return err
			}
		} else {
			if err := conn.Send(cmd.Name, cmd.Args...); err != nil {
				return err
			}
		}
	} // end of for
	if err := conn.Flush(); err != nil {
		return err
	}
	for _, cmd := range cmdArr {
		reply, err := conn.Receive()
		cmd.ReplyHandler(err, reply, cmd.ReplyArgs...)
	}
	return nil
}
