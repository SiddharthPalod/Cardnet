package rules

import (
	"cardnet/internal/risk/config"
	"cardnet/internal/risk/models"
)

func GeoRule(tx models.Transaction, ctx models.RiskContext) (int, []string) {
	if ctx.LastCountry == "" || ctx.LastCountry == tx.Country {
		return 0, nil
	}

	timeDiffHours := float64(tx.Timestamp-ctx.LastTxTime) / 3600
	if timeDiffHours <= 0 {
		return config.GeoImpossibleScore, []string{"GEO_TIME_ANOMALY"}
	}
	
	distanceKm := 8000.0

	speed := distanceKm / timeDiffHours

	if speed > config.GeoImpossibleSpeedKmph {
		return config.GeoImpossibleScore, []string{"IMPOSSIBLE_GEO_TRAVEL"}
	}

	return config.GeoSuspiciousScore, []string{"GEO_MISMATCH"}
}
