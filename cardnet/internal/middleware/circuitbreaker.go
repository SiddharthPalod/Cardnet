package middleware

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failureCount     int
	failureThreshold int
	resetTimeout     time.Duration
	lastFailure      time.Time
}

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: threshold,
		resetTimeout:     resetTimeout,
	}
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
	cb.checkState()

	if cb.isOpen() {
		return ErrCircuitOpen
	}

	err := operation()

	cb.reportResult(err)

	return err
}

func (cb *CircuitBreaker) checkState() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateOpen && time.Since(cb.lastFailure) > cb.resetTimeout {
		cb.state = StateHalfOpen
	}
}

func (cb *CircuitBreaker) isOpen() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state == StateOpen
}

func (cb *CircuitBreaker) reportResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()
		if cb.state == StateHalfOpen || cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
		}
	} else {
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
			cb.failureCount = 0
		} else if cb.state == StateClosed {
			cb.failureCount = 0
		}
	}
}
