package rules

import (
	"cardnet/internal/risk/config"
	"cardnet/internal/risk/models"
)

func BINRule(tx models.Transaction, ctx models.RiskContext) (int, []string) {
	if ctx.BINRiskScore >= config.BINHighRiskScore {
		return config.BINHighScore, []string{"HIGH_RISK_BIN"}
	}

	if ctx.BINRiskScore >= config.BINMediumRiskScore {
		return config.BINMediumScore, []string{"MEDIUM_RISK_BIN"}
	}

	return 0, nil
}
