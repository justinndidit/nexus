package ledger

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/justinndidit/nexus/ledger/internal/ledger/repository"
	"github.com/rs/zerolog"
)

type LedgerService struct {
	validator validator.Validate
	logger    *zerolog.Logger
	txManager repository.TransactionManager
	stores    *repository.PostgresStores
}

func NewLegerService(tm repository.TransactionManager, stores *repository.PostgresStores, v validator.Validate, log *zerolog.Logger) *LedgerService {
	return &LedgerService{
		validator: v,
		logger:    log,
		txManager: tm,
		stores:    stores,
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
	if req.DestinationAccountID == req.FromAccountID {
		s.logger.Error().Str("func", "IntranBankTransfer").Msg("self transfer is not allowed")
		return fmt.Errorf("self transfer is not allowed")
	}
	description := ""
	if rawDescription, ok := req.Meta["Description"]; ok {
		description, _ = rawDescription.(string)
	}

	/*
		Lock Ordering: To prevent deadlocks
		Basically we sort the account numbers in ascending order, so regardless of transaction direction,
		we lock the pair of accounts in the same order, giving a deterministic behavior.
	*/
	firstAccount, secondAccount := sortAccount(req.DestinationAccountID.String(), req.FromAccountID.String())
	return s.txManager.WithTransaction(ctx, func(stores repository.PostgresStores) error {
		accounts := make(map[string]*domain.Account)

		for _, accountID := range []string{firstAccount, secondAccount} {
			account, err := stores.AccountStore.GetAccountForUpdate(ctx, accountID)

			if err != nil {
				s.logger.Error().Err(err).Msgf("failed to fetch account with id %s for update", accountID)
				return err
			}
			accounts[accountID] = account
		}

		txCurrency := string(req.Money.CurrencyCode)
		sender := accounts[req.FromAccountID.String()]
		receiver := accounts[req.DestinationAccountID.String()]

		existingTx, err := stores.TransactionStore.GetTransactionBySessionID(ctx, req.IdempotencyKey)
		if err != nil {
			s.logger.Error().Err(err).Str("func", "IntraBankTransfer").Msgf("failed to fetch transaction with session id: %s", req.IdempotencyKey)
			return err
		}

		if existingTx != nil {
			s.logger.Error().Str("func", "IntraBankTransfer").Msgf("transfer with session id %s already exists", req.IdempotencyKey)
			return fmt.Errorf("duplicate transaction")
		}

		if sender.Currency != txCurrency || receiver.Currency != txCurrency {
			s.logger.Error().Str("func", "IntrabankTransfer").Msg("invalid currency transaction")
			return fmt.Errorf("invalid currency transaction")
		}

		if sender.AvailableBalanceMinorUnits < req.Money.AmountMinorUnits {
			s.logger.Error().Str("func", "IntraBankTransfer").Msgf("sender %s has insufficient balance", sender.ID)
			return errors.New("insufficient funds")
		}

		if err := stores.AccountStore.UpdateBalance(ctx, req.FromAccountID.String(), (-1 * req.Money.AmountMinorUnits)); err != nil {
			s.logger.Error().Err(err).Msgf("failed to update account %s", req.FromAccountID)
			return err
		}

		if err := stores.AccountStore.UpdateBalance(ctx, req.DestinationAccountID.String(), req.Money.AmountMinorUnits); err != nil {
			s.logger.Error().Err(err).Msgf("failed to update account %s", req.DestinationAccountID)
			return err
		}

		tx := domain.CreateTransactionRequest{
			FromAccountID:        req.FromAccountID,
			DestinationAccountID: req.DestinationAccountID,
			IdempotencyKey:       req.IdempotencyKey,
			Currency:             string(req.Money.CurrencyCode),
			Description:          description,
			Status:               string(domain.TRANSACTION_COMPLETED),
			AmountMinorUnits:     req.Money.AmountMinorUnits,
		}

		newTx, err := stores.TransactionStore.CreateTransaction(ctx, tx)
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
				Currency:             string(req.Money.CurrencyCode),
			},
			Status:         domain.OutboxEventPending,
			IdempotencyKey: req.IdempotencyKey,
			Priority:       domain.OutboxPriorityHigh,
			//TODO:Convert this to a variable in config
			Producer: "ledger service",
		}
		if err := stores.OutboxStore.CreateOutBoxEvent(ctx, outboxEvent); err != nil {
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
			Currency:         string(req.Money.CurrencyCode),
		}

		//represents recipient
		entries[1] = domain.LedgerEntry{
			TransactionID:    newTx.ID,
			AccountID:        req.DestinationAccountID,
			EntryType:        string(domain.TRANSACTION_CREDIT),
			AmountMinorUnits: req.Money.AmountMinorUnits,
			Currency:         string(req.Money.CurrencyCode),
		}

		return stores.LedgerEntryStore.CreateLedgerEntry(ctx, entries)
	})
}

func (s LedgerService) InterBankTransfer(ctx context.Context, req domain.TransferRequest) error {
	if req.DestinationAccountID == req.FromAccountID {
		s.logger.Error().Str("func", "InterBankTransfer").Msg("self transfer is not allowed")
		return fmt.Errorf("self transfer is not allowed")
	}
	return nil
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

func sortAccount(recipient, sender string) (string, string) {
	if strings.Compare(recipient, sender) < 0 {
		return recipient, sender
	}

	return sender, recipient
}

func calcTransferCharges(amount int64) int {
	return 10
}
