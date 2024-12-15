package config

import (
	"context"

	"github.com/0x5d/hash/api/http"
	"github.com/0x5d/hash/log"
	"github.com/0x5d/hash/persistence"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	HTTP http.Config        `env:", prefix=HTTP_"`
	DB   persistence.Config `env:", prefix=DB_"`
	Log  log.Config         `env:", prefix=LOG_"`
}

func LoadFromEnv(ctx context.Context) (Config, error) {
	c := Config{
		HTTP: http.DefaultConfig(),
		Log:  log.DefaultConfig(),
	}
	err := envconfig.Process(ctx, &c)
	return c, err
}
