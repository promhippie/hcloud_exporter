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

// PricingCollector collects metrics about the prices for additional products: Image, FloatingIP, Traffic and ServerBackup.
type PricingCollector struct {
	client   *hcloud.Client
	logger   *slog.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Image                *prometheus.Desc
	Traffic              *prometheus.Desc
	ServerBackup         *prometheus.Desc
	Volume               *prometheus.Desc
	FloatingIP           *prometheus.Desc
	PrimaryIP            *prometheus.Desc
	Servers              *prometheus.Desc
	ServersTraffic       *prometheus.Desc
	LoadBalancers        *prometheus.Desc
	LoadBalancersTraffic *prometheus.Desc
}

// NewPricingCollector returns a new PricingCollector.
func NewPricingCollector(logger *slog.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *PricingCollector {
	if failures != nil {
		failures.WithLabelValues("pricing").Add(0)
	}

	labels := []string{"currency", "vat"}
	return &PricingCollector{
		client:   client,
		logger:   logger.With("collector", "pricing"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Image: prometheus.NewDesc(
			"hcloud_pricing_image",
			"The cost of an image per GB/month",
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
		FloatingIP: prometheus.NewDesc(
			"hcloud_pricing_floating_ip",
			"The cost of one floating IP per month",
			append(labels, "type", "location"),
			nil,
		),
		PrimaryIP: prometheus.NewDesc(
			"hcloud_pricing_primary_ip",
			"The cost of one primary IP per month",
			append(labels, "type", "location"),
			nil,
		),
		Servers: prometheus.NewDesc(
			"hcloud_pricing_server_type",
			"The costs of a server by type and location per month",
			append(labels, "type", "location"),
			nil,
		),
		ServersTraffic: prometheus.NewDesc(
			"hcloud_pricing_server_type_traffic",
			"The costs of additional traffic per TB for a server by type and location per month",
			append(labels, "type", "location"),
			nil,
		),
		LoadBalancers: prometheus.NewDesc(
			"hcloud_pricing_loadbalancer_type",
			"The costs of a load balancer by type and location per month",
			append(labels, "type", "location"),
			nil,
		),
		LoadBalancersTraffic: prometheus.NewDesc(
			"hcloud_pricing_loadbalancer_type_traffic",
			"The costs of additional traffic per TB for a load balancer by type and location per month",
			append(labels, "type", "location"),
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *PricingCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Image,
		c.ServerBackup,
		c.Volume,
		c.FloatingIP,
		c.PrimaryIP,
		c.Servers,
		c.ServersTraffic,
		c.LoadBalancers,
		c.LoadBalancersTraffic,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *PricingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Image
	ch <- c.ServerBackup
	ch <- c.Volume
	ch <- c.FloatingIP
	ch <- c.PrimaryIP
	ch <- c.Servers
	ch <- c.ServersTraffic
	ch <- c.LoadBalancers
	ch <- c.LoadBalancersTraffic
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *PricingCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	pricing, _, err := c.client.Pricing.Get(ctx)

	if err != nil {
		c.logger.Error("Failed to fetch pricing",
			"err", err,
		)

		c.failures.WithLabelValues("pricing").Inc()
		return
	}

	c.logger.Debug("Fetched pricing")

	if gross, err := strconv.ParseFloat(pricing.Image.PerGBMonth.Gross, 64); err != nil {
		c.logger.Error("Failed to parse image costs",
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
		c.logger.Error("Failed to parse image costs",
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

	if backup, err := strconv.ParseFloat(pricing.ServerBackup.Percentage, 64); err != nil {
		c.logger.Error("Failed to parse server backup costs",
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
		c.logger.Error("Failed to parse volume costs",
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
		c.logger.Error("Failed to parse volume costs",
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

	for _, floatingType := range pricing.FloatingIPs {
		for _, floatingPricing := range floatingType.Pricings {
			if gross, err := strconv.ParseFloat(floatingPricing.Monthly.Gross, 64); err != nil {
				c.logger.Error("Failed to parse floating IP costs",
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.FloatingIP,
					prometheus.GaugeValue,
					gross,
					floatingPricing.Monthly.Currency,
					"gross",
					string(floatingType.Type),
					floatingPricing.Location.Name,
				)
			}

			if net, err := strconv.ParseFloat(floatingPricing.Monthly.Net, 64); err != nil {
				c.logger.Error("Failed to parse floating IP costs",
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.FloatingIP,
					prometheus.GaugeValue,
					net,
					floatingPricing.Monthly.Currency,
					"net",
					string(floatingType.Type),
					floatingPricing.Location.Name,
				)
			}

		}
	}

	for _, primaryType := range pricing.PrimaryIPs {
		for _, primaryPricing := range primaryType.Pricings {
			if gross, err := strconv.ParseFloat(primaryPricing.Monthly.Gross, 64); err != nil {
				c.logger.Error("Failed to parse primary IP costs",
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.PrimaryIP,
					prometheus.GaugeValue,
					gross,
					"",
					"gross",
					string(primaryType.Type),
					primaryPricing.Location,
				)
			}

			if net, err := strconv.ParseFloat(primaryPricing.Monthly.Net, 64); err != nil {
				c.logger.Error("Failed to parse primary IP costs",
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.PrimaryIP,
					prometheus.GaugeValue,
					net,
					"",
					"net",
					string(primaryType.Type),
					primaryPricing.Location,
				)
			}

		}
	}

	for _, serverType := range pricing.ServerTypes {
		for _, serverPricing := range serverType.Pricings {
			if gross, err := strconv.ParseFloat(serverPricing.PerTBTraffic.Gross, 64); err != nil {
				c.logger.Error("Failed to parse traffic costs",
					"service", "server",
					"type", serverType.ServerType.Name,
					"locations", serverPricing.Location.Name,
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.ServersTraffic,
					prometheus.GaugeValue,
					gross,
					serverPricing.PerTBTraffic.Currency,
					"gross",
					serverType.ServerType.Name,
					serverPricing.Location.Name,
				)
			}

			if net, err := strconv.ParseFloat(serverPricing.PerTBTraffic.Net, 64); err != nil {
				c.logger.Error("Failed to parse traffic costs",
					"service", "server",
					"type", serverType.ServerType.Name,
					"locations", serverPricing.Location.Name,
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.ServersTraffic,
					prometheus.GaugeValue,
					net,
					serverPricing.PerTBTraffic.Currency,
					"net",
					serverType.ServerType.Name,
					serverPricing.Location.Name,
				)
			}

			if gross, err := strconv.ParseFloat(serverPricing.Monthly.Gross, 64); err != nil {
				c.logger.Error("Failed to parse server costs",
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
				c.logger.Error("Failed to parse server costs",
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
			if gross, err := strconv.ParseFloat(lbPricing.PerTBTraffic.Gross, 64); err != nil {
				c.logger.Error("Failed to parse traffic costs",
					"service", "server",
					"type", lbType.LoadBalancerType.Name,
					"locations", lbPricing.Location.Name,
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.LoadBalancersTraffic,
					prometheus.GaugeValue,
					gross,
					lbPricing.PerTBTraffic.Currency,
					"gross",
					lbType.LoadBalancerType.Name,
					lbPricing.Location.Name,
				)
			}

			if net, err := strconv.ParseFloat(lbPricing.PerTBTraffic.Net, 64); err != nil {
				c.logger.Error("Failed to parse traffic costs",
					"service", "server",
					"type", lbType.LoadBalancerType.Name,
					"locations", lbPricing.Location.Name,
					"err", err,
				)

				c.failures.WithLabelValues("pricing").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.LoadBalancersTraffic,
					prometheus.GaugeValue,
					net,
					lbPricing.PerTBTraffic.Currency,
					"net",
					lbType.LoadBalancerType.Name,
					lbPricing.Location.Name,
				)
			}

			if gross, err := strconv.ParseFloat(lbPricing.Monthly.Gross, 64); err != nil {
				c.logger.Error("Failed to parse load balancer costs",
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
				c.logger.Error("Failed to parse load balancer costs",
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

	c.logger.Debug("Processed pricing collector",
		"duration", time.Since(now),
	)

	c.duration.WithLabelValues("pricing").Observe(time.Since(now).Seconds())
}
