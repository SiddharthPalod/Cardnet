package service

import (
	"cardnet/internal/settlement/domain"
	"cardnet/internal/settlement/reporting"
	"context"
	"fmt"
	"log"
	"os"
)

type ReconciliationRepository interface {
	FindDuplicates(context.Context, string) ([]string, error)
	HasCriticalNulls(context.Context, string) (bool, error)
	CountAuthorizations(context.Context, string) (int64, error)
	SaveIssue(context.Context, string, string) error
	MarkBatchStatus(context.Context, string, string) error
}

type ReconciliationService struct {
	repo      ReconciliationRepository
	generator *reporting.Generator
}

func NewReconciliationService(
	repo ReconciliationRepository,
	gen *reporting.Generator,
) *ReconciliationService {
	return &ReconciliationService{
		repo:      repo,
		generator: gen,
	}
}

func (s *ReconciliationService) ReconcileBatch(ctx context.Context, batchID string) error {

	count, err := s.repo.CountAuthorizations(ctx, batchID)
	if err != nil {
		return err
	}

	if count == 0 {
		s.repo.SaveIssue(ctx, batchID, "empty batch")
		return s.repo.MarkBatchStatus(ctx, batchID, string(domain.Failed))
	}

	dups, err := s.repo.FindDuplicates(ctx, batchID)
	if err != nil {
		return err
	}

	for _, dup := range dups {
		s.repo.SaveIssue(ctx, batchID, fmt.Sprintf("duplicate authorization: %s", dup))
	}

	if len(dups) > 0 {
		return s.repo.MarkBatchStatus(ctx, batchID, string(domain.Failed))
	}

	hasNulls, err := s.repo.HasCriticalNulls(ctx, batchID)
	if err != nil {
		return err
	}

	if hasNulls {
		s.repo.SaveIssue(ctx, batchID, "critical null fields detected")
		return s.repo.MarkBatchStatus(ctx, batchID, string(domain.Failed))
	}

	if err := s.repo.MarkBatchStatus(ctx, batchID, string(domain.Reconciled)); err != nil {
		return err
	}

	// Generate Report
	if s.generator != nil {
		report, err := s.generator.GenerateBatchReport(ctx, batchID)
		if err != nil {
			log.Printf("Failed to generate report for batch %s: %v", batchID, err)
			return nil // Don't fail the reconciliation if reporting fails
		}

		_ = os.MkdirAll("reports", 0755)
		filename := fmt.Sprintf("reports/batch_%s.csv", batchID)
		if err := reporting.ExportCSV(report, filename); err != nil {
			log.Printf("Failed to export report for batch %s: %v", batchID, err)
		} else {
			log.Printf("Report generated: %s", filename)
		}
	}

	return nil
}
