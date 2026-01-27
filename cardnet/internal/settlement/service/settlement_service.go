package service

import (
	"context"

	"cardnet/internal/settlement/domain"
)

type SettlementRepository interface {
	Insert(context.Context, domain.SettlementAuthorization) error
}

type SettlementService struct {
	repo SettlementRepository
}

func NewSettlementService(repo SettlementRepository) *SettlementService {
	return &SettlementService{repo: repo}
}

func (s *SettlementService) ProcessAuthorization(
	ctx context.Context,
	auth domain.SettlementAuthorization,
) error {
	return s.repo.Insert(ctx, auth)
}
