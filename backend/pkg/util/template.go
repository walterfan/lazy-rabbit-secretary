package util

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// PromptTemplate represents a template configuration
type PromptTemplate struct {
	Description  string `yaml:"description"`
	SystemPrompt string `yaml:"system_prompt"`
	UserPrompt   string `yaml:"user_prompt"`
	Tags         string `yaml:"tags"`
}

// TemplateData represents template variables
type TemplateData map[string]string

// PromptConfig represents the entire prompts configuration
type PromptConfig struct {
	Prompts map[string]PromptTemplate `yaml:"prompts"`
}

// LoadPromptTemplates loads prompt templates from YAML file
func LoadPromptTemplates(filepath string) (*PromptConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompts file: %w", err)
	}

	var config PromptConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompts YAML: %w", err)
	}

	return &config, nil
}

// GetPromptTemplate retrieves a specific prompt template
func (pc *PromptConfig) GetPromptTemplate(name string) (PromptTemplate, error) {
	template, exists := pc.Prompts[name]
	if !exists {
		return PromptTemplate{}, fmt.Errorf("prompt template '%s' not found", name)
	}
	return template, nil
}

// RenderTemplate renders a template string with provided data
func RenderTemplate(template string, data TemplateData) string {
	result := template
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// RenderPromptTemplate renders both system and user prompts with data
func (pt *PromptTemplate) RenderPromptTemplate(data TemplateData) (string, string) {
	systemPrompt := RenderTemplate(pt.SystemPrompt, data)
	userPrompt := RenderTemplate(pt.UserPrompt, data)
	return systemPrompt, userPrompt
}

// ListAvailableTemplates returns a list of available template names
func (pc *PromptConfig) ListAvailableTemplates() []string {
	var names []string
	for name := range pc.Prompts {
		names = append(names, name)
	}
	return names
}

// GetTemplatesByTag returns templates that have the specified tag
func (pc *PromptConfig) GetTemplatesByTag(tag string) map[string]PromptTemplate {
	result := make(map[string]PromptTemplate)
	for name, template := range pc.Prompts {
		if strings.Contains(template.Tags, tag) {
			result[name] = template
		}
	}
	return result
}
