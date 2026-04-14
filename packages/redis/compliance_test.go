package redis_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/internal/mock"
)

// This file is the Laravel-parity ledger. Each subtest name maps to a
// test case in laravel/framework 13.x (tests/Redis or tests/Integration/
// Redis). The subtests assert Go-side parity using the in-memory fake;
// the integration build tag re-runs a superset against a real Redis.
func TestCompliance_IlluminateRedis(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// RedisConnectionTest.php :: testItGetsAndSetsKeys
	t.Run("RedisConnectionTest/GetsAndSetsKeys", func(t *testing.T) {
		c := redis.NewConnection("default", mock.New())
		if err := c.Set(ctx, "name", "taylor", 0); err != nil {
			t.Fatal(err)
		}
		if v, _ := c.Get(ctx, "name"); v != "taylor" {
			t.Fatalf("got %q", v)
		}
	})

	// RedisConnectionTest.php :: testItDeletesKeys
	t.Run("RedisConnectionTest/DeletesKeys", func(t *testing.T) {
		c := redis.NewConnection("default", mock.New())
		_ = c.Set(ctx, "a", "1", 0)
		if n, _ := c.Del(ctx, "a"); n != 1 {
			t.Fatalf("Del=%d", n)
		}
	})

	// RedisConnectionTest.php :: testItReturnsCommandResultsForIncr
	t.Run("RedisConnectionTest/IncrementsAndDecrements", func(t *testing.T) {
		c := redis.NewConnection("default", mock.New())
		_, _ = c.Incr(ctx, "counter")
		n, _ := c.IncrBy(ctx, "counter", 4)
		if n != 5 {
			t.Fatalf("IncrBy=%d", n)
		}
		n, _ = c.Decr(ctx, "counter")
		if n != 4 {
			t.Fatalf("Decr=%d", n)
		}
	})

	// RedisConnectionTest.php :: testItHandlesHashes
	t.Run("RedisConnectionTest/HashOperations", func(t *testing.T) {
		c := redis.NewConnection("default", mock.New())
		_, _ = c.HSet(ctx, "h", "name", "taylor")
		v, _ := c.HGet(ctx, "h", "name")
		if v != "taylor" {
			t.Fatalf("HGet=%q", v)
		}
	})

	// RedisConnectionTest.php :: testItHandlesLists
	t.Run("RedisConnectionTest/ListOperations", func(t *testing.T) {
		c := redis.NewConnection("default", mock.New())
		_, _ = c.RPush(ctx, "l", "a", "b", "c")
		vals, _ := c.LRange(ctx, "l", 0, -1)
		if len(vals) != 3 {
			t.Fatalf("LRange=%+v", vals)
		}
	})

	// RedisEventsTest.php :: testCommandExecutedEventIsFired
	t.Run("RedisEventsTest/CommandExecutedFired", func(t *testing.T) {
		c := redis.NewConnection("default", mock.New())
		var n int
		c.Listen(func(redis.CommandExecuted) { n++ })
		_, _ = c.Get(ctx, "k")
		if n != 1 {
			t.Fatalf("events fired=%d", n)
		}
	})

	// RedisManagerExtensionTest.php :: testExtensionIsRegistered
	t.Run("RedisManagerExtensionTest/ExtensionIsRegistered", func(t *testing.T) {
		m := redis.NewManager("p", map[string]redis.ConnectionConfig{"p": {}})
		m.Extend("default", func(redis.ConnectionConfig) (redis.Client, error) {
			return mock.New(), nil
		})
		if _, err := m.Connection("p"); err != nil {
			t.Fatal(err)
		}
	})

	// RedisConnectionTest.php :: testItPipelines
	t.Run("RedisConnectionTest/Pipelines", func(t *testing.T) {
		c := redis.NewConnection("default", mock.New())
		cmds, err := c.Pipeline(ctx, func(p redis.Pipeliner) error {
			p.Do(ctx, "SET", "k", "v")
			p.Do(ctx, "GET", "k")
			return nil
		})
		if err != nil || len(cmds) != 2 {
			t.Fatalf("Pipeline err=%v len=%d", err, len(cmds))
		}
	})
}
