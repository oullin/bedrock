//go:build integration

// This file is compiled only with `-tags=integration`. It exercises the
// package against a real Redis instance given by the REDIS_URL env var.
// It is the Go analogue of tests/Integration/Redis/PredisConnectionTest.php.
package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/limiters"
)

func newIntegrationConn(t *testing.T) *redis.Connection {
	t.Helper()
	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL not set; skipping integration test")
	}
	opts, err := goredis.ParseURL(url)
	if err != nil {
		t.Fatalf("ParseURL: %v", err)
	}
	c := redis.NewConnection("integration", redis.NewGoRedisClient(goredis.NewClient(opts)))
	t.Cleanup(func() { _ = c.Close() })
	_ = c.FlushDB(context.Background())
	return c
}

func TestIntegration_StringRoundTrip(t *testing.T) {
	c := newIntegrationConn(t)
	ctx := context.Background()

	if err := c.Set(ctx, "k", "hello", time.Second); err != nil {
		t.Fatal(err)
	}
	if v, _ := c.Get(ctx, "k"); v != "hello" {
		t.Fatalf("Get=%q", v)
	}
}

func TestIntegration_Pipeline(t *testing.T) {
	c := newIntegrationConn(t)
	ctx := context.Background()
	cmds, err := c.Pipeline(ctx, func(p redis.Pipeliner) error {
		p.Do(ctx, "SET", "k1", "1")
		p.Do(ctx, "SET", "k2", "2")
		p.Do(ctx, "GET", "k1")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := cmds[2].Result(); v != "1" {
		t.Fatalf("GET result=%v", v)
	}
}

func TestIntegration_PubSub(t *testing.T) {
	c := newIntegrationConn(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan string, 1)
	go func() {
		_ = c.Subscribe(ctx, []string{"ch"}, func(_, payload string) {
			select {
			case done <- payload:
			default:
			}
			cancel()
		})
	}()

	// Give the subscriber a moment to attach.
	time.Sleep(100 * time.Millisecond)
	pub := newIntegrationConn(t)
	if _, err := pub.Command(context.Background(), "PUBLISH", "ch", "hi"); err != nil {
		t.Fatal(err)
	}

	select {
	case p := <-done:
		if p != "hi" {
			t.Fatalf("payload=%q", p)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no message received")
	}
}

func TestIntegration_ConcurrencyLimiter(t *testing.T) {
	c := newIntegrationConn(t)
	err := limiters.NewConcurrencyBuilder(c, "job").
		Limit(2).
		ReleaseAfter(2 * time.Second).
		Block(500 * time.Millisecond).
		Sleep(25 * time.Millisecond).
		Then(context.Background(), func() error { return nil }, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIntegration_DurationLimiter(t *testing.T) {
	c := newIntegrationConn(t)
	ctx := context.Background()
	lim := limiters.NewDurationLimiter(c, "rate", 2, 2*time.Second)
	for i := 0; i < 2; i++ {
		ok, err := lim.Acquire(ctx)
		if err != nil || !ok {
			t.Fatalf("attempt %d ok=%v err=%v", i, ok, err)
		}
	}
	ok, _ := lim.Acquire(ctx)
	if ok {
		t.Fatal("3rd acquire should fail")
	}
}
