package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestRedisDriverPush(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}

	n, _ := client.LLen(context.Background(), "queues:default")

	if n != 1 {
		t.Errorf("expected 1 item in list, got %d", n)
	}
}

func TestRedisDriverPushDelayed(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, err := drv.PushDelayed(context.Background(), "default", []byte("payload"), 5*time.Second)

	if err != nil {
		t.Fatal(err)
	}

	n, _ := client.ZCard(context.Background(), "queues:default:delayed")

	if n != 1 {
		t.Errorf("expected 1 delayed, got %d", n)
	}
}

func TestRedisDriverPushMultiple(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	ids, err := drv.PushMultiple(context.Background(), "default", [][]byte{
		[]byte("a"), []byte("b"), []byte("c"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %d", len(ids))
	}

	n, _ := client.LLen(context.Background(), "queues:default")

	if n != 3 {
		t.Errorf("expected 3 items, got %d", n)
	}
}

func TestRedisDriverPopFromMainQueue(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "default", []byte("test-payload"))

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if string(job.Payload()) != "test-payload" {
		t.Errorf("expected 'test-payload', got %q", job.Payload())
	}

	if job.GetQueue() != "default" {
		t.Errorf("expected queue 'default', got %q", job.GetQueue())
	}

	if job.GetConnectionName() != "redis" {
		t.Errorf("expected connection 'redis', got %q", job.GetConnectionName())
	}
}

func TestRedisDriverPopMigratesDueDelayedJobs(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	// Add a delayed job with a past score (already due).
	_ = client.ZAdd(context.Background(), "queues:default:delayed", float64(time.Now().Add(-time.Minute).Unix()), "due-payload")

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if string(job.Payload()) != "due-payload" {
		t.Errorf("expected 'due-payload', got %q", job.Payload())
	}

	// Delayed set should be empty after migration.
	n, _ := client.ZCard(context.Background(), "queues:default:delayed")

	if n != 0 {
		t.Errorf("expected 0 delayed after migration, got %d", n)
	}
}

func TestRedisDriverPopEmptyQueueReturnsErrNoJob(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, err := drv.Pop(context.Background(), "empty")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestRedisDriverJobRelease(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "default", []byte("release-me"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Release(0)

	if err != nil {
		t.Fatal(err)
	}

	// Should be back on the main queue.
	n, _ := client.LLen(context.Background(), "queues:default")

	if n != 1 {
		t.Errorf("expected 1 item after release, got %d", n)
	}
}

func TestRedisDriverJobReleaseWithDelay(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "default", []byte("delay-me"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Release(5 * time.Second)

	if err != nil {
		t.Fatal(err)
	}

	// Should be in the delayed set.
	n, _ := client.ZCard(context.Background(), "queues:default:delayed")

	if n != 1 {
		t.Errorf("expected 1 delayed after release, got %d", n)
	}
}

func TestRedisDriverJobDelete(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "default", []byte("delete-me"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Delete()

	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisDriverJobFail(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "default", []byte("fail-me"))

	job, _ := drv.Pop(context.Background(), "default")

	err := job.Fail(errors.New("something broke"))

	if err != nil {
		t.Fatal(err)
	}

	// Should have an entry in the failed key.
	n, _ := client.LLen(context.Background(), "queues:default:failed")

	if n != 1 {
		t.Errorf("expected 1 failed entry, got %d", n)
	}
}

func TestRedisDriverSize(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "default", []byte("a"))
	_, _ = drv.Push(context.Background(), "default", []byte("b"))

	n, err := drv.Size(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 2 {
		t.Errorf("expected size 2, got %d", n)
	}
}

func TestRedisDriverPendingSize(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "default", []byte("a"))

	n, err := drv.PendingSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Errorf("expected pending size 1, got %d", n)
	}
}

func TestRedisDriverDelayedSize(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.PushDelayed(context.Background(), "default", []byte("a"), time.Hour)

	n, err := drv.DelayedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Errorf("expected delayed size 1, got %d", n)
	}
}

func TestRedisDriverReservedSizeAlwaysZero(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	n, err := drv.ReservedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestRedisDriverConnectionName(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "my-redis")

	if drv.ConnectionName() != "my-redis" {
		t.Errorf("expected 'my-redis', got %q", drv.ConnectionName())
	}
}
