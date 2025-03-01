# Security Documentation

This document provides detailed information about the security features, best practices, and scanning tools used in the fullstack application.

## Security Features

### Authentication Security

- **PASETO Tokens**: The application uses [PASETO (Platform-Agnostic Security Tokens)](https://github.com/paragonie/paseto) for secure authentication, which provides stronger security guarantees than JWT.
- **Password Hashing**: User passwords are hashed using Argon2id, a modern password hashing algorithm designed to be resistant to both brute-force and side-channel attacks.
- **Rate Limiting**: The application implements rate limiting on authentication endpoints to prevent brute force attacks.
- **Account Locking**: After multiple failed login attempts, user accounts are temporarily locked to prevent unauthorized access.
- **Session Management**: The application implements secure session management with token rotation and expiration.

### API Security

- **Input Validation**: All user input is validated using strict schemas to prevent injection attacks and ensure data integrity.
- **Error Handling**: The application implements proper error handling to prevent information leakage.
- **Security Headers**: HTTP security headers are set by default to protect against various attacks:
  - Content-Security-Policy
  - X-Content-Type-Options
  - X-Frame-Options
  - X-XSS-Protection
  - Strict-Transport-Security
- **CSRF Protection**: Cross-Site Request Forgery protection is implemented for all state-changing operations.
- **XSS Protection**: Cross-Site Scripting protection is implemented through proper output encoding and content security policies.
- **SQL Injection Protection**: Parameterized queries and ORM tools prevent SQL injection attacks.
- **Timeout Middleware**: Prevents long-running requests that could lead to denial of service.

### Data Security

- **TLS Encryption**: All communication between client and server is encrypted using TLS.
- **Sensitive Data Handling**: Sensitive data is encrypted at rest and in transit.
- **Database Security**: The database is configured with secure defaults and access is restricted.
- **Logging**: Security-relevant events are logged for audit purposes, but sensitive information is never logged.

## Security Scanning Tools

The application includes automated security scanning tools to identify vulnerabilities:

### Gosec

[Gosec](https://github.com/securego/gosec) is a static analysis tool for Go code that finds security issues. It scans the codebase for common security issues such as:

- Hardcoded credentials
- SQL injection
- Command injection
- Insecure file operations
- Weak cryptography
- And more

To run Gosec:

```bash
cd backend
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec -quiet -exclude-dir=vendor ./...
```

### Nancy

[Nancy](https://github.com/sonatype-nexus-community/nancy) is a tool to check for vulnerabilities in Go dependencies. It scans the dependency tree for known vulnerabilities in the National Vulnerability Database.

To run Nancy:

```bash
cd backend
go install github.com/sonatype-nexus-community/nancy@latest
go list -json -deps ./... | nancy sleuth
```

### Security Check Script

The application includes a security check script (`backend/scripts/security_check.sh`) that automates the process of running these security scanning tools. The script:

1. Checks if the required tools are installed, and installs them if necessary
2. Runs Gosec to scan the codebase for security issues
3. Runs Nancy to check dependencies for known vulnerabilities

To run the security check script:

```bash
cd backend
make security-check
```

## Security Best Practices

The application follows these security best practices:

### Code Security

- **Dependency Management**: Dependencies are regularly updated to include security patches.
- **Code Reviews**: All code changes undergo security-focused code reviews.
- **Secure Defaults**: The application is configured with secure defaults.
- **Principle of Least Privilege**: Components only have access to the resources they need.

### Authentication and Authorization

- **Strong Password Policies**: The application enforces strong password policies.
- **Multi-Factor Authentication**: Support for multi-factor authentication is available.
- **Role-Based Access Control**: Access to resources is controlled based on user roles.
- **Session Management**: Sessions are managed securely with proper expiration and rotation.

### Infrastructure Security

- **Container Security**: Docker images are built with minimal attack surface.
- **Network Security**: Services communicate over a private network.
- **Environment Isolation**: Development, staging, and production environments are isolated.
- **Secret Management**: Sensitive configuration values are managed securely.

## Security Response

If you discover a security vulnerability in this project, please follow these steps:

1. **Do not disclose the vulnerability publicly** until it has been addressed by the maintainers.
2. Email the details to security@example.com or open a private security advisory on GitHub.
3. Include as much information as possible, including:
   - A description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fixes (if any)

The maintainers will acknowledge receipt of your report within 48 hours and provide a more detailed response within 72 hours, including:

- Confirmation of the vulnerability
- The planned timeline for addressing the issue
- Any additional questions or information needed

## Security Updates

Security updates are released as soon as possible after a vulnerability is discovered and fixed. Users are encouraged to:

1. Keep dependencies up to date
2. Monitor the project's security advisories
3. Apply security patches promptly

## Compliance

The application is designed with security and privacy in mind, following best practices for:

- GDPR compliance
- OWASP Top 10 security risks
- NIST Cybersecurity Framework

## Additional Resources

- [OWASP Go Security Cheatsheet](https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Go_Security_Cheatsheet.md)
- [Go Security Best Practices](https://blog.sqreen.com/go-security-best-practices/)
- [Docker Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)

## Environment Variables and Sensitive Data

- **Never commit sensitive data**: API keys, tokens, passwords, or other credentials should never be committed to the repository.
- **Use .env files**: Store sensitive data in .env files which are ignored by git.
- **Use .env.example files**: Provide example .env files with placeholders instead of real credentials.
- **Check for leaks**: Regularly audit the codebase for accidentally committed sensitive data.

If you discover that you've accidentally committed sensitive data:
1. Immediately change any exposed credentials
2. Remove the sensitive data from git history using git-filter-repo
3. Force push the changes to the repository

## Reporting a Vulnerability

If you discover a security vulnerability within this project, please send an email to [security@example.com](mailto:security@example.com). All security vulnerabilities will be promptly addressed.

## Security Checks

This project uses the following tools for security checks:

- **gosec**: Static analysis for Go code to find potential security issues
- **nancy**: Checks for vulnerabilities in dependencies

To run security checks locally:

```bash
cd backend
make security-check
```

## Secure Coding Practices

This project follows these secure coding practices:

1. **Input Validation**: All user inputs are validated before processing
2. **Proper Error Handling**: Errors are handled appropriately without leaking sensitive information
3. **Secure Password Storage**: Passwords are hashed using bcrypt
4. **Rate Limiting**: API endpoints are protected against brute force attacks
5. **HTTPS Only**: All communications use HTTPS
6. **Content Security Policy**: Implemented to prevent XSS attacks
7. **Regular Dependency Updates**: Dependencies are regularly updated to patch security vulnerabilities

## Security Features

The application includes the following security features:

- PASETO tokens for secure authentication
- Email verification for new accounts
- Two-factor authentication (optional)
- Account lockout after multiple failed login attempts
- Secure password reset flow
- Login notification emails
- Session management and revocation
- Audit logging for security events 