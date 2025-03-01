package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nanayaw/fullstack/internal/config"
	"github.com/nanayaw/fullstack/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleService struct {
	config *oauth2.Config
}

func NewGoogleService(cfg *config.Config) *GoogleService {
	return &GoogleService{
		config: &oauth2.Config{
			ClientID:     cfg.OAuth.Google.ClientID,
			ClientSecret: cfg.OAuth.Google.ClientSecret,
			RedirectURL:  cfg.OAuth.Google.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (s *GoogleService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *GoogleService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

func (s *GoogleService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*models.OAuthUserInfo, error) {
	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", body)
	}

	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &models.OAuthUserInfo{
		Provider:      "google",
		ProviderID:    userInfo.ID,
		Email:         userInfo.Email,
		EmailVerified: userInfo.VerifiedEmail,
		Name:          userInfo.Name,
		Picture:       userInfo.Picture,
	}, nil
}
