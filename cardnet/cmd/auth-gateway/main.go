package main

import (
	"cardnet/internal/issuer/config"
	"cardnet/internal/issuer/routing"
	"log"
	"net"
	"os"

	pb "cardnet/internal/api/grpc"
	issuerpb "cardnet/internal/api/issuer"
	ledgerpb "cardnet/internal/api/ledger"
	riskpb "cardnet/internal/api/risk"
	"cardnet/internal/auth"
	"cardnet/internal/db"
	"cardnet/internal/gates"
	"cardnet/internal/network"
	"cardnet/internal/transaction"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const issuerConfigPath = "internal/issuer/config/issuers.yaml"

func main() {
	log.Println("Starting CardNet Auth Gateway...")

	// Init DB
	database, err := db.NewPostgres()
	if err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	// Initialize transaction state schema
	txSchema := `
		CREATE TABLE IF NOT EXISTS transaction_states (
			auth_id TEXT NOT NULL,
			state TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY (auth_id, state),
			CONSTRAINT valid_state CHECK (state IN ('INITIATED', 'RISK_EVALUATED', 'ISSUER_REQUESTED', 'ISSUER_RESPONDED', 'FINALIZED'))
		);
		CREATE INDEX IF NOT EXISTS idx_transaction_states_auth_id ON transaction_states(auth_id);
		CREATE INDEX IF NOT EXISTS idx_transaction_states_timestamp ON transaction_states(timestamp);
	`
	if _, err := database.Exec(txSchema); err != nil {
		log.Printf("Warning: Failed to initialize transaction_states schema (may already exist): %v", err)
	}

	// CRITICAL FIX: Initialize authorizations table (was missing, causing DB errors)
	authSchema := `
		CREATE TABLE IF NOT EXISTS authorizations (
			auth_id TEXT PRIMARY KEY,
			request_id TEXT UNIQUE NOT NULL,
			merchant_id TEXT NOT NULL,
			card_token TEXT NOT NULL,
			amount BIGINT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_authorizations_request_id ON authorizations(request_id);
		CREATE INDEX IF NOT EXISTS idx_authorizations_merchant_id ON authorizations(merchant_id);
		CREATE INDEX IF NOT EXISTS idx_authorizations_created_at ON authorizations(created_at);
	`
	if _, err := database.Exec(authSchema); err != nil {
		log.Printf("Warning: Failed to initialize authorizations schema (may already exist): %v", err)
	}

	// Initialize risk_audit table (used by SaveRiskAudit)
	riskAuditSchema := `
		CREATE TABLE IF NOT EXISTS risk_audit (
			auth_id TEXT NOT NULL,
			risk_score INTEGER NOT NULL,
			reasons TEXT[],
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY (auth_id)
		);
		CREATE INDEX IF NOT EXISTS idx_risk_audit_auth_id ON risk_audit(auth_id);
	`
	if _, err := database.Exec(riskAuditSchema); err != nil {
		log.Printf("Warning: Failed to initialize risk_audit schema (may already exist): %v", err)
	}

	// Wire layers
	repo := auth.NewRepository(database)
	txStateRepo := transaction.NewRepository(database)

	// Load gates configuration
	gatesConfig, err := gates.LoadConfig("")
	if err != nil {
		log.Printf("Failed to load gates config, using defaults: %v", err)
		gatesConfig = gates.DefaultConfig()
	}

	// Load BIN routing table (required for ResolveIssuer to work in auth service)
	table, err := config.LoadIssuers(issuerConfigPath)
	if err != nil {
		log.Printf("Warning: Failed to load issuer config: %v", err)
		log.Println("BIN routing will not work until issuer config is loaded")
	} else {
		routing.Load(table)
		log.Println("Loaded BIN routing table")
	}

	// Create shared network metrics (for issuer health tracking)
	networkMetrics := network.NewMetrics()

	// 2️⃣ CONNECT TO RATE-LIMITER
	rateAddr := os.Getenv("RATE_LIMITER_ADDR")
	if rateAddr == "" {
		rateAddr = "localhost:60051"
	}
	rateConn, err := grpc.NewClient(
		rateAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("rate-limiter connection failed: %v", err)
	}

	defer func(rateConn *grpc.ClientConn) {
		err := rateConn.Close()
		if err != nil {
			log.Fatalf("failed to create rate-limiter client: %v", err)
		}
	}(rateConn)

	rateClient := pb.NewRateLimiterClient(rateConn)

	riskAddr := os.Getenv("RISK_ENGINE_ADDR")
	if riskAddr == "" {
		riskAddr = "localhost:60052"
	}
	riskConn, err := grpc.NewClient(
		riskAddr, // risk-engine
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("risk-engine connection failed: %v", err)
	}
	defer riskConn.Close()
	riskClient := riskpb.NewRiskEngineClient(riskConn)

	// 5️⃣ CONNECT TO EVENT LEDGER (IMMUTABLE AUDIT LOG)
	ledgerAddr := os.Getenv("EVENT_LEDGER_ADDR")
	if ledgerAddr == "" {
		ledgerAddr = "localhost:60055"
	}
	ledgerConn, err := grpc.NewClient(
		ledgerAddr, // event-ledger
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("event-ledger connection failed: %v", err)
	}
	defer ledgerConn.Close()
	ledgerClient := ledgerpb.NewEventLedgerServiceClient(ledgerConn)

	// 4️⃣ CONNECT TO CARD NETWORK (ROUTING LAYER)
	issuerAddr := os.Getenv("CARD_NETWORK_ADDR")
	if issuerAddr == "" {
		issuerAddr = "localhost:60053"
	}
	issuerConn, err := grpc.NewClient(
		issuerAddr, // card-network
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("card-network connection failed: %v", err)
	}
	defer issuerConn.Close()
	issuerClient := issuerpb.NewIssuerServiceClient(issuerConn)

	service := auth.NewService(repo, txStateRepo, gatesConfig, networkMetrics, rateClient, riskClient, issuerClient, ledgerClient)
	handler := auth.NewHandler(service)

	// Start gRPC server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthGatewayServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Println("Listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
