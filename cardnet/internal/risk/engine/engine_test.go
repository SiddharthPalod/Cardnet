package engine

import (
	"reflect"
	"testing"

	"cardnet/internal/risk/models"
)

func TestRiskEngine_GoldenCases(t *testing.T) {
	re := New()

	tests := []struct {
		name     string
		tx       models.Transaction
		ctx      models.RiskContext
		expected models.RiskResult
	}{
		{
			name: "clean transaction",
			tx: models.Transaction{
				Amount:    1000,
				Country:   "IN",
				Timestamp: 10_000,
			},
			ctx: models.RiskContext{
				TxCount1Min:       1,
				UserAvgAmount:     1200,
				UserP95Amount:     3000,
				MerchantRiskScore: 10,
				BINRiskScore:      10,
				LastCountry:       "IN",
				LastTxTime:        9_000,
			},
			expected: models.RiskResult{
				Decision:  "APPROVE",
				RiskScore: 0,
				Reasons:   nil, // IMPORTANT: engine returns nil slice
			},
		},
		{
			name: "velocity + amount deviation",
			tx: models.Transaction{
				Amount:    6000,
				Country:   "IN",
				Timestamp: 10_000,
			},
			ctx: models.RiskContext{
				TxCount1Min:   6,    // +30
				UserAvgAmount: 1000, // +10
				UserP95Amount: 2000,
				LastCountry:   "IN",
				LastTxTime:    9_000,
			},
			expected: models.RiskResult{
				Decision:  "SOFT_DECLINE",
				RiskScore: 40,
				Reasons:   []string{"HIGH_VELOCITY", "AMOUNT_DEVIATION"},
			},
		},
		{
			name: "geo impossible + risky merchant",
			tx: models.Transaction{
				Amount:    3000,
				Country:   "US",
				Timestamp: 10_000,
			},
			ctx: models.RiskContext{
				LastCountry:       "IN",  // +25
				LastTxTime:        9_940, // ~1 min
				MerchantRiskScore: 80,    // +15
				UserAvgAmount:     5000,  // prevent amount rule
				UserP95Amount:     8000,
			},
			expected: models.RiskResult{
				Decision:  "SOFT_DECLINE",
				RiskScore: 40,
				Reasons:   []string{"IMPOSSIBLE_GEO_TRAVEL", "HIGH_RISK_MERCHANT"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := re.Evaluate(tt.tx, tt.ctx)

			if result.Decision != tt.expected.Decision ||
				result.RiskScore != tt.expected.RiskScore ||
				!reflect.DeepEqual(
					normalizeReasons(result.Reasons),
					normalizeReasons(tt.expected.Reasons),
				) {

				t.Fatalf(
					"\nexpected: %+v\nactual:   %+v",
					tt.expected,
					result,
				)
			}
		})
	}
}

func normalizeReasons(r []string) []string {
	if len(r) == 0 {
		return nil
	}
	return r
}
