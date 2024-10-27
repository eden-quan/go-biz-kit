package injection

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"go.uber.org/fx"
)

type Injector struct {
	cache      injectCacheIns
	app        *fx.App
	middleware *MiddlewareCollector
	count      int

	bizCache map[interface{}]any
}

// NewInjector 创建独立的注入器，用于小范围的依赖注入管理, 如无需进行自定义作用范围的依赖管理，则应使用 Package 级别的注入器
func NewInjector() Injector {
	return Injector{
		cache:    injectCacheIns{},
		bizCache: map[interface{}]any{},
	}
}

var globalInjector Injector = Injector{
	cache:    injectCacheIns{},
	bizCache: make(map[interface{}]any),
}

func init() {
	globalInjector.cache.cache = make([]interface{}, 0)
}

type Option interface {
	Provide() fx.Option
}

// DoIt 按照已收集的依赖信息启动注入, 该函数会阻塞当前线程，如启动后需要执行其他功能，请使用 goroutine 执行
func (inj *Injector) DoIt(opts ...Option) {
	injectCache := inj.cache.injectCache()
	injectCache = append(injectCache,
		fx.Annotate(
			DefaultLoggerProvider,
			fx.ResultTags(`name:"logger"`),
		),
		fx.Annotate(
			collectServiceProvider,
			fx.ParamTags(`name:"logger"`, `group:"http_register"`, `group:"grpc_register"`),
		),
		fx.Annotate(
			inj.NewMiddlewareCollector,
		),
		fx.Annotate(
			inj.NewCollectorTrigger,
			fx.ParamTags(`group:"middleware"`),
		),
	)

	options := []fx.Option{
		fx.Provide(
			injectCache...,
		),
	}

	options = append(options, fx.Invoke(injectEntrance))

	for _, iv := range inj.cache.invokeCache {
		options = append(options, iv.Provide())
	}

	for _, o := range opts {
		options = append(options, o.Provide())
	}

	//options = append(options, fx.RecoverFromPanics())

	inj.app = fx.New(
		options...,
	)

	inj.app.Run()
}

type Middle struct {
	Middleware middleware.Middleware
	MidType    string
}

type MiddlewareCollector struct {
	Middlewares []Middle
}

func (inj *Injector) NewMiddlewareCollector(_ CollectorTrigger) *MiddlewareCollector {
	return inj.middleware
}

type CollectorTrigger struct{}

func (inj *Injector) NewCollectorTrigger(middle []Middle) CollectorTrigger {
	inj.middleware = &MiddlewareCollector{Middlewares: make([]Middle, 0)}
	inj.middleware.Middlewares = append(inj.middleware.Middlewares, middle...)
	return CollectorTrigger{}
}

func (inj *Injector) Replace(replacer interface{}) {
	inj.Invoke(WithReplace(replacer))
}
