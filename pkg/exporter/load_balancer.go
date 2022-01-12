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

// LoadBalancerCollector collects metrics about the load balancers.
type LoadBalancerCollector struct {
	client   *hcloud.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Created                 *prometheus.Desc
	Connections             *prometheus.Desc
	MaxConnections          *prometheus.Desc
	ConnectionsPerSecond    *prometheus.Desc
	RequestsPerSecond       *prometheus.Desc
	IncomingBandwidth       *prometheus.Desc
	OutgoingBandwidth       *prometheus.Desc
	Services                *prometheus.Desc
	MaxServices             *prometheus.Desc
	Targets                 *prometheus.Desc
	MaxTargets              *prometheus.Desc
	TargetsHealthy          *prometheus.Desc
	TargetsUnhealthy        *prometheus.Desc
	AssignedCertificates    *prometheus.Desc
	MaxAssignedCertificates *prometheus.Desc
	IngoingTraffic          *prometheus.Desc
	OutgoingTraffic         *prometheus.Desc
	IncludedTraffic         *prometheus.Desc
	PriceHourly             *prometheus.Desc
	PriceMonthly            *prometheus.Desc
}

// NewLoadBalancerCollector returns a new LoadBalancerCollector.
func NewLoadBalancerCollector(logger log.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *LoadBalancerCollector {
	if failures != nil {
		failures.WithLabelValues("load-balancer").Add(0)
	}

	labels := []string{"id", "name", "datacenter"}
	pricingLabels := append(labels, "vat")

	return &LoadBalancerCollector{
		client:   client,
		logger:   log.With(logger, "collector", "load-balancer"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Created: prometheus.NewDesc(
			"hcloud_loadbalancer_created_timestamp",
			"Timestamp when the load balancer have been created",
			labels,
			nil,
		),
		Connections: prometheus.NewDesc(
			"hcloud_loadbalancer_open_connections",
			"The number of open connections",
			labels,
			nil,
		),
		MaxConnections: prometheus.NewDesc(
			"hcloud_loadbalancer_max_open_connections",
			"The maximum number of open connections",
			labels,
			nil,
		),
		ConnectionsPerSecond: prometheus.NewDesc(
			"hcloud_loadbalancer_connections_per_second",
			"The number of new connections per second",
			labels,
			nil,
		),
		RequestsPerSecond: prometheus.NewDesc(
			"hcloud_loadbalancer_requests_per_second",
			"The number of requests per second",
			labels,
			nil,
		),
		IncomingBandwidth: prometheus.NewDesc(
			"hcloud_loadbalancer_open_connections_bandwidth_in",
			"The incoming bandwidth in bytes per second",
			labels,
			nil,
		),
		OutgoingBandwidth: prometheus.NewDesc(
			"hcloud_loadbalancer_open_connections_bandwidth_out",
			"The outgoing bandwidth in bytes per second",
			labels,
			nil,
		),
		Services: prometheus.NewDesc(
			"hcloud_loadbalancer_services",
			"The number of configured services",
			labels,
			nil,
		),
		MaxServices: prometheus.NewDesc(
			"hcloud_loadbalancer_max_services",
			"The maximum number of services that can be configured",
			labels,
			nil,
		),
		Targets: prometheus.NewDesc(
			"hcloud_loadbalancer_targets",
			"The number of targets",
			labels,
			nil,
		),
		MaxTargets: prometheus.NewDesc(
			"hcloud_loadbalancer_max_targets",
			"The maximum number of targets",
			labels,
			nil,
		),
		TargetsHealthy: prometheus.NewDesc(
			"hcloud_loadbalancer_targets_healthy",
			"The number of healthy targets",
			labels,
			nil,
		),
		TargetsUnhealthy: prometheus.NewDesc(
			"hcloud_loadbalancer_targets_unhealthy",
			"The number of unhealthy targets",
			labels,
			nil,
		),
		AssignedCertificates: prometheus.NewDesc(
			"hcloud_loadbalancer_assigned_certificates",
			"The number of assigned certificates",
			labels,
			nil,
		),
		MaxAssignedCertificates: prometheus.NewDesc(
			"hcloud_loadbalancer_max_assigned_certificates",
			"The maximum number of certificates that can be assigned",
			labels,
			nil,
		),
		IngoingTraffic: prometheus.NewDesc(
			"hcloud_loadbalancer_ingoing_traffic",
			"The total amount of ingoing traffic in bytes",
			labels,
			nil,
		),
		OutgoingTraffic: prometheus.NewDesc(
			"hcloud_loadbalancer_outgoing_traffic",
			"The total amount of outgoing traffic in bytes",
			labels,
			nil,
		),
		IncludedTraffic: prometheus.NewDesc(
			"hcloud_loadbalancer_included_traffic",
			"The amount of traffic that is included for the load balancer type in bytes",
			labels,
			nil,
		),
		PriceHourly: prometheus.NewDesc(
			"hcloud_loadbalancer_price_hourly",
			"Price of the load balancer billed hourly in €",
			pricingLabels,
			nil,
		),
		PriceMonthly: prometheus.NewDesc(
			"hcloud_loadbalancer_price_monthly",
			"Price of the load balancer billed monthly in €",
			pricingLabels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *LoadBalancerCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Created,
		c.Connections,
		c.MaxConnections,
		c.ConnectionsPerSecond,
		c.RequestsPerSecond,
		c.IncomingBandwidth,
		c.OutgoingBandwidth,
		c.Services,
		c.MaxServices,
		c.Targets,
		c.MaxTargets,
		c.TargetsHealthy,
		c.TargetsUnhealthy,
		c.AssignedCertificates,
		c.MaxAssignedCertificates,
		c.IngoingTraffic,
		c.OutgoingTraffic,
		c.IncludedTraffic,
		c.PriceHourly,
		c.PriceMonthly,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *LoadBalancerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Created
	ch <- c.Connections
	ch <- c.MaxConnections
	ch <- c.ConnectionsPerSecond
	ch <- c.RequestsPerSecond
	ch <- c.IncomingBandwidth
	ch <- c.OutgoingBandwidth
	ch <- c.Services
	ch <- c.MaxServices
	ch <- c.Targets
	ch <- c.MaxTargets
	ch <- c.TargetsHealthy
	ch <- c.TargetsUnhealthy
	ch <- c.AssignedCertificates
	ch <- c.MaxAssignedCertificates
	ch <- c.IngoingTraffic
	ch <- c.OutgoingTraffic
	ch <- c.IncludedTraffic
	ch <- c.PriceHourly
	ch <- c.PriceMonthly
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *LoadBalancerCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	loadBalancers, err := c.client.LoadBalancer.All(ctx)
	c.duration.WithLabelValues("load-balancer").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch load balancers",
			"err", err,
		)

		c.failures.WithLabelValues("load-balancer").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched load balancer",
		"count", len(loadBalancers),
	)

	for _, loadBalancer := range loadBalancers {
		labels := []string{
			strconv.Itoa(loadBalancer.ID),
			loadBalancer.Name,
			loadBalancer.Location.Name,
		}

		loadBalancerMetrics, _, err := c.client.LoadBalancer.GetMetrics(ctx, loadBalancer, hcloud.LoadBalancerGetMetricsOpts{
			Types: []hcloud.LoadBalancerMetricType{
				hcloud.LoadBalancerMetricOpenConnections,
				hcloud.LoadBalancerMetricConnectionsPerSecond,
				hcloud.LoadBalancerMetricRequestsPerSecond,
				hcloud.LoadBalancerMetricBandwidth,
			},
			Start: now,
			End:   now,
			Step:  1,
		})
		if err != nil {
			level.Error(c.logger).Log(
				"msg", "Failed to fetch load balancer metrics",
				"err", err,
			)

			c.failures.WithLabelValues("load-balancer").Inc()
			return
		}

		ch <- prometheus.MustNewConstMetric(
			c.Created,
			prometheus.GaugeValue,
			float64(loadBalancer.Created.Unix()),
			labels...,
		)

		c.addMetric(ch, c.Connections, "open_connections", *loadBalancerMetrics, labels)

		ch <- prometheus.MustNewConstMetric(
			c.MaxConnections,
			prometheus.GaugeValue,
			float64(loadBalancer.LoadBalancerType.MaxConnections),
			labels...,
		)

		c.addMetric(ch, c.ConnectionsPerSecond, "connections_per_second", *loadBalancerMetrics, labels)

		c.addMetric(ch, c.RequestsPerSecond, "requests_per_second", *loadBalancerMetrics, labels)

		c.addMetric(ch, c.IncomingBandwidth, "bandwidth_in", *loadBalancerMetrics, labels)

		c.addMetric(ch, c.OutgoingBandwidth, "bandwidth_out", *loadBalancerMetrics, labels)

		ch <- prometheus.MustNewConstMetric(
			c.Services,
			prometheus.GaugeValue,
			float64(len(loadBalancer.Services)),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MaxServices,
			prometheus.GaugeValue,
			float64(loadBalancer.LoadBalancerType.MaxServices),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Targets,
			prometheus.GaugeValue,
			float64(len(loadBalancer.Targets)),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MaxTargets,
			prometheus.GaugeValue,
			float64(loadBalancer.LoadBalancerType.MaxTargets),
			labels...,
		)

		targetsHealthy := 0
		targetsUnhealthy := 0
		for _, target := range loadBalancer.Targets {
			isHealthy := true
			for _, healthStatus := range target.HealthStatus {
				if healthStatus.Status != "healthy" {
					isHealthy = false
				}
			}

			if isHealthy {
				targetsHealthy += 1
			} else {
				targetsUnhealthy += 1
			}
		}

		ch <- prometheus.MustNewConstMetric(
			c.TargetsHealthy,
			prometheus.GaugeValue,
			float64(targetsHealthy),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.TargetsUnhealthy,
			prometheus.GaugeValue,
			float64(targetsUnhealthy),
			labels...,
		)

		assignedCertificates := 0
		for _, service := range loadBalancer.Services {
			if service.Protocol == "https" {
				assignedCertificates += len(service.HTTP.Certificates)
			}
		}

		ch <- prometheus.MustNewConstMetric(
			c.AssignedCertificates,
			prometheus.GaugeValue,
			float64(assignedCertificates),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MaxAssignedCertificates,
			prometheus.GaugeValue,
			float64(loadBalancer.LoadBalancerType.MaxAssignedCertificates),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IngoingTraffic,
			prometheus.GaugeValue,
			float64(loadBalancer.IngoingTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.OutgoingTraffic,
			prometheus.GaugeValue,
			float64(loadBalancer.OutgoingTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IncludedTraffic,
			prometheus.GaugeValue,
			float64(loadBalancer.IncludedTraffic),
			labels...,
		)

		labelsNet := append(labels, "net")
		labelsGross := append(labels, "gross")

		for _, pricing := range loadBalancer.LoadBalancerType.Pricings {
			if loadBalancer.Location.Name == pricing.Location.Name {
				hourlyNet, _ := strconv.ParseFloat(pricing.Hourly.Net, 64)
				hourlyGross, _ := strconv.ParseFloat(pricing.Hourly.Gross, 64)

				ch <- prometheus.MustNewConstMetric(
					c.PriceHourly,
					prometheus.GaugeValue,
					hourlyNet,
					labelsNet...,
				)

				ch <- prometheus.MustNewConstMetric(
					c.PriceHourly,
					prometheus.GaugeValue,
					hourlyGross,
					labelsGross...,
				)

				monthlyNet, _ := strconv.ParseFloat(pricing.Monthly.Net, 64)
				monthlyGross, _ := strconv.ParseFloat(pricing.Monthly.Gross, 64)

				ch <- prometheus.MustNewConstMetric(
					c.PriceMonthly,
					prometheus.GaugeValue,
					monthlyNet,
					labelsNet...,
				)

				ch <- prometheus.MustNewConstMetric(
					c.PriceMonthly,
					prometheus.GaugeValue,
					monthlyGross,
					labelsGross...,
				)
			}
		}
	}
}

func (c *LoadBalancerCollector) addMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, metric string, loadBalancerMetrics hcloud.LoadBalancerMetrics, labels []string) {
	if metric, ok := loadBalancerMetrics.TimeSeries[metric]; ok {
		if len(metric) > 0 {
			value, err := strconv.ParseFloat(metric[0].Value, 64)
			if err != nil {
				level.Error(c.logger).Log(
					"msg", "Failed to parse value of open_connections",
					"err", err,
				)

				c.failures.WithLabelValues("load-balancer").Inc()
				return
			}

			ch <- prometheus.MustNewConstMetric(
				desc,
				prometheus.GaugeValue,
				value,
				labels...,
			)
		}
	}
}
