package postgres

import (
	"cardnet/internal/analytics/aggregations"
	"context"
	"database/sql"
)

func (r *Repository) GetMerchantStats(ctx context.Context, merchantID string) (*aggregations.MerchantStats, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT merchant_id, tx_count, approved_count, total_amount
		FROM merchant_stats
		WHERE merchant_id = $1
	`, merchantID)

	var s aggregations.MerchantStats
	if err := row.Scan(&s.MerchantID, &s.TxCount, &s.ApprovedTx, &s.TotalAmount); err != nil {
		if err == sql.ErrNoRows {
			return &aggregations.MerchantStats{MerchantID: merchantID}, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetBinStats(ctx context.Context, bin string) (*aggregations.BINStats, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT bin, tx_count, approved_count, total_amount
		FROM bin_stats
		WHERE bin = $1
	`, bin)

	var s aggregations.BINStats
	if err := row.Scan(&s.BIN, &s.TxCount, &s.ApprovedTx, &s.TotalAmount); err != nil {
		if err == sql.ErrNoRows {
			return &aggregations.BINStats{BIN: bin}, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetRuleStats(ctx context.Context, ruleID string) (*aggregations.RuleStats, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT rule_id, triggered_count, decline_count
		FROM rule_stats
		WHERE rule_id = $1
	`, ruleID)

	var s aggregations.RuleStats
	if err := row.Scan(&s.RuleID, &s.TriggeredCount, &s.DeclineCount); err != nil {
		if err == sql.ErrNoRows {
			return &aggregations.RuleStats{RuleID: ruleID}, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetNetworkMetrics(ctx context.Context) (*aggregations.ApprovalMetrics, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT total_tx, approved_tx
		FROM network_metrics
		LIMIT 1
	`)

	var s aggregations.ApprovalMetrics
	if err := row.Scan(&s.TotalTx, &s.ApprovedTx); err != nil {
		if err == sql.ErrNoRows {
			return &aggregations.ApprovalMetrics{}, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetAllMerchantStats(ctx context.Context) ([]*aggregations.MerchantStats, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT merchant_id, tx_count, approved_count, total_amount
		FROM merchant_stats
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*aggregations.MerchantStats
	for rows.Next() {
		var s aggregations.MerchantStats
		if err := rows.Scan(&s.MerchantID, &s.TxCount, &s.ApprovedTx, &s.TotalAmount); err != nil {
			return nil, err
		}
		results = append(results, &s)
	}
	return results, nil
}

func (r *Repository) GetAllBinStats(ctx context.Context) ([]*aggregations.BINStats, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT bin, tx_count, approved_count, total_amount
		FROM bin_stats
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*aggregations.BINStats
	for rows.Next() {
		var s aggregations.BINStats
		if err := rows.Scan(&s.BIN, &s.TxCount, &s.ApprovedTx, &s.TotalAmount); err != nil {
			return nil, err
		}
		results = append(results, &s)
	}
	return results, nil
}