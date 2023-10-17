package predis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// newOptions
//
//	@Description: build default Options
//	@return *Options
func newOptions() *Options {
	return &Options{
		MaxIdle:         2,
		MaxActive:       2,
		IdleTimeout:     1 * time.Minute,
		Wait:            true,
		MaxConnLifetime: 0,
		DialOptions: []redis.DialOption{
			redis.DialConnectTimeout(10 * time.Second),
			redis.DialReadTimeout(5 * time.Second),
			redis.DialWriteTimeout(5 * time.Second),
		},
	}
}

// Options
//
//	@Description: options for configuring redis pool
type Options struct {
	// Maximum number of idle connections in the pool.
	MaxIdle int
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int
	// stop connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration
	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool
	// stop connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration
	// specifies options for dialing a Redis server.
	DialOptions []redis.DialOption
}

// PoolOption configures Options
type PoolOption func(*Options)

func WithPoolOption(options *Options) PoolOption {
	return func(opts *Options) {
		*opts = *options
	}
}

func WithMaxIdle(maxIdle int) PoolOption {
	return func(opts *Options) {
		opts.MaxIdle = maxIdle
	}
}

func WithMaxActive(maxActive int) PoolOption {
	return func(opts *Options) {
		opts.MaxActive = maxActive
	}
}

func WithIdleTimeout(idleTimeout time.Duration) PoolOption {
	return func(opts *Options) {
		opts.IdleTimeout = idleTimeout
	}
}

func WithWait(wait bool) PoolOption {
	return func(opts *Options) {
		opts.Wait = wait
	}
}

func WithMaxConnLifetime(maxConnLifetime time.Duration) PoolOption {
	return func(opts *Options) {
		opts.MaxConnLifetime = maxConnLifetime
	}
}

func WithDialOptions(dialOptions ...redis.DialOption) PoolOption {
	return func(opts *Options) {
		opts.DialOptions = dialOptions
	}
}
