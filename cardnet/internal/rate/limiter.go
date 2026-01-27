package rate

import "sync"

type Limiter struct {
	merchantBuckets map[string]*TokenBucket
	binBuckets      map[string]*TokenBucket
	mccBuckets      map[string]*TokenBucket
	mu              sync.Mutex
}

func NewLimiter() *Limiter {
	return &Limiter{
		merchantBuckets: make(map[string]*TokenBucket),
		binBuckets:      make(map[string]*TokenBucket),
		mccBuckets:      make(map[string]*TokenBucket),
	}
}

func (l *Limiter) getBuckets(merchantId, bin, mcc string) (*TokenBucket, *TokenBucket, *TokenBucket) {
	l.mu.Lock()
	defer l.mu.Unlock()

	mb, ok := l.merchantBuckets[merchantId]
	if !ok {
		mb = NewBucket(100, 100)
		l.merchantBuckets[merchantId] = mb
	}

	bb, ok := l.binBuckets[bin]
	if !ok {
		bb = NewBucket(20, 20)
		l.binBuckets[bin] = bb
	}

	mccb, ok := l.mccBuckets[mcc]
	if !ok {
		// Default limit for MCC, e.g., 500
		mccb = NewBucket(500, 500)
		l.mccBuckets[mcc] = mccb
	}

	return mb, bb, mccb
}

func (l *Limiter) Allow(merchantId, bin, mcc string) (bool, string) {
	// Retrieve buckets under a single map lock
	mb, bb, mccb := l.getBuckets(merchantId, bin, mcc)

	// Check limits sequentially (buckets handle their own locking)
	if !mb.Allow() {
		return false, "merchant rate limit exceeded"
	}

	if !bb.Allow() {
		return false, "BIN rate limit exceeded"
	}

	if !mccb.Allow() {
		// Note: Merchant and BIN tokens were consumed. This is acceptable for failing closed.
		return false, "MCC rate limit exceeded"
	}

	return true, ""
}
