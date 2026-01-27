package rules

import "cardnet/internal/risk/models"

type Rule func(
	tx models.Transaction,
	ctx models.RiskContext,
) (score int, reasons []string)
