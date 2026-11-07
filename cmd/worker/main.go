package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	invwdatabase "github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
	"github.com/CSCE-3701-Amara/amara/internal/platform/database"
	"github.com/CSCE-3701-Amara/amara/internal/platform/logging"
)

func main() {
	if err := run(); err != nil {
		logging.New(os.Stderr, slog.LevelError, "text").Error("Worker stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger := logging.New(os.Stdout, cfg.LogLevel, cfg.LogFormat).With(
		"service", "amara-worker",
		"environment", cfg.Environment,
	)
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbSvc, err := database.NewService(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer dbSvc.Close()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		if err := invwdatabase.SpawnAssessmentInstances(ctx, dbSvc.Pool()); err != nil && ctx.Err() == nil {
			logger.Error("assessment instance sweep failed", "error", err)
		}
		if err := invwdatabase.ReconcileAdvancement(ctx, dbSvc.Pool()); err != nil && ctx.Err() == nil {
			logger.Error("stage advancement sweep failed", "error", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}

}
