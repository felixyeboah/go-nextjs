package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResendClient struct {
	apiKey    string
	fromEmail string
	fromName  string
	baseURL   string
}

type EmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

type EmailResponse struct {
	ID string `json:"id"`
}

func NewResendClient(apiKey, fromEmail, fromName string) *ResendClient {
	return &ResendClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
		baseURL:   "https://api.resend.com",
	}
}

func (c *ResendClient) SendEmail(to []string, subject, htmlContent string) error {
	from := fmt.Sprintf("%s <%s>", c.fromName, c.fromEmail)

	emailReq := EmailRequest{
		From:    from,
		To:      to,
		Subject: subject,
		Html:    htmlContent,
	}

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return fmt.Errorf("error marshaling email request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var emailResp EmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&emailResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email to a new user
func (c *ResendClient) SendWelcomeEmail(to, username string) error {
	subject := "Welcome to Our Platform!"
	htmlContent := fmt.Sprintf(`
		<h1>Welcome to Our Platform, %s!</h1>
		<p>We're excited to have you on board. Here are some next steps to get started:</p>
		<ul>
			<li>Complete your profile</li>
			<li>Explore our features</li>
			<li>Connect with others</li>
		</ul>
		<p>If you have any questions, feel free to reach out to our support team.</p>
	`, username)

	return c.SendEmail([]string{to}, subject, htmlContent)
}

// SendPasswordResetEmail sends a password reset email
func (c *ResendClient) SendPasswordResetEmail(to, resetToken string) error {
	subject := "Password Reset Request"
	htmlContent := fmt.Sprintf(`
		<h1>Password Reset Request</h1>
		<p>We received a request to reset your password. Click the link below to proceed:</p>
		<p><a href="%s/reset-password?token=%s">Reset Password</a></p>
		<p>If you didn't request this, please ignore this email.</p>
		<p>This link will expire in 1 hour.</p>
	`, c.baseURL, resetToken)

	return c.SendEmail([]string{to}, subject, htmlContent)
}
