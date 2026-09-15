package exporter

import (
	"net/http"
	"testing"
	"time"

	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const serverCollectorPayload = `{
  "servers": [
    {
      "id": 42,
      "name": "example-server",
      "status": "running",
      "location": {"name": "nbg1"},
      "server_type": {
        "name": "cx22",
        "cores": 2,
        "memory": 4,
        "disk": 40,
        "prices": [
          {
            "location": "nbg1",
            "price_hourly": {"net": "0.0060", "gross": "0.0071"},
            "price_monthly": {"net": "3.7900", "gross": "4.5101"}
          }
        ]
      },
      "labels": {"environment": "production"},
      "created": "2024-01-01T00:00:00+00:00"
    }
  ],
  "meta": {"pagination": {"page": 1, "per_page": 50, "total_entries": 1}}
}`

func serverCollector(t *testing.T, labels []string) *ServerCollector {
	t.Helper()

	return NewServerCollector(
		testLogger(),
		testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/servers" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(serverCollectorPayload))
		}),
		testFailures(),
		testDuration(),
		config.Target{
			Timeout: 5 * time.Second,
			Servers: config.Servers{Labels: labels},
		},
	)
}

// The label values get built by appending in a loop, which leaves spare capacity
// for most label counts, the default three included. Appending the vat label in
// place would let the gross price overwrite the net one and fail the whole
// scrape with a duplicate series.
func TestServerCollectorKeepsPricingLabelsApart(t *testing.T) {
	series := gatherLabels(t, serverCollector(t, config.ServerLabels()), "hcloud_server_price_hourly")

	vat := []string{}

	for _, labels := range series {
		vat = append(vat, labels["vat"])
	}

	assert.ElementsMatch(t, []string{"net", "gross"}, vat,
		"the hourly price should be exported once as net and once as gross")
}

func TestServerCollectorResolvesResourceLabels(t *testing.T) {
	series := gatherLabels(t, serverCollector(t, []string{"id", "name", "environment"}), "hcloud_server_running")

	require.Len(t, series, 1)
	assert.Equal(t, "production", series[0]["environment"],
		"a label of the server should be exported with its value")
}
