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

// StorageBoxCollector collects metrics about the SSH keys.
type StorageBoxCollector struct {
	client   *hcloud.Client
	logger   *slog.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Active       *prometheus.Desc
	Locked       *prometheus.Desc
	GeneralSize  *prometheus.Desc
	DataSize     *prometheus.Desc
	SnapshotSize *prometheus.Desc
	Quota        *prometheus.Desc
	External     *prometheus.Desc
	Samba        *prometheus.Desc
	SSH          *prometheus.Desc
	Webdav       *prometheus.Desc
	ZFS          *prometheus.Desc
	Protection   *prometheus.Desc
	Created      *prometheus.Desc
}

// NewStorageBoxCollector returns a new StorageBoxCollector.
func NewStorageBoxCollector(logger *slog.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *StorageBoxCollector {
	if failures != nil {
		failures.WithLabelValues("storagebox").Add(0)
	}

	labels := cfg.StorageBoxes.Labels
	return &StorageBoxCollector{
		client:   client,
		logger:   logger.With("collector", "storagebox"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Active: prometheus.NewDesc(
			"hcloud_storagebox_active",
			"If 1 the storagebox have been active, 0 otherwise",
			labels,
			nil,
		),
		Locked: prometheus.NewDesc(
			"hcloud_storagebox_locked",
			"If 1 the storagebox have been locked, 0 otherwise",
			labels,
			nil,
		),
		GeneralSize: prometheus.NewDesc(
			"hcloud_storagebox_general_size",
			"Overall size for the storagebox in bytes",
			labels,
			nil,
		),
		DataSize: prometheus.NewDesc(
			"hcloud_storagebox_data_size",
			"Data size for the storagebox in bytes",
			labels,
			nil,
		),
		SnapshotSize: prometheus.NewDesc(
			"hcloud_storagebox_snapshot_size",
			"Snapshot size for the storagebox in bytes",
			labels,
			nil,
		),
		Quota: prometheus.NewDesc(
			"hcloud_storagebox_quota",
			"Available storage for the storagebox in bytes",
			labels,
			nil,
		),
		External: prometheus.NewDesc(
			"hcloud_storagebox_external",
			"If 1 the storagebox can be accessed from external, 0 otherwise",
			labels,
			nil,
		),
		Samba: prometheus.NewDesc(
			"hcloud_storagebox_samba",
			"If 1 the storagebox can be accessed via samba, 0 otherwise",
			labels,
			nil,
		),
		SSH: prometheus.NewDesc(
			"hcloud_storagebox_ssh",
			"If 1 the storagebox can be accessed via ssh, 0 otherwise",
			labels,
			nil,
		),
		Webdav: prometheus.NewDesc(
			"hcloud_storagebox_webdav",
			"If 1 the storagebox can be accessed via webdav, 0 otherwise",
			labels,
			nil,
		),
		ZFS: prometheus.NewDesc(
			"hcloud_storagebox_zfs",
			"If 1 the zfs directory is visible, 0 otherwise",
			labels,
			nil,
		),
		Protection: prometheus.NewDesc(
			"hcloud_storagebox_protection",
			"If 1 the storagebox is protected, 0 otherwise",
			labels,
			nil,
		),
		Created: prometheus.NewDesc(
			"hcloud_storagebox_created",
			"Timestamp when the storagebox have been created",
			labels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *StorageBoxCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Active,
		c.Locked,
		c.GeneralSize,
		c.DataSize,
		c.SnapshotSize,
		c.Quota,
		c.External,
		c.Samba,
		c.SSH,
		c.Webdav,
		c.ZFS,
		c.Protection,
		c.Created,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *StorageBoxCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Active
	ch <- c.Locked
	ch <- c.GeneralSize
	ch <- c.DataSize
	ch <- c.SnapshotSize
	ch <- c.Quota
	ch <- c.External
	ch <- c.Samba
	ch <- c.SSH
	ch <- c.Webdav
	ch <- c.ZFS
	ch <- c.Protection
	ch <- c.Created
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *StorageBoxCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	storageboxes, err := c.client.StorageBox.All(ctx)
	c.duration.WithLabelValues("storagebox").Observe(time.Since(now).Seconds())

	if err != nil {
		c.logger.Error("Failed to fetch storageboxes",
			"err", err,
		)

		c.failures.WithLabelValues("storagebox").Inc()
		return
	}

	c.logger.Debug("Fetched storageboxes",
		"count", len(storageboxes),
	)

	for _, record := range storageboxes {
		if record.Status == "initializing" {
			c.logger.Warn("Skipping initializing storagebox",
				"id", record.ID,
				"name", record.Name,
			)

			continue
		}

		labels := []string{}

		for _, label := range c.config.StorageBoxes.Labels {
			labels = append(
				labels,
				c.byLabel(record, label),
			)
		}

		ch <- prometheus.MustNewConstMetric(
			c.Active,
			prometheus.GaugeValue,
			boolToFloat64(record.Status == "active"),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Locked,
			prometheus.GaugeValue,
			boolToFloat64(record.Status == "locked"),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.GeneralSize,
			prometheus.GaugeValue,
			float64(record.Stats.Size),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DataSize,
			prometheus.GaugeValue,
			float64(record.Stats.SizeData),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.SnapshotSize,
			prometheus.GaugeValue,
			float64(record.Stats.SizeSnapshots),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Quota,
			prometheus.GaugeValue,
			float64(record.StorageBoxType.Size),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.External,
			prometheus.GaugeValue,
			boolToFloat64(record.AccessSettings.ReachableExternally),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Samba,
			prometheus.GaugeValue,
			boolToFloat64(record.AccessSettings.SambaEnabled),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.SSH,
			prometheus.GaugeValue,
			boolToFloat64(record.AccessSettings.SSHEnabled),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Webdav,
			prometheus.GaugeValue,
			boolToFloat64(record.AccessSettings.WebDAVEnabled),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ZFS,
			prometheus.GaugeValue,
			boolToFloat64(record.AccessSettings.ZFSEnabled),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Protection,
			prometheus.GaugeValue,
			boolToFloat64(record.Protection.Delete),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Created,
			prometheus.GaugeValue,
			float64(record.Created.Unix()),
			labels...,
		)
	}
}

func (c *StorageBoxCollector) byLabel(record *hcloud.StorageBox, label string) string {
	switch label {
	case "id":
		return strconv.FormatInt(record.ID, 10)
	case "name":
		return record.Name
	case "location":
		return record.Location.Name
	case "server":
		return record.Server
	case "system":
		return record.System
	case "type":
		return record.StorageBoxType.Name
	case "username":
		return record.Username
	default:
		if val, ok := record.Labels[label]; ok {
			return val
		}

		return ""
	}
}
