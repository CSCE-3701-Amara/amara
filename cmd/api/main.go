package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
	"github.com/CSCE-3701-Amara/amara/internal/platform/database"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
	"github.com/CSCE-3701-Amara/amara/internal/platform/logging"
)

func main() {
	if err := run(); err != nil {
		logging.New(os.Stderr, slog.LevelError).Error("API stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger := logging.New(os.Stdout, cfg.LogLevel).With(
		"service", "amara-api",
		"environment", cfg.Environment,
	)
	slog.SetDefault(logger)

	gin.SetMode(gin.ReleaseMode)
	if cfg.Environment == "development" {
		gin.SetMode(gin.DebugMode)
	}

	// Construct module services and pass their HTTP route registrars here.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Fail fast: the API must not accept traffic it cannot serve.
	pool, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()
	logger.Info("database connection ready")

	server, err := httpserver.New(cfg.HTTP, logger)
	if err != nil {
		return fmt.Errorf("create HTTP server: %w", err)
	}

	return server.Run(ctx)
}
