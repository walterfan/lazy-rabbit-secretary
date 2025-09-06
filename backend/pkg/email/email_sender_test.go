package email

import (
	"os"
	"testing"
)

func TestNewEmailSender(t *testing.T) {
	// Set up test environment variables
	os.Setenv("MAIL_SERVER", "smtp.163.com")
	os.Setenv("MAIL_PORT", "587")
	os.Setenv("MAIL_USE_SSL", "false")
	os.Setenv("MAIL_USE_TLS", "true")
	os.Setenv("MAIL_USERNAME", "test@163.com")
	os.Setenv("MAIL_PASSWORD", "testpass")
	os.Setenv("MAIL_SENDER", "test@163.com")
	os.Setenv("MAIL_RECEIVER", "receiver@example.com")

	sender, err := NewEmailSender()
	if err != nil {
		t.Fatalf("Failed to create email sender: %v", err)
	}

	config := sender.GetConfig()
	if config.SMTPHost != "smtp.163.com" {
		t.Errorf("Expected SMTP host 'smtp.163.com', got '%s'", config.SMTPHost)
	}

	if config.SMTPPort != 587 {
		t.Errorf("Expected SMTP port 587, got %d", config.SMTPPort)
	}

	if config.UseSSL != false {
		t.Errorf("Expected UseSSL false, got %t", config.UseSSL)
	}

	if config.UseTLS != true {
		t.Errorf("Expected UseTLS true, got %t", config.UseTLS)
	}

	// Clean up
	os.Unsetenv("MAIL_SERVER")
	os.Unsetenv("MAIL_PORT")
	os.Unsetenv("MAIL_USE_SSL")
	os.Unsetenv("MAIL_USE_TLS")
	os.Unsetenv("MAIL_USERNAME")
	os.Unsetenv("MAIL_PASSWORD")
	os.Unsetenv("MAIL_SENDER")
	os.Unsetenv("MAIL_RECEIVER")
}

func TestEmailMessage(t *testing.T) {
	config := &EmailConfig{
		SMTPHost:    "smtp.163.com",
		SMTPPort:    587,
		UseSSL:      false,
		UseTLS:      true,
		Username:    "test@163.com",
		Password:    "testpass",
		DefaultFrom: "test@163.com",
		DefaultTo:   "receiver@example.com",
	}

	sender := NewEmailSenderWithConfig(config)

	// Test building email content
	message := &EmailMessage{
		Subject:  "Test Subject",
		Body:     "Test Body",
		ToAddr:   []string{"to@example.com"},
		CCAddr:   []string{"cc@example.com"},
		FromAddr: "from@example.com",
	}

	content := sender.buildEmailContent(message)
	contentStr := string(content)

	if !contains(contentStr, "Subject: Test Subject") {
		t.Error("Email content should contain subject")
	}

	if !contains(contentStr, "To: to@example.com") {
		t.Error("Email content should contain To address")
	}

	if !contains(contentStr, "Cc: cc@example.com") {
		t.Error("Email content should contain CC address")
	}

	if !contains(contentStr, "From: from@example.com") {
		t.Error("Email content should contain From address")
	}

	if !contains(contentStr, "Test Body") {
		t.Error("Email content should contain body")
	}
}

func TestLoadConfigFromEnv_MissingRequired(t *testing.T) {
	// Clear all environment variables
	os.Unsetenv("MAIL_SERVER")
	os.Unsetenv("MAIL_USERNAME")
	os.Unsetenv("MAIL_PASSWORD")

	_, err := loadConfigFromEnv()
	if err == nil {
		t.Error("Expected error when required environment variables are missing")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

