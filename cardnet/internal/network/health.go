package network

import "sync"

type HealthScore struct {
	mu sync.Mutex

	successes int
	failures  int
}

func NewHealthScore() *HealthScore {
	return &HealthScore{}
}

func (h *HealthScore) RecordSuccess() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.successes++
}

func (h *HealthScore) RecordFailure() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.failures++
}

func (h *HealthScore) Score() int {
	h.mu.Lock()
	defer h.mu.Unlock()

	total := h.successes + h.failures
	if total == 0 {
		return 100
	}

	successRate := float64(h.successes) / float64(total)
	return int(successRate * 100)
}
