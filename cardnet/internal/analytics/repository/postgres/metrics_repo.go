package postgres

import (
	"cardnet/internal/analytics/aggregations"
	"context"
)

func (r *Repository) RecordDecision(ctx context.Context, approved bool) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if approved {
		if _, err := tx.ExecContext(ctx, `UPDATE network_metrics SET approved_tx = approved_tx + 1`); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE network_metrics SET total_tx = total_tx + 1`); err != nil {
		return err
	}

	return tx.Commit()
}

// RecordNetworkHistory stores a single final outcome authorization event for later querying.
func (r *Repository) RecordNetworkHistory(ctx context.Context, item *aggregations.NetworkHistoryItem) error {
	const query = `
		INSERT INTO network_history (auth_id, merchant_id, amount, currency, status, reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (auth_id) DO UPDATE
		SET merchant_id = EXCLUDED.merchant_id,
		    amount = EXCLUDED.amount,
		    currency = EXCLUDED.currency,
		    status = EXCLUDED.status,
		    reason = EXCLUDED.reason,
		    created_at = EXCLUDED.created_at
	`

	_, err := r.db.ExecContext(ctx, query,
		item.AuthID,
		item.MerchantID,
		item.Amount,
		item.Currency,
		item.Status,
		item.Reason,
		item.CreatedAt,
	)
	return err
}
