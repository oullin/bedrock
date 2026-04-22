package jobqueue

import (
	"errors"
	"testing"
)

func TestRedisPrefixInventoryParity(t *testing.T) {
	t.Run("RedisPrefixTest::test_prefix_can_be_configured", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Prefix: "custom:"})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if connection.Prefix != "custom:" {
			t.Fatalf("prefix = %q, want custom:", connection.Prefix)
		}
	})

	t.Run("RedisPrefixTest::test_cluster_connection_uses_hash_tagged_prefix", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Prefix: "jobqueue:", Kind: RedisCluster, SupportsCluster: true})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if connection.Prefix != "{jobqueue}:" {
			t.Fatalf("prefix = %q, want {jobqueue}:", connection.Prefix)
		}
	})

	t.Run("RedisPrefixTest::test_cluster_prefix_is_not_double_tagged", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Prefix: "{jobqueue}:", Kind: RedisCluster, SupportsCluster: true})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if connection.Prefix != "{jobqueue}:" {
			t.Fatalf("prefix = %q, want unchanged hash tag", connection.Prefix)
		}
	})

	t.Run("RedisPrefixTest::test_standalone_connection_prefix_is_unchanged", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Prefix: "jobqueue:", Kind: RedisStandalone})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if connection.Prefix != "jobqueue:" {
			t.Fatalf("prefix = %q, want jobqueue:", connection.Prefix)
		}
	})

	t.Run("RedisPrefixTest::test_cluster_connection_uses_fallback_prefix", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("fallback:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Kind: RedisCluster, SupportsCluster: true})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if connection.Prefix != "{fallback}:" {
			t.Fatalf("prefix = %q, want {fallback}:", connection.Prefix)
		}
	})

	t.Run("RedisPrefixTest::test_cluster_connection_preserves_nodes", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Kind: RedisCluster, SupportsCluster: true, Nodes: []string{"127.0.0.1:6379", "127.0.0.1:6380"}})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if len(connection.Nodes) != 2 || connection.Nodes[1] != "127.0.0.1:6380" {
			t.Fatalf("nodes = %#v", connection.Nodes)
		}
	})

	t.Run("RedisPrefixTest::test_use_throws_for_unknown_connection", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")

		_, err := registry.Use("missing")
		if !errors.Is(err, ErrUnknownRedisConnection) {
			t.Fatalf("Use error = %v, want ErrUnknownRedisConnection", err)
		}
	})

	t.Run("RedisPrefixTest::test_cluster_connection_falls_back_to_standalone_when_unsupported", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Kind: RedisCluster, SupportsCluster: false})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if connection.Kind != RedisStandalone {
			t.Fatalf("kind = %q, want standalone", connection.Kind)
		}
	})

	t.Run("RedisPrefixTest::test_cluster_connection_registers_under_clusters_not_standalone", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Kind: RedisCluster, SupportsCluster: true})

		if !registry.Cluster("jobqueue") || registry.Standalone("jobqueue") {
			t.Fatalf("cluster = %v, standalone = %v", registry.Cluster("jobqueue"), registry.Standalone("jobqueue"))
		}
	})

	t.Run("RedisPrefixTest::test_standalone_connection_does_not_register_under_clusters", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Kind: RedisStandalone})

		if registry.Cluster("jobqueue") || !registry.Standalone("jobqueue") {
			t.Fatalf("cluster = %v, standalone = %v", registry.Cluster("jobqueue"), registry.Standalone("jobqueue"))
		}
	})

	t.Run("RedisPrefixTest::test_cluster_connection_preserves_existing_options", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Kind: RedisCluster, SupportsCluster: true, Options: map[string]string{"read_timeout": "1s"}})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if connection.Options["read_timeout"] != "1s" {
			t.Fatalf("options = %#v", connection.Options)
		}

		connection.Options["read_timeout"] = "mutated"
		again, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if again.Options["read_timeout"] != "1s" {
			t.Fatalf("options were not cloned: %#v", again.Options)
		}
	})

	t.Run("RedisPrefixTest::test_cluster_takes_precedence_over_standalone_with_same_name", func(t *testing.T) {
		registry := NewRedisConnectionRegistry("jobqueue:")
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Prefix: "standalone:", Kind: RedisStandalone})
		registry.Register(RedisConnectionConfig{Name: "jobqueue", Prefix: "cluster:", Kind: RedisCluster, SupportsCluster: true})

		connection, err := registry.Use("jobqueue")
		if err != nil {
			t.Fatalf("Use returned error: %v", err)
		}
		if connection.Kind != RedisCluster || connection.Prefix != "{cluster}:" {
			t.Fatalf("connection = %#v, want cluster precedence", connection)
		}
	})
}
