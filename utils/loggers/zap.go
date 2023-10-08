package loggers

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

// ZapLogger
//
//	@Description: zap日志
type ZapLogger interface {
	Debug(msg string, fields ...zap.Field)

	Info(msg string, fields ...zap.Field)

	Warn(msg string, fields ...zap.Field)

	Error(msg string, fields ...zap.Field)

	Fatal(msg string, fields ...zap.Field)

	Panic(msg string, fields ...zap.Field)
}

func GetZapLogger() ZapLogger {
	return zapLogger
}

func SetZapLogger(logger ZapLogger) {
	zapLogger = logger
}

var zapLogger ZapLogger = &fmtZapLogger{}

type fmtZapLogger struct {
}

func (logger *fmtZapLogger) Debug(msg string, fields ...zap.Field) {}

func (logger *fmtZapLogger) Info(msg string, fields ...zap.Field) {
	fmt.Print(msg)
}

func (logger *fmtZapLogger) Warn(msg string, fields ...zap.Field) {
	fmt.Print(msg)
}

func (logger *fmtZapLogger) Error(msg string, fields ...zap.Field) {
	fmt.Print(msg)
}

func (logger *fmtZapLogger) Fatal(msg string, fields ...zap.Field) {
	fmt.Print(msg)
	os.Exit(1)
}

func (logger *fmtZapLogger) Panic(msg string, fields ...zap.Field) {
	panic(msg)
}
