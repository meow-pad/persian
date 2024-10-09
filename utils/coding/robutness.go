package coding

import (
	"github.com/meow-pad/persian/utils/loggers"
	"go.uber.org/zap"
)

// CatchPanicError
//
//	@Description: 捕捉panic错误
//	@param errMsg string
//	@param fields ...zap.Field
func CatchPanicError(errMsg string, errFunc func(), fields ...zap.Field) {
	if err := recover(); err != nil {
		if len(fields) > 0 {
			loggers.Error(errMsg, append(fields, zap.Any("error", err), zap.StackSkip("stack", 1))...)
		} else {
			loggers.Error(errMsg, zap.Any("error", err), zap.StackSkip("stack", 1))
		}
		if errFunc != nil {
			errFunc()
		}
	} // end of if
}

// HandlePanicError
//
//	@Description: 处理panic错误
//	@param errMsg
//	@param handleFunc
//	@param fields
func HandlePanicError(errMsg string, handleFunc func(err any), fields ...zap.Field) {
	if err := recover(); err != nil {
		if len(fields) > 0 {
			loggers.Error(errMsg, append(fields, zap.Any("error", err), zap.StackSkip("stack", 1))...)
		} else {
			loggers.Error(errMsg, zap.Any("error", err), zap.StackSkip("stack", 1))
		}
		if handleFunc != nil {
			handleFunc(err)
		}
	} else {
		if handleFunc != nil {
			handleFunc(nil)
		}
	}
}
