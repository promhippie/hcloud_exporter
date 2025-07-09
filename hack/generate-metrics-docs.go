//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/promhippie/hcloud_exporter/pkg/exporter"
)

type metric struct {
	Name   string
	Help   string
	Labels []string
}

func main() {
	collectors := make([]*prometheus.Desc, 0)

	cfg := config.Load().Target
	cfg.StorageBoxes.Labels = config.StorageBoxLabels()

	collectors = append(
		collectors,
		exporter.NewFloatingIPCollector(slog.Default(), nil, nil, nil, cfg).Metrics()...,
	)

	collectors = append(
		collectors,
		exporter.NewImageCollector(slog.Default(), nil, nil, nil, cfg).Metrics()...,
	)

	collectors = append(
		collectors,
		exporter.NewPricingCollector(slog.Default(), nil, nil, nil, cfg).Metrics()...,
	)

	collectors = append(
		collectors,
		exporter.NewServerCollector(slog.Default(), nil, nil, nil, cfg).Metrics()...,
	)

	collectors = append(
		collectors,
		exporter.NewServerMetricsCollector(slog.Default(), nil, nil, nil, cfg).Metrics()...,
	)

	collectors = append(
		collectors,
		exporter.NewSSHKeyCollector(slog.Default(), nil, nil, nil, cfg).Metrics()...,
	)

	collectors = append(
		collectors,
		exporter.NewVolumeCollector(slog.Default(), nil, nil, nil, cfg).Metrics()...,
	)

	collectors = append(
		collectors,
		exporter.NewStorageBoxCollector(slog.Default(), nil, nil, nil, cfg).Metrics()...,
	)

	metrics := make([]metric, 0)

	metrics = append(metrics, metric{
		Name:   "hcloud_request_duration_seconds",
		Help:   "Histogram of latencies for requests to the api per collector",
		Labels: []string{"collector"},
	})

	metrics = append(metrics, metric{
		Name:   "hcloud_request_failures_total",
		Help:   "Total number of failed requests to the api per collector",
		Labels: []string{"collector"},
	})

	for _, desc := range collectors {
		m := metric{
			Name:   reflect.ValueOf(desc).Elem().FieldByName("fqName").String(),
			Help:   reflect.ValueOf(desc).Elem().FieldByName("help").String(),
			Labels: make([]string, 0),
		}

		labels := reflect.Indirect(
			reflect.ValueOf(desc).Elem().FieldByName("variableLabels"),
		).FieldByName("names")

		for i := 0; i < labels.Len(); i++ {
			m.Labels = append(m.Labels, labels.Index(i).String())
		}

		metrics = append(metrics, m)
	}

	sort.SliceStable(metrics, func(i, j int) bool {
		return metrics[i].Name < metrics[j].Name
	})

	f, err := os.Create("docs/partials/metrics.md")

	if err != nil {
		fmt.Printf("failed to create file")
		os.Exit(1)
	}

	defer f.Close()

	last := metrics[len(metrics)-1]
	for _, m := range metrics {
		f.WriteString(fmt.Sprintf(
			"%s{%s}\n",
			m.Name,
			strings.Join(
				m.Labels,
				", ",
			),
		))

		f.WriteString(fmt.Sprintf(
			": %s\n",
			m.Help,
		))

		if m.Name != last.Name {
			f.WriteString("\n")
		}
	}
}
