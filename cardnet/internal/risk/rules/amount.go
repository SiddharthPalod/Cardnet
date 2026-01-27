package rules

import (
	"cardnet/internal/risk/config"
	"cardnet/internal/risk/models"
)

func AmountRule(tx models.Transaction, ctx models.RiskContext) (int, []string) {
	if tx.Amount > 3*ctx.UserP95Amount {
		return config.AmountHighScore, []string{"AMOUNT_OUTLIER"}
	}
	if tx.Amount > 2*ctx.UserAvgAmount {
		return config.AmountMediumScore, []string{"AMOUNT_DEVIATION"}
	}
	return 0, nil
}
