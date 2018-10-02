package exporter

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// namespace defines the namespace for the exporter.
	namespace = "hcloud"
)

// GeneralCollector collects metrics, mostly runtime, about this exporter in general.
type GeneralCollector struct {
	version   string
	revision  string
	buildDate string
	goVersion string
	startTime time.Time

	StartTime *prometheus.Desc
	BuildInfo *prometheus.Desc
}

// NewGeneralCollector returns a new GeneralCollector.
func NewGeneralCollector(version string, revision string, buildDate string, goVersion string, startTime time.Time) *GeneralCollector {
	return &GeneralCollector{
		version:   version,
		revision:  revision,
		buildDate: buildDate,
		goVersion: goVersion,
		startTime: startTime,

		StartTime: prometheus.NewDesc(
			"hcloud_start_time",
			"Unix timestamp of the start time",
			nil,
			nil,
		),

		BuildInfo: prometheus.NewDesc(
			"hcloud_build_info",
			"A metric with a constant '1' value labeled by version, revision, and branch from which the exporter was built.",
			[]string{"version", "revision", "builddate", "goversion"},
			nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *GeneralCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.StartTime
	ch <- c.BuildInfo
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *GeneralCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.StartTime,
		prometheus.GaugeValue,
		float64(c.startTime.Unix()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.BuildInfo,
		prometheus.GaugeValue,
		1.0,
		c.version,
		c.revision,
		c.buildDate,
		c.goVersion,
	)
}
