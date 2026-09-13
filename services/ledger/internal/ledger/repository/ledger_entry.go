package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/rs/zerolog"
)

type LedgerEntryRepo struct {
	repo   Repository
	logger *zerolog.Logger
}

func NewLedgerEntryRepo(r Repository, logger *zerolog.Logger) *LedgerEntryRepo {
	return &LedgerEntryRepo{
		repo:   r,
		logger: logger,
	}
}

func (l *LedgerEntryRepo) CreateLedgerEntry(ctx context.Context, entries []domain.LedgerEntry) error {
	if len(entries) == 0 {
		l.logger.Warn().Str("func", "CreateLedgerEntries").Msg("empty ledger entries")
		return fmt.Errorf("no ledger entries provided")
	}

	if err := validateEntries(entries); err != nil {
		l.logger.Error().Err(err).Str("func", "CreateLedgerEntries").Msg("ledger entries do not align")
		return err
	}
	return l.createLedgerEntryBulk(ctx, entries)
}

func (l *LedgerEntryRepo) createLedgerEntryBulk(ctx context.Context, entries []domain.LedgerEntry) error {
	copyCount, err := l.repo.CopyFrom(ctx, pgx.Identifier{"ledger_entries"}, []string{
		"transaction_id", "account_id", "amount", "entry_type",
		"currency_code",
	},
		pgx.CopyFromSlice(len(entries), func(i int) ([]any, error) {
			return []any{entries[i].TransactionID, entries[i].AccountID,
				entries[i].AmountMinorUnits, entries[i].EntryType, entries[i].Currency}, nil
		}))

	if err != nil {
		l.logger.Error().Err(err).Msg("failed to bulk insert ledger entries")
		return err
	}

	if copyCount < int64(len(entries)) {
		l.logger.Error().Msg("failed to insert some entries")
		return fmt.Errorf("failed to insert all entries")
	}

	return nil
}

func validateEntries(entries []domain.LedgerEntry) error {
	var debits, credits int64
	currency := entries[0].Currency

	for _, entry := range entries {
		if entry.Currency != currency {
			return fmt.Errorf("mixed currencies in one transaction: %s VS %s", currency, entry.Currency)
		}
		switch domain.TransactionType(entry.EntryType) {
		case domain.TRANSACTION_CREDIT:
			credits += entry.AmountMinorUnits
		case domain.TRANSACTION_DEBIT:
			debits += entry.AmountMinorUnits
		default:
			return fmt.Errorf("unknown entry type: %s", entry.EntryType)
		}
	}
	if debits != credits {
		return fmt.Errorf("unbalanced ledger entries: debits=%d credits=%d", debits, credits)
	}
	return nil
}
