package bus_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/bus"
)

type testCommand struct{ Value string }

func TestDispatcherSyncDispatch(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)
	d.Map(testCommand{}, func(ctx context.Context, cmd any) (any, error) {
		c := cmd.(testCommand)

		return c.Value + "_handled", nil
	})

	result, err := d.Dispatch(context.Background(), testCommand{Value: "hello"})
	if err != nil {
		t.Fatal(err)
	}

	if result != "hello_handled" {
		t.Errorf("got %v, want hello_handled", result)
	}
}

func TestDispatcherPipelineMiddleware(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	order := []string{}
	d.PipeThrough(
		func(ctx context.Context, cmd any, next bus.Handler) (any, error) {
			order = append(order, "pipe1_before")
			result, err := next(ctx, cmd)
			order = append(order, "pipe1_after")

			return result, err
		},
		func(ctx context.Context, cmd any, next bus.Handler) (any, error) {
			order = append(order, "pipe2_before")
			result, err := next(ctx, cmd)
			order = append(order, "pipe2_after")

			return result, err
		},
	)

	d.Map(testCommand{}, func(ctx context.Context, _ any) (any, error) {
		order = append(order, "handler")

		return nil, nil
	})

	_, _ = d.Dispatch(context.Background(), testCommand{})

	expected := []string{"pipe1_before", "pipe2_before", "handler", "pipe2_after", "pipe1_after"}
	for i, e := range expected {
		if i >= len(order) || order[i] != e {
			t.Errorf("pipeline order[%d]: got %q, want %q", i, func() string {
				if i < len(order) {
					return order[i]
				}

				return "<missing>"
			}(), e)
		}
	}
}

func TestDispatcherDeferredFlush(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	count := 0
	d.Map(testCommand{}, func(_ context.Context, _ any) (any, error) {
		count++

		return nil, nil
	})

	_ = d.DispatchAfterResponse(context.Background(), testCommand{})
	_ = d.DispatchAfterResponse(context.Background(), testCommand{})

	if count != 0 {
		t.Error("deferred commands should not run before Flush")
	}

	if err := d.FlushDeferred(context.Background()); err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Errorf("expected 2 deferred commands run, got %d", count)
	}
}

func TestDispatcherNoHandlerError(t *testing.T) {
	d := bus.NewDispatcher(nil, nil)

	_, err := d.Dispatch(context.Background(), testCommand{})
	if err == nil {
		t.Error("expected error for unregistered command")
	}
}
