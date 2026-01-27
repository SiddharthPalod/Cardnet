package integration

import (
	"cardnet/internal/ledger/domain"
	"cardnet/internal/ledger/repository/cassandra"
	"context"
	"testing"
	"time"

	ledgerpb "cardnet/internal/api/ledger"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLedgerRepository_StoreEvent(t *testing.T) {
	// ---- Cassandra connection ----
	session, err := cassandra.NewSession(cassandra.Config{
		Hosts:       []string{"127.0.0.1"},
		Keyspace:    "cardnet_ledger",
		Consistency: gocql.LocalQuorum,
		Timeout:     2 * time.Second,
	})
	require.NoError(t, err)
	defer session.Close()

	repo := cassandra.NewLedgerRepo(session)

	// ---- Test event ----
	eventTime := time.Now().UTC()
	event := &domain.AuthEvent{
		EventID:    gocql.UUID(uuid.New()),
		AuthID:     "auth_test_123",
		MerchantID: "merchant_abc",
		CardHash:   "cardhash_xyz",
		EventType:  ledgerpb.EventType_RISK_DECISION,
		Decision:   ledgerpb.Decision_APPROVED,
		Payload:    []byte(`{"score":0.91}`),
		CreatedAt:  eventTime,
	}

	err = repo.Store(context.Background(), event)
	require.NoError(t, err)

	// ---- Assert: verify one table----
	var (
		authID   string
		eventTyp string
	)
	err = session.Query(`
		SELECT auth_id, event_type
		FROM auth_events_by_merchant
		WHERE merchant_id = ?
		LIMIT 1
	`, event.MerchantID).
		Scan(&authID, &eventTyp)

	require.NoError(t, err)
	require.Equal(t, event.AuthID, authID)
	require.Equal(t, event.EventType.String(), eventTyp)
}
