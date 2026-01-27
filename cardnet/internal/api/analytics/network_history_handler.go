package analytics

import (
	"context"
)

// GetNetworkHistory serves network-wide authorization history for analytics UI.
func (s *Server) GetNetworkHistory(ctx context.Context, req *GetNetworkHistoryRequest) (*GetNetworkHistoryResponse, error) {
	items, err := s.repo.GetNetworkHistory(ctx, req.GetMerchantId(), req.GetStatus(), req.GetLimit())
	if err != nil {
		return nil, err
	}

	out := make([]*NetworkHistoryItem, 0, len(items))
	for _, it := range items {
		out = append(out, &NetworkHistoryItem{
			AuthId:     it.AuthID,
			MerchantId: it.MerchantID,
			Amount:     it.Amount,
			Currency:   it.Currency,
			Status:     it.Status,
			Reason:     it.Reason,
			CreatedAt:  it.CreatedAt.Unix(),
		})
	}

	return &GetNetworkHistoryResponse{Items: out}, nil
}
