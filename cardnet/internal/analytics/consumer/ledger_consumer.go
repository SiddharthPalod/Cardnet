package consumer

import (
	"cardnet/internal/analytics/aggregations"
	"cardnet/internal/analytics/config"
	"cardnet/internal/analytics/observability"
	"cardnet/internal/analytics/repository/postgres"
	ledgerpb "cardnet/internal/api/ledger"
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type ConsumerAuthEvent struct {
	AuthID     string             `json:"AuthID"`
	MerchantID string             `json:"MerchantID"`
	CardHash   string             `json:"CardHash"`
	EventType  ledgerpb.EventType `json:"EventType"`
	Decision   ledgerpb.Decision  `json:"Decision"`
	Payload    []byte             `json:"Payload"`
	CreatedAt  time.Time          `json:"CreatedAt"`
}

type AuthorizationPayload struct {
	Amount   interface{} `json:"amount"` // Can be int64 (cents) or float64 (dollars)
	Currency string      `json:"currency"`
	Status   string      `json:"status"`
	Reason   string      `json:"reason,omitempty"`
}

type LedgerConsumer struct {
	reader *kafka.Reader
	repo   *postgres.Repository
}

func NewLedgerConsumer(cfg config.KafkaConfig, repo *postgres.Repository) *LedgerConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   "auth.events",
	})
	return &LedgerConsumer{reader: r, repo: repo}
}

func (c *LedgerConsumer) Start() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("ledger read error:", err)
			continue
		}

		observability.IncLedgerMessages()

		var event ConsumerAuthEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("ledger decode error:", err)
			continue
		}

		if event.EventType != ledgerpb.EventType_FINAL_OUTCOME {
			continue
		}

		var payload AuthorizationPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			log.Println("payload decode error:", err)
			continue
		}

		// Use a detached context or short timeout for DB operations to avoid blocking consumer indefinitely
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// Convert amount to float64 (handle both int64 cents and float64 dollars)
		var amountFloat float64
		switch v := payload.Amount.(type) {
		case float64:
			amountFloat = v
		case int64:
			amountFloat = float64(v) / 100.0 // Convert cents to dollars
		case int:
			amountFloat = float64(v) / 100.0 // Convert cents to dollars
		default:
			// Try to convert via JSON number
			if num, ok := v.(json.Number); ok {
				if f, err := num.Float64(); err == nil {
					amountFloat = f
				} else if i, err := num.Int64(); err == nil {
					amountFloat = float64(i) / 100.0
				}
			}
		}

		approved := payload.Status == "APPROVED"
		
		// Record network-level metrics
		if err := c.repo.RecordDecision(ctx, approved); err != nil {
			log.Println("repo record network metrics error:", err)
		}
		
		// Record merchant-level metrics
		if err := c.repo.RecordTransaction(ctx, event.MerchantID, amountFloat, approved); err != nil {
			log.Println("repo record merchant stats error:", err)
		}

		// Extract BIN from CardHash (which is the card token, not a hash)
		// Use routing.ExtractBIN to handle dashes/spaces properly
		if len(event.CardHash) >= 6 {
			// Normalize card number by removing non-digits, then extract first 6 digits
			bin := extractBINFromCard(event.CardHash)
			if len(bin) >= 6 {
				if err := c.repo.RecordBINTransaction(ctx, bin, amountFloat, approved); err != nil {
					log.Println("repo record bin stats error:", err)
				}
			}
		}

		// Record raw history row for analytics UI
		historyItem := &aggregations.NetworkHistoryItem{
			AuthID:     event.AuthID,
			MerchantID: event.MerchantID,
			Amount:     amountFloat,
			Currency:   payload.Currency,
			Status:     payload.Status,
			Reason:     payload.Reason,
			CreatedAt:  event.CreatedAt,
		}
		if err := c.repo.RecordNetworkHistory(ctx, historyItem); err != nil {
			log.Println("repo record network history error:", err)
		}

		cancel()
	}
}

// extractBINFromCard extracts the first 6 digits from a card number (handles dashes/spaces)
func extractBINFromCard(cardNumber string) string {
	var digits strings.Builder
	for _, r := range cardNumber {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	normalized := digits.String()
	if len(normalized) >= 6 {
		return normalized[:6]
	}
	return normalized
}
