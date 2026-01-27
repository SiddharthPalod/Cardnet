package inference

import (
	"cardnet/internal/analytics/feedback"
	"cardnet/internal/analytics/ml/models"
	"context"
	"math"
	"time"
)

func ScoreMerchant(
	ctx context.Context,
	model *models.MerchantRiskModel,
	merchantID string,
	features map[string]float64,
	publisher *feedback.Publisher,
) float64 {

	score := model.Bias
	for k, v := range features {
		if w, ok := model.Weights[k]; ok {
			score += w * v
		}
	}

	probability := 1 / (1 + math.Exp(-score))

	if probability > 0.8 && publisher != nil {
		_ = publisher.PublishRiskSuggestion(ctx, feedback.RiskSuggestion{
			EntityType:     "merchant",
			EntityID:       merchantID,
			SuggestedScore: float64(int(probability * 100)),
			Confidence:     0.85,
			ModelVersion:   model.Version,
			Reason:         "High merchant risk score calculated",
			CreatedAt:      time.Now(),
		})
	}

	return probability
}
