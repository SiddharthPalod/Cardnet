package simulator

import (
	"cardnet/internal/issuer/model"
	"context"
	"time"

	issuerpb "cardnet/internal/api/issuer"
)

func Execute(
	req *issuerpb.IssuerAuthRequest,
	issuer model.IssuerProfile,
) (*issuerpb.IssuerAuthResponse, error) {
	rng := rngFromSeed(req.AuthId)
	time.Sleep(time.Duration(issuer.LatencyMs) * time.Millisecond)

	if rng.Float64() < issuer.TimeoutRate {
		time.Sleep(2 * time.Second)
		return nil, context.DeadlineExceeded
	}

	if rng.Float64() < issuer.DeclineRate {
		return &issuerpb.IssuerAuthResponse{
			Approved:    false,
			DeclineCode: "INSUFFICIENT_FUNDS",
		}, nil
	}

	// Chaos: forced latency
	if issuer.Chaos.ExtraLatencyMs > 0 {
		time.Sleep(time.Duration(issuer.Chaos.ExtraLatencyMs) * time.Millisecond)
	}
	// Chaos: forced timeout
	if issuer.Chaos.ForceTimeout {
		time.Sleep(2 * time.Second)
		return nil, context.DeadlineExceeded
	}
	// Chaos: forced decline
	if issuer.Chaos.ForceDecline {
		return &issuerpb.IssuerAuthResponse{
			Approved:    false,
			DeclineCode: "CHAOS_DECLINE",
		}, nil
	}

	return &issuerpb.IssuerAuthResponse{
		Approved: true,
	}, nil
}
