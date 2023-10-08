package plog

import (
	"fmt"
	"github.com/meow-pad/persian/utils/loggers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	consoleSeparator = " | "
	defaultLogDir    = "/var/logs/"
)

var (
	// 初始化异常退出日志函数
	initFatalLog func(logDirectory string)
	// 默认日志
	defaultLogger Logger
)

// Config
//
//	@Description: 日志配置定义
type Config struct {
	AppName             string          // 应用名
	AppId               string          // 应用编号
	Development         bool            // 开发模式(该模式下会使用更容易阅读的日志形式,但性能也更差)
	LogLevel            zapcore.Level   // 记录日志级别
	LogsDirectory       string          // 日志文件上层目录,一般该目录为默认目录的"AppName"下
	LogsFilenamePattern string          // 日志文件名样式
	LogsMaxAge          time.Duration   // 日志文件保存最大时间
	LogsRotationCount   uint            // 日志文件个数限制(与最大保存时间不同时生效)
	LogsRotationTime    time.Duration   // 日志文件分割周期
	LogsRotationSize    int64           // 日志文件分割大小(单位字节)
	ZapLogEncoder       zapcore.Encoder // zap日志编码器
}

// NewDevConfig
//
//	@Description: 构建开发模式默认配置
//	@param appName string
//	@param appId string
//	@return *Config
func NewDevConfig(appName, appId string) *Config {
	return &Config{
		AppName:     appName,
		AppId:       appId,
		Development: true,
		LogLevel:    zapcore.DebugLevel,
		//LogsDirectory:    "./logs/"+appName+"/",
		//LogsFilenamePattern: "%Y-%m-%d.plog",
		//LogsMaxAge:       time.Hour * 24 * 7,
		//LogsRotationSize: 50_000_000,
		ZapLogEncoder: newZapEncoder(true),
	}
}

// NewProConfig
//
//	@Description: 构建生产模式的默认配置
//	@param appName string
//	@param appId string
//	@return *Config
func NewProConfig(appName, appId string) *Config {
	return &Config{
		AppName:             appName,
		AppId:               appId,
		Development:         false,
		LogLevel:            zapcore.InfoLevel,
		LogsDirectory:       defaultLogDir + appName + "/",
		LogsFilenamePattern: "%Y-%m-%d.plog",
		LogsMaxAge:          time.Hour * 24 * 30,
		LogsRotationTime:    time.Hour * 24,
		ZapLogEncoder:       newZapEncoder(false),
	}
}

func newZapEncoder(development bool) zapcore.Encoder {
	var encoder zapcore.Encoder
	if development {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		encoderConfig.EncodeTime = colorTimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.ConsoleSeparator = consoleSeparator
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		encoderConfig.EncodeTime = timeEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	return encoder
}

// Init
//
//	@Description: 初始化默认日志系统，使用日志前需要先自行初始化
//	@param config *Config
//	@param opts ...zap.Option
func Init(config *Config, opts ...zap.Option) {
	err := defaultLogger.init(config, completeOptions(config, opts...)...)
	if err != nil {
		panic(err)
	}
	// 设置工具模块日志
	loggers.SetZapLogger(CoreLogger())
}

// completeOptions
//
//	@Description: 补全日志配置
//	@param config *Config
//	@param opts ...zap.Option
//	@return loggerOpts []zap.Option
func completeOptions(config *Config, opts ...zap.Option) (loggerOpts []zap.Option) {
	loggerOpts = append(loggerOpts, zap.AddCaller())
	// 自动添加开发选项
	if config.Development {
		loggerOpts = append(loggerOpts, zap.Development())
	}
	// 添加过滤调用层
	loggerOpts = append(loggerOpts, zap.AddCallerSkip(1))
	// 设置应用编号
	loggerOpts = append(loggerOpts, zap.Fields(zap.String("appId", fmt.Sprintf("%s_%s", config.AppName, config.AppId))))
	if len(opts) > 0 {
		loggerOpts = append(loggerOpts, opts...)
	}
	return
}

func init() {
	// 默认将日志初始化到 控制台 输出，使用前应先自行初始化默认日志系统
	Init(&Config{
		AppName:       "dummy",
		AppId:         "0",
		Development:   false,
		LogLevel:      zapcore.InfoLevel,
		ZapLogEncoder: newZapEncoder(false),
	})
}

func CoreLogger() *zap.Logger {
	return defaultLogger.inner
}

func LoggerLevel() zapcore.Level {
	return defaultLogger.LoggerLevel()
}
