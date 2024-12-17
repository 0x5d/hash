package config

import (
	"context"

	"github.com/0x5d/hash/api/http"
	"github.com/0x5d/hash/cache"
	"github.com/0x5d/hash/log"
	"github.com/0x5d/hash/persistence"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	// The HTTP server/ API config.
	HTTP http.Config `env:", prefix=HTTP_"`
	// The database config.
	DB persistence.Config `env:", prefix=DB_"`
	// The logging config.
	Log log.Config `env:", prefix=LOG_"`
	// The cache config.
	Cache cache.Config `env:", prefix=CACHE_"`
}

// Loads the config from the fields' respective env vars.
func LoadFromEnv(ctx context.Context) (Config, error) {
	c := Config{
		HTTP:  http.DefaultConfig(),
		Log:   log.DefaultConfig(),
		Cache: cache.DefaultConfig(),
	}
	err := envconfig.Process(ctx, &c)
	return c, err
}
