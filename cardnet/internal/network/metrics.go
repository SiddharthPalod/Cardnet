package network

import (
	"sync"
	"time"
)

type BreakerState string

const (
	StateClosed   BreakerState = "CLOSED"
	StateOpen     BreakerState = "OPEN"
	StateHalfOpen BreakerState = "HALF_OPEN"
)

type IssuerMetrics struct {
	Failures int64
	State    BreakerState
	Health   int
}

type Metrics struct {
	mu sync.Mutex

	totalCalls   int64
	successCalls int64
	failedCalls  int64
	totalLatency time.Duration

	issuers map[string]*IssuerMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		issuers: make(map[string]*IssuerMetrics),
	}

}

func (m *Metrics) RecordSuccess(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalCalls++
	m.successCalls++
	m.totalLatency += latency
}

func (m *Metrics) RecordFailure(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalCalls++
	m.failedCalls++
	m.totalLatency += latency
}

func (m *Metrics) ensureIssuer(issuer string) *IssuerMetrics {
	im, ok := m.issuers[issuer]
	if !ok {
		im = &IssuerMetrics{
			State:  StateClosed,
			Health: 100,
		}
		m.issuers[issuer] = im
	}
	return im
}

func (m *Metrics) RecordIssuerFailure(issuer string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	im := m.ensureIssuer(issuer)
	im.Failures++
	if im.Health > 0 {
		im.Health -= 5
	}
}

func (m *Metrics) RecordIssuerSuccess(issuer string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	im := m.ensureIssuer(issuer)
	if im.Health < 100 {
		im.Health += 1
	}
}

func (m *Metrics) SetBreakerState(issuer string, state BreakerState) {
	m.mu.Lock()
	defer m.mu.Unlock()

	im := m.ensureIssuer(issuer)
	im.State = state
}

func (m *Metrics) Snapshot() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	avgLatency := time.Duration(0)
	if m.totalCalls > 0 {
		avgLatency = m.totalLatency / time.Duration(m.totalCalls)
	}

	issuers := make(map[string]IssuerMetrics)
	for k, v := range m.issuers {
		issuers[k] = *v
	}

	return map[string]interface{}{
		"total_calls":    m.totalCalls,
		"success_calls":  m.successCalls,
		"failed_calls":   m.failedCalls,
		"avg_latency_ms": avgLatency.Milliseconds(),
		"issuers":        issuers,
	}
}
