package models

import (
	"testing"
)

func TestLLMRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request LLMRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: LLMRequest{
				Model: "gpt-4",
				Contents: []Content{
					{Role: "user", Message: "Hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing model",
			request: LLMRequest{
				Contents: []Content{
					{Role: "user", Message: "Hello"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty contents",
			request: LLMRequest{
				Model:    "gpt-4",
				Contents: []Content{},
			},
			wantErr: true,
		},
		{
			name: "valid complex request",
			request: LLMRequest{
				Model: "gpt-4",
				Contents: []Content{
					{Role: "system", Message: "You are a helpful assistant"},
					{Role: "user", Message: "Hello"},
				},
				Config: &GenerateContentConfig{
					Temperature: 0.7,
					MaxTokens:   100,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppendInstructions(t *testing.T) {
	request := &LLMRequest{
		Model: "gpt-4",
		Contents: []Content{
			{Role: "user", Message: "Hello"},
		},
	}

	// Test initial append
	request.AppendInstructions("Be helpful")
	if request.Config == nil || request.Config.SystemInstruction != "Be helpful" {
		t.Errorf("Expected system instruction 'Be helpful', got %q",
			request.Config.SystemInstruction)
	}

	// Test appending additional instructions
	request.AppendInstructions("Be concise")
	expected := "Be helpful\n\nBe concise"
	if request.Config.SystemInstruction != expected {
		t.Errorf("Expected system instruction %q, got %q",
			expected, request.Config.SystemInstruction)
	}

	// Test appending multiple instructions
	request.AppendInstructions("Avoid jargon", "Use simple examples")
	expected = "Be helpful\n\nBe concise\n\nAvoid jargon\n\nUse simple examples"
	if request.Config.SystemInstruction != expected {
		t.Errorf("Expected system instruction %q, got %q",
			expected, request.Config.SystemInstruction)
	}
}

type mockTool struct {
	name string
	decl string
	err  error
}

func (m mockTool) Name() string {
	return m.name
}

func (m mockTool) Declaration() (string, error) {
	return m.decl, m.err
}

func TestAppendTools(t *testing.T) {
	request := &LLMRequest{
		Model: "gpt-4",
		Contents: []Content{
			{Role: "user", Message: "Hello"},
		},
	}

	tool1 := mockTool{name: "tool1", decl: `{"name":"tool1","description":"Tool 1"}`, err: nil}
	tool2 := mockTool{name: "tool2", decl: `{"name":"tool2","description":"Tool 2"}`, err: nil}

	// Test adding tools
	err := request.AppendTools(tool1, tool2)
	if err != nil {
		t.Errorf("AppendTools() error = %v", err)
	}

	// Check tools were added to the ToolsDict
	if len(request.ToolsDict) != 2 {
		t.Errorf("Expected 2 tools in ToolsDict, got %d", len(request.ToolsDict))
	}

	// Check tool declarations were added to Config
	if request.Config == nil || len(request.Config.Tools) != 1 {
		t.Errorf("Expected 1 ToolDeclaration in Config, got %d",
			len(request.Config.Tools))
	}

	if len(request.Config.Tools[0].FunctionDeclarations) != 2 {
		t.Errorf("Expected 2 function declarations, got %d",
			len(request.Config.Tools[0].FunctionDeclarations))
	}
}

func TestSetOutputSchema(t *testing.T) {
	request := &LLMRequest{
		Model: "gpt-4",
		Contents: []Content{
			{Role: "user", Message: "Hello"},
		},
	}

	type testSchema struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	schema := testSchema{Field1: "test", Field2: 123}
	request.SetOutputSchema(schema)

	if request.Config == nil {
		t.Fatal("Config is nil after SetOutputSchema")
	}

	if request.Config.ResponseMimeType != "application/json" {
		t.Errorf("Expected ResponseMimeType to be 'application/json', got %q",
			request.Config.ResponseMimeType)
	}

	if request.Config.ResponseSchema == nil {
		t.Error("ResponseSchema is nil")
	}
}
