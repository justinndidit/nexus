package mocks

import "errors"

type MockBcrypt struct {
	Hashed      bool
	FailHashing bool
}

func NewMockBcrypt() *MockBcrypt {
	return &MockBcrypt{}
}

func (b *MockBcrypt) Hash(password string) (string, error) {
	if b.FailHashing {
		return "", errors.New("hashing failed")
	}
	b.Hashed = true
	return "hashed-password", nil
}

func (b *MockBcrypt) Compare(hash, password string) bool {
	return true
}
