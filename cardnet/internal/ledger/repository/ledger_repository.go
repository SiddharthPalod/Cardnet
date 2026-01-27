package repository

import (
	"cardnet/internal/ledger/domain"
	"context"
)

type LedgerRepository interface {
	Store(ctx context.Context, event *domain.AuthEvent) error
}
