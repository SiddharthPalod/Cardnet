package cassandra

import (
	"cardnet/internal/ledger/domain"
	"context"
	"time"

	"github.com/gocql/gocql"
)

type LedgerRepo struct {
	session *gocql.Session
}

func NewLedgerRepo(session *gocql.Session) *LedgerRepo {
	return &LedgerRepo{session: session}
}

func (r *LedgerRepo) Store(ctx context.Context, e *domain.AuthEvent) error {
	batch := r.session.NewBatch(gocql.LoggedBatch)

	batch.Query(insertByMerchant,
		e.MerchantID, e.CreatedAt, e.AuthID, e.EventID,
		e.EventType.String(), e.Decision.String(), e.Payload,
	)

	batch.Query(insertByCard,
		e.CardHash, e.CreatedAt, e.AuthID, e.MerchantID,
		e.EventID, e.EventType.String(), e.Decision.String(), e.Payload,
	)

	batch.Query(insertByTime,
		e.CreatedAt.Truncate(24*time.Hour),
		e.CreatedAt, e.AuthID, e.MerchantID,
		e.CardHash, e.EventType.String(), e.Decision.String(),
		e.Payload, e.EventID,
	)

	return r.session.ExecuteBatch(batch)
}
