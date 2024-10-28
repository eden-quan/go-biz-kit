package servers

import (
	"context"

	authpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/auth"
	middlewarepkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/middleware"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"

	errorutil "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/error"
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

			header := tr.RequestHeader()
			token := header.Get(authpkg.AuthorizationKey)

			if token != "" {
				ctx = context.WithValue(ctx, authpkg.AuthorizationKey, token)
			}

			return handler(ctx, req)
		}
	}
}

// DefaultServerMiddlewares 中间件
func DefaultServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(recovery.WithHandler(middlewareutil.RecoveryHandler())),
		metadata.Server(),
		tracing.Server(),
		errorutil.ErrorResultMiddleware(),
		middlewarepkg.RequestAndResponseHeader(),
		AuthorizationMiddleware(),
		middlewareutil.Validator(),
	}
}

func DefaultGrpcServerMiddlewares() []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(recovery.WithHandler(middlewareutil.RecoveryHandler())),
		metadata.Server(),
		tracing.Server(),
		errorutil.ErrorResultMiddleware(),
		middlewareutil.Validator(),
		middlewarepkg.RequestAndResponseHeader(),
		AuthorizationMiddleware(),
	}
}
