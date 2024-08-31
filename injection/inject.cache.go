package injection

import (
	"go.uber.org/fx"
)

type injectCacheIns struct {
	cache       []interface{}
	invokeCache []Option
}

func newInjectCacheIns() *injectCacheIns {
	return &injectCacheIns{
		cache:       make([]interface{}, 0),
		invokeCache: make([]Option, 0),
	}
}

func (cache *injectCacheIns) AddCache(anno interface{}) {
	cache.cache = append(cache.cache, anno)
}

func (cache *injectCacheIns) Caches() []interface{} {
	return cache.cache
}

// Inject 接收任意类型的函数，每个函数通过参数及返回值来告知依赖注入引擎这个函数需要什么信息(参数)，以及能够提供什么信息(返回值)
// 注入后的信息会在最终调用 DoIt 时进行解析
func (cache *injectCacheIns) Inject(anno interface{}) {
	switch a := anno.(type) {
	case []interface{}:
		for _, i := range a {
			cache.Inject(i)
		}
	default:
		cache.AddCache(anno)
	}

}

// InjectAs 注入信息为自定义类型, 如存在以下两个 interface
// interface LoggerA 和 LoggerB
// InjectAs(func() LoggerA { return logger }, new(LoggerB)}
func (cache *injectCacheIns) InjectAs(anno interface{}, t interface{}) {
	cache.Inject(fx.Annotate(
		anno,
		fx.As(t),
	))
}

// InjectMany 批量添加注入函数到依赖注入引擎中, 具体定义见 Inject 定义
func (cache *injectCacheIns) InjectMany(anno ...interface{}) {
	for _, a := range anno {
		cache.Inject(a)
	}
}

// InjectCache 获取所有需要注入的信息, 该接口一般只提供给依赖注入内部使用
func (cache *injectCacheIns) injectCache() []interface{} {
	return cache.cache
}

// Invoke 用于注册依赖注入解析完成后的启动函数, 注册后的函数会在解析完成后启动
func (cache *injectCacheIns) Invoke(option Option) {
	cache.invokeCache = append(cache.invokeCache, option)
}
