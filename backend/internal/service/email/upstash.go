package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nanayaw/fullstack/internal/config"
)

// UpstashWorkflowService implements the EmailService interface using Upstash Workflow
type UpstashWorkflowService struct {
	config           *config.EmailConfig
	workflowURL      string
	workflowToken    string
	verificationTTL  time.Duration
	passwordResetTTL time.Duration
}

// WorkflowRequest represents a request to the Upstash Workflow API
type WorkflowRequest struct {
	Name    string                 `json:"name"`
	Data    map[string]interface{} `json:"data"`
	Delay   *int                   `json:"delay,omitempty"`
	Cron    *string                `json:"cron,omitempty"`
	Retries *int                   `json:"retries,omitempty"`
}

// NewUpstashWorkflowService creates a new UpstashWorkflowService
func NewUpstashWorkflowService(cfg *config.EmailConfig) (*UpstashWorkflowService, error) {
	if cfg.UpstashWorkflowURL == "" {
		return nil, fmt.Errorf("upstash workflow URL is required")
	}
	if cfg.UpstashWorkflowToken == "" {
		return nil, fmt.Errorf("upstash workflow token is required")
	}

	return &UpstashWorkflowService{
		config:           cfg,
		workflowURL:      cfg.UpstashWorkflowURL,
		workflowToken:    cfg.UpstashWorkflowToken,
		verificationTTL:  cfg.VerificationTTL,
		passwordResetTTL: cfg.PasswordResetTTL,
	}, nil
}

// SendVerificationEmail sends a verification email to the user
func (s *UpstashWorkflowService) SendVerificationEmail(ctx context.Context, to string, token string) error {
	verificationURL := fmt.Sprintf("%s?token=%s", s.config.VerificationURL, token)

	data := map[string]interface{}{
		"to":      to,
		"from":    s.config.FromEmail,
		"subject": "Verify your email address",
		"body":    fmt.Sprintf("Please verify your email by clicking on the following link: %s", verificationURL),
		"html":    fmt.Sprintf("<p>Please verify your email by clicking on the following link: <a href=\"%s\">Verify Email</a></p>", verificationURL),
	}

	return s.triggerWorkflow("send-email", data)
}

// SendPasswordResetEmail sends a password reset email to the user
func (s *UpstashWorkflowService) SendPasswordResetEmail(ctx context.Context, to string, token string) error {
	resetURL := fmt.Sprintf("%s?token=%s", s.config.PasswordResetURL, token)

	data := map[string]interface{}{
		"to":      to,
		"from":    s.config.FromEmail,
		"subject": "Reset your password",
		"body":    fmt.Sprintf("Please reset your password by clicking on the following link: %s", resetURL),
		"html":    fmt.Sprintf("<p>Please reset your password by clicking on the following link: <a href=\"%s\">Reset Password</a></p>", resetURL),
	}

	return s.triggerWorkflow("send-email", data)
}

// triggerWorkflow sends a request to the Upstash Workflow API
func (s *UpstashWorkflowService) triggerWorkflow(name string, data map[string]interface{}) error {
	workflowReq := WorkflowRequest{
		Name: name,
		Data: data,
	}

	jsonData, err := json.Marshal(workflowReq)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow request: %w", err)
	}

	req, err := http.NewRequest("POST", s.workflowURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.workflowToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("workflow API returned error status: %d", resp.StatusCode)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email
func (s *UpstashWorkflowService) SendWelcomeEmail(ctx context.Context, to, userName string) error {
	data := map[string]interface{}{
		"userName": userName,
	}

	return s.triggerWorkflow("welcome", data)
}

// SendLoginNotificationEmail sends a login notification email
func (s *UpstashWorkflowService) SendLoginNotificationEmail(ctx context.Context, to, deviceInfo, location string) error {
	if !s.config.LoginNotification {
		return nil
	}

	data := map[string]interface{}{
		"deviceInfo": deviceInfo,
		"location":   location,
		"time":       time.Now().Format(time.RFC1123),
	}

	return s.triggerWorkflow("login_notification", data)
}

// SendPasswordChangedEmail sends a password changed notification email
func (s *UpstashWorkflowService) SendPasswordChangedEmail(ctx context.Context, to string) error {
	data := map[string]interface{}{
		"time": time.Now().Format(time.RFC1123),
	}

	return s.triggerWorkflow("password_changed", data)
}

// ValidateEmailAddress validates an email address
func (s *UpstashWorkflowService) ValidateEmailAddress(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	return true
}

// ParseTemplate parses a template with data
func (s *UpstashWorkflowService) ParseTemplate(templateName string, data interface{}) (string, error) {
	// Not needed for Upstash Workflow as templates are managed by the workflow
	return "", nil
}
