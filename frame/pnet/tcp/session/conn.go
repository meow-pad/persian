package session

import (
	"github.com/panjf2000/gnet/v2"
	"io"
	"net"
	"persian/errdef"
	"sync/atomic"
)

type Writer interface {
	io.Writer     // not goroutine-safe
	io.ReaderFrom // not goroutine-safe

	// Writev writes multiple byte slices to peer synchronously, it's not goroutine-safe,
	// you must invoke it within any method in EventHandler.
	Writev(bs [][]byte) (n int, err error)

	// Flush writes any buffered data to the underlying connection, it's not goroutine-safe,
	// you must invoke it within any method in EventHandler.
	Flush() (err error)

	// OutboundBuffered returns the number of bytes that can be read from the current buffer.
	// it's not goroutine-safe, you must invoke it within any method in EventHandler.
	OutboundBuffered() (n int)

	// AsyncWrite writes bytes to peer asynchronously, it's goroutine-safe,
	// you don't have to invoke it within any method in EventHandler,
	// usually you would call it in an individual goroutine.
	//
	// Note that it will go synchronously with UDP, so it is needless to call
	// this asynchronous method, we may disable this method for UDP and just
	// return ErrUnsupportedOp in the future, therefore, please don't rely on
	// this method to do something important under UDP, if you're working with UDP,
	// just call Conn.Write to send back your data.
	AsyncWrite(buf []byte, callback func(c Conn, err error) error) (err error)

	// AsyncWritev writes multiple byte slices to peer asynchronously,
	// you don't have to invoke it within any method in EventHandler,
	// usually you would call it in an individual goroutine.
	AsyncWritev(bs [][]byte, callback func(c Conn, err error) error) (err error)
}

// Conn
//
//	@Description: connection
type Conn interface {
	net.Conn
	gnet.Reader
	Writer

	// Hash
	//	@Description: get conn hash code
	//	@return uint64
	//
	Hash() uint64

	// Context returns a user-defined context, it's not goroutine-safe,
	// you must invoke it within any method in EventHandler.
	Context() (ctx any)

	// SetContext sets a user-defined context, it's not goroutine-safe,
	// you must invoke it within any method in EventHandler.
	SetContext(ctx any)

	// IsClosed
	//	@Description: get connection status and closed reason
	//	@return bool closed or not
	//  @return error closed reason
	//
	IsClosed() (bool, error)
}

type BaseConn struct {
	hash     uint64
	closed   atomic.Bool
	closeErr error
}

func (conn *BaseConn) Init(pConn net.Conn, fromClient bool) error {
	if pConn == nil {
		return errdef.ErrNilValue
	}
	addr := pConn.RemoteAddr().String()
	if fromClient {
		addr = pConn.LocalAddr().String()
	}
	for _, ch := range addr {
		conn.hash = 31*conn.hash + uint64(uint32(ch))
	}
	conn.closed.Store(false)
	return nil
}

func (conn *BaseConn) Hash() uint64 {
	return conn.hash
}

func (conn *BaseConn) IsClosed() (bool, error) {
	return conn.closed.Load(), conn.closeErr
}

func (conn *BaseConn) ToClosed(reason error) bool {
	if conn.closed.CompareAndSwap(false, true) {
		conn.closeErr = reason
		return true
	}
	return false
}
