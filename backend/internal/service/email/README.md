# Email Service

This package provides email service implementations for the application. It supports multiple email providers through a common interface.

## Supported Providers

### Resend

[Resend](https://resend.com/) is a modern email API for developers. It provides a simple way to send emails from your application.

### Upstash Workflow

[Upstash Workflow](https://upstash.com/docs/workflow/overall/getstarted) is a serverless workflow engine that can be used to send emails and perform other tasks. It's particularly useful for scheduling emails and handling complex email workflows.

## Configuration

The email service is configured through the `EmailConfig` struct in `internal/config/config.go`. The following environment variables are used:

```env
# Email (Resend)
RESEND_API_KEY=your-resend-api-key
EMAIL_FROM_ADDRESS=noreply@yourdomain.com
EMAIL_FROM_NAME=Your App Name
EMAIL_VERIFICATION_URL=http://localhost:3000/verify-email
EMAIL_PASSWORD_RESET_URL=http://localhost:3000/reset-password
EMAIL_LOGIN_NOTIFICATION=true

# Upstash Workflow (for email workflows)
UPSTASH_WORKFLOW_URL=your-upstash-workflow-url
UPSTASH_WORKFLOW_TOKEN=your-upstash-workflow-token
```

## Usage

The email service is initialized in `cmd/api/main.go` using the factory function:

```go
emailService, err := email.NewEmailService(&cfg.Email)
if err != nil {
    log.Fatalf("Failed to initialize email service: %v", err)
}
```

The factory function will choose the appropriate email service implementation based on the configuration:

- If Upstash Workflow is configured (both URL and token are provided), it will use the Upstash Workflow implementation.
- Otherwise, if Resend is configured (API key is provided), it will use the Resend implementation.
- If neither is configured, it will return an error.

## Email Templates

Both implementations provide basic email templates for common use cases:

- Verification emails
- Password reset emails
- Welcome emails
- Login notification emails
- Password changed notification emails

## Extending

To add a new email service implementation:

1. Create a new file in the `email` package (e.g., `sendgrid.go`).
2. Implement the `EmailService` interface defined in `internal/service/interfaces.go`.
3. Add a new factory function in `factory.go` to create the new implementation. 