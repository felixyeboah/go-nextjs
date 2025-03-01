package templates

// Plain text email templates for better accessibility and plain text email clients

const verificationTextTemplate = `Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},

Thank you for signing up for {{.AppName}}. To complete your registration and verify your email address, please visit the following link:

{{.VerificationURL}}

This verification link will expire in {{.ExpiresIn}}.

If you didn't create an account with {{.AppName}}, you can safely ignore this email.

Need help? Contact our support team at {{.SupportEmail}}.

© {{.Year}} {{.AppName}}. All rights reserved.`

// #nosec G101 - This is a template for password reset emails, not a hardcoded credential
const passwordResetTextTemplate = `Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},

We received a request to reset your password for your {{.AppName}} account. To reset your password, please visit the following link:

{{.ResetURL}}

This password reset link will expire in {{.ExpiresIn}}.

Security Information:
- This request was made from: {{.RequestedFrom}}
- Time of request: {{.RequestTime}}

If you didn't request a password reset, please ignore this email or contact our support team if you have concerns about your account security.

Need help? Contact our support team at {{.SupportEmail}}.

© {{.Year}} {{.AppName}}. All rights reserved.`

const welcomeTextTemplate = `Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},

Thank you for joining {{.AppName}}! We're excited to have you on board.

Here are a few things you can do with your new account:
- Complete your profile - Add your information to get the most out of our platform.
- Explore our features - Discover all the tools and services we offer.
- Connect with others - Build your network and collaborate with like-minded individuals.

Visit your dashboard: {{.BaseURL}}/dashboard

Need help getting started? Check out our Help Center at {{.BaseURL}}/help or contact our support team at {{.SupportEmail}}.

© {{.Year}} {{.AppName}}. All rights reserved.`

const loginNotificationTextTemplate = `Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},

We detected a new login to your {{.AppName}} account. If this was you, no action is needed.

Login Details:
- Device: {{.DeviceInfo}}
- Location: {{.Location}}
- IP Address: {{.IPAddress}}
- Time: {{.Time}}
- Browser: {{.UserAgent}}

If you don't recognize this activity, please visit {{.BaseURL}}/account/security to secure your account.

If you didn't authorize this login, please change your password immediately and contact our support team at {{.SupportEmail}}.

© {{.Year}} {{.AppName}}. All rights reserved.`

// #nosec G101 - This is a template for password changed emails, not a hardcoded credential
const passwordChangedTextTemplate = `Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},

This email confirms that your password for your {{.AppName}} account has been successfully changed.

Change Details:
- Device: {{.DeviceInfo}}
- Location: {{.Location}}
- Time: {{.Time}}

If you didn't make this change, please visit {{.BaseURL}}/account/security to secure your account.

If you didn't authorize this password change, please contact our support team immediately at {{.SupportEmail}}.

© {{.Year}} {{.AppName}}. All rights reserved.`

const accountLockedTextTemplate = `Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},

For your security, we've temporarily locked your {{.AppName}} account due to multiple failed login attempts.

Lock Details:
- Failed Login Attempts: {{.FailedLogins}}
- Account Will Unlock: {{.UnlockTime}}

If you were trying to log in, please wait until the account unlocks and try again with the correct password. If you've forgotten your password, you can reset it by visiting:

{{.BaseURL}}/reset-password

If you didn't attempt to log in and believe someone else might be trying to access your account, please contact our support team immediately at {{.SupportEmail}}.

© {{.Year}} {{.AppName}}. All rights reserved.`

const suspiciousActivityTextTemplate = `Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},

We've detected suspicious activity on your {{.AppName}} account that requires your immediate attention.

Activity Details:
- Activity: {{.ActivityType}}
- Device: {{.DeviceInfo}}
- Location: {{.Location}}
- IP Address: {{.IPAddress}}
- Time: {{.Time}}

For your security, we recommend taking immediate action by visiting:
{{.BaseURL}}/account/security

If you recognize this activity, you can safely ignore this email. If not, please change your password immediately and contact our support team at {{.SupportEmail}}.

© {{.Year}} {{.AppName}}. All rights reserved.`
