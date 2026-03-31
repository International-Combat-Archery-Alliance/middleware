package middleware

import (
	"context"
	"log/slog"

	"github.com/International-Combat-Archery-Alliance/auth"
	"github.com/google/uuid"
)

const (
	ctxRequestIdKey      = "REQUEST_ID"
	ctxLoggerKey         = "LOGGER"
	ctxJWTKey            = "JWT"
	ctxRefreshTokenIDKey = "REFRESH_TOKEN_ID"
)

func CtxWithRequestId(ctx context.Context, requestId uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxRequestIdKey, requestId)
}

func GetRequestIdFromCtx(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(ctxRequestIdKey)
	if v == nil {
		return uuid.UUID{}, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
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
