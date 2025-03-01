package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/nanayaw/fullstack/internal/config"
	"github.com/resendlabs/resend-go"
)

// ResendService implements the EmailService interface using Resend
type ResendService struct {
	client *resend.Client
	config *config.EmailConfig
}

// NewResendService creates a new ResendService
func NewResendService(cfg *config.EmailConfig) (*ResendService, error) {
	if cfg.ResendAPIKey == "" {
		return nil, fmt.Errorf("resend API key is required")
	}

	client := resend.NewClient(cfg.ResendAPIKey)

	return &ResendService{
		client: client,
		config: cfg,
	}, nil
}

// SendVerificationEmail sends a verification email to the user
func (s *ResendService) SendVerificationEmail(ctx context.Context, to string, token string) error {
	verificationURL := fmt.Sprintf("%s?token=%s", s.config.VerificationURL, token)

	fromEmail := s.config.FromEmail
	if s.config.FromName != "" {
		fromEmail = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{to},
		Subject: "Verify your email address",
		Html:    fmt.Sprintf("<p>Please verify your email by clicking on the following link: <a href=\"%s\">Verify Email</a></p>", verificationURL),
		Text:    fmt.Sprintf("Please verify your email by clicking on the following link: %s", verificationURL),
	}

	_, err := s.client.Emails.Send(params)
	return err
}

// SendPasswordResetEmail sends a password reset email to the user
func (s *ResendService) SendPasswordResetEmail(ctx context.Context, to string, token string) error {
	resetURL := fmt.Sprintf("%s?token=%s", s.config.PasswordResetURL, token)

	fromEmail := s.config.FromEmail
	if s.config.FromName != "" {
		fromEmail = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{to},
		Subject: "Reset your password",
		Html:    fmt.Sprintf("<p>Please reset your password by clicking on the following link: <a href=\"%s\">Reset Password</a></p>", resetURL),
		Text:    fmt.Sprintf("Please reset your password by clicking on the following link: %s", resetURL),
	}

	_, err := s.client.Emails.Send(params)
	return err
}

func (s *ResendService) SendWelcomeEmail(ctx context.Context, to, userName string) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail),
		To:      []string{to},
		Subject: "Welcome to Our Platform",
		Html:    s.getWelcomeEmailTemplate(userName),
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

func (s *ResendService) SendLoginNotificationEmail(ctx context.Context, to, deviceInfo, location string) error {
	if !s.config.LoginNotification {
		return nil
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail),
		To:      []string{to},
		Subject: "New Login Detected",
		Html:    s.getLoginNotificationTemplate(deviceInfo, location),
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send login notification: %w", err)
	}

	return nil
}

func (s *ResendService) SendPasswordChangedEmail(ctx context.Context, to string) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail),
		To:      []string{to},
		Subject: "Your Password Has Been Changed",
		Html:    s.getPasswordChangedTemplate(),
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send password changed notification: %w", err)
	}

	return nil
}

func (s *ResendService) ValidateEmailAddress(email string) bool {
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

// Email templates
func (s *ResendService) getVerificationEmailTemplate(link string) string {
	return fmt.Sprintf(`
		<h2>Verify Your Email Address</h2>
		<p>Please click the link below to verify your email address:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>If you didn't request this, you can safely ignore this email.</p>
	`, link)
}

func (s *ResendService) getPasswordResetEmailTemplate(link string) string {
	return fmt.Sprintf(`
		<h2>Reset Your Password</h2>
		<p>You recently requested to reset your password. Click the link below to proceed:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>If you didn't request this, you can safely ignore this email.</p>
		<p>This link will expire in 1 hour.</p>
	`, link)
}

func (s *ResendService) getWelcomeEmailTemplate(userName string) string {
	return fmt.Sprintf(`
		<h2>Welcome, %s!</h2>
		<p>Thank you for joining our platform. We're excited to have you on board!</p>
		<p>If you have any questions, feel free to reach out to our support team.</p>
	`, userName)
}

func (s *ResendService) getLoginNotificationTemplate(deviceInfo, location string) string {
	return fmt.Sprintf(`
		<h2>New Login Detected</h2>
		<p>We detected a new login to your account from:</p>
		<p>Device: %s</p>
		<p>Location: %s</p>
		<p>If this wasn't you, please change your password immediately and contact support.</p>
	`, deviceInfo, location)
}

func (s *ResendService) getPasswordChangedTemplate() string {
	return `
		<h2>Password Changed Successfully</h2>
		<p>Your password has been successfully changed.</p>
		<p>If you didn't make this change, please contact support immediately.</p>
	`
}

func (s *ResendService) ParseTemplate(templateName string, data interface{}) (string, error) {
	// TODO: Implement proper template parsing
	// For now, return hardcoded templates
	switch templateName {
	case "verification":
		return s.getVerificationEmailTemplate(data.(string)), nil
	case "password_reset":
		return s.getPasswordResetEmailTemplate(data.(string)), nil
	case "welcome":
		return s.getWelcomeEmailTemplate(data.(string)), nil
	case "login_notification":
		d := data.(map[string]string)
		return s.getLoginNotificationTemplate(d["deviceInfo"], d["location"]), nil
	case "password_changed":
		return s.getPasswordChangedTemplate(), nil
	default:
		return "", fmt.Errorf("unknown template: %s", templateName)
	}
}
