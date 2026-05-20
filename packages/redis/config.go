package redis

import "time"

// ConnectionConfig is the configuration for a single-node connection.
// Field names mirror the upstream connector config keys.
type ConnectionConfig struct {
	Name     string
	Host     string
	Port     int
	URL      string // optional, overrides Host/Port if set
	Password string
	Username string
	Database int
	Prefix   string
	Scheme   string // "tcp", "tls", "unix"

	Timeout      time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	MaxRetries   int
	MinIdleConns int
	PoolSize     int

	// Cluster, if non-nil, promotes this connection to a cluster.
	Cluster *ClusterConfig

	// Sentinel, if non-nil, promotes this connection to Sentinel failover.
	Sentinel *SentinelConfig
}

// ClusterConfig is the upstream "cluster" options block.
type ClusterConfig struct {
	Addrs    []string
	Password string
	Username string
	Prefix   string

	ReadOnly       bool
	RouteByLatency bool
	RouteRandomly  bool
}

// SentinelConfig is the upstream "sentinel" options block.
type SentinelConfig struct {
	MasterName       string
	SentinelAddrs    []string
	SentinelPassword string
	Password         string
	Username         string
	Database         int
}
