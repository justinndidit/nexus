package service

import (
	"context"

	"github.com/justinndidit/nexus/user/internal/platform/broker"
	"github.com/justinndidit/nexus/user/internal/user/domain"
	"github.com/justinndidit/nexus/user/internal/user/repository"
	"github.com/rs/zerolog"
)

type AuthService struct {
	logger         *zerolog.Logger
	userRepo       repository.UserRepository
	publisher      broker.Publisher
	passwordHasher PasswordHasher
	jwtService     JwtService
}

func NewAuthService(logger *zerolog.Logger, userRepo repository.UserRepository, publisher broker.Publisher, passwordHasher PasswordHasher, jwtService JwtService) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		publisher:      publisher,
		passwordHasher: passwordHasher,
		jwtService:     jwtService,
		logger:         logger,
	}
}

func (a *AuthService) Register(ctx context.Context, userData *domain.RegisterUserDTO) (*domain.UserDTO, error) {
	user, err := a.userRepo.FetchUserByEmail(ctx, userData.Email)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to fetch user from db")
		return nil, err
	}

	if !(user == nil) {
		a.logger.Error().Msgf("user with email %s already exists", userData.Email)
		return nil, domain.ErrEmailAlreadyExists
	}

	passwordHash, err := a.passwordHasher.Hash(userData.Password)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to hash user password")
		return nil, err
	}
	userData.Password = passwordHash

	newUser, err := a.userRepo.CreateUser(ctx, *userData)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to create user in db")
		return nil, err
	}

	if err = a.publisher.Publish("mwssage"); err != nil {
		a.logger.Error().Err(err).Msg("failed to publish user signed up event")
		return nil, err
	}

	return newUser, nil
}
