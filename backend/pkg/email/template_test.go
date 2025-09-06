package email

import (
	"os"
	"testing"
	"time"
)

func TestEmailTemplateManager(t *testing.T) {
	// Create a temporary template file for testing
	tempConfig := `
templates:
  test_template:
    subject: "Test Subject for {{.Username}}"
    body: |
      Hello {{.Username}},
      
      This is a test email.
      
      App: {{.AppName}}
      URL: {{.TestURL}}

app:
  name: "Test App"
  base_url: "http://test.com"

defaults:
  test_default: "default_value"
`

	// Write temp file
	tempFile := "test_email_template.yaml"
	err := os.WriteFile(tempFile, []byte(tempConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}
	defer os.Remove(tempFile)

	// Create template manager
	manager, err := NewEmailTemplateManager(tempFile)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Test template rendering
	data := map[string]interface{}{
		"Username": "testuser",
		"TestURL":  "http://example.com/test",
	}

	message, err := manager.RenderTemplate("test_template", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Verify subject
	expectedSubject := "Test Subject for testuser"
	if message.Subject != expectedSubject {
		t.Errorf("Expected subject '%s', got '%s'", expectedSubject, message.Subject)
	}

	// Verify body contains expected content
	if !contains(message.Body, "Hello testuser") {
		t.Errorf("Body should contain 'Hello testuser', got: %s", message.Body)
	}

	if !contains(message.Body, "App: Test App") {
		t.Errorf("Body should contain app name, got: %s", message.Body)
	}

	if !contains(message.Body, "URL: http://example.com/test") {
		t.Errorf("Body should contain test URL, got: %s", message.Body)
	}
}

func TestEmailConfirmationTemplate(t *testing.T) {
	// Test with the actual config file if it exists
	configPath := "../config/email_template.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "../../config/email_template.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("Email template config not found, skipping test")
	}

	manager, err := NewEmailTemplateManager(configPath)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Test email confirmation template
	data := map[string]interface{}{
		"Username":   "john_doe",
		"Email":      "john@example.com",
		"ConfirmURL": "http://localhost:8080/api/v1/auth/confirm?token=abc123",
		"ExpiryTime": time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05 MST"),
		"ToAddr":     []string{"john@example.com"},
	}

	message, err := manager.RenderTemplate("email_confirmation", data)
	if err != nil {
		t.Fatalf("Failed to render email confirmation template: %v", err)
	}

	// Verify subject
	if message.Subject == "" {
		t.Error("Subject should not be empty")
	}

	// Verify body contains expected content
	if !contains(message.Body, "john_doe") {
		t.Errorf("Body should contain username, got: %s", message.Body)
	}

	if !contains(message.Body, "http://localhost:8080/api/v1/auth/confirm?token=abc123") {
		t.Errorf("Body should contain confirmation URL, got: %s", message.Body)
	}

	if !contains(message.Body, "Lazy Rabbit Reminder") {
		t.Errorf("Body should contain app name, got: %s", message.Body)
	}

	t.Logf("Rendered subject: %s", message.Subject)
	t.Logf("Rendered body: %s", message.Body)
}

func TestRegistrationApprovalTemplate(t *testing.T) {
	// Test with the actual config file if it exists
	configPath := "../config/email_template.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "../../config/email_template.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("Email template config not found, skipping test")
	}

	manager, err := NewEmailTemplateManager(configPath)
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Test approval template
	approvalData := map[string]interface{}{
		"Username": "jane_doe",
		"Email":    "jane@example.com",
		"Approved": true,
		"Status":   "Approved",
		"ToAddr":   []string{"jane@example.com"},
	}

	message, err := manager.RenderTemplate("registration_approval", approvalData)
	if err != nil {
		t.Fatalf("Failed to render registration approval template: %v", err)
	}

	if !contains(message.Body, "Congratulations") {
		t.Errorf("Approval body should contain 'Congratulations', got: %s", message.Body)
	}

	// Test denial template
	denialData := map[string]interface{}{
		"Username": "bob_smith",
		"Email":    "bob@example.com",
		"Approved": false,
		"Status":   "Denied",
		"Reason":   "Invalid email domain",
		"ToAddr":   []string{"bob@example.com"},
	}

	message, err = manager.RenderTemplate("registration_approval", denialData)
	if err != nil {
		t.Fatalf("Failed to render registration denial template: %v", err)
	}

	if !contains(message.Body, "regret to inform") {
		t.Errorf("Denial body should contain 'regret to inform', got: %s", message.Body)
	}

	if !contains(message.Body, "Invalid email domain") {
		t.Errorf("Denial body should contain reason, got: %s", message.Body)
	}

	t.Logf("Denial subject: %s", message.Subject)
	t.Logf("Denial body: %s", message.Body)
}

// Note: Helper functions contains and containsSubstring are defined in email_sender_test.go
