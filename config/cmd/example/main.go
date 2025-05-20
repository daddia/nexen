// Example application for the Nexen config module
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nexen/config"
)

func main() {
	if len(os.Args) > 1 {
		// If a service name is provided, load service-specific config
		serviceName := os.Args[1]
		cfg, err := config.LoadServiceConfig(serviceName)
		if err != nil {
			log.Fatalf("Failed to load configuration for service %s: %v", serviceName, err)
		}

		fmt.Printf("Configuration for service: %s\n", cfg.ServiceName)
		fmt.Printf("Environment: %s\n", cfg.Environment)

		// Display service-specific configuration
		switch serviceName {
		case "gateway":
			fmt.Printf("Gateway REST API enabled: %t\n", cfg.Gateway.EnableREST)
			fmt.Printf("Gateway gRPC API enabled: %t\n", cfg.Gateway.EnableGRPC)
			fmt.Printf("Gateway cache TTL: %v\n", cfg.Gateway.CacheTTL)
		case "model-selection":
			fmt.Printf("Model selection strategy: %s\n", cfg.ModelSelection.Strategy)
			fmt.Printf("Max cost per request: $%.4f\n", cfg.ModelSelection.MaxCostPerRequest)
			fmt.Printf("Max latency: %dms\n", cfg.ModelSelection.MaxLatencyMs)
		default:
			fmt.Println("No service-specific configuration to display")
		}
	} else {
		// Load general configuration
		cfg, err := config.New()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		fmt.Println("General Nexen Configuration:")
		fmt.Printf("Server port: %d\n", cfg.Server.Port)
		fmt.Printf("Server host: %s\n", cfg.Server.Host)
		fmt.Printf("Logging level: %s\n", cfg.Logging.Level)
		fmt.Printf("Redis address: %s\n", cfg.Redis.Address)

		if cfg.Telemetry.Enabled {
			fmt.Printf("Telemetry enabled with collector at: %s\n", cfg.Telemetry.CollectorAddr)
		} else {
			fmt.Println("Telemetry disabled")
		}
	}
}
