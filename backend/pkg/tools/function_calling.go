package tools

import (
	"fmt"
	"time"
)

// FunctionCall represents a function call request
type FunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// FunctionDefinition defines a function that can be called
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Handler     FunctionHandler        `json:"-"`
}

// FunctionHandler is the actual function implementation
type FunctionHandler func(args map[string]interface{}) (interface{}, error)

// FunctionRegistry manages available functions
type FunctionRegistry struct {
	functions map[string]FunctionDefinition
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]FunctionDefinition),
	}
}

// RegisterFunction registers a new function
func (fr *FunctionRegistry) RegisterFunction(def FunctionDefinition) {
	fr.functions[def.Name] = def
}

// GetFunction retrieves a function definition
func (fr *FunctionRegistry) GetFunction(name string) (FunctionDefinition, bool) {
	def, exists := fr.functions[name]
	return def, exists
}

// ListFunctions returns all available function definitions (for LLM)
func (fr *FunctionRegistry) ListFunctions() []FunctionDefinition {
	var defs []FunctionDefinition
	for _, def := range fr.functions {
		defs = append(defs, def)
	}
	return defs
}

// ExecuteFunction executes a function call
func (fr *FunctionRegistry) ExecuteFunction(call FunctionCall) (interface{}, error) {
	def, exists := fr.functions[call.Name]
	if !exists {
		return nil, fmt.Errorf("function '%s' not found", call.Name)
	}

	return def.Handler(call.Arguments)
}

// GetFunctionDefinitionsForLLM returns function definitions in OpenAI format
func (fr *FunctionRegistry) GetFunctionDefinitionsForLLM() []map[string]interface{} {
	var functions []map[string]interface{}

	for _, def := range fr.functions {
		functions = append(functions, map[string]interface{}{
			"name":        def.Name,
			"description": def.Description,
			"parameters":  def.Parameters,
		})
	}

	return functions
}

// DefaultRegistry provides commonly used functions
func DefaultRegistry() *FunctionRegistry {
	registry := NewFunctionRegistry()

	// Register built-in functions
	registry.RegisterFunction(GetCurrentDateFunction())
	registry.RegisterFunction(GetCurrentTimeFunction())
	registry.RegisterFunction(GetWeatherFunction())

	return registry
}

// GetCurrentDateFunction returns current date
func GetCurrentDateFunction() FunctionDefinition {
	return FunctionDefinition{
		Name:        "get_current_date",
		Description: "Get the current date in specified format",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"format": map[string]interface{}{
					"type":        "string",
					"description": "Date format (default: 2006-01-02)",
					"default":     "2006-01-02",
				},
				"timezone": map[string]interface{}{
					"type":        "string",
					"description": "Timezone (default: Local)",
					"default":     "Local",
				},
			},
		},
		Handler: func(args map[string]interface{}) (interface{}, error) {
			format := "2006-01-02"
			if f, ok := args["format"].(string); ok && f != "" {
				format = f
			}

			timezone := "Local"
			if tz, ok := args["timezone"].(string); ok && tz != "" {
				timezone = tz
			}

			var loc *time.Location
			var err error
			if timezone == "Local" {
				loc = time.Local
			} else {
				loc, err = time.LoadLocation(timezone)
				if err != nil {
					return nil, fmt.Errorf("invalid timezone: %v", err)
				}
			}

			now := time.Now().In(loc)
			return now.Format(format), nil
		},
	}
}

// GetCurrentTimeFunction returns current time
func GetCurrentTimeFunction() FunctionDefinition {
	return FunctionDefinition{
		Name:        "get_current_time",
		Description: "Get the current time in specified format",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"format": map[string]interface{}{
					"type":        "string",
					"description": "Time format (default: 15:04:05)",
					"default":     "15:04:05",
				},
				"timezone": map[string]interface{}{
					"type":        "string",
					"description": "Timezone (default: Local)",
					"default":     "Local",
				},
			},
		},
		Handler: func(args map[string]interface{}) (interface{}, error) {
			format := "15:04:05"
			if f, ok := args["format"].(string); ok && f != "" {
				format = f
			}

			timezone := "Local"
			if tz, ok := args["timezone"].(string); ok && tz != "" {
				timezone = tz
			}

			var loc *time.Location
			var err error
			if timezone == "Local" {
				loc = time.Local
			} else {
				loc, err = time.LoadLocation(timezone)
				if err != nil {
					return nil, fmt.Errorf("invalid timezone: %v", err)
				}
			}

			now := time.Now().In(loc)
			return now.Format(format), nil
		},
	}
}
