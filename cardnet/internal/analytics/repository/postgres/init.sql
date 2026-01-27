CREATE TABLE merchant_stats (
  merchant_id TEXT PRIMARY KEY,
  tx_count INT,
  approved_count INT DEFAULT 0,
  total_amount NUMERIC
);

CREATE TABLE bin_stats (
  bin TEXT PRIMARY KEY,
  tx_count INT DEFAULT 0,
  approved_count INT DEFAULT 0,
  total_amount NUMERIC DEFAULT 0
);

CREATE TABLE rule_stats (
  rule_id TEXT PRIMARY KEY,
  triggered_count INT DEFAULT 0,
  decline_count INT DEFAULT 0
);

CREATE TABLE network_metrics (
  total_tx INT DEFAULT 0,
  approved_tx INT DEFAULT 0
);

INSERT INTO network_metrics DEFAULT VALUES;

-- Batch Aggregation Tables for Historical Analytics
CREATE TABLE hourly_merchant_aggregations (
  id SERIAL PRIMARY KEY,
  merchant_id TEXT NOT NULL,
  hour_bucket TIMESTAMP NOT NULL,
  tx_count INT DEFAULT 0,
  approved_count INT DEFAULT 0,
  total_amount NUMERIC DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(merchant_id, hour_bucket)
);

CREATE TABLE hourly_bin_aggregations (
  id SERIAL PRIMARY KEY,
  bin TEXT NOT NULL,
  hour_bucket TIMESTAMP NOT NULL,
  tx_count INT DEFAULT 0,
  approved_count INT DEFAULT 0,
  total_amount NUMERIC DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(bin, hour_bucket)
);

CREATE TABLE hourly_network_aggregations (
  id SERIAL PRIMARY KEY,
  hour_bucket TIMESTAMP NOT NULL UNIQUE,
  total_tx INT DEFAULT 0,
  approved_tx INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE daily_merchant_aggregations (
  id SERIAL PRIMARY KEY,
  merchant_id TEXT NOT NULL,
  day_bucket DATE NOT NULL,
  tx_count INT DEFAULT 0,
  approved_count INT DEFAULT 0,
  total_amount NUMERIC DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(merchant_id, day_bucket)
);

CREATE TABLE daily_bin_aggregations (
  id SERIAL PRIMARY KEY,
  bin TEXT NOT NULL,
  day_bucket DATE NOT NULL,
  tx_count INT DEFAULT 0,
  approved_count INT DEFAULT 0,
  total_amount NUMERIC DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(bin, day_bucket)
);

CREATE TABLE daily_network_aggregations (
  id SERIAL PRIMARY KEY,
  day_bucket DATE NOT NULL UNIQUE,
  total_tx INT DEFAULT 0,
  approved_tx INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for efficient querying
CREATE INDEX idx_hourly_merchant_hour ON hourly_merchant_aggregations(hour_bucket);
CREATE INDEX idx_hourly_merchant_id ON hourly_merchant_aggregations(merchant_id);
CREATE INDEX idx_hourly_bin_hour ON hourly_bin_aggregations(hour_bucket);
CREATE INDEX idx_daily_merchant_day ON daily_merchant_aggregations(day_bucket);
CREATE INDEX idx_daily_merchant_id ON daily_merchant_aggregations(merchant_id);
CREATE INDEX idx_daily_bin_day ON daily_bin_aggregations(day_bucket);

-- Raw network authorization history for analytics UI
CREATE TABLE IF NOT EXISTS network_history (
  auth_id TEXT PRIMARY KEY,
  merchant_id TEXT NOT NULL,
  amount NUMERIC NOT NULL,
  currency TEXT NOT NULL,
  status TEXT NOT NULL,
  reason TEXT,
  created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_network_history_created_at ON network_history(created_at);
CREATE INDEX IF NOT EXISTS idx_network_history_merchant_id ON network_history(merchant_id);
CREATE INDEX IF NOT EXISTS idx_network_history_status ON network_history(status);
