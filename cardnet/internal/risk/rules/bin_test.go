package rules

import (
	"testing"

	"cardnet/internal/risk/models"
)

func TestBINRule(t *testing.T) {
	tests := []struct {
		name     string
		binRisk  int
		expected int
	}{
		{"high risk bin", 90, 10},
		{"medium risk bin", 60, 5},
		{"low risk bin", 20, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := BINRule(
				models.Transaction{},
				models.RiskContext{BINRiskScore: tt.binRisk},
			)

			if score != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, score)
			}
		})
	}
}
