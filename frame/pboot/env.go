package pboot

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
