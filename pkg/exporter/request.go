package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// RequestDuration tracks the duration of API requests by collector.
func RequestDuration() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "request_duration_seconds",
			Help:      "Histogram of latencies for requests to the HetznerCloud API per collector.",
			Buckets:   []float64{0.001, 0.01, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0},
		},
		[]string{"collector"},
	)
}

// RequestFailures tracks the number of API failures by collector.
func RequestFailures() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "request_failures_total",
			Help:      "Total number of failed requests to the HetznerCloud API per collector.",
		},
		[]string{"collector"},
	)
}
