// Package logger provides a centralized setup for the application's slog logger.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Setup initializes the global slog logger with the given level.
func Setup(logLevel string) *slog.Logger {
	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // Default to INFO
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
