package mocks

import "errors"

type MockJwtService struct {
	GeneratedToken        bool
	GeneratedRefreshToken bool

	FailGenerateToken        bool
	FailGenerateRefreshToken bool
}

func NewMockJwtService() *MockJwtService {
	return &MockJwtService{}
}

func (m *MockJwtService) Generate(payload any) (string, error) {
	if m.FailGenerateToken {
		return "", errors.New("failed to generate access token")
	}

	m.GeneratedToken = true
	return "mock-access-token", nil
}

func (m *MockJwtService) GenerateRefresh(payload any) (string, error) {
	if m.FailGenerateRefreshToken {
		return "", errors.New("failed to generate refresh token")
	}

	m.GeneratedRefreshToken = true
	return "mock-refresh-token", nil
}
