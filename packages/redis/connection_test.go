package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/internal/mock"
)

func newConn(t *testing.T) (*redis.Connection, *mock.Client) {
	t.Helper()
	m := mock.New()

	return redis.NewConnection("default", m), m
}

// ---- strings ---------------------------------------------------------

func TestConnectionGetSet(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	if err := c.Set(ctx, "k", "v", 0); err != nil {
		t.Fatal(err)
	}

	got, err := c.Get(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if got != "v" {
		t.Fatalf("got %q", got)
	}
}

func TestConnectionGetMissing(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	_, err := c.Get(context.Background(), "missing")

	if !errors.Is(err, redis.ErrNil) {
		t.Fatalf("expected ErrNil, got %v", err)
	}
}

func TestConnectionSetWithExpiration(t *testing.T) {
	t.Parallel()
	c, m := newConn(t)
	ctx := context.Background()

	fake := time.Unix(1_700_000_000, 0)
	m.SetClock(func() time.Time { return fake })

	_ = c.Set(ctx, "k", "v", 10*time.Millisecond)
	fake = fake.Add(100 * time.Millisecond)
	_, err := c.Get(ctx, "k")

	if !errors.Is(err, redis.ErrNil) {
		t.Fatalf("expected ErrNil after TTL, got %v", err)
	}
}

func TestConnectionSetNX(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	ok, err := c.SetNX(ctx, "k", "v1", 0)

	if err != nil || !ok {
		t.Fatalf("first SetNX should succeed, got ok=%v err=%v", ok, err)
	}

	ok, err = c.SetNX(ctx, "k", "v2", 0)

	if err != nil || ok {
		t.Fatalf("second SetNX should fail, got ok=%v err=%v", ok, err)
	}
}

func TestConnectionMGet(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()
	_ = c.Set(ctx, "a", "1", 0)
	_ = c.Set(ctx, "b", "2", 0)
	vals, err := c.MGet(ctx, "a", "b", "missing")

	if err != nil {
		t.Fatal(err)
	}

	if len(vals) != 3 || vals[0] != "1" || vals[1] != "2" || vals[2] != nil {
		t.Fatalf("unexpected MGET reply: %+v", vals)
	}
}

func TestConnectionIncrIncrByDecr(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	n, _ := c.Incr(ctx, "counter")

	if n != 1 {
		t.Fatalf("Incr want 1 got %d", n)
	}

	n, _ = c.IncrBy(ctx, "counter", 9)

	if n != 10 {
		t.Fatalf("IncrBy want 10 got %d", n)
	}

	n, _ = c.Decr(ctx, "counter")

	if n != 9 {
		t.Fatalf("Decr want 9 got %d", n)
	}
}

func TestConnectionDelExists(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()
	_ = c.Set(ctx, "a", "1", 0)
	n, _ := c.Exists(ctx, "a", "b")

	if n != 1 {
		t.Fatalf("Exists want 1 got %d", n)
	}

	d, _ := c.Del(ctx, "a")

	if d != 1 {
		t.Fatalf("Del want 1 got %d", d)
	}
}

// ---- hashes ----------------------------------------------------------

func TestConnectionHashRoundTrip(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	if _, err := c.HSet(ctx, "h", "f1", "v1", "f2", "v2"); err != nil {
		t.Fatal(err)
	}

	v, err := c.HGet(ctx, "h", "f1")

	if err != nil || v != "v1" {
		t.Fatalf("HGet=%q err=%v", v, err)
	}

	vals, err := c.HMGet(ctx, "h", "f1", "missing")

	if err != nil {
		t.Fatal(err)
	}

	if vals[0] != "v1" || vals[1] != nil {
		t.Fatalf("HMGET %+v", vals)
	}

	all, _ := c.HGetAll(ctx, "h")

	if all["f2"] != "v2" {
		t.Fatalf("HGETALL %+v", all)
	}

	ok, _ := c.HSetNX(ctx, "h", "f1", "nope")

	if ok {
		t.Fatal("HSETNX should fail on existing field")
	}

	n, _ := c.HDel(ctx, "h", "f1")

	if n != 1 {
		t.Fatalf("HDEL want 1 got %d", n)
	}
}

// ---- lists -----------------------------------------------------------

func TestConnectionListRoundTrip(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	_, _ = c.RPush(ctx, "l", "a", "b", "c")
	_, _ = c.LPush(ctx, "l", "z")
	vals, _ := c.LRange(ctx, "l", 0, -1)

	if len(vals) != 4 || vals[0] != "z" || vals[3] != "c" {
		t.Fatalf("LRange %+v", vals)
	}

	v, _ := c.LPop(ctx, "l")

	if v != "z" {
		t.Fatalf("LPop got %q", v)
	}

	v, _ = c.RPop(ctx, "l")

	if v != "c" {
		t.Fatalf("RPop got %q", v)
	}

	n, _ := c.LRem(ctx, "l", 0, "a")

	if n != 1 {
		t.Fatalf("LRem want 1 got %d", n)
	}
}

// ---- sets ------------------------------------------------------------

func TestConnectionSetRoundTrip(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	_, _ = c.SAdd(ctx, "s", "a", "b", "c")
	ok, _ := c.SIsMember(ctx, "s", "b")

	if !ok {
		t.Fatal("b should be member")
	}

	members, _ := c.SMembers(ctx, "s")

	if len(members) != 3 {
		t.Fatalf("want 3 members, got %+v", members)
	}

	n, _ := c.SRem(ctx, "s", "a")

	if n != 1 {
		t.Fatal("SRem should remove 1")
	}
}

// ---- sorted sets -----------------------------------------------------

func TestConnectionZAddRange(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()
	_, _ = c.ZAdd(ctx, "z",
		redis.ZMember{Score: 3, Member: "c"},
		redis.ZMember{Score: 1, Member: "a"},
		redis.ZMember{Score: 2, Member: "b"},
	)
	vals, _ := c.ZRange(ctx, "z", 0, -1)

	if len(vals) != 3 || vals[0] != "a" || vals[2] != "c" {
		t.Fatalf("ZRange %+v", vals)
	}
}

// ---- server + raw ----------------------------------------------------

func TestConnectionPingFlushDB(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	if err := c.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	_ = c.Set(ctx, "k", "v", 0)

	if err := c.FlushDB(ctx); err != nil {
		t.Fatal(err)
	}

	n, _ := c.Exists(ctx, "k")

	if n != 0 {
		t.Fatalf("expected empty after FlushDB, got Exists=%d", n)
	}
}

func TestConnectionExecuteRaw(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()
	_, err := c.ExecuteRaw(ctx, []any{"SET", "k", "v"})

	if err != nil {
		t.Fatal(err)
	}

	v, _ := c.Get(ctx, "k")

	if v != "v" {
		t.Fatalf("got %q", v)
	}
}
