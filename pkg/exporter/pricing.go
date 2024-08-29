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

// PricingCollector collects metrics about the prices for additional products: Image, FloatingIP, Traffic and ServerBackup.
type PricingCollector struct {
	client   *hcloud.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Image         *prometheus.Desc
	FloatingIP    *prometheus.Desc
	Traffic       *prometheus.Desc
	ServerBackup  *prometheus.Desc
	Volume        *prometheus.Desc
	Servers       *prometheus.Desc
	LoadBalancers *prometheus.Desc
}

// NewPricingCollector returns a new PricingCollector.
func NewPricingCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *PricingCollector {
	if failures != nil {
		failures.WithLabelValues("pricing").Add(0)
	}

	labels := []string{"currency", "vat"}
	return &PricingCollector{
		client:   client,
		logger:   log.With(logger, "collector", "pricing"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Image: prometheus.NewDesc(
			"hcloud_pricing_image",
			"The cost of an image per GB/month",
			labels,
			nil,
		),
		FloatingIP: prometheus.NewDesc(
			"hcloud_pricing_floating_ip",
			"The cost of one floating IP per month",
			labels,
			nil,
		),
		Traffic: prometheus.NewDesc(
			"hcloud_pricing_traffic",
			"The cost of additional traffic per TB",
			labels,
			nil,
		),
		ServerBackup: prometheus.NewDesc(
			"hcloud_pricing_server_backup",
			"Will increase base server costs by specific percentage if server backups are enabled",
			nil,
			nil,
		),
		Volume: prometheus.NewDesc(
			"hcloud_pricing_volume",
			"The cost of a volume per GB/month",
			labels,
			nil,
		),
		Servers: prometheus.NewDesc(
			"hcloud_pricing_server_type",
			"The costs of a server by type and location per month",
			append(labels, "type", "location"),
			nil,
		),
		LoadBalancers: prometheus.NewDesc(
			"hcloud_pricing_loadbalancer_type",
			"The costs of a load balancer by type and location per month",
			append(labels, "type", "location"),
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *PricingCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Image,
		c.FloatingIP,
		c.Traffic,
		c.ServerBackup,
		c.Volume,
		c.Servers,
		c.LoadBalancers,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *PricingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Image
	ch <- c.FloatingIP
	ch <- c.Traffic
	ch <- c.ServerBackup
	ch <- c.Volume
	ch <- c.Servers
	ch <- c.LoadBalancers
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *PricingCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	pricing, _, err := c.client.Pricing.Get(ctx)

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch pricing",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched pricing",
	)

	if gross, err := strconv.ParseFloat(pricing.Image.PerGBMonth.Gross, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse image costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.Image,
			prometheus.GaugeValue,
			gross,
			pricing.Image.PerGBMonth.Currency,
			"gross",
		)
	}

	if net, err := strconv.ParseFloat(pricing.Image.PerGBMonth.Net, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse image costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.Image,
			prometheus.GaugeValue,
			net,
			pricing.Image.PerGBMonth.Currency,
			"net",
		)
	}

	if gross, err := strconv.ParseFloat(pricing.FloatingIP.Monthly.Gross, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse floating IP costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.FloatingIP,
			prometheus.GaugeValue,
			gross,
			pricing.FloatingIP.Monthly.Currency,
			"gross",
		)
	}

	if net, err := strconv.ParseFloat(pricing.FloatingIP.Monthly.Net, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse floating IP costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.FloatingIP,
			prometheus.GaugeValue,
			net,
			pricing.FloatingIP.Monthly.Currency,
			"net",
		)
	}

	if backup, err := strconv.ParseFloat(pricing.ServerBackup.Percentage, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse server backup costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.ServerBackup,
			prometheus.GaugeValue,
			backup/100,
		)
	}

	if gross, err := strconv.ParseFloat(pricing.Volume.PerGBMonthly.Gross, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse volume costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.Volume,
			prometheus.GaugeValue,
			gross,
			pricing.Volume.PerGBMonthly.Currency,
			"gross",
		)
	}

	if net, err := strconv.ParseFloat(pricing.Volume.PerGBMonthly.Net, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse volume costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.Volume,
			prometheus.GaugeValue,
			net,
			pricing.Volume.PerGBMonthly.Currency,
			"net",
		)
	}

	if gross, err := strconv.ParseFloat(pricing.Traffic.PerTB.Gross, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse traffic costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.Traffic,
			prometheus.GaugeValue,
			gross,
			pricing.Traffic.PerTB.Currency,
			"gross",
		)
	}

	if net, err := strconv.ParseFloat(pricing.Traffic.PerTB.Net, 64); err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to parse traffic costs",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
	} else {
		ch <- prometheus.MustNewConstMetric(
			c.Traffic,
			prometheus.GaugeValue,
			net,
			pricing.Traffic.PerTB.Currency,
			"net",
		)
	}

	for _, serverType := range pricing.ServerTypes {
		for _, serverPricing := range serverType.Pricings {
			if gross, err := strconv.ParseFloat(serverPricing.Monthly.Gross, 64); err != nil {
				level.Error(c.logger).Log(
					"msg", "Failed to parse server costs",
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.Servers,
					prometheus.GaugeValue,
					gross,
					serverPricing.Monthly.Currency,
					"gross",
					serverType.ServerType.Name,
					serverPricing.Location.Name,
				)
			}

			if net, err := strconv.ParseFloat(serverPricing.Monthly.Net, 64); err != nil {
				level.Error(c.logger).Log(
					"msg", "Failed to parse server costs",
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.Servers,
					prometheus.GaugeValue,
					net,
					serverPricing.Monthly.Currency,
					"net",
					serverType.ServerType.Name,
					serverPricing.Location.Name,
				)
			}
		}
	}

	for _, lbType := range pricing.LoadBalancerTypes {
		for _, lbPricing := range lbType.Pricings {
			if gross, err := strconv.ParseFloat(lbPricing.Monthly.Gross, 64); err != nil {
				level.Error(c.logger).Log(
					"msg", "Failed to parse server costs",
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.LoadBalancers,
					prometheus.GaugeValue,
					gross,
					lbPricing.Monthly.Currency,
					"gross",
					lbType.LoadBalancerType.Name,
					lbPricing.Location.Name,
				)
			}

			if net, err := strconv.ParseFloat(lbPricing.Monthly.Net, 64); err != nil {
				level.Error(c.logger).Log(
					"msg", "Failed to parse server costs",
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.LoadBalancers,
					prometheus.GaugeValue,
					net,
					lbPricing.Monthly.Currency,
					"net",
					lbType.LoadBalancerType.Name,
					lbPricing.Location.Name,
				)
			}
		}
	}

	level.Debug(c.logger).Log(
		"msg", "Processed pricing collector",
		"duration", time.Since(now),
	)

	c.duration.WithLabelValues("pricing").Observe(time.Since(now).Seconds())
}
