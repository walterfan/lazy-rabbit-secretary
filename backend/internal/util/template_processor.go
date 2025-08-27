package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/walterfan/lazy-rabbit-reminder/internal/tools"
)

// TemplateProcessor handles template rendering with function calling support
type TemplateProcessor struct {
	registry *tools.FunctionRegistry
}

// NewTemplateProcessor creates a new template processor
func NewTemplateProcessor() *TemplateProcessor {
	return &TemplateProcessor{
		registry: tools.DefaultRegistry(),
	}
}

// ProcessTemplate processes a template with both static data and function calls
func (tp *TemplateProcessor) ProcessTemplate(template string, data TemplateData) (string, error) {
	// First, handle function calls
	processedTemplate, err := tp.processFunctionCalls(template)
	if err != nil {
		return "", fmt.Errorf("failed to process function calls: %v", err)
	}

	// Then handle static template variables
	result := RenderTemplate(processedTemplate, data)

	return result, nil
}

// processFunctionCalls finds and executes function calls in the template
func (tp *TemplateProcessor) processFunctionCalls(template string) (string, error) {
	// Regex to match function calls like {{function_name(arg1="value1", arg2="value2")}}
	funcCallRegex := regexp.MustCompile(`\{\{(\w+)\((.*?)\)\}\}`)

	result := template
	matches := funcCallRegex.FindAllStringSubmatch(template, -1)

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		fullMatch := match[0]
		functionName := match[1]
		argsString := match[2]

		// Parse arguments
		args, err := tp.parseArguments(argsString)
		if err != nil {
			return "", fmt.Errorf("failed to parse arguments for function %s: %v", functionName, err)
		}

		// Execute function
		functionResult, err := tp.registry.ExecuteFunction(tools.FunctionCall{
			Name:      functionName,
			Arguments: args,
		})
		if err != nil {
			return "", fmt.Errorf("failed to execute function %s: %v", functionName, err)
		}

		// Convert result to string
		resultString := tp.resultToString(functionResult)

		// Replace the function call with the result
		result = strings.ReplaceAll(result, fullMatch, resultString)
	}

	return result, nil
}

// parseArguments parses function arguments from string format
func (tp *TemplateProcessor) parseArguments(argsString string) (map[string]interface{}, error) {
	args := make(map[string]interface{})

	if strings.TrimSpace(argsString) == "" {
		return args, nil
	}

	// Simple argument parsing - handles key="value" format
	argRegex := regexp.MustCompile(`(\w+)=["']([^"']*)["']`)
	matches := argRegex.FindAllStringSubmatch(argsString, -1)

	for _, match := range matches {
		if len(match) == 3 {
			key := match[1]
			value := match[2]
			args[key] = value
		}
	}

	// Also handle simple key=value without quotes
	simpleArgRegex := regexp.MustCompile(`(\w+)=([^,\s]+)`)
	simpleMatches := simpleArgRegex.FindAllStringSubmatch(argsString, -1)

	for _, match := range simpleMatches {
		if len(match) == 3 {
			key := match[1]
			value := match[2]
			// Skip if already processed with quotes
			if _, exists := args[key]; !exists {
				args[key] = value
			}
		}
	}

	return args, nil
}

// resultToString converts function result to string representation
func (tp *TemplateProcessor) resultToString(result interface{}) string {
	switch v := result.(type) {
	case string:
		return v
	case map[string]interface{}:
		// For complex objects, format them nicely
		if weather, ok := v["weather"].(string); ok {
			if city, ok := v["city"].(string); ok {
				if temp, ok := v["temperature"].(string); ok {
					return fmt.Sprintf("%s, %s, %s°C", city, weather, temp)
				}
			}
		}
		// Fallback to JSON
		if jsonBytes, err := json.Marshal(v); err == nil {
			return string(jsonBytes)
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// AddCustomFunction adds a custom function to the registry
func (tp *TemplateProcessor) AddCustomFunction(def tools.FunctionDefinition) {
	tp.registry.RegisterFunction(def)
}

// GetAvailableFunctions returns list of available function names
func (tp *TemplateProcessor) GetAvailableFunctions() []string {
	var names []string
	for _, def := range tp.registry.ListFunctions() {
		names = append(names, def.Name)
	}
	return names
}

// ProcessTemplateWithLLM processes template with both function calls and LLM enhancement
func (tp *TemplateProcessor) ProcessTemplateWithLLM(template string, data TemplateData, llmProcessor func(string) (string, error)) (string, error) {
	// First process function calls
	processed, err := tp.ProcessTemplate(template, data)
	if err != nil {
		return "", err
	}

	// Then optionally enhance with LLM
	if llmProcessor != nil {
		enhanced, err := llmProcessor(processed)
		if err != nil {
			return processed, nil // Return processed version if LLM fails
		}
		return enhanced, nil
	}

	return processed, nil
}
