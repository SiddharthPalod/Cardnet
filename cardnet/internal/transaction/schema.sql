-- Transaction State Ledger Table
-- Stores canonical transaction states with idempotency guarantees
CREATE TABLE IF NOT EXISTS transaction_states (
    auth_id TEXT NOT NULL,
    state TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (auth_id, state),
    CONSTRAINT valid_state CHECK (state IN ('INITIATED', 'RISK_EVALUATED', 'ISSUER_REQUESTED', 'ISSUER_RESPONDED', 'FINALIZED'))
);

CREATE INDEX IF NOT EXISTS idx_transaction_states_auth_id ON transaction_states(auth_id);
CREATE INDEX IF NOT EXISTS idx_transaction_states_timestamp ON transaction_states(timestamp);
