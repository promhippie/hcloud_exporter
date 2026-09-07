package exporter

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The volumes endpoint embeds the attached server as a bare ID, so the name has
// to come from the servers endpoint. Payloads trimmed from real API responses.
const volumesPayload = `{
  "volumes": [
    {
      "id": 1001,
      "name": "attached-volume",
      "server": 42,
      "size": 10,
      "status": "available",
      "location": {"name": "nbg1"},
      "protection": {"delete": false},
      "created": "2024-01-01T00:00:00+00:00"
    },
    {
      "id": 1002,
      "name": "unattached-volume",
      "server": null,
      "size": 10,
      "status": "available",
      "location": {"name": "nbg1"},
      "protection": {"delete": false},
      "created": "2024-01-01T00:00:00+00:00"
    }
  ],
  "meta": {"pagination": {"page": 1, "per_page": 50, "total_entries": 2}}
}`

const serversPayload = `{
  "servers": [
    {
      "id": 42,
      "name": "example-server",
      "status": "running",
      "datacenter": {"name": "nbg1-dc3"},
      "server_type": {"name": "cx22", "cores": 2, "memory": 4, "disk": 40},
      "created": "2024-01-01T00:00:00+00:00"
    }
  ],
  "meta": {"pagination": {"page": 1, "per_page": 50, "total_entries": 1}}
}`

func volumeServerLabels(t *testing.T, handler http.HandlerFunc) map[string]string {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := hcloud.NewClient(
		hcloud.WithEndpoint(server.URL),
		hcloud.WithToken("token"),
	)

	collector := NewVolumeCollector(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		client,
		prometheus.NewCounterVec(prometheus.CounterOpts{Name: "failures"}, []string{"collector"}),
		prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "duration"}, []string{"collector"}),
		config.Target{Timeout: 5 * time.Second},
	)

	registry := prometheus.NewPedanticRegistry()
	require.NoError(t, registry.Register(collector))

	families, err := registry.Gather()
	require.NoError(t, err)

	labels := make(map[string]string)

	for _, family := range families {
		if family.GetName() != "hcloud_volume_size" {
			continue
		}

		for _, metric := range family.GetMetric() {
			var name, srv string

			for _, label := range metric.GetLabel() {
				switch label.GetName() {
				case "name":
					name = label.GetValue()
				case "server":
					srv = label.GetValue()
				}
			}

			labels[name] = srv
		}
	}

	return labels
}

func TestVolumeCollectorResolvesServerName(t *testing.T) {
	labels := volumeServerLabels(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/volumes":
			_, _ = w.Write([]byte(volumesPayload))
		case "/servers":
			_, _ = w.Write([]byte(serversPayload))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	assert.Equal(t, "example-server", labels["attached-volume"],
		"an attached volume should carry the name of its server")
	assert.Empty(t, labels["unattached-volume"],
		"an unattached volume should carry an empty server label")
}

func TestVolumeCollectorKeepsCollectingWhenServersFail(t *testing.T) {
	labels := volumeServerLabels(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/volumes":
			_, _ = w.Write([]byte(volumesPayload))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	assert.Len(t, labels, 2, "volume metrics should still be exported")
	assert.Empty(t, labels["attached-volume"])
}
