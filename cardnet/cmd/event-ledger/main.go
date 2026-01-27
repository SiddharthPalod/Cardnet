package main

import (
	ledgerpb "cardnet/internal/api/ledger"
	"cardnet/internal/ledger/publisher/kafka"
	"cardnet/internal/ledger/repository/cassandra"
	"cardnet/internal/ledger/service"
	ledgergrpc "cardnet/internal/ledger/transport/grpc"
	"log"
	"net"
	"os"
	"time"

	"github.com/gocql/gocql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	if cassandraHost == "" {
		cassandraHost = "127.0.0.1"
	}

	session, err := cassandra.NewSession(cassandra.Config{
		Hosts:       []string{cassandraHost},
		Keyspace:    "cardnet_ledger",
		Consistency: gocql.LocalQuorum,
		Timeout:     2 * time.Second,
	})
	if err != nil {
		log.Fatalf("failed to connect to cassandra: %v", err)
	}
	defer session.Close()

	repo := cassandra.NewLedgerRepo(session)

	kafkaAddr := os.Getenv("KAFKA_ADDR")
	if kafkaAddr == "" {
		kafkaAddr = "localhost:9092"
	}
	producer := kafka.NewProducer([]string{kafkaAddr})

	svc := service.NewLedgerService(repo, producer)
	handler := ledgergrpc.NewHandler(svc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "60055"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on :"+port+": %v", err)
	}

	grpcServer := grpc.NewServer()
	ledgerpb.RegisterEventLedgerServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Println("event-ledger listening on :" + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
