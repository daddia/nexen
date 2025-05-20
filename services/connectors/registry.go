package connectors

import (
    "fmt"
    "regexp"
    "sync"
)

// LLM represents the generic interface for any LLM client.
// Concrete implementations should satisfy this interface.
type LLM interface {
    // Call sends the request and returns a response or error.
    Call(request any) (any, error)
}

// constructorFn defines a function that creates an LLM given a model name and config.
type constructorFn func(model string, opts ...Option) (LLM, error)

// registry holds mappings from model-name regexes to LLM constructors.
var (
    mu               sync.RWMutex
    registry         = make(map[string]constructorFn)
    resolveCache     = make(map[string]constructorFn)
)

// Register associates a model-name regex with an LLM constructor.
// Call this in each connector's init() function or setup.
func Register(modelRegex string, constructor constructorFn) error {
    mu.Lock()
    defer mu.Unlock()
    if _, exists := registry[modelRegex]; exists {
        // Overwriting existing registration
    }
    registry[modelRegex] = constructor
    // clear cache so new registrations are considered
    resolveCache = make(map[string]constructorFn)
    return nil
}

// Resolve returns the constructor for the given model name, matching against registered regexes.
// It caches resolved constructors for performance.
func Resolve(model string) (constructorFn, error) {
    mu.RLock()
    if ctor, cached := resolveCache[model]; cached {
        mu.RUnlock()
        return ctor, nil
    }
    mu.RUnlock()

    mu.Lock()
    defer mu.Unlock()
    // Double-check cache under write-lock
    if ctor, cached := resolveCache[model]; cached {
        return ctor, nil
    }

    for regex, ctor := range registry {
        matched, err := regexp.MatchString(regex, model)
        if err != nil {
            return nil, fmt.Errorf("invalid regex %s: %w", regex, err)
        }
        if matched {
            resolveCache[model] = ctor
            return ctor, nil
        }
    }
    return nil, fmt.Errorf("no LLM constructor found for model %s", model)
}

// NewLLM creates an LLM instance for the given model name using the resolved constructor.
func NewLLM(model string, opts ...Option) (LLM, error) {
    ctor, err := Resolve(model)
    if err != nil {
        return nil, err
    }
    return ctor(model, opts...)
}
