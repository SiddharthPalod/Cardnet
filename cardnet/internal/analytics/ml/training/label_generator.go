package training

func GenerateMerchantLabels(
	approvalRates map[string]float64,
) map[string]int {

	labels := make(map[string]int)

	for merchant, rate := range approvalRates {
		if rate < 0.7 {
			labels[merchant] = 1 // risky
		} else {
			labels[merchant] = 0
		}
	}
	return labels
}
