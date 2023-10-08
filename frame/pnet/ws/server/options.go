package server

import (
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/pruntime"
	"github.com/panjf2000/gnet/v2"
	"time"
)

// NewOptions
//
//	@Description: 创建 Options
//	@param opts
//	@return Options
func NewOptions(opts ...Option) (*Options, error) {
	options := &Options{
		GNetOptions: gnet.Options{
			// NumEventLoop is set up to start the given number of event-loop goroutine.
			NumEventLoop: pruntime.MaxProcess(),
			// LB represents the load-balancing algorithm used when assigning new connections.
			LB: gnet.LeastConnections,
			// ReuseAddr indicates whether to set up the SO_REUSEADDR socket option.
			ReuseAddr: false,
			// ReusePort indicates whether to set up the SO_REUSEPORT socket option.
			ReusePort: false,
			// MulticastInterfaceIndex is the index of the interface name where the multicast UDP addresses will be bound to.
			MulticastInterfaceIndex: 1,
			// ============================= Options for both server-side and client-side =============================
			// ReadBufferCap is the maximum number of bytes that can be read from the peer when the readable event comes.
			// The default value is 64KB, it can either be reduced to avoid starving the subsequent connections or increased
			// to read more data from a socket.
			ReadBufferCap: 16 * 1024,
			// WriteBufferCap is the maximum number of bytes that a static outbound buffer can hold,
			// if the data exceeds this value, the overflow will be stored in the elastic linked list buffer.
			// The default value is 64KB.
			WriteBufferCap: 32 * 1024,
			// LockOSThread is used to determine whether each I/O event-loop is associated to an OS thread, it is useful when you
			// need some kind of mechanisms like thread local storage, or invoke certain C libraries (such as graphics lib: GLib)
			// that require thread-level manipulation via cgo, or want all I/O event-loops to actually run in parallel for a
			// potential higher performance.
			LockOSThread: false,
			// Ticker indicates whether the ticker has been set up.
			Ticker: true,
			// TCPKeepAlive sets up a duration for (SO_KEEPALIVE) socket option.
			TCPKeepAlive: 5 * time.Minute,
			// TCPNoDelay controls whether the operating system should delay
			// packet transmission in hopes of sending fewer packets (Nagle's algorithm).
			TCPNoDelay: gnet.TCPNoDelay,
			// SocketRecvBuffer sets the maximum socket receive buffer in bytes.
			SocketRecvBuffer: 16 * 1024,
			// SocketSendBuffer sets the maximum socket send buffer in bytes.
			SocketSendBuffer: 256 * 1024,
			// Logger is the customized logger for logging info
			Logger: plog.SugarLogger(),
		},
		UnregisterSessionLife: 20,
		CheckSessionInterval:  30 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options, nil
}

type Options struct {
	GNetOptions gnet.Options
	// 未注册session的存活时间，单位秒
	UnregisterSessionLife int64
	// 检查session间隔
	CheckSessionInterval time.Duration
}

type Option func(options *Options)

func WithGNetOption(options ...gnet.Option) Option {
	return func(opts *Options) {
		for _, option := range options {
			option(&opts.GNetOptions)
		}
	}
}
