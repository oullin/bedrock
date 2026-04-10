package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestSyncDriverProcessesImmediately(t *testing.T) {
	processed := 0
	handler := queue.HandlerFunc(func(_ context.Context, job queue.Job) error {
		processed++

		return nil
	})

	drv := drivers.NewSyncDriver("default", handler)
	_, err := drv.Push(context.Background(), "default", []byte(`{"job":"test"}`))

	if err != nil {
		t.Fatal(err)
	}

	if processed != 1 {
		t.Errorf("expected 1 processed job, got %d", processed)
	}
}

func TestNullDriverDiscardsJobs(t *testing.T) {
	drv := drivers.NewNullDriver("null")
	id, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}

	_ = id

	n, _ := drv.Size(context.Background(), "default")

	if n != 0 {
		t.Errorf("null driver should always report size 0, got %d", n)
	}
}

func TestWorkerStopsOnEmpty(t *testing.T) {
	drv := drivers.NewNullDriver("default")
	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })

	w := queue.NewWorker(drv, handler, nil, queue.WorkerOptions{
		StopOnEmpty: true,
		Sleep:       time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	if err := w.Run(ctx, "default"); err != nil {
		t.Errorf("worker returned error: %v", err)
	}
}

func TestFailoverDriverFallsBack(t *testing.T) {
	// A null driver always succeeds (even though it discards). FailoverDriver
	// should return the first successful push.
	d1 := drivers.NewNullDriver("d1")
	d2 := drivers.NewNullDriver("d2")
	drv := drivers.NewFailoverDriver("failover", d1, d2)

	_, err := drv.Push(context.Background(), "q", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestDeferredDriverFlush(t *testing.T) {
	flushed := 0
	drv := drivers.NewDeferredDriver("deferred", func(_ context.Context, _ drivers.DeferredEntry) error {
		flushed++

		return nil
	})

	_, _ = drv.Push(context.Background(), "q", []byte("a"))
	_, _ = drv.Push(context.Background(), "q", []byte("b"))

	n, _ := drv.Size(context.Background(), "q")

	if n != 2 {
		t.Errorf("expected 2 deferred, got %d", n)
	}

	if err := drv.Flush(context.Background()); err != nil {
		t.Fatal(err)
	}

	if flushed != 2 {
		t.Errorf("expected 2 flushed, got %d", flushed)
	}

	n2, _ := drv.Size(context.Background(), "q")

	if n2 != 0 {
		t.Errorf("expected 0 after flush, got %d", n2)
	}
}
