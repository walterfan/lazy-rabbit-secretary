package log_test

import (
	"os"
	"testing"

	"github.com/walterfan/lazy-rabbit-reminder/pkg/log"
)

func TestLoggerWithDifferentLevels(t *testing.T) {
	// Test with debug level
	config := &log.LogConfig{
		Level:      "debug",
		Format:     "console",
		EnableFile: false, // Disable file for testing
	}

	err := log.InitLoggerWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	logger := log.GetLogger()

	// Test different log levels
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Test with structured fields
	logger.Infow("User logged in",
		"user_id", "123",
		"ip_address", "192.168.1.1",
		"timestamp", "2025-09-17T14:00:00Z",
	)
}

func TestLoggerWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("LOG_LEVEL", "warn")
	os.Setenv("LOG_FORMAT", "json")
	defer os.Unsetenv("LOG_LEVEL")
	defer os.Unsetenv("LOG_FORMAT")

	err := log.InitLogger()
	if err != nil {
		t.Fatalf("Failed to initialize logger with env vars: %v", err)
	}

	logger := log.GetLogger()

	// Only warn and error should be visible
	logger.Debug("This debug message should not appear")
	logger.Info("This info message should not appear")
	logger.Warn("This warning message should appear")
	logger.Error("This error message should appear")
}

func TestLoggerJSONFormat(t *testing.T) {
	config := &log.LogConfig{
		Level:      "info",
		Format:     "json",
		EnableFile: false,
	}

	err := log.InitLoggerWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to initialize JSON logger: %v", err)
	}

	logger := log.GetLogger()

	logger.Infow("JSON formatted message",
		"component", "test",
		"action", "logging",
		"success", true,
		"count", 42,
	)
}
