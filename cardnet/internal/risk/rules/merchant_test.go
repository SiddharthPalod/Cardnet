package rules

import (
	"testing"

	"cardnet/internal/risk/models"
)

func TestMerchantRule(t *testing.T) {
	tests := []struct {
		name     string
		risk     int
		expected int
	}{
		{"high risk merchant", 80, 15},
		{"medium risk merchant", 50, 8},
		{"low risk merchant", 20, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := MerchantRule(
				models.Transaction{},
				models.RiskContext{MerchantRiskScore: tt.risk},
			)

			if score != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, score)
			}
		})
	}
}
