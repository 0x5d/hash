package main

import (
	"context"

	"github.com/0x5d/hash/api/http"
	"github.com/0x5d/hash/config"
	"github.com/0x5d/hash/core"
	"github.com/0x5d/hash/log"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	c, err := config.LoadFromEnv(ctx)
	logger := log.NewLogger(c.Log)
	if err != nil {
		logger.Fatal("Failed to load config from env", zap.Error(err))
	}
	shortener := core.NewShortener(c.Core)
	s := http.NewServer(c.HTTP, shortener, logger)
	s.Start(ctx)
}
