# Unified Logger Package

This package provides a unified, configurable logging solution for the Lazy Rabbit Secretary application using Uber's Zap logger.

## Features

- **Configurable Log Levels**: debug, info, warn, error, fatal, panic
- **Multiple Output Formats**: JSON (structured) or Console (human-readable)
- **Dual Output**: Console and file logging
- **Configuration Sources**: Config file, environment variables, and programmatic configuration
- **Structured Logging**: Support for key-value pairs and contextual information
- **Performance**: Built on Zap for high-performance logging

## Configuration

### Config File (config.yaml)

```yaml
log:
  file: "lazy-rabbit-secretary.log"  # Log file path
  level: "info"                     # Log level: debug, info, warn, error, fatal, panic
  format: "json"                    # Output format: json or console
  enable_file: true                 # Enable/disable file logging
  max_size: 100                     # Maximum log file size in MB
```

### Environment Variables

Environment variables take precedence over config file settings:

- `LOG_LEVEL`: Set log level (debug, info, warn, error, fatal, panic)
- `LOG_FILE`: Set log file path
- `LOG_FORMAT`: Set output format (json, console)

### Programmatic Configuration

```go
import "github.com/walterfan/lazy-rabbit-secretary/pkg/log"

// Custom configuration
config := &log.LogConfig{
    Level:      "debug",
    File:       "custom.log",
    Format:     "console",
    EnableFile: true,
    MaxSize:    50,
}

err := log.InitLoggerWithConfig(config)
if err != nil {
    panic(err)
}
```

## Usage

### Basic Logging

```go
import "github.com/walterfan/lazy-rabbit-secretary/pkg/log"

// Initialize logger (usually done once at application startup)
err := log.InitLogger()
if err != nil {
    panic(err)
}

// Get logger instance
logger := log.GetLogger()

// Basic logging
logger.Debug("Debug message")
logger.Info("Info message")
logger.Warn("Warning message")
logger.Error("Error message")
logger.Fatal("Fatal message") // Exits the program
```

### Structured Logging

```go
// Structured logging with key-value pairs
logger.Infow("User logged in",
    "user_id", "12345",
    "ip_address", "192.168.1.100",
    "timestamp", time.Now(),
    "success", true,
)

// Formatted logging
logger.Infof("Processing %d items for user %s", itemCount, userID)
```

### Error Logging with Stack Traces

```go
if err != nil {
    logger.Errorw("Database operation failed",
        "operation", "user_create",
        "error", err,
        "user_id", userID,
    )
}
```

## Log Formats

### JSON Format (Production)

```json
{
  "level": "info",
  "timestamp": "2025-09-17T14:13:49.587+0800",
  "caller": "service/user.go:45",
  "msg": "User created successfully",
  "user_id": "12345",
  "email": "user@example.com"
}
```

### Console Format (Development)

```
2025-09-17T14:13:49.587+0800    INFO    service/user.go:45    User created successfully    {"user_id": "12345", "email": "user@example.com"}
```

## Log Levels

| Level | Description | Use Case |
|-------|-------------|----------|
| `debug` | Detailed information for debugging | Development, troubleshooting |
| `info` | General information about program execution | Normal operations, milestones |
| `warn` | Warning messages for potentially harmful situations | Recoverable errors, deprecated usage |
| `error` | Error messages for error conditions | Errors that don't stop the program |
| `fatal` | Fatal messages that cause program termination | Critical errors, startup failures |
| `panic` | Panic messages that cause panic | Unrecoverable errors |

## Best Practices

### 1. Use Appropriate Log Levels

```go
// Good: Use info for normal operations
logger.Info("Server started successfully")

// Good: Use error for recoverable errors
logger.Errorw("Failed to send email", "error", err, "recipient", email)

// Good: Use debug for detailed debugging info
logger.Debugw("Processing request", "method", "POST", "path", "/api/users")
```

### 2. Include Context

```go
// Good: Include relevant context
logger.Infow("User action completed",
    "user_id", userID,
    "action", "profile_update",
    "duration_ms", duration.Milliseconds(),
)

// Avoid: Vague messages without context
logger.Info("Action completed")
```

### 3. Use Structured Logging

```go
// Good: Structured logging
logger.Infow("Database query executed",
    "table", "users",
    "operation", "SELECT",
    "duration_ms", 150,
    "rows_returned", 5,
)

// Avoid: String concatenation
logger.Infof("Database query on %s took %dms and returned %d rows", table, duration, rows)
```

### 4. Handle Errors Properly

```go
if err != nil {
    logger.Errorw("Operation failed",
        "operation", "user_creation",
        "error", err,
        "user_data", userData,
    )
    return err
}
```

## Environment-Specific Configuration

### Development
```yaml
log:
  level: "debug"
  format: "console"
  enable_file: false
```

### Production
```yaml
log:
  level: "info"
  format: "json"
  enable_file: true
  file: "/var/log/lazy-rabbit-secretary.log"
  max_size: 100
```

## Integration with Application

The logger is automatically initialized when the application starts via `cmd/root.go`. It reads configuration from:

1. Default values
2. Config file (`config.yaml`)
3. Environment variables (highest priority)

## Performance Considerations

- JSON format is optimized for log aggregation systems
- Console format is optimized for human readability
- Structured logging is more efficient than string formatting
- Debug level logging has minimal overhead when disabled

## Log Rotation

Currently, the logger supports basic file size limits. For production environments, consider using external log rotation tools like `logrotate` or implementing log rotation libraries.

## Monitoring and Alerting

The structured JSON format makes it easy to:
- Parse logs with tools like ELK stack, Fluentd, or Promtail
- Set up alerts based on error rates
- Create dashboards for application monitoring
- Perform log analysis and debugging
