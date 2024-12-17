package http

import "time"

// Config holds an HTTP server's configuration.
type Config struct {
	// The address <ip>:<port> to bind to.
	Addr string `env:"ADDR"`
	// The advertised address, which will be used in responses when creating shortened URLs.
	AdvertisedAddr string `env:"ADV_ADDR"`
	// The timeout when waiting for the server to shut down.
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}
