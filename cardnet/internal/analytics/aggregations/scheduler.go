package aggregations

import (
	"context"
	"log"
	"time"
)

// Scheduler manages batch aggregation jobs
type Scheduler struct {
	batchService *BatchAggregationService
	ctx          context.Context
}

// NewScheduler creates a new batch aggregation scheduler
func NewScheduler(batchService *BatchAggregationService) *Scheduler {
	return &Scheduler{
		batchService: batchService,
		ctx:          context.Background(),
	}
}

// Start begins running hourly and daily batch aggregation jobs
func (s *Scheduler) Start() {
	// Run hourly aggregation every hour
	hourlyTicker := time.NewTicker(time.Hour)
	go func() {
		// Run immediately on start (for testing)
		s.runHourlyJob()
		
		for range hourlyTicker.C {
			s.runHourlyJob()
		}
	}()
	
	// Run daily aggregation every 24 hours at midnight UTC
	dailyTicker := time.NewTicker(24 * time.Hour)
	go func() {
		// Wait until next midnight UTC
		now := time.Now().UTC()
		nextMidnight := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
		time.Sleep(time.Until(nextMidnight))
		
		// Run immediately after waiting
		s.runDailyJob()
		
		for range dailyTicker.C {
			s.runDailyJob()
		}
	}()
	
	log.Println("Batch aggregation scheduler started (hourly + daily)")
}

func (s *Scheduler) runHourlyJob() {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Minute)
	defer cancel()
	
	if err := s.batchService.RunHourlyAggregation(ctx); err != nil {
		log.Printf("Hourly batch aggregation failed: %v", err)
	}
}

func (s *Scheduler) runDailyJob() {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Minute)
	defer cancel()
	
	if err := s.batchService.RunDailyAggregation(ctx); err != nil {
		log.Printf("Daily batch aggregation failed: %v", err)
	}
}
