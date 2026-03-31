package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CorsConfig holds configuration for CORS middleware
type CorsConfig struct {
	LocalOrigins []string
	ProdOrigins  []string
	IsProduction bool
}

// DefaultCorsConfig returns a default CORS configuration for ICAA services
func DefaultCorsConfig() CorsConfig {
	return CorsConfig{
		LocalOrigins: []string{"http://localhost:4173", "http://localhost:5173"},
		ProdOrigins: []string{
			"https://www.icaa.world",
			"https://icaa.world",
			"https://www.*-icaa-world.curly-sound-f2cd.workers.dev",
			"https://*-icaa-world.curly-sound-f2cd.workers.dev",
		},
		IsProduction: false,
	}
}

// CorsMiddleware creates a CORS middleware with the given configuration
func CorsMiddleware(config CorsConfig) MiddlewareFunc {
	var allowedOrigins []string
	if config.IsProduction {
		allowedOrigins = config.ProdOrigins
	} else {
		allowedOrigins = config.LocalOrigins
	}

	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		MaxAge:           300,
		AllowCredentials: true,
	})

	return c.Handler
}
