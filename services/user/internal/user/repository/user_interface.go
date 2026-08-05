package repository

import (
	"context"

	"github.com/justinndidit/nexus/user/internal/user/domain"
)

type UserRepository interface {
	CreateUser(context.Context, domain.RegisterUserDTO) (*domain.UserDTO, error)
	FetchUserByEmail(context.Context, string) (*domain.UserDTO, error)
}
