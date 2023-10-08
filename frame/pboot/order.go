package pboot

const (
	OrderInternal = 1    // 内部优先级
	OrderConfig   = 500  // 配置优先级
	OrderDB       = 600  // 数据库优先级
	OrderTools    = 700  // 工具类优先级
	OrderCustom   = 1000 // 自定义
	OrderMax      = 5000 // 最大值
)
