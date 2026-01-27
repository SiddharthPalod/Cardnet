package auth

import (
	"cardnet/internal/gates"
	"sync"
	"time"
)

// Ensure IssuerHealthTracker implements gates.IssuerHealthGetter
var _ gates.IssuerHealthGetter = (*IssuerHealthTracker)(nil)

// IssuerHealthTracker tracks issuer health metrics locally
type IssuerHealthTracker struct {
	mu sync.RWMutex
	// Map of issuer key -> health metrics
	health map[string]*IssuerHealth
}

type IssuerHealth struct {
	Successes   int64
	Failures    int64
	Timeouts    int64
	LastUpdated time.Time
}

func NewIssuerHealthTracker() *IssuerHealthTracker {
	return &IssuerHealthTracker{
		health: make(map[string]*IssuerHealth),
	}
}

func (t *IssuerHealthTracker) RecordSuccess(issuerKey string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	h := t.getOrCreate(issuerKey)
	h.Successes++
	h.LastUpdated = time.Now()
}

func (t *IssuerHealthTracker) RecordFailure(issuerKey string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	h := t.getOrCreate(issuerKey)
	h.Failures++
	h.LastUpdated = time.Now()
}

func (t *IssuerHealthTracker) RecordTimeout(issuerKey string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	h := t.getOrCreate(issuerKey)
	h.Timeouts++
	h.LastUpdated = time.Now()
}

func (t *IssuerHealthTracker) getOrCreate(issuerKey string) *IssuerHealth {
	h, ok := t.health[issuerKey]
	if !ok {
		h = &IssuerHealth{}
		t.health[issuerKey] = h
	}
	return h
}

// GetTimeoutRate returns the timeout rate for an issuer (0.0 to 1.0)
func (t *IssuerHealthTracker) GetTimeoutRate(issuerKey string) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	h, ok := t.health[issuerKey]
	if !ok {
		return 0.0
	}

	total := h.Successes + h.Failures + h.Timeouts
	if total == 0 {
		return 0.0
	}

	return float64(h.Timeouts) / float64(total)
}

// GetHealthScore returns a health score (0-100) for an issuer
func (t *IssuerHealthTracker) GetHealthScore(issuerKey string) int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	h, ok := t.health[issuerKey]
	if !ok {
		return 100 // Assume healthy if no data
	}

	total := h.Successes + h.Failures + h.Timeouts
	if total == 0 {
		return 100
	}

	successRate := float64(h.Successes) / float64(total)
	return int(successRate * 100)
}
