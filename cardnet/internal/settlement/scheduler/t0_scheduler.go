package scheduler

import (
	"context"

	"cardnet/internal/settlement/service"
)

type T0Scheduler struct {
	processor *service.BatchProcessor
}

func NewT0Scheduler(p *service.BatchProcessor) *T0Scheduler {
	return &T0Scheduler{processor: p}
}

func (s *T0Scheduler) Run(ctx context.Context) error {
	return s.processor.RunT0Batch(ctx)
}
