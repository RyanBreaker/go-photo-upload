package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	if os.Getenv("ENV") == "production" {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	} else {
		options := &slog.HandlerOptions{Level: slog.LevelDebug}
		Logger = slog.New(slog.NewTextHandler(os.Stdout, options))
	}
}
