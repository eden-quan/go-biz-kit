package errorutil

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

func TruncateErrorMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {

			defer func() {
				// TODO: we need move this truncate to another middleware
				// check if we need truncate error to nil
				newErr, ok := IsTruncateToEmptyError(err)
				if ok {
					err = newErr.IsTruncateToEmpty()
				}

			}()

			return handler(ctx, req)
		}
	}
}
