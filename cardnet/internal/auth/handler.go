package auth

import (
	"context"

	pb "cardnet/internal/api/grpc"
)

type Handler struct {
	pb.UnimplementedAuthGatewayServer
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Authorize(
	ctx context.Context,
	req *pb.AuthRequest,
) (*pb.AuthResponse, error) {

	authID, status, err := h.service.Authorize(
		ctx,
		req.RequestId,
		req.MerchantId,
		req.CardToken,
		req.Amount,
		req.Currency,
		req.Mcc,
	)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		AuthId: authID,
		Status: status,
	}, nil
}
