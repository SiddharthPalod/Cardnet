package postgres

import (
	"context"
	"database/sql"

	"cardnet/internal/settlement/domain"
)

type SettlementRepository struct {
	db *sql.DB
}

func NewSettlementRepository(db *sql.DB) *SettlementRepository {
	return &SettlementRepository{db: db}
}

func (r *SettlementRepository) Insert(
	ctx context.Context,
	auth domain.SettlementAuthorization,
) error {

	query := `
	INSERT INTO settlement_authorizations (
		authorization_id,
		merchant_id,
		card_hash,
		amount,
		currency,
		approved,
		event_time,
		created_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	ON CONFLICT (authorization_id) DO NOTHING
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		auth.AuthorizationID,
		auth.MerchantID,
		auth.CardHash,
		auth.Amount,
		auth.Currency,
		auth.Approved,
		auth.EventTime,
		auth.CreatedAt,
	)

	return err
}
