package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Money is an immutable-by-convention, always-non-negative monetary amount
// stored as minor units (e.g. cents), matching how amounts are persisted
// (AmountMinorUnits int64). Direction/sign lives on EntryType or
// TransactionType, not on Money — Money is always a magnitude.
type Money struct {
	AmountMinorUnits int64
	CurrencyCode     Currency
}

const (
	Naira Currency = "NGN"
	USD   Currency = "USD"
)

var fractionDigitsTable = map[Currency]int{
	Naira: 2,
	USD:   2,
}

func fractionDigits(currencyCode Currency) (int, error) {
	digits, ok := fractionDigitsTable[currencyCode]
	if !ok {
		return 0, fmt.Errorf("unknown fraction digits for currency %v", currencyCode)
	}
	return digits, nil
}

// NewMoney validates that the currency is known and the amount is non-negative.
func NewMoney(amountMinorUnits int64, currencyCode Currency) (Money, error) {
	if _, err := fractionDigits(currencyCode); err != nil {
		return Money{}, err
	}
	if amountMinorUnits < 0 {
		return Money{}, fmt.Errorf("amount must be non-negative; direction belongs on EntryType/TransactionType")
	}
	return Money{AmountMinorUnits: amountMinorUnits, CurrencyCode: currencyCode}, nil
}

func NewZero(currencyCode Currency) (Money, error) {
	return NewMoney(0, currencyCode)
}

// FromAmount builds a Money from a decimal amount, e.g. decimal.NewFromString("19.99").
// Boundary conversion only — decimal.Decimal never appears in the arithmetic path.
func FromAmount(amount decimal.Decimal, currencyCode Currency) (Money, error) {
	digits, err := fractionDigits(currencyCode)
	if err != nil {
		return Money{}, err
	}
	if amount.IsNegative() {
		return Money{}, fmt.Errorf("amount must be non-negative; direction belongs on EntryType/TransactionType")
	}

	unit := decimal.New(1, int32(digits))
	totalMinorUnits := amount.Mul(unit)
	if !totalMinorUnits.Equal(totalMinorUnits.Truncate(0)) {
		return Money{}, fmt.Errorf("amount %s has more precision than %v allows", amount, currencyCode)
	}

	return NewMoney(totalMinorUnits.IntPart(), currencyCode)
}

// ToAmount converts to a decimal amount for display/serialization, e.g. 19.99.
func (m Money) ToAmount() (decimal.Decimal, error) {
	digits, err := fractionDigits(m.CurrencyCode)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return decimal.New(m.AmountMinorUnits, -int32(digits)), nil
}

// Combine adds two Money values of the same currency. Pure integer op.
func (m Money) Combine(other Money) (Money, error) {
	if m.CurrencyCode != other.CurrencyCode {
		return Money{}, fmt.Errorf("currency mismatch: %v vs %v", m.CurrencyCode, other.CurrencyCode)
	}
	return NewMoney(m.AmountMinorUnits+other.AmountMinorUnits, m.CurrencyCode)
}

// Subtract subtracts other from m. Returns an error if it would go negative —
// Money can't represent a negative magnitude; a caller wanting a debit past
// zero balance is a business-rule decision, not something this type decides.
func (m Money) Subtract(other Money) (Money, error) {
	if m.CurrencyCode != other.CurrencyCode {
		return Money{}, fmt.Errorf("currency mismatch: %v vs %v", m.CurrencyCode, other.CurrencyCode)
	}
	result := m.AmountMinorUnits - other.AmountMinorUnits
	return NewMoney(result, m.CurrencyCode)
}

func (m Money) String() string {
	amount, err := m.ToAmount()
	if err != nil {
		return fmt.Sprintf("<invalid Money: %v>", err)
	}
	return fmt.Sprintf("%s %v", amount.String(), m.CurrencyCode)
}
