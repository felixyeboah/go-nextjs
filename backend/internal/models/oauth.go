package models

// OAuthUserInfo represents user information returned by OAuth providers
type OAuthUserInfo struct {
	Provider      string `json:"provider"`
	ProviderID    string `json:"provider_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}
