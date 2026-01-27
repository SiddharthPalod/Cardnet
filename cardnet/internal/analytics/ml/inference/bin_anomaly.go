package inference

import (
	"cardnet/internal/analytics/feedback"
	"cardnet/internal/analytics/ml/models"
	"context"
	"time"
)

func IsBINAnomalous(
	ctx context.Context,
	model *models.BINAnomalyModel,
	bin string,
	value float64,
	publisher *feedback.Publisher,
) bool {
	isAnomalous := value > model.Threshold

	if isAnomalous && publisher != nil {
		_ = publisher.PublishRiskSuggestion(ctx, feedback.RiskSuggestion{
			EntityType:     "bin",
			EntityID:       bin,
			SuggestedScore: 100, // High risk
			Confidence:     0.9,
			ModelVersion:   model.Version,
			Reason:         "BIN anomaly detected (high decline/volume)",
			CreatedAt:      time.Now(),
		})
	}

	return isAnomalous
}
