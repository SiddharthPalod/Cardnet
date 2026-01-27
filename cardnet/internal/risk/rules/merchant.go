package rules

import (
	"cardnet/internal/risk/config"
	"cardnet/internal/risk/models"
)

func MerchantRule(tx models.Transaction, ctx models.RiskContext) (int, []string) {
	if ctx.MerchantRiskScore >= config.MerchantHighRiskScore {
		return config.MerchantHighScore, []string{"HIGH_RISK_MERCHANT"}
	}

	if ctx.MerchantRiskScore >= config.MerchantMediumRiskScore {
		return config.MerchantMediumScore, []string{"MEDIUM_RISK_MERCHANT"}
	}

	return 0, nil
}
