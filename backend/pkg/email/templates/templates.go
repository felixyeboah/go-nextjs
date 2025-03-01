package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// EmailTemplate represents an email template with subject and body
type EmailTemplate struct {
	Subject string
	HTML    string
	Text    string
}

// TemplateData contains common data for all email templates
type TemplateData struct {
	AppName      string
	Year         int
	SupportEmail string
	BaseURL      string
}

// VerificationData contains data for verification email
type VerificationData struct {
	TemplateData
	UserName        string
	VerificationURL string
	ExpiresIn       string
}

// PasswordResetData contains data for password reset email
type PasswordResetData struct {
	TemplateData
	UserName      string
	ResetURL      string
	ExpiresIn     string
	RequestedFrom string
	RequestTime   string
}

// WelcomeData contains data for welcome email
type WelcomeData struct {
	TemplateData
	UserName string
}

// LoginNotificationData contains data for login notification email
type LoginNotificationData struct {
	TemplateData
	UserName   string
	DeviceInfo string
	Location   string
	IPAddress  string
	Time       string
	UserAgent  string
}

// PasswordChangedData contains data for password changed email
type PasswordChangedData struct {
	TemplateData
	UserName   string
	DeviceInfo string
	Location   string
	Time       string
}

// AccountLockedData contains data for account locked email
type AccountLockedData struct {
	TemplateData
	UserName     string
	UnlockTime   string
	FailedLogins int
}

// SuspiciousActivityData contains data for suspicious activity email
type SuspiciousActivityData struct {
	TemplateData
	UserName     string
	ActivityType string
	DeviceInfo   string
	Location     string
	Time         string
	IPAddress    string
}

// NewTemplateData creates a new TemplateData with default values
func NewTemplateData(appName, supportEmail, baseURL string) TemplateData {
	return TemplateData{
		AppName:      appName,
		Year:         time.Now().Year(),
		SupportEmail: supportEmail,
		BaseURL:      baseURL,
	}
}

// RenderTemplate renders a template with the given data
func RenderTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return buf.String(), nil
}

// GetVerificationEmail returns the verification email template
func GetVerificationEmail(data VerificationData) (EmailTemplate, error) {
	htmlTmpl, err := template.New("verification_html").Parse(verificationHTMLTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	textTmpl, err := template.New("verification_text").Parse(verificationTextTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse text template: %w", err)
	}

	html, err := RenderTemplate(htmlTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	text, err := RenderTemplate(textTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	return EmailTemplate{
		Subject: fmt.Sprintf("Verify your email address for %s", data.AppName),
		HTML:    html,
		Text:    text,
	}, nil
}

// GetPasswordResetEmail returns the password reset email template
func GetPasswordResetEmail(data PasswordResetData) (EmailTemplate, error) {
	htmlTmpl, err := template.New("password_reset_html").Parse(passwordResetHTMLTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	textTmpl, err := template.New("password_reset_text").Parse(passwordResetTextTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse text template: %w", err)
	}

	html, err := RenderTemplate(htmlTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	text, err := RenderTemplate(textTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	return EmailTemplate{
		Subject: fmt.Sprintf("Reset your password for %s", data.AppName),
		HTML:    html,
		Text:    text,
	}, nil
}

// GetWelcomeEmail returns the welcome email template
func GetWelcomeEmail(data WelcomeData) (EmailTemplate, error) {
	htmlTmpl, err := template.New("welcome_html").Parse(welcomeHTMLTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	textTmpl, err := template.New("welcome_text").Parse(welcomeTextTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse text template: %w", err)
	}

	html, err := RenderTemplate(htmlTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	text, err := RenderTemplate(textTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	return EmailTemplate{
		Subject: fmt.Sprintf("Welcome to %s", data.AppName),
		HTML:    html,
		Text:    text,
	}, nil
}

// GetLoginNotificationEmail returns the login notification email template
func GetLoginNotificationEmail(data LoginNotificationData) (EmailTemplate, error) {
	htmlTmpl, err := template.New("login_notification_html").Parse(loginNotificationHTMLTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	textTmpl, err := template.New("login_notification_text").Parse(loginNotificationTextTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse text template: %w", err)
	}

	html, err := RenderTemplate(htmlTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	text, err := RenderTemplate(textTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	return EmailTemplate{
		Subject: fmt.Sprintf("New login to your %s account", data.AppName),
		HTML:    html,
		Text:    text,
	}, nil
}

// GetPasswordChangedEmail returns the password changed email template
func GetPasswordChangedEmail(data PasswordChangedData) (EmailTemplate, error) {
	htmlTmpl, err := template.New("password_changed_html").Parse(passwordChangedHTMLTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	textTmpl, err := template.New("password_changed_text").Parse(passwordChangedTextTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse text template: %w", err)
	}

	html, err := RenderTemplate(htmlTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	text, err := RenderTemplate(textTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	return EmailTemplate{
		Subject: fmt.Sprintf("Your %s password has been changed", data.AppName),
		HTML:    html,
		Text:    text,
	}, nil
}

// GetAccountLockedEmail returns the account locked email template
func GetAccountLockedEmail(data AccountLockedData) (EmailTemplate, error) {
	htmlTmpl, err := template.New("account_locked_html").Parse(accountLockedHTMLTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	textTmpl, err := template.New("account_locked_text").Parse(accountLockedTextTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse text template: %w", err)
	}

	html, err := RenderTemplate(htmlTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	text, err := RenderTemplate(textTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	return EmailTemplate{
		Subject: fmt.Sprintf("Your %s account has been temporarily locked", data.AppName),
		HTML:    html,
		Text:    text,
	}, nil
}

// GetSuspiciousActivityEmail returns the suspicious activity email template
func GetSuspiciousActivityEmail(data SuspiciousActivityData) (EmailTemplate, error) {
	htmlTmpl, err := template.New("suspicious_activity_html").Parse(suspiciousActivityHTMLTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	textTmpl, err := template.New("suspicious_activity_text").Parse(suspiciousActivityTextTemplate)
	if err != nil {
		return EmailTemplate{}, fmt.Errorf("failed to parse text template: %w", err)
	}

	html, err := RenderTemplate(htmlTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	text, err := RenderTemplate(textTmpl, data)
	if err != nil {
		return EmailTemplate{}, err
	}

	return EmailTemplate{
		Subject: fmt.Sprintf("Suspicious activity detected on your %s account", data.AppName),
		HTML:    html,
		Text:    text,
	}, nil
}
