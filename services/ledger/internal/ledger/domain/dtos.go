package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransferMetaData map[string]interface{}
type OutboxEventType string
type OutboxEventStatus string
type OutboxEventPriority int

const (
	EventMoneyTransfer OutboxEventType = "event_money_transfer"

	OutboxEventQueued     OutboxEventStatus = "queued"
	OutboxEventPending    OutboxEventStatus = "pending"
	OutboxEventProcessing OutboxEventStatus = "processing"
	OutboxEventCompleted  OutboxEventStatus = "completed"
	OutboxEventFailed     OutboxEventStatus = "failed"

	OutboxPriorityHigh   OutboxEventPriority = 10
	OutboxPriorityMedium OutboxEventPriority = 5
	OutboxPriorityLow    OutboxEventPriority = 0
)

type Money struct {
	Currency string `json:"currency_code" db:"currency_code" validate:"required,len=3"`
	Unit     int64  `json:"unit" db:"unit" validate:"required"`
	Nanos    int64  `json:"nanos" db:"nanos" validate:"required,gte=0,lte=999999999"`
}

type TransferRequest struct {
	FromAccountID        uuid.UUID        `json:"from_account_id" db:"from_account_id" validate:"required"`
	DestinationAccountID uuid.UUID        `json:"destination_account_id" db:"destination_account_id" validate:"required"`
	IdempotencyKey       string           `json:"idempotency_key" db:"idempotency_key" validate:"required"`
	Money                Money            `json:"money" validate:"required"`
	Meta                 TransferMetaData `json:"meta,omitempty" db:"meta"`
}

type TransferResponse struct {
	TransactionID string            `json:"transaction_id" db:"transaction_id"`
	Status        TransactionStatus `json:"status" db:"status"`
}

type CreateTransactionRequest struct {
	FromAccountID        uuid.UUID       `db:"from_account_id" json:"from_account_id"`
	DestinationAccountID uuid.UUID       `db:"destination_Account_id" json:"destination_Account_id"`
	Reference            string          `db:"reference" json:"reference,omitempty"`
	SessionID            string          `db:"session_id" json:"session_id"`
	Currency             string          `db:"currency_code" json:"currency_code"`
	Description          string          `db:"description" json:"description omitempty"`
	Status               string          `db:"status" json:"status"`
	Amount               decimal.Decimal `db:"amount" json:"amount"`
}

type CreateOutboxEventRequest struct {
	EventType      OutboxEventType     `db:"event_type" json:"event_type"`
	Payload        interface{}         `db:"payload" json:"payload"`
	Status         OutboxEventStatus   `db:"status" json:"status"`
	IdempotencyKey string              `db:"idempotency_key" json:"idempotency_key"`
	Priority       OutboxEventPriority `db:"priority" json:"priority"`
	Producer       string              `db:"producer" json:"producer"`
}

type TransactionEventPayload struct {
	FromAccountID        uuid.UUID       `json:"from_account_id" db:"from_account_id"`
	DestinationAccountID uuid.UUID       `json:"destination_account_id" db:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount" db:"amount"`
	Currency             string          `json:"currency_code" db:"currency_code"`
}
