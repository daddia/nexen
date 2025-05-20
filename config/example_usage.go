// config/example_usage.go
package config

import (
	"fmt"
	"log"
)

// ExampleUsage demonstrates how to load and use configuration in a service
func ExampleUsage() {
	// Load configuration for the gateway service
	cfg, err := LoadServiceConfig("gateway")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Access configuration values
	fmt.Printf("Service: %s\n", cfg.ServiceName)
	fmt.Printf("Environment: %s\n", cfg.Environment)
	fmt.Printf("Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("Logging Level: %s\n", cfg.Logging.Level)

	// Check if telemetry is enabled
	if cfg.Telemetry.Enabled {
		fmt.Printf("Telemetry enabled with collector at: %s\n", cfg.Telemetry.CollectorAddr)
	}

	// Use gateway-specific settings
	fmt.Printf("Gateway Cache TTL: %v\n", cfg.Gateway.CacheTTL)
	fmt.Printf("Gateway Request Timeout: %v\n", cfg.Gateway.RequestTimeout)

	// Example of how to handle Redis connection
	if cfg.Redis.Password != "" {
		fmt.Printf("Connecting to Redis at %s with password\n", cfg.Redis.Address)
	} else {
		fmt.Printf("Connecting to Redis at %s without password\n", cfg.Redis.Address)
	}
}

// ExampleModelSelectionUsage shows how to use configuration in the model-selection service
func ExampleModelSelectionUsage() {
	// Load configuration for the model-selection service
	cfg, err := LoadServiceConfig("model-selection")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Access model selection specific configuration
	fmt.Printf("Model Selection Strategy: %s\n", cfg.ModelSelection.Strategy)
	fmt.Printf("Max Cost Per Request: $%.4f\n", cfg.ModelSelection.MaxCostPerRequest)
	fmt.Printf("Max Latency: %dms\n", cfg.ModelSelection.MaxLatencyMs)

	// Access standard configuration
	fmt.Printf("Service: %s\n", cfg.ServiceName)
	fmt.Printf("Server Port: %d\n", cfg.Server.Port)
}
