package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"cardnet/internal/db"
	"cardnet/internal/settlement/domain"
	"cardnet/internal/settlement/reporting"
	"cardnet/internal/settlement/repository/postgres"
	"cardnet/internal/settlement/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSettlementPipeline(t *testing.T) {
	// 1. Setup DB Connection
	// Ensure DB_DSN is set or use default localhost
	if os.Getenv("DB_DSN") == "" {
		os.Setenv("DB_DSN", "postgres://cardnet:cardnet@localhost:5432/cardnet?sslmode=disable")
	}

	pgDB, err := db.NewPostgres()
	if err != nil {
		t.Skipf("Skipping integration test: failed to connect to postgres: %v", err)
	}
	defer pgDB.Close()

	// Initialize Schema (idempotent)
	initSettlementSchema(t, pgDB)

	// Clean up tables before test
	cleanTables(t, pgDB)
	defer cleanTables(t, pgDB)

	ctx := context.Background()

	// 2. Wiring
	stmRepo := postgres.NewSettlementRepository(pgDB)
	batchRepo := postgres.NewBatchRepository(pgDB)
	reconRepo := postgres.NewReconciliationRepository(pgDB)
	reportGen := reporting.NewGenerator(pgDB)

	stmService := service.NewSettlementService(stmRepo)
	batchProcessor := service.NewBatchProcessor(batchRepo)
	reconService := service.NewReconciliationService(reconRepo, reportGen)

	// 3. Prepare Data
	// "Yesterday" relative to T0 execution logic, but here we just pick a specific date
	// and ensure we create authorizations for that time window.
	targetDate := time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)

	// Create a few authorizations
	auth1 := domain.SettlementAuthorization{
		AuthorizationID: uuid.New().String(),
		MerchantID:      "merch_001",
		CardHash:        "hash_1234",
		Amount:          1000, // 10.00
		Currency:        "USD",
		Approved:        true,
		EventTime:       targetDate.Add(1 * time.Hour), // 01:00 AM on target date
		CreatedAt:       time.Now(),
	}

	auth2 := domain.SettlementAuthorization{
		AuthorizationID: uuid.New().String(),
		MerchantID:      "merch_001",
		CardHash:        "hash_5678",
		Amount:          2550, // 25.50
		Currency:        "USD",
		Approved:        true,
		EventTime:       targetDate.Add(23 * time.Hour), // 11:00 PM on target date
		CreatedAt:       time.Now(),
	}

	// Insert them
	require.NoError(t, stmService.ProcessAuthorization(ctx, auth1))
	require.NoError(t, stmService.ProcessAuthorization(ctx, auth2))

	// 4. Run Batch Processing (T1 Job logic)
	// This should pick up everything from targetDate to targetDate+24h
	batchID, err := batchProcessor.RunForDate(ctx, targetDate)
	require.NoError(t, err)
	require.NotEmpty(t, batchID)
	t.Logf("Created Batch ID: %s", batchID)

	// Verify Batch Status is READY
	var status string
	err = pgDB.QueryRow("SELECT status FROM settlement_batches WHERE id=$1", batchID).Scan(&status)
	require.NoError(t, err)
	require.Equal(t, "READY", status)

	// Verify Authorizations are linked
	var count int
	err = pgDB.QueryRow("SELECT COUNT(*) FROM settlement_authorizations WHERE batch_id=$1", batchID).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	// 5. Run Reconciliation
	err = reconService.ReconcileBatch(ctx, batchID)
	require.NoError(t, err)

	// Verify Reconciliation Status
	var reconStatus string
	err = pgDB.QueryRow("SELECT reconciliation_status FROM settlement_batches WHERE id=$1", batchID).Scan(&reconStatus)
	require.NoError(t, err)
	require.Equal(t, "Reconciled", reconStatus)

	// Verify Report Generation
	// Check if file exists in reports/
	reportPath := fmt.Sprintf("reports/batch_%s.csv", batchID)
	_, err = os.Stat(reportPath)
	require.NoError(t, err, "Report file should exist")

	// Cleanup report file
	_ = os.Remove(reportPath)
}

func initSettlementSchema(t *testing.T, db *sql.DB) {
	// Basic Schema required for settlement
	queries := []string{
		`CREATE TABLE IF NOT EXISTS settlement_batches (
			id TEXT PRIMARY KEY,
			batch_date DATE NOT NULL UNIQUE,
			status TEXT NOT NULL,
			reconciliation_status TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS settlement_authorizations (
			authorization_id TEXT PRIMARY KEY,
			merchant_id TEXT NOT NULL,
			card_hash TEXT NOT NULL,
			amount BIGINT NOT NULL,
			currency TEXT NOT NULL,
			approved BOOLEAN NOT NULL,
			event_time TIMESTAMP NOT NULL,
			batch_id TEXT REFERENCES settlement_batches(id),
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_settlement_auth_event_time ON settlement_authorizations(event_time);`,
		`CREATE INDEX IF NOT EXISTS idx_settlement_auth_batch_id ON settlement_authorizations(batch_id);`,
		`CREATE TABLE IF NOT EXISTS settlement_reconciliation_issues (
			id SERIAL PRIMARY KEY,
			batch_id TEXT NOT NULL REFERENCES settlement_batches(id),
			reason TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		require.NoError(t, err)
	}
}

func cleanTables(t *testing.T, db *sql.DB) {
	tables := []string{
		"settlement_reconciliation_issues",
		"settlement_authorizations",
		"settlement_batches",
	}
	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		require.NoError(t, err)
	}
}
