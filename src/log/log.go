// Package log contains logging utilities.
package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	levelErr   = "err"
	levelWarn  = "warn"
	levelInfo  = "info"
	levelDebug = "debug"
)

// Config encapsulates the logging config.
type Config struct {
	Level string `env:"LEVEL"`
}

// DefaultConfig returns a default log config.
func DefaultConfig() Config {
	return Config{Level: "info"}
}

// NewLogger creates a new Zap logger and configures it with the given Config.
func NewLogger(c Config) *zap.Logger {
	var level zapcore.Level
	switch strings.ToLower(c.Level) {
	case levelErr:
		level = zap.ErrorLevel
	case levelWarn:
		level = zap.WarnLevel
	case levelDebug:
		level = zap.DebugLevel
	case levelInfo:
		fallthrough
	default:
		level = zap.InfoLevel
	}
	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(level),
	))
	return log
}
