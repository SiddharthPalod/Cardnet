package domain

import "time"

type SettlementBatch struct {
	ID        string
	BatchDate time.Time
	Status    BatchStatus
	CreatedAt time.Time
}

type BatchStatus string

const (
	BatchOpen  BatchStatus = "OPEN"
	BatchReady BatchStatus = "READY"
)
