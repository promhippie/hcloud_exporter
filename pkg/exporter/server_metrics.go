package exporter

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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
		failures.WithLabelValues("server_metrics").Add(0)
	}

	labels := []string{"id", "name", "datacenter"}
	diskLabels := append(labels, "disk")
	networkLabels := append(labels, "interface")
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
			diskLabels,
			nil,
		),

		DiskWriteIops: prometheus.NewDesc(
			"hcloud_server_metrics_disk_write_iops",
			"Server disk write iop/s metric",
			diskLabels,
			nil,
		),

		DiskReadBps: prometheus.NewDesc(
			"hcloud_server_metrics_disk_read_bps",
			"Server disk write bytes/s metric",
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
	servers, _, err := c.client.Server.List(ctx, hcloud.ServerListOpts{
		Status: []hcloud.ServerStatus{
			hcloud.ServerStatusRunning,
		},
	})

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch servers",
			"err", err,
		)

		c.failures.WithLabelValues("server_metrics").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched online servers",
		"count", len(servers),
	)

	var (
		wg sync.WaitGroup
	)

	for _, server := range servers {
		labels := []string{
			strconv.FormatInt(server.ID, 10),
			server.Name,
			server.Datacenter.Name,
		}

		wg.Add(1)

		go func(server *hcloud.Server) {
			defer wg.Done()
			logger := log.With(c.logger, "server", server.Name)

			metrics, _, err := c.client.Server.GetMetrics(
				ctx,
				server,
				hcloud.ServerGetMetricsOpts{
					Types: []hcloud.ServerMetricType{
						hcloud.ServerMetricCPU,
						hcloud.ServerMetricDisk,
						hcloud.ServerMetricNetwork,
					},
					Start: now,
					End:   now,
				},
			)

			if err != nil {
				level.Error(logger).Log(
					"msg", "Failed to fetch server metrics",
					"err", err,
				)

				c.failures.WithLabelValues("server_metrics").Inc()
				return
			}

			diskLabels := append(labels, "0")
			networkLabels := append(labels, "0")

			if len(metrics.TimeSeries["cpu"]) > 0 {
				cpuUsage, _ := strconv.ParseFloat(metrics.TimeSeries["cpu"][len(metrics.TimeSeries["cpu"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.CPU,
					prometheus.GaugeValue,
					cpuUsage,
					labels...,
				)
			}

			if len(metrics.TimeSeries["disk.0.iops.read"]) > 0 {
				diskReadIops, _ := strconv.ParseFloat(metrics.TimeSeries["disk.0.iops.read"][len(metrics.TimeSeries["disk.0.iops.read"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.DiskReadIops,
					prometheus.GaugeValue,
					diskReadIops,
					diskLabels...,
				)
			}

			if len(metrics.TimeSeries["disk.0.iops.write"]) > 0 {
				diskWriteIops, _ := strconv.ParseFloat(metrics.TimeSeries["disk.0.iops.write"][len(metrics.TimeSeries["disk.0.iops.write"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.DiskWriteIops,
					prometheus.GaugeValue,
					diskWriteIops,
					diskLabels...,
				)
			}

			if len(metrics.TimeSeries["disk.0.bandwidth.read"]) > 0 {
				diskReadBps, _ := strconv.ParseFloat(metrics.TimeSeries["disk.0.bandwidth.read"][len(metrics.TimeSeries["disk.0.bandwidth.read"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.DiskReadBps,
					prometheus.GaugeValue,
					diskReadBps,
					diskLabels...,
				)
			}

			if len(metrics.TimeSeries["disk.0.bandwidth.write"]) > 0 {
				diskWriteBps, _ := strconv.ParseFloat(metrics.TimeSeries["disk.0.bandwidth.write"][len(metrics.TimeSeries["disk.0.bandwidth.write"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.DiskWriteBps,
					prometheus.GaugeValue,
					diskWriteBps,
					diskLabels...,
				)
			}

			if len(metrics.TimeSeries["network.0.pps.in"]) > 0 {
				networkInPps, _ := strconv.ParseFloat(metrics.TimeSeries["network.0.pps.in"][len(metrics.TimeSeries["network.0.pps.in"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.NetworkInPps,
					prometheus.GaugeValue,
					networkInPps,
					networkLabels...,
				)
			}

			if len(metrics.TimeSeries["network.0.pps.out"]) > 0 {
				networkOutPps, _ := strconv.ParseFloat(metrics.TimeSeries["network.0.pps.out"][len(metrics.TimeSeries["network.0.pps.out"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.NetworkOutPps,
					prometheus.GaugeValue,
					networkOutPps,
					networkLabels...,
				)
			}

			if len(metrics.TimeSeries["network.0.bandwidth.in"]) > 0 {
				networkInBps, _ := strconv.ParseFloat(metrics.TimeSeries["network.0.bandwidth.in"][len(metrics.TimeSeries["network.0.bandwidth.in"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.NetworkInBps,
					prometheus.GaugeValue,
					networkInBps,
					networkLabels...,
				)
			}

			if len(metrics.TimeSeries["network.0.bandwidth.out"]) > 0 {
				networkOutBps, _ := strconv.ParseFloat(metrics.TimeSeries["network.0.bandwidth.out"][len(metrics.TimeSeries["network.0.bandwidth.out"])-1].Value, 64)
				ch <- prometheus.MustNewConstMetric(
					c.NetworkOutBps,
					prometheus.GaugeValue,
					networkOutBps,
					networkLabels...,
				)
			}
		}(server)
	}

	wg.Wait()

	level.Debug(c.logger).Log(
		"msg", "Processed server metrics collector",
		"duration", time.Since(now),
	)

	c.duration.WithLabelValues("server_metrics").Observe(time.Since(now).Seconds())
}
