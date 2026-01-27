package scheduler

import (
	"cardnet/internal/settlement/service"
	"context"
	"time"
)

type T1Scheduler struct {
	batchProcessor *service.BatchProcessor
	reconService   *service.ReconciliationService
}

func NewT1Scheduler(
	p *service.BatchProcessor,
	r *service.ReconciliationService,
) *T1Scheduler {
	return &T1Scheduler{
		batchProcessor: p,
		reconService:   r,
	}
}

func (s *T1Scheduler) Run(ctx context.Context) error {
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Truncate(24 * time.Hour)

	batchID, err := s.batchProcessor.RunForDate(ctx, yesterday)
	if err != nil {
		return err
	}

	return s.reconService.ReconcileBatch(ctx, batchID)
}
