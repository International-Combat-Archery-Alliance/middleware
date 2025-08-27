package middleware

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed swagger-ui/dist/*
var swaggerUI embed.FS

func HostSwaggerUI(basePath string, spec *openapi3.T) (MiddlewareFunc, error) {
	openapiServer := http.NewServeMux()

	swaggerUIPath, err := url.JoinPath(basePath, "/swagger-ui")
	if err != nil {
		return nil, fmt.Errorf("error joining basepath with swagger-ui: %w", err)
	}

	openApiJsonPath, err := url.JoinPath(basePath, "/openapi.json")
	if err != nil {
		return nil, fmt.Errorf("error joining basepath with openapi.json: %w", err)
	}

	swaggerUiSubFS, err := fs.Sub(swaggerUI, "swagger-ui/dist")
	if err != nil {
		return nil, fmt.Errorf("error getting swagger ui sub fs: %w", err)
	}

	openapiServer.Handle(swaggerUIPath, http.StripPrefix(swaggerUIPath, http.FileServer(http.FS(swaggerUiSubFS))))
	openapiServer.HandleFunc(openApiJsonPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(spec)
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler, matchedPath := openapiServer.Handler(r)

			if matchedPath == "" {
				next.ServeHTTP(w, r)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}, nil
}
