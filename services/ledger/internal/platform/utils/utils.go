package utils

import (
	"strings"
)

func SortAccount(recipient, sender string) (string, string) {
	if strings.Compare(recipient, sender) < 0 {
		return recipient, sender
	}

	return sender, recipient
}
