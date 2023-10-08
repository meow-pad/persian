package plog

import (
	"errors"
	"fmt"
	"github.com/1set/gut/yos"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
)

// NewLogger
//
//	@Description: 构建新的日志系统
//	@param config *Config
//	@param opts ...zap.Option
//	@return *Logger
func NewLogger(config *Config, opts ...zap.Option) (*Logger, error) {
	logger := &Logger{}
	err := logger.init(config, opts...)
	return logger, err
}

// Logger
//
//	@Description: 日志封装对象
type Logger struct {
	logCfg *Config     // 日志配置
	inner  *zap.Logger // 日志对象
	sugar  sugarLogger // 语法糖日志
}

// init
//
//	@Description: 初始化日志系统
//	@receiver wrapper
//	@param config
//	@param opts
//	@return error
func (logger *Logger) init(config *Config, opts ...zap.Option) error {
	var encoder = config.ZapLogEncoder
	if encoder == nil {
		return errors.New("less ZapLogEncoder in config")
	}
	logLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= config.LogLevel
	})
	var logCore zapcore.Core
	if rotateWriter := logger.rotateLogsWriter(config); rotateWriter != nil {
		//logCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(
		//	zapcore.AddSync(rotateWriter), zapcore.AddSync(os.Stdout)), logLevel)
		// 配置了日志文件,则去除控制台输出
		logCore = zapcore.NewCore(encoder, zapcore.AddSync(rotateWriter), logLevel)
	} else {
		logCore = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel)
	}
	log := zap.New(zapcore.NewTee(logCore))
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	logger.logCfg = config
	logger.inner = log
	// 设置语法糖日志核心
	logger.sugar.core = logger.inner.WithOptions(zap.AddCallerSkip(1)) // 额外跳过一层调用
	// 默认日志进行额外初始化
	if logger == &defaultLogger {
		// 触发Fatal日志配置
		if initFatalLog != nil {
			initFatalLog(logger.logCfg.LogsDirectory)
		}
	} // end of if
	return nil
}

func (logger *Logger) rotateLogsWriter(config *Config) io.Writer {
	if len(config.LogsDirectory) <= 0 {
		return nil
	}
	if err := yos.MakeDir(config.LogsDirectory); err != nil {
		if logger.inner == nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("make logs directory(%s) error:%v", config.LogsDirectory, err))
			os.Exit(1)
		} else {
			logger.inner.Error("make logs directory error", zap.String("dir", config.LogsDirectory), zap.Error(err))
		}
	} // end of if
	writer, err := rotateLogs.New(
		filepath.Join(config.LogsDirectory, config.LogsFilenamePattern),
		rotateLogs.WithMaxAge(config.LogsMaxAge),
		rotateLogs.WithRotationTime(config.LogsRotationTime),
		rotateLogs.WithRotationCount(config.LogsRotationCount),
		rotateLogs.WithRotationSize(config.LogsRotationSize),
		rotateLogs.WithLinkName(filepath.Join(config.LogsDirectory, "current.log")),
	)
	if err != nil {
		if logger.inner == nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("build file-rotateLogs logs error:%v", err))
			os.Exit(1)
		} else {
			logger.inner.Error("build file-rotateLogs logs error", zap.Error(err))
		}
	}
	return writer
}

// Enabled
//
//	@Description: 判定是否是可输出的日志等级
//	@receiver wrapper *Logger
//	@param level zapcore.Level 日志等级
//	@return bool
func (logger *Logger) Enabled(level zapcore.Level) bool {
	if logger.inner != nil {
		return logger.inner.Core().Enabled(level)
	} else {
		return level >= zap.DPanicLevel
	}
}

// Debug
//
//	@Description: 记录 DebugLevel 级别的日志消息.
//	@receiver wrapper *Logger
//	@param msg string
//	@param fields ...zap.Field
func (logger *Logger) Debug(msg string, fields ...zap.Field) {
	if logger.inner != nil {
		if logger.logCfg.Development {
			if logger.inner.Core().Enabled(zap.DebugLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		logger.inner.Debug(msg, fields...)
	} // end of if
}

// Info
//
//	@Description: 记录 InfoLevel 级别的日志消息.
//	@receiver wrapper *Logger
//	@param msg string
//	@param fields ...zap.Field
func (logger *Logger) Info(msg string, fields ...zap.Field) {
	if logger.inner != nil {
		if logger.logCfg.Development {
			if logger.inner.Core().Enabled(zap.InfoLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		logger.inner.Info(msg, fields...)
	} // end of if
}

// Warn
//
//	@Description: 记录 WarnLevel 级别的日志消息.
//	@receiver wrapper *Logger
//	@param msg string
//	@param fields ...zap.Field
func (logger *Logger) Warn(msg string, fields ...zap.Field) {
	if logger.inner != nil {
		if logger.logCfg.Development {
			if logger.inner.Core().Enabled(zap.WarnLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		logger.inner.Warn(msg, fields...)
	} // end of if
}

// Error
//
//	@Description: 记录 ErrorLevel 级别的日志消息.
//	@receiver wrapper *Logger
//	@param msg string
//	@param fields ...zap.Field
func (logger *Logger) Error(msg string, fields ...zap.Field) {
	if logger.inner != nil {
		if logger.logCfg.Development {
			if logger.inner.Core().Enabled(zap.ErrorLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		logger.inner.Error(msg, fields...)
	} // end of if
}

// DPanic
//
//	@Description:
//		记录 DPanicLevel 级别的日志消息.
//		如果日志系统处于development模式,该方法将触发panics(DPanic意味着"development panic").
//		这主要用于捕捉可自恢复,但不应出现的错误.
//	@receiver wrapper *Logger
//	@param msg string
//	@param fields ...zap.Field
func (logger *Logger) DPanic(msg string, fields ...zap.Field) {
	if logger.inner != nil {
		if logger.logCfg.Development {
			if logger.inner.Core().Enabled(zap.DPanicLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		logger.inner.DPanic(msg, fields...)
	} // end of if
}

// Panic
//
//	@Description:
//		记录 PanicLevel 级别的日志消息.
//		调用后,日志系统随即会panics,即使PanicLevel级别不可用.
//	@receiver wrapper *Logger
//	@param msg string
//	@param fields ...zap.Field
func (logger *Logger) Panic(msg string, fields ...zap.Field) {
	if logger.inner != nil {
		if logger.logCfg.Development {
			if logger.inner.Core().Enabled(zap.PanicLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		logger.inner.Panic(msg, fields...)
	} else {
		panic(msg)
	}
}

// Fatal
//
//	@Description:
//		记录 FatalLevel 级别的日志消息.
//		日志系统随后会调用os.Exit(1), 即使FatalLevel级别不可用
//	@receiver wrapper *Logger
//	@param msg string
//	@param fields ...zap.Field
func (logger *Logger) Fatal(msg string, fields ...zap.Field) {
	if logger.inner != nil {
		if logger.logCfg.Development {
			if logger.inner.Core().Enabled(zap.FatalLevel) {
				msg, fields = encodeErrorFieldLayout(msg, fields...)
			}
		} // end of if
		logger.inner.Fatal(msg, fields...)
	} else {
		fmt.Print(msg)
		os.Exit(1)
	}
}

// Sync
//
//	@Description:
//		调用底层同步方法,将所有缓存日志刷新到输出媒介.
//		应用程序关闭前,应先通过Sync来确保日志不遗漏.
//	@receiver wrapper *Logger
func (logger *Logger) Sync() {
	if logger.inner != nil {
		if err := logger.inner.Sync(); err != nil {
			logger.inner.Error("inner sync error:", zap.Error(err))
		}
	}
}

func (logger *Logger) LoggerLevel() zapcore.Level {
	return logger.logCfg.LogLevel
}
