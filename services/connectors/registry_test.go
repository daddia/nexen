package connectors

import (
	"context"
	"testing"

	"github.com/nexen/models"
	"github.com/nexen/services/connectors/common"
)

// mockLLM is a test implementation of the LLM interface.
type mockLLM struct{}

func (m *mockLLM) Call(ctx context.Context, request *models.LLMRequest) (*models.LLMResponse, error) {
	return &models.LLMResponse{}, nil
}

func (m *mockLLM) BatchCall(ctx context.Context, requests []*models.LLMRequest) ([]*models.LLMResponse, error) {
	responses := make([]*models.LLMResponse, len(requests))
	for i := range requests {
		responses[i] = &models.LLMResponse{}
	}
	return responses, nil
}

func (m *mockLLM) SupportedModels() []string {
	return []string{"test-model"}
}

// mockConstructor is a test constructor function.
func mockConstructor(model string, opts ...common.Option) (common.LLM, error) {
	return &mockLLM{}, nil
}

func TestRegistry(t *testing.T) {
	// Clear the registry before testing
	mu.Lock()
	registry = make(map[string]constructorFn)
	resolveCache = make(map[string]constructorFn)
	mu.Unlock()

	// Test Register
	err := Register("test-.*", mockConstructor)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Test Resolve - positive case
	ctor, err := Resolve("test-model")
	if err != nil {
		t.Fatalf("Resolve failed for valid model: %v", err)
	}
	if ctor == nil {
		t.Fatal("Resolve returned nil constructor")
	}

	// Test Resolve - cache hit
	ctor, err = Resolve("test-model")
	if err != nil {
		t.Fatalf("Resolve cache hit failed: %v", err)
	}

	// Test Resolve - negative case
	_, err = Resolve("unknown-model")
	if err == nil {
		t.Fatal("Resolve should have failed for unknown model")
	}

	// Test NewLLM
	llm, err := NewLLM("test-model")
	if err != nil {
		t.Fatalf("NewLLM failed: %v", err)
	}
	if llm == nil {
		t.Fatal("NewLLM returned nil LLM")
	}

	// Test NewLLM - negative case
	_, err = NewLLM("unknown-model")
	if err == nil {
		t.Fatal("NewLLM should have failed for unknown model")
	}
}
