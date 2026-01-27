package analytics

import (
	"context"
)

func (s *Server) GetMerchantStats(ctx context.Context, req *GetMerchantStatsRequest) (*GetMerchantStatsResponse, error) {
	stats, err := s.repo.GetMerchantStats(ctx, req.GetMerchantId())
	if err != nil {
		return nil, err
	}
	return &GetMerchantStatsResponse{
		MerchantId:   stats.MerchantID,
		TxCount:      int64(stats.TxCount),
		ApprovedTx:   int64(stats.ApprovedTx),
		TotalAmount:  stats.TotalAmount,
		ApprovalRate: stats.ApprovalRate(),
	}, nil
}
