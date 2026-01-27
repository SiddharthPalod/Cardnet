package scheduler

import (
	"context"
	"log"
	"time"
)

type Runner struct {
	t0 *T0Scheduler
	t1 *T1Scheduler
}

func NewRunner(
	t0 *T0Scheduler,
	t1 *T1Scheduler,
) *Runner {
	return &Runner{
		t0: t0,
		t1: t1,
	}
}

func (r *Runner) Start(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Run immediately on start for testing
	log.Println("Running startup settlement jobs...")
	_ = r.t0.Run(ctx)
	// _ = r.t1.Run(ctx) // Optional on startup

	for {
		select {
		case <-ticker.C:
			log.Println("Running daily settlement jobs")
			if err := r.t0.Run(ctx); err != nil {
				log.Printf("T0 job failed: %v", err)
			}
			if err := r.t1.Run(ctx); err != nil {
				log.Printf("T1 job failed: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
