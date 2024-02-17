package pboot

import (
	"fmt"
	"github.com/meow-pad/persian/errdef"
)

const (
	OrderInternal   = 1     // 内部优先级
	OrderConfig     = 500   // 配置优先级
	OrderDB         = 600   // 数据库优先级
	OrderTools      = 700   // 工具类优先级
	OrderCustomBase = 1000  // 自定义基础优先级
	OrderCustom     = 2000  // 自定义
	OrderMax        = 10000 // 最大值
)

var baseOrders = []float32{OrderInternal, OrderConfig, OrderDB, OrderTools, OrderCustomBase, OrderCustom}

func getOrder(baseOrder float32, value float32) (float32, error) {
	if value < 0 {
		return 0, errdef.ErrInvalidParams
	}
	boIndex := -1
	for i, bo := range baseOrders {
		if bo == baseOrder {
			boIndex = i
			break
		}
	}
	if boIndex < 0 {
		return 0, fmt.Errorf("invalid base order:%f", baseOrder)
	}
	max := float32(0)
	if boIndex >= len(baseOrders)-1 {
		max = OrderMax
	} else {
		max = baseOrders[boIndex+1] - 1
	}
	order := baseOrder + value
	if order > max {
		return max, nil
	}
	return order, nil
}
