package http

import "time"

// Config holds an HTTP server's configuration.
type Config struct {
	Addr            string        `env:"ADDR"`
	AdvertisedAddr  string        `env:"ADV_ADDR"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}
