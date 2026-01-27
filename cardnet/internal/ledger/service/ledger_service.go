package service

import (
	"cardnet/internal/ledger/publisher"
	"context"
	"log"

	"cardnet/internal/ledger/domain"
	"cardnet/internal/ledger/repository"
)

type LedgerService struct {
	repo      repository.LedgerRepository
	publisher publisher.EventPublisher
}

func NewLedgerService(repo repository.LedgerRepository,
	publisher publisher.EventPublisher) *LedgerService {
	return &LedgerService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *LedgerService) Record(_ context.Context, e *domain.AuthEvent) {
	// Detach from request lifecycle
	go func() {
		ctx := context.Background()

		if err := s.repo.Store(ctx, e); err != nil {
			log.Printf("ledger write failed: %v", err)
		}

		if err := s.publisher.Publish(ctx, e); err != nil {
			log.Printf("kafka publish failed: %v", err)
		}
	}()
}
