package inference

import (
	"cardnet/internal/analytics/feedback"
	"context"
	"time"
)

func SuggestRuleThreshold(
	ctx context.Context,
	ruleID string,
	current int,
	signal float64,
	publisher *feedback.Publisher,
) int {
	suggested := current
	if signal > 0.8 {
		suggested = current - 2
	} else if signal < 0.3 {
		suggested = current + 1
	}

	if suggested != current && publisher != nil {
		_ = publisher.PublishThresholdUpdate(ctx, feedback.ThresholdUpdate{
			RuleName:       ruleID,
			EntityType:     "global",
			SuggestedValue: suggested,
			Confidence:     0.75,
			ModelVersion:   "rule-optimizer-v1",
			CreatedAt:      time.Now(),
		})
	}

	return suggested
}
