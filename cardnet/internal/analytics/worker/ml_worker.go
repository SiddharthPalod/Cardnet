package worker

import (
	"cardnet/internal/analytics/ml/inference"
	"cardnet/internal/analytics/ml/models"
	"cardnet/internal/analytics/repository/postgres"
	"cardnet/internal/gates"
	"context"
	"log"
	"time"
)

type MLWorker struct {
	repo  *postgres.Repository
	model *models.MerchantRiskModel
	gates *gates.Gates
}

func NewMLWorker(repo *postgres.Repository, gatesInstance *gates.Gates) *MLWorker {
	// Initialize a dummy model
	// In a real system, this would load from a file or registry
	model := &models.MerchantRiskModel{
		Version: "v1.0.0",
		Bias:    -2.0, // Base negative bias (assume good)
		Weights: map[string]float64{
			"tx_count":      0.01,
			"approval_rate": -5.0, // Low approval rate increases risk score (negative * negative vs positive)
			// Wait, approval_rate is 0.0-1.0.
			// If approval_rate is 0.1 (low), we want HIGH score.
			// If coefficient is negative: -5 * 0.1 = -0.5.
			// If approval_rate is 0.9 (high), -5 * 0.9 = -4.5.
			// So lower approval rate = higher score (less negative). Correct.
			"avg_amount": 0.001,
		},
	}
	return &MLWorker{
		repo:  repo,
		model: model,
		gates: gatesInstance,
	}
}

func (w *MLWorker) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			w.runScoringJob()
		}
	}()
	log.Println("ML Worker started (1m interval)")
}

func (w *MLWorker) runScoringJob() {
	ctx := context.Background()
	
	// Check ML safety gate - skip if ML is disabled or data is stale
	if w.gates != nil && w.gates.ShouldIgnoreMLSuggestions() {
		log.Println("ML Scoring Job skipped - ML safety gate closed (data stale or ML disabled)")
		return
	}
	
	log.Println("Running ML Scoring Job...")

	stats, err := w.repo.GetAllMerchantStats(ctx)
	if err != nil {
		log.Printf("Failed to fetch merchant stats: %v", err)
		return
	}

	for _, s := range stats {
		features := map[string]float64{
			"tx_count":      float64(s.TxCount),
			"approval_rate": s.ApprovalRate(),
			"avg_amount":    s.TotalAmount / float64(s.TxCount), // approximate
		}
		if s.TxCount == 0 {
			features["avg_amount"] = 0
		}

		score := inference.ScoreMerchant(ctx, w.model, s.MerchantID, features, nil)
		if score > 0.7 {
			log.Printf("[RISK ALERT] Merchant %s has high risk score: %.4f", s.MerchantID, score)
		}
	}
	log.Printf("Scored %d merchants", len(stats))
}
