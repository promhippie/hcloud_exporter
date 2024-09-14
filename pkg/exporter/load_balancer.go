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

// LoadBalancerCollector collects metrics about the load balancers.
type LoadBalancerCollector struct {
	client   *hcloud.Client
	logger   *slog.Logger
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
	TargetsUnknown          *prometheus.Desc
	AssignedCertificates    *prometheus.Desc
	MaxAssignedCertificates *prometheus.Desc
	IngoingTraffic          *prometheus.Desc
	OutgoingTraffic         *prometheus.Desc
	IncludedTraffic         *prometheus.Desc
}

// NewLoadBalancerCollector returns a new LoadBalancerCollector.
func NewLoadBalancerCollector(logger *slog.Logger, client *hcloud.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *LoadBalancerCollector {
	if failures != nil {
		failures.WithLabelValues("load_balancer").Add(0)
	}

	labels := []string{"id", "name", "datacenter"}
	return &LoadBalancerCollector{
		client:   client,
		logger:   logger.With("collector", "load-balancer"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Created: prometheus.NewDesc(
			"hcloud_loadbalancer_created_timestamp",
			"Timestamp when the load balancer have been created",
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
		TargetsUnknown: prometheus.NewDesc(
			"hcloud_loadbalancer_targets_unknown",
			"The number of unknown targets",
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
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *LoadBalancerCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Created,
		c.Services,
		c.MaxServices,
		c.Targets,
		c.MaxTargets,
		c.TargetsHealthy,
		c.TargetsUnhealthy,
		c.TargetsUnknown,
		c.AssignedCertificates,
		c.MaxAssignedCertificates,
		c.IngoingTraffic,
		c.OutgoingTraffic,
		c.IncludedTraffic,
		c.Connections,
		c.MaxConnections,
		c.ConnectionsPerSecond,
		c.RequestsPerSecond,
		c.IncomingBandwidth,
		c.OutgoingBandwidth,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *LoadBalancerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Created
	ch <- c.Services
	ch <- c.MaxServices
	ch <- c.Targets
	ch <- c.MaxTargets
	ch <- c.TargetsHealthy
	ch <- c.TargetsUnhealthy
	ch <- c.TargetsUnknown
	ch <- c.AssignedCertificates
	ch <- c.MaxAssignedCertificates
	ch <- c.IngoingTraffic
	ch <- c.OutgoingTraffic
	ch <- c.IncludedTraffic
	ch <- c.Connections
	ch <- c.MaxConnections
	ch <- c.ConnectionsPerSecond
	ch <- c.RequestsPerSecond
	ch <- c.IncomingBandwidth
	ch <- c.OutgoingBandwidth
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *LoadBalancerCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	lbs, err := c.client.LoadBalancer.All(ctx)

	if err != nil {
		c.logger.Error("Failed to fetch load balancers",
			"err", err,
		)

		c.failures.WithLabelValues("load_balancer").Inc()
		return
	}

	c.logger.Debug("Fetched load balancers",
		"count", len(lbs),
	)

	for _, lb := range lbs {
		var (
			targetsHealthy       int
			targetsUnhealthy     int
			targetsUnknown       int
			assignedCertificates int
		)

		labels := []string{
			strconv.FormatInt(lb.ID, 10),
			lb.Name,
			lb.Location.Name,
		}

		ch <- prometheus.MustNewConstMetric(
			c.Created,
			prometheus.GaugeValue,
			float64(lb.Created.Unix()),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Services,
			prometheus.GaugeValue,
			float64(len(lb.Services)),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MaxServices,
			prometheus.GaugeValue,
			float64(lb.LoadBalancerType.MaxServices),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Targets,
			prometheus.GaugeValue,
			float64(len(lb.Targets)),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MaxTargets,
			prometheus.GaugeValue,
			float64(lb.LoadBalancerType.MaxTargets),
			labels...,
		)

		for _, target := range lb.Targets {
			match := false

			for _, healthStatus := range target.HealthStatus {
				if healthStatus.Status == hcloud.LoadBalancerTargetHealthStatusStatusHealthy {
					match = true
				}
			}

			if match {
				targetsHealthy++
			}
		}

		ch <- prometheus.MustNewConstMetric(
			c.TargetsHealthy,
			prometheus.GaugeValue,
			float64(targetsHealthy),
			labels...,
		)

		for _, target := range lb.Targets {
			match := false

			for _, healthStatus := range target.HealthStatus {
				if healthStatus.Status == hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy {
					match = true
				}
			}

			if match {
				targetsUnhealthy++
			}
		}

		ch <- prometheus.MustNewConstMetric(
			c.TargetsUnhealthy,
			prometheus.GaugeValue,
			float64(targetsUnhealthy),
			labels...,
		)

		for _, target := range lb.Targets {
			match := false

			for _, healthStatus := range target.HealthStatus {
				if healthStatus.Status == hcloud.LoadBalancerTargetHealthStatusStatusUnknown {
					match = true
				}
			}

			if match {
				targetsUnknown++
			}
		}

		ch <- prometheus.MustNewConstMetric(
			c.TargetsUnknown,
			prometheus.GaugeValue,
			float64(targetsUnknown),
			labels...,
		)

		for _, service := range lb.Services {
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
			float64(lb.LoadBalancerType.MaxAssignedCertificates),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IngoingTraffic,
			prometheus.GaugeValue,
			float64(lb.IngoingTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.OutgoingTraffic,
			prometheus.GaugeValue,
			float64(lb.OutgoingTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IncludedTraffic,
			prometheus.GaugeValue,
			float64(lb.IncludedTraffic),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.MaxConnections,
			prometheus.GaugeValue,
			float64(lb.LoadBalancerType.MaxConnections),
			labels...,
		)

		metrics, _, err := c.client.LoadBalancer.GetMetrics(ctx, lb, hcloud.LoadBalancerGetMetricsOpts{
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
			c.logger.Error("Failed to fetch load balancer metrics",
				"load-balancer", lb.Name,
				"err", err,
			)

			c.failures.WithLabelValues("load_balancer").Inc()
			return
		}

		c.addTimeSeries(
			ch,
			c.Connections,
			"open_connections",
			lb,
			*metrics,
			labels,
		)

		c.addTimeSeries(
			ch,
			c.ConnectionsPerSecond,
			"connections_per_second",
			lb,
			*metrics,
			labels,
		)

		c.addTimeSeries(
			ch,
			c.RequestsPerSecond,
			"requests_per_second",
			lb,
			*metrics,
			labels,
		)

		c.addTimeSeries(
			ch,
			c.IncomingBandwidth,
			"bandwidth.in",
			lb,
			*metrics,
			labels,
		)

		c.addTimeSeries(
			ch,
			c.OutgoingBandwidth,
			"bandwidth.out",
			lb,
			*metrics,
			labels,
		)
	}

	c.logger.Debug("Processed load balancer collector",
		"duration", time.Since(now),
	)

	c.duration.WithLabelValues("load_balancer").Observe(time.Since(now).Seconds())
}

func (c *LoadBalancerCollector) addTimeSeries(
	ch chan<- prometheus.Metric,
	desc *prometheus.Desc,
	name string,
	lb *hcloud.LoadBalancer,
	metrics hcloud.LoadBalancerMetrics,
	labels []string,
) {
	if metric, ok := metrics.TimeSeries[name]; ok {
		if len(metric) > 0 {
			value, err := strconv.ParseFloat(metric[0].Value, 64)

			if err != nil {
				c.logger.Error("Failed to parse load balancer metric",
					"load-balancer", lb.Name,
					"name", name,
					"err", err,
				)

				c.failures.WithLabelValues("load_balancer").Inc()
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
