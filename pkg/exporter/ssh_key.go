package exporter

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/prometheus/client_golang/prometheus"
)

// SSHKeyCollector collects metrics about the SSH keys.
type SSHKeyCollector struct {
	client   *hcloud.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	timeout  time.Duration

	Key *prometheus.Desc
}

// NewSSHKeyCollector returns a new SSHKeyCollector.
func NewSSHKeyCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, timeout time.Duration) *SSHKeyCollector {
	failures.WithLabelValues("ssh_key").Add(0)

	labels := []string{"id", "name", "fingerprint"}
	return &SSHKeyCollector{
		client:   client,
		logger:   log.With(logger, "collector", "ssh_key"),
		failures: failures,
		duration: duration,
		timeout:  timeout,

		Key: prometheus.NewDesc(
			"hcloud_ssh_key",
			"Information about SSH keys in your HetznerCloud project",
			labels,
			nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *SSHKeyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Key
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *SSHKeyCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	now := time.Now()
	keys, err := c.client.SSHKey.All(ctx)
	c.duration.WithLabelValues("ssh_key").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch SSH keys",
			"err", err,
		)

		c.failures.WithLabelValues("ssh_key").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched SSH keys",
		"count", len(keys),
	)

	for _, key := range keys {
		labels := []string{
			strconv.Itoa(key.ID),
			key.Name,
			key.Fingerprint,
		}

		ch <- prometheus.MustNewConstMetric(
			c.Key,
			prometheus.GaugeValue,
			1.0,
			labels...,
		)
	}
}
