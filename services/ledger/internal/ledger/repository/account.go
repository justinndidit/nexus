package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/rs/zerolog"
)

type AccountRepo struct {
	repo   Repository
	logger *zerolog.Logger
}

func NewAccountRepo(r Repository, logger *zerolog.Logger) *AccountRepo {
	return &AccountRepo{
		repo:   r,
		logger: logger,
	}
}

func (a *AccountRepo) GetAccountForUpdate(ctx context.Context, accountID string) (*domain.Account, error) {
	stmt := `
		SELECT id, available_balance, ledger_balance, version
		FROM accounts
		WHERE id = @account_id
		FOR UPDATE
	`
	rows, err := a.repo.Query(ctx, stmt, pgx.NamedArgs{
		"account_id": accountID,
	})

	if err != nil {
		a.logger.Error().Err(err).Msg("failed to retrieve account")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("account with id %s does not exist", accountID)
		}
		return nil, err
	}

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Account])
	if err != nil {
		a.logger.Error().Err(err).Str("func", "GetAccountForUpdate").Msg("failed to collect account")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("account with id %s does not exist", accountID)
		}
		return nil, err
	}

	return &account, nil
}

func (a AccountRepo) UpdateBalance(ctx context.Context, accountID string, amount int64) error {
	stmt := `
		UPDATE accounts SET available_balance = available_balance + @amount
		WHERE id = @accountID
	`
	cmd, err := a.repo.Exec(ctx, stmt, pgx.NamedArgs{
		"amount":    amount,
		"accountID": accountID,
	})
	if err != nil {
		a.logger.Error().Err(err).Str("func", "UpdateBalance").Msgf("failed to update account with id: %s", accountID)
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("account with id %s does not exist", accountID)
		}
		return err
	}

	if cmd.RowsAffected() == 0 {
		a.logger.Error().Str("func", "UpdateBalance").Msgf("failed to update account with id %s", accountID)
		return fmt.Errorf("failed to update account with id: %s ", accountID)
	}

	return nil
}
