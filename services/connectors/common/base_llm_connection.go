package common

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

// DefaultRetryStatusCodes defines the standard HTTP status codes to retry.
var DefaultRetryStatusCodes = []int{
	http.StatusTooManyRequests,     // 429
	http.StatusInternalServerError, // 500
	http.StatusBadGateway,          // 502
	http.StatusServiceUnavailable,  // 503
	http.StatusGatewayTimeout,      // 504
}

// DefaultTimeoutSeconds is the default request timeout (30 seconds).
const DefaultTimeoutSeconds = 30

// DefaultRetryConfig provides sensible defaults for retry behavior.
var DefaultRetryConfig = RetryConfig{
	MaxRetries:         3,
	MinBackoff:         100,  // 100ms
	MaxBackoff:         5000, // 5s
	StatusCodesToRetry: DefaultRetryStatusCodes,
}

// DefaultLLMConfig provides a default configuration for LLM clients.
func DefaultLLMConfig() *LLMConfig {
	return &LLMConfig{
		Timeout:     DefaultTimeoutSeconds,
		RetryConfig: DefaultRetryConfig,
		RegionRouting: RegionRouting{
			EnableRegionRouting: false,
			FailoverStrategy:    "sequential",
		},
		CustomOptions: make(map[string]interface{}),
	}
}

// ApplyOptions applies options to a configuration.
func ApplyOptions(config *LLMConfig, opts ...Option) error {
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return err
		}
	}
	return nil
}

// NewHTTPClientWithTimeout creates an HTTP client with the specified timeout.
func NewHTTPClientWithTimeout(timeoutSec int) *http.Client {
	if timeoutSec <= 0 {
		timeoutSec = DefaultTimeoutSeconds
	}

	return &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}
}

// CalculateBackoff determines the backoff duration for a retry attempt.
func CalculateBackoff(attempt int, config RetryConfig) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	// Exponential backoff with jitter: min(maxBackoff, minBackoff * 2^attempt) with 20% jitter
	backoff := float64(config.MinBackoff) * math.Pow(2, float64(attempt))
	backoff = math.Min(backoff, float64(config.MaxBackoff))

	// Add jitter (±20%)
	jitter := backoff * 0.2 * (2*rand() - 1)
	backoff = backoff + jitter

	return time.Duration(backoff) * time.Millisecond
}

// rand returns a random float in [0,1)
func rand() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}

// ShouldRetry determines if a request should be retried based on status code.
func ShouldRetry(statusCode int, config RetryConfig) bool {
	for _, code := range config.StatusCodesToRetry {
		if statusCode == code {
			return true
		}
	}
	return false
}

// WithContext applies a context timeout to an existing context.
func WithContext(parent context.Context, timeoutSec int) (context.Context, context.CancelFunc) {
	if timeoutSec <= 0 {
		timeoutSec = DefaultTimeoutSeconds
	}
	return context.WithTimeout(parent, time.Duration(timeoutSec)*time.Second)
}

// CreateEndpointURL constructs an endpoint URL based on region settings.
func CreateEndpointURL(baseEndpoint string, config *LLMConfig) string {
	if config.EndpointOverride != "" {
		return config.EndpointOverride
	}

	if !config.RegionRouting.EnableRegionRouting || len(config.RegionRouting.PreferredRegions) == 0 {
		return baseEndpoint
	}

	// Here we would implement region routing logic based on the FailoverStrategy
	// For simplicity, we'll just use the first preferred region in this example
	region := config.RegionRouting.PreferredRegions[0]

	// This is a simplified example; actual implementation would depend on provider URL structure
	return fmt.Sprintf("%s/%s", baseEndpoint, region)
}
