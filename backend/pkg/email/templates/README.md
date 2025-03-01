# Email Templates Package

This package provides a collection of responsive, accessible email templates for various authentication and security-related communications with users.

## Available Templates

The package includes the following email templates:

1. **Verification Email** - Sent to users to verify their email address during registration
2. **Password Reset Email** - Sent when a user requests a password reset
3. **Welcome Email** - Sent to users after successful registration and verification
4. **Login Notification** - Alerts users about new logins to their account
5. **Password Changed** - Notifies users when their password has been changed
6. **Account Locked** - Informs users when their account has been temporarily locked due to failed login attempts
7. **Suspicious Activity** - Alerts users about potentially suspicious activity on their account

## Features

- **Responsive Design** - All HTML templates are responsive and work well on mobile devices
- **Accessibility** - Each template includes both HTML and plain text versions for better accessibility
- **Customizable** - Templates can be customized with your application's name, support email, and other details
- **Security-Focused** - Templates include security information such as device details, location, and IP addresses where relevant

## Usage

### Basic Usage

```go
import (
    "github.com/yourusername/yourproject/pkg/email/templates"
)

// Create template data with your app details
baseData := templates.NewTemplateData(
    "Your App Name",
    "support@yourapp.com",
    "https://yourapp.com",
)

// Create verification email data
verificationData := templates.VerificationData{
    TemplateData:     baseData,
    UserName:         "John Doe",
    VerificationURL:  "https://yourapp.com/verify?token=abc123",
    ExpiresIn:        "24 hours",
}

// Get the email template
emailTemplate, err := templates.GetVerificationEmail(verificationData)
if err != nil {
    // Handle error
}

// Use the template in your email service
emailSubject := emailTemplate.Subject
htmlBody := emailTemplate.HTML
textBody := emailTemplate.Text

// Send the email using your email service
// ...
```

### Available Template Data Structures

Each template has its own data structure with specific fields:

- `VerificationData` - For email verification
- `PasswordResetData` - For password reset
- `WelcomeData` - For welcome emails
- `LoginNotificationData` - For login notifications
- `PasswordChangedData` - For password change notifications
- `AccountLockedData` - For account locked notifications
- `SuspiciousActivityData` - For suspicious activity alerts

## Customization

You can customize the templates by modifying the HTML and text template constants in the `html_templates.go` and `text_templates.go` files. The templates use Go's standard template package for rendering.

## Best Practices

1. **Always include both HTML and text versions** when sending emails
2. **Personalize emails** with the user's name when available
3. **Include clear security information** in security-related emails
4. **Provide clear actions** for users to take if they didn't initiate the action
5. **Include contact information** for your support team 