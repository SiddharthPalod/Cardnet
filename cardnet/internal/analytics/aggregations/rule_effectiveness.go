package aggregations

type RuleStats struct {
	RuleID         string
	TriggeredCount int
	DeclineCount   int
}

func (r *RuleStats) DeclineRate() float64 {
	if r.TriggeredCount == 0 {
		return 0
	}
	return float64(r.DeclineCount) / float64(r.TriggeredCount)
}
