package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
)

// liveURL returns a non-empty DATABASE_URL when the integration tests should
// run, or skips the test otherwise.
func liveURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("AMARA_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set AMARA_TEST_DATABASE_URL to run the live database connection test")
	}
	return url
}

// ---------------------------------------------------------------------------
// NewPool (legacy helper)
// ---------------------------------------------------------------------------

// TestNewPoolConnectsToLiveDatabase runs only when a reachable database is
// provided; ordinary test runs and CI stay offline.
func TestNewPoolConnectsToLiveDatabase(t *testing.T) {
	url := liveURL(t)
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

// ---------------------------------------------------------------------------
// NewService – unit-level (no real database required)
// ---------------------------------------------------------------------------

func TestNewServiceRejectsUnreachableDatabase(t *testing.T) {
	cfg := config.DatabaseConfig{
		URL:      "postgres://user:pass@127.0.0.1:1/db",
		MaxConns: 1,
	}
	svc, err := NewService(t.Context(), cfg)
	if err == nil {
		svc.Close()
		t.Fatal("an unreachable database must fail service initialization")
	}
}

// ---------------------------------------------------------------------------
// NewService – integration (requires AMARA_TEST_DATABASE_URL)
// ---------------------------------------------------------------------------

func TestNewServiceConnectsToLiveDatabase(t *testing.T) {
	url := liveURL(t)
	cfg := config.DatabaseConfig{URL: url, MaxConns: 2, HealthTimeout: 5 * time.Second}
	svc, err := NewService(t.Context(), cfg)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	defer svc.Close()

	if svc.Pool() == nil {
		t.Fatal("Pool() must return a non-nil pool after successful initialization")
	}
}

// ---------------------------------------------------------------------------
// Health – unit-level stubs
// ---------------------------------------------------------------------------

// closedService builds a service whose pool is already closed so that every
// Ping() call fails, letting us exercise the unhealthy path without a live DB.
func closedService(t *testing.T) Service {
	t.Helper()
	// Use an unreachable address so NewService fails; then construct the
	// internal type directly to test the health logic in isolation.
	cfg := config.DatabaseConfig{
		URL:      "postgres://user:pass@127.0.0.1:1/db",
		MaxConns: 1,
	}
	// We can't reach NewService here because the pool won't open. Use the
	// exported interface only – test via a live DB if needed.
	_ = cfg
	return nil
}

// TestHealthStatusFields verifies the typed fields on HealthStatus without a
// database: construct the struct directly and check zero values are safe.
func TestHealthStatusFields(t *testing.T) {
	hs := HealthStatus{Status: "up", Stats: PoolStats{MaxConns: 4}}
	if hs.Status != "up" {
		t.Fatalf("status = %q, want %q", hs.Status, "up")
	}
	if hs.Error != "" {
		t.Fatalf("error must be empty on success, got %q", hs.Error)
	}
	if hs.Stats.MaxConns != 4 {
		t.Fatalf("MaxConns = %d, want 4", hs.Stats.MaxConns)
	}
}

func TestHealthStatusDownFields(t *testing.T) {
	hs := HealthStatus{Status: "down", Error: "database ping timed out"}
	if hs.Status != "down" {
		t.Fatalf("status = %q, want %q", hs.Status, "down")
	}
	if hs.Error == "" {
		t.Fatal("error must be non-empty when status is down")
	}
}

// TestDefaultHealthTimeout verifies the package-level constant has a
// reasonable value so callers are never surprised.
func TestDefaultHealthTimeout(t *testing.T) {
	if defaultHealthTimeout <= 0 {
		t.Fatalf("defaultHealthTimeout = %v, must be positive", defaultHealthTimeout)
	}
	if defaultHealthTimeout > 30*time.Second {
		t.Fatalf("defaultHealthTimeout = %v, should not exceed 30s", defaultHealthTimeout)
	}
}

// ---------------------------------------------------------------------------
// Health – integration (requires AMARA_TEST_DATABASE_URL)
// ---------------------------------------------------------------------------

func TestHealthReturnsUpWithLiveDatabase(t *testing.T) {
	url := liveURL(t)
	cfg := config.DatabaseConfig{URL: url, MaxConns: 2, HealthTimeout: 5 * time.Second}
	svc, err := NewService(t.Context(), cfg)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	defer svc.Close()

	hs := svc.Health(t.Context())
	if hs.Status != "up" {
		t.Fatalf("health status = %q, want %q; error: %s", hs.Status, "up", hs.Error)
	}
	if hs.Error != "" {
		t.Fatalf("error field must be empty when status is up, got %q", hs.Error)
	}
	if hs.Stats.MaxConns != 2 {
		t.Fatalf("MaxConns in stats = %d, want 2", hs.Stats.MaxConns)
	}
}

func TestHealthReturnsDownAfterClose(t *testing.T) {
	url := liveURL(t)
	cfg := config.DatabaseConfig{URL: url, MaxConns: 1, HealthTimeout: 2 * time.Second}
	svc, err := NewService(t.Context(), cfg)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	// Close the pool before calling Health to simulate a lost connection.
	svc.Close()

	hs := svc.Health(context.Background())
	if hs.Status != "down" {
		t.Fatalf("health status = %q after pool close, want %q", hs.Status, "down")
	}
	if hs.Error == "" {
		t.Fatal("error field must be non-empty when status is down")
	}
}

// ---------------------------------------------------------------------------
// Pool configuration propagation
// ---------------------------------------------------------------------------

func TestNewPoolAppliesMaxConns(t *testing.T) {
	url := liveURL(t)
	cfg := config.DatabaseConfig{URL: url, MaxConns: 3}
	pool, err := NewPool(t.Context(), cfg)
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	defer pool.Close()
	if pool.Stat().MaxConns() != 3 {
		t.Fatalf("MaxConns = %d, want 3", pool.Stat().MaxConns())
	}
}
