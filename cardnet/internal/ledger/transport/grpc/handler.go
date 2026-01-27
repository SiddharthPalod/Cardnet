package grpc

import (
	ledgerpb "cardnet/internal/api/ledger"
	"cardnet/internal/ledger/domain"
	"cardnet/internal/ledger/service"
	"context"
)

type Handler struct {
	ledgerpb.UnimplementedEventLedgerServiceServer
	svc *service.LedgerService
}

func NewHandler(svc *service.LedgerService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RecordEvent(
	ctx context.Context,
	req *ledgerpb.AuthEvent) (*ledgerpb.RecordEventResponse, error) {
	event := domain.FromProto(req)

	h.svc.Record(ctx, event)
	return &ledgerpb.RecordEventResponse{
		Accepted: true,
	}, nil
}
