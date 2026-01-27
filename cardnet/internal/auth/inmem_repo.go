package auth

import (
	"context"
	"sync"
)

type inMemoryRepository struct {
	mu    sync.Mutex
	store map[string]struct {
		authID string
		status string
	}
}

func NewInMemoryRepo() Repository {
	return &inMemoryRepository{
		store: make(map[string]struct {
			authID string
			status string
		}),
	}
}

func (r *inMemoryRepository) SaveAuthorization(
	ctx context.Context,
	authID string,
	requestID string,
	merchantID string,
	cardToken string,
	amount int64,
	status string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[requestID]; exists {
		return ErrDuplicateRequest
	}

	r.store[requestID] = struct {
		authID string
		status string
	}{
		authID: authID,
		status: status,
	}
	return nil
}

func (r *inMemoryRepository) GetByRequestID(
	ctx context.Context,
	requestID string,
) (string, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rec, ok := r.store[requestID]
	if !ok {
		return "", "", ErrNotFound
	}
	return rec.authID, rec.status, nil
}

func (r *inMemoryRepository) SaveRiskAudit(
	ctx context.Context,
	authID string,
	riskScore int,
	reasons []string,
) error {
	// best-effort no-op
	return nil
}
