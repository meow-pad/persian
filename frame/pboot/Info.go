package pboot

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"time"
)

type AppInfo interface {
	// Id
	//  @Description: 应用编号
	//  @return string
	//
	Id() string

	// Name
	//  @Description: 应用名
	//  @return string
	//
	Name() string

	// Env
	//  @Description: 运行环境
	//  @return Env
	//
	Env() Env

	// EnvName
	//  @Description: 运行环境名
	//  @return string
	//
	EnvName() string

	// Cluster
	//
	//	@Description: 获取服务集群
	//	@return string
	Cluster() string

	// TimeZone
	//
	//	@Description: 当前服务时区
	//	@return *time.Location
	TimeZone() *time.Location
}

type BaseAppInfo struct {
	id      string
	name    string
	env     Env
	envName string
	cluster string
}

func (info *BaseAppInfo) Init() error {
	// 获取环境数据
	info.envName = os.Getenv(osEnvNameEnvironment)
	switch info.envName {
	case envLocalName:
		info.env = EnvLocal
	case envDevelopmentName:
		info.env = EnvDevelopment
	case envTestName:
		info.env = EnvTest
	case envProductName:
		info.env = EnvProduct
	default:
		info.env = EnvInvalid
		return fmt.Errorf("unknown running environment:%s", info.envName)
	}
	// 获取集群
	info.cluster = os.Getenv(osEnvNameCluster)
	// 获取服务器时区
	serverTZName := os.Getenv(osEnvNameTimeZone)
	if len(serverTZName) > 0 {
		var err error
		serverTimeZone, err := time.LoadLocation(serverTZName)
		if err != nil {
			return errors.WithMessage(err, "invalid timezone:"+serverTZName)
		}
		time.Local = serverTimeZone
	}
	return nil
}

func (info *BaseAppInfo) Id() string {
	return info.id
}

func (info *BaseAppInfo) Name() string {
	return info.name
}

func (info *BaseAppInfo) Env() Env {
	return info.env
}

func (info *BaseAppInfo) EnvName() string {
	return info.envName
}

func (info *BaseAppInfo) Cluster() string {
	return info.cluster
}

func (info *BaseAppInfo) TimeZone() *time.Location {
	return time.Local
}
