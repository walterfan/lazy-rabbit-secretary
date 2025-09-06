package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// EmailConfig holds email configuration
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

// EmailMessage represents an email message
type EmailMessage struct {
	Subject  string
	Body     string
	ToAddr   []string
	CCAddr   []string
	FromAddr string
}

// EmailSender handles sending emails
type EmailSender struct {
	config *EmailConfig
}

// NewEmailSender creates a new email sender with configuration from environment variables
func NewEmailSender() (*EmailSender, error) {
	config, err := loadConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load email config: %w", err)
	}

	return &EmailSender{
		config: config,
	}, nil
}

// NewEmailSenderWithConfig creates a new email sender with provided configuration
func NewEmailSenderWithConfig(config *EmailConfig) *EmailSender {
	return &EmailSender{
		config: config,
	}
}

// loadConfigFromEnv loads email configuration from environment variables
func loadConfigFromEnv() (*EmailConfig, error) {
	// Required fields
	smtpHost := os.Getenv("MAIL_SERVER")
	if smtpHost == "" {
		return nil, fmt.Errorf("MAIL_SERVER environment variable is required")
	}

	portStr := os.Getenv("MAIL_PORT")
	if portStr == "" {
		portStr = "587" // Default SMTP port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MAIL_PORT: %s", portStr)
	}

	username := os.Getenv("MAIL_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("MAIL_USERNAME environment variable is required")
	}

	password := os.Getenv("MAIL_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("MAIL_PASSWORD environment variable is required")
	}

	// Optional fields with defaults
	useSSLStr := os.Getenv("MAIL_USE_SSL")
	useSSL, _ := strconv.ParseBool(useSSLStr) // defaults to false if not set or invalid

	useTLSStr := os.Getenv("MAIL_USE_TLS")
	useTLS := true // default to true
	if useTLSStr != "" {
		useTLS, _ = strconv.ParseBool(useTLSStr)
	}

	defaultFrom := os.Getenv("MAIL_SENDER")
	if defaultFrom == "" {
		defaultFrom = username // Use username as default sender
	}

	defaultTo := os.Getenv("MAIL_RECEIVER")

	return &EmailConfig{
		SMTPHost:    smtpHost,
		SMTPPort:    port,
		UseSSL:      useSSL,
		UseTLS:      useTLS,
		Username:    username,
		Password:    password,
		DefaultFrom: defaultFrom,
		DefaultTo:   defaultTo,
	}, nil
}

// SendEmail sends an email message
func (es *EmailSender) SendEmail(message *EmailMessage) error {
	// Validate message
	if message.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if message.Body == "" {
		return fmt.Errorf("email body is required")
	}
	if len(message.ToAddr) == 0 {
		if es.config.DefaultTo == "" {
			return fmt.Errorf("no recipients specified and no default receiver configured")
		}
		message.ToAddr = []string{es.config.DefaultTo}
	}

	// Set default from address if not specified
	if message.FromAddr == "" {
		message.FromAddr = es.config.DefaultFrom
	}

	// Build recipient list
	var recipients []string
	recipients = append(recipients, message.ToAddr...)
	recipients = append(recipients, message.CCAddr...)

	// Create email content
	emailContent := es.buildEmailContent(message)

	// Send email
	return es.sendSMTP(message.FromAddr, recipients, emailContent)
}

// SendSimpleEmail sends a simple email with subject and body to default recipient
func (es *EmailSender) SendSimpleEmail(subject, body string) error {
	message := &EmailMessage{
		Subject: subject,
		Body:    body,
	}
	return es.SendEmail(message)
}

// SendEmailToRecipients sends an email to specific recipients
func (es *EmailSender) SendEmailToRecipients(subject, body string, toAddr []string, ccAddr []string) error {
	message := &EmailMessage{
		Subject: subject,
		Body:    body,
		ToAddr:  toAddr,
		CCAddr:  ccAddr,
	}
	return es.SendEmail(message)
}

// buildEmailContent creates the email content with proper headers
func (es *EmailSender) buildEmailContent(message *EmailMessage) []byte {
	var content strings.Builder

	// Headers
	content.WriteString(fmt.Sprintf("From: %s\r\n", message.FromAddr))
	content.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.ToAddr, ", ")))

	if len(message.CCAddr) > 0 {
		content.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(message.CCAddr, ", ")))
	}

	content.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))
	content.WriteString("MIME-Version: 1.0\r\n")
	content.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	content.WriteString("\r\n")

	// Body
	content.WriteString(message.Body)

	return []byte(content.String())
}

// sendSMTP sends email via SMTP
func (es *EmailSender) sendSMTP(from string, to []string, content []byte) error {
	addr := fmt.Sprintf("%s:%d", es.config.SMTPHost, es.config.SMTPPort)

	// Create SMTP client
	var client *smtp.Client
	var err error

	if es.config.UseSSL {
		// SSL connection
		tlsConfig := &tls.Config{
			ServerName: es.config.SMTPHost,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect with SSL: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, es.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
	} else {
		// Plain connection (with optional TLS upgrade)
		client, err = smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
	}
	defer client.Close()

	// Start TLS if required and not using SSL
	if es.config.UseTLS && !es.config.UseSSL {
		tlsConfig := &tls.Config{
			ServerName: es.config.SMTPHost,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	if es.config.Username != "" && es.config.Password != "" {
		auth := smtp.PlainAuth("", es.config.Username, es.config.Password, es.config.SMTPHost)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send email content
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	if _, err = writer.Write(content); err != nil {
		return fmt.Errorf("failed to write email content: %w", err)
	}

	return nil
}

// GetConfig returns the current email configuration (without sensitive data)
func (es *EmailSender) GetConfig() EmailConfig {
	config := *es.config
	config.Password = "***" // Hide password
	return config
}

// TestConnection tests the SMTP connection without sending an email
func (es *EmailSender) TestConnection() error {
	addr := fmt.Sprintf("%s:%d", es.config.SMTPHost, es.config.SMTPPort)

	var client *smtp.Client
	var err error

	if es.config.UseSSL {
		tlsConfig := &tls.Config{
			ServerName: es.config.SMTPHost,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("SSL connection failed: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, es.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("SMTP client creation failed: %w", err)
		}
	} else {
		client, err = smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("SMTP connection failed: %w", err)
		}
	}
	defer client.Close()

	// Test TLS if enabled
	if es.config.UseTLS && !es.config.UseSSL {
		tlsConfig := &tls.Config{
			ServerName: es.config.SMTPHost,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("TLS upgrade failed: %w", err)
		}
	}

	// Test authentication
	if es.config.Username != "" && es.config.Password != "" {
		auth := smtp.PlainAuth("", es.config.Username, es.config.Password, es.config.SMTPHost)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	return nil
}
