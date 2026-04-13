package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// newTestJob creates a BaseJob-backed Job via the SyncDriver for testing.
func newTestJob() queue.Job {
	var captured queue.Job

	handler := queue.HandlerFunc(func(_ context.Context, job queue.Job) error {
		captured = job

		return nil
	})

	drv := drivers.NewSyncDriver("test", handler)
	_, _ = drv.Push(context.Background(), "q", []byte("payload"))

	return captured
}

func TestBaseJobDefaultValues(t *testing.T) {
	t.Parallel()

	job := newTestJob()

	if job.UUID() != "" {
		t.Errorf("expected empty UUID, got %q", job.UUID())
	}

	if job.GetJobID() != "" {
		t.Errorf("expected empty job ID, got %q", job.GetJobID())
	}

	if job.Attempts() != 0 {
		t.Errorf("expected 0 attempts, got %d", job.Attempts())
	}

	if job.MaxTries() != 0 {
		t.Errorf("expected 0 max tries, got %d", job.MaxTries())
	}

	if job.MaxExceptions() != 0 {
		t.Errorf("expected 0 max exceptions, got %d", job.MaxExceptions())
	}

	if job.Timeout() != 0 {
		t.Errorf("expected 0 timeout, got %v", job.Timeout())
	}

	if len(job.Backoff()) != 0 {
		t.Errorf("expected empty backoff, got %v", job.Backoff())
	}

	if job.RetryUntil() != nil {
		t.Errorf("expected nil retry until, got %v", job.RetryUntil())
	}

	if job.IsDeleted() {
		t.Error("expected not deleted")
	}

	if job.IsReleased() {
		t.Error("expected not released")
	}

	if job.HasFailed() {
		t.Error("expected not failed")
	}
}

func TestBaseJobPayload(t *testing.T) {
	t.Parallel()

	job := newTestJob()

	if string(job.Payload()) != "payload" {
		t.Errorf("expected 'payload', got %q", job.Payload())
	}
}

func TestBaseJobGetQueue(t *testing.T) {
	t.Parallel()

	job := newTestJob()

	if job.GetQueue() != "q" {
		t.Errorf("expected 'q', got %q", job.GetQueue())
	}
}

func TestBaseJobGetConnectionName(t *testing.T) {
	t.Parallel()

	job := newTestJob()

	if job.GetConnectionName() != "test" {
		t.Errorf("expected 'test', got %q", job.GetConnectionName())
	}
}

func TestBaseJobFireWithFunc(t *testing.T) {
	t.Parallel()

	fired := false
	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		fired = true

		return nil
	})

	drv := drivers.NewSyncDriver("test", handler)
	_, _ = drv.Push(context.Background(), "q", []byte("p"))

	if !fired {
		t.Error("expected fire to be called")
	}
}

func TestBaseJobReleaseSetsFlag(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "q", []byte("p"))

	job, _ := drv.Pop(context.Background(), "q")

	if job.IsReleased() {
		t.Error("expected not released before Release()")
	}

	_ = job.Release(0)

	if !job.IsReleased() {
		t.Error("expected released after Release()")
	}
}

func TestBaseJobDeleteSetsFlag(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "q", []byte("p"))

	job, _ := drv.Pop(context.Background(), "q")

	if job.IsDeleted() {
		t.Error("expected not deleted before Delete()")
	}

	_ = job.Delete()

	if !job.IsDeleted() {
		t.Error("expected deleted after Delete()")
	}
}

func TestBaseJobFailSetsFlag(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "q", []byte("p"))

	job, _ := drv.Pop(context.Background(), "q")

	if job.HasFailed() {
		t.Error("expected not failed before Fail()")
	}

	_ = job.Fail(errors.New("err"))

	if !job.HasFailed() {
		t.Error("expected failed after Fail()")
	}
}

func TestBaseJobMarkAsFailedIsAlias(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "q", []byte("p"))

	job, _ := drv.Pop(context.Background(), "q")

	err := job.MarkAsFailed(errors.New("err"))

	if err != nil {
		t.Fatal(err)
	}

	if !job.HasFailed() {
		t.Error("expected failed after MarkAsFailed()")
	}
}

func TestBaseJobFireWithoutFunc(t *testing.T) {
	t.Parallel()

	// Pop from redis gives a job without fireFunc set.
	client := newMockRedisClient()
	drv := drivers.NewRedisDriver(client, "redis")

	_, _ = drv.Push(context.Background(), "q", []byte("p"))

	job, _ := drv.Pop(context.Background(), "q")

	err := job.Fire(context.Background())

	if err != nil {
		t.Errorf("expected nil error from Fire without func, got %v", err)
	}
}

func TestBaseJobReleaseWithoutFunc(t *testing.T) {
	t.Parallel()

	// SyncDriver pops return ErrNoJob, but we can test via the null driver approach.
	// Create a job that has no releaseFunc by using the sync driver's captured job.
	job := newTestJob()

	err := job.Release(time.Second)

	if err != nil {
		t.Errorf("expected nil error from Release without func, got %v", err)
	}
}
