package observability

import "sync"

type LagMonitor struct {
	mu sync.RWMutex

	ledgerLag int64
	issuerLag int64
	riskLag   int64
}

var lagMonitor = &LagMonitor{}

func UpdateLedgerLag(lag int64) {
	lagMonitor.mu.Lock()
	defer lagMonitor.mu.Unlock()
	lagMonitor.ledgerLag = lag
}

func UpdateIssuerLag(lag int64) {
	lagMonitor.mu.Lock()
	defer lagMonitor.mu.Unlock()
	lagMonitor.issuerLag = lag
}

func UpdateRiskLag(lag int64) {
	lagMonitor.mu.Lock()
	defer lagMonitor.mu.Unlock()
	lagMonitor.riskLag = lag
}

func GetLagSnapshot() (ledger, issuer, risk int64) {
	lagMonitor.mu.RLock()
	defer lagMonitor.mu.RUnlock()
	return lagMonitor.ledgerLag, lagMonitor.issuerLag, lagMonitor.riskLag
}
