package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	authdatabase "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/database"
	authhttp "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/http"
	authpassword "github.com/CSCE-3701-Amara/amara/internal/auth/adapters/password"
	authapplication "github.com/CSCE-3701-Amara/amara/internal/auth/application"
	canddatabase "github.com/CSCE-3701-Amara/amara/internal/cand/adapters/database"
	candhttp "github.com/CSCE-3701-Amara/amara/internal/cand/adapters/http"
	candapplication "github.com/CSCE-3701-Amara/amara/internal/cand/application"
	invwdatabase "github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
	invwhttp "github.com/CSCE-3701-Amara/amara/internal/invw/adapters/http"
	invwapplication "github.com/CSCE-3701-Amara/amara/internal/invw/application"
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

	authRepo, err := authdatabase.NewRegistrationRepository(dbSvc.Pool())
	if err != nil {
		return fmt.Errorf("create auth repository: %w", err)
	}
	passwordChecker, err := authpassword.NewHIBPChecker(&http.Client{Timeout: 10 * time.Second})
	if err != nil {
		return fmt.Errorf("create password compromise checker: %w", err)
	}
	companyRegistration, err := authapplication.NewCompanyRegistrationService(authRepo, passwordChecker)
	if err != nil {
		return fmt.Errorf("create company registration service: %w", err)
	}
	companyRegistrationHandler, err := authhttp.NewCompanyRegistrationHandler(companyRegistration)
	if err != nil {
		return fmt.Errorf("create company registration handler: %w", err)
	}
	mfaService, err := authapplication.NewMFAService(authRepo, cfg.Auth.MFAEncryptionKey)
	if err != nil {
		return fmt.Errorf("create MFA service: %w", err)
	}
	loginService, err := authapplication.NewLoginService(authRepo, authRepo, mfaService)
	if err != nil {
		return fmt.Errorf("create login service: %w", err)
	}
	loginHandler, err := authhttp.NewLoginHandlerWithLogger(loginService, logger)
	if err != nil {
		return fmt.Errorf("create login handler: %w", err)
	}
	logoutService, err := authapplication.NewLogoutService(authRepo)
	if err != nil {
		return fmt.Errorf("create logout service: %w", err)
	}
	logoutHandler, err := authhttp.NewLogoutHandler(logoutService, authRepo, logger)
	if err != nil {
		return fmt.Errorf("create logout handler: %w", err)
	}
	accountActionService, err := authapplication.NewAccountActionService(authRepo)
	if err != nil {
		return fmt.Errorf("create account action service: %w", err)
	}
	accountActionHandler, err := authhttp.NewAccountActionHandler(accountActionService, authRepo, logger)
	if err != nil {
		return fmt.Errorf("create account action handler: %w", err)
	}

	jobLifecycleRepo := jobdatabase.NewPostingLifecycleRepository(dbSvc.Pool())
	jobLifecycle := jobapplication.NewPostingLifecycle(jobLifecycleRepo)

	postingListRepo := jobdatabase.NewCompanyPostingListRepository(dbSvc.Pool())
	postingListHandler := jobhttp.NewPostingListHandler(jobapplication.NewCompanyPostingList(postingListRepo, jobLifecycle), authRepo)
	postingRepo := jobdatabase.NewPostingManagementRepository(dbSvc.Pool())
	postingService := jobapplication.NewPostingManagement(postingRepo)
	// AUTH integration placeholder: the current checkout has no session middleware.
	// The handler fails closed until AUTH supplies a middleware that sets its principal.
	postingHandler := jobhttp.NewPostingHandler(postingService, nil)
	applicationRepo := jobdatabase.NewApplicationRepository(dbSvc.Pool())
	applicationService := jobapplication.NewApplicationAdmission(applicationRepo)
	applicationHandler := jobhttp.NewApplicationHandler(applicationService, nil)

	feedRepo := canddatabase.NewJobFeedRepository(dbSvc.Pool())
	feedHandler := candhttp.NewJobFeedHandler(
		candapplication.NewJobFeed(feedRepo, jobLifecycle),
		authhttp.CandidateAuth(authRepo),
	)
	pipelineHandler := invwhttp.NewPipelineHandler(invwdatabase.NewPipelineRepository(dbSvc.Pool()), authRepo)
	moduleHandler := invwhttp.NewModuleHandler(invwdatabase.NewModuleRepository(dbSvc.Pool()), authRepo)
	trackingRepo := canddatabase.NewApplicationTrackingRepository(dbSvc.Pool())
	trackingHandler := candhttp.NewApplicationTrackingHandler(trackingRepo, authRepo)

	pipelineProgressRepo := invwdatabase.NewCandidatePipelineProgressRepository(dbSvc.Pool())
	pipelineProgressHandler := invwhttp.NewCandidatePipelineProgressHandler(
		invwapplication.NewCandidatePipelineProgress(pipelineProgressRepo),
		nil,
	)

	// Candidate/interviewer registration needs a verification sender. Until that
	// provider is configured, reject requests before creating unusable accounts.
	server, err := httpserver.New(cfg.HTTP, logger, dbSvc,
		authhttp.RegisterUnavailableRegistrationRoutes,
		companyRegistrationHandler.RegisterRoutes,
		loginHandler.RegisterRoutes,
		logoutHandler.RegisterRoutes,
		accountActionHandler.RegisterRoutes,
		feedHandler.RegisterRoutes,
		postingHandler.RegisterRoutes,
		postingListHandler.RegisterRoutes,
		applicationHandler.RegisterRoutes,
		moduleHandler.RegisterRoutes,
		trackingHandler.RegisterRoutes,
		pipelineHandler.RegisterRoutes,
		pipelineProgressHandler.RegisterRoutes,
	)
	if err != nil {
		return fmt.Errorf("create HTTP server: %w", err)
	}

	return server.Run(ctx)
}
