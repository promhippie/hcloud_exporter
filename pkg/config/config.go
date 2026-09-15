package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

// Server defines the general server configuration.
type Server struct {
	Addr    string
	Path    string
	Timeout time.Duration
	Web     string
	Pprof   bool
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string
	Pretty bool
}

// FloatingIPs defines the floating IP specific configuration.
type FloatingIPs struct {
	Labels []string
}

// Images defines the image specific configuration.
type Images struct {
	Labels []string
}

// Servers defines the server specific configuration.
type Servers struct {
	Labels []string
}

// ServerMetrics defines the server metrics specific configuration.
type ServerMetrics struct {
	Labels []string
}

// LoadBalancers defines the load balancer specific configuration.
type LoadBalancers struct {
	Labels []string
}

// SSHKeys defines the SSH key specific configuration.
type SSHKeys struct {
	Labels []string
}

// Volumes defines the volume specific configuration.
type Volumes struct {
	Labels []string
}

// StorageBoxes defines the storagebox specific configuration.
type StorageBoxes struct {
	Labels []string
}

// Target defines the target specific configuration.
type Target struct {
	Token         string
	Timeout       time.Duration
	FloatingIPs   FloatingIPs
	Images        Images
	Servers       Servers
	ServerMetrics ServerMetrics
	LoadBalancers LoadBalancers
	SSHKeys       SSHKeys
	Volumes       Volumes
	StorageBoxes  StorageBoxes
}

// Collector defines the collector specific configuration.
type Collector struct {
	FloatingIPs   bool
	Images        bool
	Pricing       bool
	Servers       bool
	ServerMetrics bool
	LoadBalancers bool
	SSHKeys       bool
	Volumes       bool
	StorageBoxes  bool
}

// Config is a combination of all available configurations.
type Config struct {
	Server    Server
	Logs      Logs
	Target    Target
	Collector Collector
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{}
}

// FloatingIPLabels defines the default labels used by floating IP collector.
func FloatingIPLabels() []string {
	return []string{
		"id",
		"server",
		"location",
		"type",
		"ip",
	}
}

// ImageLabels defines the default labels used by image collector.
func ImageLabels() []string {
	return []string{
		"id",
		"name",
		"type",
		"server",
		"flavor",
		"version",
	}
}

// ServerLabels defines the default labels used by server collector.
func ServerLabels() []string {
	return []string{
		"id",
		"name",
		"datacenter",
	}
}

// ServerMetricsLabels defines the default labels used by server metrics collector.
func ServerMetricsLabels() []string {
	return []string{
		"id",
		"name",
		"datacenter",
	}
}

// LoadBalancerLabels defines the default labels used by load balancer collector.
func LoadBalancerLabels() []string {
	return []string{
		"id",
		"name",
		"datacenter",
	}
}

// SSHKeyLabels defines the default labels used by SSH key collector.
func SSHKeyLabels() []string {
	return []string{
		"id",
		"name",
		"fingerprint",
	}
}

// VolumeLabels defines the default labels used by volume collector.
func VolumeLabels() []string {
	return []string{
		"id",
		"server",
		"location",
		"name",
	}
}

// StorageBoxLabels defines the default labels used by storagebox collector.
func StorageBoxLabels() []string {
	return []string{
		"id",
		"name",
		"location",
	}
}

// Value returns the config value based on a DSN.
func Value(val string) (string, error) {
	if strings.HasPrefix(val, "file://") {
		content, err := os.ReadFile(
			strings.TrimPrefix(val, "file://"),
		)

		if err != nil {
			return "", fmt.Errorf("failed to parse secret file: %w", err)
		}

		return string(content), nil
	}

	if strings.HasPrefix(val, "base64://") {
		content, err := base64.StdEncoding.DecodeString(
			strings.TrimPrefix(val, "base64://"),
		)

		if err != nil {
			return "", fmt.Errorf("failed to parse base64 value: %w", err)
		}

		return string(content), nil
	}

	return val, nil
}
