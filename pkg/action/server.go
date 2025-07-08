package action

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/promhippie/hcloud_exporter/pkg/exporter"
	"github.com/promhippie/hcloud_exporter/pkg/middleware"
	"github.com/promhippie/hcloud_exporter/pkg/version"
)

// Server handles the server sub-command.
func Server(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("Launching Hetzner Cloud Exporter",
		"version", version.String,
		"revision", version.Revision,
		"date", version.Date,
		"go", version.Go,
	)

	token, err := config.Value(cfg.Target.Token)

	if err != nil {
		logger.Error("Failed to load token from file",
			"err", err,
		)

		return err
	}

	client := hcloud.NewClient(
		hcloud.WithToken(
			token,
		),
		hcloud.WithApplication(
			"hcloud_exporter",
			version.String,
		),
	)

	var gr run.Group

	{
		server := &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      handler(cfg, logger, client),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: cfg.Server.Timeout,
		}

		gr.Add(func() error {
			logger.Info("Starting metrics server",
				"address", cfg.Server.Addr,
			)

			return web.ListenAndServe(
				server,
				&web.FlagConfig{
					WebListenAddresses: sliceP([]string{cfg.Server.Addr}),
					WebSystemdSocket:   boolP(false),
					WebConfigFile:      stringP(cfg.Server.Web),
				},
				logger,
			)
		}, func(reason error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Failed to shutdown metrics gracefully",
					"err", err,
				)

				return
			}

			logger.Info("Metrics shutdown gracefully",
				"reason", reason,
			)
		})
	}

	{
		stop := make(chan os.Signal, 1)

		gr.Add(func() error {
			signal.Notify(stop, os.Interrupt)

			<-stop

			return nil
		}, func(_ error) {
			close(stop)
		})
	}

	return gr.Run()
}

func handler(cfg *config.Config, logger *slog.Logger, client *hcloud.Client) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer(logger))
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Timeout)
	mux.Use(middleware.Cache)

	if cfg.Server.Pprof {
		mux.Mount("/debug", middleware.Profiler())
	}

	if cfg.Collector.FloatingIPs {
		logger.Debug("Floating IP collector registered")

		registry.MustRegister(exporter.NewFloatingIPCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.Images {
		logger.Debug("Image collector registered")

		registry.MustRegister(exporter.NewImageCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.Pricing {
		logger.Debug("Pricing collector registered")

		registry.MustRegister(exporter.NewPricingCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.Servers {
		logger.Debug("Server collector registered")

		registry.MustRegister(exporter.NewServerCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.ServerMetrics {
		logger.Debug("Server metrics collector registered")

		registry.MustRegister(exporter.NewServerMetricsCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.LoadBalancers {
		logger.Debug("Load balancer collector registered")

		registry.MustRegister(exporter.NewLoadBalancerCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.SSHKeys {
		logger.Debug("SSH key collector registered")

		registry.MustRegister(exporter.NewSSHKeyCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	if cfg.Collector.Volumes {
		logger.Debug("Volumes collector registered")

		registry.MustRegister(exporter.NewVolumeCollector(
			logger,
			client,
			requestFailures,
			requestDuration,
			cfg.Target,
		))
	}

	reg := promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			ErrorLog: promLogger{logger},
		},
	)

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, cfg.Server.Path, http.StatusMovedPermanently)
	})

	mux.Route("/", func(root chi.Router) {
		root.Get(cfg.Server.Path, func(w http.ResponseWriter, r *http.Request) {
			reg.ServeHTTP(w, r)
		})

		root.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
		})

		root.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
		})
	})

	return mux
}
