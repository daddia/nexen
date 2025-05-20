package anthropic

import (
	"context"
	"testing"

	"github.com/nexen/models"
	"github.com/nexen/services/connectors/common"
)

func TestAnthropicClientCreation(t *testing.T) {
	// Test client creation with missing API key
	_, err := NewAnthropicClient("claude-3-sonnet")
	if err == nil {
		t.Fatal("Expected error for missing API key, got nil")
	}

	// Test client creation with API key
	client, err := NewAnthropicClient("claude-3-sonnet", common.WithAPIKey("test-api-key"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("Client is nil")
	}

	// Check model name mapping
	anthClient, ok := client.(*AnthropicClient)
	if !ok {
		t.Fatal("Client is not an AnthropicClient")
	}
	if anthClient.modelName != "claude-3-sonnet" {
		t.Fatalf("Expected model name 'claude-3-sonnet', got '%s'", anthClient.modelName)
	}
}

func TestMapToAnthropicModel(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"claude-3-opus", "claude-3-opus-latest"},
		{"claude-3-sonnet", "claude-3-sonnet-20240229"},
		{"claude-3-haiku", "claude-3-haiku-20240307"},
		{"claude-3.5-sonnet", "claude-3-5-sonnet-latest"},
		{"unknown-model", "unknown-model"}, // Should return the original name
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := mapToAnthropicModel(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestContentToMessageParams(t *testing.T) {
	testContents := []models.Content{
		{
			Role:    "user",
			Message: "Hello, world!",
		},
		{
			Role:    "assistant",
			Message: "Hi there!",
		},
	}

	messages := contentToMessageParams(testContents)

	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	if messages[0].Role != "user" {
		t.Errorf("Expected 'user' role, got '%s'", messages[0].Role)
	}

	if messages[1].Role != "assistant" {
		t.Errorf("Expected 'assistant' role, got '%s'", messages[1].Role)
	}
}

func TestMockCall(t *testing.T) {
	// Create a client with a mock API key
	client, err := NewAnthropicClient("claude-3-sonnet", common.WithAPIKey("test-api-key"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// This is just a mock test since we can't actually call the API without real credentials
	// In a real project, you would use mocking to simulate the API response
	request := &models.LLMRequest{
		Model: "claude-3-sonnet",
		Contents: []models.Content{
			{
				Role:    "user",
				Message: "Hello, world!",
			},
		},
	}

	// We expect this to actually fail due to the invalid API key
	// The test just verifies the request handling logic
	_, err = client.Call(context.Background(), request)
	if err == nil {
		t.Fatal("Expected error for invalid API key, got nil")
	}
}
