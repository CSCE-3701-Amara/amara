package config

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func fromValues(values map[string]string) (Config, error) {
	merged := map[string]string{
		// Every process requires a database connection, so tests carry a
		// valid default and override it when validating the database rules.
		"DATABASE_URL": "postgres://user:pass@localhost:5432/db",
	}
	for key, value := range values {
		merged[key] = value
	}
	return load(func(key string) (string, bool) {
		value, ok := merged[key]
		return value, ok
	})
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := fromValues(nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environment != "development" || cfg.LogLevel != slog.LevelInfo || cfg.HTTP.Address != ":8080" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.HTTP.ReadHeaderTimeout <= 0 || cfg.HTTP.ReadTimeout <= 0 ||
		cfg.HTTP.WriteTimeout <= 0 || cfg.HTTP.IdleTimeout <= 0 || cfg.HTTP.ShutdownTimeout <= 0 {
		t.Fatal("all HTTP timeouts must have positive defaults")
	}
	if len(cfg.HTTP.AllowedOrigins) != 2 || len(cfg.HTTP.TrustedProxies) != 0 {
		t.Fatalf("unexpected origins or proxies: %+v", cfg.HTTP)
	}
}

func TestLoadOverrides(t *testing.T) {
	cfg, err := fromValues(map[string]string{
		"APP_ENV":                  "production",
		"LOG_LEVEL":                "debug",
		"HTTP_ADDR":                "127.0.0.1:9090",
		"HTTP_READ_HEADER_TIMEOUT": "2s",
		"HTTP_READ_TIMEOUT":        "10s",
		"HTTP_WRITE_TIMEOUT":       "20s",
		"HTTP_IDLE_TIMEOUT":        "1m",
		"HTTP_SHUTDOWN_TIMEOUT":    "8s",
		"CORS_ALLOWED_ORIGINS":     " https://amara.example, https://admin.amara.example ",
		"HTTP_TRUSTED_PROXIES":     " 127.0.0.1, 10.0.0.0/8, ::1 ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environment != "production" || cfg.LogLevel != slog.LevelDebug || cfg.HTTP.Address != "127.0.0.1:9090" {
		t.Fatalf("overrides were not loaded: %+v", cfg)
	}
	if cfg.HTTP.ReadHeaderTimeout != 2*time.Second || cfg.HTTP.ReadTimeout != 10*time.Second ||
		cfg.HTTP.WriteTimeout != 20*time.Second || cfg.HTTP.IdleTimeout != time.Minute || cfg.HTTP.ShutdownTimeout != 8*time.Second {
		t.Fatalf("timeout overrides were not loaded: %+v", cfg.HTTP)
	}
	if len(cfg.HTTP.AllowedOrigins) != 2 || cfg.HTTP.AllowedOrigins[0] != "https://amara.example" ||
		len(cfg.HTTP.TrustedProxies) != 3 || cfg.HTTP.TrustedProxies[0] != "127.0.0.1" {
		t.Fatalf("lists were not trimmed and parsed: %+v", cfg.HTTP)
	}
}

func TestLoadOnlyAllowsLocalOriginsByDefaultInDevelopment(t *testing.T) {
	for _, environment := range []string{"test", "staging", "production"} {
		t.Run(environment, func(t *testing.T) {
			cfg, err := fromValues(map[string]string{"APP_ENV": environment})
			if err != nil {
				t.Fatal(err)
			}
			if len(cfg.HTTP.AllowedOrigins) != 0 {
				t.Fatal("non-development environments must explicitly allow cross-origin requests")
			}
		})
	}
	cfg, err := fromValues(map[string]string{"CORS_ALLOWED_ORIGINS": ""})
	if err != nil || len(cfg.HTTP.AllowedOrigins) != 0 {
		t.Fatalf("an explicit empty origin list must disable cross-origin access: %+v, %v", cfg, err)
	}
}

func TestLoadRejectsInvalidConfiguration(t *testing.T) {
	for _, tc := range []struct{ name, value string }{
		{"APP_ENV", "unknown"},
		{"APP_ENV", ""},
		{"LOG_LEVEL", "verbose"},
		{"HTTP_ADDR", "8080"},
		{"HTTP_ADDR", ":65536"},
		{"HTTP_ADDR", ":-1"},
		{"HTTP_ADDR", ":http"},
		{"CORS_ALLOWED_ORIGINS", "*"},
		{"CORS_ALLOWED_ORIGINS", "https://*.amara.example"},
		{"CORS_ALLOWED_ORIGINS", "https://amara.example/path"},
		{"CORS_ALLOWED_ORIGINS", "https://amara.example?token=value"},
		{"CORS_ALLOWED_ORIGINS", "https://user:password@amara.example"},
		{"CORS_ALLOWED_ORIGINS", "localhost:5173"},
		{"HTTP_TRUSTED_PROXIES", "localhost"},
		{"HTTP_TRUSTED_PROXIES", "10.0.0.0/99"},
	} {
		t.Run(tc.name+"/"+tc.value, func(t *testing.T) {
			if _, err := fromValues(map[string]string{tc.name: tc.value}); err == nil || !strings.Contains(err.Error(), tc.name) {
				t.Fatalf("expected an error naming %s, got %v", tc.name, err)
			}
		})
	}
	for _, name := range []string{
		"HTTP_READ_HEADER_TIMEOUT", "HTTP_READ_TIMEOUT", "HTTP_WRITE_TIMEOUT", "HTTP_IDLE_TIMEOUT", "HTTP_SHUTDOWN_TIMEOUT",
	} {
		for _, value := range []string{"", "0s", "-1s", "invalid-sensitive-value"} {
			t.Run(name+"/"+value, func(t *testing.T) {
				_, err := fromValues(map[string]string{name: value})
				if err == nil || !strings.Contains(err.Error(), name) {
					t.Fatalf("expected an error naming %s, got %v", name, err)
				}
				if strings.Contains(err.Error(), "invalid-sensitive-value") {
					t.Fatal("configuration errors must not echo environment values")
				}
			})
		}
	}
}
