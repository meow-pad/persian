package pnet

import "errors"

var (
	ErrNilMessage        = errors.New("nil message")
	ErrEmptyEncodeBuffer = errors.New("encoded buffer length is 0")
	ErrInvalidMagic      = errors.New("invalid magic")
	ErrMessageTooLarge   = errors.New("message is too large")
	ErrClosedConn        = errors.New("connection is closed")
	ErrRegisteredSession = errors.New("session is registered")
	ErrInvalidSessionId  = errors.New("invalid session id")
	ErrOutOfReadCap      = errors.New("out of read capacity")
	ErrOutOfWriteCap     = errors.New("out of write capacity")
	ErrClosedClient      = errors.New("client is closed")
	ErrWriteQueueFull    = errors.New("write queue is full")
	ErrClientActiveClose = errors.New("client active close")
)
