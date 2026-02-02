package ledger

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/rs/zerolog"
)

type PostgresRepository struct {
	pool   *pgxpool.Pool
	tx     pgx.Tx //pgx.Tx is an interface
	logger *zerolog.Logger
}

type PostgresTransactionManager struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

func NewPostgresRepo(pool *pgxpool.Pool, logger *zerolog.Logger, tx pgx.Tx) *PostgresRepository {
	return &PostgresRepository{
		pool:   pool,
		logger: logger,
		tx:     tx,
	}
}

func NewPostgresTransactionManager(pool *pgxpool.Pool, logger *zerolog.Logger) *PostgresTransactionManager {
	return &PostgresTransactionManager{
		pool:   pool,
		logger: logger,
	}
}

func (pr *PostgresRepository) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if pr.tx != nil {
		return pr.tx.Exec(ctx, sql, args...)
	}

	return pr.pool.Exec(ctx, sql, args...)
}

func (pr *PostgresRepository) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if pr.tx != nil {
		return pr.tx.Query(ctx, sql, args...)
	}

	return pr.pool.Query(ctx, sql, args...)
}

func (pr *PostgresRepository) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if pr.tx != nil {
		return pr.tx.QueryRow(ctx, sql, args...)
	}

	return pr.pool.QueryRow(ctx, sql, args...)
}

func (pr *PostgresRepository) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	if pr.tx != nil {
		return pr.tx.CopyFrom(ctx, tableName, columnNames, rowSrc)
	}

	return pr.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)

}

func (pr *PostgresRepository) CreateTransaction(ctx context.Context, transaction domain.CreateTransactionRequest) (*domain.Transaction, error) {
	stmt := `
		INSERT INTO transactions(from_account_id, destination_account_id,reference, session_id, currency_code, description, status, amount)
		VALUES(@accountID, @destinationAccountID,@reference, @sessionID, @currencyCode, @description, @status, @amount)
		RETURNING *
	`

	rows, err := pr.Query(ctx, stmt, pgx.NamedArgs{
		"fromAccountID":        transaction.FromAccountID,
		"destinationAccountID": transaction.DestinationAccountID,
		"reference":            transaction.Reference,
		"sessionID":            transaction.SessionID,
		"currencyCode":         transaction.Currency,
		"description":          transaction.Description,
		"status":               transaction.Status,
		"amount":               transaction.AmountMinDenomination,
	})
	if err != nil {
		pr.logger.Error().Err(err).Msg("failed to execute sql statement")
		return nil, err
	}

	newTx, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Transaction])
	if err != nil {
		pr.logger.Error().Err(err).Msg("failed to collect row")
		return nil, err
	}

	return &newTx, nil

}

func (pr *PostgresRepository) CreateLedgerEntry(ctx context.Context, entries []domain.LedgerEntry) error {
	stmt := `
		INSERT INTO Ledger_entries(transaction_id, account_id, entry_type, amount, currency, idempotency_key)
		VALUES (@transaction_id, @account_id, @entry_type, @amount, @currency, @idempotency_key)
	`
	for _, entry := range entries {
		cmd, err := pr.Exec(ctx, stmt, pgx.NamedArgs{
			"transaction_id":  entry.TransactionID,
			"account_id":      entry.AccountID,
			"amount":          entry.AmountMinDenomination,
			"currency":        entry.Currency,
			"idempotency_key": entry.IdempotencyKey,
		})

		if err != nil {
			pr.logger.Error().Err(err).Msgf("failed to execute sql: %s", stmt)
			return err
		}
		if cmd.RowsAffected() == 0 {
			return fmt.Errorf("failed to insert ledger entry: %v", entry)
		}
	}
	return nil
}

func (pr *PostgresRepository) CreateLedgerEntryBulk(ctx context.Context, entries []domain.LedgerEntry) error {

	copyCount, err := pr.CopyFrom(ctx, pgx.Identifier{"ledger_entry"}, []string{
		"transaction_id", "session_id", "account_id", "amount", "entry_type",
		"currency", "status",
	},
		pgx.CopyFromSlice(len(entries), func(i int) ([]any, error) {
			return []any{entries[i].TransactionID, entries[i].IdempotencyKey, entries[i].AccountID,
				entries[i].AmountMinDenomination, entries[i].EntryType, entries[i].Currency, entries[i].Status}, nil
		}))

	if err != nil {
		pr.logger.Error().Err(err).Msg("failed to bulk insert ledger entries")
		return err
	}

	if copyCount < int64(len(entries)) {
		pr.logger.Error().Msg("failed to insert some entries")
		return fmt.Errorf("failed to insert all entries")
	}

	return nil
}

func (pr *PostgresRepository) GetAccountForUpdate(ctx context.Context, accountID string) (*domain.Account, error) {
	var account domain.Account
	stmt := `
		SELECT (id,available_balance, ledger_balance, version)
		FROM accounts
		WHERE account_id = @account_id
		FOR UPDATE
	`
	err := pr.QueryRow(ctx, stmt, pgx.NamedArgs{
		"account_id": accountID,
	}).Scan(&account)

	if err != nil {
		pr.logger.Error().Err(err).Msg("failed to retrieve account")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("account with id %s does not exist", accountID)
		}
		return nil, err
	}

	return &account, nil
}

func (pr *PostgresRepository) UpdateBalance(ctx context.Context, accountID string, amount int64) error {
	stmt := `
		UPDATE accounts SET available_balance = available_balance + @amount
		WHERE id = @accountID
	`
	cmd, err := pr.Exec(ctx, stmt, pgx.NamedArgs{
		"amount":    amount,
		"accountID": accountID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("account with id %s does not exist", accountID)
		}
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("failed to update account with id: %s ", accountID)
	}

	return nil
}

func (pr *PostgresRepository) CreateOutBoxEvent(ctx context.Context, payload domain.CreateOutboxEventRequest) error {
	stmt := `
		INSERT INTO outbox_events(event_type, payload, status, idempotency_key, priority, producer)
		VALUES (@eventType, @payload, @status, @idempotency_key, @priority, @producer)
	`
	cmd, err := pr.Exec(ctx, stmt, pgx.NamedArgs{
		"payload":         payload.Payload,
		"idempotency_key": payload.IdempotencyKey,
		"eventType":       payload.EventType,
		"producer":        payload.Producer,
		"priority":        int(payload.Priority),
		"status":          payload.Status,
	})
	if err != nil {
		pr.logger.Error().Err(err).Msg("failed to insert event")
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected")
	}

	return nil
}

// func (pr *PostgresRepository) GetOutBoxEvent(ctx context.Context, ID string) (*domain.OutBoxEvent, error) {
// 	return nil, nil
// }

func (pr *PostgresRepository) GetOutBoxEventsForUpdate(ctx context.Context) ([]domain.OutBoxEvent, error) {
	stmt := `
			UPDATE outbox_events
			SET status = 'PROCESSING',
					locked_at = NOW()
			WHERE id IN (
					SELECT id FROM outbox_events
					WHERE status = 'PENDING'
					ORDER BY created_at ASC
					LIMIT 100
					FOR UPDATE SKIP LOCKED
			)
			RETURNING id, event_type, payload, status, idempotency_key, priority, producer, created_at, updated_at
	`
	rows, err := pr.Query(ctx, stmt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pr.logger.Info().Msg("no pending tasks....")
			return nil, nil
		}
		pr.logger.Error().Err(err).Msgf("failed to execute statement %s", stmt)
		return nil, err
	}

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.OutBoxEvent])
	if err != nil {
		pr.logger.Error().Err(err).Msg("failed to collect events")
		return nil, err
	}

	return events, nil
}

func (pr *PostgresRepository) IncrementRetryCount(ctx context.Context, id string, err_string string) error {
	stmt := `
		UPDATE outbox_events
		SET retry_count = retry_count + 1,
		error_string = @err_string
		status = 'PENDING
		WHERE id = @id
	`

	rows, err := pr.Exec(ctx, stmt, pgx.NamedArgs{
		"id":         id,
		"err_string": err_string,
	})
	if err != nil {
		pr.logger.Error().Err(err).Msgf("failed to update %s for retry", id)
		return err
	}

	if rows.RowsAffected() == 0 {
		pr.logger.Warn().Msg("no rows updated")
		return fmt.Errorf("failed to update %s for retry", id)
	}

	return nil

}

func (pr *PostgresRepository) MarkEventProcessed(ctx context.Context, id string) error {
	stmt := `
			UPDATE outbox_events
			SET status = 'PROCESSED',
					locked_at = NOW()
			WHERE id = @id
	`
	rows, err := pr.Exec(ctx, stmt, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		pr.logger.Error().Err(err).Msgf("failed to update %s to \"PROCESSED\"", id)
		return err
	}

	if rows.RowsAffected() == 0 {
		pr.logger.Warn().Msg("no rows updated")
		return fmt.Errorf("failed to mark %s as processed", id)
	}

	return nil
}

func (tm *PostgresTransactionManager) WithTransaction(ctx context.Context, fn func(repo Repository) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		tm.logger.Error().Err(err).Msg("failed to begin transaction")
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	repo := NewPostgresRepo(tm.pool, tm.logger, tx)

	if err = fn(repo); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
