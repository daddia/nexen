# Logging Library (`libs/logging`)

A shared, pre-configured logging module based on [zerolog](https://github.com/rs/zerolog), providing structured, high-performance JSON logging for all services in the Nexen platform.

## Features

* **Zero-allocation JSON logs** for minimal overhead in production.
* **Global logger** configured with RFC3339 timestamps and stderr output by default.
* **Dynamic log levels** via `SetLevel(level zerolog.Level)`.
* **Custom output** destinations via `SetOutput(io.Writer)` (files, buffers, etc.).
* **Service tagging** helper `WithService(name string)` to include a `service` field in every entry.
* **Comprehensive configuration** with `Configure(Config)` for setting up production-ready logging.
* **Context integration** for passing loggers through your application contexts.
* **Source line tracing** to identify the exact source of log entries.
* **Human-readable console format** available for development environments.
* **Test coverage** with `log_test.go` verifying JSON structure and level filtering.

## Installation

1. Ensure your project's `go.work` or `go.mod` includes the `libs/logging` module:

   ```bash
   go get github.com/yourorg/nexen/libs/logging@latest
   ```

2. Import in your Go code:
   ```go
   import (
    "github.com/yourorg/nexen/libs/logging"
    "github.com/rs/zerolog"
   )
   ```

## Usage

### Basic Usage

```go
package main

import (
    "github.com/rs/zerolog"
    "github.com/yourorg/nexen/libs/logging"
)

func main() {
    // Service-scoped logger
    log := logging.WithService("gateway")

    log.Info().Str("user", "alice").Msg("User login event")
    log.Error().Err(err).Msg("Operation failed")
}
```

### Advanced Configuration

```go
package main

import (
    "os"
    "github.com/rs/zerolog"
    "github.com/yourorg/nexen/libs/logging"
)

func main() {
    // Configure with custom options for production
    logging.Configure(logging.Config{
        Level:         zerolog.InfoLevel,
        ConsoleOutput: false,  // Use JSON format for production
        WithCaller:    true,   // Add source file:line to logs
        Output:        os.Stdout,
    })

    // For development, you can enable console output
    if os.Getenv("ENV") == "development" {
        logging.EnableConsoleOutput()
    }

    // Create service logger and use
    log := logging.WithService("api-service")
    log.Info().Str("version", "1.2.3").Msg("Service started")
}
```

### Context Integration

```go
package main

import (
    "context"
    "net/http"

    "github.com/yourorg/nexen/libs/logging"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Create a logger with request info
    log := logging.WithService("api")
        .With()
        .Str("method", r.Method)
        .Str("path", r.URL.Path)
        .Str("remote", r.RemoteAddr)
        .Logger()

    // Store in context
    ctx := logging.WithContext(r.Context(), log)

    // Pass context to other functions
    processRequest(ctx, r)
}

func processRequest(ctx context.Context, r *http.Request) {
    // Get logger from context
    log := logging.FromContext(ctx)
    log.Info().Msg("Processing request")

    // Do work...
}
```

## API Reference

### Core Functions

* **`func SetLevel(level zerolog.Level)`**
  Set the global minimum log level (e.g. `zerolog.DebugLevel`, `zerolog.InfoLevel`).

* **`func SetOutput(w io.Writer)`**
  Override the log output destination (e.g. to a file or buffer).

* **`func WithService(name string) zerolog.Logger`**
  Returns a child logger with the `service` field set to `name`.

* **`var Logger zerolog.Logger`**
  The base logger instance; you can use it directly for simple, unscoped logs.

### Production Features

* **`func Configure(cfg Config)`**
  Configure the logger with production-ready settings using the Config struct.

* **`func EnableConsoleOutput()`**
  Switch to human-readable console output for development environments.

* **`func FromContext(ctx context.Context) zerolog.Logger`**
  Retrieve a logger from a context, or return the default logger if none exists.

* **`func WithContext(ctx context.Context, logger zerolog.Logger) context.Context`**
  Store a logger in a context for passing through your application.

* **`func Caller(skip int) string`**
  Helper to get file:line information for custom usage.

### Configuration Options

```go
type Config struct {
    // Level sets the minimum log level
    Level zerolog.Level
    // ConsoleOutput determines if logs should be formatted for human readability
    ConsoleOutput bool
    // WithCaller adds caller information (file:line) to logs
    WithCaller bool
    // CallerSkipFrameCount sets the number of stack frames to skip when capturing caller info
    CallerSkipFrameCount int
    // Output sets the destination for log entries
    Output io.Writer
}
```

## Testing

Run unit tests to verify logging behavior:

```bash
cd libs/logging
go test ./... -v
```

The tests check that:

* Logs are emitted as valid JSON.
* `service` and custom fields appear correctly.
* Log level filtering works as expected.
* Context integration functions correctly.
* Console output formatting works as expected.
* Caller information is correctly captured.

---
