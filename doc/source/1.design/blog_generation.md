# Blog Generator Command

This document describes the Go implementation of the daily technical blog generator, which is equivalent to the original Python script you provided.

## Overview

The blog generator is a command-line tool that uses LLM (Large Language Model) services to generate daily technical blog content based on customizable templates and ideas. It supports both English and Chinese content generation.

## Features

- **Template-based generation**: Uses configurable YAML templates for consistent blog structure
- **Multi-language support**: Generate blogs in English, Chinese, or both
- **Customizable prompts**: Override system prompts and templates
- **Flexible output**: Print to console or save to file
- **LLM integration**: Works with OpenAI-compatible APIs
- **Command-line interface**: Easy to use and integrate into workflows

## Prerequisites

1. **Environment Setup**: Set up your LLM API credentials using either environment variables or a `.env` file

   **Option A: Environment Variables**
   ```bash
   export LLM_BASE_URL="https://api.openai.com/v1"
   export LLM_API_KEY="your-api-key-here"
   export LLM_MODEL="gpt-3.5-turbo"
   export LLM_TEMPERATURE="1.0"
   ```

   **Option B: .env File** (recommended for development)
   Create a `.env` file in the backend directory:
   ```bash
   # .env file for OpenAI
   LLM_BASE_URL=https://api.openai.com/v1
   LLM_API_KEY=your-openai-api-key-here
   LLM_MODEL=gpt-3.5-turbo
   LLM_TEMPERATURE=1.0
   
   # Other configurations for different providers:
   
   # For local LLM (e.g., Ollama)
   # LLM_BASE_URL=http://localhost:11434/v1
   # LLM_API_KEY=not-required-for-local
   # LLM_MODEL=llama2
   
   # For Azure OpenAI
   # LLM_BASE_URL=https://your-resource.openai.azure.com/openai/deployments/your-deployment/
   # LLM_API_KEY=your-azure-api-key
   # LLM_MODEL=gpt-35-turbo
   ```

2. **Build the Application**:
   ```bash
   cd backend
   go build -mod=mod -o lazy-rabbit-secretary .
   ```

## Usage

### Basic Usage

Generate a blog with a technical idea:

```bash
./lazy-rabbit-secretary blog --idea "用 webrtc 和 pion 打造一款网络录音机"
```

### Language Options

```bash
# Generate Chinese blog only
./lazy-rabbit-secretary blog --idea "Building a distributed cache with Redis" --lang cn

# Generate English blog only
./lazy-rabbit-secretary blog --idea "Building a distributed cache with Redis" --lang en

# Generate both languages (default)
./lazy-rabbit-secretary blog --idea "Building a distributed cache with Redis" --lang both
```

### Save to File

```bash
./lazy-rabbit-secretary blog \
    --idea "Microservices with Go and Docker" \
    --output "blog_$(date +%Y-%m-%d).md"
```

### Custom Title and Template

```bash
./lazy-rabbit-secretary blog \
    --idea "Machine Learning in Production" \
    --title "My ML Journey - $(date +%Y-%m-%d)" \
    --template "write_daily_blog_en" \
    --lang en
```

### Override LLM Settings

```bash
./lazy-rabbit-secretary blog \
    --idea "Kubernetes best practices" \
    --base-url "https://your-llm-endpoint.com/v1" \
    --api-key "your-custom-key" \
    --model "gpt-4" \
    --temperature 0.7
```

## Configuration Precedence

The application uses the following order of precedence for configuration (highest to lowest):

1. **Command-line flags** (highest priority)
2. **Environment variables** (from `.env` file or system environment)
3. **Default values** (lowest priority)

For example:
```bash
# This will use the API key from command line, overriding any .env or environment variable
./lazy-rabbit-secretary blog --idea "test" --api-key "override-key"

# This will use settings from .env file if no command-line flags are provided
./lazy-rabbit-secretary blog --idea "test"
```

## Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--idea` | `-i` | **Required.** The main technical idea for the blog | - |
| `--title` | `-t` | Custom blog title | Auto-generated with date |
| `--lang` | `-l` | Language: 'en', 'cn', or 'both' | `both` |
| `--output` | `-o` | Output file path | Print to stdout |
| `--system` | `-s` | Custom system prompt | Use template default |
| `--prompts-file` | - | Path to prompts YAML file | `config/prompts.yaml` |
| `--template` | - | Specific template name | Auto-selected by language |
| `--base-url` | - | LLM API base URL | From `LLM_BASE_URL` env |
| `--api-key` | - | LLM API key | From `LLM_API_KEY` env |
| `--model` | - | LLM model name | From `LLM_MODEL` env |
| `--temperature` | - | LLM temperature | From `LLM_TEMPERATURE` env |

## Templates

The system includes three built-in blog templates:

### 1. `write_daily_blog` (Both Languages)
- **Description**: Generate daily technical blog content
- **Languages**: Both English and Chinese
- **Use case**: Default template for bilingual content

### 2. `write_daily_blog_en` (English Only)
- **Description**: Generate daily technical blog content in English
- **Languages**: English only
- **Use case**: English-focused technical content

### 3. `write_daily_blog_cn` (Chinese Only)
- **Description**: Generate daily technical blog content in Chinese
- **Languages**: Chinese only
- **Use case**: Chinese-focused technical content

## Blog Structure

All templates generate blogs with the following structure:

```markdown
# [Blog Title]

## "[Your Technical Idea]"
### what
[Explanation of what the technology/concept is]

### why
[Why it's important/useful]

### how
[How to implement/use it]

### example
[Practical examples]

### summary
[Summary and conclusions]

### reference
[References and links]

## Daily recommendation of 1 github hot project
[GitHub project recommendation]

## Daily recommendation of 1 best practice in software development and AI applications
[Best practice recommendation]

## Daily practice of 1 leetcode algorithm question by one of the following languages: go/java/python/typescript/rust
[Algorithm practice]

## Daily practice of 1 classic design pattern by one of the following languages: go/java/python/typescript/rust
[Design pattern practice]

## Daily recitation of 10 English quotes
[Inspirational quotes]
```

## Examples

### Example 1: WebRTC Network Recorder

```bash
./lazy-rabbit-secretary blog \
    --idea "用 webrtc 和 pion 打造一款网络录音机" \
    --lang cn \
    --title "技术探索 - WebRTC 录音机实现"
```

This generates a Chinese blog about building a network audio recorder using WebRTC and Pion.

### Example 2: Go Microservices

```bash
./lazy-rabbit-secretary blog \
    --idea "Building scalable microservices with Go and gRPC" \
    --lang en \
    --output "microservices-blog.md"
```

This generates an English blog about microservices and saves it to a file.

### Example 3: Daily Automation

```bash
#!/bin/bash
# daily-blog.sh - Automated daily blog generation

IDEA=$(curl -s "https://api.github.com/search/repositories?q=stars:>1000+language:go&sort=stars&order=desc" | jq -r '.items[0].description')

./lazy-rabbit-secretary blog \
    --idea "Exploring: $IDEA" \
    --output "blogs/$(date +%Y-%m-%d)-daily-blog.md" \
    --lang both
```

## Integration with Original Python Script

This Go implementation provides the same functionality as your original Python script:

**Python Script Features** → **Go Implementation**
- ✅ Jinja2 templating → Custom template rendering system
- ✅ Environment variable support → Complete env var support
- ✅ LLM integration → Uses existing LLM service
- ✅ Date/time handling → Built-in time formatting
- ✅ Template data injection → Template data system
- ✅ Async LLM calls → Synchronous calls (can be extended)

## Troubleshooting

### Common Issues

1. **"Failed to load prompt templates"**
   - Ensure `config/prompts.yaml` exists and is valid YAML
   - Check file permissions

2. **"Template not found"**
   - Verify template name exists in prompts.yaml
   - Use `--template` flag with correct template name

3. **LLM API errors**
   - Check your API key and base URL
   - Verify the model name is correct
   - Ensure you have sufficient API credits

4. **Empty or invalid responses**
   - Try reducing temperature (0.1-0.7)
   - Check if the model supports your language
   - Verify your prompt templates are correctly formatted

### Debug Mode

For debugging, you can check the generated prompts by examining the log output:

```bash
./lazy-rabbit-secretary blog --idea "test" --lang en 2>&1 | grep "Using LLM settings"
```

## Extending the System

### Adding New Templates

1. Edit `config/prompts.yaml`
2. Add your new template following the existing pattern:

```yaml
prompts:
  my_custom_blog:
    description: "My custom blog template"
    system_prompt: "You are a specialized writer for {{domain}}"
    user_prompt: |
      Generate a blog about {{idea}} focusing on {{aspect}}.
      
      Structure:
      - Introduction
      - Main content
      - Conclusion
    tags: blog, custom
```

3. Use with `--template my_custom_blog`

### Adding New Languages

1. Create language-specific templates in `prompts.yaml`
2. Update the language selection logic in `cmd/blog.go`
3. Add appropriate system prompts for the new language

## Performance Considerations

- **Template loading**: Templates are loaded once per execution
- **LLM calls**: Single synchronous call per blog generation
- **Memory usage**: Minimal, suitable for scripting and automation
- **File I/O**: Efficient file writing with proper error handling

## Security Notes

- API keys are passed via environment variables (secure)
- No API keys are logged or exposed in output
- Template injection is prevented through controlled variable substitution
- File paths are validated for output operations

---

This Go implementation provides a robust, production-ready alternative to the Python script with enhanced configurability and integration capabilities.
