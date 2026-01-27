package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Kafka    KafkaConfig
	Postgres PostgresConfig
	Port     int
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

type PostgresConfig struct {
	DSN string
}

func Load() *Config {
	viper.SetConfigName("analytics")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("port", 60057)
	viper.SetDefault("kafka.brokers", []string{"kafka:9092"})
	viper.SetDefault("kafka.group_id", "analytics-group")
	viper.SetDefault("postgres.dsn", "postgres://cardnet:cardnet@postgres:5432/cardnet?sslmode=disable")

	// Try to read config file, but don't fail if it doesn't exist (env vars will be used)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: config file not found, using defaults and environment variables: %v", err)
	}

	// Override with environment variables if set
	if port := viper.GetInt("PORT"); port != 0 {
		viper.Set("port", port)
	}
	if dsn := viper.GetString("POSTGRES_DSN"); dsn != "" {
		viper.Set("postgres.dsn", dsn)
	}
	if brokers := viper.GetString("KAFKA_BROKERS"); brokers != "" {
		viper.Set("kafka.brokers", []string{brokers})
	}
	if groupID := viper.GetString("KAFKA_GROUP_ID"); groupID != "" {
		viper.Set("kafka.group_id", groupID)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("config unmarshal failed: %v", err)
	}

	return &cfg
}
