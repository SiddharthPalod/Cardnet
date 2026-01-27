package reporting

import (
	"context"
	"database/sql"
)

type Generator struct {
	db *sql.DB
}

func NewGenerator(db *sql.DB) *Generator {
	return &Generator{db: db}
}

func (g *Generator) GenerateBatchReport(
	ctx context.Context,
	batchID string,
) (*BatchReport, error) {

	row := g.db.QueryRowContext(
		ctx,
		`
		SELECT
			b.id,
			b.batch_date,
			b.status,
			COUNT(a.authorization_id),
			COALESCE(SUM(a.amount), 0),
			MAX(a.currency)
		FROM settlement_batches b
		LEFT JOIN settlement_authorizations a
			ON a.batch_id = b.id
		WHERE b.id = $1
		GROUP BY b.id, b.batch_date, b.status
		`,
		batchID,
	)

	var r BatchReport
	err := row.Scan(
		&r.BatchID,
		&r.BatchDate,
		&r.Status,
		&r.TotalCount,
		&r.TotalAmount,
		&r.Currency,
	)

	return &r, err
}
