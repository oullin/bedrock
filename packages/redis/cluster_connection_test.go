package redis_test

import (
	"context"
	"sort"
	"testing"

	"github.com/bedrock/packages/redis"
	"github.com/bedrock/packages/redis/internal/mock"
)

type clusterAwareFake struct {
	*mock.Client
	masters []*mock.Client
	calls   int
}

func newClusterAwareFake(masters ...*mock.Client) *clusterAwareFake {
	return &clusterAwareFake{
		Client:  mock.New(),
		masters: masters,
	}
}

func (c *clusterAwareFake) ForEachMaster(ctx context.Context, fn func(redis.Client) error) error {
	c.calls++

	for _, master := range c.masters {
		if err := fn(master); err != nil {
			return err
		}
	}

	return nil
}

// PhpRedisClusterConnectionTest::testItScansUsingDefaultNode
// PhpRedisClusterConnectionTest::testItOnlyFetchesDefaultNodeOnce
// PhpRedisClusterConnectionTest::testItScansUsingOptionNode
func TestPhpRedisClusterConnectionScansAllMasters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	m1 := mock.New()
	m2 := mock.New()
	c1 := redis.NewConnection("master-1", m1)
	c2 := redis.NewConnection("master-2", m2)

	_ = c1.Set(ctx, "user:1", "a", 0)
	_ = c2.Set(ctx, "user:2", "b", 0)

	fake := newClusterAwareFake(m1, m2)
	conn := redis.NewConnection("cluster", fake)

	res, err := conn.ClusterScan(ctx, 0, "user:*", 10)
	if err != nil {
		t.Fatalf("ClusterScan err=%v", err)
	}

	if res.Cursor != 0 {
		t.Fatalf("Cursor=%d, want 0", res.Cursor)
	}

	sort.Strings(res.Values)

	if len(res.Values) != 2 || res.Values[0] != "user:1" || res.Values[1] != "user:2" {
		t.Fatalf("Values=%v", res.Values)
	}

	if fake.calls != 1 {
		t.Fatalf("ForEachMaster calls=%d, want 1", fake.calls)
	}
}

// PhpRedisClusterConnectionTest::testItReturnsFalseWhenCursorIsZeroAndResultIsEmpty
func TestPhpRedisClusterConnectionReturnsEmptyScanWhenNoMasters(t *testing.T) {
	t.Parallel()

	conn := redis.NewConnection("cluster", newClusterAwareFake())

	res, err := conn.ClusterScan(context.Background(), 0, "missing:*", 10)
	if err != nil {
		t.Fatalf("ClusterScan err=%v", err)
	}

	if res.Cursor != 0 || len(res.Values) != 0 {
		t.Fatalf("Scan result=%+v", res)
	}
}

// PhpRedisClusterConnectionTest::testItFlushesAllMasterNodes
func TestPhpRedisClusterConnectionFlushesAllMasters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	m1 := mock.New()
	m2 := mock.New()
	c1 := redis.NewConnection("master-1", m1)
	c2 := redis.NewConnection("master-2", m2)

	_ = c1.Set(ctx, "alpha", "1", 0)
	_ = c2.Set(ctx, "beta", "2", 0)

	conn := redis.NewConnection("cluster", newClusterAwareFake(m1, m2))
	if err := conn.ClusterFlushDB(ctx); err != nil {
		t.Fatalf("ClusterFlushDB err=%v", err)
	}

	if got := len(m1.Data()); got != 0 {
		t.Fatalf("master1 keys=%d, want 0", got)
	}

	if got := len(m2.Data()); got != 0 {
		t.Fatalf("master2 keys=%d, want 0", got)
	}
}
