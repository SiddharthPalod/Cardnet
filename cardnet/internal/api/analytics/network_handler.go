package analytics

import (
	"context"
)

func (s *Server) GetNetworkStats(ctx context.Context, req *GetNetworkStatsRequest) (*GetNetworkStatsResponse, error) {
	stats, err := s.repo.GetNetworkMetrics(ctx)
	if err != nil {
		return nil, err
	}
	return &GetNetworkStatsResponse{
		TotalTx:      int64(stats.TotalTx),
		ApprovedTx:   int64(stats.ApprovedTx),
		ApprovalRate: stats.Rate(),
	}, nil
}
