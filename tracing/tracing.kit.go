package tracing

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/eden-quan/go-biz-kit/config"
	"github.com/eden-quan/go-biz-kit/config/def"
)

var lock sync.Once

type TracerInitializer struct{}

// InitTracing 设置全局链路跟踪器
func InitTracing(conf *def.Configuration, logger log.Logger, local *config.LocalConfigure) (*TracerInitializer, error) {
	var err error = nil

	lock.Do(func() {

		tracingConf := &conf.Tracing

		if !tracingConf.GetEnable() {
			exp := tracetest.NewNoopExporter()
			tp := trace.NewTracerProvider(trace.WithBatcher(exp))
			otel.SetTracerProvider(tp)
			return
		}

		exp, err := otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithURLPath(tracingConf.GetUrlPath()),
			otlptracehttp.WithEndpoint(tracingConf.GetEndpoint()),
		)

		if err != nil {
			log.NewHelper(logger).Errorf("creating tracing client with error %s", err)
			return
		}

		tp := trace.NewTracerProvider(
			// Configure the simple rate
			trace.WithSampler(
				trace.ParentBased(
					trace.TraceIDRatioBased(tracingConf.GetSimpleRate()),
				)),
			// keep batch process for performance
			trace.WithBatcher(
				exp,
				trace.WithMaxExportBatchSize(int(tracingConf.GetMaxBatchSize())),
				trace.WithMaxQueueSize(int(tracingConf.GetMaxQueueSize())),
			),
			// record service basic info
			trace.WithResource(resource.NewSchemaless(
				semconv.ServiceNameKey.String(local.APP.Name),
				semconv.ServiceVersionKey.String(local.APP.Version),
				attribute.String("exporter", tracingConf.GetType()),
			)),
		)

		otel.SetTracerProvider(tp)

	})

	return &TracerInitializer{}, err
}
