package config

// OAuthProviderConfig holds configuration for a single OAuth provider
type OAuthProviderConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

// GetRedirectURL returns the full redirect URL for a provider
func GetOAuthRedirectURL(baseURL, provider string) string {
	return baseURL + "/api/auth/" + provider + "/callback"
}
