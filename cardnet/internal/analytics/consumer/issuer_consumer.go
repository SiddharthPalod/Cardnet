package consumer

import (
	"cardnet/internal/analytics/config"
	"cardnet/internal/analytics/observability"
	"cardnet/internal/analytics/repository/postgres"
	ledgerpb "cardnet/internal/api/ledger"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type IssuerConsumer struct {
	reader *kafka.Reader
	repo   *postgres.Repository
}

func NewIssuerConsumer(cfg config.KafkaConfig, repo *postgres.Repository) *IssuerConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   "issuer.responses",
	})
	return &IssuerConsumer{reader: r, repo: repo}
}

func (c *IssuerConsumer) Start() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("issuer read error:", err)
			continue
		}

		observability.IncIssuerMessages()

		var event ConsumerAuthEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("issuer decode error:", err)
			continue
		}

		if event.EventType != ledgerpb.EventType_ISSUER_RESPONSE {
			continue
		}

		approved := event.Decision == ledgerpb.Decision_APPROVED

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := c.repo.RecordDecision(ctx, approved); err != nil {
			log.Println("repo record decision error:", err)
		}
		cancel()
	}
}
