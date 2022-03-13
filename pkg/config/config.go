package config

import (
	"time"
)

// Server defines the general server configuration.
type Server struct {
	Addr    string
	Path    string
	Timeout time.Duration
	Web     string
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string
	Pretty bool
}

// Target defines the target specific configuration.
type Target struct {
	Token   string
	Timeout time.Duration
}

// Collector defines the collector specific configuration.
type Collector struct {
	FloatingIPs   bool
	Images        bool
	Pricing       bool
	Servers       bool
	LoadBalancers bool
	SSHKeys       bool
	Volumes       bool
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
