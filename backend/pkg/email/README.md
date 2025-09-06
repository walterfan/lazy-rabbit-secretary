# Email Package

This package provides a simple and flexible email sending functionality for Go applications using SMTP.

## Features

- ✅ SMTP authentication support
- ✅ TLS/SSL encryption support  
- ✅ Configuration via environment variables
- ✅ Multiple recipients (To, CC)
- ✅ Custom email headers
- ✅ Connection testing
- ✅ Error handling and validation
- ✅ Support for various SMTP providers (Gmail, 163.com, etc.)

## Configuration

The email sender can be configured using environment variables:

### Required Environment Variables

```bash
MAIL_SERVER=smtp.163.com          # SMTP server hostname
MAIL_PORT=587                     # SMTP server port
MAIL_USERNAME=your-email@163.com  # SMTP username
MAIL_PASSWORD=your-password       # SMTP password
```

### Optional Environment Variables

```bash
MAIL_USE_SSL=false               # Use SSL connection (default: false)
MAIL_USE_TLS=true                # Use TLS encryption (default: true)
MAIL_SENDER=your-email@163.com   # Default sender address (default: MAIL_USERNAME)
MAIL_RECEIVER=default@email.com  # Default recipient address
```

## Usage

### Basic Usage

```go
package main

import (
    "log"
    "your-project/pkg/email"
)

func main() {
    // Create email sender from environment variables
    sender, err := email.NewEmailSender()
    if err != nil {
        log.Fatal(err)
    }

    // Test connection
    if err := sender.TestConnection(); err != nil {
        log.Fatalf("Connection test failed: %v", err)
    }

    // Send simple email
    err = sender.SendSimpleEmail(
        "Test Subject",
        "This is the email body.",
    )
    if err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}
```

### Advanced Usage

```go
// Send email to specific recipients
err = sender.SendEmailToRecipients(
    "Important Notice",
    "This is an important message.",
    []string{"user1@example.com", "user2@example.com"}, // To
    []string{"manager@example.com"},                     // CC
)

// Send custom email message
message := &email.EmailMessage{
    Subject:  "Custom Subject",
    Body:     "Custom email body",
    ToAddr:   []string{"recipient@example.com"},
    CCAddr:   []string{"cc@example.com"},
    FromAddr: "custom-sender@example.com",
}
err = sender.SendEmail(message)
```

### Custom Configuration

```go
config := &email.EmailConfig{
    SMTPHost:    "smtp.gmail.com",
    SMTPPort:    587,
    UseSSL:      false,
    UseTLS:      true,
    Username:    "your-email@gmail.com",
    Password:    "your-app-password",
    DefaultFrom: "your-email@gmail.com",
    DefaultTo:   "default@example.com",
}

sender := email.NewEmailSenderWithConfig(config)
```

## SMTP Provider Examples

### Gmail Configuration

```bash
MAIL_SERVER=smtp.gmail.com
MAIL_PORT=587
MAIL_USE_SSL=false
MAIL_USE_TLS=true
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-app-password  # Use App Password, not regular password
```

### 163.com Configuration

```bash
MAIL_SERVER=smtp.163.com
MAIL_PORT=587
MAIL_USE_SSL=false
MAIL_USE_TLS=true
MAIL_USERNAME=your-email@163.com
MAIL_PASSWORD=your-password
```

### Outlook/Hotmail Configuration

```bash
MAIL_SERVER=smtp.live.com
MAIL_PORT=587
MAIL_USE_SSL=false
MAIL_USE_TLS=true
MAIL_USERNAME=your-email@outlook.com
MAIL_PASSWORD=your-password
```

## API Reference

### Types

```go
type EmailConfig struct {
    SMTPHost    string
    SMTPPort    int
    UseSSL      bool
    UseTLS      bool
    Username    string
    Password    string
    DefaultFrom string
    DefaultTo   string
}

type EmailMessage struct {
    Subject  string
    Body     string
    ToAddr   []string
    CCAddr   []string
    FromAddr string
}

type EmailSender struct {
    // private fields
}
```

### Functions

```go
// Create new email sender from environment variables
func NewEmailSender() (*EmailSender, error)

// Create new email sender with custom configuration
func NewEmailSenderWithConfig(config *EmailConfig) *EmailSender

// Send email message
func (es *EmailSender) SendEmail(message *EmailMessage) error

// Send simple email to default recipient
func (es *EmailSender) SendSimpleEmail(subject, body string) error

// Send email to specific recipients
func (es *EmailSender) SendEmailToRecipients(subject, body string, toAddr []string, ccAddr []string) error

// Test SMTP connection
func (es *EmailSender) TestConnection() error

// Get configuration (password hidden)
func (es *EmailSender) GetConfig() EmailConfig
```

## Error Handling

The package provides detailed error messages for common issues:

- Missing required environment variables
- Invalid port numbers
- SMTP connection failures
- Authentication failures
- TLS/SSL connection issues
- Invalid email addresses

## Security Notes

1. **App Passwords**: For Gmail and other providers, use App Passwords instead of regular passwords
2. **Environment Variables**: Store sensitive credentials in environment variables, not in code
3. **TLS Encryption**: Always use TLS encryption for production environments
4. **Password Protection**: The `GetConfig()` method hides the password for security

## Testing

Run tests with:

```bash
go test ./pkg/email
```

Test with environment variables:

```bash
MAIL_SERVER=smtp.163.com \
MAIL_PORT=587 \
MAIL_USERNAME=test@163.com \
MAIL_PASSWORD=testpass \
go test ./pkg/email
```

## Integration with Registration System

The package includes helper functions for common use cases:

```go
// Send registration approval email
err := email.SendRegistrationApprovalEmail("user@example.com", "username")

// Send registration denial email
err := email.SendRegistrationDenialEmail("user@example.com", "username", "reason")

// Send password reset email
err := email.SendPasswordResetEmail("user@example.com", "username", "reset-token")

// Send system notification
err := email.SendSystemNotificationEmail("Alert", "System message")
```

## Troubleshooting

### Common Issues

1. **Authentication Failed**: Check username/password and enable "Less secure app access" if using Gmail
2. **Connection Timeout**: Verify SMTP server and port settings
3. **TLS Handshake Failed**: Check TLS/SSL settings for your provider
4. **Invalid Recipient**: Ensure email addresses are properly formatted

### Debug Mode

Enable debug logging to troubleshoot issues:

```go
sender, err := email.NewEmailSender()
if err != nil {
    log.Fatal(err)
}

// Test connection first
if err := sender.TestConnection(); err != nil {
    log.Printf("Connection test failed: %v", err)
    // Handle error or retry
}
```

