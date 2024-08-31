package contextutil

import (
	"context"

	authpkg "github.com/eden/go-kratos-pkg/auth"
	contextpkg "github.com/eden/go-kratos-pkg/context"
	headerpkg "github.com/eden/go-kratos-pkg/header"
)

// GetTraceID ...
func GetTraceID(ctx context.Context) (string, bool) {
	tr, ok := contextpkg.FromServerContext(ctx)
	if !ok {
		return "", false
	}
	traceID := tr.RequestHeader().Get(headerpkg.RequestID)

	return traceID, traceID != ""
}

// WithAuthorizationToken 创建一个带有 Token 信息的上下文
func WithAuthorizationToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authpkg.AuthorizationKey, token)
}

// GetAuthorizationToken 尝试从上下文中获取 token 信息
func GetAuthorizationToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(authpkg.AuthorizationKey).(string)
	if !ok {
		return "", false
	}

	return token, true
}
