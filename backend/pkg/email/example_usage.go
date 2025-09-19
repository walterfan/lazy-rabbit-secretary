package email

import (
	"fmt"

	"github.com/walterfan/lazy-rabbit-secretary/pkg/log"
)

// ExampleUsage demonstrates how to use the EmailSender
func ExampleUsage() {
	// Method 1: Create email sender from environment variables
	sender, err := NewEmailSender()
	if err != nil {
		log.GetLogger().Fatalf("Failed to create email sender: %v", err)
	}

	// Test connection first
	if err := sender.TestConnection(); err != nil {
		log.GetLogger().Errorf("Email connection test failed: %v", err)
		return
	}
	log.GetLogger().Info("Email connection test successful!")

	// Example 1: Send simple email to default recipient
	err = sender.SendSimpleEmail(
		"Test Subject",
		"This is a test email body.",
	)
	if err != nil {
		log.GetLogger().Errorf("Failed to send simple email: %v", err)
	} else {
		log.GetLogger().Info("Simple email sent successfully!")
	}

	// Example 2: Send email to specific recipients
	err = sender.SendEmailToRecipients(
		"Important Notification",
		"This is an important notification email.\n\nBest regards,\nSystem Admin",
		[]string{"user1@example.com", "user2@example.com"}, // To addresses
		[]string{"manager@example.com"},                    // CC addresses
	)
	if err != nil {
		log.GetLogger().Errorf("Failed to send email to recipients: %v", err)
	} else {
		log.GetLogger().Info("Email sent to specific recipients successfully!")
	}

	// Example 3: Send custom email message
	message := &EmailMessage{
		Subject:  "Registration Approval",
		Body:     "Your registration has been approved. Welcome to our system!",
		ToAddr:   []string{"newuser@example.com"},
		CCAddr:   []string{"admin@example.com"},
		FromAddr: "noreply@yourcompany.com",
	}

	err = sender.SendEmail(message)
	if err != nil {
		log.GetLogger().Errorf("Failed to send custom email: %v", err)
	} else {
		log.GetLogger().Info("Custom email sent successfully!")
	}

	// Method 2: Create email sender with custom configuration
	customConfig := &EmailConfig{
		SMTPHost:    "smtp.gmail.com",
		SMTPPort:    587,
		UseSSL:      false,
		UseTLS:      true,
		Username:    "your-email@gmail.com",
		Password:    "your-app-password",
		DefaultFrom: "your-email@gmail.com",
		DefaultTo:   "default-recipient@example.com",
	}

	customSender := NewEmailSenderWithConfig(customConfig)

	// Use custom sender
	err = customSender.SendSimpleEmail(
		"Custom Config Email",
		"This email was sent using custom configuration.",
	)
	if err != nil {
		log.GetLogger().Errorf("Failed to send email with custom config: %v", err)
	} else {
		log.GetLogger().Info("Email with custom config sent successfully!")
	}
}

// SendRegistrationApprovalEmail sends an email when a user registration is approved
func SendRegistrationApprovalEmail(userEmail, username string) error {
	sender, err := NewEmailSender()
	if err != nil {
		return fmt.Errorf("failed to create email sender: %w", err)
	}

	subject := "Registration Approved - Welcome!"
	body := fmt.Sprintf(`Dear %s,

Congratulations! Your registration has been approved.

You can now log in to the system using your credentials.

Best regards,
The Admin Team`, username)

	message := &EmailMessage{
		Subject: subject,
		Body:    body,
		ToAddr:  []string{userEmail},
	}

	return sender.SendEmail(message)
}

// SendRegistrationDenialEmail sends an email when a user registration is denied
func SendRegistrationDenialEmail(userEmail, username, reason string) error {
	sender, err := NewEmailSender()
	if err != nil {
		return fmt.Errorf("failed to create email sender: %w", err)
	}

	subject := "Registration Status Update"
	body := fmt.Sprintf(`Dear %s,

We regret to inform you that your registration has been denied.

Reason: %s

If you have any questions, please contact our support team.

Best regards,
The Admin Team`, username, reason)

	message := &EmailMessage{
		Subject: subject,
		Body:    body,
		ToAddr:  []string{userEmail},
	}

	return sender.SendEmail(message)
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(userEmail, username, resetToken string) error {
	sender, err := NewEmailSender()
	if err != nil {
		return fmt.Errorf("failed to create email sender: %w", err)
	}

	subject := "Password Reset Request"
	body := fmt.Sprintf(`Dear %s,

You have requested a password reset for your account.

Please use the following token to reset your password:
%s

This token will expire in 24 hours.

If you did not request this password reset, please ignore this email.

Best regards,
The Admin Team`, username, resetToken)

	message := &EmailMessage{
		Subject: subject,
		Body:    body,
		ToAddr:  []string{userEmail},
	}

	return sender.SendEmail(message)
}

// SendSystemNotificationEmail sends a system notification to administrators
func SendSystemNotificationEmail(title, message string) error {
	sender, err := NewEmailSender()
	if err != nil {
		return fmt.Errorf("failed to create email sender: %w", err)
	}

	subject := fmt.Sprintf("System Notification: %s", title)
	body := fmt.Sprintf(`System Notification

Title: %s

Message:
%s

Timestamp: %s

This is an automated message from the system.`, title, message, "timestamp_placeholder")

	// Send to default admin email (configured in MAIL_RECEIVER)
	return sender.SendSimpleEmail(subject, body)
}
