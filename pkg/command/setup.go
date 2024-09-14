package command

import (
	"log/slog"
	"os"
	"strings"

	"github.com/promhippie/hcloud_exporter/pkg/config"
)

func setupLogger(cfg *config.Config) *slog.Logger {
	if cfg.Logs.Pretty {
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: loggerLevel(cfg),
			}),
		)
	}

	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: loggerLevel(cfg),
		}),
	)
}

func loggerLevel(cfg *config.Config) slog.Leveler {
	switch strings.ToLower(cfg.Logs.Level) {
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	}

	return slog.LevelInfo
}
