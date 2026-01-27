package domain

import (
	ledgerpb "cardnet/internal/api/ledger"
	"time"

	"github.com/gocql/gocql"
)

type AuthEvent struct {
	EventID    gocql.UUID
	AuthID     string
	MerchantID string
	CardHash   string
	EventType  ledgerpb.EventType
	Decision   ledgerpb.Decision
	Payload    []byte
	CreatedAt  time.Time
}

func FromProto(p *ledgerpb.AuthEvent) *AuthEvent {
	id, _ := gocql.ParseUUID(p.EventId)

	return &AuthEvent{
		EventID:    id,
		AuthID:     p.AuthId,
		MerchantID: p.MerchantId,
		CardHash:   p.CardHash,
		EventType:  p.EventType,
		Decision:   p.Decision,
		Payload:    p.Payload,
		CreatedAt:  time.Unix(p.CreatedAt, 0),
	}
}
