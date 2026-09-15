package exporter

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/hcloud_exporter/pkg/config"
)

// FloatingIPCollector collects metrics about the floating IPs.
type FloatingIPCollector struct {
	client   *hcloud.Client
	logger   *slog.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Active *prometheus.Desc
}

// NewFloatingIPCollector returns a new FloatingIPCollector.
func NewFloatingIPCollector(logger *slog.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *FloatingIPCollector {
	if failures != nil {
		failures.WithLabelValues("floating_ip").Add(0)
	}

	labels := cfg.FloatingIPs.Labels
	return &FloatingIPCollector{
		client:   client,
		logger:   logger.With("collector", "floating-ip"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Active: prometheus.NewDesc(
			"hcloud_floating_ip_active",
			"If 1 the floating IP is used by a server, 0 otherwise",
			labels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *FloatingIPCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Active,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *FloatingIPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Active
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *FloatingIPCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	ips, err := c.client.FloatingIP.All(ctx)

	if err != nil {
		c.logger.Error("Failed to fetch floating IPs",
			"err", err,
		)

		c.failures.WithLabelValues("floating_ip").Inc()
		return
	}

	c.logger.Debug("Fetched floating IPs",
		"count", len(ips),
	)

	var names map[int64]string

	for _, ip := range ips {
		if ip.Server != nil && ip.Server.Name == "" {
			if names, err = resolveServerNames(ctx, c.client); err != nil {
				c.logger.Warn("Failed to resolve server names",
					"err", err,
				)
			}

			break
		}
	}

	for _, ip := range ips {
		var (
			active float64
			name   string
		)

		if ip.Server != nil {
			active = 1.0
			name = ip.Server.Name

			if name == "" {
				name = names[ip.Server.ID]
			}
		}

		labels := []string{}

		for _, label := range c.config.FloatingIPs.Labels {
			labels = append(
				labels,
				c.byLabel(ip, name, label),
			)
		}

		ch <- prometheus.MustNewConstMetric(
			c.Active,
			prometheus.GaugeValue,
			active,
			labels...,
		)
	}

	c.logger.Debug("Processed floating IP collector",
		"duration", time.Since(now),
	)

	c.duration.WithLabelValues("floating_ip").Observe(time.Since(now).Seconds())
}

func (c *FloatingIPCollector) byLabel(record *hcloud.FloatingIP, server, label string) string {
	switch label {
	case "id":
		return strconv.FormatInt(record.ID, 10)
	case "server":
		return server
	case "location":
		return record.HomeLocation.Name
	case "type":
		return string(record.Type)
	case "ip":
		return record.IP.String()
	default:
		if val, ok := record.Labels[label]; ok {
			return val
		}

		return ""
	}
}
