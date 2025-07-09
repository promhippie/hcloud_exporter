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

// StorageBoxes defines the storagebox specific configuration.
type StorageBoxes struct {
	Labels []string
}

// Target defines the target specific configuration.
type Target struct {
	Token        string
	Timeout      time.Duration
	StorageBoxes StorageBoxes
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
