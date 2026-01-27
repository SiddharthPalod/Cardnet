package main

import (
	ledgerpb "cardnet/internal/api/ledger"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func main() {
	payload, _ := json.Marshal(map[string]any{
		"amount":   1000,
		"currency": "USD",
	})

	event := &ledgerpb.AuthEvent{
		EventId:    "11111111-1111-1111-1111-111111111111",
		AuthId:     "auth-test-001",
		MerchantId: "merchant-123",
		CardHash:   "hash-abc",
		EventType:  ledgerpb.EventType_AUTH_REQUEST,
		Decision:   ledgerpb.Decision_APPROVED,
		Payload:    payload,
		CreatedAt:  time.Now().Unix(),
	}

	bytes, _ := proto.Marshal(event)

	w := kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "ledger-events",
	}

	_ = w.WriteMessages(context.Background(), kafka.Message{
		Value: bytes,
	})
}
