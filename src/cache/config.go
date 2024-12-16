package cache

import "time"

type Config struct {
	Address      string        `env:"ADDR"`
	Password     string        `env:"PASS"`
	TTL          time.Duration `env:"TTL"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
}

func DefaultConfig() Config {
	return Config{
		TTL:          24 * time.Hour,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 1 * time.Second,
	}
}
