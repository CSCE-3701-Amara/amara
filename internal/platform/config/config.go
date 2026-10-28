package config

import (
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds process configuration. Secrets are supplied through environment
// variables by Infisical; the application does not read local .env files.
type Config struct {
	Environment string
	LogLevel    slog.Level
	HTTP        HTTPConfig
	Database    DatabaseConfig
	LogFormat   string
}

// HTTPConfig contains transport settings shared by the API server and router.
type HTTPConfig struct {
	Address           string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	AllowedOrigins    []string
	TrustedProxies    []string
}

// Load reads and validates configuration once at process startup.
func Load() (Config, error) {
	return load(os.LookupEnv)
}

func load(lookup func(string) (string, bool)) (Config, error) {
	value := func(name, fallback string) string {
		if v, ok := lookup(name); ok {
			return strings.TrimSpace(v)
		}
		return fallback
	}

	cfg := Config{
		Environment: value("APP_ENV", "development"),
		LogFormat:   value("LOG_FORMAT", "text"),
		HTTP: HTTPConfig{
			Address:           value("HTTP_ADDR", ":8080"),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
			ShutdownTimeout:   10 * time.Second,
			TrustedProxies:    splitList(value("HTTP_TRUSTED_PROXIES", "")),
		},
	}

	// TODO: Adjust environements configuration for different load. Assumed only development now.
	// Comeback at Sprint 2
	switch cfg.Environment {
	case "development", "test", "staging", "production":
	default:
		return Config{}, fmt.Errorf("APP_ENV must be development, test, staging, or production")
	}

	databaseCfg, err := newDatabaseConfig(lookup)
	if err != nil {
		return Config{}, err
	}
	cfg.Database = databaseCfg
	switch cfg.LogFormat {
	case "json", "text":
	default:
		return Config{}, fmt.Errorf("LOG_FORMAT must be json or text")
	}

	if err := cfg.LogLevel.UnmarshalText([]byte(value("LOG_LEVEL", "info"))); err != nil {
		return Config{}, fmt.Errorf("LOG_LEVEL must be a valid slog level, such as debug, info, warn, or error")
	}

	_, port, err := net.SplitHostPort(cfg.HTTP.Address)
	if err != nil {
		return Config{}, fmt.Errorf("HTTP_ADDR must be a host:port address")
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 0 || portNumber > 65535 {
		return Config{}, fmt.Errorf("HTTP_ADDR must contain a port between 0 and 65535")
	}

	for _, setting := range []struct {
		name   string
		target *time.Duration
	}{
		{"HTTP_READ_HEADER_TIMEOUT", &cfg.HTTP.ReadHeaderTimeout},
		{"HTTP_READ_TIMEOUT", &cfg.HTTP.ReadTimeout},
		{"HTTP_WRITE_TIMEOUT", &cfg.HTTP.WriteTimeout},
		{"HTTP_IDLE_TIMEOUT", &cfg.HTTP.IdleTimeout},
		{"HTTP_SHUTDOWN_TIMEOUT", &cfg.HTTP.ShutdownTimeout},
	} {
		duration, err := time.ParseDuration(value(setting.name, setting.target.String()))
		if err != nil || duration <= 0 {
			return Config{}, fmt.Errorf("%s must be a positive duration, such as 5s", setting.name)
		}
		*setting.target = duration
	}

	origins := ""
	if cfg.Environment == "development" {
		origins = "http://localhost:5173,http://127.0.0.1:5173"
	}
	cfg.HTTP.AllowedOrigins = splitList(value("CORS_ALLOWED_ORIGINS", origins))
	for _, origin := range cfg.HTTP.AllowedOrigins {
		u, err := url.Parse(origin)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") ||
			u.Hostname() == "" || u.User != nil || u.Path != "" ||
			u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || strings.Contains(origin, "*") {
			return Config{}, fmt.Errorf("CORS_ALLOWED_ORIGINS must contain exact HTTP(S) origins without paths or wildcards")
		}
	}

	for _, proxy := range cfg.HTTP.TrustedProxies {
		if net.ParseIP(proxy) != nil {
			continue
		}
		if _, _, err := net.ParseCIDR(proxy); err != nil {
			return Config{}, fmt.Errorf("HTTP_TRUSTED_PROXIES must contain IP addresses or CIDR ranges")
		}
	}

	return cfg, nil
}

func splitList(value string) []string {
	var values []string
	for item := range strings.SplitSeq(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			values = append(values, item)
		}
	}
	return values
}
