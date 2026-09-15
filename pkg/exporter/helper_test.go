package exporter

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testFailures() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{Name: "failures"}, []string{"collector"})
}

func testDuration() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "duration"}, []string{"collector"})
}

// testClient returns a client talking to a fake API served by the handler.
func testClient(t *testing.T, handler http.HandlerFunc) *hcloud.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return hcloud.NewClient(
		hcloud.WithEndpoint(server.URL),
		hcloud.WithToken("token"),
	)
}

// gatherLabels collects through a pedantic registry and returns the labels of
// every series within the given metric family.
func gatherLabels(t *testing.T, collector prometheus.Collector, family string) []map[string]string {
	t.Helper()

	registry := prometheus.NewPedanticRegistry()
	require.NoError(t, registry.Register(collector))

	families, err := registry.Gather()
	require.NoError(t, err)

	result := []map[string]string{}

	for _, f := range families {
		if f.GetName() != family {
			continue
		}

		for _, metric := range f.GetMetric() {
			labels := map[string]string{}

			for _, label := range metric.GetLabel() {
				labels[label.GetName()] = label.GetValue()
			}

			result = append(result, labels)
		}
	}

	return result
}

// spareCapacity leaves room in the backing array of the labels, the same way
// urfave/cli does for lists parsed from a flag or an environment variable.
func spareCapacity(labels []string) []string {
	return append(make([]string, 0, len(labels)+4), labels...)
}
