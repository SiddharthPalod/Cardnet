package postgres

import "context"

func (r *Repository) RecordTransaction(ctx context.Context, merchantID string, amount float64, approved bool) error {
	approvedInt := 0
	approvedAmount := 0.0
	if approved {
		approvedInt = 1
		approvedAmount = amount
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO merchant_stats (merchant_id, tx_count, approved_count, total_amount)
		VALUES ($1, 1, $2, $3)
		ON CONFLICT (merchant_id)
		DO UPDATE SET
			tx_count = merchant_stats.tx_count + 1,
			approved_count = merchant_stats.approved_count + $2,
			total_amount = merchant_stats.total_amount + $3
	`, merchantID, approvedInt, approvedAmount)
	return err
}
