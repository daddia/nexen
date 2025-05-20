// libs/logging/log.go
package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

// Config holds logging configuration options
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

// DefaultConfig provides sensible defaults for production logging
var DefaultConfig = Config{
	Level:                zerolog.InfoLevel,
	ConsoleOutput:        false,
	WithCaller:           true,
	CallerSkipFrameCount: 2,
	Output:               os.Stderr,
}

// Logger is the package‐level zerolog.Logger instance pre‐configured with
// timestamping.  Call SetLevel or SetOutput to customise.
var Logger zerolog.Logger

// Initialize the logger with default settings
func init() {
	// use RFC3339 for human‐readable timestamps
	zerolog.TimeFieldFormat = time.RFC3339
	// Set default global level
	zerolog.SetGlobalLevel(DefaultConfig.Level)
	// Initialize with defaults
	Configure(DefaultConfig)
}

// Configure sets up the logger with the given configuration
func Configure(cfg Config) {
	output := cfg.Output
	if output == nil {
		output = os.Stderr
	}

	// Apply console formatting if requested
	if cfg.ConsoleOutput {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
	}

	// Build the logger
	ctx := zerolog.New(output).With().Timestamp()

	// Add caller info if requested
	if cfg.WithCaller {
		ctx = ctx.CallerWithSkipFrameCount(cfg.CallerSkipFrameCount)
	}

	Logger = ctx.Logger()
}

// SetLevel sets the global logging level (e.g. zerolog.DebugLevel).
func SetLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

// SetOutput redirects logs to the given io.Writer (e.g. a file or buffer).
func SetOutput(w io.Writer) {
	Logger = Logger.Output(w)
}

// WithService returns a child logger with a "service" field pre‐populated.
func WithService(name string) zerolog.Logger {
	return Logger.With().Str("service", name).Logger()
}

// FromContext retrieves a logger from context or returns the default logger
func FromContext(ctx context.Context) zerolog.Logger {
	if ctx == nil {
		return Logger
	}
	if l := zerolog.Ctx(ctx); l != nil && l != zerolog.DefaultContextLogger {
		return *l
	}
	return Logger
}

// WithContext returns a new context with the logger attached
func WithContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}

// EnableConsoleOutput switches to human-readable console output format
func EnableConsoleOutput() {
	Logger = Logger.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})
}

// Caller returns a formatted string with file:line of caller (for custom usage)
func Caller(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown:0"
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}
