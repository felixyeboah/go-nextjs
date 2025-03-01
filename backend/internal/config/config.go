package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Auth        AuthConfig
	Email       EmailConfig
	OAuth       OAuthConfig
	PASETO      PASETOConfig
	Security    SecurityConfig
	App         AppConfig
}

type ServerConfig struct {
	Host         string        `mapstructure:"SERVER_HOST"`
	Port         int           `mapstructure:"SERVER_PORT"`
	ReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	URL          string `mapstructure:"DATABASE_URL"`
	AuthToken    string `mapstructure:"DATABASE_AUTH_TOKEN"`
	MaxOpenConns int    `mapstructure:"DATABASE_MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"DATABASE_MAX_IDLE_CONNS"`
}

type RedisConfig struct {
	URL string `mapstructure:"REDIS_URL"`
}

type AuthConfig struct {
	PublicKey          string        `mapstructure:"AUTH_PUBLIC_KEY"`
	PrivateKey         string        `mapstructure:"AUTH_PRIVATE_KEY"`
	AccessTokenTTL     time.Duration `mapstructure:"AUTH_ACCESS_TOKEN_TTL"`
	RefreshTokenTTL    time.Duration `mapstructure:"AUTH_REFRESH_TOKEN_TTL"`
	VerificationTTL    time.Duration `mapstructure:"AUTH_VERIFICATION_TTL"`
	PasswordResetTTL   time.Duration `mapstructure:"AUTH_PASSWORD_RESET_TTL"`
	MaxLoginAttempts   int           `mapstructure:"AUTH_MAX_LOGIN_ATTEMPTS"`
	LockoutDuration    time.Duration `mapstructure:"AUTH_LOCKOUT_DURATION"`
	SessionMaxLifetime time.Duration `mapstructure:"AUTH_SESSION_MAX_LIFETIME"`
}

type EmailConfig struct {
	// Resend configuration
	ResendAPIKey      string `mapstructure:"RESEND_API_KEY"`
	FromEmail         string `mapstructure:"EMAIL_FROM_ADDRESS"`
	FromName          string `mapstructure:"EMAIL_FROM_NAME"`
	VerificationURL   string `mapstructure:"EMAIL_VERIFICATION_URL"`
	PasswordResetURL  string `mapstructure:"EMAIL_PASSWORD_RESET_URL"`
	LoginNotification bool   `mapstructure:"EMAIL_LOGIN_NOTIFICATION"`

	// Upstash Workflow configuration
	UpstashWorkflowURL   string `mapstructure:"UPSTASH_WORKFLOW_URL"`
	UpstashWorkflowToken string `mapstructure:"UPSTASH_WORKFLOW_TOKEN"`

	// TTL values (shared with AuthConfig)
	VerificationTTL  time.Duration `mapstructure:"AUTH_VERIFICATION_TTL"`
	PasswordResetTTL time.Duration `mapstructure:"AUTH_PASSWORD_RESET_TTL"`
}

type OAuthConfig struct {
	Google struct {
		ClientID     string `mapstructure:"OAUTH_GOOGLE_CLIENT_ID"`
		ClientSecret string `mapstructure:"OAUTH_GOOGLE_CLIENT_SECRET"`
		RedirectURL  string `mapstructure:"OAUTH_GOOGLE_REDIRECT_URL"`
	}
	GitHub struct {
		ClientID     string `mapstructure:"OAUTH_GITHUB_CLIENT_ID"`
		ClientSecret string `mapstructure:"OAUTH_GITHUB_CLIENT_SECRET"`
		RedirectURL  string `mapstructure:"OAUTH_GITHUB_REDIRECT_URL"`
	}
}

type PASETOConfig struct {
	PublicKeyPath  string `mapstructure:"PASETO_PUBLIC_KEY_PATH"`
	PrivateKeyPath string `mapstructure:"PASETO_PRIVATE_KEY_PATH"`
}

// SecurityConfig contains security-related configuration
type SecurityConfig struct {
	// Maximum number of failed login attempts before account is locked
	MaxLoginAttempts int `mapstructure:"max_login_attempts"`
	// Duration for which an account is locked after too many failed attempts
	AccountLockDuration time.Duration `mapstructure:"account_lock_duration"`
	// Whether to enable login notifications for new devices/locations
	EnableLoginNotifications bool `mapstructure:"enable_login_notifications"`
	// Whether to enable suspicious activity detection
	EnableSuspiciousActivityDetection bool `mapstructure:"enable_suspicious_activity_detection"`
	// Whether to enable rate limiting
	EnableRateLimiting bool `mapstructure:"enable_rate_limiting"`
	// Global rate limit (requests per minute)
	GlobalRateLimit int `mapstructure:"global_rate_limit"`
	// Auth rate limit (login attempts per 15 minutes)
	AuthRateLimit int `mapstructure:"auth_rate_limit"`
}

// AppConfig contains application-level configuration
type AppConfig struct {
	// Application name
	Name string `mapstructure:"name"`
	// Application URL
	URL string `mapstructure:"url"`
	// Environment (development, staging, production)
	Environment string `mapstructure:"environment"`
	// Debug mode
	Debug bool `mapstructure:"debug"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	fmt.Println("Loading config from path:", path)

	// First try to read from .env file directly
	envFile := filepath.Join(path, ".env")
	if _, err := os.Stat(envFile); err == nil {
		fmt.Println("Found .env file at:", envFile)
		viper.SetConfigFile(envFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading .env file: %w", err)
		}
		fmt.Println(".env file loaded successfully")
	} else {
		// Fall back to standard Viper config loading
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}

			fmt.Println("Config file not found, trying .env file")

			// Try to read from .env file
			viper.SetConfigName(".env")
			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					return nil, fmt.Errorf("error reading .env file: %w", err)
				}
				fmt.Println(".env file not found either")
			} else {
				fmt.Println(".env file found and loaded")
			}
		} else {
			fmt.Println("Config file found and loaded")
		}
	}

	// Set environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults
	setDefaults()

	// Debug output
	fmt.Println("Database URL from viper:", viper.GetString("DATABASE_URL"))
	fmt.Println("Database Auth Token from viper:", viper.GetString("DATABASE_AUTH_TOKEN"))

	// Unmarshal the config
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Debug output after unmarshal
	fmt.Println("Database URL after unmarshal:", config.Database.URL)
	fmt.Println("Database Auth Token after unmarshal:", config.Database.AuthToken)
	fmt.Println("Redis URL after unmarshal:", config.Redis.URL)

	// If database URL is empty after unmarshal, use environment variable directly
	if config.Database.URL == "" {
		config.Database.URL = os.Getenv("DATABASE_URL")
		fmt.Println("Using DATABASE_URL from environment:", config.Database.URL)
	}

	// If database auth token is empty after unmarshal, use environment variable directly
	if config.Database.AuthToken == "" {
		config.Database.AuthToken = os.Getenv("DATABASE_AUTH_TOKEN")
		fmt.Println("Using DATABASE_AUTH_TOKEN from environment:", config.Database.AuthToken)
	}

	// If Redis URL is empty after unmarshal, use environment variable directly
	if config.Redis.URL == "" {
		config.Redis.URL = os.Getenv("REDIS_URL")
		fmt.Println("Using REDIS_URL from environment:", config.Redis.URL)
	}

	// Ensure PASETO key paths are absolute
	if config.PASETO.PublicKeyPath != "" && !filepath.IsAbs(config.PASETO.PublicKeyPath) {
		config.PASETO.PublicKeyPath = filepath.Join(path, config.PASETO.PublicKeyPath)
	}
	if config.PASETO.PrivateKeyPath != "" && !filepath.IsAbs(config.PASETO.PrivateKeyPath) {
		config.PASETO.PrivateKeyPath = filepath.Join(path, config.PASETO.PrivateKeyPath)
	}

	// Validate required configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_READ_TIMEOUT", "15s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "15s")

	// Database defaults
	viper.SetDefault("DATABASE_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DATABASE_MAX_IDLE_CONNS", 25)

	// Redis defaults
	viper.SetDefault("REDIS_URL", "localhost:6379")

	// Auth defaults
	viper.SetDefault("AUTH_ACCESS_TOKEN_TTL", "15m")
	viper.SetDefault("AUTH_REFRESH_TOKEN_TTL", "7d")
	viper.SetDefault("AUTH_VERIFICATION_TTL", "24h")
	viper.SetDefault("AUTH_PASSWORD_RESET_TTL", "1h")
	viper.SetDefault("AUTH_MAX_LOGIN_ATTEMPTS", 5)
	viper.SetDefault("AUTH_LOCKOUT_DURATION", "15m")
	viper.SetDefault("AUTH_SESSION_MAX_LIFETIME", "30d")

	// Email defaults
	viper.SetDefault("EMAIL_LOGIN_NOTIFICATION", true)

	// Security defaults
	viper.SetDefault("max_login_attempts", 5)
	viper.SetDefault("account_lock_duration", 30*time.Minute)
	viper.SetDefault("enable_login_notifications", true)
	viper.SetDefault("enable_suspicious_activity_detection", true)
	viper.SetDefault("enable_rate_limiting", true)
	viper.SetDefault("global_rate_limit", 100)
	viper.SetDefault("auth_rate_limit", 5)

	// App defaults
	viper.SetDefault("name", "Go+Next Fullstack App")
	viper.SetDefault("url", "http://localhost:3000")
	viper.SetDefault("environment", "development")
	viper.SetDefault("debug", true)
}

func validateConfig(config *Config) error {
	if config.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}

	// Only require Redis URL in non-development environments
	if config.Redis.URL == "" && config.Environment != "development" {
		return fmt.Errorf("redis URL is required")
	}

	// Set a default Redis URL in development mode if it's empty
	if config.Redis.URL == "" && config.Environment == "development" {
		config.Redis.URL = "redis:6379"
		fmt.Println("Using default Redis URL in development mode:", config.Redis.URL)
	}

	// Only require email API key in non-development environments
	if config.Email.ResendAPIKey == "" && config.Environment != "development" {
		return fmt.Errorf("email API key is required")
	}

	// Set a default email API key in development mode if it's empty
	if config.Email.ResendAPIKey == "" && config.Environment == "development" {
		config.Email.ResendAPIKey = "test_api_key"
		fmt.Println("Using default email API key in development mode:", config.Email.ResendAPIKey)
	}

	// Only require email from address in non-development environments
	if config.Email.FromEmail == "" && config.Environment != "development" {
		return fmt.Errorf("email from address is required")
	}

	// Set a default email from address in development mode if it's empty
	if config.Email.FromEmail == "" && config.Environment == "development" {
		config.Email.FromEmail = "test@example.com"
		fmt.Println("Using default email from address in development mode:", config.Email.FromEmail)
	}

	// Skip PASETO key file checks in development mode
	if config.Environment != "development" {
		// Check if PASETO key files exist
		if _, err := os.Stat(config.PASETO.PublicKeyPath); err != nil {
			return fmt.Errorf("PASETO public key file not found: %w", err)
		}
		if _, err := os.Stat(config.PASETO.PrivateKeyPath); err != nil {
			return fmt.Errorf("PASETO private key file not found: %w", err)
		}
	}

	return nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Environment: "development",
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Database: DatabaseConfig{
			URL:          "libsql://your-database.turso.io",
			AuthToken:    "turso_auth_token",
			MaxOpenConns: 25,
			MaxIdleConns: 25,
		},
		Redis: RedisConfig{
			URL: "localhost:6379",
		},
		Auth: AuthConfig{
			PublicKey:          "public_key",
			PrivateKey:         "private_key",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    7 * 24 * time.Hour,
			VerificationTTL:    24 * time.Hour,
			PasswordResetTTL:   1 * time.Hour,
			MaxLoginAttempts:   5,
			LockoutDuration:    15 * time.Minute,
			SessionMaxLifetime: 30 * 24 * time.Hour,
		},
		Email: EmailConfig{
			ResendAPIKey:         "resend_api_key",
			FromEmail:            "noreply@example.com",
			FromName:             "Go+Next",
			VerificationURL:      "http://localhost:3000/verify",
			PasswordResetURL:     "http://localhost:3000/reset",
			LoginNotification:    true,
			UpstashWorkflowURL:   "https://api.upstash.com/workflows/workflow_id",
			UpstashWorkflowToken: "upstash_workflow_token",
		},
		OAuth: OAuthConfig{
			Google: struct {
				ClientID     string `mapstructure:"OAUTH_GOOGLE_CLIENT_ID"`
				ClientSecret string `mapstructure:"OAUTH_GOOGLE_CLIENT_SECRET"`
				RedirectURL  string `mapstructure:"OAUTH_GOOGLE_REDIRECT_URL"`
			}{
				ClientID:     "google_client_id",
				ClientSecret: "google_client_secret",
				RedirectURL:  "http://localhost:3000/auth/google/callback",
			},
			GitHub: struct {
				ClientID     string `mapstructure:"OAUTH_GITHUB_CLIENT_ID"`
				ClientSecret string `mapstructure:"OAUTH_GITHUB_CLIENT_SECRET"`
				RedirectURL  string `mapstructure:"OAUTH_GITHUB_REDIRECT_URL"`
			}{
				ClientID:     "github_client_id",
				ClientSecret: "github_client_secret",
				RedirectURL:  "http://localhost:3000/auth/github/callback",
			},
		},
		PASETO: PASETOConfig{
			PublicKeyPath:  "public_key.pem",
			PrivateKeyPath: "private_key.pem",
		},
		Security: SecurityConfig{
			MaxLoginAttempts:                  5,
			AccountLockDuration:               30 * time.Minute,
			EnableLoginNotifications:          true,
			EnableSuspiciousActivityDetection: true,
			EnableRateLimiting:                true,
			GlobalRateLimit:                   100,
			AuthRateLimit:                     5,
		},
		App: AppConfig{
			Name:        "Go+Next Fullstack App",
			URL:         "http://localhost:3000",
			Environment: "development",
			Debug:       true,
		},
	}
}
