package feedback

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Publisher) PublishRiskSuggestion(
	ctx context.Context,
	suggestion RiskSuggestion) error {

	payload, err := json.Marshal(suggestion)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: payload,
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *Publisher) PublishThresholdUpdate(
	ctx context.Context,
	update ThresholdUpdate,
) error {

	payload, err := json.Marshal(update)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Publisher) Close() {
	if err := p.writer.Close(); err != nil {
		log.Println("feedback publisher close error:", err)
	}
}
