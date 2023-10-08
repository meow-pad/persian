package loggers

import (
	"fmt"
	"os"
)

// StdLogger
//
//	@Description: 标准日志接口
type StdLogger interface {
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

func GetStdLogger() StdLogger {
	return stdLogger
}

func SetStdLogger(logger StdLogger) {
	stdLogger = logger
}

type fmtStdLogger struct{}

var stdLogger StdLogger = &fmtStdLogger{}

func (fLogger *fmtStdLogger) Print(args ...any) {
	fmt.Print(args...)
}

func (fLogger *fmtStdLogger) Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (fLogger *fmtStdLogger) Debug(args ...any)                 {}
func (fLogger *fmtStdLogger) Debugf(format string, args ...any) {}

func (fLogger *fmtStdLogger) Info(args ...any) {
	fmt.Print(args...)
}
func (fLogger *fmtStdLogger) Infof(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (fLogger *fmtStdLogger) Warn(args ...any) {
	fmt.Print(args...)
}
func (fLogger *fmtStdLogger) Warnf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (fLogger *fmtStdLogger) Error(args ...any) {
	fmt.Print(args...)
}
func (fLogger *fmtStdLogger) Errorf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (fLogger *fmtStdLogger) Fatal(args ...any) {
	fmt.Print(args...)
	os.Exit(1)
}
func (fLogger *fmtStdLogger) Fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func (fLogger *fmtStdLogger) Panic(args ...any) {
	panic(fmt.Sprint(args...))
}

func (fLogger *fmtStdLogger) Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}
