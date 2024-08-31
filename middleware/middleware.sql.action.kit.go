package middlewareutil

import (
	"context"
	"reflect"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/eden/go-biz-kit/setup"
)

// SQLActionMiddleware 在执行具体的业务逻辑前触发，他根据上下文信息获取当前请求中是否配置了 SQLAction, 如果能够找到对应的 SQLAction
// 则在执行后 SQLAction 后取消业务逻辑的执行，直接使用 SQLAction 作为执行的返回结果
// 如需执行业务逻辑，则需要通过 SQLAction 中的 Inject 类型来注入需要执行的函数名
func SQLActionMiddleware(manager *setup.ActionManager) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				op := tr.Operation()
				action := manager.Find(op)
				if action != nil {
					reply, err := action.ExecuteQuery(ctx, req)
					if reflect.ValueOf(reply).IsNil() {
						reply = nil
					}
					return reply, err
				}
			}

			return handler(ctx, req)
		}
	}
}
