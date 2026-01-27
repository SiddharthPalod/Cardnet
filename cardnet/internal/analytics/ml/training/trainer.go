package training

import "cardnet/internal/analytics/ml/models"

func TrainMerchantRiskModel(
	rows []TrainingRow,
) *models.MerchantRiskModel {

	// Placeholder weights (v1 baseline)
	return &models.MerchantRiskModel{
		Version: "v1",
		Weights: map[string]float64{
			"tx_count":      -0.01,
			"approval_rate": -2.0,
			"avg_amount":    0.001,
		},
		Bias: 0.5,
	}
}
