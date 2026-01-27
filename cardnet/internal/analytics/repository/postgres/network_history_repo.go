package postgres

import (
	"cardnet/internal/analytics/aggregations"
	"context"
	"fmt"
	"strings"
)

// GetNetworkHistory returns recent network-wide authorization history with optional filters.
func (r *Repository) GetNetworkHistory(ctx context.Context, merchantID, status string, limit int32) ([]*aggregations.NetworkHistoryItem, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	baseQuery := `
		SELECT auth_id, merchant_id, amount, currency, status, COALESCE(reason, ''), created_at
		FROM network_history
	`

	args := []any{}
	clauses := []string{}
	argIdx := 1

	if merchantID != "" {
		clauses = append(clauses, fmt.Sprintf("merchant_id = $%d", argIdx))
		args = append(args, merchantID)
		argIdx++
	}
	if status != "" {
		clauses = append(clauses, fmt.Sprintf("UPPER(status) = UPPER($%d)", argIdx))
		args = append(args, status)
		argIdx++
	}

	query := baseQuery
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*aggregations.NetworkHistoryItem
	for rows.Next() {
		var item aggregations.NetworkHistoryItem
		if err := rows.Scan(
			&item.AuthID,
			&item.MerchantID,
			&item.Amount,
			&item.Currency,
			&item.Status,
			&item.Reason,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, &item)
	}

	return out, rows.Err()
}
