package rate

import (
	"testing"
	"time"
)

func TestTokenBucket_RefillPrecision(t *testing.T) {
	// Rate: 10 tokens per second
	// Capacity: 10
	b := NewBucket(10, 10)

	// Consume all tokens
	for i := 0; i < 10; i++ {
		if !b.Allow() {
			t.Fatalf("Failed to consume initial token %d", i)
		}
	}

	if b.Allow() {
		t.Fatal("Bucket should be empty")
	}

	// Wait 0.1s -> should gain 1 token
	time.Sleep(110 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("Should allow 1 token after 0.1s check")
	}

	if b.Allow() {
		t.Fatal("Should not allow 2nd token yet")
	}
}

func TestTokenBucket_NoLostTime(t *testing.T) {
	// Rate: 1 token per second
	b := NewBucket(1, 1)
	
	// Drain
	b.Allow() 
	
	// Wait 0.5s (0.5 tokens accumulated)
	time.Sleep(500 * time.Millisecond)
	if b.Allow() {
		t.Fatal("Should not allow token at 0.5s")
	}
	
	// Wait another 0.6s (total > 1.0s)
	time.Sleep(600 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("Should allow token after total 1.1s")
	}
}
