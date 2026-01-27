package postgres

import (
	ledgerpb "cardnet/internal/api/ledger"
	"context"
)

func (r *Repository) RecordRuleHit(ctx context.Context, ruleID string, decision ledgerpb.Decision) error {
	isDecline := 0
	if decision == ledgerpb.Decision_DECLINED || decision == ledgerpb.Decision_REVIEW {
		isDecline = 1
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO rule_stats (rule_id, triggered_count, decline_count)
		VALUES ($1, 1, $2)
		ON CONFLICT (rule_id)
		DO UPDATE SET
			triggered_count = rule_stats.triggered_count + 1,
			decline_count = rule_stats.decline_count + $2
	`, ruleID, isDecline)
	return err
}
