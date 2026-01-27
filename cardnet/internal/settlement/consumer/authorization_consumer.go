package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	ledgerpb "cardnet/internal/api/ledger"
	"cardnet/internal/settlement/domain"
	"cardnet/internal/settlement/service"
)

type AuthorizationConsumer struct {
	service *service.SettlementService
}

// ConsumerAuthEvent Matches internal/ledger/domain/auth_event.go JSON structure
type ConsumerAuthEvent struct {
	AuthID     string             `json:"AuthID"`
	MerchantID string             `json:"MerchantID"`
	CardHash   string             `json:"CardHash"`
	EventType  ledgerpb.EventType `json:"EventType"`
	Decision   ledgerpb.Decision  `json:"Decision"`
	Payload    []byte             `json:"Payload"`
	CreatedAt  time.Time          `json:"CreatedAt"`
}

func NewAuthorizationConsumer(
	svc *service.SettlementService,
) *AuthorizationConsumer {
	return &AuthorizationConsumer{
		service: svc,
	}
}

func (c *AuthorizationConsumer) HandleMessage(
	ctx context.Context,
	msg []byte,
) error {

	var event ConsumerAuthEvent
	// 1. Unmarshal the Kafka message (JSON of AuthEvent)
	if err := json.Unmarshal(msg, &event); err != nil {
		log.Printf("Failed to unmarshal Kafka message: %v", err)
		return nil // Skip malformed messages
	}

	// We only care about FINAL_OUTCOME for settlement
	// Other event types (AUTH_REQUEST, RISK_DECISION, ISSUER_RESPONSE) are expected
	// and should be silently ignored - no logging needed as this is normal behavior
	if event.EventType != ledgerpb.EventType_FINAL_OUTCOME {
		// Silently ignore - this is expected behavior, not an error
		// No logging to reduce noise - settlement service only processes FINAL_OUTCOME
		return nil
	}
	log.Printf("Processing FINAL_OUTCOME for AuthID: %s", event.AuthID)

	var payload domain.AuthorizationPayload

	// 2. Unmarshal the inner payload (JSON of FinalOutcomePayload)
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Printf("Failed to unmarshal event payload: %v", err)
		return nil
	}

	return c.service.ProcessAuthorization(ctx, domain.SettlementAuthorization{
		AuthorizationID: event.AuthID,
		MerchantID:      event.MerchantID,
		CardHash:        event.CardHash,
		Amount:          payload.Amount,
		Currency:        payload.Currency,
		Approved:        payload.Status == "APPROVED", // from auth service
		EventTime:       event.CreatedAt,
		CreatedAt:       time.Now(),
	})
}
