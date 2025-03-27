package srv

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ServerConfig defines configuration for launching the HTTP server.
type ServerConfig struct {
	TLSCertPath string // Path to the TLS certificate file
	TLSKeyPath  string // Path to the TLS key file
	Port        int    // Port to bind the HTTP server to
	TLSEnabled  bool   // Enable TLS if true
}

// CleanupHandler is an interface that allows components (e.g., DB, workers)
// to hook into the server shutdown and clean up resources gracefully.
type CleanupHandler interface {
	Shutdown(ctx context.Context) error
}

// server wraps an http.Handler and its configuration.
type server struct {
	handler http.Handler
	config  ServerConfig
}

// NewServer creates a new server instance with the given handler and configuration.
func NewServer(handler http.Handler, config ServerConfig) *server {
	return &server{
		handler: handler,
		config:  config,
	}
}

// StartWithGracefulShutdown starts the HTTP server and listens for SIGINT/SIGTERM
// to shut down gracefully. It runs cleanup handlers in parallel and shuts down
// the HTTP server within the given timeout.
//
// Accepts a parent context for integration into external lifecycle systems.
func (s *server) StartWithGracefulShutdown(
	parentCtx context.Context,
	timeout time.Duration,
	handlers ...CleanupHandler,
) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.handler,
	}

	// Setup TLS if enabled
	if s.config.TLSEnabled {
		cert, err := tls.LoadX509KeyPair(s.config.TLSCertPath, s.config.TLSKeyPath)
		if err != nil {
			slog.Error("Failed to load TLS certificate", slog.String("component", "http-server"), slog.Any("error", err))
			os.Exit(1)
		}
		srv.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	// Create context that cancels on interrupt signals
	ctx, stop := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the server in background
	go func() {
		slog.Info("Starting server", slog.String("component", "http-server"), slog.String("addr", srv.Addr))
		var err error
		if s.config.TLSEnabled {
			err = srv.ListenAndServeTLS("", "")
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", slog.String("component", "http-server"), slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received", slog.String("component", "http-server"))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			slog.Error("Graceful shutdown timed out, forcing exit", slog.String("component", "http-server"))
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	for _, h := range handlers {
		wg.Add(1)
		go func(handler CleanupHandler) {
			defer wg.Done()
			if err := handler.Shutdown(shutdownCtx); err != nil {
				slog.Error("Cleanup handler failed", slog.String("component", "http-server"), slog.Any("error", err))
			}
		}(h)
	}
	wg.Wait()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown failed", slog.String("component", "http-server"), slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Shutdown complete", slog.String("component", "http-server"))
}

// Start starts the HTTP server without signal handling or graceful shutdown.
// It is intended for use in test scenarios.
func (s *server) Start() (*http.Server, error) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.handler,
	}

	go func() {
		slog.Info("Starting server (test mode)", slog.String("component", "http-server"), slog.String("addr", srv.Addr))
		var err error
		if s.config.TLSEnabled {
			err = srv.ListenAndServeTLS(s.config.TLSCertPath, s.config.TLSKeyPath)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Test server failed", slog.String("component", "http-server"), slog.Any("error", err))
		}
	}()

	return srv, nil
}
