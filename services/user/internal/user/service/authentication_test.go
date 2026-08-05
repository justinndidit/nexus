package service_test

import (
	"context"
	"testing"

	"github.com/justinndidit/nexus/user/internal/user/domain"
	"github.com/justinndidit/nexus/user/internal/user/service"
	"github.com/justinndidit/nexus/user/tests/mocks"
	"github.com/rs/zerolog"
)

func TestAuthService_Register(t *testing.T) {
	userData := &domain.RegisterUserDTO{
		Email:       "testing@test.com",
		Password:    "password",
		PhoneNumber: "08152281909",
		Username:    "testMonkey",
	}
	t.Run("successful signup", func(t *testing.T) {

		mockRepo := mocks.NewMockUserRepo()
		publisher := mocks.NewMockPublisher()
		mockPasswordHasher := mocks.NewMockBcrypt()
		jwtService := mocks.NewMockJwtService()
		logger := &zerolog.Logger{}
		authService := service.NewAuthService(logger, mockRepo, publisher, mockPasswordHasher, jwtService)

		newUserData, err := authService.Register(context.Background(), userData)

		if err != nil {
			t.Fatalf("expected no errors, got error: %v", err)
		}

		if mockRepo.UserCreated != true {
			t.Fatal("service must create user")
		}

		if !mockRepo.FetchUser {
			t.Fatalf("service must check if user: %s exists", userData.Email)
		}
		if newUserData == nil {
			t.Fatal("service must return new user data")
		}

		if !(mockPasswordHasher.Hashed) {
			t.Fatal("password was not not hashed before storing")
		}
		if newUserData.Password != "" {
			t.Fatal("password must never be returned")
		}

		if newUserData.ID.String() == "" {
			t.Fatal("user id must be generated")
		}
		if !publisher.Published {
			t.Fatal("sign up notification not published")
		}

	})

	// t.Run("user alreaddy exists", func(t *testing.T) {

	// })

}

// func UserAdminSignup(t *testing.T) {

// }

// func UserVerifyOTP(t *testing.T) {

// }

// func UserLoginTest(t *testing.T) {

// }

// func UserRefreshSession(t *testing.T) {

// }

// func UserResetPassword(t *testing.T) {

// }
