// config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew_LoadsFromFileAndEnv(t *testing.T) {
	// create a temp directory and write a custom nexen.json
	tmp := t.TempDir()
	cfgFile := filepath.Join(tmp, "nexen.json")
	content := `{
		"server": {
			"port": 9090,
			"host": "127.0.0.1",
			"read_timeout": 7,
			"write_timeout": 14
		},
		"logging": {
			"level": "debug",
			"pretty": true,
			"prefix": "[TEST]",
			"syslog": true,
			"stdout": true
		},
		"redis": {
			"address": "redis.test:6379",
			"db": 2,
			"timeout": 8,
			"password": "redispass"
		},
		"telemetry": {
			"enabled": true,
			"collector_addr": "otel.test:4317"
		},
		"model_selection": {
			"strategy": "cost",
			"max_cost_per_request": 0.01,
			"max_latency_ms": 2000
		},
		"gateway": {
			"enable_grpc": true,
			"enable_rest": true,
			"cache_ttl": "7200s",
			"request_timeout": "15s"
		},
		"environment": "testing"
	}`
	if err := os.WriteFile(cfgFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// chdir into temp so Viper finds the file
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(tmp)

	// override one via ENV
	os.Setenv("NEXEN_SERVER_PORT", "9191")
	defer os.Unsetenv("NEXEN_SERVER_PORT")

	cfg, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	// ENV override should win
	if cfg.Server.Port != 9191 {
		t.Errorf("expected port=9191, got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 7*time.Second {
		t.Errorf("expected read_timeout=7s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Logging.Level != "debug" || !cfg.Logging.Pretty {
		t.Errorf("unexpected logging cfg: %+v", cfg.Logging)
	}
	if cfg.Redis.Address != "redis.test:6379" || cfg.Redis.DB != 2 {
		t.Errorf("unexpected redis cfg: %+v", cfg.Redis)
	}

	// Test new fields
	if cfg.Redis.Password != "redispass" {
		t.Errorf("expected redis password=redispass, got %s", cfg.Redis.Password)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host=127.0.0.1, got %s", cfg.Server.Host)
	}

	if !cfg.Telemetry.Enabled {
		t.Errorf("expected telemetry to be enabled")
	}

	if cfg.Telemetry.CollectorAddr != "otel.test:4317" {
		t.Errorf("expected collector_addr=otel.test:4317, got %s", cfg.Telemetry.CollectorAddr)
	}

	if cfg.ModelSelection.Strategy != "cost" {
		t.Errorf("expected strategy=cost, got %s", cfg.ModelSelection.Strategy)
	}

	if cfg.Gateway.CacheTTL != 7200*time.Second {
		t.Errorf("expected cache_ttl=7200s, got %v", cfg.Gateway.CacheTTL)
	}

	if cfg.Environment != "testing" {
		t.Errorf("expected environment=testing, got %s", cfg.Environment)
	}
}

func TestLoadServiceConfig(t *testing.T) {
	// create a temp directory and write a custom nexen.json
	tmp := t.TempDir()
	cfgFile := filepath.Join(tmp, "nexen.json")
	content := `{
		"server": {
			"port": 9090
		},
		"telemetry": {
			"enabled": true
		}
	}`
	if err := os.WriteFile(cfgFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// chdir into temp so Viper finds the file
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(tmp)

	serviceName := "model-selection"
	cfg, err := LoadServiceConfig(serviceName)
	if err != nil {
		t.Fatalf("LoadServiceConfig() error: %v", err)
	}

	if cfg.ServiceName != serviceName {
		t.Errorf("expected ServiceName=%s, got %s", serviceName, cfg.ServiceName)
	}

	if cfg.Telemetry.ServiceName != serviceName {
		t.Errorf("expected Telemetry.ServiceName=%s, got %s", serviceName, cfg.Telemetry.ServiceName)
	}
}
