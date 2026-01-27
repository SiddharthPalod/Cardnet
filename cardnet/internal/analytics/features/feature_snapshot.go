package features

import (
	"cardnet/internal/analytics/observability"
	"context"
	"database/sql"
	"time"
)

type FeatureSnapshot struct {
	Merchants    []MerchantFeatures
	Bins         []BINFeatures
	Transactions []TransactionFeatures
}

func BuildFeatureSnapshot(ctx context.Context, db *sql.DB) (*FeatureSnapshot, error) {
	merchants, err := BuildMerchantFeatures(ctx, db)
	if err != nil {
		return nil, err
	}

	bins, err := BuildBINFeatures(ctx, db)
	if err != nil {
		return nil, err
	}

	transactions, err := BuildTransactionFeatures(ctx, db)
	if err != nil {
		return nil, err
	}

	observability.RecordFeatureBuild(time.Now())

	return &FeatureSnapshot{
		Merchants:    merchants,
		Bins:         bins,
		Transactions: transactions,
	}, nil
}
