// config/config.go
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// RedisConfig holds connection settings for Redis.
type RedisConfig struct {
	Address  string        `mapstructure:"address"`
	DB       int           `mapstructure:"db"`
	Timeout  time.Duration `mapstructure:"timeout"` // in seconds
	Password string        `mapstructure:"password"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // seconds
	WriteTimeout time.Duration `mapstructure:"write_timeout"` // seconds
	Host         string        `mapstructure:"host"`
}

// LoggingConfig holds log‐level and format settings.
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Prefix string `mapstructure:"prefix"`
	Syslog bool   `mapstructure:"syslog"`
	Stdout bool   `mapstructure:"stdout"`
}

// TelemetryConfig holds OpenTelemetry settings
type TelemetryConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	CollectorAddr string `mapstructure:"collector_addr"`
	ServiceName   string `mapstructure:"service_name"`
}

// ModelSelectionConfig holds settings for model selection service
type ModelSelectionConfig struct {
	Strategy           string  `mapstructure:"strategy"` // e.g., "cost", "performance", "balanced"
	MaxCostPerRequest  float64 `mapstructure:"max_cost_per_request"`
	MaxLatencyMs       int     `mapstructure:"max_latency_ms"`
	ModelSelectionPort int     `mapstructure:"model_selection_port"`
}

// GatewayConfig holds settings specific to the API gateway
type GatewayConfig struct {
	EnableGRPC        bool          `mapstructure:"enable_grpc"`
	EnableREST        bool          `mapstructure:"enable_rest"`
	CacheTTL          time.Duration `mapstructure:"cache_ttl"`
	RequestTimeout    time.Duration `mapstructure:"request_timeout"`
	RateLimitRequests int           `mapstructure:"rate_limit_requests"`
	RateLimitPeriod   time.Duration `mapstructure:"rate_limit_period"`
}

// Config is your application's root configuration.
type Config struct {
	Server         ServerConfig         `mapstructure:"server"`
	Logging        LoggingConfig        `mapstructure:"logging"`
	Redis          RedisConfig          `mapstructure:"redis"`
	Telemetry      TelemetryConfig      `mapstructure:"telemetry"`
	ModelSelection ModelSelectionConfig `mapstructure:"model_selection"`
	Gateway        GatewayConfig        `mapstructure:"gateway"`
	ServiceName    string               `mapstructure:"service_name"`
	Environment    string               `mapstructure:"environment"`
}

// New reads configuration from nexen.json + ENV vars and returns a Config.
// ENV variables override file values, with prefix NEXEN_ (e.g. NEXEN_REDIS_ADDRESS).
func New() (*Config, error) {
	v := viper.New()
	v.SetConfigName("nexen")
	v.SetConfigType("json")
	v.AddConfigPath(".") // look in repo root
	v.SetEnvPrefix("nexen")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// sensible defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", 5)
	v.SetDefault("server.write_timeout", 10)

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.pretty", false)
	v.SetDefault("logging.prefix", "[NEXEN]")
	v.SetDefault("logging.syslog", false)
	v.SetDefault("logging.stdout", true)

	v.SetDefault("redis.address", "localhost:6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.timeout", 5)
	v.SetDefault("redis.password", "")

	v.SetDefault("telemetry.enabled", false)
	v.SetDefault("telemetry.collector_addr", "localhost:4317")

	v.SetDefault("gateway.enable_grpc", true)
	v.SetDefault("gateway.enable_rest", true)
	v.SetDefault("gateway.cache_ttl", "3600s")
	v.SetDefault("gateway.request_timeout", "30s")
	v.SetDefault("gateway.rate_limit_requests", 100)
	v.SetDefault("gateway.rate_limit_period", "1m")

	v.SetDefault("model_selection.strategy", "balanced")
	v.SetDefault("model_selection.max_cost_per_request", 0.05)
	v.SetDefault("model_selection.max_latency_ms", 5000)
	v.SetDefault("model_selection.model_selection_port", 8081)

	v.SetDefault("environment", "development")

	if err := v.ReadInConfig(); err != nil {
		// only error if config file missing *and* not a use-case for defaults/ENV
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	// Convert numeric time values to seconds for duration fields
	cfg.Server.ReadTimeout = time.Duration(v.GetInt("server.read_timeout")) * time.Second
	cfg.Server.WriteTimeout = time.Duration(v.GetInt("server.write_timeout")) * time.Second
	cfg.Redis.Timeout = time.Duration(v.GetInt("redis.timeout")) * time.Second

	// Handle string duration values
	if cacheTTL, err := time.ParseDuration(v.GetString("gateway.cache_ttl")); err == nil {
		cfg.Gateway.CacheTTL = cacheTTL
	}

	if requestTimeout, err := time.ParseDuration(v.GetString("gateway.request_timeout")); err == nil {
		cfg.Gateway.RequestTimeout = requestTimeout
	}

	if rateLimitPeriod, err := time.ParseDuration(v.GetString("gateway.rate_limit_period")); err == nil {
		cfg.Gateway.RateLimitPeriod = rateLimitPeriod
	}

	return &cfg, nil
}

// LoadServiceConfig loads configuration with specific service name
func LoadServiceConfig(serviceName string) (*Config, error) {
	cfg, err := New()
	if err != nil {
		return nil, err
	}

	cfg.ServiceName = serviceName
	cfg.Telemetry.ServiceName = serviceName

	return cfg, nil
}
