package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/rs/zerolog"
)

type TransactionRepo struct {
	repo   Repository
	logger *zerolog.Logger
}

func NewTransactionRepository(r Repository, logger *zerolog.Logger) *TransactionRepo {
	return &TransactionRepo{
		repo:   r,
		logger: logger,
	}
}

func (t *TransactionRepo) CreateTransaction(ctx context.Context, transaction domain.CreateTransactionRequest) (*domain.Transaction, error) {
	stmt := `
		INSERT INTO transactions(from_account_id, destination_account_id,reference, idempotency_key, currency_code, description, status, amount)
		VALUES(@fromAccountID, @destinationAccountID,@reference, @idempotency_key, @currencyCode, @description, @status, @amount)
		RETURNING *
	`

	rows, err := t.repo.Query(ctx, stmt, pgx.NamedArgs{
		"fromAccountID":        transaction.FromAccountID,
		"destinationAccountID": transaction.DestinationAccountID,
		"reference":            transaction.Reference,
		"idempotencyKey":       transaction.IdempotencyKey,
		"currencyCode":         transaction.Currency,
		"description":          transaction.Description,
		"status":               transaction.Status,
		"amount":               transaction.AmountMinorUnits,
	})
	if err != nil {
		t.logger.Error().Err(err).Msg("failed to execute sql statement")
		return nil, err
	}

	newTx, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Transaction])
	if err != nil {
		t.logger.Error().Err(err).Msg("failed to collect row")
		return nil, err
	}

	return &newTx, nil
}

func (t *TransactionRepo) GetTransactionBySessionID(ctx context.Context, sessionID string) (*domain.Transaction, error) {
	stmt := `
		SELECT id, source_account_id, destination_account_id, reference,
		       idempotency_key, currency_code, description, status, amount, created_at
		FROM transactions
		WHERE idempotency_key = @idempotencyKey
	`
	rows, err := t.repo.Query(ctx, stmt, pgx.NamedArgs{
		"idempotencyKey": sessionID,
	})
	if err != nil {
		t.logger.Error().Err(err).Msg("failed to query transaction by idempotency key")
		return nil, err
	}

	tx, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Transaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		t.logger.Error().Err(err).Msg("failed to collect transaction row")
		return nil, err
	}

	return &tx, nil
}
