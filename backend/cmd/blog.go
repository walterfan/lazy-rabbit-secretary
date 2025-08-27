package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/walterfan/lazy-rabbit-reminder/internal/llm"
	"github.com/walterfan/lazy-rabbit-reminder/internal/log"
	"github.com/walterfan/lazy-rabbit-reminder/internal/tools"
	"github.com/walterfan/lazy-rabbit-reminder/internal/util"
)

var blogCmd = &cobra.Command{
	Use:   "blog",
	Short: "Generate a daily technical blog",
	Long: `Generate a daily technical blog with customizable topics and ideas.
This command creates both English and Chinese versions of the blog based on a template.`,
	Run: runBlogGenerator,
}

var (
	blogIdea     string
	blogTitle    string
	blogLang     string
	blogOutput   string
	systemPrompt string
	baseUrl      string
	apiKey       string
	model        string
	temperature  float64
	promptsFile  string
	templateName string
	city         string
)

func init() {
	rootCmd.AddCommand(blogCmd)

	// Add flags for customization
	blogCmd.Flags().StringVarP(&blogIdea, "idea", "i", "", "The main technical idea for the blog (required)")
	blogCmd.Flags().StringVarP(&blogTitle, "title", "t", "", "Custom blog title (optional, auto-generated if not provided)")
	blogCmd.Flags().StringVarP(&blogLang, "lang", "l", "cn", "Language for blog generation: 'en', 'cn', or 'both' (default: both)")
	blogCmd.Flags().StringVarP(&blogOutput, "output", "o", "", "Output file path (optional, prints to stdout if not provided)")
	blogCmd.Flags().StringVarP(&systemPrompt, "system", "s", "", "Custom system prompt (optional)")
	blogCmd.Flags().StringVar(&promptsFile, "prompts-file", "config/prompts.yaml", "Path to prompts configuration file")
	blogCmd.Flags().StringVar(&templateName, "template", "", "Template name to use (optional, auto-selected based on language)")
	blogCmd.Flags().StringVar(&city, "city", "", "City for weather information (default: none)")

	// LLM configuration flags
	blogCmd.Flags().StringVar(&baseUrl, "base-url", "", "LLM API base URL (optional, uses env LLM_BASE_URL if not provided)")
	blogCmd.Flags().StringVar(&apiKey, "api-key", "", "LLM API key (optional, uses env LLM_API_KEY if not provided)")
	blogCmd.Flags().StringVar(&model, "model", "", "LLM model name (optional, uses env LLM_MODEL if not provided)")
	blogCmd.Flags().Float64Var(&temperature, "temperature", 0, "LLM temperature (optional, uses env LLM_TEMPERATURE if not provided)")

	// Mark the idea flag as required
	blogCmd.MarkFlagRequired("idea")
}

func runBlogGenerator(cmd *cobra.Command, args []string) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// .env file is optional, so we just log a warning instead of failing
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// Initialize logger if not already done
	err = log.InitLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	logger := log.GetLogger()

	// Get current date
	today := time.Now()
	todayStr := today.Format("2006-01-02")

	// Generate title if not provided
	if blogTitle == "" {
		blogTitle = fmt.Sprintf("My Tech Blog - %s", todayStr)
	}

	// Load prompt templates
	promptConfig, err := util.LoadPromptTemplates(promptsFile)
	if err != nil {
		logger.Errorf("Failed to load prompt templates: %v", err)
		fmt.Printf("Error loading prompt templates: %v\n", err)
		return
	}

	// Determine which template to use
	var selectedTemplate util.PromptTemplate
	if templateName != "" {
		// Use explicit template name
		selectedTemplate, err = promptConfig.GetPromptTemplate(templateName)
		if err != nil {
			logger.Errorf("Failed to get template '%s': %v", templateName, err)
			fmt.Printf("Error getting template '%s': %v\n", templateName, err)
			return
		}
	} else {
		// Auto-select template based on language
		var templateKey string
		switch blogLang {
		case "en":
			templateKey = "write_daily_blog_en"
		case "cn":
			templateKey = "write_daily_blog_cn"
		default:
			templateKey = "write_daily_blog"
		}

		selectedTemplate, err = promptConfig.GetPromptTemplate(templateKey)
		if err != nil {
			logger.Errorf("Failed to get default template '%s': %v", templateKey, err)
			fmt.Printf("Error getting default template '%s': %v\n", templateKey, err)
			return
		}
	}

	// Create template data
	templateData := util.TemplateData{
		"title": blogTitle,
		"today": todayStr,
		"idea":  blogIdea,
		"date":  todayStr, // Add date for function calling
		"city":  city,     // Add city for weather function calling
	}

	// Create template processor with function calling support
	templateProcessor := util.NewTemplateProcessor()

	// Process templates with function calling support
	renderedSystemPrompt, err := templateProcessor.ProcessTemplate(selectedTemplate.SystemPrompt, templateData)
	if err != nil {
		logger.Warnf("Failed to process system prompt with functions, falling back to basic rendering: %v", err)
		renderedSystemPrompt = util.RenderTemplate(selectedTemplate.SystemPrompt, templateData)
	}

	renderedUserPrompt, err := templateProcessor.ProcessTemplate(selectedTemplate.UserPrompt, templateData)
	if err != nil {
		logger.Warnf("Failed to process user prompt with functions, falling back to basic rendering: %v", err)
		renderedUserPrompt = util.RenderTemplate(selectedTemplate.UserPrompt, templateData)
	}

	// Use custom system prompt if provided
	if systemPrompt != "" {
		renderedSystemPrompt = systemPrompt
	}

	logger.Infof("Generating blog with idea: %s", blogIdea)
	logger.Infof("Date: %s", todayStr)
	logger.Infof("Language: %s", blogLang)
	logger.Infof("Template: %s", selectedTemplate.Description)
	logger.Infof("City: %s", city)

	var result string
	// If city provided, use OpenAI function calling with weather function
	if strings.TrimSpace(city) != "" {
		// Build a minimal registry with only weather function to expose to LLM
		registry := tools.NewFunctionRegistry()
		registry.RegisterFunction(tools.GetWeatherFunction())

		functions := registry.GetFunctionDefinitionsForLLM()
		logger.Infof("Input: %s, %s, %+v", renderedSystemPrompt, renderedUserPrompt, functions)

		content, calls, errFunc := llm.AskLLMWithFunctions(renderedSystemPrompt, renderedUserPrompt, functions)
		if errFunc != nil {
			logger.Errorf("Failed during function-calling LLM request: %v", errFunc)
			fmt.Printf("Error generating blog (function calling): %v\n", errFunc)
			return
		}
		// log the input and output of AskLLMWithFunctions

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
				fmt.Printf("Error executing function %s: %v\n", call.Name, execErr)
				return
			}

			// Provide function result back to LLM in a second turn
			resultJSON, _ := json.Marshal(fnResult)
			followUpUser := renderedUserPrompt + "\n\n" + fmt.Sprintf("Function result for %s: %s\nUse this real weather data to complete the blog accurately.", call.Name, string(resultJSON))
			logger.Infof("Follow-up user prompt: %s", followUpUser)
			finalContent, errSecond := llm.AskLLM(renderedSystemPrompt, followUpUser)
			if errSecond != nil {
				logger.Errorf("Failed during follow-up LLM request: %v", errSecond)
				fmt.Printf("Error generating blog (follow-up): %v\n", errSecond)
				return
			}
			result = finalContent
		} else {
			// No function call needed; use content directly
			result = content
		}
	} else {
		// No city: plain call
		var resultBuilder strings.Builder
		err = llm.AskLLMWithStream(renderedSystemPrompt, renderedUserPrompt, func(chunk string) {
			resultBuilder.WriteString(chunk)
		})
		if err == nil {
			result = resultBuilder.String()
		}
	}
	if err != nil {
		logger.Errorf("Failed to generate blog: %v", err)
		fmt.Printf("Error generating blog: %v\n", err)
		return
	}

	// Output result
	if blogOutput != "" {
		// Write to file
		err := writeToFile(blogOutput, result)
		if err != nil {
			logger.Errorf("Failed to write to file: %v", err)
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
		fmt.Printf("Blog generated successfully and saved to: %s\n", blogOutput)
	} else {
		// Print to stdout
		fmt.Println(strings.Repeat("=", 80))
		fmt.Printf("Generated Blog - %s\n", todayStr)
		fmt.Println(strings.Repeat("=", 80))
		fmt.Println(result)
	}

	logger.Infof("Blog generation completed successfully")
}

func writeToFile(filepath, content string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
