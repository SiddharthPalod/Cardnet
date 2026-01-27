package gates

import (
	"cardnet/internal/analytics/observability"
	"cardnet/internal/network"
	"time"
)

// Gates holds all behavior gates
type Gates struct {
	config Config
	metrics *network.Metrics
	// IssuerHealthGetter is an interface for getting issuer health
	IssuerHealthGetter IssuerHealthGetter
}

// IssuerHealthGetter interface for getting issuer health metrics
type IssuerHealthGetter interface {
	GetTimeoutRate(issuerKey string) float64
	GetHealthScore(issuerKey string) int
}

// NewGates creates a new Gates instance
func NewGates(config Config, metrics *network.Metrics) *Gates {
	return &Gates{
		config:  config,
		metrics: metrics,
	}
}

// NewGatesWithHealthTracker creates a new Gates instance with a custom health tracker
func NewGatesWithHealthTracker(config Config, healthTracker IssuerHealthGetter) *Gates {
	return &Gates{
		config:             config,
		IssuerHealthGetter: healthTracker,
	}
}

// IssuerHealthGate checks if issuer should be skipped due to health issues
// Returns true if issuer should be skipped (gate closed)
func (g *Gates) IssuerHealthGate(issuerKey string) bool {
	if !g.config.IssuerRequired {
		return false // Gate open - issuer not required
	}

	// Use custom health tracker if available
	if g.IssuerHealthGetter != nil {
		timeoutRate := g.IssuerHealthGetter.GetTimeoutRate(issuerKey)
		// If timeout rate > 20%, skip issuer (gate closed)
		if timeoutRate > 0.2 {
			return true // Gate closed - skip issuer
		}
		return false
	}

	// Fallback to network metrics if available
	if g.metrics != nil {
		snapshot := g.metrics.Snapshot()
		issuers, ok := snapshot["issuers"].(map[string]interface{})
		if !ok {
			return false // Gate open - can't determine health
		}

		issuerData, ok := issuers[issuerKey].(map[string]interface{})
		if !ok {
			return false // Gate open - issuer not found
		}

		health, ok := issuerData["health"].(int)
		if !ok {
			return false
		}

		// If health is below 50%, skip issuer (gate closed)
		if health < 50 {
			return true // Gate closed - skip issuer
		}
	}

	return false // Gate open - use issuer
}

// AnalyticsFreshnessGate checks if analytics data is stale
// Returns true if analytics should be disabled (gate closed)
func (g *Gates) AnalyticsFreshnessGate() bool {
	if !g.config.AnalyticsEnabled {
		return true // Gate closed - analytics disabled
	}

	freshness := observability.CheckFreshness(time.Now())
	
	// If either aggregation or features are stale, close the gate
	if freshness.AggregationStale || freshness.FeaturesStale {
		return true // Gate closed - data stale
	}

	return false // Gate open - data fresh
}

// MLSafetyGate checks if ML suggestions should be ignored
// Returns true if ML should be ignored (gate closed)
func (g *Gates) MLSafetyGate() bool {
	if !g.config.MLEnabled {
		return true // Gate closed - ML disabled
	}

	// Check if analytics data is stale (ML depends on fresh data)
	if g.AnalyticsFreshnessGate() {
		return true // Gate closed - data stale
	}

	return false // Gate open - ML enabled and data fresh
}

// ShouldSkipIssuer returns true if issuer should be skipped
func (g *Gates) ShouldSkipIssuer(issuerKey string) bool {
	return g.IssuerHealthGate(issuerKey)
}

// ShouldDisableAnalytics returns true if analytics should be disabled
func (g *Gates) ShouldDisableAnalytics() bool {
	return g.AnalyticsFreshnessGate()
}

// ShouldIgnoreMLSuggestions returns true if ML suggestions should be ignored
func (g *Gates) ShouldIgnoreMLSuggestions() bool {
	return g.MLSafetyGate()
}
