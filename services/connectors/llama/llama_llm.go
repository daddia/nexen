package llama

import (
	"context"
	"fmt"

	"github.com/nexen/models"
	"github.com/nexen/services/connectors"
	"github.com/nexen/services/connectors/common"
)

const (
	defaultLlamaEndpoint = "http://localhost:8080/v1"
)

var (
	// List of model patterns the Llama connector supports
	supportedModelPatterns = []string{
		"llama-.*",
	}
)

// LlamaClient implements the LLM interface for locally hosted Llama models.
type LlamaClient struct {
	config    *common.LLMConfig
	modelName string
	// We would include the actual Llama client here in a real implementation
	// client *llama.Client
}

// init registers this adapter with the connectors registry.
func init() {
	for _, pattern := range supportedModelPatterns {
		connectors.Register(pattern, NewLlamaClient)
	}
}

// NewLlamaClient creates a new Llama client for the given model name.
func NewLlamaClient(model string, opts ...common.Option) (common.LLM, error) {
	config := common.DefaultLLMConfig()

	// Apply provided options
	if err := common.ApplyOptions(config, opts...); err != nil {
		return nil, fmt.Errorf("applying options: %w", err)
	}

	// For locally-hosted models, we might not need an API key
	// but we do need a valid endpoint
	if config.EndpointOverride == "" {
		config.EndpointOverride = defaultLlamaEndpoint
	}

	return &LlamaClient{
		config:    config,
		modelName: model,
		// In a real implementation, we would initialize the Llama client here
	}, nil
}

// Call implements the LLM interface Call method.
func (c *LlamaClient) Call(ctx context.Context, request *models.LLMRequest) (*models.LLMResponse, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// In a real implementation, we would:
	// 1. Transform the models.LLMRequest to Llama's request format
	// 2. Call the Llama API
	// 3. Transform the response to models.LLMResponse
	// 4. Handle errors, retries, and streaming if requested

	// For this example, we'll return a mock response
	mockResponse := &models.GenerateContentResponse{
		Candidates: []models.Candidate{
			{
				Content: &models.Content{
					Role:    "assistant",
					Message: fmt.Sprintf("This is a mock response from %s", c.modelName),
				},
				FinishReason: "stop",
			},
		},
		Usage: models.UsageMetrics{
			PromptTokens:     80,
			CompletionTokens: 30,
			TotalTokens:      110,
			LatencyMs:        1200, // Local models might be slower
			CostCents:        0,    // Local models typically have no per-token cost
		},
	}

	return &models.LLMResponse{
		Content: mockResponse.Candidates[0].Content,
		Usage:   mockResponse.Usage,
	}, nil
}

// BatchCall implements the LLM interface BatchCall method.
func (c *LlamaClient) BatchCall(ctx context.Context, requests []*models.LLMRequest) ([]*models.LLMResponse, error) {
	responses := make([]*models.LLMResponse, len(requests))
	var err error

	// Process each request sequentially
	// Local models typically can't handle many parallel requests
	for i, req := range requests {
		responses[i], err = c.Call(ctx, req)
		if err != nil {
			return responses, fmt.Errorf("error processing request %d: %w", i, err)
		}
	}

	return responses, nil
}

// SupportedModels returns a list of model names supported by this client.
func (c *LlamaClient) SupportedModels() []string {
	// In a real implementation, we might query the local server
	// to see what models are available
	return []string{
		"llama-7b",
		"llama-13b",
		"llama-70b",
	}
}
