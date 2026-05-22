package broadcastclient_test

import (
	"testing"

	"github.com/bedrock/packages/broadcastclient"
)

// TestDispatchChannelTriggersAllListeners verifies that all registered
// callbacks for an event are invoked when that event is dispatched, and
// that callbacks for other events are not triggered. This is a direct port
// of "triggers all listeners for an event" from socketio-channel.test.ts.
func TestDispatchChannelTriggersAllListeners(t *testing.T) {
	t.Parallel()

	// namespace: false equivalent — empty namespace
	ch := broadcastclient.NewDispatchChannel(broadcastclient.NewEventFormatter(""))

	var l1Called, l2Called, l3Called bool

	// Named variables are required for reflect-based callback identity.
	l1 := func(data any) { l1Called = true }
	l2 := func(data any) { l2Called = true }
	l3 := func(data any) { l3Called = true }

	ch.Listen("MyEvent", l1)
	ch.Listen("MyEvent", l2)
	ch.Listen("MyOtherEvent", l3)

	ch.Dispatch("MyEvent", map[string]any{})

	if !l1Called {
		t.Fatal("expected l1 to be called after dispatching MyEvent")
	}

	if !l2Called {
		t.Fatal("expected l2 to be called after dispatching MyEvent")
	}

	if l3Called {
		t.Fatal("expected l3 NOT to be called after dispatching MyEvent")
	}

	ch.Dispatch("MyOtherEvent", map[string]any{})

	if !l3Called {
		t.Fatal("expected l3 to be called after dispatching MyOtherEvent")
	}
}

// TestDispatchChannelCanRemoveSpecificListener verifies that a specific
// callback can be removed while leaving others intact. This is a direct port
// of "can remove a listener for an event" from socketio-channel.test.ts.
func TestDispatchChannelCanRemoveSpecificListener(t *testing.T) {
	t.Parallel()

	ch := broadcastclient.NewDispatchChannel(broadcastclient.NewEventFormatter(""))

	var l1Called, l2Called, l3Called bool

	l1 := func(data any) { l1Called = true }
	l2 := func(data any) { l2Called = true }
	l3 := func(data any) { l3Called = true }

	ch.Listen("MyEvent", l1)
	ch.Listen("MyEvent", l2)
	ch.Listen("MyOtherEvent", l3)

	ch.StopListening("MyEvent", l1)

	ch.Dispatch("MyEvent", map[string]any{})

	if l1Called {
		t.Fatal("expected l1 NOT to be called after being removed")
	}

	if !l2Called {
		t.Fatal("expected l2 to be called after dispatching MyEvent")
	}

	if l3Called {
		t.Fatal("expected l3 NOT to be called after dispatching MyEvent")
	}

	ch.Dispatch("MyOtherEvent", map[string]any{})

	if !l3Called {
		t.Fatal("expected l3 to be called after dispatching MyOtherEvent")
	}
}

// TestDispatchChannelCanRemoveAllListeners verifies that passing nil to
// StopListening removes all callbacks for that event. This is a direct port
// of "can remove all listeners for an event" from socketio-channel.test.ts.
func TestDispatchChannelCanRemoveAllListeners(t *testing.T) {
	t.Parallel()

	ch := broadcastclient.NewDispatchChannel(broadcastclient.NewEventFormatter(""))

	var l1Called, l2Called, l3Called bool

	l1 := func(data any) { l1Called = true }
	l2 := func(data any) { l2Called = true }
	l3 := func(data any) { l3Called = true }

	ch.Listen("MyEvent", l1)
	ch.Listen("MyEvent", l2)
	ch.Listen("MyOtherEvent", l3)

	// nil removes all listeners for MyEvent.
	ch.StopListening("MyEvent", nil)

	ch.Dispatch("MyEvent", map[string]any{})

	if l1Called {
		t.Fatal("expected l1 NOT to be called after all listeners removed")
	}

	if l2Called {
		t.Fatal("expected l2 NOT to be called after all listeners removed")
	}

	if l3Called {
		t.Fatal("expected l3 NOT to be called after dispatching MyEvent")
	}

	ch.Dispatch("MyOtherEvent", map[string]any{})

	if !l3Called {
		t.Fatal("expected l3 to be called after dispatching MyOtherEvent")
	}
}
