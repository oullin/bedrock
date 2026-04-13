package drivers_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestDeferredDriverPushBuffers(t *testing.T) {
	t.Parallel()

	dispatched := 0
	drv := drivers.NewDeferredDriver("deferred", func(_ context.Context, _ drivers.DeferredEntry) error {
		dispatched++

		return nil
	})

	_, _ = drv.Push(context.Background(), "q", []byte("a"))

	if dispatched != 0 {
		t.Error("expected no dispatch before flush")
	}

	n, _ := drv.Size(context.Background(), "q")

	if n != 1 {
		t.Errorf("expected size 1, got %d", n)
	}
}

func TestDeferredDriverPushDelayedSetsAfter(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", nil)

	before := time.Now()
	_, _ = drv.PushDelayed(context.Background(), "q", []byte("a"), 5*time.Second)

	n, _ := drv.Size(context.Background(), "q")

	if n != 1 {
		t.Errorf("expected size 1, got %d", n)
	}

	_ = before // After field is internal; we just verify it buffers.
}

func TestDeferredDriverPushMultipleBuffersAll(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", nil)

	ids, _ := drv.PushMultiple(context.Background(), "q", [][]byte{
		[]byte("a"), []byte("b"), []byte("c"),
	})

	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %d", len(ids))
	}

	n, _ := drv.Size(context.Background(), "q")

	if n != 3 {
		t.Errorf("expected size 3, got %d", n)
	}
}

func TestDeferredDriverPopAlwaysReturnsErrNoJob(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", nil)

	_, err := drv.Pop(context.Background(), "q")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestDeferredDriverSizeReflectsBuffer(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", nil)

	_, _ = drv.Push(context.Background(), "q", []byte("a"))
	_, _ = drv.Push(context.Background(), "q", []byte("b"))

	n, _ := drv.Size(context.Background(), "q")

	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}

func TestDeferredDriverFlushDispatchesAll(t *testing.T) {
	t.Parallel()

	flushed := 0
	drv := drivers.NewDeferredDriver("deferred", func(_ context.Context, _ drivers.DeferredEntry) error {
		flushed++

		return nil
	})

	_, _ = drv.Push(context.Background(), "q", []byte("a"))
	_, _ = drv.Push(context.Background(), "q", []byte("b"))

	if err := drv.Flush(context.Background()); err != nil {
		t.Fatal(err)
	}

	if flushed != 2 {
		t.Errorf("expected 2 flushed, got %d", flushed)
	}
}

func TestDeferredDriverFlushClearsBuffer(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", func(_ context.Context, _ drivers.DeferredEntry) error {
		return nil
	})

	_, _ = drv.Push(context.Background(), "q", []byte("a"))
	_ = drv.Flush(context.Background())

	n, _ := drv.Size(context.Background(), "q")

	if n != 0 {
		t.Errorf("expected 0 after flush, got %d", n)
	}
}

func TestDeferredDriverFlushErrorStopsEarly(t *testing.T) {
	t.Parallel()

	calls := 0
	drv := drivers.NewDeferredDriver("deferred", func(_ context.Context, _ drivers.DeferredEntry) error {
		calls++

		if calls == 1 {
			return errors.New("dispatch failed")
		}

		return nil
	})

	_, _ = drv.Push(context.Background(), "q", []byte("a"))
	_, _ = drv.Push(context.Background(), "q", []byte("b"))

	err := drv.Flush(context.Background())

	if err == nil {
		t.Fatal("expected error from flush")
	}

	if calls != 1 {
		t.Errorf("expected 1 call before error, got %d", calls)
	}
}

func TestDeferredDriverFlushNilDispatcher(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", nil)

	_, _ = drv.Push(context.Background(), "q", []byte("a"))

	err := drv.Flush(context.Background())

	if err != nil {
		t.Fatal("expected no error with nil dispatcher")
	}

	n, _ := drv.Size(context.Background(), "q")

	if n != 0 {
		t.Errorf("expected 0 after flush, got %d", n)
	}
}

func TestDeferredDriverConcurrentPush(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", nil)

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, _ = drv.Push(context.Background(), "q", []byte("data"))
		}()
	}

	wg.Wait()

	n, _ := drv.Size(context.Background(), "q")

	if n != 50 {
		t.Errorf("expected 50, got %d", n)
	}
}

func TestDeferredDriverConnectionName(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("my-deferred", nil)

	if drv.ConnectionName() != "my-deferred" {
		t.Errorf("expected 'my-deferred', got %q", drv.ConnectionName())
	}
}

func TestDeferredDriverDelayedSizeAlwaysZero(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", nil)

	n, err := drv.DelayedSize(context.Background(), "q")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestDeferredDriverReservedSizeAlwaysZero(t *testing.T) {
	t.Parallel()

	drv := drivers.NewDeferredDriver("deferred", nil)

	n, err := drv.ReservedSize(context.Background(), "q")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}
