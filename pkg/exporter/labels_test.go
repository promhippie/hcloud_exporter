package exporter

import (
	"net"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/promhippie/hcloud_exporter/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCollectorsResolveLabels(t *testing.T) {
	var (
		logger   = testLogger()
		cfg      = config.Target{}
		labels   = map[string]string{"environment": "production"}
		location = &hcloud.Location{Name: "nbg1"}
	)

	servers := NewServerCollector(logger, nil, nil, nil, cfg)
	serverMetrics := NewServerMetricsCollector(logger, nil, nil, nil, cfg)
	volumes := NewVolumeCollector(logger, nil, nil, nil, cfg)
	floatingIPs := NewFloatingIPCollector(logger, nil, nil, nil, cfg)
	images := NewImageCollector(logger, nil, nil, nil, cfg)
	loadBalancers := NewLoadBalancerCollector(logger, nil, nil, nil, cfg)
	sshKeys := NewSSHKeyCollector(logger, nil, nil, nil, cfg)

	server := &hcloud.Server{ID: 1, Name: "example-server", Location: location, Labels: labels}
	volume := &hcloud.Volume{ID: 2, Name: "example-volume", Location: location, Labels: labels}
	ip := &hcloud.FloatingIP{ID: 3, IP: net.ParseIP("192.0.2.1"), Type: hcloud.FloatingIPTypeIPv4, HomeLocation: location, Labels: labels}
	image := &hcloud.Image{ID: 4, Name: "example-image", Type: hcloud.ImageTypeSnapshot, OSFlavor: "ubuntu", OSVersion: "24.04", Labels: labels}
	lb := &hcloud.LoadBalancer{ID: 5, Name: "example-lb", Location: location, Labels: labels}
	key := &hcloud.SSHKey{ID: 6, Name: "example-key", Fingerprint: "b7:2f:30:a0", Labels: labels}

	tests := []struct {
		name     string
		byLabel  func(string) string
		builtins map[string]string
	}{
		{
			name:     "server",
			byLabel:  func(label string) string { return servers.byLabel(server, label) },
			builtins: map[string]string{"id": "1", "name": "example-server", "datacenter": "nbg1"},
		},
		{
			name:     "server metrics",
			byLabel:  func(label string) string { return serverMetrics.byLabel(server, label) },
			builtins: map[string]string{"id": "1", "name": "example-server", "datacenter": "nbg1"},
		},
		{
			name:     "volume",
			byLabel:  func(label string) string { return volumes.byLabel(volume, "example-server", label) },
			builtins: map[string]string{"id": "2", "server": "example-server", "location": "nbg1", "name": "example-volume"},
		},
		{
			name:     "floating ip",
			byLabel:  func(label string) string { return floatingIPs.byLabel(ip, "example-server", label) },
			builtins: map[string]string{"id": "3", "server": "example-server", "location": "nbg1", "type": "ipv4", "ip": "192.0.2.1"},
		},
		{
			name:     "image",
			byLabel:  func(label string) string { return images.byLabel(image, "example-server", label) },
			builtins: map[string]string{"id": "4", "name": "example-image", "type": "snapshot", "server": "example-server", "flavor": "ubuntu", "version": "24.04"},
		},
		{
			name:     "load balancer",
			byLabel:  func(label string) string { return loadBalancers.byLabel(lb, label) },
			builtins: map[string]string{"id": "5", "name": "example-lb", "datacenter": "nbg1"},
		},
		{
			name:     "ssh key",
			byLabel:  func(label string) string { return sshKeys.byLabel(key, label) },
			builtins: map[string]string{"id": "6", "name": "example-key", "fingerprint": "b7:2f:30:a0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for label, value := range tt.builtins {
				assert.Equal(t, value, tt.byLabel(label), "built-in label %q", label)
			}

			assert.Equal(t, "production", tt.byLabel("environment"),
				"any other label should be looked up in the labels of the resource")
			assert.Empty(t, tt.byLabel("missing"),
				"a label the resource does not carry should stay empty")
		})
	}
}
