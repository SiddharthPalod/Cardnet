package integration

import (
	"context"
	"net"
	"testing"
	"time"

	ratepb "cardnet/internal/api/grpc"
	issuerpb "cardnet/internal/api/issuer"
	ledgerpb "cardnet/internal/api/ledger"
	riskpb "cardnet/internal/api/risk"
	authsvc "cardnet/internal/auth"
	riskengine "cardnet/internal/risk/engine"
	riskmodels "cardnet/internal/risk/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type riskServer struct {
	riskpb.UnimplementedRiskEngineServer
	engine *riskengine.RiskEngine
}

func (s *riskServer) Evaluate(
	ctx context.Context,
	req *riskpb.RiskRequest,
) (*riskpb.RiskResponse, error) {

	result := s.engine.Evaluate(
		riskmodels.Transaction{
			Amount:    req.Transaction.Amount,
			Country:   "IN",
			Timestamp: time.Now().Unix(),
		},
		riskmodels.RiskContext{
			TxCount1Min:   10,
			UserAvgAmount: 1000,
		},
	)

	return &riskpb.RiskResponse{
		Decision:  riskpb.Decision_SOFT_DECLINE,
		RiskScore: int32(result.RiskScore),
		Reasons:   result.Reasons,
	}, nil
}

type rateServer struct {
	ratepb.UnimplementedRateLimiterServer
}

func (s *rateServer) CheckLimit(
	ctx context.Context,
	req *ratepb.RateLimitRequest,
) (*ratepb.RateLimitResponse, error) {
	return &ratepb.RateLimitResponse{
		Allowed: true,
	}, nil
}

type issuerServer struct {
	issuerpb.UnimplementedIssuerServiceServer
}

func (s *issuerServer) Authorize(
	ctx context.Context,
	req *issuerpb.IssuerAuthRequest,
) (*issuerpb.IssuerAuthResponse, error) {
	return &issuerpb.IssuerAuthResponse{
		Approved: true,
	}, nil
}

type ledgerServer struct {
	ledgerpb.UnimplementedEventLedgerServiceServer
}

func (s *ledgerServer) RecordEvent(
	ctx context.Context,
	req *ledgerpb.AuthEvent,
) (*ledgerpb.RecordEventResponse, error) {
	return &ledgerpb.RecordEventResponse{
		Accepted: true,
	}, nil
}

// Integration Test
func TestAuthRiskIntegration(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	riskpb.RegisterRiskEngineServer(
		grpcServer,
		&riskServer{engine: riskengine.New()},
	)
	ratepb.RegisterRateLimiterServer(
		grpcServer,
		&rateServer{},
	)
	issuerpb.RegisterIssuerServiceServer(
		grpcServer,
		&issuerServer{},
	)
	ledgerpb.RegisterEventLedgerServiceServer(
		grpcServer,
		&ledgerServer{},
	)

	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	riskClient := riskpb.NewRiskEngineClient(conn)
	rateClient := ratepb.NewRateLimiterClient(conn)
	issuerClient := issuerpb.NewIssuerServiceClient(conn)
	ledgerClient := ledgerpb.NewEventLedgerServiceClient(conn)

	authService := authsvc.NewService(
		authsvc.NewInMemoryRepo(),
		rateClient,
		riskClient,
		issuerClient,
		ledgerClient,
	)

	// Call Authorize
	_, status, _ := authService.Authorize(
		context.Background(),
		"req-1",
		"merchant-1",
		"4111111111111111",
		5000,
		"INR",
		"5411",
	)

	if status != "CHALLENGE_REQUIRED" {
		t.Fatalf("expected CHALLENGE_REQUIRED, got %s", status)
	}
}
