# Email Implementation Guide

This document describes the consolidated email implementations available in the `pkg/email` package.

## Overview

The `pkg/email` package provides two complementary approaches to sending emails:

1. **Direct Email Sender** - Simple, environment-based email sender (`email_sender.go`)
2. **Interface-based Email Service** - Flexible, composable email service (`email_interface.go`)

Both implementations are now consolidated in the same package for easier maintenance and usage.

## Implementation Comparison

### 1. Direct Email Sender

**File**: `pkg/email/email_sender.go`

**Architecture**:
- Direct struct-based implementation
- Environment variable configuration
- Comprehensive feature set
- Production-ready with extensive error handling

**Key Features**:
- ✅ Automatic environment variable loading
- ✅ SSL/TLS configuration options
- ✅ Connection testing
- ✅ Multiple SMTP provider support
- ✅ Comprehensive error handling
- ✅ Built-in HTML email support
- ✅ Async email sending integration

**Configuration**:
```go
// Automatic from environment variables
sender, err := email.NewEmailSender()

// Or with custom config
config := &email.EmailConfig{
    SMTPHost: "smtp.163.com",
    SMTPPort: 587,
    // ... other fields
}
sender := email.NewEmailSenderWithConfig(config)
```

**Usage**:
```go
// Simple email
err := sender.SendSimpleEmail("Subject", "Body")

// Email to specific recipients
err := sender.SendEmailToRecipients("Subject", "Body", toAddrs, ccAddrs)

// HTML email
err := sender.SendHTMLEmail("to@example.com", "Subject", "<h1>HTML Body</h1>")

// Full email message
message := &email.EmailMessage{
    Subject: "Test",
    Body: "Body",
    ToAddr: []string{"to@example.com"},
    CCAddr: []string{"cc@example.com"},
}
err := sender.SendEmail(message)
```

### 2. Interface-based Email Service

**File**: `pkg/email/email_interface.go`

**Architecture**:
- Interface-based design (`EmailService` interface)
- Flexible configuration options
- Clean separation of concerns
- Extensible and testable

**Key Features**:
- ✅ Interface-based design for better testing
- ✅ Multiple configuration methods
- ✅ Clean API with compose/send pattern
- ✅ Extensible architecture
- ✅ JSON-serializable structures
- ✅ Future-ready for additional features (receive, etc.)

**Configuration**:
```go
// Direct instantiation
service := email.NewEmailService(server, port, username, password, sender, receiver)

// From configuration map
config := map[string]interface{}{
    "email_server": "smtp.163.com",
    "port": 587,
    // ... other fields
}
service, err := email.NewEmailServiceFromConfig(config)
```

**Usage**:
```go
// Compose and send pattern
emailMsg := service.Compose(toAddrs, ccAddrs, "Subject", "Body")
err := service.Send(emailMsg)

// Direct sending methods (with type assertion)
if impl, ok := service.(*email.EmailServiceImpl); ok {
    err := impl.SendSimple("Subject", "Body")
    err := impl.SendWithCC("Subject", "Body", toAddrs, ccAddrs)
    err := impl.ComposeAndSend(toAddrs, ccAddrs, "Subject", "Body")
}
```

## CLI Command Usage

The `email` command supports both implementations:

### Direct Email Sender (Default)
```bash
# Test connection
./lazy-rabbit-secretary email --test

# Send email
./lazy-rabbit-secretary email --subject "Test" --body "Hello World"
```

### Interface-based Email Service
```bash
# Test connection with interface implementation
./lazy-rabbit-secretary email --test --interface

# Send email with interface implementation
./lazy-rabbit-secretary email --interface --subject "Test" --body "Hello World"
```

## When to Use Each Implementation

### Use Direct Email Sender When:
- ✅ **Production applications** - More mature and battle-tested
- ✅ **Quick integration** - Automatic environment configuration
- ✅ **Rich features needed** - HTML emails, SSL/TLS options, etc.
- ✅ **Registration system** - Already integrated with auth service
- ✅ **Minimal setup** - Works out of the box with environment variables

### Use Interface-based Email Service When:
- ✅ **Clean architecture** - Interface-based design preferred
- ✅ **Unit testing** - Easier to mock and test
- ✅ **Custom configuration** - Need flexible configuration options
- ✅ **Future extensions** - Planning to add receive functionality
- ✅ **Microservices** - Better for dependency injection patterns

## Environment Variables

Both implementations use the same environment variables:

```bash
# Required
MAIL_SERVER=smtp.163.com
MAIL_PORT=587
MAIL_USERNAME=fantasyamin@163.com
MAIL_PASSWORD=your-password

# Optional
MAIL_USE_SSL=false
MAIL_USE_TLS=true
MAIL_SENDER=fantasyamin@163.com
MAIL_RECEIVER=walter.fan@zoom.us
```

## Integration Examples

### Registration System Integration

**Direct Email Sender** (Currently Used):
```go
// In auth/service.go
func (a *AuthService) sendRegistrationStatusEmail(user *models.User, approved bool, reason string) error {
    sender, err := email.NewEmailSender()
    if err != nil {
        return err
    }
    
    message := &email.EmailMessage{
        Subject: subject,
        Body:    body,
        ToAddr:  []string{user.Email},
    }
    
    return sender.SendEmail(message)
}
```

**Interface-based Email Service** (Alternative):
```go
// Alternative implementation
func (a *AuthService) sendRegistrationStatusEmailInterface(user *models.User, approved bool, reason string) error {
    emailService, err := createEmailServiceFromEnv()
    if err != nil {
        return err
    }
    
    emailMsg := emailService.Compose(
        []string{user.Email},
        nil,
        subject,
        body,
    )
    
    return emailService.Send(emailMsg)
}

func createEmailServiceFromEnv() (email.EmailService, error) {
    config := map[string]interface{}{
        "email_server": os.Getenv("MAIL_SERVER"),
        "port":         getEnvIntOrDefault("MAIL_PORT", 587),
        "username":     os.Getenv("MAIL_USERNAME"),
        "password":     os.Getenv("MAIL_PASSWORD"),
        "sender":       os.Getenv("MAIL_SENDER"),
        "receiver":     os.Getenv("MAIL_RECEIVER"),
    }
    return email.NewEmailServiceFromConfig(config)
}
```

## Testing

### Direct Email Sender Testing:
```go
func TestDirectEmailSender(t *testing.T) {
    // Set environment variables
    os.Setenv("MAIL_SERVER", "smtp.test.com")
    os.Setenv("MAIL_PORT", "587")
    // ... other env vars
    
    sender, err := email.NewEmailSender()
    if err != nil {
        t.Fatalf("Failed to create sender: %v", err)
    }
    
    // Test functionality
    err = sender.SendSimpleEmail("Test Subject", "Test Body")
    // ... assertions
}
```

### Interface-based Email Service Testing:
```go
func TestInterfaceEmailService(t *testing.T) {
    // Mock the interface
    mockService := &MockEmailService{}
    
    // Test with mock
    emailMsg := mockService.Compose([]string{"test@example.com"}, nil, "Subject", "Body")
    err := mockService.Send(emailMsg)
    // ... assertions
}

type MockEmailService struct{}

func (m *MockEmailService) Send(emailMsg *email.Email) error {
    // Mock implementation
    return nil
}

func (m *MockEmailService) Compose(to, cc []string, subject, body string) *email.Email {
    return &email.Email{ToAddr: to, CCAddr: cc, Subject: subject, Body: body}
}

func (m *MockEmailService) Check() error { return nil }
func (m *MockEmailService) Receive() ([]*email.Email, error) { return nil, nil }
```

## Examples and Integration

The package includes comprehensive examples in:
- `email_interface_example.go` - Basic usage examples for the interface-based service
- `integration_examples.go` - Complete integration examples showing both implementations

To run the examples:
```go
package main

import "github.com/walterfan/lazy-rabbit-secretary/pkg/email"

func main() {
    // Run all integration examples
    email.RunAllExamples()
}
```

## Performance Considerations

### Direct Email Sender:
- **Memory**: Slightly higher due to comprehensive configuration
- **Startup**: Faster initialization with environment variables
- **Features**: More features = slightly larger binary

### Interface-based Email Service:
- **Memory**: Lower memory footprint
- **Startup**: Minimal initialization overhead
- **Features**: Lean implementation, smaller binary contribution

## Migration Guide

### Switching Between Implementations:

**From Direct to Interface-based:**
```go
// Before (Direct)
sender, err := email.NewEmailSender()
err = sender.SendSimpleEmail("Subject", "Body")

// After (Interface-based)
service, err := email.NewEmailServiceFromConfig(config)
if impl, ok := service.(*email.EmailServiceImpl); ok {
    err = impl.SendSimple("Subject", "Body")
}
```

**From Interface-based to Direct:**
```go
// Before (Interface-based)
service, err := email.NewEmailServiceFromConfig(config)
emailMsg := service.Compose(to, cc, subject, body)
err = service.Send(emailMsg)

// After (Direct)
sender, err := email.NewEmailSender()
message := &email.EmailMessage{
    ToAddr: to, CCAddr: cc, Subject: subject, Body: body,
}
err = sender.SendEmail(message)
```

## Recommendation

For **new projects** and **production use**, we recommend the **Direct Email Sender** because:

1. **Battle-tested**: Already integrated and tested in the registration system
2. **Feature-complete**: Includes SSL/TLS, connection testing, HTML support
3. **Easy setup**: Works immediately with environment variables
4. **Production-ready**: Comprehensive error handling and logging

For **testing**, **clean architecture**, or **future extensibility**, consider the **Interface-based Email Service** for its interface design and flexibility.

Both implementations are maintained and supported in the same package. Choose based on your specific needs and architectural preferences.