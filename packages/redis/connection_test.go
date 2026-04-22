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
// RedisConnectionTest::testItSetsValuesWithExpiry
// RedisConnectionTest::testItDeletesKeys
// RedisConnectionTest::testItChecksForExistence
// RedisConnectionTest::testItExpiresKeys
// RedisConnectionTest::testItSetsKeyIfNotExists

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

// RedisConnectionTest::testItSetsKeyIfNotExists
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

// RedisConnectionTest::testItGetsMultipleKeys
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
// RedisConnectionTest::testItSetsMultipleHashFields
// RedisConnectionTest::testItGetsMultipleHashFields
// RedisConnectionTest::testItSetsHashFieldIfNotExists

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
// RedisConnectionTest::testItAddsMembersToSortedSet
// RedisConnectionTest::testItCountsMembersInSortedSet
// RedisConnectionTest::testItIncrementsScoreOfSortedSet
// RedisConnectionTest::testItReturnsRangeInSortedSet
// RedisConnectionTest::testItReturnsRevRangeInSortedSet
// RedisConnectionTest::testItReturnsRangeByScoreInSortedSet
// RedisConnectionTest::testItReturnsRevRangeByScoreInSortedSet
// RedisConnectionTest::testItReturnsRankInSortedSet
// RedisConnectionTest::testItReturnsScoreInSortedSet
// RedisConnectionTest::testItCalculatesIntersectionOfSortedSetsAndStores
// RedisConnectionTest::testItCalculatesUnionOfSortedSetsAndStores
// RedisConnectionTest::testItRemovesMembersInSortedSet
// RedisConnectionTest::testItRemovesMembersByScoreInSortedSet
// RedisConnectionTest::testItRemovesMembersByRankInSortedSet

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

	if n, _ := c.ZCard(ctx, "z"); n != 3 {
		t.Fatalf("ZCard want 3 got %d", n)
	}

	if score, _ := c.ZIncrBy(ctx, "z", 3, "a"); score != 4 {
		t.Fatalf("ZIncrBy want 4 got %v", score)
	}

	rev, _ := c.ZRevRange(ctx, "z", 0, -1)

	if len(rev) != 3 || rev[0] != "a" || rev[2] != "b" {
		t.Fatalf("ZRevRange %+v", rev)
	}
}

func TestConnectionZSetMutators(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	_, _ = c.ZAdd(ctx, "z",
		redis.ZMember{Score: 1, Member: "a"},
		redis.ZMember{Score: 2, Member: "b"},
		redis.ZMember{Score: 3, Member: "c"},
	)

	rank, err := c.ZRank(ctx, "z", "b")
	if err != nil {
		t.Fatal(err)
	}

	if rank != 1 {
		t.Fatalf("ZRank want 1 got %d", rank)
	}

	score, err := c.ZScore(ctx, "z", "c")
	if err != nil {
		t.Fatal(err)
	}

	if score != 3 {
		t.Fatalf("ZScore want 3 got %v", score)
	}

	if n, _ := c.ZRem(ctx, "z", "a"); n != 1 {
		t.Fatalf("ZRem want 1 got %d", n)
	}

	if n, _ := c.ZRemRangeByScore(ctx, "z", "2", "2"); n != 1 {
		t.Fatalf("ZRemRangeByScore want 1 got %d", n)
	}

	_, _ = c.ZAdd(ctx, "z",
		redis.ZMember{Score: 1, Member: "a"},
		redis.ZMember{Score: 2, Member: "b"},
		redis.ZMember{Score: 3, Member: "c"},
	)

	if n, _ := c.ZRemRangeByRank(ctx, "z", 0, 1); n != 2 {
		t.Fatalf("ZRemRangeByRank want 2 got %d", n)
	}

	if vals, _ := c.ZRange(ctx, "z", 0, -1); len(vals) != 1 || vals[0] != "c" {
		t.Fatalf("ZRange after removals %+v", vals)
	}
}

func TestConnectionZRangeByScoreAndStore(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	_, _ = c.ZAdd(ctx, "src",
		redis.ZMember{Score: 1, Member: "a"},
		redis.ZMember{Score: 2, Member: "b"},
		redis.ZMember{Score: 3, Member: "c"},
	)

	ranged, err := c.ZRangeByScore(ctx, "src", "2", "3")
	if err != nil {
		t.Fatal(err)
	}

	if len(ranged) != 2 || ranged[0] != "b" || ranged[1] != "c" {
		t.Fatalf("ZRangeByScore %+v", ranged)
	}

	rev, err := c.ZRevRangeByScore(ctx, "src", "3", "2")
	if err != nil {
		t.Fatal(err)
	}

	if len(rev) != 2 || rev[0] != "c" || rev[1] != "b" {
		t.Fatalf("ZRevRangeByScore %+v", rev)
	}

	if n, err := c.ZInterStore(ctx, "dest-inter", "src"); err != nil || n != 3 {
		t.Fatalf("ZInterStore n=%d err=%v", n, err)
	}

	if n, err := c.ZUnionStore(ctx, "dest-union", "src"); err != nil || n != 3 {
		t.Fatalf("ZUnionStore n=%d err=%v", n, err)
	}
}

// ---- server + raw ----------------------------------------------------
// RedisConnectionTest::testItFlushes
// RedisConnectionTest::testItFlushesAsynchronous
// RedisConnectionTest::testItRunsRawCommand

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

	if err := c.FlushAllAsync(ctx); err != nil {
		t.Fatal(err)
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

// RedisConnectionTest::testItSPopsForKeys
// RedisConnectionTest::testItScansForKeys
// RedisConnectionTest::testItHscansForKeys
// RedisConnectionTest::testItSscansForKeys
// RedisConnectionTest::testItZscansForKeys
// RedisConnectionTest::testPhpRedisScanOption
func TestConnectionPopAndScanFamilies(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	_, _ = c.SAdd(ctx, "set", "a", "b", "c")
	popped, err := c.SPop(ctx, "set", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(popped) != 2 {
		t.Fatalf("SPop %+v", popped)
	}

	_ = c.Set(ctx, "user:1", "a", 0)
	_ = c.Set(ctx, "user:2", "b", 0)
	scanned, err := c.Scan(ctx, 0, "user:*", 10)
	if err != nil {
		t.Fatal(err)
	}
	if scanned.Cursor != 0 || len(scanned.Values) != 2 {
		t.Fatalf("Scan %+v", scanned)
	}

	_, _ = c.HSet(ctx, "hash", "f1", "v1", "f2", "v2")
	hscan, err := c.HScan(ctx, "hash", 0, "*", 10)
	if err != nil {
		t.Fatal(err)
	}
	if hscan.Cursor != 0 || len(hscan.Values) != 4 {
		t.Fatalf("HScan %+v", hscan)
	}

	_, _ = c.SAdd(ctx, "scan-set", "x", "y")
	sscan, err := c.SScan(ctx, "scan-set", 0, "*", 10)
	if err != nil {
		t.Fatal(err)
	}
	if sscan.Cursor != 0 || len(sscan.Values) != 2 {
		t.Fatalf("SScan %+v", sscan)
	}

	_, _ = c.ZAdd(ctx, "scan-zset", redis.ZMember{Score: 1, Member: "m1"}, redis.ZMember{Score: 2, Member: "m2"})
	zscan, err := c.ZScan(ctx, "scan-zset", 0, "*", 10)
	if err != nil {
		t.Fatal(err)
	}
	if zscan.Cursor != 0 || len(zscan.Values) != 4 {
		t.Fatalf("ZScan %+v", zscan)
	}
}

// RedisConnectionTest::testItRenamesKeys
// RedisConnectionTest::testItPersistsConnection
func TestConnectionRenameAndPersist(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	_ = c.Set(ctx, "old", "v", time.Second)

	if err := c.Rename(ctx, "old", "new"); err != nil {
		t.Fatal(err)
	}

	v, err := c.Get(ctx, "new")
	if err != nil || v != "v" {
		t.Fatalf("Rename result=%q err=%v", v, err)
	}

	ok, err := c.Persist(ctx, "new")
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("Persist should report a removed TTL")
	}
}

// RedisConnectionTest::testItRunsEval
func TestConnectionEvalReturnsValue(t *testing.T) {
	t.Parallel()
	c, _ := newConn(t)
	ctx := context.Background()

	v, err := c.Eval(ctx, "return 1", nil)
	if err != nil {
		t.Fatal(err)
	}

	if n, ok := v.(int64); !ok || n != 1 {
		t.Fatalf("Eval=%T %#v", v, v)
	}
}

func TestHasHashTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{name: "plain key", key: "queues:default", want: false},
		{name: "queue key with tag", key: "queues:{default}", want: true},
		{name: "tagged queue derivative", key: "queues:{default}:reserved", want: true},
		{name: "empty tag", key: "queues:{}", want: false},
		{name: "missing close brace", key: "queues:{default", want: false},
		{name: "missing open brace", key: "queues:default}", want: false},
		{name: "empty first tag disables hash tag", key: "queues:{}:{default}", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := redis.HasHashTag(tt.key); got != tt.want {
				t.Fatalf("HasHashTag(%q)=%v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
