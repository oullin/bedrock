package redis_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/internal/mock"
)

func TestManagerExtendAndConnection(t *testing.T) {
	t.Parallel()

	m := redis.NewManager("primary", map[string]redis.ConnectionConfig{
		"primary": {Name: "primary"},
	})

	// Custom driver that returns a fake client regardless of config.
	m.Extend("default", func(cfg redis.ConnectionConfig) (redis.Client, error) {
		return mock.New(), nil
	})

	c, err := m.Connection("primary")

	if err != nil {
		t.Fatal(err)
	}

	if c.Name() != "primary" {
		t.Fatalf("name=%q", c.Name())
	}

	// Cached on second access.
	c2, _ := m.Connection("primary")

	if c != c2 {
		t.Fatal("expected cached connection identity")
	}

	// Purge forces rebuild.
	m.Purge("primary")
	c3, _ := m.Connection("primary")

	if c == c3 {
		t.Fatal("expected new connection after Purge")
	}
}

func TestManagerUnknownConnection(t *testing.T) {
	t.Parallel()
	m := redis.NewManager("x", nil)
	_, err := m.Connection("nope")

	if !errors.Is(err, redis.ErrConnectionNotFound) {
		t.Fatalf("want ErrConnectionNotFound, got %v", err)
	}
}

func TestManagerEventsPropagateToConnections(t *testing.T) {
	t.Parallel()

	m := redis.NewManager("p", map[string]redis.ConnectionConfig{"p": {}})
	m.Extend("default", func(redis.ConnectionConfig) (redis.Client, error) { return mock.New(), nil })

	var n int

	m.Listen(func(redis.CommandExecuted) { n++ })

	c, _ := m.Connection("p")
	_ = c.Set(context.Background(), "k", "v", 0)

	if n == 0 {
		t.Fatal("manager listener did not receive event")
	}
}

func TestManagerRegisterReusesConnection(t *testing.T) {
	t.Parallel()
	m := redis.NewManager("p", nil)
	pre := redis.NewConnection("p", mock.New())
	m.Register("p", pre)

	got, err := m.Connection("p")

	if err != nil {
		t.Fatal(err)
	}

	if got != pre {
		t.Fatal("Register did not round-trip")
	}
}
