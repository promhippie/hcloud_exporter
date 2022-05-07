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

// ServerMetricsCollector collects metrics about the servers.
type ServerMetricsCollector struct {
	client          *hcloud.Client
	logger          log.Logger
	failures        *prometheus.CounterVec
	duration        *prometheus.HistogramVec
	config          config.Target
	metricsDuration time.Duration

	CPU                  *prometheus.Desc
	Disk0IOPSRead        *prometheus.Desc
	Disk0IOPSWrite       *prometheus.Desc
	Disk0BandwidthRead   *prometheus.Desc
	Disk0BandwidthWrite  *prometheus.Desc
	Network0PPSIn        *prometheus.Desc
	Network0PPSOut       *prometheus.Desc
	Network0BandwidthIn  *prometheus.Desc
	Network0BandwidthOut *prometheus.Desc
}

// avgServerMetric returns the average of the given metric.
func avgServerMetric(seriesValues []hcloud.ServerMetricsValue) float64 {
	avg := 0.0

	for _, value := range seriesValues {
		floatValue, _ := strconv.ParseFloat(value.Value, 64)
		avg += floatValue
	}

	avg /= float64(len(seriesValues))
	return avg
}

// NewServerMetricsCollector returns a new ServerMetricsCollector.
func NewServerMetricsCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target, metricsDuration time.Duration) *ServerMetricsCollector {
	if failures != nil {
		failures.WithLabelValues("server_metrics").Add(0)
	}

	labels := []string{"id", "name", "datacenter"}

	return &ServerMetricsCollector{
		client:          client,
		logger:          log.With(logger, "collector", "server_metrics"),
		failures:        failures,
		duration:        duration,
		config:          cfg,
		metricsDuration: metricsDuration,

		CPU: prometheus.NewDesc(
			"hcloud_server_metrics_cpu_percent",
			"Percent CPU usage",
			labels,
			nil,
		),

		Disk0IOPSRead: prometheus.NewDesc(
			"hcloud_server_metrics_disk_0_iops_read",
			"Number of read IO operations per second",
			labels,
			nil,
		),

		Disk0IOPSWrite: prometheus.NewDesc(
			"hcloud_server_metrics_disk_0_iops_write",
			"Number of write IO operations per second",
			labels,
			nil,
		),

		Disk0BandwidthRead: prometheus.NewDesc(
			"hcloud_server_metrics_disk_0_bandwidth_read",
			"Bytes read per second",
			labels,
			nil,
		),

		Disk0BandwidthWrite: prometheus.NewDesc(
			"hcloud_server_metrics_disk_0_bandwidth_write",
			"Bytes written per second",
			labels,
			nil,
		),

		Network0PPSIn: prometheus.NewDesc(
			"hcloud_server_metrics_network_0_pps_in",
			"Public Network interface packets per second received",
			labels,
			nil,
		),

		Network0PPSOut: prometheus.NewDesc(
			"hcloud_server_metrics_network_0_pps_out",
			"Public Network interface packets per second sent",
			labels,
			nil,
		),

		Network0BandwidthIn: prometheus.NewDesc(
			"hcloud_server_metrics_network_0_bandwidth_in",
			"Public Network interface bytes/s received",
			labels,
			nil,
		),

		Network0BandwidthOut: prometheus.NewDesc(
			"hcloud_server_metrics_network_0_bandwidth_out",
			"Public Network interface bytes/s sent",
			labels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *ServerMetricsCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.CPU,
		c.Disk0IOPSRead,
		c.Disk0IOPSWrite,
		c.Disk0BandwidthRead,
		c.Disk0BandwidthWrite,
		c.Network0PPSIn,
		c.Network0PPSOut,
		c.Network0BandwidthIn,
		c.Network0BandwidthOut,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *ServerMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CPU
	ch <- c.Disk0IOPSRead
	ch <- c.Disk0IOPSWrite
	ch <- c.Disk0BandwidthRead
	ch <- c.Disk0BandwidthWrite
	ch <- c.Network0PPSIn
	ch <- c.Network0PPSOut
	ch <- c.Network0BandwidthIn
	ch <- c.Network0BandwidthOut
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ServerMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	servers, err := c.client.Server.All(ctx)
	c.duration.WithLabelValues("server_metrics").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch servers",
			"err", err,
		)

		c.failures.WithLabelValues("server_metrics").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched servers",
		"count", len(servers),
	)

	for _, server := range servers {
		labels := []string{
			strconv.Itoa(server.ID),
			server.Name,
			server.Datacenter.Name,
		}

		metricsOpts := hcloud.ServerGetMetricsOpts{
			Types: []hcloud.ServerMetricType{hcloud.ServerMetricCPU, hcloud.ServerMetricDisk, hcloud.ServerMetricNetwork},
			Start: time.Now().Add(-c.metricsDuration),
			End:   time.Now(),
		}

		now := time.Now()
		metrics, _, err := c.client.Server.GetMetrics(ctx, server, metricsOpts)

		c.duration.WithLabelValues("server_metrics").Observe(time.Since(now).Seconds())

		if err != nil {
			level.Error(c.logger).Log(
				"msg", "Failed to fetch server metrics",
				"server", server.Name,
				"err", err,
			)

			c.failures.WithLabelValues("server_metrics").Inc()
			return
		}

		// Series names are documented at https://docs.hetzner.cloud/#servers-get-metrics-for-a-server
		seriesNamesToDesc := map[string]*prometheus.Desc{
			"cpu":                     c.CPU,
			"disk.0.iops.read":        c.Disk0IOPSRead,
			"disk.0.iops.write":       c.Disk0IOPSWrite,
			"disk.0.bandwidth.read":   c.Disk0BandwidthRead,
			"disk.0.bandwidth.write":  c.Disk0BandwidthWrite,
			"network.0.pps.in":        c.Network0PPSIn,
			"network.0.pps.out":       c.Network0PPSOut,
			"network.0.bandwidth.in":  c.Network0BandwidthIn,
			"network.0.bandwidth.out": c.Network0BandwidthOut,
		}

		for seriesName, seriesValues := range metrics.TimeSeries {
			if desc, ok := seriesNamesToDesc[seriesName]; ok {
				ch <- prometheus.MustNewConstMetric(
					desc,
					prometheus.GaugeValue,
					avgServerMetric(seriesValues),
					labels...,
				)
			} else {
				level.Warn(c.logger).Log(
					"msg", "Unknown series name",
					"seriesName", seriesName,
				)
			}
		}
	}
}
