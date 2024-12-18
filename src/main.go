package main

import (
	"context"

	"github.com/0x5d/hash/api/http"
	"github.com/0x5d/hash/cache"
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
	urlRepo, err := persistence.NewPGURLRepo(logger, ctx, c.DB)
	if err != nil {
		logger.Fatal("Failed to start URL repo", zap.Error(err))
	}
	cache, err := cache.NewRedisCache(ctx, c.Cache)
	if err != nil {
		logger.Fatal("Failed to connect to cache", zap.Error(err))
	}
	logger.Info("Table", zap.String("table", c.DB.MigrationsTable))
	err = urlRepo.Migrate(ctx)
	if err != nil {
		logger.Fatal("Failed to apply migrations", zap.Error(err))
	}
	urlSvc := core.NewURLService(logger, urlRepo, cache)
	s := http.NewServer(c.HTTP, urlSvc, logger.Named("http"))
	s.Start(ctx)
}
