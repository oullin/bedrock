package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestNullDriverPushDiscardsJob(t *testing.T) {
	t.Parallel()

	drv := drivers.NewNullDriver("null")
	id, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}

	if id != "" {
		t.Errorf("expected empty id, got %q", id)
	}
}

func TestNullDriverPushDelayedDiscardsJob(t *testing.T) {
	t.Parallel()

	drv := drivers.NewNullDriver("null")
	id, err := drv.PushDelayed(context.Background(), "default", []byte("payload"), 5*time.Second)

	if err != nil {
		t.Fatal(err)
	}

	if id != "" {
		t.Errorf("expected empty id, got %q", id)
	}
}

func TestNullDriverPushMultipleReturnsCorrectIDCount(t *testing.T) {
	t.Parallel()

	drv := drivers.NewNullDriver("null")
	ids, err := drv.PushMultiple(context.Background(), "default", [][]byte{
		[]byte("a"), []byte("b"), []byte("c"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %d", len(ids))
	}
}

func TestNullDriverPopReturnsErrNoJob(t *testing.T) {
	t.Parallel()

	drv := drivers.NewNullDriver("null")
	_, err := drv.Pop(context.Background(), "default")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestNullDriverAllSizesReturnZero(t *testing.T) {
	t.Parallel()

	drv := drivers.NewNullDriver("null")
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

func TestNullDriverConnectionName(t *testing.T) {
	t.Parallel()

	drv := drivers.NewNullDriver("my-null")

	if drv.ConnectionName() != "my-null" {
		t.Errorf("expected 'my-null', got %q", drv.ConnectionName())
	}
}
