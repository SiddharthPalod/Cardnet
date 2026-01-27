package kafka

import (
	"fmt"
	ledgerpb "cardnet/internal/api/ledger"
	"cardnet/internal/ledger/domain"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	// Initialize topics on startup
	go ensureTopicsExist(brokers)

	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
			Async:    true, // fire-and-forget
		},
	}
}

// ensureTopicsExist creates Kafka topics if they don't exist
// RedPanda auto-creates topics, but we try to create them explicitly for better control
func ensureTopicsExist(brokers []string) {
	topics := []string{"auth.events", "risk.decisions", "issuer.responses", "settlement.ready"}
	
	// Retry logic: try up to 5 times with exponential backoff
	maxRetries := 5
	baseDelay := 2 * time.Second
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := baseDelay * time.Duration(1<<uint(attempt-1)) // Exponential backoff
			log.Printf("Retrying topic creation (attempt %d/%d) after %v...", attempt+1, maxRetries, delay)
			time.Sleep(delay)
		}
		
		// Try to connect to the broker
		conn, err := kafka.Dial("tcp", brokers[0])
		if err != nil {
			if attempt == maxRetries-1 {
				log.Printf("WARN: Failed to connect to Kafka/RedPanda after %d attempts: %v (topics will be auto-created on first write)", maxRetries, err)
			}
			continue
		}
		
		// Get controller (for Kafka) or use connection directly (for RedPanda)
		controller, err := conn.Controller()
		if err != nil {
			conn.Close()
			if attempt == maxRetries-1 {
				log.Printf("WARN: Failed to get controller after %d attempts: %v (topics will be auto-created on first write)", maxRetries, err)
			}
			continue
		}
		
		controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
		conn.Close()
		if err != nil {
			if attempt == maxRetries-1 {
				log.Printf("WARN: Failed to connect to controller after %d attempts: %v (topics will be auto-created on first write)", maxRetries, err)
			}
			continue
		}
		
		// Try to create topics
		allCreated := true
		for _, topic := range topics {
			topicConfig := kafka.TopicConfig{
				Topic:             topic,
				NumPartitions:     3,
				ReplicationFactor: 1,
			}
			
			err := controllerConn.CreateTopics(topicConfig)
			if err != nil {
				// Check if error is "topic already exists" - that's fine
				errStr := err.Error()
				if errStr != "topic already exists" && !contains(errStr, "already exists") {
					log.Printf("WARN: Failed to create topic %s: %v", topic, err)
					allCreated = false
				} else {
					log.Printf("Topic %s already exists, skipping", topic)
				}
			} else {
				log.Printf("✅ Created Kafka topic: %s", topic)
			}
		}
		
		controllerConn.Close()
		
		// If we got here and created topics, we're done
		if allCreated {
			log.Printf("✅ All Kafka topics initialized successfully")
			return
		}
	}
	
	log.Printf("WARN: Topic creation completed with some failures, but topics will be auto-created on first write")
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (p *Producer) Publish(ctx context.Context, event *domain.AuthEvent) error {
	topic := topicForEvent(event)

	if topic == "" {
		return nil
	}

	value, _ := json.Marshal(event)

	// CRITICAL FIX: Include Topic in message (writer doesn't have default topic)
	// Use a short timeout to avoid blocking if Kafka is down
	publishCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := p.writer.WriteMessages(publishCtx, kafka.Message{
		Topic: topic, // Required - writer has no default topic
		Key:   []byte(event.AuthID),
		Value: value,
	})

	// Log errors but don't fail the request (fire-and-forget)
	// RedPanda auto-creates topics, so "Unknown Topic" errors should be transient
	if err != nil {
		errStr := err.Error()
		// Suppress "Unknown Topic Or Partition" errors - they're transient and will resolve
		// when RedPanda auto-creates the topic
		if !contains(errStr, "Unknown Topic") && !contains(errStr, "Unknown Partition") {
			log.Printf("kafka publish failed: %v", err)
		}
		// Return nil to not block the auth path
		return nil
	}

	return nil
}

func topicForEvent(e *domain.AuthEvent) string {

	switch e.EventType {
	case ledgerpb.EventType_AUTH_REQUEST,
		ledgerpb.EventType_FINAL_OUTCOME:
		return "auth.events"
	case ledgerpb.EventType_RISK_DECISION:
		return "risk.decisions"
	case ledgerpb.EventType_ISSUER_RESPONSE:
		return "issuer.responses"
	default:
		return ""
	}
}
