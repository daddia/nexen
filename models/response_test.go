package models

import (
	"testing"
)

func TestCreateLLMResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *GenerateContentResponse
		wantErr  bool
		checkMsg string
	}{
		{
			name: "successful response",
			input: &GenerateContentResponse{
				Candidates: []Candidate{
					{
						Content: &Content{
							Role:    "assistant",
							Message: "Hello, I'm an AI assistant.",
						},
						FinishReason: "STOP",
					},
				},
				Usage: UsageMetrics{
					PromptTokens:     10,
					CompletionTokens: 20,
					LatencyMs:        150.5,
					CostCents:        0.002,
				},
			},
			wantErr:  false,
			checkMsg: "Hello, I'm an AI assistant.",
		},
		{
			name: "content with parts",
			input: &GenerateContentResponse{
				Candidates: []Candidate{
					{
						Content: &Content{
							Role:  "assistant",
							Parts: []any{"Hello, I'm an AI assistant."},
						},
						FinishReason: "STOP",
					},
				},
				Usage: UsageMetrics{
					PromptTokens:     10,
					CompletionTokens: 20,
				},
			},
			wantErr: false,
		},
		{
			name: "candidate with error",
			input: &GenerateContentResponse{
				Candidates: []Candidate{
					{
						Content:       nil,
						FinishReason:  "MAX_TOKENS",
						FinishMessage: "Response exceeded maximum token limit",
					},
				},
				Usage: UsageMetrics{
					PromptTokens:     10,
					CompletionTokens: 0,
				},
			},
			wantErr:  true,
			checkMsg: "MAX_TOKENS",
		},
		{
			name: "prompt feedback error",
			input: &GenerateContentResponse{
				Candidates: []Candidate{},
				PromptFeedback: &PromptFeedback{
					BlockReason:        "SAFETY",
					BlockReasonMessage: "Prompt contains unsafe content",
				},
				Usage: UsageMetrics{
					PromptTokens:     10,
					CompletionTokens: 0,
				},
			},
			wantErr:  true,
			checkMsg: "SAFETY",
		},
		{
			name: "empty response",
			input: &GenerateContentResponse{
				Candidates: []Candidate{},
				Usage: UsageMetrics{
					PromptTokens:     10,
					CompletionTokens: 0,
				},
			},
			wantErr:  true,
			checkMsg: "UNKNOWN_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateLLMResponse(tt.input)

			// Verify total tokens calculation
			expectedTotal := tt.input.Usage.PromptTokens + tt.input.Usage.CompletionTokens
			if result.Usage.TotalTokens != expectedTotal {
				t.Errorf("Expected total tokens %d, got %d",
					expectedTotal, result.Usage.TotalTokens)
			}

			// Check error status
			isError := result.IsError()
			if isError != tt.wantErr {
				t.Errorf("IsError() = %v, want %v", isError, tt.wantErr)
			}

			// For success case, check content
			if !tt.wantErr && tt.checkMsg != "" {
				if result.Content == nil || result.Content.Message != tt.checkMsg {
					t.Errorf("Expected message %q, got %v",
						tt.checkMsg, result.Content)
				}
			}

			// For error case, check error code
			if tt.wantErr && tt.checkMsg != "" {
				if result.ErrorCode == nil || *result.ErrorCode != tt.checkMsg {
					code := "nil"
					if result.ErrorCode != nil {
						code = *result.ErrorCode
					}
					t.Errorf("Expected error code %q, got %q", tt.checkMsg, code)
				}
			}
		})
	}
}

func TestLLMResponseErrorInterface(t *testing.T) {
	tests := []struct {
		name     string
		response LLMResponse
		expected string
	}{
		{
			name: "error message present",
			response: LLMResponse{
				ErrorMessage: strPtr("Something went wrong"),
			},
			expected: "Something went wrong",
		},
		{
			name: "error code only",
			response: LLMResponse{
				ErrorCode: strPtr("ERROR_CODE"),
			},
			expected: "Error: ERROR_CODE",
		},
		{
			name:     "no error info",
			response: LLMResponse{},
			expected: "Unknown error",
		},
		{
			name: "both error code and message",
			response: LLMResponse{
				ErrorCode:    strPtr("ERROR_CODE"),
				ErrorMessage: strPtr("Detailed error message"),
			},
			expected: "Detailed error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Error() method
			if tt.response.Error() != tt.expected {
				t.Errorf("Error() = %q, want %q", tt.response.Error(), tt.expected)
			}

			// Test IsError() method
			shouldBeError := tt.response.ErrorCode != nil || tt.response.ErrorMessage != nil
			if tt.response.IsError() != shouldBeError {
				t.Errorf("IsError() = %v, want %v", tt.response.IsError(), shouldBeError)
			}
		})
	}
}

// Helper to create string pointers
func strPtr(s string) *string {
	return &s
}
