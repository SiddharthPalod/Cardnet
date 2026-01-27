package features

import (
	"context"
	"database/sql"
)

type MerchantFeatures struct {
	MerchantID      string
	TxCount24h      int
	ApprovalRate24h float64
	AvgAmount24h    float64
	DeclineRate24h  float64
}

func BuildMerchantFeatures(
	ctx context.Context,
	db *sql.DB,
) ([]MerchantFeatures, error) {

	rows, err := db.QueryContext(ctx, `
		SELECT
			merchant_id,
			tx_count,
			CAST(approved_count AS FLOAT) / NULLIF(tx_count, 0) as approval_rate,
			total_amount / NULLIF(tx_count, 0) as avg_amount
		FROM merchant_stats
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MerchantFeatures
	for rows.Next() {
		var f MerchantFeatures
		var rate, avg sql.NullFloat64
		if err := rows.Scan(
			&f.MerchantID,
			&f.TxCount24h,
			&rate,
			&avg,
		); err != nil {
			return nil, err
		}
		if rate.Valid {
			f.ApprovalRate24h = rate.Float64
			f.DeclineRate24h = 1.0 - f.ApprovalRate24h
		}
		if avg.Valid {
			f.AvgAmount24h = avg.Float64
		}
		result = append(result, f)
	}
	return result, nil
}
