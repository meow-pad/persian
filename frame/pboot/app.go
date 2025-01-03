package pboot

import (
	"github.com/go-spring/spring-core/gs"
	"github.com/go-spring/spring-core/gs/arg"
	"github.com/go-spring/spring-core/web"
	"net/http"
)

var (
	app *gs.App
)

func initApp() {
	app = gs.NewApp()
}

// Run 启动程序。
func Run() error {
	return app.Run()
}

// RunWithWeb 带web启动。
func RunWithWeb() error {
	gs.Object(new(gs.WebStarter)).Export((*gs.AppEvent)(nil))
	return app.Run()
}

// ShutDown 停止程序。
func ShutDown(msg ...string) {
	app.ShutDown(msg...)
}

// Banner 参考 App.Banner 的解释。
func Banner(banner string) {
	app.Banner(banner)
}

// OnProperty 参考 App.OnProperty 的解释。
func OnProperty(key string, fn any) {
	app.OnProperty(key, fn)
}

// Property 参考 Container.Property 的解释。
func Property(key string, value any) {
	app.Property(key, value)
}

func InternalObject(i any) *Bean {
	return setupLifeCycleModule(app.Object(i), OrderInternal)
}

func InternalProvide(ctor any, args ...arg.Arg) *Bean {
	return setupLifeCycleModule(app.Provide(ctor, args...), OrderInternal)
}

func ConfigObject(i any) *Bean {
	return setupLifeCycleModule(app.Object(i), OrderConfig)
}

func ConfigProvide(ctor any, args ...arg.Arg) *Bean {
	return setupLifeCycleModule(app.Provide(ctor, args...), OrderConfig)
}

func DBObject(i any) *Bean {
	return setupLifeCycleModule(app.Object(i), OrderDB)
}

func DBProvide(ctor any, args ...arg.Arg) *Bean {
	return setupLifeCycleModule(app.Provide(ctor, args...), OrderDB)
}

func ToolsObject(i any) *Bean {
	return setupLifeCycleModule(app.Object(i), OrderTools)
}

func ToolsProvide(ctor any, args ...arg.Arg) *Bean {
	return setupLifeCycleModule(app.Provide(ctor, args...), OrderTools)
}

func BaseObject(i any) *Bean {
	return setupLifeCycleModule(app.Object(i), OrderCustomBase)
}

func BaseProvide(ctor any, args ...arg.Arg) *Bean {
	return setupLifeCycleModule(app.Provide(ctor, args...), OrderCustomBase)
}

// Object 参考 Container.Object 的解释。
func Object(i any) *Bean {
	return setupLifeCycleModule(app.Object(i), OrderCustom)
}

func LastObject(i any) *Bean {
	return setupLifeCycleModule(app.Object(i), OrderCustom).Order(OrderMax)
}

// Provide 参考 Container.Provide 的解释。
func Provide(ctor any, args ...arg.Arg) *Bean {
	return setupLifeCycleModule(app.Provide(ctor, args...), OrderCustom)
}

func setupLifeCycleModule(beanDef *gs.BeanDefinition, baseOrder float32) *Bean {
	bean := newBean(beanDef, baseOrder).Order(1)
	lifeCycle, _ := beanDef.Interface().(LifeCycle)
	if lifeCycle != nil {
		addLifeCycle(lifeCycle, bean)
	}
	startListener, _ := beanDef.Interface().(StartListener)
	if startListener != nil {
		addStartListener(startListener, bean)
	}
	return bean
}

// HandleGet 参考 App.HandleGet 的解释。
func HandleGet(path string, h web.Handler) *web.Mapper {
	return app.HandleGet(path, h)
}

// GetMapping 参考 App.GetMapping 的解释。
func GetMapping(path string, fn web.HandlerFunc) *web.Mapper {
	return app.GetMapping(path, fn)
}

// GetBinding 参考 App.GetBinding 的解释。
func GetBinding(path string, fn any) *web.Mapper {
	return app.GetBinding(path, fn)
}

// HandlePost 参考 App.HandlePost 的解释。
func HandlePost(path string, h web.Handler) *web.Mapper {
	return app.HandlePost(path, h)
}

// PostMapping 参考 App.PostMapping 的解释。
func PostMapping(path string, fn web.HandlerFunc) *web.Mapper {
	return app.PostMapping(path, fn)
}

// PostBinding 参考 App.PostBinding 的解释。
func PostBinding(path string, fn any) *web.Mapper {
	return app.PostBinding(path, fn)
}

// HandlePut 参考 App.HandlePut 的解释。
func HandlePut(path string, h web.Handler) *web.Mapper {
	return app.HandlePut(path, h)
}

// PutMapping 参考 App.PutMapping 的解释。
func PutMapping(path string, fn web.HandlerFunc) *web.Mapper {
	return app.PutMapping(path, fn)
}

// PutBinding 参考 App.PutBinding 的解释。
func PutBinding(path string, fn any) *web.Mapper {
	return app.PutBinding(path, fn)
}

// HandleDelete 参考 App.HandleDelete 的解释。
func HandleDelete(path string, h web.Handler) *web.Mapper {
	return app.HandleDelete(path, h)
}

// DeleteMapping 参考 App.DeleteMapping 的解释。
func DeleteMapping(path string, fn web.HandlerFunc) *web.Mapper {
	return app.DeleteMapping(path, fn)
}

// DeleteBinding 参考 App.DeleteBinding 的解释。
func DeleteBinding(path string, fn any) *web.Mapper {
	return app.DeleteBinding(path, fn)
}

// HandleRequest 参考 App.HandleRequest 的解释。
func HandleRequest(method uint32, path string, h web.Handler) *web.Mapper {
	return app.HandleRequest(method, path, h)
}

// RequestMapping 参考 App.RequestMapping 的解释。
func RequestMapping(method uint32, path string, fn web.HandlerFunc) *web.Mapper {
	return app.RequestMapping(method, path, fn)
}

// RequestBinding 参考 App.RequestBinding 的解释。
func RequestBinding(method uint32, path string, fn any) *web.Mapper {
	return app.RequestBinding(method, path, fn)
}

// File 定义单个文件资源
func File(path string, file string) *web.Mapper {
	return app.File(path, file)
}

// Static 定义一组文件资源
func Static(prefix string, dir string) *web.Mapper {
	return app.Static(prefix, dir)
}

// StaticFS 定义一组文件资源
func StaticFS(prefix string, fs http.FileSystem) *web.Mapper {
	return app.StaticFS(prefix, fs)
}

// Consume 参考 App.Consume 的解释。
func Consume(fn any, topics ...string) {
	app.Consume(fn, topics...)
}
