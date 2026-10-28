package config

import (
	"strings"
	"testing"
)

func TestDatabaseConfigAcceptsPoolerURL(t *testing.T) {
	cfg, err := fromValues(map[string]string{
		"DATABASE_URL":      "postgresql://postgres:secret@aws-0-eu-west-1.pooler.supabase.com:5432/postgres",
		"DB_POOL_MAX_CONNS": "20",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Database.URL != "postgresql://postgres:secret@aws-0-eu-west-1.pooler.supabase.com:5432/postgres" {
		t.Fatalf("unexpected URL: %s", cfg.Database.URL)
	}
	if cfg.Database.MaxConns != 20 {
		t.Fatalf("MaxConns = %d, want 20", cfg.Database.MaxConns)
	}
}

func TestDatabaseConfigDefaultsPoolSize(t *testing.T) {
	cfg, err := fromValues(nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Database.MaxConns != 4 {
		t.Fatalf("MaxConns = %d, want 4", cfg.Database.MaxConns)
	}
}

func TestLoadRejectsInvalidDatabaseConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name, url, maxConns, errContent string
	}{
		{"missing", "", "", "DATABASE_URL"},
		{"not a url", "plain-hostname", "", "absolute"},
		{"wrong scheme", "mysql://user:pass@localhost/db", "", "postgres"},
		{"no host", "postgres://user:pass@/db", "", "host"},
		{"no user", "postgres://:pass@localhost:5432/db", "", "user"},
		{"no database", "postgres://user:pass@localhost:5432/", "", "database"},
		{"zero conns", "postgres://user:pass@localhost/db", "0", "DB_POOL_MAX_CONNS"},
		{"over max conns", "postgres://user:pass@localhost/db", "101", "DB_POOL_MAX_CONNS"},
		{"garbage conns", "postgres://user:pass@localhost/db", "many", "DB_POOL_MAX_CONNS"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fromValues(map[string]string{
				"DATABASE_URL":      tc.url,
				"DB_POOL_MAX_CONNS": tc.maxConns,
			})
			if err == nil {
				t.Fatal("expected configuration to fail")
			}
			if !strings.Contains(err.Error(), tc.errContent) {
				t.Fatalf("error %q must mention %q", err, tc.errContent)
			}
		})
	}
}
