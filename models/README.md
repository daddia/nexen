# Nexen Models Package

The `models` package defines shared data types and the canonical model registry used across all nexen services. It provides a type-safe way to define and interact with different LLM models and their capabilities.

## Overview

This package is purely declarative—it contains no business logic. Its main responsibilities are:

1. Define DTOs for LLM requests and responses
2. Declare model identifiers and profiles (thinking, agent, etc.)
3. Centralize metadata (cost tiers, max tokens, context window sizes)
4. Enable compile-time safety when referring to models across services

## Key Components

### Model Registry

The registry provides a central catalog of available models with their capabilities and constraints:

```go
// Get a model's information
modelInfo, err := models.Resolve("gpt-4-turbo")
if err != nil {
    // Handle unknown model
}

// Check if a model supports a specific capability
hasRag, err := models.HasProfile("claude-3-opus", models.ProfileRAG)
```

### LLM Request/Response

Standardized structures for making requests to models and handling their responses:

```go
// Create a new request
req := &models.LLMRequest{
    Model: "gpt-4-turbo",
    Contents: []models.Content{
        {Role: "user", Message: "Explain quantum computing"},
    },
}

// Add system instructions
req.AppendInstructions("Be concise and technical")

// Set output format
req.SetOutputSchema(MyResponseSchema{})
```

## Model Profiles

Models are tagged with capability profiles:

- `chat`: Basic conversational ability
- `thinking`: Deep reasoning capability
- `agent`: Tool/function calling capability
- `rag`: Retrieval-augmented generation support
- `creative`: Creative content generation
- `code`: Code generation capability

## Usage

To use the models package:

1. Import it in your service:
   ```go
   import "github.com/nexen/models"
   ```

2. Initialize the registry in your service startup:
   ```go
   func init() {
       models.Init() // Registers standard models
   }
   ```

3. Add custom models if needed:
   ```go
   models.NewModelInfo(models.ModelInfo{
       ID:           "custom-model",
       Profiles:     []string{models.ProfileChat},
       MaxTokens:    16000,
       CostPerToken: 0.00001,
       Provider:     models.ProviderCustom,
       CostTier:     models.CostTierStandard,
   }, "custom-model.*")
   ```

## Examples

### Creating a Request with Tools

```go
req := &models.LLMRequest{
    Model: "gpt-4-turbo",
    Contents: []models.Content{
        {Role: "user", Message: "Book a flight to Tokyo"},
    },
}

err := req.AppendTools(flightSearchTool, bookingTool)
if err != nil {
    // Handle tool declaration error
}
```

### Processing a Response

```go
resp := service.CallLLM(req)

if resp.IsError() {
    fmt.Println("Error:", resp.Error())
    return
}

fmt.Println("Response:", resp.Content.Message)
fmt.Printf("Tokens used: %d, Cost: $%.4f\n", 
    resp.Usage.TotalTokens, resp.Usage.CostCents/100)
``` 