package errdef

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrInvalidParams  = errors.New("invalid params")
	ErrNilValue       = errors.New("nil value")
	ErrUnsupportedOp  = errors.New("unsupported operation")
	ErrNotFound       = errors.New("not found")
	ErrNotInitialized = errors.New("not initialized")
	ErrNotStarted     = errors.New("not started")
)
