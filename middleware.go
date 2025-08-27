package middleware

import "net/http"

type MiddlewareFunc func(next http.Handler) http.Handler

func UseMiddlewares(r *http.ServeMux, middlewares ...MiddlewareFunc) http.Handler {
	var s http.Handler
	s = r

	for _, mw := range middlewares {
		s = mw(s)
	}

	return s
}
