//go:build linux
// +build linux

package plog

import (
	"os"
	"syscall"
)

func init() {
	initFatalLog = func(logDirectory string) {
		if len(logDirectory) <= 0 {
			return
		}
		logFile := buildFatalLogFile(logDirectory)
		if logFile != nil {
			syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
		}
	}
}
