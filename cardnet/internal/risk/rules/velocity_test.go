package rules

import (
	"cardnet/internal/risk/models"
	"testing"
)

func TestVelocityRule(t *testing.T) {
	tests := []struct {
		name     string
		txCount  int
		expected int
	}{
		{"high velocity", 6, 30},
		{"medium velocity", 3, 15},
		{"normal", 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := VelocityRule(
				models.Transaction{},
				models.RiskContext{TxCount1Min: tt.txCount},
			)

			if score != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, score)
			}
		})
	}
}
