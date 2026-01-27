package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	riskpb "cardnet/internal/api/risk"
	"cardnet/internal/risk/config"
	"cardnet/internal/risk/engine"
	"cardnet/internal/risk/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type riskServer struct {
	riskpb.UnimplementedRiskEngineServer
	engine *engine.RiskEngine
}

func (s *riskServer) Evaluate(
	ctx context.Context,
	req *riskpb.RiskRequest,
) (*riskpb.RiskResponse, error) {

	// -----------------------------
	// Map Transaction
	// -----------------------------
	tx := models.Transaction{
		Amount:    req.Transaction.Amount,
		Country:   "IN", // TODO: geo enrichment later
		Timestamp: time.Now().Unix(),
	}

	// -----------------------------
	// Map RiskContext (IMPORTANT)
	// -----------------------------
	var rctx models.RiskContext
	if req.Context != nil {
		rctx = models.RiskContext{
			TxCount1Min:       int(req.Context.TxCount_1M),
			TxCount5Min:       int(req.Context.TxCount_5M),
			LastCountry:       req.Context.LastCountry,
			LastTxTime:        req.Context.LastTxTs,
			UserAvgAmount:     req.Context.AvgAmount,
			UserP95Amount:     req.Context.P95Amount,
			MerchantRiskScore: int(req.Context.MerchantRiskScore),
			BINRiskScore:      int(req.Context.BinRiskScore),
		}
	}

	// -----------------------------
	// Evaluate Risk
	// -----------------------------
	result := s.engine.Evaluate(tx, rctx)

	return &riskpb.RiskResponse{
		Decision:  mapDecisionToProto(result.Decision),
		RiskScore: int32(result.RiskScore),
		Reasons:   result.Reasons,
	}, nil
}

func mapDecisionToProto(decision string) riskpb.Decision {
	switch decision {
	case "APPROVE":
		return riskpb.Decision_APPROVE
	case "SOFT_DECLINE":
		return riskpb.Decision_SOFT_DECLINE
	case "HARD_DECLINE":
		return riskpb.Decision_HARD_DECLINE
	default:
		return riskpb.Decision_DECISION_UNSPECIFIED
	}
}

func main() {
	config.Load(config.RuntimeConfig{
		VelocityHighTx1Min: 5,
		VelocityHighScore:  30,
	})

	re := engine.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "60052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	riskpb.RegisterRiskEngineServer(
		grpcServer,
		&riskServer{engine: re},
	)
	reflection.Register(grpcServer)

	log.Println("Risk Engine gRPC running on :60052")
	log.Fatal(grpcServer.Serve(lis))
}
