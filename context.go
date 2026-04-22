package middleware

import (
	"context"
	"log/slog"

	"github.com/International-Combat-Archery-Alliance/auth"
)

const (
	ctxTraceIDKey       = "TRACE_ID"
	ctxSpanIDKey        = "SPAN_ID"
	ctxLoggerKey        = "LOGGER"
	ctxJWTKey           = "JWT"
	ctxRefreshTokenIDKey = "REFRESH_TOKEN_ID"
)

func CtxWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxTraceIDKey, traceID)
}

func GetTraceIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxTraceIDKey)
	if v == nil {
		return "", false
	}
	traceID, ok := v.(string)
	return traceID, ok
}

func CtxWithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, ctxSpanIDKey, spanID)
}

func GetSpanIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxSpanIDKey)
	if v == nil {
		return "", false
	}
	spanID, ok := v.(string)
	return spanID, ok
}

func CtxWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

func GetLoggerFromCtx(ctx context.Context) (*slog.Logger, bool) {
	v := ctx.Value(ctxLoggerKey)
	if v == nil {
		return nil, false
	}
	logger, ok := v.(*slog.Logger)
	return logger, ok
}

func CtxWithJWT(ctx context.Context, jwt auth.AuthToken) context.Context {
	return context.WithValue(ctx, ctxJWTKey, jwt)
}

func GetJWTFromCtx(ctx context.Context) (auth.AuthToken, bool) {
	v := ctx.Value(ctxJWTKey)
	if v == nil {
		return nil, false
	}
	token, ok := v.(auth.AuthToken)
	return token, ok
}

// CtxWithRefreshTokenID stores a refresh token ID in the context
func CtxWithRefreshTokenID(ctx context.Context, tokenID string) context.Context {
	return context.WithValue(ctx, ctxRefreshTokenIDKey, tokenID)
}

// GetRefreshTokenIDFromCtx retrieves a refresh token ID from the context
func GetRefreshTokenIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxRefreshTokenIDKey)
	if v == nil {
		return "", false
	}
	tokenID, ok := v.(string)
	return tokenID, ok
}
