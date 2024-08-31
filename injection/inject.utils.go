package injection

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	clientutil "github.com/eden-quan/go-biz-kit/client"
)

// HttpProvider 提供便捷的入口简化用户创建匿名函数来提供 Http 服务
func HttpProvider(server *http.Server) func() *http.Server {
	return func() *http.Server {
		return server
	}
}

// GrpcProvider 提供便捷的入口简化用户创建匿名函数来提供 GRPC 服务
func GrpcProvider(server grpc.ServiceRegistrar) func() grpc.ServiceRegistrar {
	return func() grpc.ServiceRegistrar {
		return server
	}
}

// LoggerProvider 提供便捷的入口简化用户创建匿名函数来提供 Logger
func LoggerProvider(logger log.Logger) func() log.Logger {
	return func() log.Logger {
		return logger
	}
}

// DefaultLoggerProvider 提供默认日志
func DefaultLoggerProvider() log.Logger {
	return log.DefaultLogger
}

func (inj *Injector) InjectHttpServer(server *http.Server) {
	inj.cache.Inject(HttpProvider(server))
}

func (inj *Injector) InjectGrpcServer(server grpc.ServiceRegistrar) {
	inj.cache.Inject(GrpcProvider(server))
}

func (inj *Injector) InjectLogger(logger log.Logger) {
	inj.cache.Inject(LoggerProvider(logger))
}

func (inj *Injector) Inject(anno interface{}) {
	inj.cache.Inject(anno)
}

func (inj *Injector) InjectMany(anno ...interface{}) {
	inj.cache.InjectMany(anno...)
}

func (inj *Injector) InjectAs(anno interface{}, t interface{}) {
	inj.cache.InjectAs(anno, t)
}

func (inj *Injector) Invoke(opt Option) {
	inj.cache.Invoke(opt)
}

//
//func (inj *Injector) ReplaceConfig(configuration *def.Configuration) {
//	inj.Invoke(
//		WithReplace(
//			func(repo config.ConfigureWatcherRepo) (*def.Configuration, error) {
//				err := repo.LoadAndStart(configuration)
//				if err != nil {
//					err = fmt.Errorf("[Inject] replacing service configuration with error %s", err)
//				}
//
//				return configuration, err
//			}),
//	)
//}

func (inj *Injector) InjectHTTPMiddleware(mid interface{}) {
	inj.count += 1
	midTempTag := fmt.Sprintf(`name:"middleware-%d"`, inj.count)

	inj.Inject(fx.Annotate(
		mid,
		fx.ResultTags(midTempTag),
	))
	inj.Inject(
		fx.Annotate(
			func(m middleware.Middleware) Middle {
				return Middle{
					Middleware: m,
					MidType:    "HTTP",
				}
			},
			fx.ParamTags(midTempTag),
			fx.ResultTags(`group:"middleware"`),
		),
	)
}

// InjectWithParam 提供了注入依赖的时候添加 Param 及 Result 标签的入口，通过标签可以将注入的依赖提供给指定的使用者
// 如 params 指定本依赖需要对应名称的标签参数，而 result 指定本注入的返回值将作为指定标签的返回值
// example:
// InjectWithParam(action1, []string{`name: "need_arg1"`}, nil) // 说明 action1 需要标签为 need_arg1 标签的参数
// InjectWithParam(action2, nil, []string{`name: "need_arg1"`}) // 说明 action2 的返回结果将作为为 action1 的参数
// 标签的类型有 name / group 两种，name 为单一对象参数，group 为数组对象参数
func (inj *Injector) InjectWithParam(anno interface{}, params []string, results []string) {
	var annotates []fx.Annotation = make([]fx.Annotation, 0)

	if len(params) != 0 {
		annotates = append(annotates, fx.ParamTags(params...))
	}

	if len(results) != 0 {
		annotates = append(annotates, fx.ResultTags(results...))
	}

	inj.Inject(fx.Annotate(anno, annotates...))
}

func (inj *Injector) InjectGRPCMiddleware(mid interface{}) {
	inj.count += 1
	midTempTag := fmt.Sprintf(`name:"middleware-%d"`, inj.count)

	inj.Inject(fx.Annotate(
		mid,
		fx.ResultTags(midTempTag),
	))
	inj.Inject(
		fx.Annotate(
			func(m middleware.Middleware) Middle {
				return Middle{
					Middleware: m,
					MidType:    "GRPC",
				}
			},
			fx.ParamTags(midTempTag),
			fx.ResultTags(`group:"middleware"`),
		),
	)
}

// InjectGRPCClient 注入 protocol buffer 生成的 GRPC 或 HTTP 客户端
func (inj *Injector) InjectGRPCClient(provider interface{}) {

	p := unsafe.Pointer(reflect.ValueOf(provider).Pointer())

	_, exists := inj.bizCache[p]
	if exists {
		return
	}

	pro := provider.(InjectClientProvider)
	inj.cache.Inject(
		pro(clientutil.RegisterGRPCClient),
	)
	inj.bizCache[p] = true
}

func (inj *Injector) InjectHTTPClient(provider interface{}) {

	p := unsafe.Pointer(reflect.ValueOf(provider).Pointer())

	_, exists := inj.bizCache[p]
	if exists {
		return
	}

	pro := provider.(InjectClientProvider)
	inj.cache.Inject(
		pro(clientutil.RegisterHTTPClient),
	)
	inj.bizCache[p] = true
}
