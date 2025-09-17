package jobs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/walterfan/lazy-rabbit-reminder/pkg/llm"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/log"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/tools"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/util"
)

// BlogWriteHandler implements JobHandler for writing blogs
type BlogWriteHandler struct {
	jobManager *JobManager
}

// PromptConfig represents the prompt configuration
type PromptConfig struct {
	// Add fields as needed for prompt configuration
}

// PromptTemplate represents a prompt template
type PromptTemplate struct {
	SystemPrompt string
	UserPrompt   string
	Description  string
}

// Execute generates daily technical blog content
func (h *BlogWriteHandler) Execute(params string) error {
	h.jobManager.logger.Info("Generating daily technical blog...")

	// Load prompt configuration
	_, err := h.loadPromptConfig()
	if err != nil {
		h.jobManager.logger.Errorf("Failed to load prompt config: %v", err)
		return fmt.Errorf("failed to load prompt config: %w", err)
	}

	// Create template data
	today := time.Now()
	data := map[string]interface{}{
		"title": fmt.Sprintf("Tech Blog - %s", today.Format("2006-01-02")),
		"today": today.Format("2006-01-02"),
		"date":  today.Format("2006-01-02"),
		"idea":  "Daily technical insights and learning",
	}

	// Use a default template for blog generation
	template := &PromptTemplate{
		SystemPrompt: "You are a technical blogger. Write engaging, informative content about technology, programming, and software development. Include practical insights and real-world examples.",
		UserPrompt:   "Write a technical blog post for {{.today}} about {{.idea}}. Make it informative and engaging for software developers. Include current weather context if available for {{.city}}.",
		Description:  "Daily Technical Blog Generator",
	}

	// Generate content using LLM
	content, err := h.generateContentWithLLM(template, data)
	if err != nil {
		h.jobManager.logger.Errorf("Failed to generate blog content: %v", err)
		return fmt.Errorf("failed to generate blog content: %w", err)
	}

	// Save the generated blog content
	filename := fmt.Sprintf("blog_%s.md", today.Format("2006-01-02"))
	if err := h.saveBlogContent(filename, content); err != nil {
		h.jobManager.logger.Errorf("Failed to save blog content: %v", err)
		return fmt.Errorf("failed to save blog content: %w", err)
	}

	h.jobManager.logger.Infof("Blog generation completed successfully: %s", filename)
	return nil
}

// loadPromptConfig loads prompt configuration
func (h *BlogWriteHandler) loadPromptConfig() (*PromptConfig, error) {
	// For now, we'll create a basic config
	return &PromptConfig{}, nil
}

// generateContentWithLLM generates content using LLM with templates
func (h *BlogWriteHandler) generateContentWithLLM(template *PromptTemplate, data map[string]interface{}) (string, error) {
	logger := log.GetLogger()

	// Create template data compatible with util.TemplateData
	templateData := util.TemplateData{}
	for k, v := range data {
		// Convert to string if needed
		switch val := v.(type) {
		case string:
			templateData[k] = val
		default:
			templateData[k] = fmt.Sprintf("%v", val)
		}
	}

	// Create template processor with function calling support
	templateProcessor := util.NewTemplateProcessor()

	// Process templates with function calling support
	renderedSystemPrompt, err := templateProcessor.ProcessTemplate(template.SystemPrompt, templateData)
	if err != nil {
		logger.Warnf("Failed to process system prompt with functions, falling back to basic rendering: %v", err)
		renderedSystemPrompt = util.RenderTemplate(template.SystemPrompt, templateData)
	}

	renderedUserPrompt, err := templateProcessor.ProcessTemplate(template.UserPrompt, templateData)
	if err != nil {
		logger.Warnf("Failed to process user prompt with functions, falling back to basic rendering: %v", err)
		renderedUserPrompt = util.RenderTemplate(template.UserPrompt, templateData)
	}

	// Check if city is provided in data for weather function calling
	city, hasCityData := data["city"].(string)
	var result string

	if strings.TrimSpace(city) != "" && hasCityData {
		// Build a minimal registry with only weather function to expose to LLM
		registry := tools.NewFunctionRegistry()
		registry.RegisterFunction(tools.GetWeatherFunction())

		functions := registry.GetFunctionDefinitionsForLLM()
		logger.Infof("Input: %s, %s, %+v", renderedSystemPrompt, renderedUserPrompt, functions)

		content, calls, errFunc := llm.AskLLMWithFunctions(renderedSystemPrompt, renderedUserPrompt, functions)
		if errFunc != nil {
			logger.Errorf("Failed during function-calling LLM request: %v", errFunc)
			return "", fmt.Errorf("failed during function-calling LLM request: %w", errFunc)
		}

		logger.Infof("Output: %s", content)
		logger.Infof("Function calls: %+v", calls)

		// If LLM requested a function call, execute locally and do a follow-up turn
		if len(calls) > 0 {
			call := calls[0]
			var args map[string]interface{}
			_ = json.Unmarshal([]byte(call.Arguments), &args)

			// Execute weather function locally
			fnResult, execErr := registry.ExecuteFunction(tools.FunctionCall{Name: call.Name, Arguments: args})
			if execErr != nil {
				logger.Errorf("Failed to execute function %s: %v", call.Name, execErr)
				return "", fmt.Errorf("failed to execute function %s: %w", call.Name, execErr)
			}

			// Provide function result back to LLM in a second turn
			resultJSON, _ := json.Marshal(fnResult)
			followUpUser := renderedUserPrompt + "\n\n" + fmt.Sprintf("Function result for %s: %s\nUse this real weather data to complete the content accurately.", call.Name, string(resultJSON))
			logger.Infof("Follow-up user prompt: %s", followUpUser)
			finalContent, errSecond := llm.AskLLM(renderedSystemPrompt, followUpUser)
			if errSecond != nil {
				logger.Errorf("Failed during follow-up LLM request: %v", errSecond)
				return "", fmt.Errorf("failed during follow-up LLM request: %w", errSecond)
			}
			result = finalContent
		} else {
			// No function call needed; use content directly
			result = content
		}
	} else {
		// No city: plain call with streaming
		var resultBuilder strings.Builder
		err = llm.AskLLMWithStream(renderedSystemPrompt, renderedUserPrompt, func(chunk string) {
			resultBuilder.WriteString(chunk)
		})
		if err != nil {
			logger.Errorf("Failed during streaming LLM request: %v", err)
			return "", fmt.Errorf("failed during streaming LLM request: %w", err)
		}
		result = resultBuilder.String()
	}

	return result, nil
}

// processTemplate processes template strings with data
func (h *BlogWriteHandler) processTemplate(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// callRealLLM calls the actual LLM package
func (h *BlogWriteHandler) callRealLLM(systemPrompt, userPrompt string) (string, error) {
	logger := log.GetLogger()

	// Use streaming LLM call for better user experience
	var resultBuilder strings.Builder
	err := llm.AskLLMWithStream(systemPrompt, userPrompt, func(chunk string) {
		resultBuilder.WriteString(chunk)
	})

	if err != nil {
		logger.Errorf("Failed during LLM request: %v", err)
		return "", fmt.Errorf("failed during LLM request: %w", err)
	}

	result := resultBuilder.String()
	logger.Infof("LLM response received, length: %d", len(result))
	return result, nil
}

// saveBlogContent saves blog content to file
func (h *BlogWriteHandler) saveBlogContent(filename, content string) error {
	// Create blogs directory if it doesn't exist
	if err := os.MkdirAll("blogs", 0755); err != nil {
		return fmt.Errorf("failed to create blogs directory: %w", err)
	}

	filepath := fmt.Sprintf("blogs/%s", filename)
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write blog file: %w", err)
	}

	return nil
}
