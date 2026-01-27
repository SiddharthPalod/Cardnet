package aggregations

type MerchantStats struct {
	MerchantID  string
	TxCount     int
	ApprovedTx  int
	TotalAmount float64
}

func (m *MerchantStats) ApprovalRate() float64 {
	if m.TxCount == 0 {
		return 0
	}
	return float64(m.ApprovedTx) / float64(m.TxCount)
}
