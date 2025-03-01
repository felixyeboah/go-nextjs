package oauth

import (
	"context"

	"github.com/nanayaw/fullstack/internal/models"
	"golang.org/x/oauth2"
)

// Service defines the interface for OAuth providers
type Service interface {
	// GetAuthURL returns the authorization URL for the OAuth provider
	GetAuthURL(state string) string

	// Exchange exchanges the authorization code for an access token
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)

	// GetUserInfo retrieves user information from the OAuth provider
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*models.OAuthUserInfo, error)
}
