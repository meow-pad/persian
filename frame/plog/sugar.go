package plog

import (
	"fmt"
	"go.uber.org/zap"
)

type GeneralLogger interface {
	Print(args ...any)
	Printf(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Panic(args ...any)
	Panicf(format string, args ...any)
}

func SugarLogger() GeneralLogger {
	return &defaultLogger.sugar
}

type sugarLogger struct {
	core *zap.Logger
}

func (sLogger *sugarLogger) Print(args ...any) {
	if sLogger.core != nil {
		sLogger.core.Info(fmt.Sprint(args...))
	}
}

func (sLogger *sugarLogger) Printf(format string, args ...any) {
	if sLogger.core != nil {
		sLogger.core.Info(fmt.Sprintf(format, args...))
	}
}

func (sLogger *sugarLogger) Debug(args ...any) {
	if sLogger.core != nil {
		sLogger.core.Debug(fmt.Sprint(args...))
	} // end of if
}

func (sLogger *sugarLogger) Debugf(format string, args ...any) {
	if sLogger.core != nil {
		sLogger.core.Debug(fmt.Sprintf(format, args...))
	} // end of if
}

func (sLogger *sugarLogger) Info(args ...any) {
	if sLogger.core != nil {
		sLogger.core.Info(fmt.Sprint(args...))
	} // end of if
}

func (sLogger *sugarLogger) Infof(format string, args ...any) {
	if sLogger.core != nil {
		sLogger.core.Info(fmt.Sprintf(format, args...))
	} // end of if
}

func (sLogger *sugarLogger) Warn(args ...any) {
	if sLogger.core != nil {
		sLogger.core.Warn(fmt.Sprint(args...))
	} // end of if
}

func (sLogger *sugarLogger) Warnf(format string, args ...any) {
	if sLogger.core != nil {
		sLogger.core.Warn(fmt.Sprintf(format, args...))
	} // end of if
}

func (sLogger *sugarLogger) Error(args ...any) {
	if sLogger.core != nil {
		sLogger.core.Error(fmt.Sprint(args...))
	} // end of if
}

func (sLogger *sugarLogger) Errorf(format string, args ...any) {
	if sLogger.core != nil {
		sLogger.core.Error(fmt.Sprintf(format, args...))
	} // end of if
}

func (sLogger *sugarLogger) Fatal(args ...any) {
	if sLogger.core != nil {
		sLogger.core.Fatal(fmt.Sprint(args...))
	} // end of if
}

func (sLogger *sugarLogger) Fatalf(format string, args ...any) {
	if sLogger.core != nil {
		sLogger.core.Fatal(fmt.Sprintf(format, args...))
	} // end of if
}

func (sLogger *sugarLogger) Panic(args ...any) {
	if sLogger.core != nil {
		sLogger.core.Panic(fmt.Sprint(args...))
	} // end of if
}

func (sLogger *sugarLogger) Panicf(format string, args ...any) {
	if sLogger.core != nil {
		sLogger.core.Panic(fmt.Sprintf(format, args...))
	} // end of if
}
