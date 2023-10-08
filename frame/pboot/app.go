package pboot

import (
	"fmt"
	"github.com/go-spring/spring-base/util"
	"github.com/go-spring/spring-core/gs"
	"github.com/go-spring/spring-core/gs/arg"
	"github.com/go-spring/spring-core/web"
	"github.com/meow-pad/persian/frame/passert"
	"net/http"
	"os"
	"reflect"
)

var (
	app       *gs.App
	container gs.Container
)

func initApp() {
	app = gs.NewApp()
	vApp := reflect.ValueOf(app)
	fmt.Println(vApp.FieldByName("c"))
	container = vApp.FieldByName("c").Interface().(gs.Container)
	passert.NotNil(container, "empty app container")
}

// Setenv 封装 os.Setenv 函数，如果发生 error 会 panic 。
func Setenv(key string, value string) {
	err := os.Setenv(key, value)
	util.Panic(err).When(err != nil)
}

type startup struct {
	web bool
}

func webStartup(enable bool) *startup {
	return &startup{web: enable}
}

func (s *startup) Run() error {
	if s.web {
		Object(new(gs.WebStarter)).Export((*gs.AppEvent)(nil))
	}
	return app.Run()
}

// Run 启动程序。
func Run() error {
	return webStartup(true).Run()
}

// RunWithWeb 带web启动。
func RunWithWeb() error {
	return webStartup(true).Run()
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

// Object 参考 Container.Object 的解释。
func Object(i any) *gs.BeanDefinition {
	return setupLifeCycleModule(container.Object(i))
}

// Provide 参考 Container.Provide 的解释。
func Provide(ctor any, args ...arg.Arg) *gs.BeanDefinition {
	return setupLifeCycleModule(container.Provide(ctor, args...))
}

func setupLifeCycleModule(bean *gs.BeanDefinition) *gs.BeanDefinition {
	lifeCycle := bean.Interface().(LifeCycle)
	if lifeCycle != nil {
		// 加入到Event中进行排序
		bean.Export((*gs.AppEvent)(nil))
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
