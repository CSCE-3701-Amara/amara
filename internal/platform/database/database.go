// Package database builds shared database connection pools. It contains
// connection and transaction infrastructure only; business queries belong in
// the epic modules' own repository adapters.
package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
)

// PingTimeout bounds the startup connectivity check applied to every pool.
const PingTimeout = 10 * time.Second

// defaultHealthTimeout is used when the configuration does not supply one.
const defaultHealthTimeout = 5 * time.Second

// Service manages the lifecycle of the shared database connection pool and
// exposes health information. It does not perform business queries.
type Service interface {
	// Health pings the database and returns a snapshot of pool statistics.
	// It never panics, calls os.Exit, or log.Fatal. A failed ping produces an
	// unhealthy result; it does not terminate the application.
	Health(ctx context.Context) HealthStatus

	// Pool returns the underlying connection pool for injection into
	// repository adapters. Business services must not receive this directly;
	// they depend on repository interfaces (ports) instead.
	Pool() *pgxpool.Pool

	// Close drains and closes the pool. Safe to call once during shutdown
	// after all in-flight requests have finished.
	Close()
}

// HealthStatus carries the result of a single database health check.
// Fields use typed values rather than map[string]string so callers can branch
// on them without parsing.
type HealthStatus struct {
	// Status is "up" when the database ping succeeds, "down" otherwise.
	Status string `json:"status"`
	// Error holds the ping error message when Status is "down". Empty on success.
	// Never contains connection strings or credentials.
	Error string `json:"error,omitempty"`
	// Stats is a snapshot of pool counters at the time of the check.
	Stats PoolStats `json:"stats"`
}

// PoolStats is a typed snapshot of pgxpool.Stat counters.
type PoolStats struct {
	TotalConns        int32         `json:"total_conns"`
	AcquiredConns     int32         `json:"acquired_conns"`
	IdleConns         int32         `json:"idle_conns"`
	ConstructingConns int32         `json:"constructing_conns"`
	MaxConns          int32         `json:"max_conns"`
	AcquireCount      int64         `json:"acquire_count"`
	AcquireDuration   time.Duration `json:"acquire_duration_ns"`
	EmptyAcquireCount int64         `json:"empty_acquire_count"`
}

// service is the concrete implementation of Service.
type service struct {
	pool          *pgxpool.Pool
	healthTimeout time.Duration
}

// NewService opens a connection pool, verifies connectivity, and returns a
// Service that owns the pool lifecycle. The caller must call Close() during
// shutdown after all dependent components have stopped using the pool.
func NewService(ctx context.Context, cfg config.DatabaseConfig) (Service, error) {
	pool, err := newPool(ctx, cfg)
	if err != nil {
		return nil, err
	}
	ht := cfg.HealthTimeout
	if ht <= 0 {
		ht = defaultHealthTimeout
	}
	return &service{pool: pool, healthTimeout: ht}, nil
}

// Health pings the database within a bounded timeout and captures pool statistics.
// It returns an unhealthy status on failure rather than propagating the error,
// so callers (HTTP handlers) can encode the result without aborting.
func (s *service) Health(ctx context.Context) HealthStatus {
	pingCtx, cancel := context.WithTimeout(ctx, s.healthTimeout)
	defer cancel()

	stat := s.pool.Stat()
	stats := PoolStats{
		TotalConns:        stat.TotalConns(),
		AcquiredConns:     stat.AcquiredConns(),
		IdleConns:         stat.IdleConns(),
		ConstructingConns: stat.ConstructingConns(),
		MaxConns:          stat.MaxConns(),
		AcquireCount:      stat.AcquireCount(),
		AcquireDuration:   stat.AcquireDuration(),
		EmptyAcquireCount: stat.EmptyAcquireCount(),
	}

	if err := s.pool.Ping(pingCtx); err != nil {
		// Suppress context-cancellation noise: use a short generic phrase so
		// callers can still distinguish a hard failure from a timeout.
		msg := "database ping failed"
		if errors.Is(err, context.DeadlineExceeded) {
			msg = "database ping timed out"
		}
		return HealthStatus{Status: "down", Error: msg, Stats: stats}
	}
	return HealthStatus{Status: "up", Stats: stats}
}

// Pool exposes the underlying pool for injection into repository adapters.
func (s *service) Pool() *pgxpool.Pool {
	return s.pool
}

// Close drains and closes the underlying pool.
func (s *service) Close() {
	s.pool.Close()
}

// newPool is the internal constructor shared by NewService and the legacy
// NewPool helper retained for backward compatibility.
func newPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	poolCfg.MaxConns = cfg.MaxConns
	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetime > 0 {
		poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, PingTimeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

// NewPool opens a connection pool for the configured database and verifies
// connectivity before returning. Callers own the returned pool and must Close
// it during shutdown.
//
// Prefer NewService over NewPool: it provides lifecycle management and health
// checks in addition to the raw pool.
func NewPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	return newPool(ctx, cfg)
}
