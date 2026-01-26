CREATE TABLE accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  profile_id UUID NOT NULL,

  account_number VARCHAR(20) NOT NULL UNIQUE,
  currency_code VARCHAR(3) NOT NULL DEFAULT 'NGN',
  account_type VARCHAR(25) NOT NULL DEFAULT 'SAVINGS',
  account_status VARCHAR(25) NOT NULL DEFAULT 'ACTIVE',

  ledger_balance DECIMAL(19,4) NOT NULL DEFAULT 0,
  available_balance DECIMAL(19,4) NOT NULL DEFAULT 0,

  version BIGINT NOT NULL DEFAULT 0,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);


CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  from_account_id UUID NOT NULL,
  destination_account_id UUID NOT NULL

  session_id VARCHAR(50) NOT NULL,
  reference VARCHAR(50),

  currency_code VARCHAR(3) NOT NULL,
  description VARCHAR(255),
  status VARCHAR(25) NOT NULL DEFAULT 'PENDING',
  amount DECIMAL(19,4) NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT uq_transactions_session UNIQUE (session_id),

);

CREATE UNIQUE INDEX uq_transactions_session ON transactions(session_id);
CREATE INDEX idx_transactions_account ON transactions(account_id, created_at DESC);
CREATE INDEX idx_transactions_reference ON transactions(reference);



CREATE TABLE ledger_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  transaction_id UUID NOT NULL,
  account_id UUID NOT NULL,

  amount DECIMAL(19,4) NOT NULL CHECK (amount >= 0),
  entry_type VARCHAR(6) NOT NULL CHECK (entry_type IN ('DEBIT', 'CREDIT')),
  currency_code VARCHAR(3) NOT NULL DEFAULT 'NGN',
  status VARCHAR(15) NOT NULL DEFAULT 'POSTED',

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
);

CREATE INDEX idx_ledger_account ON ledger_entries(account_id, created_at);
CREATE INDEX idx_ledger_transaction ON ledger_entries(transaction_id);



CREATE TABLE outbox_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_type VARCHAR(50) NOT NULL,
  payload JSONB NOT NULL,
  status VARCHAR(15) NOT NULL DEFAULT 'PENDING',
  idempotency_key VARCHAR(50) NOT NULL,
  priority INTEGER DEFAULT 0,
  producer VARCHAR(50) NOT NULL,

  queue_topic VARCHAR(50) NOT NULL,

  retry_count INT NOT NULL DEFAULT 0,
  error_string TEXT,

  locked_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ,

  CONSTRAINT uq_outbox_idempo UNIQUE (idempotency_key)
);
CREATE INDEX idx_outbox_pending ON outbox_events(status, priority DESC, created_at);
CREATE UNIQUE INDEX uq_outbox_idempo ON outbox_events(producer, idempotency_key);
CREATE INDEX idx_outbox_fetch_worker ON outbox_events(status, priority DESC, created_at) WHERE status = 'PENDING';