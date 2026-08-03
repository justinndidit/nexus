CREATE TABLE accounts IF NOT EXISTS(
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
    from_account_id UUID NOT NULL, -- Added type
    destination_account_id UUID NOT NULL, -- Added type
    currency_code VARCHAR(3) NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    reference VARCHAR(100) UNIQUE, -- Added this for your index
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_from_acc ON transactions(from_account_id, created_at DESC);
CREATE INDEX idx_transactions_reference ON transactions(reference);
CREATE TABLE IF NOT EXISTS processed_event(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPZ NOT NULL DEFAULT now()
);