package exporter

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/hcloud_exporter/pkg/config"
)

// ImageCollector collects metrics about the images.
type ImageCollector struct {
	client   *hcloud.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Active     *prometheus.Desc
	ImageSize  *prometheus.Desc
	DiskSize   *prometheus.Desc
	Created    *prometheus.Desc
	Deprecated *prometheus.Desc
}

// NewImageCollector returns a new ImageCollector.
func NewImageCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *ImageCollector {
	if failures != nil {
		failures.WithLabelValues("image").Add(0)
	}

	labels := []string{"id", "name", "type", "server", "flavor", "version"}
	return &ImageCollector{
		client:   client,
		logger:   log.With(logger, "collector", "image"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Active: prometheus.NewDesc(
			"hcloud_image_active",
			"If 1 the image is used by a server, 0 otherwise",
			labels,
			nil,
		),
		ImageSize: prometheus.NewDesc(
			"hcloud_image_size_bytes",
			"Size of the image in bytes",
			labels,
			nil,
		),
		DiskSize: prometheus.NewDesc(
			"hcloud_image_disk_bytes",
			"Size if the disk for the image in bytes",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"hcloud_image_created_timestamp",
			"Timestamp when the image have been created",
			labels,
			nil,
		),
		Deprecated: prometheus.NewDesc(
			"hcloud_image_deprecated_timestamp",
			"Timestamp when the image will be deprecated, 0 if not deprecated",
			labels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *ImageCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Active,
		c.ImageSize,
		c.DiskSize,
		c.Created,
		c.Deprecated,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *ImageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Active
	ch <- c.ImageSize
	ch <- c.DiskSize
	ch <- c.Created
	ch <- c.Deprecated
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *ImageCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	images, err := c.client.Image.All(ctx)

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch images",
			"err", err,
		)

		c.failures.WithLabelValues("image").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched images",
		"count", len(images),
	)

	for _, image := range images {
		var (
			active     float64
			name       string
			deprecated float64
		)

		if image.CreatedFrom != nil {
			name = image.CreatedFrom.Name
		}

		if image.BoundTo != nil && image.BoundTo.Name != "" {
			active = 1.0
			name = image.BoundTo.Name
		}

		labels := []string{
			strconv.Itoa(image.ID),
			image.Name,
			string(image.Type),
			name,
			image.OSFlavor,
			image.OSVersion,
		}

		ch <- prometheus.MustNewConstMetric(
			c.Active,
			prometheus.GaugeValue,
			active,
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ImageSize,
			prometheus.GaugeValue,
			float64(image.ImageSize*1024*1024),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DiskSize,
			prometheus.GaugeValue,
			float64(image.DiskSize*1024*1024),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Created,
			prometheus.GaugeValue,
			float64(image.Created.Unix()),
			labels...,
		)

		if !image.Deprecated.IsZero() {
			deprecated = float64(image.Deprecated.Unix())
		}

		ch <- prometheus.MustNewConstMetric(
			c.Deprecated,
			prometheus.GaugeValue,
			deprecated,
			labels...,
		)
	}

	level.Debug(c.logger).Log(
		"msg", "Processed image collector",
		"duration", time.Since(now),
	)

	c.duration.WithLabelValues("image").Observe(time.Since(now).Seconds())
}
