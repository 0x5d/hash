package main

import (
	"context"

	"github.com/0x5d/hash/api/http"
	"github.com/0x5d/hash/config"
	"github.com/0x5d/hash/core"
	"github.com/0x5d/hash/log"
	"github.com/0x5d/hash/persistence"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	c, err := config.LoadFromEnv(ctx)
	logger := log.NewLogger(c.Log)
	if err != nil {
		logger.Fatal("Failed to load config from env", zap.Error(err))
	}
	urlRepo, err := persistence.NewPGURLRepo(ctx, c.DB)
	if err != nil {
		logger.Fatal("Failed to start URL repo", zap.Error(err))
	}
	urlSvc := core.NewURLService(urlRepo)
	s := http.NewServer(c.HTTP, urlSvc, logger)
	s.Start(ctx)
}
