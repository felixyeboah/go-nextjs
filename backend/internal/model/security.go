package model

import (
	"time"
)

// LoginAttempt represents a user login attempt
type LoginAttempt struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	UserAgent   string    `json:"user_agent" db:"user_agent"`
	Location    string    `json:"location" db:"location"`
	Successful  bool      `json:"successful" db:"successful"`
	AttemptedAt time.Time `json:"attempted_at" db:"attempted_at"`
}

// SecurityEvent represents a security-related event for a user
type SecurityEvent struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	EventType   string    `json:"event_type" db:"event_type"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	UserAgent   string    `json:"user_agent" db:"user_agent"`
	Location    string    `json:"location" db:"location"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// AccountLock represents a user account lock
type AccountLock struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	LockedAt  time.Time `json:"locked_at" db:"locked_at"`
	UnlockAt  time.Time `json:"unlock_at" db:"unlock_at"`
	Reason    string    `json:"reason" db:"reason"`
	CreatedBy string    `json:"created_by" db:"created_by"`
}

// SecurityEventTypes defines constants for different types of security events
const (
	// Login-related events
	EventLoginSuccess     = "login_success"
	EventLoginFailed      = "login_failed"
	EventNewDeviceLogin   = "new_device_login"
	EventNewLocationLogin = "new_location_login"

	// Account-related events
	EventAccountCreated  = "account_created"
	EventAccountLocked   = "account_locked"
	EventAccountUnlocked = "account_unlocked"
	EventAccountDisabled = "account_disabled"
	EventAccountEnabled  = "account_enabled"

	// Password-related events
	EventPasswordChanged        = "password_changed"
	EventPasswordReset          = "password_reset"
	EventPasswordResetRequested = "password_reset_requested"

	// Email-related events
	EventEmailChanged  = "email_changed"
	EventEmailVerified = "email_verified"

	// Suspicious activity
	EventSuspiciousActivity = "suspicious_activity"

	// Admin actions
	EventAdminAction = "admin_action"
)
