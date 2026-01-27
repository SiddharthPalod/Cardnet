package observability

import (
	"sync"
	"time"
)

type Metrics struct {
	mu sync.RWMutex

	LastAggregationRun time.Time
	LastFeatureBuild   time.Time

	LedgerMessages int64
	IssuerMessages int64
	RiskMessages   int64
}

var metrics = &Metrics{}

func RecordAggregationRun(t time.Time) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.LastAggregationRun = t
}

func RecordFeatureBuild(t time.Time) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.LastFeatureBuild = t
}

func IncLedgerMessages() {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.LedgerMessages++
}

func IncIssuerMessages() {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.IssuerMessages++
}

func IncRiskMessages() {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	metrics.RiskMessages++
}

func Snapshot() Metrics {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()
	return *metrics
}
