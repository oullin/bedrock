package redis_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/bedrock/packages/redis"
)

// RedisConnectionTest::testItDispatchesQueryEvent
func TestCommandExecutedDispatchedWhenEnabled(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)

	var mu sync.Mutex

	var events []redis.CommandExecuted

	c.Listen(func(e redis.CommandExecuted) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	})

	ctx := context.Background()
	_ = c.Set(ctx, "k", "v", 0)
	_, _ = c.Get(ctx, "k")
	_, _ = c.Del(ctx, "k")

	mu.Lock()

	defer mu.Unlock()

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d: %+v", len(events), events)
	}

	names := []string{events[0].Command, events[1].Command, events[2].Command}
	want := []string{"set", "get", "del"}

	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("event[%d] command = %q, want %q", i, names[i], want[i])
		}
	}

	for _, e := range events {
		if e.ConnectionName != "default" {
			t.Fatalf("ConnectionName=%q", e.ConnectionName)
		}
	}
}

func TestEventsDisabledByDefault(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)

	var fired int
	// Register listener directly on dispatcher so we don't auto-enable.
	c.Events().Listen(func(e redis.CommandExecuted) { fired++ })

	ctx := context.Background()
	_ = c.Set(ctx, "k", "v", 0)

	if fired != 0 {
		t.Fatalf("expected no events when disabled, fired=%d", fired)
	}

	c.Events().Enable()
	_ = c.Set(ctx, "k", "v", 0)

	if fired == 0 {
		t.Fatal("expected event after Enable")
	}
}

// RedisEventsTest::testCommandFailedEventIsDispatched
// RedisEventsTest::testCommandExecutedEventIsNotDispatchedWhenCommandFails
// RedisEventsTest::testCommandFailedEventContainsConnectionName
// RedisEventsTest::testListenForFailuresRegistersCallback
func TestCommandFailedDispatchedAndListened(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)

	var executed int
	c.Listen(func(redis.CommandExecuted) { executed++ })

	var failed []redis.CommandFailed
	c.ListenForFailures(func(e redis.CommandFailed) {
		failed = append(failed, e)
	})

	ctx := context.Background()
	_, err := c.Get(ctx, "missing")

	if !errors.Is(err, redis.ErrNil) {
		t.Fatalf("expected ErrNil, got %v", err)
	}

	if len(failed) != 1 {
		t.Fatalf("expected one failure event, got %d", len(failed))
	}

	if executed != 0 {
		t.Fatalf("expected no success events, got %d", executed)
	}

	if failed[0].Command != "get" || failed[0].ConnectionName != "default" {
		t.Fatalf("failure event=%+v", failed[0])
	}
}
