package coding

import (
	"go.uber.org/zap"
	"persian/utils/loggers"
)

// CachePanicError
//
//	@Description: 捕捉panic错误
//	@param errMsg string
//	@param fields ...zap.Field
func CachePanicError(errMsg string, errFunc func(), fields ...zap.Field) {
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
