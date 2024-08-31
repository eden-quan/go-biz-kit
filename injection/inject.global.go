package injection

import (
	"go.uber.org/fx"
)

// In 用于实现将依赖的参数集中到 struct 的能力，防止依赖太多时参数列表难以维护
type In = fx.In

// DoIt 按照已收集的依赖信息启动注入
func DoIt(opts ...Option) {
	globalInjector.DoIt(opts...)
}

// Inject 接收任意类型的函数，每个函数通过参数及返回值来告知依赖注入引擎这个函数需要什么信息(参数)，以及能够提供什么信息(返回值)
// 注入后的信息会在最终调用 DoIt 时进行解析
func Inject(anno interface{}) {
	globalInjector.cache.Inject(anno)
}

type InjectClientProvider = func(interface{}) []interface{}

// InjectGRPCClient 注入 protocol buffer 生成的 GRPC 或 HTTP 客户端
func InjectGRPCClient(provider interface{}) {
	globalInjector.InjectGRPCClient(provider)
}

// InjectHTTPClient 注入 protocol buffer 生成的 GRPC 或 HTTP 客户端
func InjectHTTPClient(provider interface{}) {
	globalInjector.InjectHTTPClient(provider)
}

func InjectGRPCMiddleware(mid interface{}) {
	globalInjector.InjectGRPCMiddleware(mid)
}

func InjectHTTPMiddleware(mid interface{}) {
	globalInjector.InjectHTTPMiddleware(mid)
}

// InjectAs 注入信息为自定义类型, 如存在以下两个 interface
// interface LoggerA 和 LoggerB
// InjectAs(func() LoggerA { return logger }, new(LoggerB)}
func InjectAs(anno interface{}, t interface{}) {
	globalInjector.cache.InjectAs(anno, t)
}

// InjectMany 批量添加注入函数到依赖注入引擎中, 具体定义见 Inject 定义
func InjectMany(anno ...interface{}) {
	globalInjector.cache.InjectMany(anno...)
}

// Invoke 用于注册依赖注入容器在解析完依赖之后，启动容器时执行的函数，
// WARN: 注册的函数除非明确执行的含义，当函数阻塞时会同时阻塞当前 Goroutine
func Invoke(opt Option) {
	globalInjector.cache.Invoke(opt)
}

//
//// ReplaceConfig 便捷函数，允许用户快速替换全局配置信息
//func ReplaceConfig(configuration *def.Configuration) {
//	globalInjector.ReplaceConfig(configuration)
//}

// Replace 将为一个函数，他可以接收其他依赖注入提供的信息，并且会使用其返回值替换原有的依赖注入类型
func Replace(replacer interface{}) {
	globalInjector.Replace(replacer)
}

// InjectWithParam 提供了注入依赖的时候添加 Param 及 Result 标签的入口，通过标签可以将注入的依赖提供给指定的使用者
// 如 params 指定本依赖需要对应名称的标签参数，而 result 指定本注入的返回值将作为指定标签的返回值
// example:
// InjectWithParam(action1, []string{`name: "need_arg1"`}, nil) // 说明 action1 需要标签为 need_arg1 标签的参数
// InjectWithParam(action2, nil, []string{`name: "need_arg1"`}) // 说明 action2 的返回结果将作为为 action1 的参数
// 标签的类型有 name / group 两种，name 为单一对象参数，group 为数组对象参数
func InjectWithParam(anno interface{}, params []string, results []string) {
	globalInjector.InjectWithParam(anno, params, results)
}

func GlobalInjector() *Injector {
	return &globalInjector
}
