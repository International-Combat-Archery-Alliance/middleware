package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func AccessLogging(baseLogger *slog.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger := baseLogger
			ctx := r.Context()

			if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
				spanCtx := span.SpanContext()
				logger = baseLogger.With(
					slog.String("trace_id", spanCtx.TraceID().String()),
					slog.String("span_id", spanCtx.SpanID().String()),
				)
				ctx = CtxWithTraceID(ctx, spanCtx.TraceID().String())
				ctx = CtxWithSpanID(ctx, spanCtx.SpanID().String())
			}

			ctx = CtxWithLogger(ctx, logger)

			loggingRW := newLoggingResponseWriter(w)

			// process the request
			next.ServeHTTP(loggingRW, r.WithContext(ctx))

			attrs := []slog.Attr{
				slog.String("latency", formatDuration(time.Since(start))),
				slog.Int64("request-content-length", r.ContentLength),
				slog.Int("resp-body-size", loggingRW.responseSize),
				slog.String("host", r.Host),
				slog.String("method", r.Method),
				slog.Int("status-code", loggingRW.statusCode),
				slog.String("path", r.URL.Path),
			}

			if traceID, ok := GetTraceIDFromCtx(ctx); ok {
				attrs = append(attrs, slog.String("trace_id", traceID))
			}
			if spanID, ok := GetSpanIDFromCtx(ctx); ok {
				attrs = append(attrs, slog.String("span_id", spanID))
			}

			logger.LogAttrs(r.Context(), slog.LevelInfo, "Access log", attrs...)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		responseSize:   0,
	}
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(data)
	lrw.responseSize += size
	return size, err
}

// formatDuration formats a duration to one decimal point.
func formatDuration(d time.Duration) string {
	div := time.Duration(10)
	switch {
	case d > time.Second:
		d = d.Round(time.Second / div)
	case d > time.Millisecond:
		d = d.Round(time.Millisecond / div)
	case d > time.Microsecond:
		d = d.Round(time.Microsecond / div)
	case d > time.Nanosecond:
		d = d.Round(time.Nanosecond / div)
	}
	return d.String()
}
