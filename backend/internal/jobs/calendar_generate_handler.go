package jobs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/walterfan/lazy-rabbit-reminder/pkg/llm"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/log"
)

// CalendarGenerateHandler implements JobHandler for generating calendars
type CalendarGenerateHandler struct {
	jobManager *JobManager
}

// Execute generates daily calendar content
func (h *CalendarGenerateHandler) Execute(params string) error {
	h.jobManager.logger.Info("Generating calendar content...")

	// Create template data for calendar generation
	today := time.Now()
	data := map[string]interface{}{
		"date":  today.Format("2006-01-02"),
		"month": today.Format("January"),
		"year":  today.Format("2006"),
		"day":   today.Format("Monday"),
	}

	// Generate calendar content using simple LLM call
	content, err := h.generateSimpleContentWithLLM(
		"You are a calendar generator. Create useful daily calendar content with events, reminders, and scheduling suggestions.",
		fmt.Sprintf("Generate calendar content for %s (%s, %s %s). Include suggested daily structure, important reminders, and productivity tips.",
			data["date"], data["day"], data["month"], data["year"]),
	)
	if err != nil {
		h.jobManager.logger.Errorf("Failed to generate calendar content: %v", err)
		return fmt.Errorf("failed to generate calendar content: %w", err)
	}

	// Save the generated calendar content
	filename := fmt.Sprintf("calendar_%s.md", today.Format("2006-01-02"))
	if err := h.saveCalendarContent(filename, content); err != nil {
		h.jobManager.logger.Errorf("Failed to save calendar content: %v", err)
		return fmt.Errorf("failed to save calendar content: %w", err)
	}

	h.jobManager.logger.Infof("Calendar generation completed successfully: %s", filename)
	return nil
}

// generateSimpleContentWithLLM generates content using simple LLM call
func (h *CalendarGenerateHandler) generateSimpleContentWithLLM(systemPrompt, userPrompt string) (string, error) {
	// Use the callRealLLM method for consistency
	return h.callRealLLM(systemPrompt, userPrompt)
}

// callRealLLM calls the actual LLM package
func (h *CalendarGenerateHandler) callRealLLM(systemPrompt, userPrompt string) (string, error) {
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

// saveCalendarContent saves calendar content to file
func (h *CalendarGenerateHandler) saveCalendarContent(filename, content string) error {
	// Create calendars directory if it doesn't exist
	if err := os.MkdirAll("calendars", 0755); err != nil {
		return fmt.Errorf("failed to create calendars directory: %w", err)
	}

	filepath := fmt.Sprintf("calendars/%s", filename)
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write calendar file: %w", err)
	}

	return nil
}
