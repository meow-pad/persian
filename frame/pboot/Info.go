package pboot

import (
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
