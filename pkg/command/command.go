package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/promhippie/hcloud_exporter/pkg/action"
	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/promhippie/hcloud_exporter/pkg/version"
	"github.com/urfave/cli/v3"
)

// Run parses the command line arguments and executes the program.
func Run() error {
	cfg := config.Load()

	app := &cli.Command{
		Name:    "hcloud_exporter",
		Version: version.String,
		Usage:   "Hetzner Cloud Exporter",
		Authors: []any{
			"Thomas Boerger <thomas@webhippie.de>",
		},
		Flags: RootFlags(cfg),
		Commands: []*cli.Command{
			Health(cfg),
		},
		Action: func(_ context.Context, _ *cli.Command) error {
			logger := setupLogger(cfg)

			if cfg.Target.Token == "" {
				logger.Error("Missing required hcloud.token")
				return fmt.Errorf("missing required hcloud.token")
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

	return app.Run(context.Background(), os.Args)
}

// RootFlags defines the available root flags.
func RootFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log.level",
			Value:       "info",
			Usage:       "Only log messages with given severity",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_LOG_LEVEL"),
			Destination: &cfg.Logs.Level,
		},
		&cli.BoolFlag{
			Name:        "log.pretty",
			Value:       false,
			Usage:       "Enable pretty messages for logging",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_LOG_PRETTY"),
			Destination: &cfg.Logs.Pretty,
		},
		&cli.StringFlag{
			Name:        "web.address",
			Value:       "0.0.0.0:9501",
			Usage:       "Address to bind the metrics server",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_WEB_ADDRESS"),
			Destination: &cfg.Server.Addr,
		},
		&cli.StringFlag{
			Name:        "web.path",
			Value:       "/metrics",
			Usage:       "Path to bind the metrics server",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_WEB_PATH"),
			Destination: &cfg.Server.Path,
		},
		&cli.BoolFlag{
			Name:        "web.debug",
			Value:       false,
			Usage:       "Enable pprof debugging for server",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_WEB_PPROF"),
			Destination: &cfg.Server.Pprof,
		},
		&cli.DurationFlag{
			Name:        "web.timeout",
			Value:       10 * time.Second,
			Usage:       "Server metrics endpoint timeout",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_WEB_TIMEOUT"),
			Destination: &cfg.Server.Timeout,
		},
		&cli.StringFlag{
			Name:        "web.config",
			Value:       "",
			Usage:       "Path to web-config file",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_WEB_CONFIG"),
			Destination: &cfg.Server.Web,
		},
		&cli.DurationFlag{
			Name:        "request.timeout",
			Value:       10 * time.Second,
			Usage:       "Request timeout as duration",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_REQUEST_TIMEOUT"),
			Destination: &cfg.Target.Timeout,
		},
		&cli.StringFlag{
			Name:        "hcloud.token",
			Value:       "",
			Usage:       "Access token for the Hetzner Cloud API",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_TOKEN"),
			Destination: &cfg.Target.Token,
		},
		&cli.BoolFlag{
			Name:        "collector.floating-ips",
			Value:       true,
			Usage:       "Enable collector for floating IPs",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_FLOATING_IPS"),
			Destination: &cfg.Collector.FloatingIPs,
		},
		&cli.BoolFlag{
			Name:        "collector.images",
			Value:       true,
			Usage:       "Enable collector for images",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_IMAGES"),
			Destination: &cfg.Collector.Images,
		},
		&cli.BoolFlag{
			Name:        "collector.pricing",
			Value:       true,
			Usage:       "Enable collector for pricing",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_PRICING"),
			Destination: &cfg.Collector.Pricing,
		},
		&cli.BoolFlag{
			Name:        "collector.servers",
			Value:       true,
			Usage:       "Enable collector for servers",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_SERVERS"),
			Destination: &cfg.Collector.Servers,
		},
		&cli.BoolFlag{
			Name:        "collector.server-metrics",
			Value:       false,
			Usage:       "Enable collector for server metrics",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_SERVER_METRICS"),
			Destination: &cfg.Collector.ServerMetrics,
		},
		&cli.BoolFlag{
			Name:        "collector.load-balancers",
			Value:       true,
			Usage:       "Enable collector for load balancers",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_LOAD_BALANCERS"),
			Destination: &cfg.Collector.LoadBalancers,
		},
		&cli.BoolFlag{
			Name:        "collector.ssh-keys",
			Value:       true,
			Usage:       "Enable collector for SSH keys",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_SSH_KEYS"),
			Destination: &cfg.Collector.SSHKeys,
		},
		&cli.BoolFlag{
			Name:        "collector.volumes",
			Value:       false,
			Usage:       "Enable collector for volumes",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_VOLUMES"),
			Destination: &cfg.Collector.Volumes,
		},
		&cli.BoolFlag{
			Name:        "collector.storageboxes",
			Value:       false,
			Usage:       "Enable collector for storage boxes",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_COLLECTOR_STORAGEBOXES"),
			Destination: &cfg.Collector.StorageBoxes,
		},
		&cli.StringSliceFlag{
			Name:        "collector.storageboxes.labels",
			Value:       config.StorageBoxLabels(),
			Usage:       "List of labels used for storageboxes",
			Sources:     cli.EnvVars("HCLOUD_EXPORTER_STORAGEBOXES_LABELS"),
			Destination: &cfg.Target.StorageBoxes.Labels,
		},
	}
}
