package contextutil

import (
	"context"

	authpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/auth"
	contextpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/context"
	headerpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/header"
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
