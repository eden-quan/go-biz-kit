package tracing

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type tracingError interface {
	Error() string
	Code() int32
	BizCode() int32
}

type options struct {
	tracerName     string
	tracerProvider trace.TracerProvider
	propagator     propagation.TextMapPropagator
}

// Tracer is otel span tracer
type Tracer struct {
	tracer trace.Tracer
	kind   trace.SpanKind
	opt    *options
}

// NewTracer create tracer instance
func NewTracer(kind trace.SpanKind) *Tracer {
	op := options{
		propagator: propagation.NewCompositeTextMapPropagator(tracing.Metadata{}, propagation.Baggage{}, propagation.TraceContext{}),
		tracerName: "eden",
	}

	if op.tracerProvider == nil {
		op.tracerProvider = otel.GetTracerProvider()
	}

	switch kind {
	case trace.SpanKindClient:
		return &Tracer{tracer: op.tracerProvider.Tracer(op.tracerName), kind: kind, opt: &op}
	case trace.SpanKindServer:
		return &Tracer{tracer: op.tracerProvider.Tracer(op.tracerName), kind: kind, opt: &op}
	default:
		panic(fmt.Sprintf("unsupported span kind: %v", kind))
	}
}

// Start is start tracing span
func (t *Tracer) Start(ctx context.Context, operation string, carrier propagation.TextMapCarrier) (context.Context, trace.Span) {
	if t.kind == trace.SpanKindServer {
		ctx = t.opt.propagator.Extract(ctx, carrier)
	}
	ctx, span := t.tracer.Start(ctx,
		operation,
		trace.WithSpanKind(t.kind),
	)
	if t.kind == trace.SpanKindClient {
		t.opt.propagator.Inject(ctx, carrier)
	}
	return ctx, span
}

// End finish tracing span
func (t *Tracer) End(_ context.Context, span trace.Span, m interface{}, err error) {
	if err != nil {
		span.RecordError(err)
		var t tracingError
		if errors.As(err, &t) {
			span.SetAttributes(attribute.Key("rpc.status_code").Int64(int64(t.Code())))
			span.SetAttributes(attribute.Key("rpc.biz_code").Int64(int64(t.BizCode())))
		} else if e := errors.FromError(err); e != nil {
			span.SetAttributes(attribute.Key("rpc.biz_code").Int64(int64(0)))
			span.SetAttributes(attribute.Key("rpc.status_code").Int64(int64(e.Code)))
		}
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "OK")
	}

	if p, ok := m.(proto.Message); ok && !reflect.ValueOf(m).IsNil() {
		if t.kind == trace.SpanKindServer {
			span.SetAttributes(attribute.Key("send_msg.size").Int(proto.Size(p)))
		} else {
			span.SetAttributes(attribute.Key("recv_msg.size").Int(proto.Size(p)))
		}
	}
	span.End()
}
