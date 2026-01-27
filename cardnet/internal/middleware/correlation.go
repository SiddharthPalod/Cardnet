package middleware

import (
	"context"

	"github.com/google/uuid"
)

type correlationKeyType string

const correlationKey correlationKeyType = "correlation_id"

func WithCorrelationID(ctx context.Context) context.Context {
	return context.WithValue(ctx, correlationKey, uuid.New().String())
}

func GetCorrelationID(ctx context.Context) string {
	if v := ctx.Value(correlationKey); v != nil {
		return v.(string)
	}
	return ""
}
