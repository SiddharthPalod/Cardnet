package rules

import (
	"cardnet/internal/risk/models"
	"testing"
)

func TestAmountRule(t *testing.T) {
	ctx := models.RiskContext{
		UserAvgAmount: 1000,
		UserP95Amount: 2000,
	}

	tests := []struct {
		amount   int64
		expected int
	}{
		{7000, 20},
		{3000, 10},
		{900, 0},
	}

	for _, tt := range tests {
		score, _ := AmountRule(
			models.Transaction{Amount: tt.amount},
			ctx,
		)

		if score != tt.expected {
			t.Fatalf("expected %d, got %d", tt.expected, score)
		}
	}
}
