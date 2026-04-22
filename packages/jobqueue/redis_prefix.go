package jobqueue

import (
	"errors"
	"strings"
)

const defaultRedisPrefix = "jobqueue:"

// ErrUnknownRedisConnection is returned when a named Redis connection is not registered.
var ErrUnknownRedisConnection = errors.New("jobqueue: redis connection not registered")

// RedisConnectionKind identifies where a JobQueue Redis connection is registered.
type RedisConnectionKind string

const (
	// RedisStandalone is a non-clustered Redis connection.
	RedisStandalone RedisConnectionKind = "standalone"
	// RedisCluster is a clustered Redis connection.
	RedisCluster RedisConnectionKind = "cluster"
)

// RedisConnectionConfig stores the Redis connection metadata JobQueue needs.
type RedisConnectionConfig struct {
	Name            string
	Prefix          string
	Kind            RedisConnectionKind
	Nodes           []string
	Options         map[string]string
	SupportsCluster bool
}

// RedisConnectionRegistry stores standalone and clustered Redis connections.
type RedisConnectionRegistry struct {
	fallbackPrefix string
	standalone     map[string]RedisConnectionConfig
	clusters       map[string]RedisConnectionConfig
}

// NewRedisConnectionRegistry creates an empty Redis connection registry.
func NewRedisConnectionRegistry(fallbackPrefix string) *RedisConnectionRegistry {
	if fallbackPrefix == "" {
		fallbackPrefix = defaultRedisPrefix
	}

	return &RedisConnectionRegistry{
		fallbackPrefix: fallbackPrefix,
		standalone:     make(map[string]RedisConnectionConfig),
		clusters:       make(map[string]RedisConnectionConfig),
	}
}

// Register stores a normalized Redis connection.
func (r *RedisConnectionRegistry) Register(config RedisConnectionConfig) {
	normalized := normalizeRedisConnection(config, r.fallbackPrefix)

	if normalized.Kind == RedisCluster && normalized.SupportsCluster {
		r.clusters[normalized.Name] = normalized
		return
	}

	normalized.Kind = RedisStandalone
	r.standalone[normalized.Name] = normalized
}

// Use returns a named Redis connection, preferring clusters over standalone connections.
func (r *RedisConnectionRegistry) Use(name string) (RedisConnectionConfig, error) {
	if config, ok := r.clusters[name]; ok {
		return cloneRedisConnectionConfig(config), nil
	}
	if config, ok := r.standalone[name]; ok {
		return cloneRedisConnectionConfig(config), nil
	}

	return RedisConnectionConfig{}, ErrUnknownRedisConnection
}

// Standalone reports whether a name is registered as a standalone connection.
func (r *RedisConnectionRegistry) Standalone(name string) bool {
	_, ok := r.standalone[name]

	return ok
}

// Cluster reports whether a name is registered as a cluster connection.
func (r *RedisConnectionRegistry) Cluster(name string) bool {
	_, ok := r.clusters[name]

	return ok
}

func normalizeRedisConnection(config RedisConnectionConfig, fallbackPrefix string) RedisConnectionConfig {
	normalized := cloneRedisConnectionConfig(config)
	if normalized.Prefix == "" {
		normalized.Prefix = fallbackPrefix
	}
	if normalized.Kind == "" {
		normalized.Kind = RedisStandalone
	}
	if normalized.Kind == RedisCluster {
		normalized.Prefix = hashTaggedPrefix(normalized.Prefix)
	}

	return normalized
}

func hashTaggedPrefix(prefix string) string {
	if prefix == "" {
		prefix = defaultRedisPrefix
	}
	if strings.HasPrefix(prefix, "{") {
		return prefix
	}

	trimmed := strings.TrimSuffix(prefix, ":")

	return "{" + trimmed + "}:"
}

func cloneRedisConnectionConfig(config RedisConnectionConfig) RedisConnectionConfig {
	config.Nodes = append([]string(nil), config.Nodes...)
	config.Options = cloneStringMap(config.Options)

	return config
}

func cloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}

	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}
