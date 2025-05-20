package common

import (
	"context"

	"github.com/nexen/models"
)

// Option defines a functional option for configuring LLM instances.
type Option func(config *LLMConfig) error

// LLMConfig holds configuration parameters for LLM instances.
type LLMConfig struct {
	// APIKey is the authentication key for the provider.
	APIKey string

	// OrgID is the organization identifier for the provider.
	OrgID string

	// EndpointOverride allows using custom endpoints.
	EndpointOverride string

	// Timeout specifies the request timeout in seconds.
	Timeout int

	// RetryConfig controls retry behavior.
	RetryConfig RetryConfig

	// RegionRouting controls endpoint region selection.
	RegionRouting RegionRouting

	// CustomOptions contains provider-specific options.
	CustomOptions map[string]interface{}
}

// RetryConfig defines retry behavior for failed requests.
type RetryConfig struct {
	// MaxRetries is the number of times to retry a failed request.
	MaxRetries int

	// MinBackoff is the minimum backoff time in milliseconds.
	MinBackoff int

	// MaxBackoff is the maximum backoff time in milliseconds.
	MaxBackoff int

	// StatusCodesToRetry lists HTTP status codes that should trigger a retry.
	StatusCodesToRetry []int
}

// RegionRouting defines region selection strategy.
type RegionRouting struct {
	// EnableRegionRouting enables routing to different regions.
	EnableRegionRouting bool

	// PreferredRegions lists regions in order of preference.
	PreferredRegions []string

	// FailoverStrategy defines failover behavior (round-robin, etc.).
	FailoverStrategy string
}

// LLM defines the core interface for interacting with language models.
// All provider-specific implementations must satisfy this interface.
type LLM interface {
	// Call sends a request to the LLM and returns a response.
	Call(ctx context.Context, request *models.LLMRequest) (*models.LLMResponse, error)

	// BatchCall processes multiple requests and returns corresponding responses.
	BatchCall(ctx context.Context, requests []*models.LLMRequest) ([]*models.LLMResponse, error)

	// SupportedModels returns a list of model IDs that this implementation can handle.
	SupportedModels() []string
}

// WithAPIKey sets the API key option.
func WithAPIKey(apiKey string) Option {
	return func(config *LLMConfig) error {
		config.APIKey = apiKey
		return nil
	}
}

// WithOrgID sets the organization ID option.
func WithOrgID(orgID string) Option {
	return func(config *LLMConfig) error {
		config.OrgID = orgID
		return nil
	}
}

// WithEndpoint sets a custom endpoint URL.
func WithEndpoint(endpoint string) Option {
	return func(config *LLMConfig) error {
		config.EndpointOverride = endpoint
		return nil
	}
}

// WithTimeout sets the request timeout in seconds.
func WithTimeout(timeoutSec int) Option {
	return func(config *LLMConfig) error {
		config.Timeout = timeoutSec
		return nil
	}
}

// WithRetryConfig sets the retry configuration.
func WithRetryConfig(maxRetries, minBackoff, maxBackoff int, statusCodes []int) Option {
	return func(config *LLMConfig) error {
		config.RetryConfig = RetryConfig{
			MaxRetries:         maxRetries,
			MinBackoff:         minBackoff,
			MaxBackoff:         maxBackoff,
			StatusCodesToRetry: statusCodes,
		}
		return nil
	}
}

// WithRegionRouting sets region routing configuration.
func WithRegionRouting(enable bool, regions []string, strategy string) Option {
	return func(config *LLMConfig) error {
		config.RegionRouting = RegionRouting{
			EnableRegionRouting: enable,
			PreferredRegions:    regions,
			FailoverStrategy:    strategy,
		}
		return nil
	}
}

// WithCustomOption sets a provider-specific custom option.
func WithCustomOption(key string, value interface{}) Option {
	return func(config *LLMConfig) error {
		if config.CustomOptions == nil {
			config.CustomOptions = make(map[string]interface{})
		}
		config.CustomOptions[key] = value
		return nil
	}
}
