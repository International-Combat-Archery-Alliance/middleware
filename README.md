# ICAA Middleware

A Go HTTP middleware library for the International Combat Archery Alliance, providing reusable components for web applications including request logging, authentication context, Swagger UI hosting, and path prefix management.

## Features

### Core Middleware
- **Request Logging**: Structured access logging with request IDs, latency tracking, and response metrics
- **Context Management**: Request-scoped context for request IDs, loggers, and JWT tokens
- **Base Path Prefixes**: URL path prefix modification for service routing
- **Middleware Composition**: Utility for composing multiple middlewares

### Swagger UI Integration
- **Embedded Swagger UI**: Self-contained Swagger UI hosting with embedded assets
- **OpenAPI Specification Serving**: Automatic serving of OpenAPI JSON specifications
- **Configurable Base Paths**: Flexible path configuration for API documentation

## Installation

```bash
go get github.com/International-Combat-Archery-Alliance/middleware
```

## Usage

### Basic Middleware Setup

```go
package main

import (
    "log/slog"
    "net/http"
    
    "github.com/International-Combat-Archery-Alliance/middleware"
)

func main() {
    mux := http.NewServeMux()
    logger := slog.Default()
    
    // Add your routes
    mux.HandleFunc("/api/health", healthHandler)
    
    // Apply middlewares
    handler := middleware.UseMiddlewares(mux,
        middleware.AccessLogging(logger),
        middleware.BaseNamePrefix(logger, "/api/v1"),
    )
    
    http.ListenAndServe(":8080", handler)
}
```

### Swagger UI Hosting

```go
import (
    "github.com/getkin/kin-openapi/openapi3"
    "github.com/International-Combat-Archery-Alliance/middleware"
)

func setupSwagger(spec *openapi3.T) http.Handler {
    mux := http.NewServeMux()
    
    // Host Swagger UI at /docs/swagger-ui with spec at /docs/openapi.json
    swaggerMiddleware, err := middleware.HostSwaggerUI("/docs", spec)
    if err != nil {
        log.Fatal(err)
    }
    
    handler := middleware.UseMiddlewares(mux, swaggerMiddleware)
    return handler
}
```

### Context Usage

```go
func apiHandler(w http.ResponseWriter, r *http.Request) {
    // Get request-scoped logger
    if logger, ok := middleware.GetLoggerFromCtx(r.Context()); ok {
        logger.Info("Processing API request")
    }
    
    // Get request ID
    if reqID, ok := middleware.GetRequestIdFromCtx(r.Context()); ok {
        w.Header().Set("X-Request-ID", reqID.String())
    }
    
    // Get JWT token (if authentication middleware is used)
    if token, ok := middleware.GetJWTFromCtx(r.Context()); ok {
        // Use token for authorization
        _ = token
    }
}
```

## Middleware Components

### AccessLogging
Provides structured HTTP access logging with:
- Request ID generation and tracking
- Request/response size logging
- Latency measurement
- HTTP method, path, and status code logging
- Integration with Go's `slog` structured logging

### BaseNamePrefix
Modifies request URLs by prepending a base path, useful for:
- API versioning (`/api/v1`)
- Service routing in reverse proxy setups
- Path-based service isolation

### HostSwaggerUI
Serves Swagger UI documentation with:
- Embedded Swagger UI assets (no external dependencies)
- Automatic OpenAPI specification serving
- Configurable base paths for documentation endpoints

## Context Utilities

The library provides context management utilities for:
- **Request IDs**: UUID-based request tracking across middleware chain
- **Structured Logging**: Request-scoped loggers with automatic request ID correlation
- **JWT Tokens**: Authentication token storage and retrieval (integrates with ICAA auth library)

## Dependencies

- `github.com/International-Combat-Archery-Alliance/auth`: Authentication and JWT handling
- `github.com/getkin/kin-openapi`: OpenAPI specification parsing and validation
- `github.com/google/uuid`: UUID generation for request tracking
- Standard library packages for HTTP handling and logging

## License

Licensed under the GNU Affero General Public License v3.0. See [LICENSE](LICENSE) for details.

## Contributing

This library is part of the International Combat Archery Alliance's software infrastructure. Contributions should follow the project's coding standards and include appropriate tests.
