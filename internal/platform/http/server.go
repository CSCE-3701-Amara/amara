package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/platform/config"
	"github.com/CSCE-3701-Amara/amara/internal/platform/database"
)

// Server owns the API transport and its shutdown lifecycle.
type Server struct {
	server          *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

// New composes the router and HTTP server. Module dependencies are constructed
// in cmd/api and supplied through route registrars.
func New(cfg config.HTTPConfig, logger *slog.Logger, dbSvc database.Service, registrars ...RouteRegistrar) (*Server, error) {
	if logger == nil {
		logger = slog.Default()
	}
	router, err := newRouter(cfg, logger, dbSvc, registrars...)
	if err != nil {
		return nil, err
	}

	return &Server{
		server: &http.Server{
			Addr:              cfg.Address,
			Handler:           router,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
		},
		logger:          logger,
		shutdownTimeout: cfg.ShutdownTimeout,
	}, nil
}

// Run serves requests until ctx is canceled or the listener fails. Cancellation
// drains active requests within the configured timeout before closing them.
func (s *Server) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return nil
	}
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	return s.serve(ctx, listener)
}

func (s *Server) serve(ctx context.Context, listener net.Listener) error {
	s.logger.Info("HTTP server started", "address", listener.Addr().String())

	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- s.server.Serve(listener)
	}()

	select {
	case err := <-serveErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		s.logger.Info("HTTP server shutting down")
	}

	// The signal context is already canceled; draining needs a fresh deadline.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		closeErr := s.server.Close()
		return errors.Join(fmt.Errorf("shutdown HTTP server: %w", err), closeErr)
	}
	if err := <-serveErrors; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	s.logger.Info("HTTP server stopped")
	return nil
}
