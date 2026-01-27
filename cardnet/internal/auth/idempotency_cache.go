package auth

import (
	"sync"
	"time"
)

// IdempotencyCache provides O(1) in-memory idempotency checks
// Fast path: memory-only check before DB
type IdempotencyCache struct {
	mu    sync.RWMutex
	cache map[string]cachedResult
	// TTL for cache entries (prevent unbounded growth)
	ttl time.Duration
}

type cachedResult struct {
	authID string
	status string
	expiry time.Time
}

func NewIdempotencyCache(ttl time.Duration) *IdempotencyCache {
	if ttl == 0 {
		ttl = 5 * time.Minute // Default 5 min TTL
	}
	return &IdempotencyCache{
		cache: make(map[string]cachedResult),
		ttl:   ttl,
	}
}

// Get returns (authID, status, found) - O(1) memory lookup
func (c *IdempotencyCache) Get(requestID string) (string, string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result, ok := c.cache[requestID]
	if !ok {
		return "", "", false
	}

	// Check expiry
	if time.Now().After(result.expiry) {
		return "", "", false
	}

	return result.authID, result.status, true
}

// Set stores a result in cache - O(1) memory write
func (c *IdempotencyCache) Set(requestID, authID, status string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[requestID] = cachedResult{
		authID: authID,
		status: status,
		expiry: time.Now().Add(c.ttl),
	}
}

// Cleanup removes expired entries (call periodically)
func (c *IdempotencyCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, v := range c.cache {
		if now.After(v.expiry) {
			delete(c.cache, k)
		}
	}
}
