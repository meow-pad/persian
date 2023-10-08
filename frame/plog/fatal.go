package plog

import (
	"fmt"
	"github.com/1set/gut/yos"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

// buildFatalLogFile
//
//	@Description: 初始化系统错误日志输出文件
//	@param directory string 日志文件目录
//	@return *os.File
func buildFatalLogFile(directory string) *os.File {
	maxLogSize := int64(50_000_000) // 50M
	flag := os.O_CREATE | os.O_APPEND | os.O_RDWR
	//flag := os.O_CREATE | os.O_RDWR
	fileName := filepath.Join(directory, "fatal.log")
	if logFile, err := os.OpenFile(fileName, flag, 0660); err != nil {
		Error("create fatal log error", zap.Error(err))
		return nil
	} else {
		if fileInfo, sErr := logFile.Stat(); sErr != nil {
			Error("get fatal log file error", zap.Error(sErr))
		} else {
			// 大小超预期
			if fileInfo.Size() > maxLogSize {
				// 拷贝到备份日志文件
				if cErr := yos.CopyFile(fileName, filepath.Join(directory, "fatal_bak.log")); cErr != nil {
					Error("copy fatal log file error", zap.Error(cErr))
				}
				// 清空当前日志
				if tErr := logFile.Truncate(0); tErr != nil {
					Error("truncate fatal log file error", zap.Error(tErr))
				}
			} // end of if
		} // end of else
		if defaultLogger.logCfg != nil {
			// 仅默认日志能写fatal.log
			_, _ = logFile.WriteString(fmt.Sprintf("%v serverName=%s,serverId=%s\n", time.Now(),
				defaultLogger.logCfg.AppName, defaultLogger.logCfg.AppId))
		}
		return logFile
	}
}

// AppFatal
//
//	@Description: 记录应用致命错误
//		出现应用启动时报错、log没有正确配置或log没有完全启动时，使用log无法记录下报错日志，此时可以使用该函数进行记录
//	@param msg string 错误消息
//	@param exit bool 是否退出
//	@param exitCode int 退出码
func AppFatal(msg string, exit bool, exitCode int) {
	_, _ = os.Stderr.WriteString(msg)
	if exit {
		os.Exit(exitCode)
	}
}
