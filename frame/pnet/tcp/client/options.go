package client

import (
	"time"
)

func newOptions(opts ...Option) *Options {
	options := &Options{
		ReadBufferCap:    16 * 1024,
		WriteBufferCap:   32 * 1024,
		WriteQueueCap:    100,
		TCPKeepAlive:     5 * time.Minute,
		TCPNoDelay:       true,
		SocketRecvBuffer: 16 * 1024,
		SocketSendBuffer: 32 * 1024,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

type Options struct {
	Name string
	// 读缓冲容量
	ReadBufferCap int
	// 写缓冲容量
	WriteBufferCap int
	// 同时提交写队列上限，队列满时异步写直接返回错误
	WriteQueueCap int
	// 保活时间
	TCPKeepAlive time.Duration
	// TCPNoDelay controls whether the operating system should delay
	// packet transmission in hopes of sending fewer packets (Nagle's algorithm).
	TCPNoDelay bool
	// socket读缓冲区
	SocketRecvBuffer int
	// socket写缓冲区
	SocketSendBuffer int
}

type Option func(*Options)

func WithName(value string) Option {
	return func(options *Options) {
		options.Name = value
	}
}

func WithReadBufferCap(cap int) Option {
	return func(options *Options) {
		options.ReadBufferCap = cap
	}
}

func WithWriteBufferCap(cap int) Option {
	return func(options *Options) {
		options.WriteBufferCap = cap
	}
}

func WithWriteQueueCap(cap int) Option {
	return func(options *Options) {
		options.WriteQueueCap = cap
	}
}

func WithTCPKeepAlive(value time.Duration) Option {
	return func(options *Options) {
		options.TCPKeepAlive = value
	}
}

func WithTCPNoDelay(value bool) Option {
	return func(options *Options) {
		options.TCPNoDelay = value
	}
}

func WithSocketRecvBuffer(cap int) Option {
	return func(options *Options) {
		options.SocketRecvBuffer = cap
	}
}

func WithSocketSendBuffer(cap int) Option {
	return func(options *Options) {
		options.SocketSendBuffer = cap
	}
}
