package pruntime

import (
	"go.uber.org/zap"
	"os"
	"persian/frame/plog"
	"time"
)

type Env int32

const (
	EnvInvalid Env = 0
	// EnvLocal 本地开发环境
	EnvLocal Env = 1
	// EnvDevelopment 开发测试服
	EnvDevelopment Env = 2
	// EnvTest 外网测试服
	EnvTest Env = 3
	// EnvProduct 生产环境
	EnvProduct Env = 4
)

const (
	envLocalName       = "local"
	envDevelopmentName = "dev"
	envTestName        = "test"
	envProductName     = "pro"
)

// 环境变量名
const (
	osEnvNameEnvironment = "SRV_ENVIRONMENT"
	osEnvNameCluster     = "SRV_CLUSTER"
	osEnvNameTimeZone    = "SRV_TIMEZONE"
)

var (
	serverEnv     Env
	serverCluster string
	hostname      string
)

func init() {
	// 获取环境数据
	envName := os.Getenv(osEnvNameEnvironment)
	switch envName {
	case envLocalName:
		serverEnv = EnvLocal
	case envDevelopmentName:
		serverEnv = EnvDevelopment
	case envTestName:
		serverEnv = EnvTest
	case envProductName:
		serverEnv = EnvProduct
	default:
		serverEnv = EnvInvalid
		plog.Fatal("unknown running environment", zap.String(osEnvNameEnvironment, envName))
	}
	// 获取集群
	serverCluster = os.Getenv(osEnvNameCluster)
	// 获取服务器时区
	serverTZName := os.Getenv(osEnvNameTimeZone)
	if len(serverTZName) > 0 {
		var err error
		serverTimeZone, err := time.LoadLocation(serverTZName)
		if err != nil {
			plog.Fatal("invalid server timezone", zap.Error(err), zap.String(osEnvNameTimeZone, serverTZName))
		}
		time.Local = serverTimeZone
	}
	var err error
	// 获取主机名
	hostname, err = os.Hostname()
	if err != nil {
		plog.Fatal("get hostname error:", zap.Error(err))
	}
}

// Environment
//
//	@Description: 获取服务运行环境
//	@return Env
func Environment() Env {
	return serverEnv
}

// Cluster
//
//	@Description: 获取服务集群
//	@return string
func Cluster() string {
	return serverCluster
}

// Hostname
//
//	@Description: 获取主机名
//	@return string
func Hostname() string {
	return hostname
}

// TimeZone
//
//	@Description: 当前服务时区
//	@return *time.Location
func TimeZone() *time.Location {
	return time.Local
}
