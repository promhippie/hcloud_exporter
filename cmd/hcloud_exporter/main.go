package main

import (
	"errors"
	"os"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/joho/godotenv"
	"github.com/promhippie/hcloud_exporter/pkg/action"
	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/promhippie/hcloud_exporter/pkg/version"
	"gopkg.in/urfave/cli.v2"
)

var (
	// ErrMissingHcloudToken defines the error if hcloud.token is empty.
	ErrMissingHcloudToken = errors.New("Missing required hcloud.token")
)

func main() {
	cfg := config.Load()

	if env := os.Getenv("HCLOUD_EXPORTER_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := &cli.App{
		Name:    "hcloud_exporter",
		Version: version.Version,
		Usage:   "HetznerCloud Exporter",
		Authors: []*cli.Author{
			{
				Name:  "Thomas Boerger",
				Email: "thomas@webhippie.de",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log.level",
				Value:       "info",
				Usage:       "Only log messages with given severity",
				EnvVars:     []string{"HCLOUD_EXPORTER_LOG_LEVEL"},
				Destination: &cfg.Logs.Level,
			},
			&cli.BoolFlag{
				Name:        "log.pretty",
				Value:       false,
				Usage:       "Enable pretty messages for logging",
				EnvVars:     []string{"HCLOUD_EXPORTER_LOG_PRETTY"},
				Destination: &cfg.Logs.Pretty,
			},
			&cli.StringFlag{
				Name:        "web.address",
				Value:       "0.0.0.0:9108",
				Usage:       "Address to bind the metrics server",
				EnvVars:     []string{"HCLOUD_EXPORTER_WEB_ADDRESS"},
				Destination: &cfg.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "web.path",
				Value:       "/metrics",
				Usage:       "Path to bind the metrics server",
				EnvVars:     []string{"HCLOUD_EXPORTER_WEB_PATH"},
				Destination: &cfg.Server.Path,
			},
			&cli.DurationFlag{
				Name:        "request.timeout",
				Value:       5 * time.Second,
				Usage:       "Request timeout as duration",
				EnvVars:     []string{"HCLOUD_EXPORTER_REQUEST_TIMEOUT"},
				Destination: &cfg.Target.Timeout,
			},
			&cli.StringFlag{
				Name:        "hcloud.token",
				Value:       "",
				Usage:       "Access token for the HetznerCloud API",
				EnvVars:     []string{"HCLOUD_EXPORTER_TOKEN"},
				Destination: &cfg.Target.Token,
			},
		},
		Action: func(c *cli.Context) error {
			logger := setupLogger(cfg)

			if cfg.Target.Token == "" {
				level.Error(logger).Log(
					"msg", ErrMissingHcloudToken,
				)

				return ErrMissingHcloudToken
			}

			return action.Server(cfg, logger)
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the current version of that tool",
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
