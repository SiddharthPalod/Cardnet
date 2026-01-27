package network

import (
	"sync"
)

type HealthRegistry struct {
	mu     sync.Mutex
	scores map[string]*HealthScore
}

func NewHealthRegistry() *HealthRegistry {
	return &HealthRegistry{
		scores: make(map[string]*HealthScore),
	}
}

func (r *HealthRegistry) Get(issuer string) *HealthScore {
	r.mu.Lock()
	defer r.mu.Unlock()

	h, ok := r.scores[issuer]
	if !ok {
		h = NewHealthScore()
		r.scores[issuer] = h
	}
	return h
}
