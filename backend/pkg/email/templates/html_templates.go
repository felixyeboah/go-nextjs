package templates

// HTML email templates with modern, responsive design

const verificationHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 20px;
        }
        h1 {
            color: #2563eb;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #2563eb;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            font-weight: bold;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #1d4ed8;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #6b7280;
        }
        .expires {
            font-style: italic;
            margin: 20px 0;
            color: #6b7280;
        }
        .help {
            margin-top: 20px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email Address</h1>
        </div>
        
        <p>Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},</p>
        
        <p>Thank you for signing up for {{.AppName}}. To complete your registration and verify your email address, please click the button below:</p>
        
        <div style="text-align: center;">
            <a href="{{.VerificationURL}}" class="button">Verify Email Address</a>
        </div>
        
        <p class="expires">This verification link will expire in {{.ExpiresIn}}.</p>
        
        <p>If you didn't create an account with {{.AppName}}, you can safely ignore this email.</p>
        
        <p>If you're having trouble clicking the button, copy and paste the following URL into your web browser:</p>
        <p style="word-break: break-all; font-size: 14px;">{{.VerificationURL}}</p>
        
        <div class="help">
            <p>Need help? Contact our support team at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>.</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

// #nosec G101 - This is a template for password reset emails, not a hardcoded credential
const passwordResetHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 20px;
        }
        h1 {
            color: #2563eb;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #2563eb;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            font-weight: bold;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #1d4ed8;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #6b7280;
        }
        .expires {
            font-style: italic;
            margin: 20px 0;
            color: #6b7280;
        }
        .help {
            margin-top: 20px;
            font-size: 14px;
        }
        .security-info {
            background-color: #f9fafb;
            padding: 15px;
            border-radius: 4px;
            margin: 20px 0;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Reset Your Password</h1>
        </div>
        
        <p>Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},</p>
        
        <p>We received a request to reset your password for your {{.AppName}} account. To reset your password, please click the button below:</p>
        
        <div style="text-align: center;">
            <a href="{{.ResetURL}}" class="button">Reset Password</a>
        </div>
        
        <p class="expires">This password reset link will expire in {{.ExpiresIn}}.</p>
        
        <div class="security-info">
            <p><strong>Security Information:</strong></p>
            <p>This request was made from: {{.RequestedFrom}}</p>
            <p>Time of request: {{.RequestTime}}</p>
        </div>
        
        <p>If you didn't request a password reset, please ignore this email or contact our support team if you have concerns about your account security.</p>
        
        <p>If you're having trouble clicking the button, copy and paste the following URL into your web browser:</p>
        <p style="word-break: break-all; font-size: 14px;">{{.ResetURL}}</p>
        
        <div class="help">
            <p>Need help? Contact our support team at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>.</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const welcomeHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to {{.AppName}}</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 20px;
        }
        h1 {
            color: #2563eb;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #2563eb;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            font-weight: bold;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #1d4ed8;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #6b7280;
        }
        .feature-list {
            margin: 30px 0;
        }
        .feature-item {
            margin-bottom: 15px;
        }
        .help {
            margin-top: 20px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to {{.AppName}}!</h1>
        </div>
        
        <p>Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},</p>
        
        <p>Thank you for joining {{.AppName}}! We're excited to have you on board.</p>
        
        <div class="feature-list">
            <h3>Here are a few things you can do with your new account:</h3>
            <div class="feature-item">
                <p>✅ <strong>Complete your profile</strong> - Add your information to get the most out of our platform.</p>
            </div>
            <div class="feature-item">
                <p>✅ <strong>Explore our features</strong> - Discover all the tools and services we offer.</p>
            </div>
            <div class="feature-item">
                <p>✅ <strong>Connect with others</strong> - Build your network and collaborate with like-minded individuals.</p>
            </div>
        </div>
        
        <div style="text-align: center;">
            <a href="{{.BaseURL}}/dashboard" class="button">Go to Dashboard</a>
        </div>
        
        <div class="help">
            <p>Need help getting started? Check out our <a href="{{.BaseURL}}/help">Help Center</a> or contact our support team at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>.</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const loginNotificationHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Login Detected</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 20px;
        }
        h1 {
            color: #2563eb;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #ef4444;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            font-weight: bold;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #dc2626;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #6b7280;
        }
        .login-details {
            background-color: #f9fafb;
            padding: 20px;
            border-radius: 4px;
            margin: 20px 0;
        }
        .detail-row {
            display: flex;
            margin-bottom: 10px;
        }
        .detail-label {
            font-weight: bold;
            width: 120px;
        }
        .help {
            margin-top: 20px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>New Login to Your Account</h1>
        </div>
        
        <p>Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},</p>
        
        <p>We detected a new login to your {{.AppName}} account. If this was you, no action is needed.</p>
        
        <div class="login-details">
            <h3>Login Details:</h3>
            <div class="detail-row">
                <div class="detail-label">Device:</div>
                <div>{{.DeviceInfo}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Location:</div>
                <div>{{.Location}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">IP Address:</div>
                <div>{{.IPAddress}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Time:</div>
                <div>{{.Time}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Browser:</div>
                <div>{{.UserAgent}}</div>
            </div>
        </div>
        
        <p><strong>If you don't recognize this activity:</strong></p>
        <div style="text-align: center;">
            <a href="{{.BaseURL}}/account/security" class="button">Secure Your Account</a>
        </div>
        
        <div class="help">
            <p>If you didn't authorize this login, please change your password immediately and contact our support team at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>.</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

// #nosec G101 - This is a template for password changed emails, not a hardcoded credential
const passwordChangedHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Changed</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 20px;
        }
        h1 {
            color: #2563eb;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #ef4444;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            font-weight: bold;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #dc2626;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #6b7280;
        }
        .change-details {
            background-color: #f9fafb;
            padding: 20px;
            border-radius: 4px;
            margin: 20px 0;
        }
        .detail-row {
            display: flex;
            margin-bottom: 10px;
        }
        .detail-label {
            font-weight: bold;
            width: 120px;
        }
        .help {
            margin-top: 20px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Your Password Has Been Changed</h1>
        </div>
        
        <p>Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},</p>
        
        <p>This email confirms that your password for your {{.AppName}} account has been successfully changed.</p>
        
        <div class="change-details">
            <h3>Change Details:</h3>
            <div class="detail-row">
                <div class="detail-label">Device:</div>
                <div>{{.DeviceInfo}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Location:</div>
                <div>{{.Location}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Time:</div>
                <div>{{.Time}}</div>
            </div>
        </div>
        
        <p><strong>If you didn't make this change:</strong></p>
        <div style="text-align: center;">
            <a href="{{.BaseURL}}/account/security" class="button">Secure Your Account</a>
        </div>
        
        <div class="help">
            <p>If you didn't authorize this password change, please contact our support team immediately at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>.</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const accountLockedHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Account Temporarily Locked</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 20px;
        }
        h1 {
            color: #2563eb;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #2563eb;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            font-weight: bold;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #1d4ed8;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #6b7280;
        }
        .lock-details {
            background-color: #f9fafb;
            padding: 20px;
            border-radius: 4px;
            margin: 20px 0;
        }
        .help {
            margin-top: 20px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Your Account Has Been Temporarily Locked</h1>
        </div>
        
        <p>Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},</p>
        
        <p>For your security, we've temporarily locked your {{.AppName}} account due to multiple failed login attempts.</p>
        
        <div class="lock-details">
            <h3>Lock Details:</h3>
            <p><strong>Failed Login Attempts:</strong> {{.FailedLogins}}</p>
            <p><strong>Account Will Unlock:</strong> {{.UnlockTime}}</p>
        </div>
        
        <p>If you were trying to log in, please wait until the account unlocks and try again with the correct password. If you've forgotten your password, you can reset it using the button below:</p>
        
        <div style="text-align: center;">
            <a href="{{.BaseURL}}/reset-password" class="button">Reset Password</a>
        </div>
        
        <div class="help">
            <p>If you didn't attempt to log in and believe someone else might be trying to access your account, please contact our support team immediately at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>.</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const suspiciousActivityHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Suspicious Activity Detected</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 20px;
        }
        h1 {
            color: #ef4444;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #ef4444;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            font-weight: bold;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #dc2626;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #6b7280;
        }
        .activity-details {
            background-color: #f9fafb;
            padding: 20px;
            border-radius: 4px;
            margin: 20px 0;
        }
        .detail-row {
            display: flex;
            margin-bottom: 10px;
        }
        .detail-label {
            font-weight: bold;
            width: 120px;
        }
        .help {
            margin-top: 20px;
            font-size: 14px;
        }
        .warning {
            color: #ef4444;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Suspicious Activity Detected</h1>
        </div>
        
        <p>Hello {{if .UserName}}{{.UserName}}{{else}}there{{end}},</p>
        
        <p class="warning">We've detected suspicious activity on your {{.AppName}} account that requires your immediate attention.</p>
        
        <div class="activity-details">
            <h3>Activity Details:</h3>
            <div class="detail-row">
                <div class="detail-label">Activity:</div>
                <div>{{.ActivityType}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Device:</div>
                <div>{{.DeviceInfo}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Location:</div>
                <div>{{.Location}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">IP Address:</div>
                <div>{{.IPAddress}}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Time:</div>
                <div>{{.Time}}</div>
            </div>
        </div>
        
        <p><strong>For your security, we recommend taking immediate action:</strong></p>
        <div style="text-align: center;">
            <a href="{{.BaseURL}}/account/security" class="button">Secure Your Account</a>
        </div>
        
        <div class="help">
            <p>If you recognize this activity, you can safely ignore this email. If not, please change your password immediately and contact our support team at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>.</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`
