package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// DatabaseConfig holds the shared database connection settings used by every
// process of the backend. The URL carries credentials and is never logged.
type DatabaseConfig struct {
	URL      string
	MaxConns int32
}

func newDatabaseConfig(lookup func(string) (string, bool)) (DatabaseConfig, error) {
	cfg := DatabaseConfig{MaxConns: 4}

	raw, ok := lookup("DATABASE_URL")
	if !ok || strings.TrimSpace(raw) == "" {
		return DatabaseConfig{}, fmt.Errorf("DATABASE_URL must contain a postgresql connection string, such as postgresql://user:pass@host:5432/dbname")
	}
	cfg.URL = strings.TrimSpace(raw)

	parsed, err := url.Parse(cfg.URL)
	if err != nil || !parsed.IsAbs() {
		return DatabaseConfig{}, fmt.Errorf("DATABASE_URL must be an absolute URL, such as postgresql://user:pass@host:5432/dbname")
	}
	switch parsed.Scheme {
	case "postgres", "postgresql":
	default:
		return DatabaseConfig{}, fmt.Errorf("DATABASE_URL must use the postgres or postgresql scheme")
	}
	if parsed.Host == "" {
		return DatabaseConfig{}, fmt.Errorf("DATABASE_URL must contain a host")
	}
	if parsed.User == nil || parsed.User.Username() == "" {
		return DatabaseConfig{}, fmt.Errorf("DATABASE_URL must contain a user")
	}
	if strings.Trim(parsed.Path, "/") == "" {
		return DatabaseConfig{}, fmt.Errorf("DATABASE_URL must contain a database name, such as /postgres")
	}

	// Supabase free tiers keep a small connection budget; never cannonball
	// the pool against a shared managed instance.
	if v, ok := lookup("DB_POOL_MAX_CONNS"); ok {
		v = strings.TrimSpace(v)
		if v != "" {
			conns, err := strconv.ParseInt(v, 10, 32)
			if err != nil || conns < 1 || conns > 100 {
				return DatabaseConfig{}, fmt.Errorf("DB_POOL_MAX_CONNS must be between 1 and 100")
			}
			cfg.MaxConns = int32(conns)
		}
	}
	return cfg, nil
}
