package aggregations

import (
	"context"
	"log"
	"time"
)

// BatchAggregationService handles hourly and daily batch aggregations
type BatchAggregationService struct {
	repo BatchRepository
}

// NewBatchAggregationService creates a new batch aggregation service
func NewBatchAggregationService(repo BatchRepository) *BatchAggregationService {
	return &BatchAggregationService{repo: repo}
}

// RunHourlyAggregation aggregates current stats into hourly buckets
func (s *BatchAggregationService) RunHourlyAggregation(ctx context.Context) error {
	log.Println("Starting hourly batch aggregation...")
	
	// Get the previous hour bucket (aggregate data from last hour)
	now := time.Now().UTC()
	hourBucket := now.Truncate(time.Hour).Add(-time.Hour)
	
	// Aggregate merchant stats
	merchants, err := s.repo.GetAllMerchantStats(ctx)
	if err != nil {
		log.Printf("Failed to get merchant stats: %v", err)
		return err
	}
	
	for _, m := range merchants {
		// For hourly aggregation, we snapshot current state
		// In a real system, you'd calculate delta from previous hour
		if err := s.repo.RecordHourlyMerchantAggregation(
			ctx,
			m.MerchantID,
			hourBucket,
			m.TxCount,
			m.ApprovedTx,
			m.TotalAmount,
		); err != nil {
			log.Printf("Failed to record hourly merchant aggregation for %s: %v", m.MerchantID, err)
			// Continue with other merchants
		}
	}
	
	// Aggregate network metrics
	networkMetrics, err := s.repo.GetNetworkMetrics(ctx)
	if err != nil {
		log.Printf("Failed to get network metrics: %v", err)
		return err
	}
	
	if err := s.repo.RecordHourlyNetworkAggregation(
		ctx,
		hourBucket,
		networkMetrics.TotalTx,
		networkMetrics.ApprovedTx,
	); err != nil {
		log.Printf("Failed to record hourly network aggregation: %v", err)
		return err
	}
	
	// Aggregate BIN stats
	bins, err := s.repo.GetAllBinStats(ctx)
	if err != nil {
		log.Printf("Failed to get BIN stats: %v", err)
		// Continue without BIN aggregation
	} else {
		for _, b := range bins {
			if err := s.repo.RecordHourlyBINAggregation(
				ctx,
				b.BIN,
				hourBucket,
				b.TxCount,
				b.ApprovedTx,
				b.TotalAmount,
			); err != nil {
				log.Printf("Failed to record hourly BIN aggregation for %s: %v", b.BIN, err)
				// Continue with other BINs
			}
		}
	}
	
	log.Printf("Hourly batch aggregation completed for hour %s", hourBucket.Format(time.RFC3339))
	return nil
}

// RunDailyAggregation aggregates hourly stats into daily buckets
func (s *BatchAggregationService) RunDailyAggregation(ctx context.Context) error {
	log.Println("Starting daily batch aggregation...")
	
	// Get yesterday's date bucket
	now := time.Now().UTC()
	yesterday := now.Truncate(24 * time.Hour).AddDate(0, 0, -1)
	
	// Aggregate merchant stats from hourly aggregations
	// In a real system, you'd query hourly_merchant_aggregations and sum them
	// For simplicity, we'll snapshot current state similar to hourly
	merchants, err := s.repo.GetAllMerchantStats(ctx)
	if err != nil {
		log.Printf("Failed to get merchant stats: %v", err)
		return err
	}
	
	for _, m := range merchants {
		if err := s.repo.RecordDailyMerchantAggregation(
			ctx,
			m.MerchantID,
			yesterday,
			m.TxCount,
			m.ApprovedTx,
			m.TotalAmount,
		); err != nil {
			log.Printf("Failed to record daily merchant aggregation for %s: %v", m.MerchantID, err)
			// Continue with other merchants
		}
	}
	
	// Aggregate BIN stats
	bins, err := s.repo.GetAllBinStats(ctx)
	if err != nil {
		log.Printf("Failed to get BIN stats: %v", err)
		// Continue without BIN aggregation
	} else {
		for _, b := range bins {
			if err := s.repo.RecordDailyBINAggregation(
				ctx,
				b.BIN,
				yesterday,
				b.TxCount,
				b.ApprovedTx,
				b.TotalAmount,
			); err != nil {
				log.Printf("Failed to record daily BIN aggregation for %s: %v", b.BIN, err)
				// Continue with other BINs
			}
		}
	}
	
	// Aggregate network metrics
	networkMetrics, err := s.repo.GetNetworkMetrics(ctx)
	if err != nil {
		log.Printf("Failed to get network metrics: %v", err)
		return err
	}
	
	if err := s.repo.RecordDailyNetworkAggregation(
		ctx,
		yesterday,
		networkMetrics.TotalTx,
		networkMetrics.ApprovedTx,
	); err != nil {
		log.Printf("Failed to record daily network aggregation: %v", err)
		return err
	}
	
	log.Printf("Daily batch aggregation completed for day %s", yesterday.Format("2006-01-02"))
	return nil
}
