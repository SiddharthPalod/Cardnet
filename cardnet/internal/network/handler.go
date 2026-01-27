package network

import (
	issuerpb "cardnet/internal/api/issuer"
	"cardnet/internal/issuer/routing"
	"cardnet/internal/middleware"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	issuerpb.UnimplementedIssuerServiceServer

	client   issuerpb.IssuerServiceClient
	breakers *BreakerRegistry
	metrics  *Metrics
	health   *HealthRegistry
}

func NewHandler() *Handler {
	return &Handler{
		client:   NewIssuerClient(),
		breakers: NewBreakerRegistry(),
		metrics:  NewMetrics(),
		health:   NewHealthRegistry(),
	}
}

func (h *Handler) Authorize(
	ctx context.Context,
	req *issuerpb.IssuerAuthRequest,
) (*issuerpb.IssuerAuthResponse, error) {
	start := time.Now()
	var resp *issuerpb.IssuerAuthResponse

	// 1️⃣ RESOLVE ISSUER (For Circuit Breaking)
	profile, ok := routing.ResolveIssuer(req.CardNumber)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "unknown BIN")
	}

	// 2️⃣ KEY BY ISSUER NAME (e.g. "HDFC", "CHASE")
	issuerKey := profile.Name
	breaker := h.breakers.Get(issuerKey)
	err := breaker.Execute(func() error {
		callCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()

		return Retry(callCtx, func(c context.Context) error {
			var callErr error
			resp, callErr = h.client.Authorize(c, req)
			return callErr
		})
	})

	latency := time.Since(start)

	health_score := h.health.Get(issuerKey)
	if !IsHealthy(health_score.Score()) {
		h.metrics.SetBreakerState(issuerKey, StateOpen)
		return nil, status.Error(
			codes.Unavailable,
			"issuer temporarily unhealthy",
		)
	}

	if err != nil {
		health_score.RecordFailure()
		h.metrics.RecordFailure(latency)
		h.metrics.RecordIssuerFailure(issuerKey)
		if err == middleware.ErrCircuitOpen {
			h.metrics.SetBreakerState(issuerKey, StateOpen)
			return nil, status.Error(codes.Unavailable, "issuer circuit open")
		}

		return nil, err
	}

	health_score.RecordSuccess()
	h.metrics.RecordSuccess(latency)
	h.metrics.RecordIssuerSuccess(issuerKey)
	h.metrics.SetBreakerState(issuerKey, StateClosed)
	return resp, nil
}
