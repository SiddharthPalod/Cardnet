package publisher

import (
	"cardnet/internal/ledger/domain"
	"context"
)

type EventPublisher interface {
	Publish(ctx context.Context, event *domain.AuthEvent) error
}
