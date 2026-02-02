package ledger

import (
	"context"

	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
)

type Repository interface {
	CreateLedgerEntry(context.Context, []domain.LedgerEntry) error
	CreateTransaction(context.Context, domain.CreateTransactionRequest) (*domain.Transaction, error)
	GetAccountForUpdate(context.Context, string) (*domain.Account, error)
	UpdateBalance(context.Context, string, int64) error
	CreateOutBoxEvent(context.Context, domain.CreateOutboxEventRequest) error
	GetOutBoxEventsForUpdate(context.Context) ([]domain.OutBoxEvent, error)
	IncrementRetryCount(context.Context, string, string) error
	MarkEventProcessed(context.Context, string) error
}

type TransactionManager interface {
	WithTansaction(context.Context, func(repo Repository) error) error
}
