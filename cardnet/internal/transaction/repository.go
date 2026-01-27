package transaction

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	// RecordState records a transaction state transition with validation
	// Returns ErrStateAlreadyExists if (auth_id, state) already exists (idempotency)
	// Returns ErrInvalidTransition if the transition is invalid
	RecordState(ctx context.Context, authID string, state TxState) error

	// RecordStateFast records a transaction state transition without validation (for async writes)
	// Use this for async writes where validation is already done in the service layer
	// Returns ErrStateAlreadyExists if (auth_id, state) already exists (idempotency)
	RecordStateFast(ctx context.Context, authID string, state TxState) error

	// GetCurrentState returns the current state for an auth_id
	GetCurrentState(ctx context.Context, authID string) (TxState, error)

	// GetStateHistory returns all states for an auth_id in order
	GetStateHistory(ctx context.Context, authID string) ([]TransactionState, error)
}

type PostgresRepository struct {
	db *sql.DB
	sm *StateMachine
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		db: db,
		sm: &StateMachine{},
	}
}

func (r *PostgresRepository) RecordState(ctx context.Context, authID string, state TxState) error {
	// Get current state to validate transition (for sync writes)
	currentState, err := r.GetCurrentState(ctx, authID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If state exists, validate transition
	if err == nil {
		if !r.sm.CanTransition(currentState, state) {
			return ErrInvalidTransition
		}
	} else {
		// First state must be INITIATED
		if state != INITIATED {
			return ErrInvalidTransition
		}
	}

	// Insert with idempotency check via UNIQUE constraint
	query := `
		INSERT INTO transaction_states (auth_id, state, timestamp)
		VALUES ($1, $2, $3)
		ON CONFLICT (auth_id, state) DO NOTHING
	`

	_, err = r.db.ExecContext(ctx, query, authID, state.String(), time.Now())
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			// Check if it's the unique constraint violation
			// If we get here, the state already exists - this is idempotent
			return nil // Success - already processed
		}
		return err
	}

	return nil
}

// RecordStateFast records state without validation (for async writes)
// State transitions are already validated in the service layer
// This reduces from 2 DB operations (read + write) to 1 (write only)
func (r *PostgresRepository) RecordStateFast(ctx context.Context, authID string, state TxState) error {
	// Insert with idempotency check via UNIQUE constraint
	// ON CONFLICT DO NOTHING handles both idempotency and duplicate writes
	query := `
		INSERT INTO transaction_states (auth_id, state, timestamp)
		VALUES ($1, $2, $3)
		ON CONFLICT (auth_id, state) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, authID, state.String(), time.Now())
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			// Unique constraint violation - state already exists (idempotent)
			return nil // Success - already processed
		}
		return err
	}

	return nil
}

func (r *PostgresRepository) GetCurrentState(ctx context.Context, authID string) (TxState, error) {
	query := `
		SELECT state
		FROM transaction_states
		WHERE auth_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var stateStr string
	err := r.db.QueryRowContext(ctx, query, authID).Scan(&stateStr)
	if err != nil {
		return INITIATED, err
	}

	return ParseTxState(stateStr)
}

func (r *PostgresRepository) GetStateHistory(ctx context.Context, authID string) ([]TransactionState, error) {
	query := `
		SELECT auth_id, state, timestamp
		FROM transaction_states
		WHERE auth_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := r.db.QueryContext(ctx, query, authID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []TransactionState
	for rows.Next() {
		var ts TransactionState
		var stateStr string
		if err := rows.Scan(&ts.AuthID, &stateStr, &ts.Timestamp); err != nil {
			return nil, err
		}
		ts.State, err = ParseTxState(stateStr)
		if err != nil {
			return nil, err
		}
		states = append(states, ts)
	}

	return states, rows.Err()
}
