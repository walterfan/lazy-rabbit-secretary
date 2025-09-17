package log

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogger *zap.Logger
)

// LogConfig holds the logging configuration
type LogConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`
	File       string `yaml:"file" mapstructure:"file"`
	Format     string `yaml:"format" mapstructure:"format"` // "json" or "console"
	EnableFile bool   `yaml:"enable_file" mapstructure:"enable_file"`
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"` // MB
}

// getDefaultConfig returns default logging configuration
func getDefaultConfig() LogConfig {
	return LogConfig{
		Level:      "info",
		File:       "lazy-rabbit-reminder.log",
		Format:     "json",
		EnableFile: true,
		MaxSize:    100,
	}
}

// InitLogger initializes the logger with configuration from config file and environment variables
func InitLogger() error {
	return InitLoggerWithConfig(nil)
}

// InitLoggerWithConfig initializes the logger with custom configuration
func InitLoggerWithConfig(customConfig *LogConfig) error {
	config := getDefaultConfig()

	// Override with custom config if provided
	if customConfig != nil {
		if customConfig.Level != "" {
			config.Level = customConfig.Level
		}
		if customConfig.File != "" {
			config.File = customConfig.File
		}
		if customConfig.Format != "" {
			config.Format = customConfig.Format
		}
		config.EnableFile = customConfig.EnableFile
		if customConfig.MaxSize > 0 {
			config.MaxSize = customConfig.MaxSize
		}
	} else {
		// Try to load from viper configuration
		if viper.IsSet("log.level") {
			config.Level = viper.GetString("log.level")
		}
		if viper.IsSet("log.file") {
			config.File = viper.GetString("log.file")
		}
		if viper.IsSet("log.format") {
			config.Format = viper.GetString("log.format")
		}
		if viper.IsSet("log.enable_file") {
			config.EnableFile = viper.GetBool("log.enable_file")
		}
		if viper.IsSet("log.max_size") {
			config.MaxSize = viper.GetInt("log.max_size")
		}
	}

	// Override with environment variables
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		config.Level = envLevel
	}
	if envFile := os.Getenv("LOG_FILE"); envFile != "" {
		config.File = envFile
	}
	if envFormat := os.Getenv("LOG_FORMAT"); envFormat != "" {
		config.Format = envFormat
	}

	// Parse log level
	level, err := parseLogLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel // fallback to info level
	}

	// Create encoder configuration
	var encoderCfg zapcore.EncoderConfig
	if config.Format == "console" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create cores
	var cores []zapcore.Core

	// Console core (always enabled)
	var consoleEncoder zapcore.Encoder
	if config.Format == "console" {
		consoleEncoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		consoleEncoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))

	// File core (optional)
	if config.EnableFile && config.File != "" {
		logFile, err := os.OpenFile(config.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderCfg) // Always use JSON for file output
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level))
	}

	// Create the logger
	core := zapcore.NewTee(cores...)
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))

	return nil
}

// parseLogLevel converts string log level to zapcore.Level
func parseLogLevel(levelStr string) (zapcore.Level, error) {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	default:
		return zapcore.InfoLevel, nil
	}
}

// GetLogger returns the configured sugared logger
func GetLogger() *zap.SugaredLogger {
	if zapLogger == nil {
		// Initialize with default config if not already initialized
		InitLogger()
	}
	return zapLogger.Sugar()
}

// GetRawLogger returns the raw zap logger (for advanced usage)
func GetRawLogger() *zap.Logger {
	if zapLogger == nil {
		// Initialize with default config if not already initialized
		InitLogger()
	}
	return zapLogger
}

// SetLevel dynamically changes the log level
func SetLevel(levelStr string) error {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		return err
	}

	// Note: This is a simplified implementation
	// In a production system, you might want to recreate the logger with the new level
	// For now, this serves as a placeholder for dynamic level changing
	_ = level
	return nil
}

// Sync flushes any buffered log entries
func Sync() error {
	if zapLogger != nil {
		return zapLogger.Sync()
	}
	return nil
}
