package exporter

import (
	"net/http"
	"testing"
	"time"

	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/stretchr/testify/assert"
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
      "labels": {"environment": "production"},
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
      "labels": {},
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

func volumesAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/volumes":
		_, _ = w.Write([]byte(volumesPayload))
	case "/servers":
		_, _ = w.Write([]byte(serversPayload))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// volumeLabels maps the name of every volume to the value of the given label.
func volumeLabels(t *testing.T, labels []string, key string, handler http.HandlerFunc) map[string]string {
	t.Helper()

	collector := NewVolumeCollector(
		testLogger(),
		testClient(t, handler),
		testFailures(),
		testDuration(),
		config.Target{
			Timeout: 5 * time.Second,
			Volumes: config.Volumes{Labels: labels},
		},
	)

	result := make(map[string]string)

	for _, series := range gatherLabels(t, collector, "hcloud_volume_size") {
		result[series["name"]] = series[key]
	}

	return result
}

func TestVolumeCollectorResolvesServerName(t *testing.T) {
	labels := volumeLabels(t, config.VolumeLabels(), "server", volumesAPI)

	assert.Equal(t, "example-server", labels["attached-volume"],
		"an attached volume should carry the name of its server")
	assert.Empty(t, labels["unattached-volume"],
		"an unattached volume should carry an empty server label")
}

func TestVolumeCollectorKeepsCollectingWhenServersFail(t *testing.T) {
	labels := volumeLabels(t, config.VolumeLabels(), "server", func(w http.ResponseWriter, r *http.Request) {
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

func TestVolumeCollectorResolvesResourceLabels(t *testing.T) {
	labels := volumeLabels(t, []string{"name", "server", "environment"}, "environment", volumesAPI)

	assert.Equal(t, "production", labels["attached-volume"],
		"a label of the volume should be exported with its value")
	assert.Empty(t, labels["unattached-volume"],
		"a volume without that label should carry an empty value")
}
