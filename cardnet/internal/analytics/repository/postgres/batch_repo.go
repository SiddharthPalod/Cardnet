package postgres

import (
	"context"
	"time"
)

// RecordHourlyMerchantAggregation records hourly aggregation for a merchant
func (r *Repository) RecordHourlyMerchantAggregation(
	ctx context.Context,
	merchantID string,
	hourBucket time.Time,
	txCount, approvedCount int,
	totalAmount float64,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO hourly_merchant_aggregations 
		(merchant_id, hour_bucket, tx_count, approved_count, total_amount)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (merchant_id, hour_bucket)
		DO UPDATE SET
			tx_count = EXCLUDED.tx_count,
			approved_count = EXCLUDED.approved_count,
			total_amount = EXCLUDED.total_amount,
			created_at = NOW()
	`, merchantID, hourBucket, txCount, approvedCount, totalAmount)
	return err
}

// RecordHourlyBINAggregation records hourly aggregation for a BIN
func (r *Repository) RecordHourlyBINAggregation(
	ctx context.Context,
	bin string,
	hourBucket time.Time,
	txCount, approvedCount int,
	totalAmount float64,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO hourly_bin_aggregations 
		(bin, hour_bucket, tx_count, approved_count, total_amount)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (bin, hour_bucket)
		DO UPDATE SET
			tx_count = EXCLUDED.tx_count,
			approved_count = EXCLUDED.approved_count,
			total_amount = EXCLUDED.total_amount,
			created_at = NOW()
	`, bin, hourBucket, txCount, approvedCount, totalAmount)
	return err
}

// RecordHourlyNetworkAggregation records hourly network-wide aggregation
func (r *Repository) RecordHourlyNetworkAggregation(
	ctx context.Context,
	hourBucket time.Time,
	totalTx, approvedTx int,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO hourly_network_aggregations 
		(hour_bucket, total_tx, approved_tx)
		VALUES ($1, $2, $3)
		ON CONFLICT (hour_bucket)
		DO UPDATE SET
			total_tx = EXCLUDED.total_tx,
			approved_tx = EXCLUDED.approved_tx,
			created_at = NOW()
	`, hourBucket, totalTx, approvedTx)
	return err
}

// RecordDailyMerchantAggregation records daily aggregation for a merchant
func (r *Repository) RecordDailyMerchantAggregation(
	ctx context.Context,
	merchantID string,
	dayBucket time.Time,
	txCount, approvedCount int,
	totalAmount float64,
) error {
	dayDate := dayBucket.Truncate(24 * time.Hour)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO daily_merchant_aggregations 
		(merchant_id, day_bucket, tx_count, approved_count, total_amount)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (merchant_id, day_bucket)
		DO UPDATE SET
			tx_count = EXCLUDED.tx_count,
			approved_count = EXCLUDED.approved_count,
			total_amount = EXCLUDED.total_amount,
			created_at = NOW()
	`, merchantID, dayDate, txCount, approvedCount, totalAmount)
	return err
}

// RecordDailyBINAggregation records daily aggregation for a BIN
func (r *Repository) RecordDailyBINAggregation(
	ctx context.Context,
	bin string,
	dayBucket time.Time,
	txCount, approvedCount int,
	totalAmount float64,
) error {
	dayDate := dayBucket.Truncate(24 * time.Hour)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO daily_bin_aggregations 
		(bin, day_bucket, tx_count, approved_count, total_amount)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (bin, day_bucket)
		DO UPDATE SET
			tx_count = EXCLUDED.tx_count,
			approved_count = EXCLUDED.approved_count,
			total_amount = EXCLUDED.total_amount,
			created_at = NOW()
	`, bin, dayDate, txCount, approvedCount, totalAmount)
	return err
}

// RecordDailyNetworkAggregation records daily network-wide aggregation
func (r *Repository) RecordDailyNetworkAggregation(
	ctx context.Context,
	dayBucket time.Time,
	totalTx, approvedTx int,
) error {
	dayDate := dayBucket.Truncate(24 * time.Hour)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO daily_network_aggregations 
		(day_bucket, total_tx, approved_tx)
		VALUES ($1, $2, $3)
		ON CONFLICT (day_bucket)
		DO UPDATE SET
			total_tx = EXCLUDED.total_tx,
			approved_tx = EXCLUDED.approved_tx,
			created_at = NOW()
	`, dayDate, totalTx, approvedTx)
	return err
}
