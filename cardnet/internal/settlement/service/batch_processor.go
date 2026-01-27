package service

import (
	"context"
	"time"
)

type BatchRepository interface {
	GetOrCreateBatch(context.Context, time.Time) (string, error)
	AttachAuthorizations(context.Context, string, time.Time) error
	MarkReady(context.Context, string) error
}

type BatchProcessor struct {
	repo BatchRepository
}

func NewBatchProcessor(repo BatchRepository) *BatchProcessor {
	return &BatchProcessor{repo: repo}
}

func (p *BatchProcessor) RunT0Batch(ctx context.Context) error {

	today := time.Now().UTC().Truncate(24 * time.Hour)

	batchID, err := p.repo.GetOrCreateBatch(ctx, today)
	if err != nil {
		return err
	}

	if err := p.repo.AttachAuthorizations(ctx, batchID, today); err != nil {
		return err
	}

	return p.repo.MarkReady(ctx, batchID)
}

func (p *BatchProcessor) RunForDate(
	ctx context.Context,
	date time.Time,
) (string, error) {
	batchID, err := p.repo.GetOrCreateBatch(ctx, date)
	if err != nil {
		return "", err
	}

	if err := p.repo.AttachAuthorizations(ctx, batchID, date); err != nil {
		return "", err
	}

	return batchID, p.repo.MarkReady(ctx, batchID)
}
