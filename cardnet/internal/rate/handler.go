package rate

import (
	"context"

	pb "cardnet/internal/api/grpc"
)

type Handler struct {
	pb.UnimplementedRateLimiterServer
	limiter *Limiter
}

func NewHandler(l *Limiter) *Handler {
	return &Handler{limiter: l}
}

func (h *Handler) CheckLimit(
	ctx context.Context,
	req *pb.RateLimitRequest,
) (*pb.RateLimitResponse, error) {

	allowed, reason := h.limiter.Allow(req.MerchantId, req.Bin, req.Mcc)

	return &pb.RateLimitResponse{
		Allowed: allowed,
		Reason:  reason,
	}, nil
}
