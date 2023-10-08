package coding

// Cast
//
//	@Description: 类型转换
//	@param origin any 原始对象
//	@return value V
func Cast[V any](origin any) (value V) {
	if origin != nil {
		if cVal, ok := origin.(V); ok {
			value = cVal
		}
	}
	return
}
