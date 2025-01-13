package tracing

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	errorpkg "gitlab.lainuoniao.cn/rhinobird/backend/go-kratos-pkg.git/error"
	"go.opentelemetry.io/otel/trace"

	errorutil "gitlab.lainuoniao.cn/rhinobird/backend/go-biz-kit.git/error"
)

type tracingServerError struct {
	info *errorpkg.ErrorMetaInfo
}

func (e *tracingServerError) Error() string {
	msg := "message stack:\n\t"
	errMsg := strings.Split(e.info.Error(), "\n")
	msg += strings.Join(errMsg, "\n\t")

	leaf := e.info.Leaf()
	if leaf.Stack != "" {
		msg += "\n\ntrace stack:\n" + e.info.Leaf().Stack
	}
	return msg
}

func (e *tracingServerError) BizCode() int32 {
	return e.info.Leaf().BizCode
}

func (e *tracingServerError) Code() int32 {
	return e.info.Leaf().Code
}

// Server returns a new server middleware for OpenTelemetry.
func Server() middleware.Middleware {
	tracer := NewTracer(trace.SpanKindServer)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, tr.Operation(), tr.RequestHeader())
				setServerSpan(ctx, span, req)
				defer func() {
					// TODO: we need move this truncate to another middleware
					// check if we need truncate error to nil
					newErr, ok := errorutil.IsTruncateToEmptyError(err)
					if ok {
						err = newErr.IsTruncateToEmpty()
					}

					isEmpty := errorpkg.IsEmptyError(err)
					if isEmpty {
						err = nil
					}

					// modify error impl for formatting tracing log
					tracer.End(ctx, span, reply, err)

					if ok {
						err = nil
					}

				}()
			}
			return handler(ctx, req)
		}
	}
}
