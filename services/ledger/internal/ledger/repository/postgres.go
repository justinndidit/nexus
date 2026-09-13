package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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

type PostgresStores struct {
	TransactionStore *TransactionRepo
	AccountStore     *AccountRepo
	OutboxStore      *OutboxRepo
	LedgerEntryStore *LedgerEntryRepo
}

func NewPostgresStore(txStore *TransactionRepo,
	accStore *AccountRepo,
	outboxStore *OutboxRepo,
	ledgerStore *LedgerEntryRepo) *PostgresStores {
	return &PostgresStores{
		TransactionStore: txStore,
		AccountStore:     accStore,
		OutboxStore:      outboxStore,
		LedgerEntryStore: ledgerStore,
	}
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

func (tm *PostgresTransactionManager) WithTransaction(ctx context.Context, fn func(stores *PostgresStores) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		tm.logger.Error().Err(err).Str("func", "WithTransaction").Msg("failed to begin transaction")
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	repo := NewPostgresRepo(tm.pool, tm.logger, tx)
	stores := NewPostgresStore(
		NewTransactionRepository(repo, tm.logger),
		NewAccountRepo(repo, tm.logger),
		NewOutBoxRepo(repo, tm.logger),
		NewLedgerEntryRepo(repo, tm.logger))

	if err = fn(stores); err != nil {
		tm.logger.Error().Err(err).Str("func", "WithTransaction").Msg("transaction callback failed, rolling back")
		txErr := tx.Rollback(ctx)
		if txErr != nil {
			tm.logger.Error().Err(txErr).Str("func", "WithTransaction").Msg("failed to rollback transaction")
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		tm.logger.Error().Err(err).Str("func", "WithTransaction").Msg("failed to commit transaction")
		return err
	}
	return nil
}
