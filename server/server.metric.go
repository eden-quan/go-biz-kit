package servers

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration",
		Help:      "server requests duration(milliseconds).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"component", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"component", "operation", "code", "bizcode"})
)

func init() {
	prometheus.MustRegister(_metricSeconds, _metricRequests)
}

func MetricsMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				component string
				operation string
			)
			startTime := time.Now()
			if trans, ok := transport.FromServerContext(ctx); ok {
				component = trans.Kind().String()
				operation = trans.Operation()
				if httpTr, ok := trans.(http.Transporter); ok {
					operation = httpTr.Request().URL.Path
				}
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
			}
			_metricRequests.WithLabelValues(component, operation, strconv.Itoa(code), "").Inc()
			_metricSeconds.WithLabelValues(component, operation).Observe(float64(time.Since(startTime).Milliseconds()))
			return reply, err
		}
	}
}
