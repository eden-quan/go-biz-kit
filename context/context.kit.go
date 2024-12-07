package contextutil

import (
	"context"

	authpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/auth"
	contextpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/context"
	headerpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/header"
)

const AuthorizationKey = "Authorization"
const InnerAuthorizationKey = "InnerToken"
const HttpContextKey = "_HTTP_CONTEXT_KEY_"

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
	return context.WithValue(ctx, AuthorizationKey, token)
}

// WithInnerAuthorization 创建一个带有 InnerToken 信息的上下文
func WithInnerAuthorization(ctx context.Context, innerToken string) context.Context {
	return context.WithValue(ctx, InnerAuthorizationKey, innerToken)
}

// GetAuthorizationToken 尝试从上下文中获取 Token 信息
func GetAuthorizationToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(authpkg.AuthorizationKey).(string)
	if !ok {
		return "", false
	}

	return token, true
}

// GetInnerAuthorizationToken 尝试从上下文中获取内部 Token 信息
func GetInnerAuthorizationToken(ctx context.Context) (string, bool) {
	innerToken, ok := ctx.Value(InnerAuthorizationKey).(string)
	if !ok {
		return "", false
	}

	return innerToken, true
}
