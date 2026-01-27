package transaction

import (
	"context"
	"log"
	"sync"
	"time"
)

// AsyncStateWriter writes state transitions asynchronously
// Fast path: memory-only, durable path: async DB writes
type AsyncStateWriter struct {
	repo      Repository // For durable writes
	queue     chan stateUpdate
	wg        sync.WaitGroup
	shutdown  chan struct{}
	numWorkers int // Number of worker goroutines
}

type stateUpdate struct {
	authID string
	state  TxState
}

// NewAsyncStateWriter creates an async state writer with buffered channel
// Uses multiple workers to process the queue faster under high load
func NewAsyncWriter(repo Repository, bufferSize int) *AsyncStateWriter {
	if bufferSize == 0 {
		bufferSize = 1000 // Default buffer
	}

	numWorkers := 3 // Use 3 workers to process queue faster
	if bufferSize > 5000 {
		// Scale workers with buffer size for high-load scenarios
		numWorkers = 5
	}

	writer := &AsyncStateWriter{
		repo:       repo,
		queue:      make(chan stateUpdate, bufferSize),
		shutdown:   make(chan struct{}),
		numWorkers: numWorkers,
	}

	// Start multiple background workers for parallel processing
	for i := 0; i < numWorkers; i++ {
		writer.wg.Add(1)
		go writer.worker()
	}

	return writer
}

// RecordStateAsync queues a state transition for async write (non-blocking)
func (w *AsyncStateWriter) RecordStateAsync(authID string, state TxState) {
	select {
	case w.queue <- stateUpdate{authID: authID, state: state}:
		// Queued successfully
	default:
		// Queue full - log but don't block
		log.Printf("WARN: state writer queue full, dropping state update: %s -> %s", authID, state.String())
	}
}

// RecordStateSync writes state synchronously (use only for FINAL state)
func (w *AsyncStateWriter) RecordStateSync(ctx context.Context, authID string, state TxState) error {
	return w.repo.RecordState(ctx, authID, state)
}

// worker processes state updates in background
func (w *AsyncStateWriter) worker() {
	defer w.wg.Done()

	batch := make([]stateUpdate, 0, 100)
	// Reduce ticker interval to 50ms for faster processing under load
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case update := <-w.queue:
			batch = append(batch, update)
			// Flush if batch is full
			if len(batch) >= 100 {
				w.flushBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			// Periodic flush
			if len(batch) > 0 {
				w.flushBatch(batch)
				batch = batch[:0]
			}

		case <-w.shutdown:
			// Final flush on shutdown
			if len(batch) > 0 {
				w.flushBatch(batch)
			}
			return
		}
	}
}

func (w *AsyncStateWriter) flushBatch(batch []stateUpdate) {
	if len(batch) == 0 {
		return
	}

	// Use longer timeout for batch writes (15 seconds for up to 100 items)
	// Each write may need to check current state, so allow more time
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Process writes with individual timeouts to prevent one slow write from blocking others
	for _, update := range batch {
		// Use a shorter timeout per write (1 second) to process faster
		// Allow the batch to take up to 15s total, but individual writes should be fast
		writeCtx, writeCancel := context.WithTimeout(ctx, 1*time.Second)
		
		// Use fast path (no validation) for async writes
		// Validation is already done in the service layer
		err := w.repo.RecordStateFast(writeCtx, update.authID, update.state)
		writeCancel()
		
		if err != nil {
			// Log but continue - don't block other writes
			// Context deadline exceeded is expected under heavy load
			if err == context.DeadlineExceeded {
				log.Printf("WARN: failed to write state %s -> %s: context deadline exceeded (DB under load)", update.authID, update.state.String())
			} else {
				log.Printf("WARN: failed to write state %s -> %s: %v", update.authID, update.state.String(), err)
			}
		}
	}
}

// Close gracefully shuts down the async writer
func (w *AsyncStateWriter) Close() {
	close(w.shutdown)
	w.wg.Wait()
}
