package ledger

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/rs/zerolog"
)

type fakeTransactionManager struct {
	called bool
	repo   Repository
	err    error
}

func (f *fakeTransactionManager) WithTansaction(ctx context.Context, fn func(repo Repository) error) error {
	f.called = true
	if f.err != nil {
		return f.err
	}
	return fn(f.repo)
}

type fakeRepository struct {
	accounts map[string]*domain.Account

	getAccountCalls []string
	updateCalls     []balanceUpdate
	transactions    []domain.CreateTransactionRequest
	outboxEvents    []domain.CreateOutboxEventRequest
	ledgerEntries   [][]domain.LedgerEntry

	createTransactionResult *domain.Transaction

	getAccountErr        error
	updateBalanceErr     error
	createTransactionErr error
	createOutboxErr      error
	createLedgerErr      error
}

type balanceUpdate struct {
	accountID string
	amount    int64
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		accounts: make(map[string]*domain.Account),
	}
}

func (f *fakeRepository) CreateLedgerEntry(ctx context.Context, entries []domain.LedgerEntry) error {
	if f.createLedgerErr != nil {
		return f.createLedgerErr
	}

	cloned := make([]domain.LedgerEntry, len(entries))
	copy(cloned, entries)
	f.ledgerEntries = append(f.ledgerEntries, cloned)
	return nil
}

func (f *fakeRepository) CreateTransaction(ctx context.Context, req domain.CreateTransactionRequest) (*domain.Transaction, error) {
	if f.createTransactionErr != nil {
		return nil, f.createTransactionErr
	}

	f.transactions = append(f.transactions, req)
	if f.createTransactionResult != nil {
		return f.createTransactionResult, nil
	}

	tx := &domain.Transaction{ID: uuid.New()}
	f.createTransactionResult = tx
	return tx, nil
}

func (f *fakeRepository) GetAccountForUpdate(ctx context.Context, accountID string) (*domain.Account, error) {
	f.getAccountCalls = append(f.getAccountCalls, accountID)
	if f.getAccountErr != nil {
		return nil, f.getAccountErr
	}

	account, ok := f.accounts[accountID]
	if !ok {
		return nil, context.Canceled
	}

	return account, nil
}

func (f *fakeRepository) UpdateBalance(ctx context.Context, accountID string, amount int64) error {
	if f.updateBalanceErr != nil {
		return f.updateBalanceErr
	}

	f.updateCalls = append(f.updateCalls, balanceUpdate{accountID: accountID, amount: amount})
	return nil
}

func (f *fakeRepository) CreateOutBoxEvent(ctx context.Context, req domain.CreateOutboxEventRequest) error {
	if f.createOutboxErr != nil {
		return f.createOutboxErr
	}

	f.outboxEvents = append(f.outboxEvents, req)
	return nil
}

func (f *fakeRepository) GetOutBoxEventsForUpdate(ctx context.Context) ([]domain.OutBoxEvent, error) {
	return nil, nil
}

func (f *fakeRepository) IncrementRetryCount(ctx context.Context, eventID string, errorString string) error {
	return nil
}

func (f *fakeRepository) MarkEventProcessed(ctx context.Context, eventID string) error {
	return nil
}

func testLogger() *zerolog.Logger {
	logger := zerolog.New(io.Discard)
	return &logger
}

func testValidator() validator.Validate {
	return *validator.New()
}

func newTestService(repo *fakeRepository, txManager *fakeTransactionManager) *LedgerService {
	return NewLegerService(repo, txManager, testValidator(), testLogger())
}

func validTransferRequest() domain.TransferRequest {
	return domain.TransferRequest{
		FromAccountID:        uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		DestinationAccountID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		IdempotencyKey:       "idem-123",
		AmountMinorUnit:      1500,
		Meta: domain.TransferMetaData{
			"Description":   "invoice payment",
			"transfer_type": "interbank",
		},
	}
}

func TestLedgerService_Transfer_ValidationFails(t *testing.T) {
	repo := newFakeRepository()
	txManager := &fakeTransactionManager{repo: repo}
	service := newTestService(repo, txManager)

	_, err := service.Transfer(context.Background(), domain.TransferRequest{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if txManager.called {
		t.Fatal("transaction manager should not be called when validation fails")
	}
}

func TestLedgerService_Transfer_Success(t *testing.T) {
	request := validTransferRequest()
	repo := newFakeRepository()
	repo.accounts[request.FromAccountID.String()] = &domain.Account{
		ID:                         request.FromAccountID,
		AvailableBalanceMinorUnits: 5000,
	}
	repo.accounts[request.DestinationAccountID.String()] = &domain.Account{
		ID:                         request.DestinationAccountID,
		AvailableBalanceMinorUnits: 200,
	}

	txManager := &fakeTransactionManager{repo: repo}
	service := newTestService(repo, txManager)

	_, err := service.Transfer(context.Background(), request)
	if err != nil {
		t.Fatalf("expected transfer to succeed, got error: %v", err)
	}

	if !txManager.called {
		t.Fatal("expected transaction manager to be called")
	}

	if got := repo.getAccountCalls; len(got) != 2 {
		t.Fatalf("expected 2 locked account reads, got %d", len(got))
	} else {
		if got[0] != request.DestinationAccountID.String() || got[1] != request.FromAccountID.String() {
			t.Fatalf("expected lock ordering by account id, got %v", got)
		}
	}

	if len(repo.updateCalls) != 2 {
		t.Fatalf("expected 2 balance updates, got %d", len(repo.updateCalls))
	}

	if repo.updateCalls[0].accountID != request.FromAccountID.String() || repo.updateCalls[0].amount != -1500 {
		t.Fatalf("unexpected debit update: %#v", repo.updateCalls[0])
	}

	if repo.updateCalls[1].accountID != request.DestinationAccountID.String() || repo.updateCalls[1].amount != 1500 {
		t.Fatalf("unexpected credit update: %#v", repo.updateCalls[1])
	}

	if len(repo.transactions) != 1 {
		t.Fatalf("expected 1 transaction record, got %d", len(repo.transactions))
	}

	createdTx := repo.transactions[0]
	if createdTx.AmountMinorUnits != request.AmountMinorUnit {
		t.Fatalf("expected transaction amount %d, got %d", request.AmountMinorUnit, createdTx.AmountMinorUnits)
	}
	if createdTx.Description != "invoice payment" {
		t.Fatalf("expected transaction description to be copied from metadata, got %q", createdTx.Description)
	}

	if len(repo.outboxEvents) != 1 {
		t.Fatalf("expected 1 outbox event, got %d", len(repo.outboxEvents))
	}
	payload, ok := repo.outboxEvents[0].Payload.(domain.TransactionEventPayload)
	if !ok {
		t.Fatalf("expected outbox payload to be TransactionEventPayload, got %T", repo.outboxEvents[0].Payload)
	}
	if payload.AmountMinorUnits != request.AmountMinorUnit {
		t.Fatalf("expected outbox amount %d, got %d", request.AmountMinorUnit, payload.AmountMinorUnits)
	}

	if len(repo.ledgerEntries) != 1 {
		t.Fatalf("expected 1 ledger entry batch, got %d", len(repo.ledgerEntries))
	}
	entries := repo.ledgerEntries[0]
	if len(entries) != 2 {
		t.Fatalf("expected 2 ledger entries, got %d", len(entries))
	}
	if entries[0].EntryType != string(domain.TRANSACTION_DEBIT) || entries[1].EntryType != string(domain.TRANSACTION_CREDIT) {
		t.Fatalf("unexpected ledger entry types: %#v", entries)
	}
	for _, entry := range entries {
		if entry.AmountMinorUnits != request.AmountMinorUnit {
			t.Fatalf("expected ledger entry amount %d, got %d", request.AmountMinorUnit, entry.AmountMinorUnits)
		}
	}
}

func TestLedgerService_Transfer_InsufficientFunds(t *testing.T) {
	request := validTransferRequest()
	repo := newFakeRepository()
	repo.accounts[request.FromAccountID.String()] = &domain.Account{
		ID:                         request.FromAccountID,
		AvailableBalanceMinorUnits: 100,
	}
	repo.accounts[request.DestinationAccountID.String()] = &domain.Account{
		ID:                         request.DestinationAccountID,
		AvailableBalanceMinorUnits: 200,
	}

	txManager := &fakeTransactionManager{repo: repo}
	service := newTestService(repo, txManager)

	_, err := service.Transfer(context.Background(), request)
	if err == nil {
		t.Fatal("expected insufficient funds error, got nil")
	}

	if !strings.Contains(err.Error(), "insufficient funds") {
		t.Fatalf("expected insufficient funds error, got %v", err)
	}

	if len(repo.updateCalls) != 0 {
		t.Fatalf("expected no balance updates, got %d", len(repo.updateCalls))
	}
	if len(repo.transactions) != 0 || len(repo.outboxEvents) != 0 || len(repo.ledgerEntries) != 0 {
		t.Fatal("expected no downstream writes when funds are insufficient")
	}
}

func TestLedgerService_isInterBankTransfer(t *testing.T) {
	service := newTestService(newFakeRepository(), &fakeTransactionManager{})

	tests := []struct {
		name string
		meta domain.TransferMetaData
		want bool
	}{
		{name: "missing meta", meta: nil, want: false},
		{name: "intra bank spelling", meta: domain.TransferMetaData{"transfer_type": "intra-bank"}, want: false},
		{name: "interbank", meta: domain.TransferMetaData{"transfer_type": "interbank"}, want: true},
		{name: "inter bank dashed", meta: domain.TransferMetaData{"transfer_type": "inter-bank"}, want: true},
		{name: "inter bank underscored", meta: domain.TransferMetaData{"transfer_type": "inter_bank"}, want: true},
		{name: "whitespace and case", meta: domain.TransferMetaData{"transfer_type": "  InterBank  "}, want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := service.isInterBankTransfer(domain.TransferRequest{Meta: tc.meta})
			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
