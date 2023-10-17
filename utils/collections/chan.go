package collections

// ClearChan
//
//	@Description: 清空 chan 内的数据，如果chan已关闭则直接返回
//	@param channel
func ClearChan[T any](channel chan T) {
outer:
	for {
		select {
		case _, ok := <-channel:
			if !ok {
				break outer
			}
		default:
			break outer
		}
	}
}
