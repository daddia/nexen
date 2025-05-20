// libs/logging/log_test.go
package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// Helper function to reset logger for each test
func setupTest(t *testing.T) *bytes.Buffer {
	// Reset buffer for output capture
	buf := &bytes.Buffer{}

	// Reset to default config with our buffer
	defaultCfg := DefaultConfig
	defaultCfg.Output = buf
	Configure(defaultCfg)

	return buf
}

func TestBasicLoggingOutputsJSON(t *testing.T) {
	// Setup fresh logger
	buf := setupTest(t)

	// Use a child logger with a known service name
	log := WithService("test-svc")
	log.Info().Str("foo", "bar").Msg("hello")

	out := buf.String()
	if !strings.Contains(out, `"service":"test-svc"`) {
		t.Errorf("expected service field, got %q", out)
	}
	if !strings.Contains(out, `"foo":"bar"`) {
		t.Errorf("expected foo field, got %q", out)
	}
	if !strings.Contains(out, `"message":"hello"`) {
		t.Errorf("expected message, got %q", out)
	}

	// Ensure valid JSON
	var obj map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}
}

func TestSetLevelSuppressesLowerLogs(t *testing.T) {
	buf := setupTest(t)
	SetLevel(zerolog.WarnLevel)

	Logger.Info().Msg("this should not appear")
	Logger.Warn().Msg("this should appear")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 log line, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[0], `"level":"warn"`) {
		t.Errorf("expected warn level, got %q", lines[0])
	}

	// Reset level for other tests
	SetLevel(zerolog.InfoLevel)
}

func TestConfigureWithOptions(t *testing.T) {
	buf := &bytes.Buffer{}

	// Configure with custom options
	cfg := Config{
		Level:                zerolog.DebugLevel,
		ConsoleOutput:        false,
		WithCaller:           true,
		CallerSkipFrameCount: 1,
		Output:               buf,
	}
	Configure(cfg)

	// Ensure global level is set
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	Logger.Debug().Msg("debug message")

	out := buf.String()
	if !strings.Contains(out, `"level":"debug"`) {
		t.Errorf("expected debug level, got %q", out)
	}

	// Should include caller info
	if !strings.Contains(out, `"caller":"`) {
		t.Errorf("expected caller field, got %q", out)
	}

	// Reset to default for other tests
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	setupTest(t)
}

func TestContextSupport(t *testing.T) {
	buf := &bytes.Buffer{}

	// Create a logger with buffer output
	cfg := Config{
		Level:      zerolog.InfoLevel,
		Output:     buf,
		WithCaller: false, // Disable caller for cleaner output
	}
	Configure(cfg)

	// Create a logger and attach to context
	log := WithService("context-test")
	ctx := WithContext(context.Background(), log)

	// Get logger from context and log with it
	contextLogger := FromContext(ctx)
	contextLogger.Info().Str("source", "context").Msg("context test")

	out := buf.String()
	if !strings.Contains(out, `"service":"context-test"`) {
		t.Errorf("expected service field from context, got %q", out)
	}
	if !strings.Contains(out, `"source":"context"`) {
		t.Errorf("expected source field, got %q", out)
	}

	// Reset for other tests
	setupTest(t)
}

func TestConsoleOutput(t *testing.T) {
	buf := &bytes.Buffer{}

	// Configure with console output
	cfg := Config{
		Level:         zerolog.InfoLevel,
		ConsoleOutput: true,
		Output:        buf,
	}
	Configure(cfg)

	Logger.Info().Str("key", "value").Msg("console test")

	out := buf.String()
	// Console output should contain our message and data in a readable format
	if !strings.Contains(out, "console test") {
		t.Errorf("missing message in console output: %q", out)
	}
	if !strings.Contains(out, "key=") {
		t.Errorf("missing field formatting in console output: %q", out)
	}

	// Reset for other tests
	setupTest(t)
}

func TestCallerHelper(t *testing.T) {
	caller := Caller(0)
	if !strings.Contains(caller, "log_test.go:") {
		t.Errorf("expected caller to contain log_test.go, got %q", caller)
	}
}
