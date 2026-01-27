package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type BatchRepository struct {
	db *sql.DB
}

func NewBatchRepository(db *sql.DB) *BatchRepository {
	return &BatchRepository{db: db}
}

func (r *BatchRepository) GetOrCreateBatch(
	ctx context.Context,
	date time.Time,
) (string, error) {

	var batchID string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id FROM settlement_batches WHERE batch_date = $1`,
		date,
	).Scan(&batchID)

	if err == nil {
		return batchID, nil
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	batchID = uuid.New().String()

	_, err = r.db.ExecContext(
		ctx,
		`
		INSERT INTO settlement_batches (id, batch_date, status)
		VALUES ($1,$2,$3)
		ON CONFLICT (batch_date) DO NOTHING
		`,
		batchID,
		date,
		"OPEN",
	)

	if err != nil {
		return "", err
	}

	return batchID, nil
}

func (r *BatchRepository) AttachAuthorizations(
	ctx context.Context,
	batchID string,
	date time.Time,
) error {
	nextDay := date.Add(24 * time.Hour)

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE settlement_authorizations
		SET batch_id = $1
		WHERE batch_id IS NULL
		  AND event_time >= $2
		  AND event_time < $3
		`,
		batchID,
		date,
		nextDay,
	)

	return err
}

func (r *BatchRepository) MarkReady(
	ctx context.Context,
	batchID string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE settlement_batches
		SET status = 'READY'
		WHERE id = $1
		`,
		batchID,
	)

	return err
}
