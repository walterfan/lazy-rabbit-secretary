package email

import "fmt"

// Email represents an email message for the interface-based service
type Email struct {
	ToAddr  []string `json:"to_addr"`
	CCAddr  []string `json:"cc_addr"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// EmailService defines the interface for email operations
type EmailService interface {
	// Send sends an email
	Send(email *Email) error

	// Receive receives emails (placeholder for future implementation)
	Receive() ([]*Email, error)

	// Compose creates a new email with the given parameters
	Compose(toAddr []string, ccAddr []string, subject string, body string) *Email

	// Check tests the email service connection
	Check() error
}

// EmailServiceImpl implements the EmailService interface using EmailSender
type EmailServiceImpl struct {
	sender   *EmailSender
	receiver string // default receiver for the service
}

// NewEmailService creates a new EmailServiceImpl instance
func NewEmailService(emailServer string, port int, username, password, sender, receiver string) EmailService {
	config := &EmailConfig{
		SMTPHost:    emailServer,
		SMTPPort:    port,
		UseTLS:      true, // default to TLS
		Username:    username,
		Password:    password,
		DefaultFrom: sender,
		DefaultTo:   receiver,
	}

	emailSender := NewEmailSenderWithConfig(config)

	return &EmailServiceImpl{
		sender:   emailSender,
		receiver: receiver,
	}
}

// NewEmailServiceFromConfig creates a new EmailServiceImpl from a configuration map
func NewEmailServiceFromConfig(config map[string]interface{}) (EmailService, error) {
	emailServer, ok := config["email_server"].(string)
	if !ok || emailServer == "" {
		return nil, fmt.Errorf("email_server is required")
	}

	port, ok := config["port"].(int)
	if !ok || port <= 0 {
		port = 587 // default SMTP port
	}

	username, ok := config["username"].(string)
	if !ok || username == "" {
		return nil, fmt.Errorf("username is required")
	}

	password, ok := config["password"].(string)
	if !ok || password == "" {
		return nil, fmt.Errorf("password is required")
	}

	sender, ok := config["sender"].(string)
	if !ok || sender == "" {
		sender = username // default to username
	}

	receiver, ok := config["receiver"].(string)
	if !ok {
		receiver = "" // optional
	}

	return NewEmailService(emailServer, port, username, password, sender, receiver), nil
}

// Send implements EmailService.Send
func (es *EmailServiceImpl) Send(email *Email) error {
	if email == nil {
		return fmt.Errorf("email cannot be nil")
	}

	// Convert Email to EmailMessage
	message := &EmailMessage{
		Subject: email.Subject,
		Body:    email.Body,
		ToAddr:  email.ToAddr,
		CCAddr:  email.CCAddr,
	}

	// Set default recipients if none provided
	if len(message.ToAddr) == 0 && es.receiver != "" {
		message.ToAddr = []string{es.receiver}
	}

	return es.sender.SendEmail(message)
}

// Receive implements EmailService.Receive
// Note: This is a placeholder implementation for future POP3/IMAP support
func (es *EmailServiceImpl) Receive() ([]*Email, error) {
	// TODO: Implement email receiving functionality using POP3 or IMAP
	return nil, fmt.Errorf("receive functionality not yet implemented")
}

// Compose implements EmailService.Compose
func (es *EmailServiceImpl) Compose(toAddr []string, ccAddr []string, subject string, body string) *Email {
	return &Email{
		ToAddr:  toAddr,
		CCAddr:  ccAddr,
		Subject: subject,
		Body:    body,
	}
}

// Check implements EmailService.Check
func (es *EmailServiceImpl) Check() error {
	return es.sender.TestConnection()
}

// GetConfig returns the email service configuration (without sensitive data)
func (es *EmailServiceImpl) GetConfig() map[string]interface{} {
	config := es.sender.GetConfig()
	return map[string]interface{}{
		"email_server": config.SMTPHost,
		"port":         config.SMTPPort,
		"username":     config.Username,
		"password":     "***", // Hide password
		"sender":       config.DefaultFrom,
		"receiver":     es.receiver,
	}
}

// SendSimple sends a simple email with subject and body to default or specified recipients
func (es *EmailServiceImpl) SendSimple(subject, body string, toAddr ...string) error {
	email := &Email{
		Subject: subject,
		Body:    body,
	}

	if len(toAddr) > 0 {
		email.ToAddr = toAddr
	}

	return es.Send(email)
}

// SendWithCC sends an email with CC recipients
func (es *EmailServiceImpl) SendWithCC(subject, body string, toAddr []string, ccAddr []string) error {
	email := &Email{
		ToAddr:  toAddr,
		CCAddr:  ccAddr,
		Subject: subject,
		Body:    body,
	}

	return es.Send(email)
}

// ComposeAndSend composes and immediately sends an email
func (es *EmailServiceImpl) ComposeAndSend(toAddr []string, ccAddr []string, subject string, body string) error {
	email := es.Compose(toAddr, ccAddr, subject, body)
	return es.Send(email)
}
