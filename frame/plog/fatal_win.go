//go:build windows

package plog

import "os"

func init() {
	initFatalLog = func(logDirectory string) {
		if len(logDirectory) <= 0 {
			return
		}
		logFile := buildFatalLogFile(logDirectory)
		if logFile != nil {
			os.Stderr = logFile
		}
	}
}
