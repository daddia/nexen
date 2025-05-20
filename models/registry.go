package models

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Common model profile constants
const (
	ProfileChat     = "chat"     // General chat capability
	ProfileThinking = "thinking" // Deep reasoning capability
	ProfileAgent    = "agent"    // Tool/function calling capability
	ProfileRAG      = "rag"      // Retrieval-augmented generation
	ProfileCreative = "creative" // Creative content generation
	ProfileCode     = "code"     // Code generation capability
)

// Provider constants for major LLM vendors
const (
	ProviderOpenAI    = "openai"
	ProviderAnthropic = "anthropic"
	ProviderGoogle    = "google"
	ProviderMistral   = "mistral"
	ProviderLlama     = "llama"
	ProviderCustom    = "custom"
)

// CostTier represents pricing categories
type CostTier string

const (
	CostTierBasic    CostTier = "basic"    // Low-cost, standard models
	CostTierStandard CostTier = "standard" // Mid-range models
	CostTierPremium  CostTier = "premium"  // High-end, expensive models
)

// ModelInfo holds metadata for an LLM model (ID, profiles, token limits, etc.).
type ModelInfo struct {
	// ID is the unique model identifier (e.g. "gpt-4-turbo").
	ID string `json:"id"`

	// Profiles lists capabilities (e.g. "chat", "thinking", "agent").
	Profiles []string `json:"profiles"`

	// MaxTokens is the model's maximum context window size.
	MaxTokens int `json:"maxTokens"`

	// CostPerToken is the price per token in cents.
	CostPerToken float64 `json:"costPerToken"`

	// Provider indicates the vendor (OpenAI, Anthropic, etc).
	Provider string `json:"provider"`

	// CostTier indicates pricing category (basic, standard, premium).
	CostTier CostTier `json:"costTier"`

	// Version is semantic version of the model if available.
	Version string `json:"version,omitempty"`
}

var (
	mu       sync.RWMutex
	registry = make(map[string]ModelInfo) // regex -> ModelInfo
	cache    = make(map[string]ModelInfo) // model name -> resolved ModelInfo
)

// Register registers a ModelInfo under a model-name regex pattern.
// regexPattern should be a valid Go regexp that matches model IDs.
func Register(regexPattern string, info ModelInfo) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[regexPattern]; exists {
		// Overwrite existing registration
	}
	// Validate regex compiles
	if _, err := regexp.Compile(regexPattern); err != nil {
		return fmt.Errorf("invalid regex %q: %w", regexPattern, err)
	}
	registry[regexPattern] = info
	// Clear cache to force re-resolve
	cache = make(map[string]ModelInfo)
	return nil
}

// Resolve returns the ModelInfo whose regex matches the given model name.
// It caches resolutions for performance.
func Resolve(model string) (ModelInfo, error) {
	mu.RLock()
	if info, found := cache[model]; found {
		mu.RUnlock()
		return info, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()
	// Double-check cache under write lock
	if info, found := cache[model]; found {
		return info, nil
	}
	for pattern, info := range registry {
		matched, err := regexp.MatchString(pattern, model)
		if err != nil {
			return ModelInfo{}, fmt.Errorf("invalid regex %q during resolve: %w", pattern, err)
		}
		if matched {
			// Create a copy with the exact ID that was requested
			resolvedInfo := info
			resolvedInfo.ID = model
			cache[model] = resolvedInfo
			return resolvedInfo, nil
		}
	}
	return ModelInfo{}, fmt.Errorf("model not found: %s", model)
}

// NewModelInfo is a helper to register multiple patterns at once.
func NewModelInfo(info ModelInfo, patterns ...string) error {
	for _, p := range patterns {
		if err := Register(p, info); err != nil {
			return err
		}
	}
	return nil
}

// ListModels returns a list of all registered model patterns.
func ListModels() []string {
	mu.RLock()
	defer mu.RUnlock()

	patterns := make([]string, 0, len(registry))
	for pattern := range registry {
		patterns = append(patterns, pattern)
	}
	return patterns
}

// ListModelsByProfile returns models that support a specific profile.
func ListModelsByProfile(profile string) []ModelInfo {
	mu.RLock()
	defer mu.RUnlock()

	var models []ModelInfo
	for _, info := range registry {
		for _, p := range info.Profiles {
			if p == profile {
				models = append(models, info)
				break
			}
		}
	}
	return models
}

// ListModelsByProvider returns models from a specific provider.
func ListModelsByProvider(provider string) []ModelInfo {
	mu.RLock()
	defer mu.RUnlock()

	var models []ModelInfo
	for _, info := range registry {
		if strings.EqualFold(info.Provider, provider) {
			models = append(models, info)
		}
	}
	return models
}

// HasProfile checks if a model supports a specific profile.
func HasProfile(model, profile string) (bool, error) {
	info, err := Resolve(model)
	if err != nil {
		return false, err
	}

	for _, p := range info.Profiles {
		if p == profile {
			return true, nil
		}
	}
	return false, nil
}

// ClearRegistry removes all model registrations.
// Primarily used for testing.
func ClearRegistry() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[string]ModelInfo)
	cache = make(map[string]ModelInfo)
}

// Init registers common models with the registry.
func Init() {
	// OpenAI models
	NewModelInfo(ModelInfo{
		ID:           "gpt-4-turbo",
		Profiles:     []string{ProfileChat, ProfileThinking, ProfileAgent, ProfileRAG},
		MaxTokens:    128000,
		CostPerToken: 0.00001,
		Provider:     ProviderOpenAI,
		CostTier:     CostTierPremium,
		Version:      "1.0",
	}, "gpt-4-turbo.*")

	NewModelInfo(ModelInfo{
		ID:           "gpt-4",
		Profiles:     []string{ProfileChat, ProfileThinking, ProfileAgent},
		MaxTokens:    8192,
		CostPerToken: 0.00003,
		Provider:     ProviderOpenAI,
		CostTier:     CostTierPremium,
		Version:      "1.0",
	}, "gpt-4$", "gpt-4-.*")

	NewModelInfo(ModelInfo{
		ID:           "gpt-3.5-turbo",
		Profiles:     []string{ProfileChat, ProfileAgent},
		MaxTokens:    16385,
		CostPerToken: 0.000002,
		Provider:     ProviderOpenAI,
		CostTier:     CostTierStandard,
		Version:      "1.0",
	}, "gpt-3.5-turbo.*")

	// Anthropic models
	NewModelInfo(ModelInfo{
		ID:           "claude-3-opus",
		Profiles:     []string{ProfileChat, ProfileThinking, ProfileRAG, ProfileCreative},
		MaxTokens:    200000,
		CostPerToken: 0.00002,
		Provider:     ProviderAnthropic,
		CostTier:     CostTierPremium,
		Version:      "1.0",
	}, "claude-3-opus.*")

	NewModelInfo(ModelInfo{
		ID:           "claude-3-sonnet",
		Profiles:     []string{ProfileChat, ProfileThinking, ProfileRAG},
		MaxTokens:    200000,
		CostPerToken: 0.00001,
		Provider:     ProviderAnthropic,
		CostTier:     CostTierStandard,
		Version:      "1.0",
	}, "claude-3-sonnet.*")

	// Google models
	NewModelInfo(ModelInfo{
		ID:           "gemini-pro",
		Profiles:     []string{ProfileChat, ProfileAgent, ProfileRAG},
		MaxTokens:    32768,
		CostPerToken: 0.000005,
		Provider:     ProviderGoogle,
		CostTier:     CostTierStandard,
		Version:      "1.0",
	}, "gemini-pro.*")

	// Mistral models
	NewModelInfo(ModelInfo{
		ID:           "mistral-large",
		Profiles:     []string{ProfileChat, ProfileThinking},
		MaxTokens:    32768,
		CostPerToken: 0.000008,
		Provider:     ProviderMistral,
		CostTier:     CostTierStandard,
		Version:      "1.0",
	}, "mistral-large.*")
}
