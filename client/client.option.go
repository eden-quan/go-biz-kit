package clientutil

import (
	"context"

	authpkg "github.com/eden-quan/go-kratos-pkg/auth"
	middlewarepkg "github.com/eden-quan/go-kratos-pkg/middleware"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/eden-quan/go-biz-kit/tracing"
)

func AuthorizationMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			tr, ok := transport.FromClientContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			token, ok := ctx.Value(authpkg.AuthorizationKey).(string)
			if ok {
				header := tr.RequestHeader()
				header.Set(authpkg.AuthorizationKey, token)
			}

			return handler(ctx, req)
		}
	}
}

func DefaultClientMiddlewares(logger log.Logger) []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(),
		metadata.Client(),
		tracing.Client(),
		middlewarepkg.ClientLogging(logger),
		AuthorizationMiddleware(),
	}
}
