# Connectors Module

The connectors module implements integrations between our platform and external LLM providers (OpenAI, Anthropic, etc). It defines a standard `LLM` interface that all provider adapters must implement.

## Features

- Standard LLM interface for all providers
- Provider-specific adapters (OpenAI, Anthropic, Google, etc.)
- Model name regex matching for factory creation
- Shared connection utilities (retry, timeout, region routing)
- Configurable options via functional option pattern

## Usage

### Creating an LLM client

```go
import (
    "context"
    "github.com/nexen/models"
    "github.com/nexen/services/connectors"
    "github.com/nexen/services/connectors/common"
)

func main() {
    // Get an LLM client for a specific model
    llm, err := connectors.NewLLM("claude-3-sonnet", 
        common.WithAPIKey("your-api-key"),
        common.WithTimeout(60))
    if err != nil {
        // Handle error
    }

    // Create a request
    request := &models.LLMRequest{
        Model: "claude-3-sonnet",
        Contents: []models.Content{
            {
                Role:    "user",
                Message: "Hello, world!",
            },
        },
    }

    // Call the LLM
    ctx := context.Background()
    response, err := llm.Call(ctx, request)
    if err != nil {
        // Handle error
    }

    // Process response
    message := response.Content.Message
    // ...
}
```

### Batching Requests

```go
requests := []*models.LLMRequest{
    // multiple requests
}

responses, err := llm.BatchCall(ctx, requests)
if err != nil {
    // Handle error
}
```

## Provider Support

The connectors module currently supports the following LLM providers:

| Provider | Status | Supported Models |
|----------|--------|------------------|
| Anthropic | ✅ Complete | claude-3-opus, claude-3-sonnet, claude-3-haiku, claude-3.5-sonnet |
| OpenAI | ⚠️ WIP | gpt-4, gpt-4-turbo, gpt-3.5-turbo |
| Google | ⚠️ WIP | gemini-pro, gemini-ultra |
| Mistral | ⚠️ WIP | mistral-small, mistral-medium, mistral-large |
| Llama | ⚠️ WIP | llama-7b, llama-13b, llama-70b |
| Custom | ⚠️ WIP | custom endpoints |

## Getting Started with Development

1. **Navigate to module**

   ```bash
   cd services/connectors
   ```

2. **Fetch dependencies & tidy**

   ```bash
   go mod tidy
   ```

3. **Run unit tests**

   ```bash
   go test ./... -v
   ```

4. **Build and use the connector tool**

   ```bash
   go build -o ./bin/connector-tool ./cmd/connector-tool
   ```

## Adding a New Provider Adapter

1. Create a new directory for the provider (e.g., `services/connectors/newprovider/`)
2. Implement the `LLM` interface
3. Register model patterns in an `init()` function
4. Ensure your implementation handles:
   - Authentication
   - Request/response mapping
   - Error handling
   - Retries
