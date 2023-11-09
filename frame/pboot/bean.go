package pboot

import (
	"github.com/go-spring/spring-base/util"
	"github.com/go-spring/spring-core/gs"
	"github.com/go-spring/spring-core/gs/cond"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"reflect"
)

func newBean(inner *gs.BeanDefinition, baseOrder float32) *Bean {
	return &Bean{
		inner:     inner,
		baseOrder: baseOrder,
	}
}

type Bean struct {
	inner     *gs.BeanDefinition
	baseOrder float32
}

// Type 返回 bean 的类型。
func (d *Bean) Type() reflect.Type {
	return d.inner.Type()
}

// Value 返回 bean 的值。
func (d *Bean) Value() reflect.Value {
	return d.inner.Value()
}

// Interface 返回 bean 的真实值。
func (d *Bean) Interface() interface{} {
	return d.inner.Interface()
}

// ID 返回 bean 的 ID 。
func (d *Bean) ID() string {
	return d.inner.ID()
}

// BeanName 返回 bean 的名称。
func (d *Bean) BeanName() string {
	return d.inner.BeanName()
}

// TypeName 返回 bean 的原始类型的全限定名。
func (d *Bean) TypeName() string {
	return d.inner.TypeName()
}

// Created 返回是否已创建。
func (d *Bean) Created() bool {
	return d.inner.Created()
}

// Wired 返回 bean 是否已经注入。
func (d *Bean) Wired() bool {
	return d.Wired()
}

// FileLine 返回 bean 的注册点。
func (d *Bean) FileLine() string {
	return d.inner.FileLine()
}

func (d *Bean) String() string {
	return d.inner.String()
}

// Match 测试 bean 的类型全限定名和 bean 的名称是否都匹配。
func (d *Bean) Match(typeName string, beanName string) bool {
	return d.inner.Match(typeName, beanName)
}

// Name 设置 bean 的名称。
func (d *Bean) Name(name string) *Bean {
	d.inner.Name(name)
	return d
}

// On 设置 bean 的 Condition。
func (d *Bean) On(cond cond.Condition) *Bean {
	d.inner.On(cond)
	return d
}

// Order 设置 bean 的排序序号，值越小顺序越靠前(优先级越高)。
func (d *Bean) Order(order float32) *Bean {
	bOrder, err := getOrder(d.baseOrder, order)
	if err != nil {
		plog.Panic("build order error:", pfield.Error(err))
	}
	d.inner.Order(bOrder)
	return d
}

// DependsOn 设置 bean 的间接依赖项。
func (d *Bean) DependsOn(selectors ...util.BeanSelector) *Bean {
	d.inner.DependsOn(selectors...)
	return d
}

// Primary 设置 bean 为主版本。
func (d *Bean) Primary() *Bean {
	d.inner.Primary()
	return d
}

// Init 设置 bean 的初始化函数。
func (d *Bean) Init(fn interface{}) *Bean {
	d.inner.Init(fn)
	return d
}

// Destroy 设置 bean 的销毁函数。
func (d *Bean) Destroy(fn interface{}) *Bean {
	d.inner.Destroy(fn)
	return d
}

// Export 设置 bean 的导出接口。
func (d *Bean) Export(exports ...interface{}) *Bean {
	d.inner.Export(exports...)
	return d
}
