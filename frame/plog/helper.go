package plog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Enabled
//
//	@Description: 判定是否是可输出的日志等级
//	@param level zapcore.Level 日志等级
//	@return bool
func Enabled(level zapcore.Level) bool {
	if defaultLogger.inner != nil {
		return defaultLogger.inner.Core().Enabled(level)
	} else {
		return level >= zap.DPanicLevel
	}
}

// Debug
//
//	@Description: 记录 DebugLevel 级别的日志消息.
//	@param msg string
//	@param fields ...zap.Field
func Debug(msg string, fields ...zap.Field) {
	if defaultLogger.inner != nil {
		if defaultLogger.logCfg.Development {
			if defaultLogger.inner.Core().Enabled(zap.DebugLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		defaultLogger.inner.Debug(msg, fields...)
	} // end of if
}

// Info
//
//	@Description: 记录 InfoLevel 级别的日志消息.
//	@param msg string
//	@param fields ...zap.Field
func Info(msg string, fields ...zap.Field) {
	if defaultLogger.inner != nil {
		if defaultLogger.logCfg.Development {
			if defaultLogger.inner.Core().Enabled(zap.InfoLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		defaultLogger.inner.Info(msg, fields...)
	} // end of if
}

// Warn
//
//	@Description: 记录 WarnLevel 级别的日志消息.
//	@param msg string
//	@param fields ...zap.Field
func Warn(msg string, fields ...zap.Field) {
	if defaultLogger.inner != nil {
		if defaultLogger.logCfg.Development {
			if defaultLogger.inner.Core().Enabled(zap.WarnLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		defaultLogger.inner.Warn(msg, fields...)
	} // end of if
}

// Error
//
//	@Description: 记录 ErrorLevel 级别的日志消息.
//	@param msg string
//	@param fields ...zap.Field
func Error(msg string, fields ...zap.Field) {
	if defaultLogger.inner != nil {
		if defaultLogger.logCfg.Development {
			if defaultLogger.inner.Core().Enabled(zap.ErrorLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		defaultLogger.inner.Error(msg, fields...)
	} // end of if
}

// DPanic
//
//	@Description:
//		记录 DPanicLevel 级别的日志消息.
//		如果日志系统处于development模式,该方法将触发panics(DPanic意味着"development panic").
//		这主要用于捕捉可自恢复,但不应出现的错误.
//	@param msg string
//	@param fields ...zap.Field
func DPanic(msg string, fields ...zap.Field) {
	if defaultLogger.inner != nil {
		if defaultLogger.logCfg.Development {
			if defaultLogger.inner.Core().Enabled(zap.DPanicLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		defaultLogger.inner.DPanic(msg, fields...)
	} // end of if
}

// Panic
//
//	@Description:
//		记录 PanicLevel 级别的日志消息.
//		调用后,日志系统随即会panics,即使PanicLevel级别不可用.
//	@param msg string
//	@param fields ...zap.Field
func Panic(msg string, fields ...zap.Field) {
	if defaultLogger.inner != nil {
		if defaultLogger.logCfg.Development {
			if defaultLogger.inner.Core().Enabled(zap.PanicLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		defaultLogger.inner.Panic(msg, fields...)
	} else {
		panic(msg)
	}
}

// Fatal
//
//	@Description:
//		记录 FatalLevel 级别的日志消息.
//		日志系统随后会调用os.Exit(1), 即使FatalLevel级别不可用
//	@param msg string
//	@param fields ...zap.Field
func Fatal(msg string, fields ...zap.Field) {
	if defaultLogger.inner != nil {
		if defaultLogger.logCfg.Development {
			if defaultLogger.inner.Core().Enabled(zap.FatalLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		defaultLogger.inner.Fatal(msg, fields...)
	} else {
		if _, err := os.Stderr.WriteString(msg); err != nil {
			fmt.Printf("Fatal message:%s\nStderr write error:%v\n", msg, err)
		}
		os.Exit(1)
	}
}

// Sync
//
//	@Description:
//		调用底层同步方法,将所有缓存日志刷新到输出媒介.
//		应用程序关闭前,应先通过Sync来确保日志不遗漏.
func Sync() {
	if defaultLogger.inner != nil {
		if err := defaultLogger.inner.Sync(); err != nil {
			defaultLogger.inner.Error("logger sync error:", zap.Error(err))
		}
	}
}
