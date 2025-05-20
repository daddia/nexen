package models

import (
	"os"
	"testing"
)

// TestMain is the entry point for all tests in the package
func TestMain(m *testing.M) {
	// Setup before running tests
	setup()

	// Run all tests
	code := m.Run()

	// Teardown after all tests have completed
	teardown()

	// Exit with the test status code
	os.Exit(code)
}

// setup initializes any test requirements
func setup() {
	// Initialize common test fixtures
	ClearRegistry()
}

// teardown cleans up after tests
func teardown() {
	// Clean up any test artifacts
	ClearRegistry()
}

// boolPtr helper function for creating boolean pointers
func boolPtr(b bool) *bool {
	return &b
}

// stringPtr helper function for creating string pointers (also in response_test.go)
func stringPtr(s string) *string {
	return &s
}

// Sample test tool for testing purposes
type TestTool struct {
	name string
	decl string
}

func (t TestTool) Name() string {
	return t.name
}

func (t TestTool) Declaration() (string, error) {
	return t.decl, nil
}

// createTestResponse creates a sample response for testing
func createTestResponse(isSuccess bool) *GenerateContentResponse {
	if isSuccess {
		return &GenerateContentResponse{
			Candidates: []Candidate{
				{
					Content: &Content{
						Role:    "assistant",
						Message: "This is a test response",
					},
				},
			},
			Usage: UsageMetrics{
				PromptTokens:     10,
				CompletionTokens: 5,
			},
		}
	} else {
		return &GenerateContentResponse{
			Candidates: []Candidate{
				{
					FinishReason:  "ERROR",
					FinishMessage: "This is a test error",
				},
			},
			Usage: UsageMetrics{
				PromptTokens: 10,
			},
		}
	}
}
