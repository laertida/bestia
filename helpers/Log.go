package helpers

import (
	"fmt"
	"log/slog"
	"os"
)

func GetLogLevel(level string) slog.Level {
	var logLevel slog.Level

	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		fmt.Errorf("invalid log level: %s", level)
		os.Exit(1)
	}
	return logLevel
}
