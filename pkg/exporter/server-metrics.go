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

// ServerMetricsCollector collects the servers metrics.
type ServerMetricsCollector struct {
	client   *hcloud.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	CPU           *prometheus.Desc
	DiskReadIops  *prometheus.Desc
	DiskWriteIops *prometheus.Desc
	DiskReadBps   *prometheus.Desc
	DiskWriteBps  *prometheus.Desc
	NetworkInPps  *prometheus.Desc
	NetworkOutPps *prometheus.Desc
	NetworkInBps  *prometheus.Desc
	NetworkOutBps *prometheus.Desc
}

// NewServerMetricsCollector returns a new ServerMetricsCollector.
func NewServerMetricsCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *ServerMetricsCollector {
	if failures != nil {
		failures.WithLabelValues("server-metrics").Add(0)
	}

	labels := []string{"id", "name", "datacenter"}
	diskLabels := []string{"id", "name", "datacenter", "disk"}
	networkLabels := []string{"id", "name", "datacenter", "interface"}
	return &ServerMetricsCollector{
		client:   client,
		logger:   log.With(logger, "collector", "server-metrics"),
		failures: failures,
		duration: duration,
		config:   cfg,

		CPU: prometheus.NewDesc(
			"hcloud_server_metrics_cpu",
			"Server CPU usage metric",
			labels,
			nil,
		),

		DiskReadIops: prometheus.NewDesc(
			"hcloud_server_metrics_disk_read_iops",
			"Server disk read iop/s metric",
			labels,
			nil,
		),

		DiskWriteIops: prometheus.NewDesc(
			"hcloud_server_metrics_disk_write_iops",
			"Server disk write iop/s metric",
			diskLabels,
			nil,
		),

		DiskWriteBps: prometheus.NewDesc(
			"hcloud_server_metrics_disk_write_bps",
			"Server disk write bytes/s metric",
			diskLabels,
			nil,
		),

		NetworkInPps: prometheus.NewDesc(
			"hcloud_server_metrics_network_in_pps",
			"Server network incoming packets/s metric",
			networkLabels,
			nil,
		),

		NetworkOutPps: prometheus.NewDesc(
			"hcloud_server_metrics_network_out_pps",
			"Server network outgoing packets/s metric",
			networkLabels,
			nil,
		),

		NetworkInBps: prometheus.NewDesc(
			"hcloud_server_metrics_network_in_bps",
			"Server network incoming bytes/s metric",
			networkLabels,
			nil,
		),

		NetworkOutBps: prometheus.NewDesc(
			"hcloud_server_metrics_network_out_bps",
			"Server network outgoing bytes/s metric",
			networkLabels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *ServerMetricsCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{

		c.CPU,
		c.DiskReadIops,
		c.DiskWriteIops,
		c.DiskReadBps,
		c.DiskWriteBps,
		c.NetworkInPps,
		c.NetworkOutPps,
		c.NetworkInBps,
		c.NetworkOutBps,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *ServerMetricsCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- c.CPU
	ch <- c.DiskReadIops
	ch <- c.DiskWriteIops
	ch <- c.DiskReadBps
	ch <- c.DiskWriteBps
	ch <- c.NetworkInPps
	ch <- c.NetworkOutPps
	ch <- c.NetworkInBps
	ch <- c.NetworkOutBps
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ServerMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	opts := hcloud.ServerListOpts{
		Status: []hcloud.ServerStatus{
			hcloud.ServerStatusRunning,
		},
	}
	servers, _, err := c.client.Server.List(ctx, opts)

	level.Debug(c.logger).Log(
		"msg", "Fetched online servers",
		"count", len(servers),
	)

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch servers",
			"err", err,
		)

		c.failures.WithLabelValues("server-metrics").Inc()
		return
	}

	type empty struct{}
	sem := make(chan empty, len(servers))

	for _, server := range servers {

		labels := []string{
			strconv.Itoa(server.ID),
			server.Name,
			server.Datacenter.Name,
		}

		go func(c *ServerMetricsCollector, ctx context.Context, server *hcloud.Server) {

			metricsOpts := hcloud.ServerGetMetricsOpts{
				Types: []hcloud.ServerMetricType{
					hcloud.ServerMetricCPU,
					hcloud.ServerMetricDisk,
					hcloud.ServerMetricNetwork,
				},
				Start: time.Now(),
				End:   time.Now(),
			}

			metrics, _, err := c.client.Server.GetMetrics(ctx, server, metricsOpts)

			sem <- empty{}

			if err != nil {
				level.Error(c.logger).Log(
					"msg", "Failed to fetch server metrics",
					"err", err,
				)

				c.failures.WithLabelValues("server-metrics").Inc()
				return
			}

			// Hetzner currently only provides a single 0-indexed timeseries for each metric, so it's simply hardcoded.
			// If Hetzner ever extends this, determining the amount of returned timeseries would be better.
			diskLabels := []string{
				strconv.Itoa(server.ID),
				server.Name,
				server.Datacenter.Name,
				"0",
			}

			networkLabels := []string{
				strconv.Itoa(server.ID),
				server.Name,
				server.Datacenter.Name,
				"0",
			}

			CPU, _ := strconv.ParseFloat(metrics.TimeSeries["cpu"][0].Value, 64)
			DiskReadIops, _ := strconv.ParseFloat(metrics.TimeSeries["disk.0.iops.read"][0].Value, 64)
			DiskWriteIops, _ := strconv.ParseFloat(metrics.TimeSeries["disk.0.iops.write"][0].Value, 64)
			DiskReadBps, _ := strconv.ParseFloat(metrics.TimeSeries["disk.0.bandwidth.read"][0].Value, 64)
			DiskWriteBps, _ := strconv.ParseFloat(metrics.TimeSeries["disk.0.bandwidth.write"][0].Value, 64)
			NetworkInPps, _ := strconv.ParseFloat(metrics.TimeSeries["network.0.pps.in"][0].Value, 64)
			NetworkOutPps, _ := strconv.ParseFloat(metrics.TimeSeries["network.0.pps.out"][0].Value, 64)
			NetworkInBps, _ := strconv.ParseFloat(metrics.TimeSeries["network.0.bandwidth.in"][0].Value, 64)
			NetworkOutBps, _ := strconv.ParseFloat(metrics.TimeSeries["network.0.bandwidth.out"][0].Value, 64)

			ch <- prometheus.MustNewConstMetric(
				c.CPU,
				prometheus.GaugeValue,
				CPU,
				labels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.DiskReadIops,
				prometheus.GaugeValue,
				DiskReadIops,
				diskLabels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.DiskWriteIops,
				prometheus.GaugeValue,
				DiskWriteIops,
				diskLabels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.DiskReadBps,
				prometheus.GaugeValue,
				DiskReadBps,
				diskLabels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.DiskWriteBps,
				prometheus.GaugeValue,
				DiskWriteBps,
				diskLabels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.NetworkInPps,
				prometheus.GaugeValue,
				NetworkInPps,
				networkLabels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.NetworkOutPps,
				prometheus.GaugeValue,
				NetworkOutPps,
				networkLabels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.NetworkInBps,
				prometheus.GaugeValue,
				NetworkInBps,
				networkLabels...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.NetworkOutBps,
				prometheus.GaugeValue,
				NetworkOutBps,
				networkLabels...,
			)
		}(c, ctx, server)

	}

	// Wait for all go-routines to signal finished metrics fetch
	for i := 0; i < len(servers); i++ {
		<-sem
	}

	c.duration.WithLabelValues("server-metrics").Observe(time.Since(now).Seconds())
	level.Debug(c.logger).Log(
		"msg", "Fetched server metrics",
		"count", len(servers),
	)

}
