package models

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	FullName      string    `json:"fullName"`
	AvatarURL     string    `json:"avatarUrl"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	RefreshToken string    `json:"-"`
	UserAgent    string    `json:"userAgent"`
	ClientIP     string    `json:"clientIp"`
	IsBlocked    bool      `json:"isBlocked"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

type OAuthAccount struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Provider       string    `json:"provider"`
	ProviderUserID string    `json:"providerUserId"`
	AccessToken    string    `json:"-"`
	RefreshToken   string    `json:"-"`
	ExpiresAt      time.Time `json:"expiresAt"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type VerificationToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"-"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type AuditLog struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	Action     string    `json:"action"`
	EntityType string    `json:"entityType"`
	EntityID   string    `json:"entityId"`
	Metadata   string    `json:"metadata"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Request/Response models
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FullName  string `json:"fullName" validate:"required"`
	AvatarURL string `json:"avatarUrl"`
}

type UpdateUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=8"`
	FullName  *string `json:"fullName"`
	AvatarURL *string `json:"avatarUrl"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type OAuthLoginRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google github"`
	Code     string `json:"code" validate:"required"`
}

type PaginationParams struct {
	Page     int `query:"page" validate:"min=1"`
	PageSize int `query:"pageSize" validate:"min=1,max=100"`
}
