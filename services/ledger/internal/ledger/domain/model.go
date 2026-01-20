package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	ID                   uuid.UUID       `db:"id" json:"id"`
	FromAccountID        uuid.UUID       `db:"from_account_id" json:"from_account_id"`
	DestinationAccountID uuid.UUID       `db:"destination_Account_id" json:"destination_Account_id"`
	Reference            string          `db:"reference" json:"reference omitempty"`
	SessionID            string          `db:"session_id" json:"session_id"`
	Currency             string          `db:"currency_code" json:"currency_code"`
	Description          string          `db:"description" json:"description omitempty"`
	Status               string          `db:"status" json:"status"`
	Amount               decimal.Decimal `db:"amount" json:"amount"`
	CreatedAt            *time.Time      `db:"created_at" json:"created_at"`
}

type LedgerEntry struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	TransactionID  uuid.UUID       `json:"transaction_id" db:"transaction_id"`
	AccountID      uuid.UUID       `json:"account_id" db:"account_id"`
	EntryType      string          `json:"entry_type" db:"entry_type"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`
	Currency       string          `json:"currency_code" db:"currency_code"`
	IdempotencyKey string          `json:"idempotency_key" db:"idempotency_key"`
	Status         string          `json:"status" db:"status"`
	CreatedAt      *time.Time      `json:"created_at" db:"created_at"`
}

type Account struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	UserID           uuid.UUID       `json:"user_id" db:"user_id"`
	ProfileID        uuid.UUID       `json:"profile_id" db:"profile_id"`
	AccountNumber    string          `json:"account_number" db:"account_number"`
	Currency         string          `json:"currency_code" db:"currency_code"`
	AccountType      string          `json:"account_type" db:"account_type"`
	AccountStatus    string          `json:"account_status" db:"account_status"`
	AvailableBalance decimal.Decimal `json:"available_balance" db:"available_balance"`
	LedgerBalance    decimal.Decimal `json:"ledger_balance" db:"ledger_balance"`
	Version          int64           `json:"version" db:"version"`
	CreatedAt        *time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        *time.Time      `json:"updated_at" db:"updated_at"`
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
