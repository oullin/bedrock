package events

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// --- testBasicEventExecution ---

func TestBasicEventExecution(t *testing.T) {
	t.Parallel()

	d := New()

	var called bool
	d.Listen([]string{"user.created"}, func(_ context.Context, _ string, _ any) (any, error) {
		called = true
		return nil, nil
	})

	if err := d.Dispatch(context.Background(), "user.created", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatal("listener was not called")
	}
}

// --- testListenMultipleEvents ---

func TestListenMultipleEvents(t *testing.T) {
	t.Parallel()

	d := New()

	var count int
	d.Listen([]string{"user.created", "user.deleted"}, func(_ context.Context, _ string, _ any) (any, error) {
		count++
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "user.created", nil)
	_ = d.Dispatch(context.Background(), "user.deleted", nil)

	if count != 2 {
		t.Fatalf("want 2 calls, got %d", count)
	}
}

// --- testHaltingEventExecution ---

func TestHaltingEventExecution(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return "stopped", nil
	})

	var secondCalled bool
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		secondCalled = true
		return nil, nil
	})

	result, err := d.Until(context.Background(), "evt", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "stopped" {
		t.Fatalf("want %q, got %v", "stopped", result)
	}

	if secondCalled {
		t.Fatal("second listener should not have been called")
	}
}

// --- testResponseWhenNoListenersAreSet ---

func TestResponseWhenNoListenersAreSet(t *testing.T) {
	t.Parallel()

	d := New()

	if err := d.Dispatch(context.Background(), "nothing", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := d.Until(context.Background(), "nothing", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Fatalf("want nil, got %v", result)
	}
}

// --- testReturningFalseStopsPropagation ---

func TestReturningFalseStopsPropagation(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return false, nil
	})

	var secondCalled bool
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		secondCalled = true
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "evt", nil)

	if secondCalled {
		t.Fatal("returning false should stop propagation")
	}
}

// --- testReturningFalsyValuesContinuesPropagation ---

func TestReturningFalsyValuesContinuesPropagation(t *testing.T) {
	t.Parallel()

	d := New()

	var count int

	// Returning 0 (not bool false) should not stop propagation.
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		count++
		return 0, nil
	})

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		count++
		return "", nil
	})

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		count++
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "evt", nil)

	if count != 3 {
		t.Fatalf("falsy non-bool values should not stop propagation, want 3, got %d", count)
	}
}

// --- testQueuedEventsAreFired ---

func TestQueuedEventsAreFired(t *testing.T) {
	t.Parallel()

	d := New()

	var called bool
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		called = true
		return nil, nil
	})

	d.Push("evt", "payload")

	if called {
		t.Fatal("listener should not be called before FlushQueued")
	}

	_ = d.FlushQueued(context.Background(), "evt")

	if !called {
		t.Fatal("listener should be called after FlushQueued")
	}
}

// --- testQueuedEventsCanBeForgotten ---

func TestQueuedEventsCanBeForgotten(t *testing.T) {
	t.Parallel()

	d := New()

	var called bool
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		called = true
		return nil, nil
	})

	d.Push("evt", nil)
	d.Flush()

	_ = d.FlushQueued(context.Background(), "evt")

	if called {
		t.Fatal("flushed queued events should not fire")
	}
}

// --- testMultiplePushedEventsWillGetFlushed ---

func TestMultiplePushedEventsWillGetFlushed(t *testing.T) {
	t.Parallel()

	d := New()

	var order []int
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, payload any) (any, error) {
		order = append(order, payload.(int))
		return nil, nil
	})

	d.Push("evt", 1)
	d.Push("evt", 2)
	d.Push("evt", 3)

	_ = d.FlushQueued(context.Background(), "evt")

	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("want [1, 2, 3], got %v", order)
	}
}

// --- testPushMethodCanAcceptObjectAsPayload ---

func TestPushMethodCanAcceptObjectAsPayload(t *testing.T) {
	t.Parallel()

	d := New()

	type user struct{ Name string }

	var received any
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, payload any) (any, error) {
		received = payload
		return nil, nil
	})

	d.Push("evt", user{Name: "Alice"})
	_ = d.FlushQueued(context.Background(), "evt")

	u, ok := received.(user)
	if !ok {
		t.Fatalf("want user, got %T", received)
	}

	if u.Name != "Alice" {
		t.Fatalf("want Alice, got %s", u.Name)
	}
}

// --- testWildcardListeners ---

func TestWildcardListeners(t *testing.T) {
	t.Parallel()

	d := New()

	var received string
	d.Listen([]string{"user.*"}, func(_ context.Context, event string, _ any) (any, error) {
		received = event
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "user.created", nil)

	if received != "user.created" {
		t.Fatalf("want %q, got %q", "user.created", received)
	}
}

// --- testWildcardListenersWithResponses ---

func TestWildcardListenersWithResponses(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"user.*"}, func(_ context.Context, _ string, _ any) (any, error) {
		return "response", nil
	})

	result, err := d.Until(context.Background(), "user.created", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "response" {
		t.Fatalf("want %q, got %v", "response", result)
	}
}

// --- testListenersCanBeRemoved ---

func TestListenersCanBeRemoved(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})

	d.Forget("evt")

	if d.HasListeners("evt") {
		t.Fatal("listeners should be removed after Forget")
	}
}

// --- testWildcardListenersCanBeRemoved ---

func TestWildcardListenersCanBeRemoved(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"user.*"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})

	d.Forget("user.*")

	if d.HasListeners("user.created") {
		t.Fatal("wildcard listeners should be removed after Forget")
	}
}

// --- testHasWildcardListeners ---

func TestHasWildcardListeners(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"user.*"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})

	if !d.HasListeners("user.created") {
		t.Fatal("wildcard should match user.created")
	}

	if d.HasListeners("order.created") {
		t.Fatal("wildcard should not match order.created")
	}
}

// --- testListenersCanBeFound ---

func TestListenersCanBeFound(t *testing.T) {
	t.Parallel()

	d := New()

	if d.HasListeners("evt") {
		t.Fatal("should not have listeners before registration")
	}

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})

	if !d.HasListeners("evt") {
		t.Fatal("should have listeners after registration")
	}
}

// --- testWildcardListenersCanBeFound ---

func TestWildcardListenersCanBeFound(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"user.*"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})

	if !d.HasListeners("user.created") {
		t.Fatal("wildcard should match user.created")
	}
}

// --- testEventPassedFirstToWildcards ---

func TestEventPassedFirstToWildcards(t *testing.T) {
	t.Parallel()

	d := New()

	var receivedEvent string
	d.Listen([]string{"user.*"}, func(_ context.Context, event string, _ any) (any, error) {
		receivedEvent = event
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "user.updated", nil)

	if receivedEvent != "user.updated" {
		t.Fatalf("want %q, got %q", "user.updated", receivedEvent)
	}
}

// --- testWildcardMatchAll ---

func TestWildcardMatchAll(t *testing.T) {
	t.Parallel()

	d := New()

	var count int
	d.Listen([]string{"*"}, func(_ context.Context, _ string, _ any) (any, error) {
		count++
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "anything", nil)
	_ = d.Dispatch(context.Background(), "user.created", nil)

	if count != 2 {
		t.Fatalf("want 2 calls, got %d", count)
	}
}

// --- testEventClassesArePayload ---

func TestEventPayloadPassedToListener(t *testing.T) {
	t.Parallel()

	d := New()

	type UserCreated struct{ ID string }

	var received any
	d.Listen([]string{"user.created"}, func(_ context.Context, _ string, payload any) (any, error) {
		received = payload
		return nil, nil
	})

	evt := UserCreated{ID: "42"}
	_ = d.Dispatch(context.Background(), "user.created", evt)

	u, ok := received.(UserCreated)
	if !ok {
		t.Fatalf("want UserCreated, got %T", received)
	}

	if u.ID != "42" {
		t.Fatalf("want ID %q, got %q", "42", u.ID)
	}
}

// --- testNestedEvent ---

func TestNestedEvent(t *testing.T) {
	t.Parallel()

	d := New()

	var nestedCalled bool
	d.Listen([]string{"outer"}, func(ctx context.Context, _ string, _ any) (any, error) {
		d.Listen([]string{"inner"}, func(_ context.Context, _ string, _ any) (any, error) {
			nestedCalled = true
			return nil, nil
		})
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "outer", nil)
	_ = d.Dispatch(context.Background(), "inner", nil)

	if !nestedCalled {
		t.Fatal("listener registered during dispatch should be callable")
	}
}

// --- testDuplicateListenersWillFire ---

func TestDuplicateListenersWillFire(t *testing.T) {
	t.Parallel()

	d := New()

	var count int
	listener := Listener(func(_ context.Context, _ string, _ any) (any, error) {
		count++
		return nil, nil
	})

	d.Listen([]string{"evt"}, listener)
	d.Listen([]string{"evt"}, listener)

	_ = d.Dispatch(context.Background(), "evt", nil)

	if count != 2 {
		t.Fatalf("duplicate listeners should both fire, want 2, got %d", count)
	}
}

// --- testGetListeners ---

func TestGetListeners(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})
	d.Listen([]string{"evt.*"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})

	listeners := d.GetListeners("evt")
	if len(listeners) != 2 {
		t.Fatalf("want 2 exact listeners, got %d", len(listeners))
	}

	listeners = d.GetListeners("evt.sub")
	if len(listeners) != 1 {
		t.Fatalf("want 1 wildcard listener for evt.sub, got %d", len(listeners))
	}
}

// --- testDispatchMultipleListenersInOrder ---

func TestDispatchMultipleListenersInOrder(t *testing.T) {
	t.Parallel()

	d := New()

	var order []int
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		order = append(order, 1)
		return nil, nil
	})
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		order = append(order, 2)
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "evt", nil)

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("want [1, 2], got %v", order)
	}
}

// --- testDispatchWithError ---

func TestDispatchWithError(t *testing.T) {
	t.Parallel()

	d := New()

	sentinel := errors.New("boom")
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, sentinel
	})

	err := d.Dispatch(context.Background(), "evt", nil)
	if !errors.Is(err, sentinel) {
		t.Fatalf("want sentinel error, got %v", err)
	}
}

// --- testDispatchErrorShortCircuits ---

func TestDispatchErrorShortCircuits(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, errors.New("fail")
	})

	var secondCalled bool
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		secondCalled = true
		return nil, nil
	})

	_ = d.Dispatch(context.Background(), "evt", nil)

	if secondCalled {
		t.Fatal("second listener should not have been called after error")
	}
}

// --- testUntilNoHalt ---

func TestUntilNoHalt(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})

	result, err := d.Until(context.Background(), "evt", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Fatalf("want nil, got %v", result)
	}
}

// --- testSubscriber ---

func TestSubscriber(t *testing.T) {
	t.Parallel()

	d := New()

	sub := &testSubscriber{}
	d.Subscribe(sub)

	if !d.HasListeners("sub.event") {
		t.Fatal("subscriber should have registered listeners")
	}
}

type testSubscriber struct{}

func (s *testSubscriber) Subscribe(d Dispatcher) {
	d.Listen([]string{"sub.event"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})
}

// --- testSubscribeCanAcceptObject ---

func TestSubscribeCanAcceptObject(t *testing.T) {
	t.Parallel()

	d := New()

	var executed bool
	sub := &funcSubscriber{fn: func(d Dispatcher) {
		d.Listen([]string{"my.event"}, func(_ context.Context, _ string, _ any) (any, error) {
			executed = true
			return nil, nil
		})
	}}

	d.Subscribe(sub)
	_ = d.Dispatch(context.Background(), "my.event", nil)

	if !executed {
		t.Fatal("subscriber listener should have executed")
	}
}

type funcSubscriber struct {
	fn func(Dispatcher)
}

func (s *funcSubscriber) Subscribe(d Dispatcher) { s.fn(d) }

// --- testFlush ---

func TestFlush(t *testing.T) {
	t.Parallel()

	d := New()

	d.Listen([]string{"a", "b.*"}, func(_ context.Context, _ string, _ any) (any, error) {
		return nil, nil
	})

	d.Flush()

	if d.HasListeners("a") || d.HasListeners("b.c") {
		t.Fatal("all listeners should be removed after Flush")
	}
}

// --- testContextCancellation ---

func TestContextCancellation(t *testing.T) {
	t.Parallel()

	d := New()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var called bool
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		called = true
		return nil, nil
	})

	err := d.Dispatch(ctx, "evt", nil)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if called {
		t.Fatal("listener should not have been called with cancelled context")
	}
}

// --- testConcurrentDispatch ---

func TestConcurrentDispatch(t *testing.T) {
	t.Parallel()

	d := New()

	var count atomic.Int64
	d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
		count.Add(1)
		return nil, nil
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = d.Dispatch(context.Background(), "evt", nil)
		}()
	}

	wg.Wait()

	if count.Load() != 100 {
		t.Fatalf("want 100 calls, got %d", count.Load())
	}
}

// --- testConcurrentListenAndDispatch ---

func TestConcurrentListenAndDispatch(t *testing.T) {
	t.Parallel()

	d := New()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Listen([]string{"evt"}, func(_ context.Context, _ string, _ any) (any, error) {
				return nil, nil
			})
		}()
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = d.Dispatch(context.Background(), "evt", nil)
		}()
	}

	wg.Wait()
}

// --- testFlushQueudSelectivelyByEvent ---

func TestFlushQueuedSelectively(t *testing.T) {
	t.Parallel()

	d := New()

	var aCalled, bCalled bool
	d.Listen([]string{"a"}, func(_ context.Context, _ string, _ any) (any, error) {
		aCalled = true
		return nil, nil
	})
	d.Listen([]string{"b"}, func(_ context.Context, _ string, _ any) (any, error) {
		bCalled = true
		return nil, nil
	})

	d.Push("a", nil)
	d.Push("b", nil)

	_ = d.FlushQueued(context.Background(), "a")

	if !aCalled {
		t.Fatal("event 'a' should have been dispatched")
	}

	if bCalled {
		t.Fatal("event 'b' should not have been dispatched yet")
	}
}

// --- testFlushQueuedAll ---

func TestFlushQueuedAll(t *testing.T) {
	t.Parallel()

	d := New()

	var count int
	d.Listen([]string{"a"}, func(_ context.Context, _ string, _ any) (any, error) {
		count++
		return nil, nil
	})
	d.Listen([]string{"b"}, func(_ context.Context, _ string, _ any) (any, error) {
		count++
		return nil, nil
	})

	d.Push("a", nil)
	d.Push("b", nil)

	_ = d.FlushQueued(context.Background(), "*")

	if count != 2 {
		t.Fatalf("want 2, got %d", count)
	}
}
