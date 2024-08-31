package servers

import (
	middlewareutil "github.com/eden/go-biz-kit/middleware"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	grpc2 "google.golang.org/grpc"

	"github.com/eden/go-biz-kit/config/def"
	"github.com/eden/go-biz-kit/injection"
	setup2 "github.com/eden/go-biz-kit/setup"

	apppkg "github.com/eden/go-kratos-pkg/app"
)

// NewGRPCServer 创建一个 gRPC 服务端
func NewGRPCServer(
	configure *def.Configuration,
	manager *setup2.LoggerManager,
	logger log.Logger,
	customMiddlewares *injection.MiddlewareCollector,
	actionManage *setup2.ActionManager,
) (srv *grpc.Server, err error) {
	helper := log.NewHelper(logger)

	if configure.Server.GetGrpc() == nil {
		helper.Warn("grpc configuration's is null ")
		return nil, nil
	}

	if !configure.Server.GetGrpc().GetEnable() {
		helper.Warn("grpc configuration's is disable ")
		return nil, nil
	}

	// options
	grpcConfig := configure.Server.GetGrpc()
	var opts []grpc.ServerOption

	if grpcConfig.GetAddress() != "" {
		opts = append(opts, grpc.Address(grpcConfig.GetAddress()))
	}

	// 默认超时时间 5 分钟
	timeOut := time.Minute * 5
	if grpcConfig.GetTimeout() != nil {
		timeOut = grpcConfig.GetTimeout().AsDuration()
	}

	opts = append(opts, grpc.Timeout(timeOut))

	var middlewareSlice = DefaultGrpcServerMiddlewares()
	middleLogger, err := manager.LoggerMiddleware()
	if err != nil {
		return srv, err
	}

	if customMiddlewares != nil && customMiddlewares.Middlewares != nil {
		// 批量添加自定义GRPC中间件
		for _, middle := range customMiddlewares.Middlewares {
			if middle.MidType == "GRPC" {
				middlewareSlice = append(middlewareSlice, middle.Middleware)
			}
		}
	}

	// 日志输出, SQL Action 处理器，确保在真正的业务逻辑执行之前触发
	middlewareSlice = append(middlewareSlice,
		apppkg.ServerLog(middleLogger),
		middlewareutil.SQLActionMiddleware(actionManage),
	)

	// 中间件选项
	opts = append(opts, grpc.Middleware(middlewareSlice...))

	// 服务
	srv = grpc.NewServer(opts...)

	return srv, err
}

// NewGrpcServiceRegistrar 提供原生的 grpc 注册接口，用于自动注册生成的 grpc 服务
func NewGrpcServiceRegistrar(gs *grpc.Server) grpc2.ServiceRegistrar {
	return gs.Server
}
