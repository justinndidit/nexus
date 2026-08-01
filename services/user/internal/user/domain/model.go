package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	PhoneNumber  *string   `json:"phone_number" db:"phone_number"` // Pointer for nullable field

	// KYC & Compliance
	KYCStatus     string     `json:"kyc_status" db:"kyc_status"`
	KYCVerifiedAt *time.Time `json:"kyc_verified_at" db:"kyc_verified_at"`
	KYCVerifiedBy *uuid.UUID `json:"kyc_verified_by" db:"kyc_verified_by"`
	KYCExpiresAt  *time.Time `json:"kyc_expires_at" db:"kyc_expires_at"`

	// User Classification
	UserType string `json:"user_type" db:"user_type"`

	// Verification Status
	IsEmailVerified bool       `json:"is_email_verified" db:"is_email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" db:"email_verified_at"`
	IsPhoneVerified bool       `json:"is_phone_verified" db:"is_phone_verified"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at" db:"phone_verified_at"`

	// Security Flags
	IsBlocked     bool       `json:"is_blocked" db:"is_blocked"`
	IsSuspended   bool       `json:"is_suspended" db:"is_suspended"`
	BlockedReason *string    `json:"blocked_reason" db:"blocked_reason"`
	BlockedAt     *time.Time `json:"blocked_at" db:"blocked_at"`
	BlockedBy     *uuid.UUID `json:"blocked_by" db:"blocked_by"`

	// Password Security
	PasswordChangedAt   time.Time  `json:"password_changed_at" db:"password_changed_at"`
	PasswordExpiresAt   *time.Time `json:"password_expires_at" db:"password_expires_at"`
	FailedLoginAttempts int        `json:"failed_login_attempts" db:"failed_login_attempts"`
	LastFailedLoginAt   *time.Time `json:"last_failed_login_at" db:"last_failed_login_at"`
	LockedUntil         *time.Time `json:"locked_until" db:"locked_until"`

	// MFA
	MFAEnabled bool     `json:"mfa_enabled" db:"mfa_enabled"`
	MFAMethods []string `json:"mfa_methods" db:"mfa_methods"` // Assumes driver supports string slice for custom enum arrays

	// Compliance & Legal
	TermsAcceptedAt    *time.Time      `json:"terms_accepted_at" db:"terms_accepted_at"`
	PrivacyAcceptedAt  *time.Time      `json:"privacy_accepted_at" db:"privacy_accepted_at"`
	DataRetentionUntil *time.Time      `json:"data_retention_until" db:"data_retention_until"`
	GDPRConsent        json.RawMessage `json:"gdpr_consent" db:"gdpr_consent"`

	// Metadata
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
	LastLoginIP  *string    `json:"last_login_ip" db:"last_login_ip"` // INET maps well to string or net.IP
	Source       *string    `json:"source" db:"source"`
	ReferralCode *string    `json:"referral_code" db:"referral_code"`
	ReferredBy   *uuid.UUID `json:"referred_by" db:"referred_by"`

	// Soft Delete
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at"`
	DeletedBy      *uuid.UUID `json:"deleted_by" db:"deleted_by"`
	DeletionReason *string    `json:"deletion_reason" db:"deletion_reason"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
