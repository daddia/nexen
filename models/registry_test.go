package models

import (
	"strings"
	"testing"
)

func setupTestRegistry() {
	ClearRegistry()

	// Add test models
	NewModelInfo(ModelInfo{
		ID:           "test-model-1",
		Profiles:     []string{ProfileChat, ProfileThinking},
		MaxTokens:    4096,
		CostPerToken: 0.0001,
		Provider:     ProviderOpenAI,
		CostTier:     CostTierStandard,
	}, "test-model-1")

	NewModelInfo(ModelInfo{
		ID:           "test-model-2",
		Profiles:     []string{ProfileChat, ProfileAgent, ProfileRAG},
		MaxTokens:    8192,
		CostPerToken: 0.0002,
		Provider:     ProviderAnthropic,
		CostTier:     CostTierPremium,
	}, "test-model-2")

	NewModelInfo(ModelInfo{
		ID:           "test-regex-model",
		Profiles:     []string{ProfileChat},
		MaxTokens:    16384,
		CostPerToken: 0.0003,
		Provider:     ProviderGoogle,
		CostTier:     CostTierPremium,
	}, "test-regex-.*")
}

func TestRegister(t *testing.T) {
	ClearRegistry()

	// Test valid registration
	err := Register("model-1", ModelInfo{ID: "model-1"})
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// Test invalid regex pattern
	err = Register("[invalid", ModelInfo{ID: "invalid"})
	if err == nil {
		t.Error("Register() should fail with invalid regex pattern")
	}

	// Test overwriting existing registration
	err = Register("model-1", ModelInfo{ID: "model-1-new"})
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// Verify the model was updated through Resolve
	info, err := Resolve("model-1")
	if err != nil {
		t.Errorf("Resolve() error = %v", err)
	}
	if info.ID != "model-1" {
		t.Errorf("Expected ID model-1, got %s", info.ID)
	}
}

func TestResolve(t *testing.T) {
	setupTestRegistry()

	tests := []struct {
		name        string
		model       string
		wantErr     bool
		expectedID  string
		expectedMax int
	}{
		{
			name:        "exact match",
			model:       "test-model-1",
			wantErr:     false,
			expectedID:  "test-model-1",
			expectedMax: 4096,
		},
		{
			name:        "regex match",
			model:       "test-regex-abc",
			wantErr:     false,
			expectedID:  "test-regex-abc", // Should preserve the requested ID
			expectedMax: 16384,
		},
		{
			name:    "no match",
			model:   "nonexistent-model",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := Resolve(tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if info.ID != tt.expectedID {
					t.Errorf("Resolve() ID = %v, want %v", info.ID, tt.expectedID)
				}
				if info.MaxTokens != tt.expectedMax {
					t.Errorf("Resolve() MaxTokens = %v, want %v",
						info.MaxTokens, tt.expectedMax)
				}
			}
		})
	}

	// Test caching
	info1, _ := Resolve("test-model-1")
	info2, _ := Resolve("test-model-1")
	// Should be identical but not the same pointer
	if &info1 == &info2 {
		t.Error("Resolve() should return copies, not the same instance")
	}
}

func TestListModels(t *testing.T) {
	setupTestRegistry()

	models := ListModels()
	if len(models) != 3 {
		t.Errorf("ListModels() returned %d models, want 3", len(models))
	}

	// Check if all models are included
	foundCount := 0
	for _, m := range models {
		if m == "test-model-1" || m == "test-model-2" || m == "test-regex-.*" {
			foundCount++
		}
	}
	if foundCount != 3 {
		t.Errorf("ListModels() missing expected models, found %d of 3", foundCount)
	}
}

func TestListModelsByProfile(t *testing.T) {
	setupTestRegistry()

	// Test chat profile (all models have this)
	chatModels := ListModelsByProfile(ProfileChat)
	if len(chatModels) != 3 {
		t.Errorf("ListModelsByProfile(chat) returned %d models, want 3",
			len(chatModels))
	}

	// Test RAG profile (only model-2 has this)
	ragModels := ListModelsByProfile(ProfileRAG)
	if len(ragModels) != 1 {
		t.Errorf("ListModelsByProfile(rag) returned %d models, want 1",
			len(ragModels))
	}
	if len(ragModels) > 0 && ragModels[0].ID != "test-model-2" {
		t.Errorf("Expected test-model-2 for RAG profile, got %s", ragModels[0].ID)
	}

	// Test nonexistent profile
	noneModels := ListModelsByProfile("nonexistent")
	if len(noneModels) != 0 {
		t.Errorf("ListModelsByProfile(nonexistent) returned %d models, want 0",
			len(noneModels))
	}
}

func TestListModelsByProvider(t *testing.T) {
	setupTestRegistry()

	// Test OpenAI provider
	openaiModels := ListModelsByProvider(ProviderOpenAI)
	if len(openaiModels) != 1 {
		t.Errorf("ListModelsByProvider(openai) returned %d models, want 1",
			len(openaiModels))
	}

	// Test with case insensitivity
	anthropicModels := ListModelsByProvider("AnThRoPiC")
	if len(anthropicModels) != 1 {
		t.Errorf("ListModelsByProvider(anthropic) returned %d models, want 1",
			len(anthropicModels))
	}

	// Test nonexistent provider
	noneModels := ListModelsByProvider("nonexistent")
	if len(noneModels) != 0 {
		t.Errorf("ListModelsByProvider(nonexistent) returned %d models, want 0",
			len(noneModels))
	}
}

func TestHasProfile(t *testing.T) {
	setupTestRegistry()

	tests := []struct {
		name     string
		model    string
		profile  string
		expected bool
		wantErr  bool
	}{
		{
			name:     "has profile",
			model:    "test-model-1",
			profile:  ProfileThinking,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "doesn't have profile",
			model:    "test-model-1",
			profile:  ProfileRAG,
			expected: false,
			wantErr:  false,
		},
		{
			name:    "nonexistent model",
			model:   "nonexistent",
			profile: ProfileChat,
			wantErr: true,
		},
		{
			name:     "regex model with profile",
			model:    "test-regex-something",
			profile:  ProfileChat,
			expected: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			has, err := HasProfile(tt.model, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("HasProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && has != tt.expected {
				t.Errorf("HasProfile() = %v, want %v", has, tt.expected)
			}
		})
	}
}

func TestInit(t *testing.T) {
	ClearRegistry()
	Init()

	// Check if common models are registered
	models := ListModels()
	if len(models) < 5 {
		t.Errorf("Init() should register at least 5 models, got %d", len(models))
	}

	// Check if specific models exist
	containsModel := func(pattern string) bool {
		for _, m := range models {
			if strings.Contains(m, pattern) {
				return true
			}
		}
		return false
	}

	expectedModels := []string{
		"gpt-4", "claude", "gemini", "mistral",
	}

	for _, model := range expectedModels {
		if !containsModel(model) {
			t.Errorf("Init() should register models containing %q", model)
		}
	}
}
