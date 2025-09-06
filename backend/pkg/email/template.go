package email

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// EmailTemplate represents a single email template
type EmailTemplate struct {
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}

// EmailTemplateConfig represents the entire email template configuration
type EmailTemplateConfig struct {
	Templates map[string]EmailTemplate `yaml:"templates"`
	App       AppConfig                `yaml:"app"`
	Defaults  map[string]interface{}   `yaml:"defaults"`
}

// AppConfig represents application configuration for templates
type AppConfig struct {
	Name          string `yaml:"name"`
	BaseURL       string `yaml:"base_url"`
	AdminPanelURL string `yaml:"admin_panel_url"`
	DashboardURL  string `yaml:"dashboard_url"`
	UserGuideURL  string `yaml:"user_guide_url"`
	SupportURL    string `yaml:"support_url"`
	SupportEmail  string `yaml:"support_email"`
}

// EmailTemplateManager manages email templates
type EmailTemplateManager struct {
	config    *EmailTemplateConfig
	templates map[string]*template.Template
}

// NewEmailTemplateManager creates a new email template manager
func NewEmailTemplateManager(configPath string) (*EmailTemplateManager, error) {
	// If no config path provided, use default
	if configPath == "" {
		configPath = "config/email_template.yaml"
	}

	// Read the YAML file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read email template config: %w", err)
	}

	// Parse YAML
	var config EmailTemplateConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse email template config: %w", err)
	}

	// Create template manager
	manager := &EmailTemplateManager{
		config:    &config,
		templates: make(map[string]*template.Template),
	}

	// Parse all templates
	if err := manager.parseTemplates(); err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return manager, nil
}

// parseTemplates parses all email templates
func (etm *EmailTemplateManager) parseTemplates() error {
	for name, tmpl := range etm.config.Templates {
		// Parse subject template
		subjectTmpl, err := template.New(name + "_subject").Parse(tmpl.Subject)
		if err != nil {
			return fmt.Errorf("failed to parse subject template for %s: %w", name, err)
		}

		// Parse body template
		bodyTmpl, err := template.New(name + "_body").Parse(tmpl.Body)
		if err != nil {
			return fmt.Errorf("failed to parse body template for %s: %w", name, err)
		}

		// Store templates
		etm.templates[name+"_subject"] = subjectTmpl
		etm.templates[name+"_body"] = bodyTmpl
	}

	return nil
}

// RenderTemplate renders an email template with the given data
func (etm *EmailTemplateManager) RenderTemplate(templateName string, data map[string]interface{}) (*EmailMessage, error) {
	// Merge default variables with provided data
	templateData := etm.mergeData(data)

	// Render subject
	subjectTmpl, exists := etm.templates[templateName+"_subject"]
	if !exists {
		return nil, fmt.Errorf("template %s not found", templateName)
	}

	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, templateData); err != nil {
		return nil, fmt.Errorf("failed to render subject template: %w", err)
	}

	// Render body
	bodyTmpl, exists := etm.templates[templateName+"_body"]
	if !exists {
		return nil, fmt.Errorf("body template %s not found", templateName)
	}

	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, templateData); err != nil {
		return nil, fmt.Errorf("failed to render body template: %w", err)
	}

	// Create email message
	message := &EmailMessage{
		Subject: subjectBuf.String(),
		Body:    bodyBuf.String(),
	}

	// Set recipients if provided in data
	if toAddr, ok := data["ToAddr"].([]string); ok {
		message.ToAddr = toAddr
	}
	if ccAddr, ok := data["CCAddr"].([]string); ok {
		message.CCAddr = ccAddr
	}
	if fromAddr, ok := data["FromAddr"].(string); ok {
		message.FromAddr = fromAddr
	}

	return message, nil
}

// mergeData merges default template data with provided data
func (etm *EmailTemplateManager) mergeData(data map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Add defaults
	for k, v := range etm.config.Defaults {
		merged[k] = v
	}

	// Add app config
	merged["AppName"] = etm.config.App.Name
	merged["BaseURL"] = etm.config.App.BaseURL
	merged["AdminPanelURL"] = etm.config.App.AdminPanelURL
	merged["DashboardURL"] = etm.config.App.DashboardURL
	merged["UserGuideURL"] = etm.config.App.UserGuideURL
	merged["SupportURL"] = etm.config.App.SupportURL
	merged["SupportEmail"] = etm.config.App.SupportEmail

	// Add provided data (overwrites defaults)
	for k, v := range data {
		merged[k] = v
	}

	return merged
}

// GetTemplateNames returns all available template names
func (etm *EmailTemplateManager) GetTemplateNames() []string {
	var names []string
	for name := range etm.config.Templates {
		names = append(names, name)
	}
	return names
}

// GetAppConfig returns the app configuration
func (etm *EmailTemplateManager) GetAppConfig() AppConfig {
	return etm.config.App
}

// LoadEmailTemplateManager loads the global email template manager
// This function can be called once during application startup
func LoadEmailTemplateManager() (*EmailTemplateManager, error) {
	// Try to find the config file in common locations
	possiblePaths := []string{
		"config/email_template.yaml",
		"../config/email_template.yaml",
		"./email_template.yaml",
	}

	var configPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		return nil, fmt.Errorf("email template config file not found in any of: %v", possiblePaths)
	}

	return NewEmailTemplateManager(configPath)
}

// Global template manager instance
var globalTemplateManager *EmailTemplateManager

// InitEmailTemplates initializes the global email template manager
func InitEmailTemplates() error {
	manager, err := LoadEmailTemplateManager()
	if err != nil {
		return err
	}
	globalTemplateManager = manager
	return nil
}

// GetGlobalTemplateManager returns the global template manager
func GetGlobalTemplateManager() *EmailTemplateManager {
	return globalTemplateManager
}

// RenderEmailTemplate is a convenience function to render templates using the global manager
func RenderEmailTemplate(templateName string, data map[string]interface{}) (*EmailMessage, error) {
	if globalTemplateManager == nil {
		return nil, fmt.Errorf("email template manager not initialized. Call InitEmailTemplates() first")
	}
	return globalTemplateManager.RenderTemplate(templateName, data)
}
