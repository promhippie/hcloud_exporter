package exporter

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func boolToFloat64(val bool) float64 {
	if val {
		return 1.0
	}

	return 0.0
}

// resolveServerNames maps server IDs to their names. The list endpoints for
// volumes and floating IPs only embed the ID of the attached server, so the
// name has to be looked up separately to fill the server label.
func resolveServerNames(ctx context.Context, client *hcloud.Client) (map[int64]string, error) {
	servers, err := client.Server.All(ctx)

	if err != nil {
		return nil, err
	}

	names := make(map[int64]string, len(servers))

	for _, server := range servers {
		names[server.ID] = server.Name
	}

	return names, nil
}
