package engine

import (
	"cardnet/internal/risk/models"
	"cardnet/internal/risk/rules"
)

type RiskEngine struct {
	rules []rules.Rule
}

func New() *RiskEngine {
	return &RiskEngine{
		rules: []rules.Rule{
			rules.VelocityRule,
			rules.AmountRule,
			rules.GeoRule,
			rules.MerchantRule,
			rules.BINRule,
		},
	}
}

func (e *RiskEngine) Evaluate(
	tx models.Transaction,
	ctx models.RiskContext,
) models.RiskResult {

	score := 0
	reasons := []string{}

	for _, rule := range e.rules {
		s, r := rule(tx, ctx)
		score += s
		reasons = append(reasons, r...)
	}

	return models.RiskResult{
		RiskScore: score,
		Decision:  mapDecision(score),
		Reasons:   reasons,
	}
}
