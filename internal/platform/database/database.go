// Package database owns shared database connection and transaction
// infrastructure. Business queries belong in module database adapters.
package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
)

const (
	PingTimeout          = 10 * time.Second
	defaultHealthTimeout = 5 * time.Second
)

// Service owns the shared pool lifecycle and reports readiness information.
// Module repositories receive Pool(); business application services receive
// their module-owned repository interfaces instead.
type Service interface {
	Health(context.Context) HealthStatus
	Pool() *pgxpool.Pool
	Close()
}

type HealthStatus struct {
	Status string    `json:"status"`
	Error  string    `json:"error,omitempty"`
	Stats  PoolStats `json:"stats"`
}

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

type service struct {
	pool          *pgxpool.Pool
	healthTimeout time.Duration
}

// NewService creates and verifies the shared pool. Call Close during shutdown.
func NewService(ctx context.Context, cfg config.DatabaseConfig) (Service, error) {
	pool, err := newPool(ctx, cfg)
	if err != nil {
		return nil, err
	}
	healthTimeout := cfg.HealthTimeout
	if healthTimeout <= 0 {
		healthTimeout = defaultHealthTimeout
	}
	return &service{pool: pool, healthTimeout: healthTimeout}, nil
}

func (s *service) Health(ctx context.Context) HealthStatus {
	pingCtx, cancel := context.WithTimeout(ctx, s.healthTimeout)
	defer cancel()

	stat := s.pool.Stat()
	stats := PoolStats{
		TotalConns: stat.TotalConns(), AcquiredConns: stat.AcquiredConns(),
		IdleConns: stat.IdleConns(), ConstructingConns: stat.ConstructingConns(),
		MaxConns: stat.MaxConns(), AcquireCount: stat.AcquireCount(),
		AcquireDuration: stat.AcquireDuration(), EmptyAcquireCount: stat.EmptyAcquireCount(),
	}
	if err := s.pool.Ping(pingCtx); err != nil {
		message := "database ping failed"
		if errors.Is(err, context.DeadlineExceeded) {
			message = "database ping timed out"
		}
		return HealthStatus{Status: "down", Error: message, Stats: stats}
	}
	return HealthStatus{Status: "up", Stats: stats}
}

func (s *service) Pool() *pgxpool.Pool { return s.pool }
func (s *service) Close()              { s.pool.Close() }

func newPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
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

// NewPool is retained for existing callers. New composition code should use
// NewService so the pool lifecycle and health contract remain together.
func NewPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	return newPool(ctx, cfg)
}
