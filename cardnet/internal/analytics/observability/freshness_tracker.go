package observability

import "time"

type FreshnessStatus struct {
	LastAggregationRun time.Time
	LastFeatureBuild   time.Time

	AggregationStale bool
	FeaturesStale    bool
}

const (
	MaxAggregationDelay = 2 * time.Hour
	MaxFeatureDelay     = 2 * time.Hour
)

func CheckFreshness(now time.Time) FreshnessStatus {
	snap := Snapshot()

	aggStale := now.Sub(snap.LastAggregationRun) > MaxAggregationDelay
	featStale := now.Sub(snap.LastFeatureBuild) > MaxFeatureDelay

	return FreshnessStatus{
		LastAggregationRun: snap.LastAggregationRun,
		LastFeatureBuild:   snap.LastFeatureBuild,
		AggregationStale:   aggStale,
		FeaturesStale:      featStale,
	}
}
