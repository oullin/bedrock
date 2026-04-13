package drivers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestSyncDriverPushExecutesImmediately(t *testing.T) {
	t.Parallel()

	processed := 0
	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		processed++

		return nil
	})

	drv := drivers.NewSyncDriver("sync", handler)
	_, err := drv.Push(context.Background(), "default", []byte(`{"job":"test"}`))

	if err != nil {
		t.Fatal(err)
	}

	if processed != 1 {
		t.Errorf("expected 1 processed, got %d", processed)
	}
}

func TestSyncDriverPushMultipleExecutesAll(t *testing.T) {
	t.Parallel()

	processed := 0
	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		processed++

		return nil
	})

	drv := drivers.NewSyncDriver("sync", handler)
	ids, err := drv.PushMultiple(context.Background(), "default", [][]byte{
		[]byte("a"), []byte("b"), []byte("c"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %d", len(ids))
	}

	if processed != 3 {
		t.Errorf("expected 3 processed, got %d", processed)
	}
}

func TestSyncDriverPopAlwaysReturnsErrNoJob(t *testing.T) {
	t.Parallel()

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })
	drv := drivers.NewSyncDriver("sync", handler)

	_, err := drv.Pop(context.Background(), "default")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestSyncDriverSizeAlwaysZero(t *testing.T) {
	t.Parallel()

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })
	drv := drivers.NewSyncDriver("sync", handler)
	ctx := context.Background()

	for _, fn := range []func(context.Context, string) (int64, error){
		drv.Size, drv.PendingSize, drv.DelayedSize, drv.ReservedSize,
	} {
		n, err := fn(ctx, "default")

		if err != nil {
			t.Fatal(err)
		}

		if n != 0 {
			t.Errorf("expected 0, got %d", n)
		}
	}
}

func TestSyncDriverConnectionName(t *testing.T) {
	t.Parallel()

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })
	drv := drivers.NewSyncDriver("my-conn", handler)

	if drv.ConnectionName() != "my-conn" {
		t.Errorf("expected 'my-conn', got %q", drv.ConnectionName())
	}
}

func TestSyncDriverHandlerError(t *testing.T) {
	t.Parallel()

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		return errors.New("handler failed")
	})

	drv := drivers.NewSyncDriver("sync", handler)
	_, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err == nil {
		t.Fatal("expected error from handler")
	}

	if err.Error() != "handler failed" {
		t.Errorf("expected 'handler failed', got %q", err.Error())
	}
}

func TestSyncDriverPushJobHasCorrectPayload(t *testing.T) {
	t.Parallel()

	var captured []byte

	handler := queue.HandlerFunc(func(_ context.Context, job queue.Job) error {
		captured = job.Payload()

		return nil
	})

	drv := drivers.NewSyncDriver("sync", handler)
	payload := []byte(`{"job":"test","data":{"key":"val"}}`)
	_, _ = drv.Push(context.Background(), "q1", payload)

	if string(captured) != string(payload) {
		t.Errorf("expected payload %q, got %q", payload, captured)
	}
}

func TestSyncDriverPushJobHasCorrectQueue(t *testing.T) {
	t.Parallel()

	var queueName string

	handler := queue.HandlerFunc(func(_ context.Context, job queue.Job) error {
		queueName = job.GetQueue()

		return nil
	})

	drv := drivers.NewSyncDriver("sync", handler)
	_, _ = drv.Push(context.Background(), "emails", []byte("p"))

	if queueName != "emails" {
		t.Errorf("expected queue 'emails', got %q", queueName)
	}
}
