package config

import "os"

type Config struct {
	TimeZone     string
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

func Load() Config {
	brokers := []string{"Kafka:9092"}
	if val := os.Getenv("KAFKA_ADDR"); val != "" {
		brokers = []string{val}
	}

	return Config{
		TimeZone:     "UTC",
		KafkaBrokers: brokers,
		KafkaTopic:   "auth.events",
		KafkaGroupID: "settlement-service",
	}
}
