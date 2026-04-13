package events_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	cevents "github.com/bedrock/packages/contracts/events"
	"github.com/bedrock/packages/events"
)

// Register listener.

// Dispatch concurrently.

// Concurrently register and dispatch.

// Deferred events are dispatched after the callback.

// Just verify it can be set without panic.

// testSubscriber implements Subscriber for testing.
type testSubscriber struct {
	orderFn events.Listener
	userFn  events.Listener
}

func TestListenAndDispatchStringEvent(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var called bool

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		called = true

		return nil, nil
	})

	_, err := d.Dispatch(ctx, "order.created")

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatal("listener was not called")
	}
}

func TestListenAndDispatchStructEvent(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var received string

	d.Listen(testOrderCreated{}, func(ctx context.Context, event any) (any, error) {
		received = event.(testOrderCreated).OrderID

		return nil, nil
	})

	_, err := d.Dispatch(ctx, testOrderCreated{OrderID: "abc"})

	if err != nil {
		t.Fatal(err)
	}

	if received != "abc" {
		t.Fatalf("expected %q, got %q", "abc", received)
	}
}

func TestDispatchCollectsResponses(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		return "first", nil
	})

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		return "second", nil
	})

	responses, err := d.Dispatch(ctx, "order.created")

	if err != nil {
		t.Fatal(err)
	}

	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}

	if responses[0] != "first" || responses[1] != "second" {
		t.Fatalf("unexpected responses: %v", responses)
	}
}

func TestDispatchNilResponsesAreExcluded(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	d.Listen("e", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	d.Listen("e", func(ctx context.Context, event any) (any, error) {
		return "ok", nil
	})

	responses, err := d.Dispatch(ctx, "e")

	if err != nil {
		t.Fatal(err)
	}

	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(responses))
	}
}

func TestUntilStopsOnFirstNonNilResponse(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var secondCalled bool

	d.Listen("e", func(ctx context.Context, event any) (any, error) {
		return "halted", nil
	})

	d.Listen("e", func(ctx context.Context, event any) (any, error) {
		secondCalled = true

		return "never", nil
	})

	result, err := d.Until(ctx, "e")

	if err != nil {
		t.Fatal(err)
	}

	if result != "halted" {
		t.Fatalf("expected %q, got %v", "halted", result)
	}

	if secondCalled {
		t.Fatal("second listener should not have been called")
	}
}

func TestUntilReturnsNilWhenNoResponse(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	d.Listen("e", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	result, err := d.Until(ctx, "e")

	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestDispatchWithError(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()
	testErr := fmt.Errorf("test error")

	d.Listen("e", func(ctx context.Context, event any) (any, error) {
		return nil, testErr
	})

	_, err := d.Dispatch(ctx, "e")

	if err != testErr {
		t.Fatalf("expected %v, got %v", testErr, err)
	}
}

func TestWildcardListeners(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var calls []string

	d.Listen("order.*", func(ctx context.Context, event any) (any, error) {
		calls = append(calls, "wildcard")

		return nil, nil
	})

	d.Dispatch(ctx, "order.created")
	d.Dispatch(ctx, "order.shipped")
	d.Dispatch(ctx, "user.created")

	if len(calls) != 2 {
		t.Fatalf("expected 2 wildcard calls, got %d", len(calls))
	}
}

func TestWildcardAndDirectListeners(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var calls []string

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		calls = append(calls, "direct")

		return nil, nil
	})

	d.Listen("order.*", func(ctx context.Context, event any) (any, error) {
		calls = append(calls, "wildcard")

		return nil, nil
	})

	d.Dispatch(ctx, "order.created")

	if len(calls) != 2 {
		t.Fatalf("expected 2 calls (direct + wildcard), got %d: %v", len(calls), calls)
	}

	if calls[0] != "direct" || calls[1] != "wildcard" {
		t.Fatalf("expected [direct wildcard], got %v", calls)
	}
}

func TestHasListeners(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()

	if d.HasListeners("order.created") {
		t.Fatal("expected no listeners initially")
	}

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	if !d.HasListeners("order.created") {
		t.Fatal("expected listeners after registration")
	}
}

func TestHasListenersWithWildcard(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()

	d.Listen("order.*", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	if !d.HasListeners("order.created") {
		t.Fatal("expected wildcard to match")
	}

	if d.HasListeners("user.created") {
		t.Fatal("expected no match for unrelated event")
	}
}

func TestHasWildcardListeners(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()

	d.Listen("order.*", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	if !d.HasWildcardListeners("order.created") {
		t.Fatal("expected wildcard match")
	}

	if d.HasWildcardListeners("user.created") {
		t.Fatal("expected no wildcard match")
	}
}

func TestForget(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	d.Forget("order.created")

	if d.HasListeners("order.created") {
		t.Fatal("expected no listeners after Forget")
	}
}

func TestForgetWildcard(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()

	d.Listen("order.*", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	d.Forget("order.*")

	if d.HasWildcardListeners("order.created") {
		t.Fatal("expected no wildcard listeners after Forget")
	}
}

func TestPushAndFlush(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var called bool

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		called = true

		return nil, nil
	})

	d.Push(ctx, "order.created")

	if called {
		t.Fatal("listener should not be called before Flush")
	}

	err := d.Flush(ctx, "order.created")

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatal("listener should be called after Flush")
	}
}

func TestForgetPushed(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var called bool

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		called = true

		return nil, nil
	})

	d.Push(ctx, "order.created")
	d.ForgetPushed()

	err := d.Flush(ctx, "order.created")

	if err != nil {
		t.Fatal(err)
	}

	if called {
		t.Fatal("listener should not be called after ForgetPushed")
	}
}

func TestSubscriber(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var orderCalled, userCalled bool

	sub := &testSubscriber{
		orderFn: func(ctx context.Context, event any) (any, error) {
			orderCalled = true

			return nil, nil
		},
		userFn: func(ctx context.Context, event any) (any, error) {
			userCalled = true

			return nil, nil
		},
	}

	d.Subscribe(sub)

	d.Dispatch(ctx, "order.created")
	d.Dispatch(ctx, "user.registered")

	if !orderCalled {
		t.Fatal("order listener from subscriber was not called")
	}

	if !userCalled {
		t.Fatal("user listener from subscriber was not called")
	}
}

func TestGetListeners(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	d.Listen("order.*", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	ls := d.GetListeners("order.created")

	if len(ls) != 2 {
		t.Fatalf("expected 2 listeners, got %d", len(ls))
	}
}

func TestGetRawListeners(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	d.Listen("user.registered", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	raw := d.GetRawListeners()

	if len(raw) != 2 {
		t.Fatalf("expected 2 event entries, got %d", len(raw))
	}

	if len(raw["order.created"]) != 1 {
		t.Fatalf("expected 1 listener for order.created, got %d", len(raw["order.created"]))
	}
}

func TestListenMultipleEvents(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var count int

	d.Listen([]string{"order.created", "order.shipped"}, func(ctx context.Context, event any) (any, error) {
		count++

		return nil, nil
	})

	d.Dispatch(ctx, "order.created")
	d.Dispatch(ctx, "order.shipped")

	if count != 2 {
		t.Fatalf("expected 2 calls, got %d", count)
	}
}

func TestDispatchNoListeners(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	responses, err := d.Dispatch(ctx, "nonexistent")

	if err != nil {
		t.Fatal(err)
	}

	if len(responses) != 0 {
		t.Fatalf("expected 0 responses, got %d", len(responses))
	}
}

func TestMakeListener(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	fn := func(ctx context.Context, event any) (any, error) {
		return "result", nil
	}

	wrapped := d.MakeListener(fn, false)
	result, err := wrapped(ctx, "test")

	if err != nil {
		t.Fatal(err)
	}

	if result != "result" {
		t.Fatalf("expected %q, got %v", "result", result)
	}
}

func TestMakeListenerWildcard(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	fn := func(ctx context.Context, event any) (any, error) {
		return "wildcard-result", nil
	}

	wrapped := d.MakeListener(fn, true)
	result, err := wrapped(ctx, "test")

	if err != nil {
		t.Fatal(err)
	}

	if result != "wildcard-result" {
		t.Fatalf("expected %q, got %v", "wildcard-result", result)
	}
}

func TestConcurrentListenAndDispatch(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var count atomic.Int64

	var wg sync.WaitGroup

	d.Listen("e", func(ctx context.Context, event any) (any, error) {
		count.Add(1)

		return nil, nil
	})

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			d.Dispatch(ctx, "e")
		}()
	}

	wg.Wait()

	if count.Load() != 100 {
		t.Fatalf("expected 100 calls, got %d", count.Load())
	}
}

func TestConcurrentListenAndDispatchMultipleEvents(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)
		event := fmt.Sprintf("event.%d", i)

		go func() {
			defer wg.Done()

			d.Listen(event, func(ctx context.Context, event any) (any, error) {
				return nil, nil
			})
		}()

		go func() {
			defer wg.Done()

			d.Dispatch(ctx, event)
		}()
	}

	wg.Wait()
}

func TestDefer(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var calls []string

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		calls = append(calls, "order.created")

		return nil, nil
	})

	err := d.Defer(ctx, func(ctx context.Context) error {
		d.Dispatch(ctx, "order.created")
		d.Dispatch(ctx, "order.created")

		return nil
	}, "order.created")

	if err != nil {
		t.Fatal(err)
	}

	if len(calls) != 2 {
		t.Fatalf("expected 2 calls after defer, got %d", len(calls))
	}
}

func TestDeferAllEvents(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var calls []string

	d.Listen("a", func(ctx context.Context, event any) (any, error) {
		calls = append(calls, "a")

		return nil, nil
	})

	d.Listen("b", func(ctx context.Context, event any) (any, error) {
		calls = append(calls, "b")

		return nil, nil
	})

	err := d.Defer(ctx, func(ctx context.Context) error {
		d.Dispatch(ctx, "a")
		d.Dispatch(ctx, "b")

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(calls) != 2 {
		t.Fatalf("expected 2 deferred calls, got %d", len(calls))
	}
}

func TestDeferCallbackError(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()
	testErr := fmt.Errorf("callback error")

	d.Listen("e", func(ctx context.Context, event any) (any, error) {
		return nil, nil
	})

	err := d.Defer(ctx, func(ctx context.Context) error {
		return testErr
	}, "e")

	if err != testErr {
		t.Fatalf("expected %v, got %v", testErr, err)
	}
}

func TestDispatchPointerEvent(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var received string

	d.Listen(testOrderCreated{}, func(ctx context.Context, event any) (any, error) {
		received = event.(*testOrderCreated).OrderID

		return nil, nil
	})

	_, err := d.Dispatch(ctx, &testOrderCreated{OrderID: "ptr-123"})

	if err != nil {
		t.Fatal(err)
	}

	if received != "ptr-123" {
		t.Fatalf("expected %q, got %q", "ptr-123", received)
	}
}

func TestFlushOnlyNamedEvents(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	ctx := context.Background()

	var orderCalled, userCalled bool

	d.Listen("order.created", func(ctx context.Context, event any) (any, error) {
		orderCalled = true

		return nil, nil
	})

	d.Listen("user.registered", func(ctx context.Context, event any) (any, error) {
		userCalled = true

		return nil, nil
	})

	d.Push(ctx, "order.created")
	d.Push(ctx, "user.registered")

	err := d.Flush(ctx, "order.created")

	if err != nil {
		t.Fatal(err)
	}

	if !orderCalled {
		t.Fatal("order listener should have been called")
	}

	if userCalled {
		t.Fatal("user listener should not have been called yet")
	}
}

func TestSetQueueResolver(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	called := false

	d.SetQueueResolver(func() events.QueueBackend {
		called = true

		return nil
	})

	if called {
		t.Fatal("resolver should not be called during set")
	}
}

func TestSetTransactionManagerResolver(t *testing.T) {
	t.Parallel()

	d := events.NewDispatcher()
	called := false

	d.SetTransactionManagerResolver(func() events.TransactionManager {
		called = true

		return nil
	})

	if called {
		t.Fatal("resolver should not be called during set")
	}
}

func (s *testSubscriber) Subscribe(d cevents.Dispatcher) {
	d.Listen("order.created", s.orderFn)
	d.Listen("user.registered", s.userFn)
}
