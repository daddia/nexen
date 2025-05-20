package anthropic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/nexen/models"
	"github.com/nexen/services/connectors"
	"github.com/nexen/services/connectors/common"
)

const (
	defaultAnthropicEndpoint = "https://api.anthropic.com/v1"
	defaultMaxTokens         = 4096
)

var (
	// List of model patterns the Anthropic connector supports
	supportedModelPatterns = []string{
		"claude-3.*",
		"claude-3-5.*",
	}
)

// AnthropicClient implements the LLM interface for Anthropic's API.
type AnthropicClient struct {
	config    *common.LLMConfig
	modelName string
	client    anthropic.Client
}

// init registers this adapter with the connectors registry.
func init() {
	for _, pattern := range supportedModelPatterns {
		connectors.Register(pattern, NewAnthropicClient)
	}
}

// NewAnthropicClient creates a new Anthropic client for the given model name.
func NewAnthropicClient(model string, opts ...common.Option) (common.LLM, error) {
	config := common.DefaultLLMConfig()

	// Apply provided options
	if err := common.ApplyOptions(config, opts...); err != nil {
		return nil, fmt.Errorf("applying options: %w", err)
	}

	// Validate required config
	if config.APIKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}

	// Initialize Anthropic client with API key
	clientOpts := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
	}

	// Set organization ID if provided
	if config.OrgID != "" {
		clientOpts = append(clientOpts, option.WithHeader("Anthropic-Organization", config.OrgID))
	}

	// Set custom endpoint if provided
	if config.EndpointOverride != "" {
		clientOpts = append(clientOpts, option.WithBaseURL(config.EndpointOverride))
	}

	// Set timeout
	if config.Timeout > 0 {
		clientOpts = append(clientOpts, option.WithRequestTimeout(time.Duration(config.Timeout)*time.Second))
	}

	// Set retry configuration
	if config.RetryConfig.MaxRetries > 0 {
		clientOpts = append(clientOpts, option.WithMaxRetries(config.RetryConfig.MaxRetries))
	}

	client := anthropic.NewClient(clientOpts...)

	return &AnthropicClient{
		config:    config,
		modelName: model,
		client:    client,
	}, nil
}

// mapToAnthropicModel maps model names to Anthropic SDK model constants
func mapToAnthropicModel(modelName string) string {
	// Define mapping for known models
	modelMap := map[string]string{
		"claude-3-opus":     string(anthropic.ModelClaude3OpusLatest),
		"claude-3-sonnet":   string(anthropic.ModelClaude_3_Sonnet_20240229),
		"claude-3-haiku":    string(anthropic.ModelClaude_3_Haiku_20240307),
		"claude-3.5-sonnet": string(anthropic.ModelClaude3_5SonnetLatest),
	}

	// If we have a direct mapping, use it
	if mappedModel, ok := modelMap[modelName]; ok {
		return mappedModel
	}

	// Default: return the original model name
	return modelName
}

// contentToMessageParams converts models.Content to anthropic.MessageParam
func contentToMessageParams(contents []models.Content) []anthropic.MessageParam {
	messages := make([]anthropic.MessageParam, 0, len(contents))

	for _, content := range contents {
		// Determine role
		role := anthropic.MessageParamRoleUser
		if content.Role == "assistant" || content.Role == "model" {
			role = anthropic.MessageParamRoleAssistant
		}

		// Create content blocks
		var contentBlocks []anthropic.ContentBlockParamUnion

		if len(content.Parts) > 0 {
			// Handle parts if they exist
			for _, part := range content.Parts {
				// Type assertion to determine part type
				switch v := part.(type) {
				case string:
					// Plain text
					contentBlocks = append(contentBlocks, anthropic.NewTextBlock(v))
				case map[string]interface{}:
					// Attempt to handle function calls or other structured content
					if funcCall, ok := v["function_call"].(map[string]interface{}); ok {
						name, _ := funcCall["name"].(string)
						args, _ := funcCall["args"].(map[string]interface{})
						id := fmt.Sprintf("tool_%d", len(contentBlocks))

						contentBlocks = append(contentBlocks, anthropic.ContentBlockParamOfToolUse(id, args, name))
					}
				}
			}
		} else if content.Message != "" {
			// Use message directly if parts are empty
			contentBlocks = append(contentBlocks, anthropic.NewTextBlock(content.Message))
		}

		messages = append(messages, anthropic.MessageParam{
			Role:    role,
			Content: contentBlocks,
		})
	}

	return messages
}

// prepareFunctionTools converts tool declarations to Anthropic tool parameters
func prepareFunctionTools(config *models.GenerateContentConfig) []anthropic.ToolUnionParam {
	if config == nil || len(config.Tools) == 0 {
		return nil
	}

	var tools []anthropic.ToolUnionParam

	for _, toolDecl := range config.Tools {
		for i, _ := range toolDecl.FunctionDeclarations {
			// Basic parsing of function declaration - in real implementation would need more robust parsing
			toolParam := anthropic.ToolParam{
				Name:        fmt.Sprintf("function_%d", i),
				Description: anthropic.String("Function tool"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]map[string]interface{}{
						"input": {
							"type": "string",
						},
					},
				},
			}

			// Convert to ToolUnionParam
			tools = append(tools, anthropic.ToolUnionParam{
				OfTool: &toolParam,
			})
		}
	}

	return tools
}

// anthropicResponseToLLMResponse converts Anthropic's response to models.LLMResponse
func anthropicResponseToLLMResponse(anthResponse *anthropic.Message) *models.LLMResponse {
	// Create a content object from the response
	content := &models.Content{
		Role: "assistant",
	}

	// Process content blocks
	if len(anthResponse.Content) > 0 {
		var sb strings.Builder

		for _, block := range anthResponse.Content {
			switch block := block.AsAny().(type) {
			case anthropic.TextBlock:
				sb.WriteString(block.Text)
			case anthropic.ToolUseBlock:
				// Tool use blocks would need specialized handling
				// This is simplified
				sb.WriteString(fmt.Sprintf("[Tool Use: %s]", block.Name))
			}
		}

		content.Message = sb.String()
	}

	// Create the final response
	response := &models.LLMResponse{
		Content: content,
		Usage: models.UsageMetrics{
			PromptTokens:     int(anthResponse.Usage.InputTokens),
			CompletionTokens: int(anthResponse.Usage.OutputTokens),
			TotalTokens:      int(anthResponse.Usage.InputTokens + anthResponse.Usage.OutputTokens),
			LatencyMs:        float64(0), // Not provided directly by Anthropic
			CostCents:        0.0,        // Would need pricing calculation
		},
	}

	// Set error information if there's a stop reason that indicates an issue
	if anthResponse.StopReason == "max_tokens" {
		maxTokensErr := "MAX_TOKENS"
		response.ErrorCode = &maxTokensErr
		errMsg := "Response was cut off due to token limit"
		response.ErrorMessage = &errMsg
	}

	return response
}

// Call implements the LLM interface Call method.
func (c *AnthropicClient) Call(ctx context.Context, request *models.LLMRequest) (*models.LLMResponse, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Prepare messages
	messages := contentToMessageParams(request.Contents)

	// Create system instruction text blocks if provided
	var systemTextBlocks []anthropic.TextBlockParam
	if request.Config != nil && request.Config.SystemInstruction != "" {
		systemTextBlocks = []anthropic.TextBlockParam{
			{
				Text: request.Config.SystemInstruction,
				Type: "text",
			},
		}
	}

	// Prepare max tokens
	maxTokens := int64(defaultMaxTokens)
	if request.Config != nil && request.Config.MaxTokens > 0 {
		maxTokens = int64(request.Config.MaxTokens)
	}

	// Create base message params
	msgParams := anthropic.MessageNewParams{
		Model:     mapToAnthropicModel(c.modelName),
		Messages:  messages,
		System:    systemTextBlocks,
		MaxTokens: maxTokens,
	}

	// Set request timeout and other options
	var callOpts []option.RequestOption
	if c.config.Timeout > 0 {
		callOpts = append(callOpts, option.WithRequestTimeout(time.Duration(c.config.Timeout)*time.Second))
	}

	// Add optional parameters
	if request.Config != nil {
		// Add temperature if provided
		if request.Config.Temperature > 0 {
			callOpts = append(callOpts, option.WithJSONSet("temperature", request.Config.Temperature))
		}

		// Add top_p if provided
		if request.Config.TopP > 0 {
			callOpts = append(callOpts, option.WithJSONSet("top_p", request.Config.TopP))
		}

		// Prepare tools if applicable
		if len(request.Config.Tools) > 0 {
			toolsParam := prepareFunctionTools(request.Config)
			if len(toolsParam) > 0 {
				msgParams.Tools = toolsParam
				// Enable auto tool choice
				msgParams.ToolChoice = anthropic.ToolChoiceUnionParam{
					OfAuto: &anthropic.ToolChoiceAutoParam{
						Type: "auto",
					},
				}
			}
		}
	}

	// Make the API call
	response, err := c.client.Messages.New(ctx, msgParams, callOpts...)
	if err != nil {
		return nil, fmt.Errorf("Anthropic API call failed: %w", err)
	}

	// Convert to LLMResponse
	return anthropicResponseToLLMResponse(response), nil
}

// BatchCall implements the LLM interface BatchCall method.
func (c *AnthropicClient) BatchCall(ctx context.Context, requests []*models.LLMRequest) ([]*models.LLMResponse, error) {
	responses := make([]*models.LLMResponse, len(requests))
	var err error

	// Process each request sequentially
	// Note: The Anthropic API doesn't have a native batch endpoint, so we process sequentially
	for i, req := range requests {
		responses[i], err = c.Call(ctx, req)
		if err != nil {
			return responses, fmt.Errorf("error processing request %d: %w", i, err)
		}
	}

	return responses, nil
}

// SupportedModels returns a list of model names supported by this client.
func (c *AnthropicClient) SupportedModels() []string {
	return []string{
		"claude-3-opus",
		"claude-3-sonnet",
		"claude-3-haiku",
		"claude-3.5-sonnet",
	}
}
