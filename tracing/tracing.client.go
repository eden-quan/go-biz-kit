package tracing

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	errorpkg "github.com/eden/go-kratos-pkg/error"
	"go.opentelemetry.io/otel/trace"
)

type tracingClientError struct {
	info *errorpkg.MetaInfo
}

func (e *tracingClientError) Error() string {
	msg := "message stack:\n\t"
	errMsg := strings.Split(e.info.Message, "\n")
	msg += strings.Join(errMsg, "\n\t")
	return msg
}

func (e *tracingClientError) BizCode() int32 {
	return e.info.BizCode
}

func (e *tracingClientError) Code() int32 {
	return e.info.Code
}

// Client returns a new client middleware for OpenTelemetry.
func Client() middleware.Middleware {
	tracer := NewTracer(trace.SpanKindClient)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, tr.Operation(), tr.RequestHeader())
				setClientSpan(ctx, span, req)
				defer func() {
					// modify error impl for formatting tracing log
					if meta := errorpkg.MetaFromError(err); meta != nil {
						tracer.End(ctx, span, reply, &tracingClientError{info: meta})
					} else {
						tracer.End(ctx, span, reply, err)
					}
				}()
			}
			return handler(ctx, req)
		}
	}
}
