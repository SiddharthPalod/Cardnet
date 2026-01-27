package main

import (
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	log.Printf("Connecting to Kafka at %s...", broker)

	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	topics := []string{
		"auth.events",
		"risk.decisions",
		"issuer.responses",
		"settlement.ready",
	}

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("Failed to get controller: %v", err)
	}
	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		log.Fatalf("Failed to connect to controller: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to connect to controller: %v", err)
	}
	defer controllerConn.Close()

	for _, topic := range topics {
		topicConfig := kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		}

		err := controllerConn.CreateTopics(topicConfig)
		if err != nil {
			// Topic might already exist, check if it's a different error
			if err.Error() != "topic already exists" {
				log.Printf("Warning: Failed to create topic %s: %v (may already exist)", topic, err)
			} else {
				log.Printf("Topic %s already exists, skipping", topic)
			}
		} else {
			log.Printf("✅ Created topic: %s", topic)
		}
	}

	log.Println("Topic creation complete!")
}
