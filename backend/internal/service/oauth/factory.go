package oauth

import (
	"fmt"

	"github.com/nanayaw/fullstack/internal/config"
)

// Provider represents an OAuth provider type
type Provider string

const (
	// ProviderGoogle represents Google OAuth provider
	ProviderGoogle Provider = "google"
	// ProviderGitHub represents GitHub OAuth provider
	ProviderGitHub Provider = "github"
)

// NewOAuthService creates a new OAuth service based on the provider
func NewOAuthService(provider Provider, cfg *config.Config) (Service, error) {
	switch provider {
	case ProviderGoogle:
		return NewGoogleService(cfg), nil
	case ProviderGitHub:
		return NewGitHubService(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
}
