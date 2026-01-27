package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cardnet/internal/db"
	"cardnet/internal/settlement/config"
	"cardnet/internal/settlement/consumer"
	"cardnet/internal/settlement/reporting"
	"cardnet/internal/settlement/repository/postgres"
	"cardnet/internal/settlement/scheduler"
	"cardnet/internal/settlement/service"

	"github.com/segmentio/kafka-go"
)

func main() {
	cfg := config.Load()

	pg, err := db.NewPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	// Initialize Schema
	if err := initSchema(pg); err != nil {
		log.Printf("Warning: Failed to initialize schema: %v", err)
		// Proceeding, as it might be transient or already done
	}

	// repositories
	batchRepo := postgres.NewBatchRepository(pg)
	reconRepo := postgres.NewReconciliationRepository(pg)
	settlementRepo := postgres.NewSettlementRepository(pg)

	// reporting
	reportGen := reporting.NewGenerator(pg)

	// services
	batchProcessor := service.NewBatchProcessor(batchRepo)
	reconService := service.NewReconciliationService(reconRepo, reportGen)
	settlementService := service.NewSettlementService(settlementRepo)

	// schedulers
	t0 := scheduler.NewT0Scheduler(batchProcessor)
	t1 := scheduler.NewT1Scheduler(batchProcessor, reconService)
	runner := scheduler.NewRunner(t0, t1)

	// consumer
	authConsumer := consumer.NewAuthorizationConsumer(settlementService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Scheduler
	go runner.Start(ctx)

	// Start Kafka Consumer
	go func() {
		log.Printf("Connecting to Kafka at %v...", cfg.KafkaBrokers)
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     cfg.KafkaBrokers,
			Topic:       cfg.KafkaTopic,
			GroupID:     cfg.KafkaGroupID,
			StartOffset: kafka.FirstOffset, // Important: Read from beginning if no offset
		})
		defer r.Close()

		log.Println("Starting Kafka consumer loop...")
		for {
			select {
			case <-ctx.Done():
				log.Println("Context done, stopping consumer")
				return
			default:
				// log.Println("Reading message...") // Too noisy for loop
				m, err := r.ReadMessage(ctx)
				if err != nil {
					log.Printf("could not read message: %v", err)
					continue
				}
				// log.Printf("Received message key=%s val=%s", string(m.Key), string(m.Value))
				if err := authConsumer.HandleMessage(ctx, m.Value); err != nil {
					log.Printf("failed to handle message: %v", err)
				}
			}
		}
	}()

	log.Println("settlement-service started")

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("settlement-service shutting down")
}

func initSchema(db *sql.DB) error {
	// Read schema file
	content, err := os.ReadFile("internal/settlement/repository/postgres/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(content))
	return err
}
