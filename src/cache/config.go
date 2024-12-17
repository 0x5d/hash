package cache

import "time"

type Config struct {
	// The cache address.
	Address string `env:"ADDR"`
	// The cache password.
	Password string `env:"PASS"`
	// The cached objects' time-to-live.
	TTL time.Duration `env:"TTL"`
	// The timeout when reading from the cache.
	ReadTimeout time.Duration `env:"READ_TIMEOUT"`
	// The timeout when writing to the cache.
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
}

func DefaultConfig() Config {
	return Config{
		TTL:          24 * time.Hour,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 1 * time.Second,
	}
}
