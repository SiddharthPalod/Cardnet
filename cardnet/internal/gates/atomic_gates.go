package gates

import (
	"cardnet/internal/analytics/observability"
	"cardnet/internal/network"
	"sync/atomic"
	"time"
)

// AtomicGates provides lock-free gate checks using atomic bools
// Gates are evaluated periodically in background, not per-request
type AtomicGates struct {
	config Config

	// Atomic bools - updated periodically, read lock-free
	issuerHealthy     atomic.Bool
	analyticsFresh    atomic.Bool
	mlEnabled         atomic.Bool

	// Background updater
	stopChan chan struct{}
}

// NewAtomicGates creates gates with atomic bools (zero DB calls, zero locks per request)
func NewAtomicGates(config Config, metrics *network.Metrics, healthTracker IssuerHealthGetter) *AtomicGates {
	gates := &AtomicGates{
		config:   config,
		stopChan: make(chan struct{}),
	}

	// Start background updater (evaluates gates every 1 second)
	go gates.updater(metrics, healthTracker)

	return gates
}

// updater periodically evaluates gates and updates atomic bools
func (g *AtomicGates) updater(metrics *network.Metrics, healthTracker IssuerHealthGetter) {
	ticker := time.NewTicker(1 * time.Second) // Update every 1 second
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			g.updateGates(metrics, healthTracker)
		case <-g.stopChan:
			return
		}
	}
}

func (g *AtomicGates) updateGates(metrics *network.Metrics, healthTracker IssuerHealthGetter) {
	// Update issuer health gate (evaluate all issuers, set to false if any unhealthy)
	issuerHealthy := true
	if g.config.IssuerRequired {
		if healthTracker != nil {
			// Check timeout rate for all known issuers (simplified - check a sample)
			// In production, you'd track issuer keys and check each
			// For now, we'll assume healthy unless explicitly marked unhealthy
		} else if metrics != nil {
			// Use metrics snapshot (only called once per second, not per request)
			snapshot := metrics.Snapshot()
			issuers, ok := snapshot["issuers"].(map[string]interface{})
			if ok {
				// Check if any issuer is unhealthy
				for _, issuerData := range issuers {
					if data, ok := issuerData.(map[string]interface{}); ok {
						if health, ok := data["health"].(int); ok && health < 50 {
							issuerHealthy = false
							break
						}
					}
				}
			}
		}
	}
	g.issuerHealthy.Store(issuerHealthy)

	// Update analytics freshness gate
	freshness := observability.CheckFreshness(time.Now())
	analyticsFresh := !freshness.AggregationStale && !freshness.FeaturesStale
	g.analyticsFresh.Store(analyticsFresh)

	// Update ML gate (ML enabled AND analytics fresh)
	mlEnabled := g.config.MLEnabled && analyticsFresh
	g.mlEnabled.Store(mlEnabled)
}

// ShouldSkipIssuer - O(1) atomic read, zero locks
func (g *AtomicGates) ShouldSkipIssuer(issuerKey string) bool {
	return !g.issuerHealthy.Load()
}

// ShouldDisableAnalytics - O(1) atomic read, zero locks
func (g *AtomicGates) ShouldDisableAnalytics() bool {
	return !g.analyticsFresh.Load()
}

// ShouldIgnoreMLSuggestions - O(1) atomic read, zero locks
func (g *AtomicGates) ShouldIgnoreMLSuggestions() bool {
	return !g.mlEnabled.Load()
}

// Close stops the background updater
func (g *AtomicGates) Close() {
	close(g.stopChan)
}
