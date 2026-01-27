package aggregations

import (
	"context"
	"time"
)

// BatchRepository defines the interface for batch aggregation operations
// This interface breaks the import cycle between aggregations and postgres packages
type BatchRepository interface {
	// Read operations
	GetAllMerchantStats(ctx context.Context) ([]*MerchantStats, error)
	GetAllBinStats(ctx context.Context) ([]*BINStats, error)
	GetNetworkMetrics(ctx context.Context) (*ApprovalMetrics, error)

	// Hourly aggregation write operations
	RecordHourlyMerchantAggregation(ctx context.Context, merchantID string, hourBucket time.Time, txCount, approvedCount int, totalAmount float64) error
	RecordHourlyBINAggregation(ctx context.Context, bin string, hourBucket time.Time, txCount, approvedCount int, totalAmount float64) error
	RecordHourlyNetworkAggregation(ctx context.Context, hourBucket time.Time, totalTx, approvedTx int) error

	// Daily aggregation write operations
	RecordDailyMerchantAggregation(ctx context.Context, merchantID string, dayBucket time.Time, txCount, approvedCount int, totalAmount float64) error
	RecordDailyBINAggregation(ctx context.Context, bin string, dayBucket time.Time, txCount, approvedCount int, totalAmount float64) error
	RecordDailyNetworkAggregation(ctx context.Context, dayBucket time.Time, totalTx, approvedTx int) error
}
