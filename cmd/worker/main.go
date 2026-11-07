package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	invwdatabase "github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
	"github.com/CSCE-3701-Amara/amara/internal/invw/adapters/judge0"
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

	if baseURL := strings.TrimSpace(os.Getenv("JUDGE0_URL")); baseURL != "" {
		var languages map[string]int
		if raw := os.Getenv("JUDGE0_LANGUAGE_IDS"); raw != "" {
			if err := json.Unmarshal([]byte(raw), &languages); err != nil || len(languages) == 0 {
				return fmt.Errorf("JUDGE0_LANGUAGE_IDS must be a JSON object mapping languages to IDs")
			}
		}
		runner, err := judge0.New(baseURL, os.Getenv("JUDGE0_AUTH_TOKEN"), languages)
		if err != nil {
			return err
		}
		var grading sync.WaitGroup
		grading.Go(func() {
			gradeCode(ctx, logger, invwdatabase.NewAssessmentRepository(dbSvc.Pool()), runner)
		})
		defer grading.Wait()
	}

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

func gradeCode(ctx context.Context, logger *slog.Logger, repo *invwdatabase.AssessmentRepository, runner invwdatabase.CodeRunner) {
	for ctx.Err() == nil {
		claimed, err := repo.GradeNextCodeAttempt(ctx, runner)
		if err != nil && ctx.Err() == nil {
			logger.Error("code grading failed", "error", err)
		}
		if claimed {
			continue
		}
		timer := time.NewTimer(5 * time.Second)
		select {
		case <-ctx.Done():
		case <-timer.C:
		}
		timer.Stop()
	}
}
