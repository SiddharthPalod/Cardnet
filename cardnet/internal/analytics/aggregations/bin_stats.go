package aggregations

type BINStats struct {
	BIN         string
	TxCount     int
	ApprovedTx  int
	TotalAmount float64
}

func (b *BINStats) ApprovalRate() float64 {
	if b.TxCount == 0 {
		return 0
	}
	return float64(b.ApprovedTx) / float64(b.TxCount)
}
