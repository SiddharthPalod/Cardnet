package aggregations

import "time"

// NetworkHistoryItem is a compact representation of a finalized authorization
// used by the analytics service for network-wide history queries.
type NetworkHistoryItem struct {
	AuthID     string
	MerchantID string
	Amount     float64
	Currency   string
	Status     string
	Reason     string
	CreatedAt  time.Time
}

