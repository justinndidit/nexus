package utils

import (
	"strings"

	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/shopspring/decimal"
)

// TODO: Implement function
func MoneyIntToDeimal(domain.Money) decimal.Decimal {
	return decimal.Decimal{}
}

func SortAccount(recipient, sender string) (string, string) {
	if strings.Compare(recipient, sender) < 0 {
		return recipient, sender
	}

	return sender, recipient
}
