package rate

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity float64
	tokens   float64
	lastFill time.Time
	rate     float64
	mu       sync.Mutex
}

func NewBucket(capacity, rate int) *TokenBucket {
	// Pre-fill bucket
	return &TokenBucket{
		capacity: float64(capacity),
		tokens:   float64(capacity),
		rate:     float64(rate),
		lastFill: time.Now(),
	}
}

func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	// Calculate elapsed time in seconds (floating point)
	elapsed := now.Sub(b.lastFill).Seconds()
	
	// Add tokens based on elapsed time * rate
	added := elapsed * b.rate
	b.tokens += added
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	
	// Always best to update lastFill to now to avoid drift
	// However, to be strictly precise about "saved" time, one might only advance logic time.
	// For this level of simplicity, updating to now is standard for token buckets 
	// as long as we use float precision for tokens.
	b.lastFill = now

	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true
	}
	return false
}
