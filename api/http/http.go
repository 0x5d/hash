// Package http contains all the configuration and functionality related to the HTTP layer.
package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
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

type errResponse struct {
	Msg string `json:"msg"`
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
		advertisedAddr: s.c.AdvertisedAddr,
		urlSvc:         &s.urlSvc,
		log:            s.log.Named("url"),
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

func shiftPath(p string) (head, tail string) {
	p = path.Clean(slash + p)
	i := strings.Index(p[1:], slash) + 1
	if i <= 0 {
		return p[1:], slash
	}
	return p[1:i], p[i:]
}

func parseJSON(log *zap.Logger, res http.ResponseWriter, body io.ReadCloser, v any) error {
	err := json.NewDecoder(body).Decode(&v)
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, &http.MaxBytesError{}):
		errMsg := fmt.Sprintf("Request body should be under %dB", bodyLimit)
		log.Warn(errMsg, zap.Error(err))
		writeErrRes(res, errMsg, http.StatusRequestEntityTooLarge)
	default:
		errMsg := "Invalid JSON"
		log.Warn(errMsg, zap.Error(err))
		writeErrRes(res, "Invalid JSON", http.StatusBadRequest)
	}
	return err
}

func writeErrRes(res http.ResponseWriter, msg string, status int) {
	errRes := errResponse{Msg: msg}
	r, err := json.Marshal(&errRes)
	if err != nil {
		http.Error(res, "Failed to encode response.", http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	http.Error(res, string(r), status)
}
