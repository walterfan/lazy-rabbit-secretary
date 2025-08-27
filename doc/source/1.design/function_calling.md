# Function Calling System for Blog Generator

This document describes the function calling system implemented in the Go blog generator, which allows dynamic generation of template placeholders like `{{date}}` and `{{weather}}` using actual API calls and computed values.

## Overview

The function calling system enables templates to include dynamic content by calling functions during template processing. This is similar to the Python example you provided, but implemented as a more generic and extensible system in Go.

## Key Features

- **Dynamic Date/Time Generation**: Get current date and time in various formats
- **Weather API Integration**: Fetch real-time weather data from Gaode API
- **Extensible Function Registry**: Easy to add custom functions
- **Template Integration**: Seamless integration with existing template system
- **Error Handling**: Graceful fallback when functions fail

## Architecture

### Core Components

1. **Function Registry** (`internal/tools/function_calling.go`)
   - Manages available functions
   - Handles function execution
   - Provides function definitions for LLM integration

2. **Template Processor** (`internal/util/template_processor.go`)
   - Processes templates with function calls
   - Parses function call syntax
   - Handles both static and dynamic placeholders

3. **Built-in Functions** (`internal/tools/`)
   - Date/time functions
   - Weather API integration
   - Easy to extend with custom functions

## Available Functions

### 1. `get_current_date`
Gets the current date in specified format.

**Parameters:**
- `format` (string, optional): Date format (default: "2006-01-02")
- `timezone` (string, optional): Timezone (default: "Local")

**Examples:**
```
{{get_current_date()}}                          → 2025-08-11
{{get_current_date(format="2006年01月02日")}}      → 2025年08月11日
{{get_current_date(format="January 2, 2006")}}   → August 11, 2025
```

### 2. `get_current_time`
Gets the current time in specified format.

**Parameters:**
- `format` (string, optional): Time format (default: "15:04:05")
- `timezone` (string, optional): Timezone (default: "Local")

**Examples:**
```
{{get_current_time()}}                      → 16:57:38
{{get_current_time(format="3:04 PM")}}      → 4:57 PM
{{get_current_time(timezone="UTC")}}        → 08:57:38
```

### 3. `get_weather`
Gets weather information from Gaode API.

**Parameters:**
- `city` (string, required): City name or code (e.g., "北京", "上海")
- `ext` (string, optional): "base" for real-time, "all" for forecast (default: "base")

**Examples:**
```
{{get_weather(city="北京")}}                    → 北京, 晴, 25°C
{{get_weather(city="上海", ext="base")}}        → 上海, 多云, 23°C
{{get_weather(city="{{city}}", ext="base")}}   → Uses template variable
```

**Prerequisites:**
- Set `GAODE_KEY` environment variable with your Gaode API key
- Register at [高德开放平台](https://lbs.amap.com/) to get API key

## Usage in Templates

### Function Call Syntax

Functions are called using the syntax: `{{function_name(param1="value1", param2="value2")}}`

### Template Examples

```yaml
prompts:
  weather_blog:
    user_prompt: |
      Today is {{get_current_date(format="2006年01月02日")}}, 
      Weather: {{get_weather(city="{{city}}", ext="base")}}
      
      Write a blog about {{topic}} considering today's weather.
```

### Complex Example

```yaml
prompts:
  daily_summary:
    user_prompt: |
      # Daily Tech Summary - {{get_current_date(format="January 2, 2006")}}
      
      **Time:** {{get_current_time(format="3:04 PM")}}
      **Weather:** {{get_weather(city="{{city}}")}}
      **Topic:** {{idea}}
      
      Generate a comprehensive summary...
```

## Command Line Usage

### Basic Blog Generation with Function Calls

```bash
# Use default template with function calls
./lazy-rabbit-reminder blog --idea "WebRTC录音机" --city "北京"

# Specify custom city for weather
./lazy-rabbit-reminder blog --idea "微服务架构" --city "上海" --lang cn

# Test function calling system
./lazy-rabbit-reminder test-functions
```

### Environment Setup for Weather

```bash
# Set Gaode API key for weather functions
export GAODE_KEY="your-gaode-api-key"

# Then run blog generation
./lazy-rabbit-reminder blog --idea "云原生技术" --city "深圳"
```

## Adding Custom Functions

### Step 1: Define Function Handler

```go
func GetCustomFunction() tools.FunctionDefinition {
    return tools.FunctionDefinition{
        Name:        "get_stock_price",
        Description: "Get current stock price",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "symbol": map[string]interface{}{
                    "type":        "string",
                    "description": "Stock symbol (e.g., AAPL)",
                },
            },
            "required": []string{"symbol"},
        },
        Handler: func(args map[string]interface{}) (interface{}, error) {
            symbol := args["symbol"].(string)
            // Implement stock price fetching logic
            return fmt.Sprintf("$150.25 (%s)", symbol), nil
        },
    }
}
```

### Step 2: Register Function

```go
// In tools/function_calling.go DefaultRegistry()
registry.RegisterFunction(GetCustomFunction())
```

### Step 3: Use in Templates

```yaml
user_prompt: |
  Stock price for AAPL: {{get_stock_price(symbol="AAPL")}}
```

## Integration with LLM Function Calling

The system also supports OpenAI-style function calling for LLM integration:

```go
// Get function definitions for LLM
registry := tools.DefaultRegistry()
functions := registry.GetFunctionDefinitionsForLLM()

// Use with LLM
result, functionCalls, err := llm.AskLLMWithFunctions(
    systemPrompt, userPrompt, functions, settings
)
```

## Error Handling

The system provides graceful error handling:

1. **Function Not Found**: Returns clear error message
2. **Invalid Parameters**: Validates required parameters
3. **API Failures**: Falls back gracefully with error logging
4. **Missing Environment Variables**: Clear error messages

Example error handling in templates:

```go
// Template processor handles errors gracefully
result, err := processor.ProcessTemplate(template, data)
if err != nil {
    // Falls back to basic template rendering
    result = util.RenderTemplate(template, data)
}
```

## Testing

### Test All Functions

```bash
# Run comprehensive function test
./lazy-rabbit-reminder test-functions
```

### Test Specific Function

```bash
# Test with weather API (requires GAODE_KEY)
export GAODE_KEY="your-api-key"
./lazy-rabbit-reminder blog --idea "Test" --city "北京"
```

## Performance Considerations

- **Function Execution**: Functions are executed synchronously during template processing
- **API Calls**: Weather API calls have 10-second timeout
- **Caching**: Consider implementing caching for frequently called functions
- **Error Recovery**: Failed function calls don't break template processing

## Security Considerations

- **API Keys**: Store API keys in environment variables, never in code
- **Input Validation**: All function parameters are validated
- **Rate Limiting**: Be aware of API rate limits (Gaode: 100,000 calls/day)
- **Error Information**: Function errors don't expose sensitive information

## Comparison with Python Implementation

Your original Python script used function calling like this:

```python
functions = [
    {
        "name": "get_weather",
        "description": "根据城市代码或城市名查询高德实时或未来天气",
        "parameters": {
            "type": "object",
            "properties": {
                "city": {"type": "string", "description": "城市代码或中文城市名"},
                "ext": {"type": "string", "enum": ["base", "all"]}
            },
            "required": ["city"]
        }
    }
]
```

The Go implementation provides:

| Feature | Python | Go Implementation |
|---------|--------|-------------------|
| Function Definition | Manual JSON | Structured Go types |
| Function Execution | OpenAI handles | Local execution + LLM integration |
| Template Integration | Separate step | Built into template processor |
| Error Handling | Basic try/catch | Comprehensive error handling |
| Extensibility | Manual function addition | Plugin-like function registry |
| Type Safety | Runtime validation | Compile-time + runtime validation |

## Future Enhancements

1. **Function Caching**: Cache function results to improve performance
2. **Async Functions**: Support for asynchronous function execution
3. **Function Dependencies**: Support functions that depend on other functions
4. **Custom Parsers**: Support for different function call syntax formats
5. **Function Marketplace**: Plugin system for community-contributed functions

---

This function calling system provides a robust foundation for dynamic template generation, combining the power of real-time data with the flexibility of template-based content creation.
