package analytics

import (
	"context"
)

func (s *Server) GetBinStats(ctx context.Context, req *GetBinStatsRequest) (*GetBinStatsResponse, error) {
	stats, err := s.repo.GetBinStats(ctx, req.GetBin())
	if err != nil {
		return nil, err
	}
	return &GetBinStatsResponse{
		Bin:          stats.BIN,
		TxCount:      int64(stats.TxCount),
		ApprovedTx:   int64(stats.ApprovedTx),
		TotalAmount:  stats.TotalAmount,
		ApprovalRate: stats.ApprovalRate(),
	}, nil
}
