package transaction

import (
	"errors"
	"time"
)

// TxState represents the canonical state of a transaction
type TxState int

const (
	INITIATED TxState = iota
	RISK_EVALUATED
	ISSUER_REQUESTED
	ISSUER_RESPONDED
	FINALIZED
)

func (s TxState) String() string {
	switch s {
	case INITIATED:
		return "INITIATED"
	case RISK_EVALUATED:
		return "RISK_EVALUATED"
	case ISSUER_REQUESTED:
		return "ISSUER_REQUESTED"
	case ISSUER_RESPONDED:
		return "ISSUER_RESPONDED"
	case FINALIZED:
		return "FINALIZED"
	default:
		return "UNKNOWN"
	}
}

// ParseTxState converts a string to TxState
func ParseTxState(s string) (TxState, error) {
	switch s {
	case "INITIATED":
		return INITIATED, nil
	case "RISK_EVALUATED":
		return RISK_EVALUATED, nil
	case "ISSUER_REQUESTED":
		return ISSUER_REQUESTED, nil
	case "ISSUER_RESPONDED":
		return ISSUER_RESPONDED, nil
	case "FINALIZED":
		return FINALIZED, nil
	default:
		return INITIATED, errors.New("invalid state")
	}
}

// TransactionState represents a state entry in the ledger
type TransactionState struct {
	AuthID    string
	State     TxState
	Timestamp time.Time
}

// StateMachine enforces valid state transitions
type StateMachine struct{}

var (
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrStateAlreadyExists = errors.New("state already exists (idempotency)")
)

// CanTransition checks if a transition from oldState to newState is valid
func (sm *StateMachine) CanTransition(oldState TxState, newState TxState) bool {
	// States are append-only - can only move forward
	if newState <= oldState {
		return false
	}

	// No skipping states - must be sequential
	if newState-oldState > 1 {
		return false
	}

	return true
}

// ValidTransitions returns all valid next states from the current state
func (sm *StateMachine) ValidTransitions(current TxState) []TxState {
	if current >= FINALIZED {
		return []TxState{} // No transitions from FINALIZED
	}
	return []TxState{current + 1}
}
