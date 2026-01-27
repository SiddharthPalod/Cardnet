package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByRequestID(
	ctx context.Context,
	requestID string,
) (string, string, error) {

	query := `
		SELECT auth_id, status
		FROM authorizations
		WHERE request_id = $1
	`

	var authID, status string
	err := r.db.QueryRowContext(ctx, query, requestID).Scan(&authID, &status)
	if err != nil {
		return "", "", err
	}

	return authID, status, nil
}

func (r *PostgresRepository) SaveAuthorization(
	ctx context.Context,
	authID string,
	requestID string,
	merchantID string,
	cardToken string,
	amount int64,
	status string,
) error {

	query := `
		INSERT INTO authorizations
		(auth_id, request_id, merchant_id, card_token, amount, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		authID,
		requestID,
		merchantID,
		cardToken,
		amount,
		status,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrDuplicateRequest
		}
		return err
	}

	return nil
}

func (r *PostgresRepository) SaveRiskAudit(
	ctx context.Context,
	authID string,
	riskScore int,
	reasons []string,
) error {
	query := `
		INSERT INTO risk_audit (
			auth_id,
			risk_score,
			reasons
		) VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		authID,
		riskScore,
		reasons,
	)
	return err
}
