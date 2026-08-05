package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PostgresTxManager struct {
	pool *pgxpool.Pool
}

func (t *PostgresTxManager) WithTransaction(ctx context.Context, fn func(uow UnitOfWork) error) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	uow := NewPostgresUnitOfWork(t.pool)
	if err = fn(uow); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

type PostgresUnitOfWork struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

func NewPostgresUnitOfWork(pool *pgxpool.Pool) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{
		pool: pool,
	}
}

func (u *PostgresUnitOfWork) Users() UserRepository {
	return NewPostgresUserRepository(u.logger, u.pool)
}
