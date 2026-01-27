package main

import (
	"cardnet/internal/analytics/aggregations"
	"cardnet/internal/analytics/config"
	"cardnet/internal/analytics/consumer"
	"cardnet/internal/analytics/repository/postgres"
	"cardnet/internal/analytics/worker"
	"cardnet/internal/api/analytics"
	"cardnet/internal/gates"
	"cardnet/internal/network"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	// Start Consumers
	ledgerConsumer := consumer.NewLedgerConsumer(cfg.Kafka, db)
	issuerConsumer := consumer.NewIssuerConsumer(cfg.Kafka, db)
	riskConsumer := consumer.NewRiskConsumer(cfg.Kafka, db)

	go ledgerConsumer.Start()
	go issuerConsumer.Start()
	go riskConsumer.Start()

	// Load gates configuration
	gatesConfig, err := gates.LoadConfig("")
	if err != nil {
		log.Printf("Failed to load gates config, using defaults: %v", err)
		gatesConfig = gates.DefaultConfig()
	}

	// Create gates instance (using network metrics for issuer health, but ML gate uses analytics freshness)
	networkMetrics := network.NewMetrics()
	gatesInstance := gates.NewGates(gatesConfig, networkMetrics)

	// Start ML Worker with gates
	mlWorker := worker.NewMLWorker(db, gatesInstance)
	mlWorker.Start()

	// Start Batch Aggregation Scheduler
	batchService := aggregations.NewBatchAggregationService(db)
	batchScheduler := aggregations.NewScheduler(batchService)
	batchScheduler.Start()

	// Start gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	analyticsServer := analytics.NewServer(db)
	analytics.RegisterAnalyticsServiceServer(grpcServer, analyticsServer)

	log.Printf("analytics-service started on port %d", cfg.Port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("analytics-service shutting down")
	grpcServer.GracefulStop()
}
