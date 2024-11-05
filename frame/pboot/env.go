package pboot

import "fmt"

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

func IsEnvName(name string) bool {
	switch name {
	case envLocalName, envDevelopmentName, envTestName, envProductName:
		return true
	default:
		return false
	}
}

func ToEnv(envName string) (Env, error) {
	var env Env
	switch envName {
	case envLocalName:
		env = EnvLocal
	case envDevelopmentName:
		env = EnvDevelopment
	case envTestName:
		env = EnvTest
	case envProductName:
		env = EnvProduct
	default:
		env = EnvInvalid
		return env, fmt.Errorf("unknown environment:%s", envName)
	}
	return env, nil
}

func EnvName(env Env) (string, error) {
	switch env {
	case EnvLocal:
		return envLocalName, nil
	case EnvDevelopment:
		return envDevelopmentName, nil
	case EnvTest:
		return envTestName, nil
	case EnvProduct:
		return envProductName, nil
	default:
		return "", fmt.Errorf("unknown environment:%d", env)
	}
}
