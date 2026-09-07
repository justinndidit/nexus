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

const (
	ServiceName = "ledger service"
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

func (s LedgerService) Transfer(ctx context.Context, req domain.TransferRequest) (*domain.TransferResponse, error) {
	err := s.validator.Struct(req)
	if err != nil {
		s.logger.Error().Err(err).Str("func", "Transfer").Msg("failed to validate transfer request data")
		return nil, err
	}

	return s.baseTransfer(ctx, req, s.isInterBankTransfer(req))
}

func (s LedgerService) IntraBankTransfer(ctx context.Context, req domain.TransferRequest) (*domain.TransferResponse, error) {
	return s.baseTransfer(ctx, req, true)
}

func (s LedgerService) InterBankTransfer(ctx context.Context, req domain.TransferRequest) (*domain.TransferResponse, error) {
	return s.baseTransfer(ctx, req, false)
}

func (s LedgerService) baseTransfer(ctx context.Context, req domain.TransferRequest, isIntraBank bool) (*domain.TransferResponse, error) {
	description := ""
	if rawDescription, ok := req.Meta["Description"]; ok {
		description, _ = rawDescription.(string)
	}

	/*
		Sort account ids in order to prevent deadlocks - Lock ordering
		This guarantees locks are synchronized and always in one direction.
	*/
	firstAccount, secondAccount := utils.SortAccount(req.DestinationAccountID.String(), req.FromAccountID.String())
	var response domain.TransferResponse
	err := s.txManager.WithTansaction(ctx, func(repo Repository) error {

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

		if sender.AvailableBalanceMinorUnits < req.AmountMinorUnit {
			s.logger.Error().Msgf("sender %s has insufficient balance", sender.ID)
			return errors.New("insufficient funds")
		}

		if err := repo.UpdateBalance(ctx, req.FromAccountID.String(), (-1 * req.AmountMinorUnit)); err != nil {
			s.logger.Error().Err(err).Msgf("failed to update account %s", req.FromAccountID)
			return err
		}

		tx := domain.CreateTransactionRequest{
			FromAccountID:        req.FromAccountID,
			DestinationAccountID: req.DestinationAccountID,
			SessionID:            req.IdempotencyKey,
			Currency:             req.CurrencyCode,
			Description:          description,
			Status:               string(domain.TRANSACTION_PENDING),
			AmountMinorUnits:     req.AmountMinorUnit,
		}

		if isIntraBank {
			if err := repo.UpdateBalance(ctx, req.DestinationAccountID.String(), req.AmountMinorUnit); err != nil {
				s.logger.Error().Err(err).Msgf("failed to update account %s", req.DestinationAccountID)
				return err
			}
			tx.Status = string(domain.TRANSACTION_COMPLETED)
		} else {
			//payment gate way
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
				AmountMinorUnits:     req.AmountMinorUnit,
				Currency:             req.CurrencyCode,
			},
			Status:         domain.OutboxEventPending,
			IdempotencyKey: req.IdempotencyKey,
			Priority:       domain.OutboxPriorityHigh,
			Producer:       ServiceName,
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
			AmountMinorUnits: req.AmountMinorUnit,
			Currency:         req.CurrencyCode,
			// IdempotencyKey:   req.IdempotencyKey,
			Status: string(domain.TRANSACTION_PENDING),
		}

		//represents recipient
		entries[1] = domain.LedgerEntry{
			TransactionID:    newTx.ID,
			AccountID:        req.DestinationAccountID,
			EntryType:        string(domain.TRANSACTION_CREDIT),
			AmountMinorUnits: req.AmountMinorUnit,
			Currency:         req.CurrencyCode,
			// IdempotencyKey:   req.IdempotencyKey,
			Status: string(domain.TRANSACTION_PENDING),
		}
		response = domain.TransferResponse{
			TransactionID: newTx.ID.String(),
			SessionID:     req.IdempotencyKey,
		}
		return repo.CreateLedgerEntry(ctx, entries)
	})

	return &response, err
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
