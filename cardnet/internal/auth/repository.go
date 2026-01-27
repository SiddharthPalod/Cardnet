package auth

import (
	"context"
	"errors"
)

var (
	ErrDuplicateRequest = errors.New("duplicate request_id")
	ErrNotFound         = errors.New("not found")
)

// 🔑 Interface used by Service
type Repository interface {
	GetByRequestID(ctx context.Context, requestID string) (string, string, error)
	SaveAuthorization(
		ctx context.Context,
		authID string,
		requestID string,
		merchantID string,
		cardToken string,
		amount int64,
		status string,
	) error
	SaveRiskAudit(
		ctx context.Context,
		authID string,
		riskScore int,
		reasons []string,
	) error
}
