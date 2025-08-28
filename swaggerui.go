package middleware

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed swagger-ui/dist
var swaggerUI embed.FS

func HostSwaggerUI(basePath string, spec *openapi3.T) (MiddlewareFunc, error) {
	openapiServer := http.NewServeMux()

	swaggerUIPath, err := url.JoinPath(basePath, "/swagger-ui/")
	if err != nil {
		return nil, fmt.Errorf("error joining basepath with swagger-ui: %w", err)
	}

	swaggerInitPath, err := url.JoinPath(swaggerUIPath, "/swagger-initializer.js")
	if err != nil {
		return nil, fmt.Errorf("error joining path for swagger init: %w", err)
	}

	openApiJsonPath, err := url.JoinPath(basePath, "/openapi.json")
	if err != nil {
		return nil, fmt.Errorf("error joining basepath with openapi.json: %w", err)
	}

	swaggerUISubFS, err := fs.Sub(swaggerUI, "swagger-ui/dist")
	if err != nil {
		return nil, fmt.Errorf("error getting swagger ui sub fs: %w", err)
	}

	// template the js initialize with the actual openapi spec path
	openapiServer.HandleFunc(swaggerInitPath, func(w http.ResponseWriter, r *http.Request) {
		b, err := swaggerUI.ReadFile("swagger-ui/dist/swagger-initializer.js")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}
		t, err := template.New("init").Parse(string(b))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		t.Execute(w, map[string]any{
			"PathToSpec": openApiJsonPath,
		})
	})
	openapiServer.Handle(swaggerUIPath, http.StripPrefix(swaggerUIPath, http.FileServer(http.FS(swaggerUISubFS))))
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
