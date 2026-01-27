package rules

import (
	"testing"

	"cardnet/internal/risk/models"
)

func TestGeoRule(t *testing.T) {
	tests := []struct {
		name        string
		lastCountry string
		currCountry string
		timeDiffSec int64
		expected    int
	}{
		{"same country", "IN", "IN", 3600, 0},
		{"geo mismatch", "IN", "US", 3600 * 12, 10}, // 12 hours
		{"impossible travel", "IN", "US", 60, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := GeoRule(
				models.Transaction{
					Country:   tt.currCountry,
					Timestamp: 10_000,
				},
				models.RiskContext{
					LastCountry: tt.lastCountry,
					LastTxTime:  10_000 - tt.timeDiffSec,
				},
			)

			if score != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, score)
			}
		})
	}
}
