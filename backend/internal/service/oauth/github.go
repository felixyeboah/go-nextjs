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
	"golang.org/x/oauth2/github"
)

type GitHubService struct {
	config *oauth2.Config
}

func NewGitHubService(cfg *config.Config) *GitHubService {
	return &GitHubService{
		config: &oauth2.Config{
			ClientID:     cfg.OAuth.GitHub.ClientID,
			ClientSecret: cfg.OAuth.GitHub.ClientSecret,
			RedirectURL:  cfg.OAuth.GitHub.RedirectURL,
			Scopes: []string{
				"user:email",
				"read:user",
			},
			Endpoint: github.Endpoint,
		},
	}
}

func (s *GitHubService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *GitHubService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

func (s *GitHubService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*models.OAuthUserInfo, error) {
	client := s.config.Client(ctx, token)

	// Get user profile
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", body)
	}

	var userInfo struct {
		ID        int64  `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// If email is not public, get primary email
	if userInfo.Email == "" {
		email, err := s.getPrimaryEmail(ctx, client)
		if err != nil {
			return nil, err
		}
		userInfo.Email = email
	}

	return &models.OAuthUserInfo{
		Provider:      "github",
		ProviderID:    fmt.Sprintf("%d", userInfo.ID),
		Email:         userInfo.Email,
		EmailVerified: true, // GitHub verifies emails
		Name:          userInfo.Name,
		Picture:       userInfo.AvatarURL,
	}, nil
}

func (s *GitHubService) getPrimaryEmail(ctx context.Context, client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", fmt.Errorf("failed to get user emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get user emails: %s", body)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode user emails: %w", err)
	}

	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no primary verified email found")
}
