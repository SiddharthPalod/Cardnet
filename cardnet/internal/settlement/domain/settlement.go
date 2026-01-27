package domain

import "time"

type SettlementAuthorization struct {
	AuthorizationID string
	MerchantID      string
	CardHash        string
	Amount          int64
	Currency        string
	Approved        bool
	EventTime       time.Time
	CreatedAt       time.Time
}

type AuthorizationPayload struct {
	Status   string `json:"status"`
	Reason   string `json:"reason"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}
