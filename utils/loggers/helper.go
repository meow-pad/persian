package loggers

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

func Debug(msg string, fields ...zap.Field) {
	if zapLogger != nil {
		zapLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if zapLogger != nil {
		zapLogger.Info(msg, fields...)
	} else {
		fmt.Print(msg)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if zapLogger != nil {
		zapLogger.Warn(msg, fields...)
	} else {
		fmt.Print(msg)
	}
}

func Error(msg string, fields ...zap.Field) {
	if zapLogger != nil {
		zapLogger.Error(msg, fields...)
	} else {
		fmt.Print(msg)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if zapLogger != nil {
		zapLogger.Fatal(msg, fields...)
	} else {
		fmt.Print(msg)
		os.Exit(1)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if zapLogger != nil {
		zapLogger.Panic(msg, fields...)
	} else {
		panic(msg)
	}
}
