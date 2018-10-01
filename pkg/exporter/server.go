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

// ServerCollector collects metrics about the servers.
type ServerCollector struct {
	client   *hcloud.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	timeout  time.Duration

	Up              *prometheus.Desc
	Created         *prometheus.Desc
	IncludedTraffic *prometheus.Desc
	OutgoingTraffic *prometheus.Desc
	IngoingTraffic  *prometheus.Desc
	Cores           *prometheus.Desc
	Memory          *prometheus.Desc
	Disk            *prometheus.Desc
	PriceHourly     *prometheus.Desc
	PriceMonthly    *prometheus.Desc
}

// NewServerCollector returns a new ServerCollector.
func NewServerCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, timeout time.Duration) *ServerCollector {
	failures.WithLabelValues("server").Add(0)

	labels := []string{"id", "name", "datacenter"}
	return &ServerCollector{
		client:   client,
		logger:   logger,
		failures: failures,
		duration: duration,
		timeout:  timeout,

		Up: prometheus.NewDesc(
			"hcloud_server_running",
			"If 1 the server is running, 0 otherwise",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"hcloud_server_created_timestamp",
			"Timestamp when the server have been created",
			labels,
			nil,
		),
		IncludedTraffic: prometheus.NewDesc(
			"hcloud_server_included_traffic_bytes",
			"Included traffic for the server in bytes",
			labels,
			nil,
		),
		OutgoingTraffic: prometheus.NewDesc(
			"hcloud_server_outgoing_traffic_bytes",
			"Outgoing traffic from the server in bytes",
			labels,
			nil,
		),
		IngoingTraffic: prometheus.NewDesc(
			"hcloud_server_incming_traffic_bytes",
			"Ingoing traffic to the server in bytes",
			labels,
			nil,
		),
		Cores: prometheus.NewDesc(
			"hcloud_server_cores",
			"Server number of cores",
			labels,
			nil,
		),
		Memory: prometheus.NewDesc(
			"hcloud_server_memory_bytes",
			"Server memory in bytes",
			labels,
			nil,
		),
		Disk: prometheus.NewDesc(
			"hcloud_server_disk_bytes",
			"Server disk in bytes",
			labels,
			nil,
		),
		PriceHourly: prometheus.NewDesc(
			"hcloud_server_price_hourly",
			"Price of the server billed hourly in €",
			labels,
			nil,
		),
		PriceMonthly: prometheus.NewDesc(
			"hcloud_server_price_monthly",
			"Price of the server billed monthly in €",
			labels,
			nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
	ch <- c.Created
	ch <- c.IncludedTraffic
	ch <- c.OutgoingTraffic
	ch <- c.IngoingTraffic
	ch <- c.Cores
	ch <- c.Memory
	ch <- c.Disk
	ch <- c.PriceHourly
	ch <- c.PriceMonthly
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ServerCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	now := time.Now()
	servers, err := c.client.Server.All(ctx)
	c.duration.WithLabelValues("server").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch servers",
			"err", err,
		)

		c.failures.WithLabelValues("server").Inc()
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

		var up float64

		if server.Status == "running" {
			up = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			up,
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Created,
			prometheus.GaugeValue,
			float64(server.Created.Unix()),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IncludedTraffic,
			prometheus.GaugeValue,
			float64(server.IncludedTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.OutgoingTraffic,
			prometheus.GaugeValue,
			float64(server.OutgoingTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IngoingTraffic,
			prometheus.GaugeValue,
			float64(server.IngoingTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Cores,
			prometheus.GaugeValue,
			float64(server.ServerType.Cores),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Memory,
			prometheus.GaugeValue,
			float64(server.ServerType.Memory*1024*1024),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Disk,
			prometheus.GaugeValue,
			float64(server.ServerType.Disk*1024*1024),
			labels...,
		)

		for _, pricing := range server.ServerType.Pricings {
			if server.Datacenter.Location.Name == pricing.Location.Name {
				hourly, _ := strconv.ParseFloat(pricing.Hourly.Net, 64)

				ch <- prometheus.MustNewConstMetric(
					c.PriceHourly,
					prometheus.GaugeValue,
					hourly,
					labels...,
				)

				monthly, _ := strconv.ParseFloat(pricing.Monthly.Net, 64)

				ch <- prometheus.MustNewConstMetric(
					c.PriceMonthly,
					prometheus.GaugeValue,
					monthly,
					labels...,
				)
			}
		}
	}
}
