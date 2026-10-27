package database

import (
	"context"
	"os"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
)

// TestNewPoolConnectsToLiveDatabase runs only when a reachable database is
// provided; ordinary test runs and CI stay offline.
func TestNewPoolConnectsToLiveDatabase(t *testing.T) {
	url := os.Getenv("AMARA_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set AMARA_TEST_DATABASE_URL to run the live database connection test")
	}
	cfg := config.DatabaseConfig{URL: url, MaxConns: 1}
	pool, err := NewPool(t.Context(), cfg)
	if err != nil {
		t.Fatalf("connect and ping: %v", err)
	}
	defer pool.Close()
	ctx, cancel := context.WithTimeout(context.Background(), PingTimeout)
	defer cancel()
	var one int
	if err := pool.QueryRow(ctx, "select 1").Scan(&one); err != nil {
		t.Fatalf("query: %v", err)
	}
	if one != 1 {
		t.Fatalf("query returned %d", one)
	}
}

func TestNewPoolRejectsUnreachableDatabase(t *testing.T) {
	cfg := config.DatabaseConfig{
		URL:      "postgres://user:pass@127.0.0.1:1/db",
		MaxConns: 1,
	}
	pool, err := NewPool(t.Context(), cfg)
	if err == nil {
		pool.Close()
		t.Fatal("an unreachable database must fail the pool startup")
	}
}
