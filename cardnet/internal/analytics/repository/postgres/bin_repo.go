package postgres

import "context"

func (r *Repository) RecordBINTransaction(ctx context.Context, bin string, amount float64, approved bool) error {
	approvedInt := 0
	approvedAmount := 0.0
	if approved {
		approvedInt = 1
		approvedAmount = amount
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO bin_stats (bin, tx_count, approved_count, total_amount)
		VALUES ($1, 1, $2, $3)
		ON CONFLICT (bin)
		DO UPDATE SET
			tx_count = bin_stats.tx_count + 1,
			approved_count = bin_stats.approved_count + $2,
			total_amount = bin_stats.total_amount + $3
	`, bin, approvedInt, approvedAmount)
	return err
}
