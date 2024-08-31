package servers

import (
	middlewareutil "github.com/eden/go-biz-kit/middleware"
	"time"

	"github.com/go-kratos/kratos/v2/transport/http"

	apputil "github.com/eden/go-biz-kit/app"
	"github.com/eden/go-biz-kit/config/def"
	"github.com/eden/go-biz-kit/injection"
	setup2 "github.com/eden/go-biz-kit/setup"

	apppkg "github.com/eden/go-kratos-pkg/app"
)

// NewHTTPServer new HTTP server.
func NewHTTPServer(
	configuration *def.Configuration,
	manager *setup2.LoggerManager,
	customMiddlewares *injection.MiddlewareCollector,
	actionManage *setup2.ActionManager,
) (*http.Server, error) {

	if !configuration.Server.GetHttp().GetEnable() {
		return nil, nil
	}

	// options
	var opts []http.ServerOption

	httpConfig := configuration.Server.GetHttp()
	if httpConfig.GetAddress() != "" {
		opts = append(opts, http.Address(httpConfig.GetAddress()))
	}

	// 默认超时 5 分钟
	timeOut := time.Minute * 5
	if httpConfig.GetTimeout() != nil {
		timeOut = httpConfig.GetTimeout().AsDuration()
	}
	opts = append(opts, http.Timeout(timeOut))

	// 响应
	opts = append(opts, http.RequestDecoder(apputil.RequestDecoder))
	opts = append(opts, http.ErrorEncoder(apputil.ErrorEncoder))

	var middlewareSlice = DefaultServerMiddlewares()
	middleLogger, err := manager.LoggerMiddleware()

	if err != nil {
		return nil, err
	}

	// 日志输出, SQL Action 处理器，确保在真正的业务逻辑执行之前触发
	middlewareSlice = append(middlewareSlice,
		apppkg.ServerLog(middleLogger),
		middlewareutil.SQLActionMiddleware(actionManage),
	)

	if customMiddlewares != nil && customMiddlewares.Middlewares != nil {
		// 批量添加自定义HTTP中间件
		for _, middle := range customMiddlewares.Middlewares {
			if middle.MidType == "HTTP" {
				middlewareSlice = append(middlewareSlice, middle.Middleware)
			}
		}
	}

	// 中间件选项
	opts = append(opts, http.Middleware(middlewareSlice...))

	// 服务
	srv := http.NewServer(opts...)

	return srv, err
}
