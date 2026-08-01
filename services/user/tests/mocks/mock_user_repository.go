package mocks

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/justinndidit/nexus/user/internal/user/domain"
)

type MockUserRepository struct {
	UserCreated bool
	FetchUser   bool
	EmailExists bool
	UserExists  bool
	FailCreate  bool
	FailFetch   bool
}

func NewMockUserRepo() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, userData domain.RegisterUserDTO) (*domain.UserDTO, error) {
	if m.FailCreate {
		return nil, errors.New("create failed")
	}

	if m.EmailExists {
		return nil, domain.ErrEmailAlreadyExists
	}

	m.UserCreated = true
	return &domain.UserDTO{
		ID:       uuid.New(),
		Email:    userData.Email,
		Username: userData.Username,
	}, nil
}

func (m *MockUserRepository) FetchUserByEmail(ctx context.Context, email string) (*domain.UserDTO, error) {
	m.FetchUser = true

	if m.FailFetch {
		return nil, errors.New("fetch failed")
	}

	if !m.UserExists {
		return nil, nil
	}

	return &domain.UserDTO{
		ID:       uuid.New(),
		Email:    email,
		Username: "existingUser",
	}, nil
}
