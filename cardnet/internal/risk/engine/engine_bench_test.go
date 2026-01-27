package engine

import (
	"cardnet/internal/risk/models"
	"testing"
)

func BenchmarkRiskEngine(b *testing.B) {
	re := New()

	tx := models.Transaction{
		Amount:    5000,
		Country:   "IN",
		Timestamp: 1700000000,
	}

	ctx := models.RiskContext{
		TxCount1Min:       2,
		UserAvgAmount:     1200,
		UserP95Amount:     3000,
		MerchantRiskScore: 30,
		BINRiskScore:      20,
		LastCountry:       "IN",
		LastTxTime:        1699999900,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		re.Evaluate(tx, ctx)
	}
}
