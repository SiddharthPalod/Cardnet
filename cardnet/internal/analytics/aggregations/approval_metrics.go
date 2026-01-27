package aggregations

type ApprovalMetrics struct {
	TotalTx    int
	ApprovedTx int
}

func (a *ApprovalMetrics) Rate() float64 {
	if a.TotalTx == 0 {
		return 0.0
	}
	return float64(a.ApprovedTx) / float64(a.TotalTx)
}
