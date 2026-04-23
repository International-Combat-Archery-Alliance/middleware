package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// OTELHandler wraps the provided handler with OpenTelemetry HTTP tracing.
// It uses a span name formatter of "METHOD PATH" (e.g. "GET /login/v1/sessions").
func OTELHandler(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "",
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}

// FlushTraces returns a middleware that flushes trace spans after the request
// has been handled. The flush runs with the provided timeout.
// If flush is nil, the middleware is a no-op.
func FlushTraces(flush func(context.Context) error, logger *slog.Logger, timeout time.Duration) MiddlewareFunc {
	if flush == nil {
		return func(next http.Handler) http.Handler {
			return next
		}
	}
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err := flush(ctx); err != nil {
				logger.Error("failed to flush traces", slog.String("error", err.Error()))
			}
		})
	}
}
