package models

type MerchantRiskModel struct {
	Version string
	Weights map[string]float64
	Bias    float64
}
