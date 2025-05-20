package custom

import (
	"context"
	"fmt"

	"github.com/nexen/models"
	"github.com/nexen/services/connectors"
	"github.com/nexen/services/connectors/common"
)

var (
	// List of model patterns the Custom connector supports
	supportedModelPatterns = []string{
		"custom-.*",
	}
)

// CustomClient implements the LLM interface for custom endpoints.
type CustomClient struct {
	config    *common.LLMConfig
	modelName string
	// We would include an HTTP client or specific client here
	// client *http.Client
}

// init registers this adapter with the connectors registry.
func init() {
	for _, pattern := range supportedModelPatterns {
		connectors.Register(pattern, NewCustomClient)
	}
}

// NewCustomClient creates a new custom client for the given model name.
func NewCustomClient(model string, opts ...common.Option) (common.LLM, error) {
	config := common.DefaultLLMConfig()

	// Apply provided options
	if err := common.ApplyOptions(config, opts...); err != nil {
		return nil, fmt.Errorf("applying options: %w", err)
	}

	// Validate required config - custom models must have endpoint overrides
	if config.EndpointOverride == "" {
		return nil, fmt.Errorf("custom model requires EndpointOverride to be set")
	}

	return &CustomClient{
		config:    config,
		modelName: model,
		// In a real implementation, we would initialize the HTTP client here
	}, nil
}

// Call implements the LLM interface Call method.
func (c *CustomClient) Call(ctx context.Context, request *models.LLMRequest) (*models.LLMResponse, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// In a real implementation, we would:
	// 1. Transform the models.LLMRequest to the format expected by the custom endpoint
	// 2. Call the custom API
	// 3. Transform the response to models.LLMResponse
	// 4. Handle errors, retries, and streaming if requested

	// For this example, we'll return a mock response
	mockResponse := &models.GenerateContentResponse{
		Candidates: []models.Candidate{
			{
				Content: &models.Content{
					Role: "assistant",
					Message: fmt.Sprintf("This is a custom response from %s at %s",
						c.modelName, c.config.EndpointOverride),
				},
				FinishReason: "stop",
			},
		},
		Usage: models.UsageMetrics{
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
			LatencyMs:        800,
			CostCents:        0, // Custom models typically don't have per-token costs
		},
	}

	return &models.LLMResponse{
		Content: mockResponse.Candidates[0].Content,
		Usage:   mockResponse.Usage,
	}, nil
}

// BatchCall implements the LLM interface BatchCall method.
func (c *CustomClient) BatchCall(ctx context.Context, requests []*models.LLMRequest) ([]*models.LLMResponse, error) {
	responses := make([]*models.LLMResponse, len(requests))
	var err error

	// Process each request sequentially
	for i, req := range requests {
		responses[i], err = c.Call(ctx, req)
		if err != nil {
			return responses, fmt.Errorf("error processing request %d: %w", i, err)
		}
	}

	return responses, nil
}

// SupportedModels returns a list of model names supported by this client.
func (c *CustomClient) SupportedModels() []string {
	// Custom models could be anything, so we just return a generic list
	return []string{
		"custom-model",
	}
}
