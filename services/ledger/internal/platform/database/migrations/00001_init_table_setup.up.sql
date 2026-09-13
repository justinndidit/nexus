CREATE TABLE accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,

  account_number VARCHAR(20) NOT NULL UNIQUE,
  currency_code VARCHAR(3) NOT NULL,
  account_type VARCHAR(25) NOT NULL DEFAULT 'SAVINGS',
  account_status VARCHAR(25) NOT NULL DEFAULT 'ACTIVE',

  ledger_balance BIGINT NOT NULL DEFAULT 0 CHECK (ledger_balance >= 0),
  available_balance BIGINT NOT NULL DEFAULT 0 CHECK (available_balance >= 0),

  version BIGINT NOT NULL DEFAULT 0,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);


CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source_account_id UUID NOT NULL REFERENCES accounts(id),
  destination_account_id UUID NOT NULL REFERENCES accounts(id),

  idempotency_key VARCHAR(50) NOT NULL,
  reference VARCHAR(50),

  currency_code VARCHAR(3) NOT NULL,
  description VARCHAR(255),
  status VARCHAR(25) NOT NULL DEFAULT 'PENDING',
  amount BIGINT NOT NULL CHECK (amount > 0),

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT uq_transactions_idempotency_key UNIQUE (idempotency_key)
);
CREATE INDEX idx_transactions_source_account ON transactions(source_account_id, created_at DESC);
CREATE INDEX idx_transactions_destination_account ON transactions(destination_account_id, created_at DESC);
CREATE INDEX idx_transactions_reference ON transactions(reference);


CREATE TABLE ledger_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  transaction_id UUID NOT NULL REFERENCES transactions(id),
  account_id UUID NOT NULL REFERENCES accounts(id),

  amount BIGINT NOT NULL CHECK (amount > 0),
  entry_type TEXT NOT NULL CHECK (entry_type IN ('DEBIT', 'CREDIT')),
  currency_code VARCHAR(3) NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_ledger_account ON ledger_entries(account_id, created_at);
CREATE INDEX idx_ledger_transaction ON ledger_entries(transaction_id);

CREATE OR REPLACE FUNCTION check_ledger_balance() RETURNS TRIGGER AS $$
DECLARE
  debit_total BIGINT;
  credit_total BIGINT;
BEGIN
  SELECT
    COALESCE(SUM(amount) FILTER (WHERE entry_type = 'DEBIT'), 0),
    COALESCE(SUM(amount) FILTER (WHERE entry_type = 'CREDIT'), 0)
  INTO debit_total, credit_total
  FROM ledger_entries
  WHERE transaction_id = NEW.transaction_id;

  IF debit_total != credit_total THEN
    RAISE EXCEPTION 'unbalanced ledger entries for transaction %: debits=% credits=%',
      NEW.transaction_id, debit_total, credit_total;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER ledger_balance_check
  AFTER INSERT ON ledger_entries
  DEFERRABLE INITIALLY DEFERRED
  FOR EACH ROW
  EXECUTE FUNCTION check_ledger_balance();

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

  CONSTRAINT uq_outbox_producer_idempotency_key UNIQUE (producer, idempotency_key)
);
CREATE INDEX idx_outbox_fetch_worker ON outbox_events(status, priority DESC, created_at) WHERE status = 'PENDING';