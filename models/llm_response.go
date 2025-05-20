package models

// UsageMetrics captures resource usage for an LLM call.
type UsageMetrics struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"promptTokens"`

	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completionTokens"`

	// TotalTokens is the sum of prompt and completion tokens.
	TotalTokens int `json:"totalTokens"`

	// LatencyMs is the request latency in milliseconds.
	LatencyMs float64 `json:"latencyMs"`

	// CostCents is the estimated cost in cents.
	CostCents float64 `json:"costCents"`
}

// GroundingMetadata contains references to sources used for grounding.
type GroundingMetadata struct {
	// Citations is a list of source citations for generated content.
	Citations []Citation `json:"citations,omitempty"`

	// GroundingScore indicates confidence in grounding (0-1).
	GroundingScore float64 `json:"groundingScore,omitempty"`
}

// Citation represents a reference to a source document.
type Citation struct {
	// SourceID identifies the referenced document.
	SourceID string `json:"sourceId"`

	// Title is the title of the source document.
	Title string `json:"title,omitempty"`

	// URL is the location of the source document.
	URL string `json:"url,omitempty"`

	// StartIndex is the start position in generated text.
	StartIndex int `json:"startIndex,omitempty"`

	// EndIndex is the end position in generated text.
	EndIndex int `json:"endIndex,omitempty"`
}

// GenerateContentResponse represents the vendor-specific response.
type GenerateContentResponse struct {
	// Candidates are the potential responses from the model.
	Candidates []Candidate `json:"candidates,omitempty"`

	// PromptFeedback contains validation or safety feedback.
	PromptFeedback *PromptFeedback `json:"promptFeedback,omitempty"`

	// Usage captures resource usage metrics.
	Usage UsageMetrics `json:"usage"`
}

// Candidate represents a single completion from the model.
type Candidate struct {
	// Content is the generated text and metadata.
	Content *Content `json:"content,omitempty"`

	// FinishReason indicates why generation stopped.
	FinishReason string `json:"finishReason,omitempty"`

	// FinishMessage provides details about the finish reason.
	FinishMessage string `json:"finishMessage,omitempty"`

	// GroundingMetadata contains citation data if enabled.
	GroundingMetadata *GroundingMetadata `json:"groundingMetadata,omitempty"`
}

// PromptFeedback contains information about prompt validation.
type PromptFeedback struct {
	// BlockReason indicates why the prompt was blocked (safety, etc.).
	BlockReason string `json:"blockReason,omitempty"`

	// BlockReasonMessage provides details about the block.
	BlockReasonMessage string `json:"blockReasonMessage,omitempty"`
}

// LLMResponse represents the standardized response from an LLM call.
// It captures either a successful content candidate or error details, along with metadata.
type LLMResponse struct {
	// Content is the primary output from the model, if available.
	Content *Content `json:"content,omitempty"`

	// GroundingMetadata holds any grounding or reference information.
	GroundingMetadata *GroundingMetadata `json:"groundingMetadata,omitempty"`

	// Partial indicates whether this is part of an unfinished stream.
	Partial *bool `json:"partial,omitempty"`

	// TurnComplete indicates whether the streaming turn is complete.
	TurnComplete *bool `json:"turnComplete,omitempty"`

	// ErrorCode for failures; provider-specific code or finish reason.
	ErrorCode *string `json:"errorCode,omitempty"`

	// ErrorMessage provides human-readable error details.
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// Interrupted signals if generation was interrupted (e.g. user cancel).
	Interrupted *bool `json:"interrupted,omitempty"`

	// CustomMetadata holds arbitrary, JSON-serializable metadata.
	CustomMetadata map[string]any `json:"customMetadata,omitempty"`

	// Usage captures tokens used, latency, and cost details.
	Usage UsageMetrics `json:"usage"`
}

// CreateLLMResponse constructs an LLMResponse from a provider-specific response.
// It extracts the first candidate if present, otherwise populates error fields.
func CreateLLMResponse(resp *GenerateContentResponse) LLMResponse {
	var result LLMResponse
	usage := UsageMetrics{
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.PromptTokens + resp.Usage.CompletionTokens,
		LatencyMs:        resp.Usage.LatencyMs,
		CostCents:        resp.Usage.CostCents,
	}
	result.Usage = usage

	if len(resp.Candidates) > 0 {
		cand := resp.Candidates[0]
		if cand.Content != nil && (len(cand.Content.Parts) > 0 || cand.Content.Message != "") {
			result.Content = cand.Content
			result.GroundingMetadata = cand.GroundingMetadata
			return result
		}
		// Candidate present but no content parts: treat as error
		result.ErrorCode = &cand.FinishReason
		result.ErrorMessage = &cand.FinishMessage
		return result
	}

	// No candidates: check prompt feedback
	if resp.PromptFeedback != nil {
		pf := resp.PromptFeedback
		result.ErrorCode = &pf.BlockReason
		result.ErrorMessage = &pf.BlockReasonMessage
		return result
	}

	// Fallback unknown error
	unknown := "UNKNOWN_ERROR"
	msg := "Unknown error occurred."
	result.ErrorCode = &unknown
	result.ErrorMessage = &msg
	return result
}

// IsError returns true if the response contains an error.
func (r *LLMResponse) IsError() bool {
	return r.ErrorCode != nil || r.ErrorMessage != nil
}

// Error implements the error interface.
func (r *LLMResponse) Error() string {
	if r.ErrorMessage != nil {
		return *r.ErrorMessage
	}
	if r.ErrorCode != nil {
		return "Error: " + *r.ErrorCode
	}
	return "Unknown error"
}
