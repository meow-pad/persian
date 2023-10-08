package pruntime

import (
	_ "go.uber.org/automaxprocs"
	"runtime"
)

// MaxProcess
//
//	@Description: 获取最大可用核心数
//		该函数支持docker环境
//	@return int
func MaxProcess() int {
	return runtime.GOMAXPROCS(0)
}

// NumGoroutine
//
//	@Description: 返回当前存在的goroutines数量
//	@return int
func NumGoroutine() int {
	return runtime.NumGoroutine()
}
