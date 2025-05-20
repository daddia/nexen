# Nexen Configuration Module

This module provides configuration management for all Nexen services.

## Features

- Unified configuration handling for all Nexen services
- Configuration from files (JSON) with environment variable overrides
- Sensible defaults for all configuration options
- Support for service-specific configuration

## Usage

### Loading Configuration

```go
import "github.com/nexen/config"

// For a specific service
cfg, err := config.LoadServiceConfig("gateway")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

// Or general configuration
cfg, err := config.New()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

### Accessing Configuration

```go
// Server settings
port := cfg.Server.Port
host := cfg.Server.Host

// Redis connection
redisAddr := cfg.Redis.Address
redisPwd := cfg.Redis.Password

// Service-specific settings
if cfg.ServiceName == "gateway" {
    cacheTTL := cfg.Gateway.CacheTTL
    requestTimeout := cfg.Gateway.RequestTimeout
}

if cfg.ServiceName == "model-selection" {
    strategy := cfg.ModelSelection.Strategy
    maxCost := cfg.ModelSelection.MaxCostPerRequest
}
```

## Configuration File

The configuration is loaded from a `nexen.json` file in the project root:

```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0",
    "read_timeout": 5,
    "write_timeout": 10
  },
  "logging": {
    "level": "info",
    "pretty": false,
    "prefix": "[NEXEN]"
  },
  "redis": {
    "address": "localhost:6379",
    "db": 0,
    "timeout": 5,
    "password": ""
  },
  "telemetry": {
    "enabled": true,
    "collector_addr": "localhost:4317"
  },
  "model_selection": {
    "strategy": "balanced",
    "max_cost_per_request": 0.05
  },
  "gateway": {
    "enable_grpc": true,
    "enable_rest": true,
    "cache_ttl": "3600s",
    "request_timeout": "30s"
  }
}
```

## Environment Variables

All configuration can be overridden with environment variables using the prefix `NEXEN_` and uppercase keys with underscores:

```
NEXEN_SERVER_PORT=9000
NEXEN_LOGGING_LEVEL=debug
NEXEN_REDIS_ADDRESS=redis.internal:6379
NEXEN_TELEMETRY_ENABLED=true
NEXEN_MODEL_SELECTION_STRATEGY=cost
```

## Configuration Structure

The configuration structure includes:

- `Server`: HTTP server settings
- `Logging`: Logging configuration
- `Redis`: Redis connection settings
- `Telemetry`: OpenTelemetry configuration
- `ModelSelection`: Model selection service settings
- `Gateway`: API gateway settings
- `ServiceName`: Name of the current service
- `Environment`: Deployment environment (development, staging, production)

See `config.go` for the complete structure definition and default values.
