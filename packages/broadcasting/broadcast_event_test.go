package broadcasting_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bedrock/packages/broadcasting"
)

// BroadcastEventTest::testBasicEventBroadcastParameterFormatting

// BroadcastEventTest::testManualParameterSpecification

// BroadcastEventTest::testSpecificBroadcasterGiven

// BroadcastEventTest::testSpecificChannelsPerConnection

// BroadcastEventTest::testMiddlewareProxiesMiddlewareFromUnderlyingEvent
// BroadcastEventTest::testMiddlewareProxiesFailedHandlerFromUnderlyingEvent

type testArrayable struct {
	value map[string]any
}

type testBroadcastEvent struct {
	FirstName  string
	LastName   string
	Collection testArrayable
	title      string
}

type testManualBroadcastEvent struct{}

type testSpecificConnectionEvent struct{}

type testPerConnectionEvent struct{}

type testMiddlewareEvent struct {
	failed error
}

func TestBroadcastEventBuildsDefaultPayload(t *testing.T) {
	t.Parallel()

	broadcaster := &recordingBroadcaster{}
	manager := broadcasting.NewManager().Extend("", broadcaster)
	event := testBroadcastEvent{FirstName: "Taylor", LastName: "Otwell", Collection: testArrayable{value: map[string]any{"foo": "bar"}}, title: "Developer"}

	if err := broadcasting.NewBroadcastEvent(event).Handle(context.Background(), manager); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	requireEqual(t, broadcaster.calls[0].channels, []string{"test-channel"})
	requireEqual(t, broadcaster.calls[0].payload, map[string]any{
		"FirstName":  "Taylor",
		"LastName":   "Otwell",
		"Collection": map[string]any{"foo": "bar"},
	})
}

func TestBroadcastEventUsesManualPayload(t *testing.T) {
	t.Parallel()

	broadcaster := &recordingBroadcaster{}
	manager := broadcasting.NewManager().Extend("", broadcaster)

	if err := broadcasting.NewBroadcastEvent(testManualBroadcastEvent{}).Handle(context.Background(), manager); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	requireEqual(t, broadcaster.calls[0].payload, map[string]any{"name": "Taylor", "socket": nil})
}

func TestBroadcastEventUsesSpecificBroadcaster(t *testing.T) {
	t.Parallel()

	broadcaster := &recordingBroadcaster{}
	manager := broadcasting.NewManager().Extend("log", broadcaster)

	if err := broadcasting.NewBroadcastEvent(testSpecificConnectionEvent{}).Handle(context.Background(), manager); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	requireEqual(t, broadcaster.calls[0].channels, []string{"test-channel"})
}

func TestBroadcastEventUsesSpecificChannelsPerConnection(t *testing.T) {
	t.Parallel()

	first := &recordingBroadcaster{}
	second := &recordingBroadcaster{}
	manager := broadcasting.NewManager().Extend("first_connection", first).Extend("second_connection", second)

	if err := broadcasting.NewBroadcastEvent(testPerConnectionEvent{}).Handle(context.Background(), manager); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	requireEqual(t, first.calls[0].channels, []string{"first-channel"})
	requireEqual(t, first.calls[0].payload, map[string]any{"firstName": "Taylor", "lastName": "Otwell"})
	requireEqual(t, second.calls[0].channels, []string{"second-channel"})
	requireEqual(t, second.calls[0].payload, map[string]any{"firstName": "Taylor"})
}

func TestBroadcastEventProxiesMiddlewareAndFailed(t *testing.T) {
	t.Parallel()

	event := &testMiddlewareEvent{}
	job := broadcasting.NewBroadcastEvent(event)

	requireEqual(t, job.Middleware(), []any{"foo", "bar"})

	err := errors.New("failed")
	job.Failed(err)

	if event.failed != err {
		t.Fatalf("failed error = %v, want %v", event.failed, err)
	}
}

func (a testArrayable) ToArray() any { return a.value }

func (testBroadcastEvent) BroadcastOn() any { return []string{"test-channel"} }

func (testManualBroadcastEvent) BroadcastOn() any { return []string{"test-channel"} }

func (testManualBroadcastEvent) BroadcastWith() map[string]any {
	return map[string]any{"name": "Taylor"}
}

func (testSpecificConnectionEvent) BroadcastOn() any { return []string{"test-channel"} }

func (testSpecificConnectionEvent) BroadcastConnections() []string {
	return []string{"log"}
}

func (testPerConnectionEvent) BroadcastConnections() []string {
	return []string{"first_connection", "second_connection"}
}

func (testPerConnectionEvent) BroadcastOn() any {
	return map[string][]string{
		"first_connection":  {"first-channel"},
		"second_connection": {"second-channel"},
	}
}

func (testPerConnectionEvent) BroadcastWith() map[string]any {
	return map[string]any{
		"first_connection":  map[string]any{"firstName": "Taylor", "lastName": "Otwell"},
		"second_connection": map[string]any{"firstName": "Taylor"},
	}
}

func (e *testMiddlewareEvent) Middleware() []any { return []any{"foo", "bar"} }

func (e *testMiddlewareEvent) Failed(err error) { e.failed = err }
