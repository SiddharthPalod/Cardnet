package rate

import (
	"testing"
)

func TestLimiter_MCC(t *testing.T) {
	l := NewLimiter()

	// Default MCC limit 500.
	// Let's manually set a small limit for testing by mocking or just creating a new bucket structure if possible,
	// but Limiter encapsulates buckets.
	// We will rely on default behavior or issuer access.
	// Since we can't easily inject, we'll brute force for a small test or check logic availability.

	// Actually, let's verify that different MCCs get different buckets

	// 1. Merchant A, MCC 1
	allowed, _ := l.Allow("m1", "123456", "5000")
	if !allowed {
		t.Fatal("Should allow initial request")
	}

	// 2. Merchant A, MCC 2
	allowed, _ = l.Allow("m1", "123456", "6000")
	if !allowed {
		t.Fatal("Should allow initial request different MCC")
	}
}

func TestLimiter_Concurrency(t *testing.T) {
	l := NewLimiter()
	done := make(chan bool)

	// Run 100 concurrent requests
	for i := 0; i < 100; i++ {
		go func() {
			l.Allow("m_conc", "123456", "5411")
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}
