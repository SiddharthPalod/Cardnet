package analytics

import (
	"context"
)

func (s *Server) GetRuleEffectiveness(ctx context.Context, req *GetRuleEffectivenessRequest) (*GetRuleEffectivenessResponse, error) {
	stats, err := s.repo.GetRuleStats(ctx, req.GetRuleId())
	if err != nil {
		return nil, err
	}
	return &GetRuleEffectivenessResponse{
		RuleId:         stats.RuleID,
		TriggeredCount: int64(stats.TriggeredCount),
		DeclineCount:   int64(stats.DeclineCount),
		DeclineRate:    stats.DeclineRate(),
	}, nil
}
