package injection

import (
	"github.com/eden-quan/go-biz-kit/injection"
	servers "github.com/eden-quan/go-biz-kit/server"
)

/*
Inject 为微服务框架注入 Web/GRPC 的基础支持
 1. NewGRPCServer 依赖配置中心及日志库，为开发框架提供了 GRPC 基础支持，可通过 *grpc.Server 获取
 2. NewHTTPServer 依赖配置中心及日志库，为开发框架提供了 HTTP 基础支持，可通过 *http.Server 获取
 3. NewApp 根据配置中心信息，提供 GRPC 及 HTTP 服务，开发框架当前基于 Kratos 提供微服务基础能力，
    后续可通过替换该组件提供其他微服务框架作为底层支撑
*/
func Inject() {
	InjectIns(injection.GlobalInjector())
}

/*
InjectIns 使用创建实例的方式为微服务框架注入 Web/GRPC 的基础支持
 1. NewGRPCServer 依赖配置中心及日志库，为开发框架提供了 GRPC 基础支持，可通过 *grpc.Server 获取
 2. NewHTTPServer 依赖配置中心及日志库，为开发框架提供了 HTTP 基础支持，可通过 *http.Server 获取
 3. NewApp 根据配置中心信息，提供 GRPC 及 HTTP 服务，开发框架当前基于 Kratos 提供微服务基础能力，
    后续可通过替换该组件提供其他微服务框架作为底层支撑
*/
func InjectIns(inj *injection.Injector) {
	inj.InjectMany(
		servers.NewGRPCServer,
		servers.NewGrpcServiceRegistrar,
		servers.NewHTTPServer,
		servers.NewApp,
	)

	inj.Invoke(
		injection.WithInvoke(servers.StartKratosApp),
	)
}
