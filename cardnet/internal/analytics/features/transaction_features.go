package features

import (
	"context"
	"database/sql"
	"time"
)

type TransactionFeatures struct {
	TransactionID string
	MerchantID    string
	BIN           string
	Amount        float64
	Currency      string
	Timestamp     time.Time
}

// BuildTransactionFeatures is a placeholder for retrieving historical transaction features.
// Currently returns empty as we don't store raw transactions in postgres for analytics lookup yet.
func BuildTransactionFeatures(ctx context.Context, db *sql.DB) ([]TransactionFeatures, error) {
	return []TransactionFeatures{}, nil
}

// ExtractTransactionFeatures creates a feature set from raw transaction data.
// This is useful for real-time inference where we don't query the DB.
func ExtractTransactionFeatures(
	txID, merchantID, bin string,
	amount float64,
	currency string,
	timestamp time.Time,
) TransactionFeatures {
	return TransactionFeatures{
		TransactionID: txID,
		MerchantID:    merchantID,
		BIN:           bin,
		Amount:        amount,
		Currency:      currency,
		Timestamp:     timestamp,
	}
}
