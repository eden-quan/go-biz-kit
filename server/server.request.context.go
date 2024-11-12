package servers

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	contextutil "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/context"
)

func HttpContextMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			ctx = context.WithValue(ctx, contextutil.HttpContextKey, ctx)
			return handler(ctx, req)
		}
	}
}
