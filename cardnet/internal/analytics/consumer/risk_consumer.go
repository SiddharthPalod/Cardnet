package consumer

import (
	"cardnet/internal/analytics/config"
	"cardnet/internal/analytics/observability"
	"cardnet/internal/analytics/repository/postgres"
	ledgerpb "cardnet/internal/api/ledger"
	riskpb "cardnet/internal/api/risk"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type ConsumerRiskResponse struct {
	Decision  riskpb.Decision `json:"decision"`
	RiskScore int32           `json:"risk_score"`
	Reasons   []string        `json:"reasons"`
}

type RiskConsumer struct {
	reader *kafka.Reader
	repo   *postgres.Repository
}

func NewRiskConsumer(cfg config.KafkaConfig, repo *postgres.Repository) *RiskConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   "risk.decisions",
	})
	return &RiskConsumer{reader: r, repo: repo}
}

func (c *RiskConsumer) Start() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("risk read error:", err)
			continue
		}

		observability.IncRiskMessages()

		var event ConsumerAuthEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("risk decode error:", err)
			continue
		}

		if event.EventType != ledgerpb.EventType_RISK_DECISION {
			continue
		}

		var payload ConsumerRiskResponse
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			log.Println("risk payload decode error:", err)
			continue
		}

		// Map decision to Ledger Decision for repository
		var ledgerDecision ledgerpb.Decision
		switch payload.Decision {
		case riskpb.Decision_APPROVE:
			ledgerDecision = ledgerpb.Decision_APPROVED
		case riskpb.Decision_SOFT_DECLINE:
			ledgerDecision = ledgerpb.Decision_REVIEW
		case riskpb.Decision_HARD_DECLINE:
			ledgerDecision = ledgerpb.Decision_DECLINED
		default:
			ledgerDecision = ledgerpb.Decision_DECISION_UNSPECIFIED
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		for _, rule := range payload.Reasons {
			if err := c.repo.RecordRuleHit(ctx, rule, ledgerDecision); err != nil {
				log.Printf("repo record rule hit error (rule: %s): %v", rule, err)
			}
		}
		cancel()

		log.Printf("Risk decision for tx %s: %s (Score: %d, Reasons: %v)",
			event.AuthID, payload.Decision, payload.RiskScore, payload.Reasons)
	}
}
