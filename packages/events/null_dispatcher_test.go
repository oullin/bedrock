package events_test

import (
	"context"
	"testing"

	cevents "github.com/bedrock/packages/contracts/events"
	"github.com/bedrock/packages/events"
)

// Push via null should not store anything.

// Flush on inner should have nothing to dispatch.

// nullTestSubscriber is a test subscriber for NullDispatcher tests.
type nullTestSubscriber struct{}

func TestNullDispatcher_DispatchIsNoOp(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)
	ctx := context.Background()

	var called bool

	inner.Listen("e", func(ctx context.Context, event any) (any, error) {
		called = true

		return nil, nil
	})

	responses, err := null.Dispatch(ctx, "e")

	if err != nil {
		t.Fatal(err)
	}

	if len(responses) != 0 {
		t.Fatalf("expected no responses, got %d", len(responses))
	}

	if called {
		t.Fatal("listener should not be called via NullDispatcher")
	}
}

func TestNullDispatcher_UntilIsNoOp(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)
	ctx := context.Background()

	inner.Listen("e", func(ctx context.Context, event any) (any, error) {
		return "should-not-reach", nil
	})

	result, err := null.Until(ctx, "e")

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestNullDispatcher_PushIsNoOp(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)
	ctx := context.Background()

	null.Push(ctx, "e")

	var called bool

	inner.Listen("e", func(ctx context.Context, event any) (any, error) {
		called = true

		return nil, nil
	})

	err := inner.Flush(ctx, "e")

	if err != nil {
		t.Fatal(err)
	}

	if called {
		t.Fatal("nothing should have been pushed")
	}
}

func TestNullDispatcher_FlushIsNoOp(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)
	ctx := context.Background()

	err := null.Flush(ctx, "e")

	if err != nil {
		t.Fatal(err)
	}
}

func TestNullDispatcher_ListenDelegates(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)

	null.Listen("e", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	if !inner.HasListeners("e") {
		t.Fatal("expected listener to be registered on inner dispatcher")
	}
}

func TestNullDispatcher_HasListenersDelegates(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)

	inner.Listen("e", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	if !null.HasListeners("e") {
		t.Fatal("expected HasListeners to delegate to inner")
	}
}

func TestNullDispatcher_HasWildcardListenersDelegates(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)

	inner.Listen("order.*", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	if !null.HasWildcardListeners("order.created") {
		t.Fatal("expected HasWildcardListeners to delegate to inner")
	}
}

func TestNullDispatcher_ForgetDelegates(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)

	inner.Listen("e", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	null.Forget("e")

	if inner.HasListeners("e") {
		t.Fatal("expected Forget to remove listeners from inner")
	}
}

func TestNullDispatcher_ForgetPushedDelegates(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)
	ctx := context.Background()

	inner.Push(ctx, "e")
	null.ForgetPushed()

	var called bool

	inner.Listen("e", func(ctx context.Context, event any) (any, error) {
		called = true

		return nil, nil
	})

	inner.Flush(ctx, "e")

	if called {
		t.Fatal("expected ForgetPushed to clear pushed events on inner")
	}
}

func TestNullDispatcher_GetListenersDelegates(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)

	inner.Listen("e", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	ls := null.GetListeners("e")

	if len(ls) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(ls))
	}
}

func TestNullDispatcher_SubscribeDelegates(t *testing.T) {
	t.Parallel()

	inner := events.NewDispatcher()
	null := events.NewNullDispatcher(inner)

	sub := &nullTestSubscriber{}
	null.Subscribe(sub)

	if !inner.HasListeners("sub.event") {
		t.Fatal("expected subscriber to register listeners on inner dispatcher")
	}
}

func (s *nullTestSubscriber) Subscribe(d cevents.Dispatcher) {
	d.Listen("sub.event", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})
}
