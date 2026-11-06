package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	canddatabase "github.com/CSCE-3701-Amara/amara/internal/cand/adapters/database"
	candhttp "github.com/CSCE-3701-Amara/amara/internal/cand/adapters/http"
	candapplication "github.com/CSCE-3701-Amara/amara/internal/cand/application"
	jobdatabase "github.com/CSCE-3701-Amara/amara/internal/job/adapters/database"
	jobhttp "github.com/CSCE-3701-Amara/amara/internal/job/adapters/http"
	jobapplication "github.com/CSCE-3701-Amara/amara/internal/job/application"
	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
	"github.com/CSCE-3701-Amara/amara/internal/platform/database"
	httpserver "github.com/CSCE-3701-Amara/amara/internal/platform/http"
	"github.com/CSCE-3701-Amara/amara/internal/platform/logging"
)

func main() {
	if err := run(); err != nil {
		// Configuration may be unreadable; aggregation still parses JSON.
		logging.New(os.Stderr, slog.LevelError, "text").Error("API stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger := logging.New(os.Stdout, cfg.LogLevel, cfg.LogFormat).With(
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
	dbSvc, err := database.NewService(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer dbSvc.Close()
	logger.Info("database connection ready")

	jobLifecycleRepo := jobdatabase.NewPostingLifecycleRepository(dbSvc.Pool())
	jobLifecycle := jobapplication.NewPostingLifecycle(jobLifecycleRepo)
	postingRepo := jobdatabase.NewPostingManagementRepository(dbSvc.Pool())
	postingService := jobapplication.NewPostingManagement(postingRepo)
	// AUTH integration placeholder: the current checkout has no session middleware.
	// The handler fails closed until AUTH supplies a middleware that sets its principal.
	postingHandler := jobhttp.NewPostingHandler(postingService, nil)
	applicationRepo := jobdatabase.NewApplicationRepository(dbSvc.Pool())
	applicationService := jobapplication.NewApplicationAdmission(applicationRepo)
	applicationHandler := jobhttp.NewApplicationHandler(applicationService, nil)
	feedRepo := canddatabase.NewJobFeedRepository(dbSvc.Pool())
	feedHandler := candhttp.NewJobFeedHandler(candapplication.NewJobFeed(feedRepo, jobLifecycle), nil)
	// Frontend integration placeholder: this REST API is the contract for the future posting editor.
	server, err := httpserver.New(cfg.HTTP, logger, dbSvc, feedHandler.RegisterRoutes, postingHandler.RegisterRoutes, applicationHandler.RegisterRoutes)
	if err != nil {
		return fmt.Errorf("create HTTP server: %w", err)
	}

	ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return server.Run(ctx)
}
