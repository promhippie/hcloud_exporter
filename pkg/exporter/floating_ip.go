package exporter

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/hcloud_exporter/pkg/config"
)

// FloatingIPCollector collects metrics about the floating IPs.
type FloatingIPCollector struct {
	client   *hcloud.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Active *prometheus.Desc
}

// NewFloatingIPCollector returns a new FloatingIPCollector.
func NewFloatingIPCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *FloatingIPCollector {
	if failures != nil {
		failures.WithLabelValues("floating_ip").Add(0)
	}

	labels := []string{"id", "server", "location", "type", "ip"}
	return &FloatingIPCollector{
		client:   client,
		logger:   log.With(logger, "collector", "floating_ip"),
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
		level.Error(c.logger).Log(
			"msg", "Failed to fetch floating IPs",
			"err", err,
		)

		c.failures.WithLabelValues("floating_ip").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched floating IPs",
		"count", len(ips),
	)

	for _, ip := range ips {
		var (
			active float64
			name   string
		)

		if ip.Server != nil {
			active = 1.0
			name = ip.Server.Name
		}

		labels := []string{
			strconv.Itoa(ip.ID),
			name,
			ip.HomeLocation.Name,
			string(ip.Type),
			ip.IP.String(),
		}

		ch <- prometheus.MustNewConstMetric(
			c.Active,
			prometheus.GaugeValue,
			active,
			labels...,
		)
	}

	c.duration.WithLabelValues("floating_ip").Observe(time.Since(now).Seconds())
}
