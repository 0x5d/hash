// Package http contains all the configuration and functionality related to the HTTP layer.
package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0x5d/hash/core"
	"go.uber.org/zap"
)

// Server is an HTTP server.
type Server struct {
	c      Config
	urlSvc core.URLService
	log    *zap.Logger
}

// Config holds an HTTP server's configuration.
type Config struct {
	Addr            string        `env:"ADDR"`
	AdvertisedAddr  string        `env:"ADV_ADDR"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{Addr: ":8080", ShutdownTimeout: 10 * time.Second}
}

// NewServer returns a new Server instance configured with the given Config.
func NewServer(c Config, urlSvc core.URLService, log *zap.Logger) *Server {
	return &Server{c: c, urlSvc: urlSvc, log: log}
}

// Start starts the server. It handles shutdown gracefully by not accepting new connections after
// it receives a SIGINT or SIGTERM.
func (s *Server) Start(ctx context.Context) {
	server := &http.Server{
		Addr: s.c.Addr,
	}

	http.Handle("/url", &urlRouter{
		urlSvc: &s.urlSvc,
		log:    s.log.Named("url"),
	})

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal("HTTP server shutdown unexpectedly", zap.Error(err))
		}
		s.log.Info("HTTP server stopped")
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	shutdownCtx, cancel := context.WithTimeout(ctx, s.c.ShutdownTimeout)
	defer cancel()

	s.log.Info("HTTP server Shutting down")
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		s.log.Fatal("HTTP server failed to shutdown", zap.Error(err))
	}
}
