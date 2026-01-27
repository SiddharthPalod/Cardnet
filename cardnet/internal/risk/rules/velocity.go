package rules

import (
	"cardnet/internal/risk/config"
	"cardnet/internal/risk/models"
)

func VelocityRule(tx models.Transaction, ctx models.RiskContext) (int, []string) {
	if ctx.TxCount1Min > 5 {
		return config.VelocityHighScore, []string{"HIGH_VELOCITY"}
	}
	if ctx.TxCount1Min >= 3 {
		return config.VelocityMediumScore, []string{"MEDIUM_VELOCITY"}
	}
	return 0, nil
}
