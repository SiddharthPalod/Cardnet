package issuer

import (
	"context"

	issuerpb "cardnet/internal/api/issuer"
	"cardnet/internal/issuer/routing"
	"cardnet/internal/issuer/simulator"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	issuerpb.UnimplementedIssuerServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Authorize(
	ctx context.Context,
	req *issuerpb.IssuerAuthRequest,
) (*issuerpb.IssuerAuthResponse, error) {

	profile, ok := routing.ResolveIssuer(req.CardNumber)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "unknown BIN")
	}

	return simulator.Execute(req, profile)
}
