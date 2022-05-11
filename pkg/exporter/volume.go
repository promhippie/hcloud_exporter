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

// VolumeCollector collects metrics about the volumes.
type VolumeCollector struct {
	client   *hcloud.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Status     *prometheus.Desc
	Size       *prometheus.Desc
	Protection *prometheus.Desc
	Created    *prometheus.Desc
}

// NewVolumeCollector returns a new VolumeCollector.
func NewVolumeCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *VolumeCollector {
	if failures != nil {
		failures.WithLabelValues("volume").Add(0)
	}

	labels := []string{"id", "server", "location", "name"}
	return &VolumeCollector{
		client:   client,
		logger:   log.With(logger, "collector", "volume"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Status: prometheus.NewDesc(
			"hcloud_volume_status",
			"If 1 the volume is availabel, 0 otherwise",
			labels,
			nil,
		),
		Size: prometheus.NewDesc(
			"hcloud_volume_size",
			"Size of the volume in GB",
			labels,
			nil,
		),
		Protection: prometheus.NewDesc(
			"hcloud_volume_protection",
			"If 1 the volume is protected, 0 otherwise",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"hcloud_volume_created",
			"Timestamp when the volume have been created",
			labels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *VolumeCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Status,
		c.Size,
		c.Protection,
		c.Created,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *VolumeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Status
	ch <- c.Size
	ch <- c.Protection
	ch <- c.Created
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *VolumeCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	volumes, err := c.client.Volume.All(ctx)

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch volumes",
			"err", err,
		)

		c.failures.WithLabelValues("volume").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched volumes",
		"count", len(volumes),
	)

	for _, volume := range volumes {
		var (
			status     float64
			protection float64
			name       string
		)

		if volume.Server != nil {
			name = volume.Server.Name
		}

		labels := []string{
			strconv.Itoa(volume.ID),
			name,
			volume.Location.Name,
			volume.Name,
		}

		if volume.Status == "available" {
			status = 1.0
		}

		if volume.Protection.Delete {
			protection = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.Status,
			prometheus.GaugeValue,
			status,
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Size,
			prometheus.GaugeValue,
			float64(volume.Size),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Protection,
			prometheus.GaugeValue,
			protection,
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Created,
			prometheus.GaugeValue,
			float64(volume.Created.Unix()),
			labels...,
		)
	}

	c.duration.WithLabelValues("volume").Observe(time.Since(now).Seconds())
}
