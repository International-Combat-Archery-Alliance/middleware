package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func AccessLogging(baseLogger *slog.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := uuid.New()
			start := time.Now()

			requestLogger := baseLogger.With(slog.String("request-id", requestId.String()))
			ctx := CtxWithRequestId(r.Context(), requestId)
			ctx = CtxWithLogger(ctx, requestLogger)

			loggingRW := newLoggingResponseWriter(w)

			// process the request
			next.ServeHTTP(loggingRW, r.WithContext(ctx))

			requestLogger.InfoContext(r.Context(),
				"Access log",
				slog.String("latency", formatDuration(time.Since(start))),
				slog.Int64("request-content-length", r.ContentLength),
				slog.Int("resp-body-size", loggingRW.responseSize),
				slog.String("host", r.Host),
				slog.String("method", r.Method),
				slog.Int("status-code", loggingRW.statusCode),
				slog.String("path", r.URL.Path),
			)
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
