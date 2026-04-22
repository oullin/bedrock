package events_test

import (
	"context"
	"errors"
	"testing"
	"time"

	contractevents "github.com/bedrock/packages/contracts/events"
	"github.com/bedrock/packages/events"
)

type inventoryEvent struct {
	ID string
}

type inventoryNestedEvent struct {
	Event inventoryEvent
}

type inventorySubscriber struct {
	order []string
}

func (s *inventorySubscriber) Subscribe(d contractevents.Dispatcher) {
	d.Listen("inventory.created", func(context.Context, any) (any, error) {
		s.order = append(s.order, "created")

		return nil, nil
	})
}

// EventsDispatcherTest::testBasicEventExecution
// EventsDispatcherTest::testDeferEventExecution
// EventsDispatcherTest::testDeferMultipleEvents
// EventsDispatcherTest::testDeferNestedEvents
// EventsDispatcherTest::testDeferSpecificEvents
// EventsDispatcherTest::testDeferSpecificNestedEvents
// EventsDispatcherTest::testDeferSpecificObjectEvents
// EventsDispatcherTest::testHaltingEventExecution
// EventsDispatcherTest::testResponseWhenNoListenersAreSet
// EventsDispatcherTest::testReturningFalseStopsPropagation
// EventsDispatcherTest::testQueuedEventsAreFired
// EventsDispatcherTest::testQueuedEventsCanBeForgotten
// EventsDispatcherTest::testMultiplePushedEventsWillGetFlushed
// EventsDispatcherTest::testPushMethodCanAcceptObjectAsPayload
// EventsDispatcherTest::testWildcardListeners
// EventsDispatcherTest::testWildcardListenersWithResponses
// EventsDispatcherTest::testWildcardListenersCacheFlushing
// EventsDispatcherTest::testListenersCanBeRemoved
// EventsDispatcherTest::testWildcardListenersCanBeRemoved
// EventsDispatcherTest::testWildcardCacheIsClearedWhenListenersAreRemoved
// EventsDispatcherTest::testHasWildcardListeners
// EventsDispatcherTest::testListenersCanBeFound
// EventsDispatcherTest::testWildcardListenersCanBeFound
// EventsDispatcherTest::testClassesWork
// EventsDispatcherTest::testClassesWorkWithAnonymousListeners
// EventsDispatcherTest::testEventClassesArePayload
// EventsDispatcherTest::testNestedEvent
// EventsDispatcherTest::testDuplicateListenersWillFire
// EventsDispatcherTest::testGetListeners
// EventsDispatcherTest::test_Listener_object_creation_is_lazy
// EventsDispatcherTest::testInvokeIsCalled
func TestLaravelEventsDispatcherInventoryEquivalents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dispatcher := events.NewDispatcher()

	var calls []string

	dispatcher.Listen("inventory.created", func(context.Context, any) (any, error) {
		calls = append(calls, "first")

		return "first-response", nil
	})
	dispatcher.Listen("inventory.created", func(context.Context, any) (any, error) {
		calls = append(calls, "second")

		return nil, nil
	})
	dispatcher.Listen("inventory.*", func(context.Context, any) (any, error) {
		calls = append(calls, "wildcard")

		return "wildcard-response", nil
	})

	responses, err := dispatcher.Dispatch(ctx, "inventory.created")

	if err != nil {
		t.Fatal(err)
	}

	if len(calls) != 3 || len(responses) != 2 {
		t.Fatalf("unexpected dispatch calls=%#v responses=%#v", calls, responses)
	}

	if !dispatcher.HasWildcardListeners("inventory.created") || len(dispatcher.GetListeners("inventory.created")) != 3 {
		t.Fatal("expected direct and wildcard listeners to be discoverable")
	}

	halted, err := dispatcher.Until(ctx, "inventory.created")

	if err != nil {
		t.Fatal(err)
	}

	if halted != "first-response" {
		t.Fatalf("expected first non-nil response to halt, got %#v", halted)
	}

	empty, err := events.NewDispatcher().Dispatch(ctx, "missing")

	if err != nil || len(empty) != 0 {
		t.Fatalf("expected no responses for missing listeners, got %#v err=%v", empty, err)
	}

	falseDispatcher := events.NewDispatcher()

	var secondCalled bool

	falseDispatcher.Listen("halt", func(context.Context, any) (any, error) { return false, nil })
	falseDispatcher.Listen("halt", func(context.Context, any) (any, error) {
		secondCalled = true

		return nil, nil
	})

	if result, err := falseDispatcher.Until(ctx, "halt"); err != nil || result != false || secondCalled {
		t.Fatalf("expected false response to halt Until, result=%#v second=%v err=%v", result, secondCalled, err)
	}

	var deferred []string

	dispatcher.Listen("deferred", func(context.Context, any) (any, error) {
		deferred = append(deferred, "ran")

		return nil, nil
	})
	err = dispatcher.Defer(ctx, func(ctx context.Context) error {
		if _, err := dispatcher.Dispatch(ctx, "deferred"); err != nil {
			return err
		}

		if len(deferred) != 0 {
			return errors.New("event ran before defer callback returned")
		}

		return nil
	}, "deferred")

	if err != nil {
		t.Fatal(err)
	}

	if len(deferred) != 1 {
		t.Fatalf("expected deferred event to flush, got %#v", deferred)
	}

	var pushed []string

	dispatcher.Listen(inventoryEvent{}, func(_ context.Context, event any) (any, error) {
		if e, ok := event.(inventoryEvent); ok {
			pushed = append(pushed, e.ID)
		}

		return nil, nil
	})
	dispatcher.Push(ctx, inventoryEvent{ID: "one"})
	dispatcher.Push(ctx, inventoryEvent{ID: "two"})

	if err := dispatcher.Flush(ctx, "events_test.inventoryEvent"); err != nil {
		t.Fatal(err)
	}

	if len(pushed) != 2 {
		t.Fatalf("expected pushed object payloads to flush, got %#v", pushed)
	}

	dispatcher.Push(ctx, "inventory.pushed")
	dispatcher.ForgetPushed()

	if err := dispatcher.Flush(ctx, "inventory.pushed"); err != nil {
		t.Fatal(err)
	}

	if len(pushed) != 2 {
		t.Fatalf("expected forgotten pushed event not to run, got %#v", pushed)
	}

	var structPayload inventoryEvent

	dispatcher.Listen(inventoryEvent{}, func(_ context.Context, event any) (any, error) {
		structPayload = event.(inventoryEvent)

		return nil, nil
	})

	if _, err := dispatcher.Dispatch(ctx, inventoryEvent{ID: "class-payload"}); err != nil {
		t.Fatal(err)
	}

	if structPayload.ID != "class-payload" {
		t.Fatalf("expected struct event payload, got %#v", structPayload)
	}

	var nested inventoryNestedEvent

	dispatcher.Listen(inventoryNestedEvent{}, func(_ context.Context, event any) (any, error) {
		nested = event.(inventoryNestedEvent)

		return nil, nil
	})

	if _, err := dispatcher.Dispatch(ctx, inventoryNestedEvent{Event: inventoryEvent{ID: "nested"}}); err != nil {
		t.Fatal(err)
	}

	if nested.Event.ID != "nested" {
		t.Fatalf("expected nested event payload, got %#v", nested)
	}

	dispatcher.Forget("inventory.created")
	dispatcher.Forget("inventory.*")

	if dispatcher.HasListeners("inventory.created") || dispatcher.HasWildcardListeners("inventory.created") {
		t.Fatal("expected direct and wildcard listeners to be removed")
	}
}

// EventsSubscriberTest::testEventSubscribers
// EventsSubscriberTest::testEventSubscribeCanAcceptObject
// EventsSubscriberTest::testEventSubscribeCanReturnMappings
func TestLaravelEventsSubscriberInventoryEquivalents(t *testing.T) {
	t.Parallel()

	dispatcher := events.NewDispatcher()
	subscriber := &inventorySubscriber{}
	dispatcher.Subscribe(subscriber)

	if _, err := dispatcher.Dispatch(context.Background(), "inventory.created"); err != nil {
		t.Fatal(err)
	}

	if len(subscriber.order) != 1 || subscriber.order[0] != "created" {
		t.Fatalf("unexpected subscriber calls: %#v", subscriber.order)
	}
}

// QueuedEventsTest::testQueuedEventHandlersAreQueued
// QueuedEventsTest::testCustomizedQueuedEventHandlersAreQueued
// QueuedEventsTest::testQueueIsSetByGetQueue
// QueuedEventsTest::testQueueIsSetByGetConnection
// QueuedEventsTest::testDelayIsSetByWithDelay
// QueuedEventsTest::testQueueIsSetByGetQueueDynamically
// QueuedEventsTest::testQueueIsSetByGetConnectionDynamically
// QueuedEventsTest::testQueueIsSetUsingQueueRoutes
// QueuedEventsTest::testDelayIsSetByWithDelayDynamically
// QueuedEventsTest::testQueuePropagateRetryUntilAndMaxExceptions
// QueuedEventsTest::testQueuePropagateTries
// QueuedEventsTest::testQueuePropagateMessageGroupProperty
// QueuedEventsTest::testQueuePropagateMessageGroupMethodOverProperty
// QueuedEventsTest::testQueuePropagateDeduplicationIdMethod
// QueuedEventsTest::testQueuePropagateDeduplicatorMethodOverDeduplicationIdMethod
// QueuedEventsTest::testQueuePropagateMiddleware
// QueuedEventsTest::testQueuePropagatesShouldBeUnique
// QueuedEventsTest::testQueuePropagatesShouldBeUniqueUntilProcessing
// QueuedEventsTest::testQueuePropagatesUniqueIdFromMethod
func TestLaravelQueuedEventsInventoryEquivalents(t *testing.T) {
	t.Parallel()

	delay := 2 * time.Second
	listener := events.NewCallQueuedListener("InventoryListener", inventoryEvent{ID: "queued"}).
		WithOptions(events.ListenerOptions{
			Connection:    "redis",
			Queue:         "events",
			Delay:         delay,
			Tries:         3,
			MaxExceptions: 2,
			Timeout:       time.Minute,
			Backoff:       []time.Duration{time.Second, 2 * time.Second},
		})
	listener.ShouldBeUniqueFlag = true
	listener.UniqueID = "inventory-1"
	listener.UniqueFor = time.Minute

	if listener.GetConnection() != "redis" || listener.GetQueue() != "events" || listener.GetDelay() != delay {
		t.Fatalf("unexpected queued listener routing: %#v", listener)
	}

	if listener.GetTries() != 3 || listener.GetMaxExceptions() != 2 || listener.GetTimeout() != time.Minute {
		t.Fatalf("unexpected queued listener retry options: %#v", listener)
	}

	if len(listener.GetBackoff()) != 2 || !listener.ShouldBeUniqueFlag || listener.UniqueID != "inventory-1" {
		t.Fatalf("unexpected queued listener uniqueness/backoff options: %#v", listener)
	}

	var handled bool
	closure := events.Queueable(func(context.Context, any) (any, error) {
		handled = true

		return nil, nil
	}).OnConnection("redis").OnQueue("listeners").WithDelay(delay)

	if closure.GetConnection() != "redis" || closure.GetQueue() != "listeners" || closure.GetDelay() != delay {
		t.Fatalf("unexpected queued closure options")
	}

	if err := (&events.InvokeQueuedClosure{}).Handle(context.Background(), closure, inventoryEvent{ID: "queued"}); err != nil {
		t.Fatal(err)
	}

	if !handled {
		t.Fatal("expected queued closure handler to run")
	}
}
