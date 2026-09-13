package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string
type TransactionType string
type Currency string

const (
	TRANSACTION_PENDING    TransactionStatus = "pending"
	TRANSACTION_FAILED     TransactionStatus = "failed"
	TRANSACTION_COMPLETED  TransactionStatus = "completed"
	TRANSACTION_PROCESSING TransactionStatus = "processing"
	TRANSACTION_CANCELLED  TransactionStatus = "cancelled"

	TRANSACTION_CREDIT TransactionType = "credit"
	TRANSACTION_DEBIT  TransactionType = "debit"
)

type Transaction struct {
	ID                   uuid.UUID  `db:"id"`
	SourceAccountID      uuid.UUID  `db:"source_account_id"`
	DestinationAccountID uuid.UUID  `db:"destination_account_id"`
	Reference            string     `db:"reference"`
	IdempotencyKey       string     `db:"idempotency_key"`
	CurrencyCode         string     `db:"currency_code"`
	Description          string     `db:"description"`
	Status               string     `db:"status"`
	AmountMinorUnits     int64      `db:"amount"`
	CreatedAt            *time.Time `db:"created_at"`
}

type LedgerEntry struct {
	ID               uuid.UUID  `db:"id"`
	TransactionID    uuid.UUID  `db:"transaction_id"`
	AccountID        uuid.UUID  `db:"account_id"`
	EntryType        string     `db:"entry_type"`
	AmountMinorUnits int64      `db:"amount"`
	Currency         string     `db:"currency_code"`
	CreatedAt        *time.Time `db:"created_at"`
}

type Account struct {
	ID                         uuid.UUID  `json:"id" db:"id"`
	UserID                     uuid.UUID  `json:"user_id" db:"user_id"`
	AccountNumber              string     `json:"account_number" db:"account_number"`
	Currency                   string     `json:"currency_code" db:"currency_code"`
	AccountType                string     `json:"account_type" db:"account_type"`
	AccountStatus              string     `json:"account_status" db:"account_status"`
	AvailableBalanceMinorUnits int64      `json:"available_balance" db:"available_balance"`
	LedgerBalanceMinorUnits    int64      `json:"ledger_balance" db:"ledger_balance"`
	Version                    int64      `json:"version" db:"version"`
	CreatedAt                  *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                  *time.Time `json:"updated_at" db:"updated_at"`
}

type OutBoxEvent struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	EventType      string      `json:"event_type" db:"event_type"`
	Payload        interface{} `json:"payload" db:"payload"`
	Status         string      `json:"status" db:"status"`
	IdempotencyKey string      `json:"idempotency_key" db:"idempotency_key"`
	Topic          string      `db:"queue_topic" json:"queue_topic"`
	Priority       int         `json:"priority" db:"priority"`
	Producer       string      `json:"producer" db:"producer"`
	CreatedAt      *time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time  `json:"updated_at" db:"updated_at"`
}
