package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// Add a delayed job with a past score (already due).

// Delayed set should be empty after migration.

// Should be back on the main queue.

// Should be in the delayed set.

// Should have an entry in the failed key.

// The default mockRedisClient does not satisfy RedisScanner,
// RedisListRanger, or RedisSortedSetRanger.

// ReservedJobs always returns ErrNotSupported for the Redis driver.

// rangerRedisClient extends mockRedisClient with the optional Scanner /
// ListRanger / SortedSetRanger contracts used by the Redis driver's
// inspection methods.
type rangerRedisClient struct {
	*mockRedisClient
	keys      []string
	scanErr   error
	lrangeErr error
	zrangeErr error
}

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

	_ = client.ZAdd(context.Background(), "queues:default:delayed", float64(time.Now().Add(-time.Minute).Unix()), "due-payload")

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if string(job.Payload()) != "due-payload" {
		t.Errorf("expected 'due-payload', got %q", job.Payload())
	}

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

func TestRedisDriverInspectionWithoutCapabilityReturnsErrNotSupported(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")
	ctx := context.Background()

	if _, err := drv.QueueNames(ctx); !errors.Is(err, queue.ErrNotSupported) {
		t.Errorf("QueueNames: want ErrNotSupported, got %v", err)
	}

	if _, err := drv.PendingJobs(ctx, "default"); !errors.Is(err, queue.ErrNotSupported) {
		t.Errorf("PendingJobs: want ErrNotSupported, got %v", err)
	}

	if _, err := drv.DelayedJobs(ctx, "default"); !errors.Is(err, queue.ErrNotSupported) {
		t.Errorf("DelayedJobs: want ErrNotSupported, got %v", err)
	}

	if _, err := drv.ReservedJobs(ctx, "default"); !errors.Is(err, queue.ErrNotSupported) {
		t.Errorf("ReservedJobs: want ErrNotSupported, got %v", err)
	}
}

func (c *rangerRedisClient) ScanMatch(_ context.Context, _ string) ([]string, error) {
	if c.scanErr != nil {
		return nil, c.scanErr
	}

	return c.keys, nil
}

func (c *rangerRedisClient) LRange(_ context.Context, key string, _, _ int64) ([]string, error) {
	if c.lrangeErr != nil {
		return nil, c.lrangeErr
	}

	return c.lists[key], nil
}

func (c *rangerRedisClient) ZRange(_ context.Context, key string, _, _ int64) ([]string, error) {
	if c.zrangeErr != nil {
		return nil, c.zrangeErr
	}

	entries := c.sorted[key]
	out := make([]string, 0, len(entries))

	for _, e := range entries {
		out = append(out, e.member)
	}

	return out, nil
}

func TestRedisDriverQueueNamesDedupesAndUnwraps(t *testing.T) {
	t.Parallel()

	client := &rangerRedisClient{
		mockRedisClient: newMockRedisClient(),
		keys: []string{
			"queues:default",
			"queues:default:delayed",
			"queues:emails",
			"queues:{cluster}",
			"queues:", // empty after trim — must be skipped
		},
	}
	drv := drivers.NewRedisDriver(client, "redis")

	names, err := drv.QueueNames(context.Background())

	if err != nil {
		t.Fatalf("QueueNames: %v", err)
	}

	want := map[string]struct{}{"default": {}, "emails": {}, "cluster": {}}

	if len(names) != len(want) {
		t.Fatalf("got %v, want %d unique names", names, len(want))
	}

	for _, n := range names {
		if _, ok := want[n]; !ok {
			t.Errorf("unexpected name %q in %v", n, names)
		}
	}
}

func TestRedisDriverQueueNamesScannerErrorPropagates(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("scan boom")
	client := &rangerRedisClient{
		mockRedisClient: newMockRedisClient(),
		scanErr:         wantErr,
	}
	drv := drivers.NewRedisDriver(client, "redis")

	_, err := drv.QueueNames(context.Background())

	if !errors.Is(err, wantErr) {
		t.Errorf("want %v, got %v", wantErr, err)
	}
}

func TestRedisDriverPendingJobsReturnsSnapshots(t *testing.T) {
	t.Parallel()

	client := &rangerRedisClient{mockRedisClient: newMockRedisClient()}
	drv := drivers.NewRedisDriver(client, "redis")
	ctx := context.Background()

	// Push two real payloads via the driver so the queue key matches.
	_, _ = drv.Push(ctx, "default", []byte(`{"uuid":"u1","displayName":"Job1","createdAt":1000000}`))
	_, _ = drv.Push(ctx, "default", []byte(`not-json`))

	jobs, err := drv.PendingJobs(ctx, "default")

	if err != nil {
		t.Fatalf("PendingJobs: %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("got %d jobs, want 2", len(jobs))
	}

	// Find the JSON one and assert decoded fields.
	var decoded *queue.InspectedJob

	for i := range jobs {
		if jobs[i].UUID == "u1" {
			decoded = &jobs[i]

			break
		}
	}

	if decoded == nil {
		t.Fatal("expected a job with UUID u1")
	}

	if decoded.Name != "Job1" || decoded.Connection != "redis" || decoded.CreatedAt.Unix() != 1000000 {
		t.Errorf("decoded mismatch: %+v", decoded)
	}
}

func TestRedisDriverDelayedJobsReturnsSnapshots(t *testing.T) {
	t.Parallel()

	client := &rangerRedisClient{mockRedisClient: newMockRedisClient()}
	drv := drivers.NewRedisDriver(client, "redis")
	ctx := context.Background()

	_, _ = drv.PushDelayed(ctx, "default", []byte(`{"uuid":"d1"}`), 24*time.Hour)

	jobs, err := drv.DelayedJobs(ctx, "default")

	if err != nil {
		t.Fatalf("DelayedJobs: %v", err)
	}

	if len(jobs) != 1 || jobs[0].UUID != "d1" || jobs[0].Queue != "default" {
		t.Errorf("got %+v, want 1 job with UUID d1 on default", jobs)
	}
}

func TestRedisDriverInspectionRangerErrorPropagates(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("range boom")
	client := &rangerRedisClient{
		mockRedisClient: newMockRedisClient(),
		lrangeErr:       wantErr,
		zrangeErr:       wantErr,
	}
	drv := drivers.NewRedisDriver(client, "redis")
	ctx := context.Background()

	if _, err := drv.PendingJobs(ctx, "default"); !errors.Is(err, wantErr) {
		t.Errorf("PendingJobs: want %v, got %v", wantErr, err)
	}

	if _, err := drv.DelayedJobs(ctx, "default"); !errors.Is(err, wantErr) {
		t.Errorf("DelayedJobs: want %v, got %v", wantErr, err)
	}
}
