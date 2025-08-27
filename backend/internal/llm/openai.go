package llm

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/walterfan/lazy-rabbit-reminder/internal/log"
)

// LLMSettings represents configuration for LLM requests
type LLMSettings struct {
	BaseUrl     string  `json:"base_url"`
	ApiKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
}

// ChatRequest represents the request payload for chat completions
type ChatRequest struct {
	Model        string                   `json:"model"`
	Messages     []ChatEntry              `json:"messages"`
	Stream       bool                     `json:"stream"`
	Temperature  float64                  `json:"temperature,omitempty"`
	Functions    []map[string]interface{} `json:"functions,omitempty"`
	FunctionCall interface{}              `json:"function_call,omitempty"`
}

// ChatEntry represents a message in the conversation
type ChatEntry struct {
	Role         string        `json:"role"`
	Content      string        `json:"content,omitempty"`
	Name         string        `json:"name,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

// FunctionCall represents a function call request
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatResponse represents the response from chat completions API
type ChatResponse struct {
	Choices []struct {
		Message ChatEntry `json:"message"`
	} `json:"choices"`
}

// ChatMessage represents a message in conversation history
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMClient handles all LLM interactions
type LLMClient struct {
	httpClient *http.Client
	logger     interface {
		Infof(string, ...interface{})
		Errorf(string, ...interface{})
		Warnf(string, ...interface{})
	}
	settings LLMSettings
}

// loadSettingsFromEnv loads settings from environment variables
func loadSettingsFromEnv() LLMSettings {
	settings := LLMSettings{}
	settings.BaseUrl = os.Getenv("LLM_BASE_URL")
	settings.ApiKey = os.Getenv("LLM_API_KEY")
	settings.Model = os.Getenv("LLM_MODEL")
	if tempStr := os.Getenv("LLM_TEMPERATURE"); tempStr != "" {
		if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
			settings.Temperature = temp
		}
	}
	if settings.Temperature == 0 {
		settings.Temperature = 1.0
	}
	return settings
}

// NewLLMClient creates a new LLM client. If no settings are provided, it loads from environment variables.
func NewLLMClient(optionalSettings ...LLMSettings) *LLMClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resolved := LLMSettings{}
	if len(optionalSettings) > 0 {
		resolved = optionalSettings[0]
		if resolved.Temperature == 0 {
			resolved.Temperature = 1.0
		}
	} else {
		resolved = loadSettingsFromEnv()
	}

	return &LLMClient{
		httpClient: &http.Client{Transport: transport},
		logger:     log.GetLogger(),
		settings:   resolved,
	}
}

// resolveSettings merges provided settings with environment variables
func (c *LLMClient) resolveSettings(settings LLMSettings) LLMSettings {
	resolved := settings

	// Use environment variables as fallback
	if resolved.BaseUrl == "" {
		resolved.BaseUrl = os.Getenv("LLM_BASE_URL")
	}
	if resolved.ApiKey == "" {
		resolved.ApiKey = os.Getenv("LLM_API_KEY")
	}
	if resolved.Model == "" {
		resolved.Model = os.Getenv("LLM_MODEL")
	}
	if resolved.Temperature == 0 {
		if tempStr := os.Getenv("LLM_TEMPERATURE"); tempStr != "" {
			if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
				resolved.Temperature = temp
			}
		}
		if resolved.Temperature == 0 {
			resolved.Temperature = 1.0
		}
	}

	return resolved
}

// createRequest creates an HTTP request for the LLM API
func (c *LLMClient) createRequest(req *ChatRequest) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", c.settings.BaseUrl), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.settings.ApiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	return httpReq, nil
}

// executeRequest executes an HTTP request and returns the response
func (c *LLMClient) executeRequest(httpReq *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.logger.Errorf("LLM request failed with status code: %d, Response Body: %s", resp.StatusCode, bodyBytes)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

// buildMessages creates a message array from system prompt, user prompt, and optional history
func buildMessages(systemPrompt, userPrompt string, history []ChatMessage) []ChatEntry {
	messages := []ChatEntry{{Role: "system", Content: systemPrompt}}

	// Add conversation history
	for _, msg := range history {
		messages = append(messages, ChatEntry{Role: msg.Role, Content: msg.Content})
	}

	// Add current user prompt
	messages = append(messages, ChatEntry{Role: "user", Content: userPrompt})

	return messages
}

// AskLLM sends a simple chat completion request
func (c *LLMClient) AskLLM(systemPrompt, userPrompt string) (string, error) {
	c.logger.Infof("Using LLM settings - BaseUrl: %s, Model: %s, Temperature: %.1f", c.settings.BaseUrl, c.settings.Model, c.settings.Temperature)

	req := &ChatRequest{
		Model:       c.settings.Model,
		Messages:    buildMessages(systemPrompt, userPrompt, nil),
		Stream:      false,
		Temperature: c.settings.Temperature,
	}

	httpReq, err := c.createRequest(req)
	if err != nil {
		return "", err
	}

	resp, err := c.executeRequest(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		c.logger.Errorf("Decode error: %v", err)
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(out.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return out.Choices[0].Message.Content, nil
}

// AskLLMWithMemory sends a chat completion request with conversation history
func (c *LLMClient) AskLLMWithMemory(systemPrompt, userPrompt string, history []ChatMessage) (string, error) {
	c.logger.Infof("Using LLM settings with memory - BaseUrl: %s, Model: %s, Temperature: %.1f, History: %d messages", c.settings.BaseUrl, c.settings.Model, c.settings.Temperature, len(history))

	req := &ChatRequest{
		Model:       c.settings.Model,
		Messages:    buildMessages(systemPrompt, userPrompt, history),
		Stream:      false,
		Temperature: c.settings.Temperature,
	}

	httpReq, err := c.createRequest(req)
	if err != nil {
		return "", err
	}

	resp, err := c.executeRequest(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		c.logger.Errorf("Decode error: %v", err)
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(out.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return out.Choices[0].Message.Content, nil
}

// AskLLMWithFunctions sends a chat completion request with function calling support
func (c *LLMClient) AskLLMWithFunctions(systemPrompt, userPrompt string, functions []map[string]interface{}) (string, []FunctionCall, error) {
	c.logger.Infof("Using LLM settings with functions - BaseUrl: %s, Model: %s, Temperature: %.1f, Functions: %d", c.settings.BaseUrl, c.settings.Model, c.settings.Temperature, len(functions))

	req := &ChatRequest{
		Model:       c.settings.Model,
		Messages:    buildMessages(systemPrompt, userPrompt, nil),
		Stream:      false,
		Temperature: c.settings.Temperature,
	}

	// Add functions if provided
	if len(functions) > 0 {
		req.Functions = functions
		req.FunctionCall = "auto"
	}

	httpReq, err := c.createRequest(req)
	if err != nil {
		return "", nil, err
	}

	resp, err := c.executeRequest(httpReq)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var out ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		c.logger.Errorf("Decode error: %v", err)
		return "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(out.Choices) == 0 {
		return "", nil, fmt.Errorf("no choices in response")
	}

	message := out.Choices[0].Message
	var functionCalls []FunctionCall

	// Check if the model wants to call functions
	if message.FunctionCall != nil {
		functionCalls = append(functionCalls, *message.FunctionCall)
	}

	return message.Content, functionCalls, nil
}

// processStreamingResponse handles streaming response processing
func (c *LLMClient) processStreamingResponse(resp *http.Response, processChunk func(string)) error {
	// Send opening tag first
	processChunk("<answer>")

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			// Send closing tag on error before returning
			processChunk("</answer>")
			return err
		}

		trimmedLine := bytes.TrimSpace(line)
		if !bytes.HasPrefix(trimmedLine, []byte("data: ")) {
			continue
		}

		data := trimmedLine[6:]
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) {
			continue
		}

		var chunk map[string]interface{}
		if err := json.Unmarshal(data, &chunk); err != nil {
			c.logger.Errorf("JSON decode error: %v (raw data: %s)", err, data)
			continue
		}

		if choices, ok := chunk["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if delta, ok := choice["delta"].(map[string]interface{}); ok {
					if content, ok := delta["content"].(string); ok && content != "" {
						processChunk(content)
					}
				}
			}
		}
	}

	// Send closing tag at the end
	processChunk("</answer>")
	return nil
}

// AskLLMWithStream sends a streaming chat completion request
func (c *LLMClient) AskLLMWithStream(systemPrompt, userPrompt string, processChunk func(string)) error {
	c.logger.Infof("Using LLM settings - BaseUrl: %s, Model: %s, Temperature: %.1f", c.settings.BaseUrl, c.settings.Model, c.settings.Temperature)

	fmt.Println("userPrompt:", userPrompt)

	req := &ChatRequest{
		Model:       c.settings.Model,
		Messages:    buildMessages(systemPrompt, userPrompt, nil),
		Stream:      true,
		Temperature: c.settings.Temperature,
	}

	httpReq, err := c.createRequest(req)
	if err != nil {
		return err
	}

	resp, err := c.executeRequest(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.processStreamingResponse(resp, processChunk)
}

// AskLLMWithStreamAndMemory sends a streaming chat completion request with conversation history
func (c *LLMClient) AskLLMWithStreamAndMemory(systemPrompt, userPrompt string, history []ChatMessage, processChunk func(string)) error {
	c.logger.Infof("Using LLM settings with memory - BaseUrl: %s, Model: %s, Temperature: %.1f, History: %d messages", c.settings.BaseUrl, c.settings.Model, c.settings.Temperature, len(history))

	req := &ChatRequest{
		Model:       c.settings.Model,
		Messages:    buildMessages(systemPrompt, userPrompt, history),
		Stream:      true,
		Temperature: c.settings.Temperature,
	}

	httpReq, err := c.createRequest(req)
	if err != nil {
		return err
	}

	resp, err := c.executeRequest(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.processStreamingResponse(resp, processChunk)
}

// Note: We no longer keep a persistent global client to avoid premature
// initialization issues. Callers should construct a client per settings
// or use the helper wrapper functions below which create a short-lived client.

// Backward compatibility functions that use short-lived clients.
// Optional settings can be provided to override env configuration.
func AskLLM(systemPrompt, userPrompt string, optionalSettings ...LLMSettings) (string, error) {
	var client *LLMClient
	if len(optionalSettings) > 0 {
		client = NewLLMClient(optionalSettings[0])
	} else {
		client = NewLLMClient()
	}
	return client.AskLLM(systemPrompt, userPrompt)
}

func AskLLMWithMemory(systemPrompt, userPrompt string, history []ChatMessage, optionalSettings ...LLMSettings) (string, error) {
	var client *LLMClient
	if len(optionalSettings) > 0 {
		client = NewLLMClient(optionalSettings[0])
	} else {
		client = NewLLMClient()
	}
	return client.AskLLMWithMemory(systemPrompt, userPrompt, history)
}

func AskLLMWithFunctions(systemPrompt, userPrompt string, functions []map[string]interface{}, optionalSettings ...LLMSettings) (string, []FunctionCall, error) {
	var client *LLMClient
	if len(optionalSettings) > 0 {
		client = NewLLMClient(optionalSettings[0])
	} else {
		client = NewLLMClient()
	}
	return client.AskLLMWithFunctions(systemPrompt, userPrompt, functions)
}

func AskLLMWithStream(systemPrompt, userPrompt string, processChunk func(string), optionalSettings ...LLMSettings) error {
	var client *LLMClient
	if len(optionalSettings) > 0 {
		client = NewLLMClient(optionalSettings[0])
	} else {
		client = NewLLMClient()
	}
	return client.AskLLMWithStream(systemPrompt, userPrompt, processChunk)
}

func AskLLMWithStreamAndMemory(systemPrompt, userPrompt string, history []ChatMessage, processChunk func(string), optionalSettings ...LLMSettings) error {
	var client *LLMClient
	if len(optionalSettings) > 0 {
		client = NewLLMClient(optionalSettings[0])
	} else {
		client = NewLLMClient()
	}
	return client.AskLLMWithStreamAndMemory(systemPrompt, userPrompt, history, processChunk)
}

// ProcessFunctionCalls executes function calls and returns results
func ProcessFunctionCalls(functionCalls []FunctionCall, registry interface{}) ([]ChatEntry, error) {
	// This would interact with the tools registry to execute functions
	// For now, return a placeholder implementation
	var messages []ChatEntry

	for _, fc := range functionCalls {
		// This would call the actual function implementation
		result := fmt.Sprintf("Function %s called with args: %s", fc.Name, fc.Arguments)
		messages = append(messages, ChatEntry{
			Role:    "function",
			Name:    fc.Name,
			Content: result,
		})
	}

	return messages, nil
}
