package ledger

import (
	"context"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/justinndidit/nexus/ledger/internal/platform/utils"
	"github.com/rs/zerolog"
)

type LedgerService struct {
	repo      Repository
	txManager TransactionManager
	validator validator.Validate
	logger    *zerolog.Logger
}

func NewLegerService(r Repository, txManager TransactionManager, v validator.Validate, log *zerolog.Logger) *LedgerService {
	return &LedgerService{
		repo:      r,
		txManager: txManager,
		validator: v,
		logger:    log,
	}
}

func (s LedgerService) Transfer(ctx context.Context, req domain.TransferRequest) error {
	err := s.validator.Struct(req)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to validate transfer request data")
		return err
	}

	if s.isInterBankTransfer(req) {
		return s.InterBankTransfer(ctx, req)
	}

	return s.IntraBankTransfer(ctx, req)
}

func (s LedgerService) IntraBankTransfer(ctx context.Context, req domain.TransferRequest) error {
	return s.baseTransfer(ctx, req)
}

func (s LedgerService) InterBankTransfer(ctx context.Context, req domain.TransferRequest) error {
	return s.baseTransfer(ctx, req)
}

func (s LedgerService) baseTransfer(ctx context.Context, req domain.TransferRequest) error {
	description := ""
	if rawDescription, ok := req.Meta["Description"]; ok {
		description, _ = rawDescription.(string)
	}

	/*
		Sort account ids in order to prevent deadlocks - Lock ordering
	*/
	firstAccount, secondAccount := utils.SortAccount(req.DestinationAccountID.String(), req.FromAccountID.String())

	return s.txManager.WithTansaction(ctx, func(repo Repository) error {

		accounts := make(map[string]*domain.Account)

		for _, accountID := range []string{firstAccount, secondAccount} {
			account, err := repo.GetAccountForUpdate(ctx, accountID)

			if err != nil {
				s.logger.Error().Err(err).Msgf("failed to fetch account with id %s for update", accountID)
				return err
			}
			accounts[accountID] = account
		}

		sender := accounts[req.FromAccountID.String()]

		if sender.AvailableBalanceMinorUnits < req.Money.AmountMinorUnits {
			s.logger.Error().Msgf("sender %s has insufficient balance", sender.ID)
			return errors.New("insufficient funds")
		}

		if err := repo.UpdateBalance(ctx, req.FromAccountID.String(), (-1 * req.Money.AmountMinorUnits)); err != nil {
			s.logger.Error().Err(err).Msgf("failed to update account %s", req.FromAccountID)
			return err
		}

		if err := repo.UpdateBalance(ctx, req.DestinationAccountID.String(), req.Money.AmountMinorUnits); err != nil {
			s.logger.Error().Err(err).Msgf("failed to update account %s", req.DestinationAccountID)
			return err
		}

		tx := domain.CreateTransactionRequest{
			FromAccountID:        req.FromAccountID,
			DestinationAccountID: req.DestinationAccountID,
			SessionID:            req.IdempotencyKey,
			Currency:             req.Money.Currency,
			Description:          description,
			Status:               string(domain.TRANSACTION_PENDING),
			AmountMinorUnits:     req.Money.AmountMinorUnits,
		}

		newTx, err := repo.CreateTransaction(ctx, tx)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to create transaction record")
			return err
		}
		outboxEvent := domain.CreateOutboxEventRequest{
			EventType: domain.EventMoneyTransfer,
			Payload: domain.TransactionEventPayload{
				TransactionID:        newTx.ID,
				FromAccountID:        req.FromAccountID,
				DestinationAccountID: req.DestinationAccountID,
				AmountMinorUnits:     req.Money.AmountMinorUnits,
				Currency:             req.Money.Currency,
			},
			Status:         domain.OutboxEventPending,
			IdempotencyKey: req.IdempotencyKey,
			Priority:       domain.OutboxPriorityHigh,
			//TODO:Convert this to a variable in config
			Producer: "ledger service",
		}
		if err := repo.CreateOutBoxEvent(ctx, outboxEvent); err != nil {
			s.logger.Error().Err(err).Msg("failed to create event")
			return err
		}
		entries := make([]domain.LedgerEntry, 2)

		//represents sender
		entries[0] = domain.LedgerEntry{
			TransactionID:    newTx.ID,
			AccountID:        req.FromAccountID,
			EntryType:        string(domain.TRANSACTION_DEBIT),
			AmountMinorUnits: req.Money.AmountMinorUnits,
			Currency:         req.Money.Currency,
			IdempotencyKey:   req.IdempotencyKey,
			Status:           string(domain.TRANSACTION_PENDING),
		}

		//represents recipient
		entries[1] = domain.LedgerEntry{
			TransactionID:    newTx.ID,
			AccountID:        req.DestinationAccountID,
			EntryType:        string(domain.TRANSACTION_CREDIT),
			AmountMinorUnits: req.Money.AmountMinorUnits,
			Currency:         req.Money.Currency,
			IdempotencyKey:   req.IdempotencyKey,
			Status:           string(domain.TRANSACTION_PENDING),
		}

		return repo.CreateLedgerEntry(ctx, entries)
	})
}

func (s LedgerService) isInterBankTransfer(req domain.TransferRequest) bool {
	rawTransferType, ok := req.Meta["transfer_type"]
	if !ok {
		return false
	}

	transferType, ok := rawTransferType.(string)
	if !ok {
		return false
	}

	normalized := strings.TrimSpace(strings.ToLower(transferType))
	return normalized == "interbank" || normalized == "inter-bank" || normalized == "inter_bank"
}
