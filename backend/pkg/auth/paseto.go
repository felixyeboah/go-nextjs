package auth

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto/v2"
)

type TokenManager struct {
	symmetricKey []byte
}

type TokenClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewTokenManager(key string) (*TokenManager, error) {
	if len(key) < 32 {
		return nil, fmt.Errorf("symmetric key must be at least 32 bytes long")
	}

	return &TokenManager{
		symmetricKey: []byte(key),
	}, nil
}

func (tm *TokenManager) CreateToken(claims TokenClaims) (string, error) {
	jsonToken := paseto.JSONToken{
		Audience:   "",
		Issuer:     "",
		Jti:        "",
		Subject:    claims.UserID,
		IssuedAt:   claims.IssuedAt,
		Expiration: claims.ExpiresAt,
		NotBefore:  claims.IssuedAt,
	}

	// Add custom claims
	jsonToken.Set("email", claims.Email)
	jsonToken.Set("role", claims.Role)

	// Encrypt the token
	return paseto.NewV2().Encrypt(tm.symmetricKey, jsonToken, nil)
}

func (tm *TokenManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	var jsonToken paseto.JSONToken
	var footer string

	err := paseto.NewV2().Decrypt(tokenString, tm.symmetricKey, &jsonToken, &footer)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if jsonToken.Expiration.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	// Get custom claims
	var email string
	if err := jsonToken.Get("email", &email); err != nil && err != paseto.ErrClaimNotFound {
		return nil, fmt.Errorf("failed to get email claim: %w", err)
	}

	var role string
	if err := jsonToken.Get("role", &role); err != nil && err != paseto.ErrClaimNotFound {
		return nil, fmt.Errorf("failed to get role claim: %w", err)
	}

	claims := &TokenClaims{
		UserID:    jsonToken.Subject,
		Email:     email,
		Role:      role,
		IssuedAt:  jsonToken.IssuedAt,
		ExpiresAt: jsonToken.Expiration,
	}

	return claims, nil
}

func (tm *TokenManager) CreateRefreshToken(userID string, duration time.Duration) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return tm.CreateToken(claims)
}
