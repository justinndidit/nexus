package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justinndidit/nexus/user/internal/user/domain"
	"github.com/rs/zerolog"
)

type PostgresUserRepository struct {
	logger *zerolog.Logger
	pool   *pgxpool.Pool
	tx     *pgxpool.Tx
}

func NewPostgresUserRepository(logger *zerolog.Logger, pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		logger: logger,
		pool:   pool,
	}
}

func (pr *PostgresUserRepository) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if pr.tx != nil {
		return pr.tx.Exec(ctx, sql, args...)
	}

	return pr.pool.Exec(ctx, sql, args...)
}

func (pr *PostgresUserRepository) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if pr.tx != nil {
		return pr.tx.Query(ctx, sql, args...)
	}

	return pr.pool.Query(ctx, sql, args...)
}

func (pr *PostgresUserRepository) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if pr.tx != nil {
		return pr.tx.QueryRow(ctx, sql, args...)
	}

	return pr.pool.QueryRow(ctx, sql, args...)
}

func (pr *PostgresUserRepository) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	if pr.tx != nil {
		return pr.tx.CopyFrom(ctx, tableName, columnNames, rowSrc)
	}

	return pr.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)

}

func (r PostgresUserRepository) CreateUser(ctx context.Context, user domain.RegisterUserDTO) (*domain.UserDTO, error) {
	stmt := `
		INSERT INTO users(email, username, phone_number,password)
		VALUES(@email, @username, @phone_number, @password)
		RETURNING *
	`

	rows, err := r.Query(ctx, stmt, pgx.NamedArgs{
		"email":        user.Email,
		"password":     user.Password,
		"phone_number": user.PhoneNumber,
		"username":     user.Username,
	})

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to execute sql statement")
		return nil, err
	}

	newUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.UserDTO])
	if err != nil {
		r.logger.Error().Err(err).Msg("error collecting row")
		return nil, err
	}

	return &newUser, nil
}

func (r PostgresUserRepository) FetchUserByEmail(ctx context.Context, email string) (*domain.UserDTO, error) {
	stmt := `
		SELECT (id, email,phone_number, kyc_status, is_email_verified, is_phone_verified, is_blocked, is_suspended, blocked_reason, blocked_by, blocked_at, last_login_at)
		FROM users
		WHERE email = @email
	`
	rows, err := r.Query(ctx, stmt, pgx.NamedArgs{
		"email": email,
	})

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to execute sql statement")
		return nil, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.UserDTO])
	if err != nil {
		r.logger.Error().Err(err).Msg("error collecting row")
		return nil, err
	}
	return &user, nil
}
