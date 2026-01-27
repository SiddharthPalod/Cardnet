package network

import (
	"cardnet/internal/middleware"
	"sync"
	"time"
)

type BreakerRegistry struct {
	mu       sync.Mutex
	breakers map[string]*middleware.CircuitBreaker
}

func NewBreakerRegistry() *BreakerRegistry {
	return &BreakerRegistry{
		breakers: make(map[string]*middleware.CircuitBreaker),
	}
}

func (r *BreakerRegistry) Get(issuer string) *middleware.CircuitBreaker {
	r.mu.Lock()
	defer r.mu.Unlock()

	cb, exists := r.breakers[issuer]
	if !exists {
		cb = middleware.NewCircuitBreaker(5, 10*time.Second)
		r.breakers[issuer] = cb
	}
	return cb
}
