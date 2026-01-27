CREATE TABLE IF NOT EXISTS settlement_batches (
    id TEXT PRIMARY KEY,
    batch_date DATE NOT NULL UNIQUE,
    status TEXT NOT NULL,
    reconciliation_status TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS settlement_authorizations (
    authorization_id TEXT PRIMARY KEY,
    merchant_id TEXT NOT NULL,
    card_hash TEXT NOT NULL,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    approved BOOLEAN NOT NULL,
    event_time TIMESTAMP NOT NULL,
    batch_id TEXT REFERENCES settlement_batches(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_settlement_auth_event_time ON settlement_authorizations(event_time);
CREATE INDEX IF NOT EXISTS idx_settlement_auth_batch_id ON settlement_authorizations(batch_id);

CREATE TABLE IF NOT EXISTS settlement_reconciliation_issues (
    id SERIAL PRIMARY KEY,
    batch_id TEXT NOT NULL REFERENCES settlement_batches(id),
    reason TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
