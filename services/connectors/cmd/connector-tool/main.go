package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nexen/models"
	"github.com/nexen/services/connectors"
	"github.com/nexen/services/connectors/common"

	// Import all connectors to register them
	_ "github.com/nexen/services/connectors/anthropic"
	_ "github.com/nexen/services/connectors/custom"
	_ "github.com/nexen/services/connectors/google"
	_ "github.com/nexen/services/connectors/llama"
	_ "github.com/nexen/services/connectors/mistral"
	_ "github.com/nexen/services/connectors/openai"
)

func main() {
	// Command-line flags
	modelFlag := flag.String("model", "gpt-4", "Model ID to test")
	promptFlag := flag.String("prompt", "Hello, world!", "Prompt to send")
	apiKeyFlag := flag.String("apikey", "", "API key (can also use env var)")
	timeoutFlag := flag.Int("timeout", 30, "Timeout in seconds")
	listFlag := flag.Bool("list", false, "List available registered model patterns")

	flag.Parse()

	// Handle list command
	if *listFlag {
		fmt.Println("Registered model patterns:")
		for _, pattern := range connectors.ListModelPatterns() {
			fmt.Printf("  - %s\n", pattern)
		}
		return
	}

	// Get API key
	apiKey := *apiKeyFlag
	if apiKey == "" {
		// Try to get from environment
		apiKey = os.Getenv("API_KEY")
	}

	// Create client with options
	llm, err := connectors.NewLLM(*modelFlag,
		common.WithAPIKey(apiKey),
		common.WithTimeout(*timeoutFlag),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}

	// Create request
	request := &models.LLMRequest{
		Model: *modelFlag,
		Contents: []models.Content{
			{
				Role:    "user",
				Message: *promptFlag,
			},
		},
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutFlag)*time.Second)
	defer cancel()

	// Call the LLM
	fmt.Printf("Sending request to %s...\n", *modelFlag)
	start := time.Now()

	response, err := llm.Call(ctx, request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling LLM: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)

	// Print response
	fmt.Println("\nResponse:")
	if response.Content != nil {
		fmt.Printf("%s\n", response.Content.Message)
	} else {
		fmt.Println("No content in response")
	}

	// Print metadata
	fmt.Printf("\nMetadata:\n")
	fmt.Printf("  Elapsed time: %v\n", elapsed)
	fmt.Printf("  Prompt tokens: %d\n", response.Usage.PromptTokens)
	fmt.Printf("  Completion tokens: %d\n", response.Usage.CompletionTokens)
	fmt.Printf("  Total tokens: %d\n", response.Usage.TotalTokens)
	fmt.Printf("  Cost (cents): %.5f\n", response.Usage.CostCents)

	// Print full JSON response for debugging
	fmt.Println("\nFull response (JSON):")
	jsonBytes, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(jsonBytes))
}
