package postgres

import (
	"context"
	"database/sql"
)

type ReconciliationRepository struct {
	db *sql.DB
}

func NewReconciliationRepository(db *sql.DB) *ReconciliationRepository {
	return &ReconciliationRepository{db: db}
}

func (r *ReconciliationRepository) FindDuplicates(
	ctx context.Context,
	batchID string,
) ([]string, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT authorization_id
		FROM settlement_authorizations
		WHERE batch_id = $1
		GROUP BY authorization_id
		HAVING COUNT(*) > 1
		`,
		batchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dups []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		dups = append(dups, id)
	}

	return dups, nil
}

func (r *ReconciliationRepository) HasCriticalNulls(
	ctx context.Context,
	batchID string,
) (bool, error) {

	var count int
	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT COUNT(*)
		FROM settlement_authorizations
		WHERE batch_id = $1
		  AND (amount IS NULL OR currency IS NULL)
		`,
		batchID,
	).Scan(&count)

	return count > 0, err
}

func (r *ReconciliationRepository) CountAuthorizations(
	ctx context.Context,
	batchID string,
) (int64, error) {

	var count int64
	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT COUNT(*)
		FROM settlement_authorizations
		WHERE batch_id = $1
		`,
		batchID,
	).Scan(&count)

	return count, err
}

func (r *ReconciliationRepository) SaveIssue(
	ctx context.Context,
	batchID, reason string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO settlement_reconciliation_issues (batch_id, reason)
		VALUES ($1,$2)
		`,
		batchID,
		reason,
	)

	return err
}

func (r *ReconciliationRepository) MarkBatchStatus(
	ctx context.Context,
	batchID, status string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE settlement_batches
		SET reconciliation_status = $2
		WHERE id = $1
		`,
		batchID,
		status,
	)

	return err
}
