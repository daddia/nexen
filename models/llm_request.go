package models

import (
	"fmt"
)

// BaseTool defines the interface for tools that can be attached to an LLMRequest.
// Each tool must provide a Name() and a Declaration() string.
// Implementations should live in the connectors or tools package.
type BaseTool interface {
	Name() string
	Declaration() (string, error)
}

// Content represents a single piece of content to send to the model.
// Contains role (system/user/assistant) and text.
type Content struct {
	Role    string `json:"role"`
	Message string `json:"message"`
	// Parts can contain multiple content segments (text, images, etc.)
	Parts []any `json:"parts,omitempty"`
}

// GenerateContentConfig holds additional generation parameters, tools, and schema.
type GenerateContentConfig struct {
	SystemInstruction string            `json:"systemInstruction,omitempty"`
	Tools             []ToolDeclaration `json:"tools,omitempty"`
	ResponseSchema    any               `json:"responseSchema,omitempty"`
	ResponseMimeType  string            `json:"responseMimeType,omitempty"`
	Temperature       float64           `json:"temperature,omitempty"`
	TopP              float64           `json:"topP,omitempty"`
	MaxTokens         int               `json:"maxTokens,omitempty"`
	StopSequences     []string          `json:"stopSequences,omitempty"`
}

// LiveConnectConfig holds live connection settings for streaming or other integrations.
type LiveConnectConfig struct {
	EnableStreaming bool           `json:"enableStreaming,omitempty"`
	StreamTimeout   int            `json:"streamTimeout,omitempty"`
	CallbackURI     string         `json:"callbackUri,omitempty"`
	CustomConfig    map[string]any `json:"customConfig,omitempty"`
}

// ToolDeclaration represents a tool's function declaration for the model.
type ToolDeclaration struct {
	FunctionDeclarations []string `json:"functionDeclarations,omitempty"`
}

// LLMRequest defines the structure for a single call to an LLM service.
// It includes the prompt contents, generation config, and attached tools.
type LLMRequest struct {
	// Model is the identifier of the LLM to use (e.g. "gpt-4-turbo").
	Model string `json:"model"`

	// Contents holds the ordered messages (system, user, assistant) to send.
	Contents []Content `json:"contents"`

	// Config holds model generation settings (system instructions, tools, schema).
	Config *GenerateContentConfig `json:"config,omitempty"`

	// LiveConnect holds optional live-streaming or other live connection settings.
	LiveConnect LiveConnectConfig `json:"liveConnect,omitempty"`

	// ToolsDict maps tool names to instances for post-processing.
	// It is populated when tools are declared on the request.
	ToolsDict map[string]BaseTool `json:"-"` // Not serialized
}

// AppendInstructions appends one or more instructions to the system instruction in Config.
func (r *LLMRequest) AppendInstructions(instructions ...string) {
	if r.Config == nil {
		r.Config = &GenerateContentConfig{}
	}
	joined := ""
	for i, instr := range instructions {
		if i > 0 {
			joined += "\n\n"
		}
		joined += instr
	}
	if r.Config.SystemInstruction != "" {
		r.Config.SystemInstruction += "\n\n" + joined
	} else {
		r.Config.SystemInstruction = joined
	}
}

// AppendTools attaches tool declarations to the request and records instances in ToolsDict.
func (r *LLMRequest) AppendTools(tools ...BaseTool) error {
	if len(tools) == 0 {
		return nil
	}
	if r.Config == nil {
		r.Config = &GenerateContentConfig{}
	}
	if r.ToolsDict == nil {
		r.ToolsDict = make(map[string]BaseTool)
	}
	var decls []string
	for _, tool := range tools {
		declaration, err := tool.Declaration()
		if err != nil {
			return fmt.Errorf("getting declaration from tool %s: %w", tool.Name(), err)
		}
		decls = append(decls, declaration)
		r.ToolsDict[tool.Name()] = tool
	}
	r.Config.Tools = append(r.Config.Tools, ToolDeclaration{FunctionDeclarations: decls})
	return nil
}

// SetOutputSchema configures the expected output schema and mime type for the response.
func (r *LLMRequest) SetOutputSchema(schema any) {
	if r.Config == nil {
		r.Config = &GenerateContentConfig{}
	}
	r.Config.ResponseSchema = schema
	r.Config.ResponseMimeType = "application/json"
}

// Validate ensures the request has all required fields and is properly formatted.
func (r *LLMRequest) Validate() error {
	if r.Model == "" {
		return fmt.Errorf("model ID is required")
	}
	if len(r.Contents) == 0 {
		return fmt.Errorf("request must contain at least one content message")
	}
	return nil
}
