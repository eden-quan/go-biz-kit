package servers

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	middlewarepkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/middleware"

	contextkit "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/context"
	middlewareutil "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/middleware"
	"gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/tracing"
)

func AuthorizationMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			// 默认 Token 处理
			header := tr.RequestHeader()
			token := header.Get(contextkit.AuthorizationKey)

			if token != "" {
				ctx = context.WithValue(ctx, contextkit.AuthorizationKey, token)
			}

			// 内部鉴权 Token 处理
			innerToken := header.Get(contextkit.InnerAuthorizationKey)
			if innerToken != "" {
				ctx = context.WithValue(ctx, contextkit.InnerAuthorizationKey, innerToken)
			}

			return handler(ctx, req)
		}
	}
}

// DefaultServerMiddlewares 中间件
func DefaultServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		HttpContextMiddleware(),
		MetricsMiddleware(),
		recovery.Recovery(recovery.WithHandler(middlewareutil.RecoveryHandler())),
		metadata.Server(),
		tracing.Server(),
		middlewarepkg.RequestAndResponseHeader(),
		AuthorizationMiddleware(),
	}
}

func DefaultGrpcServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		HttpContextMiddleware(),
		recovery.Recovery(recovery.WithHandler(middlewareutil.RecoveryHandler())),
		metadata.Server(),
		tracing.Server(),
		middlewarepkg.RequestAndResponseHeader(),
		AuthorizationMiddleware(),
	}
}
