package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID              uuid.UUID  `json:"id" validate:"required" db:"id"`
	Email           string     `json:"email" validate:"required" db:"email"`
	PhoneNumber     *string    `json:"phone_number" validate:"required" db:"phone_number"`
	Username        string     `json:"username" validate:"required" db:"username"`
	Password        string     `json:"-" db:"password"`
	KYCStatus       string     `json:"kyc_status" db:"kyc_status"`
	UserType        string     `json:"user_type" db:"user_type"`
	IsEmailVerified bool       `json:"is_email_verified" db:"is_email_verified"`
	IsPhoneVerified bool       `json:"is_phone_verified" db:"is_phone_verified"`
	LockedUntil     *time.Time `json:"locked_until" db:"locked_until"`
	LastLoginAt     *time.Time `json:"last_login_at" db:"last_login_at"`
	LastLoginIP     *string    `json:"last_login_ip" db:"last_login_ip"` // INET maps well to string or net.IP
}

type RegisterUserDTO struct {
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Username    string `json:"username" validate:"required"`
}

type HttpResponse struct {
	Status     string                 `json:"status"`
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"string"`
	Data       map[string]interface{} `json:"data"`
	Errors     map[string]interface{} `json:"errors"`
	Action     map[string]interface{} `json:"action omitempty"`
}
