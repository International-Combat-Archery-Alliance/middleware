package middleware

import (
	"log/slog"
	"net/http"
	"net/url"
)

func BaseNamePrefix(logger *slog.Logger, baseName string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlWithBasePath, err := url.JoinPath(baseName, r.URL.Path)
			if err != nil {
				logger.Error("url.JoinPath returned an error somehow?", slog.String("error", err.Error()))
			} else {
				r.URL.Path = urlWithBasePath
			}

			next.ServeHTTP(w, r)
		})
	}
}
