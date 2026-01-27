package features

import (
	"context"
	"database/sql"
)

type BINFeatures struct {
	BIN             string
	TxCount24h      int
	ApprovalRate24h float64
}

func BuildBINFeatures(ctx context.Context, db *sql.DB) ([]BINFeatures, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT
			bin,
			tx_count,
			CAST(approved_count AS FLOAT) / NULLIF(tx_count, 0) as approval_rate
		FROM bin_stats
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bins []BINFeatures
	for rows.Next() {
		var b BINFeatures
		var rate sql.NullFloat64
		if err := rows.Scan(&b.BIN, &b.TxCount24h, &rate); err != nil {
			return nil, err
		}
		if rate.Valid {
			b.ApprovalRate24h = rate.Float64
		}
		bins = append(bins, b)
	}

	return bins, nil
}
