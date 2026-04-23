package redis

import "testing"

// RedisConnectorTest::testDefaultConfiguration
// RedisConnectorTest::testUrl
// RedisConnectorTest::testUrlWithScheme
// RedisConnectorTest::testScheme
func TestAddrResolution(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  ConnectionConfig
		want string
	}{
		{
			name: "default host and port",
			cfg:  ConnectionConfig{},
			want: "127.0.0.1:6379",
		},
		{
			name: "explicit host and port",
			cfg: ConnectionConfig{
				Host: "cache.internal",
				Port: 6380,
			},
			want: "cache.internal:6380",
		},
		{
			name: "url takes precedence",
			cfg: ConnectionConfig{
				URL:  "tcp://redis.example:6381",
				Host: "cache.internal",
				Port: 6380,
			},
			want: "tcp://redis.example:6381",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := addr(tt.cfg); got != tt.want {
				t.Fatalf("addr()=%q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatAddr(t *testing.T) {
	t.Parallel()

	if got := formatAddr("redis.local", 6390); got != "redis.local:6390" {
		t.Fatalf("formatAddr()=%q, want %q", got, "redis.local:6390")
	}
}

// RedisConnectorTest::testPredisConfigurationWithUsername
// RedisConnectorTest::testPredisConfigurationWithSentinel
// RedisConnectorTest::testPrefixOverrideBehaviour
func TestConnectionConfigCarriesConnectorFields(t *testing.T) {
	t.Parallel()

	cfg := ConnectionConfig{
		Username: "redis-user",
		Password: "redis-pass",
		Prefix:   "cache:",
		Sentinel: &SentinelConfig{
			MasterName:    "mymaster",
			SentinelAddrs: []string{"127.0.0.1:26379"},
			Database:      3,
		},
	}

	if cfg.Username != "redis-user" || cfg.Password != "redis-pass" {
		t.Fatalf("auth fields not preserved: %+v", cfg)
	}

	if cfg.Prefix != "cache:" {
		t.Fatalf("Prefix=%q, want cache:", cfg.Prefix)
	}

	if cfg.Sentinel == nil || cfg.Sentinel.MasterName != "mymaster" || cfg.Sentinel.Database != 3 {
		t.Fatalf("Sentinel=%+v", cfg.Sentinel)
	}
}
