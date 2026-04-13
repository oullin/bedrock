package log_test

import (
	"context"
	"sync"
	"testing"

	cevents "github.com/bedrock/packages/contracts/events"
	"github.com/bedrock/packages/log"
)

// spyHandler captures all handled records for assertions.
type spyHandler struct {
	mu      sync.Mutex
	records []log.Record
	level   log.Level
}

func newSpyHandler(level log.Level) *spyHandler {
	return &spyHandler{level: level}
}

func (h *spyHandler) Handle(record log.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.records = append(h.records, record)

	return nil
}

func (h *spyHandler) IsHandling(level log.Level) bool {
	return level >= h.level
}

func (h *spyHandler) Close() error {
	return nil
}

func (h *spyHandler) Records() []log.Record {
	h.mu.Lock()
	defer h.mu.Unlock()

	cp := make([]log.Record, len(h.records))
	copy(cp, h.records)

	return cp
}

func (h *spyHandler) LastRecord() log.Record {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.records[len(h.records)-1]
}

// spyDispatcher captures dispatched events for assertions.
type spyDispatcher struct {
	mu        sync.Mutex
	events    []any
	listeners map[string][]cevents.Listener
}

func newSpyDispatcher() *spyDispatcher {
	return &spyDispatcher{
		listeners: make(map[string][]cevents.Listener),
	}
}

func (d *spyDispatcher) Listen(events any, listeners ...cevents.Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()

	name := "default"
	if s, ok := events.(string); ok {
		name = s
	}

	d.listeners[name] = append(d.listeners[name], listeners...)
}

func (d *spyDispatcher) HasListeners(_ any) bool  { return false }
func (d *spyDispatcher) HasWildcardListeners(_ any) bool { return false }
func (d *spyDispatcher) Subscribe(_ cevents.Subscriber) {}
func (d *spyDispatcher) Until(_ context.Context, _ any) (any, error) { return nil, nil }

func (d *spyDispatcher) Dispatch(_ context.Context, event any) ([]any, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.events = append(d.events, event)

	return nil, nil
}

func (d *spyDispatcher) Push(_ context.Context, _ any) {}
func (d *spyDispatcher) Flush(_ context.Context, _ string) error { return nil }
func (d *spyDispatcher) Forget(_ any) {}
func (d *spyDispatcher) ForgetPushed() {}
func (d *spyDispatcher) GetListeners(_ any) []cevents.Listener { return nil }

func (d *spyDispatcher) Events() []any {
	d.mu.Lock()
	defer d.mu.Unlock()

	cp := make([]any, len(d.events))
	copy(cp, d.events)

	return cp
}

func TestLoggerEmergency(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test")

	logger.Emergency("system down")

	records := spy.Records()
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	if records[0].Level != log.LevelEmergency {
		t.Fatalf("expected LevelEmergency, got %d", records[0].Level)
	}

	if records[0].Message != "system down" {
		t.Fatalf("expected 'system down', got %q", records[0].Message)
	}
}

func TestLoggerAllLevels(t *testing.T) {
	t.Parallel()

	levels := []struct {
		name   string
		level  log.Level
		method func(*log.Logger, string, ...map[string]any)
	}{
		{"Emergency", log.LevelEmergency, (*log.Logger).Emergency},
		{"Alert", log.LevelAlert, (*log.Logger).Alert},
		{"Critical", log.LevelCritical, (*log.Logger).Critical},
		{"Error", log.LevelError, (*log.Logger).Error},
		{"Warning", log.LevelWarning, (*log.Logger).Warning},
		{"Notice", log.LevelNotice, (*log.Logger).Notice},
		{"Info", log.LevelInfo, (*log.Logger).Info},
		{"Debug", log.LevelDebug, (*log.Logger).Debug},
	}

	for _, tt := range levels {
		spy := newSpyHandler(log.LevelDebug)
		logger := log.NewLogger(spy, "test")

		tt.method(logger, "test message")

		records := spy.Records()
		if len(records) != 1 {
			t.Fatalf("%s: expected 1 record, got %d", tt.name, len(records))
		}

		if records[0].Level != tt.level {
			t.Fatalf("%s: expected level %d, got %d", tt.name, tt.level, records[0].Level)
		}
	}
}

func TestLoggerLog(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test")

	logger.Log(log.LevelWarning, "generic log")

	record := spy.LastRecord()
	if record.Level != log.LevelWarning {
		t.Fatalf("expected LevelWarning, got %d", record.Level)
	}
}

func TestLoggerContext(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test")

	withCtx := logger.WithContext(map[string]any{"user_id": 42})
	withCtx.Info("with context")

	record := spy.LastRecord()
	if record.Context["user_id"] != 42 {
		t.Fatalf("expected user_id = 42, got %v", record.Context["user_id"])
	}
}

func TestLoggerWithoutContext(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test", log.WithLoggerContext(map[string]any{"a": 1, "b": 2}))

	cleared := logger.WithoutContext()
	cleared.Info("no context")

	record := spy.LastRecord()
	if len(record.Context) != 0 {
		t.Fatalf("expected empty context, got %v", record.Context)
	}
}

func TestLoggerContextMerge(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test", log.WithLoggerContext(map[string]any{"a": 1}))

	logger.Info("merged", map[string]any{"b": 2})

	record := spy.LastRecord()
	if record.Context["a"] != 1 {
		t.Fatalf("expected a = 1, got %v", record.Context["a"])
	}

	if record.Context["b"] != 2 {
		t.Fatalf("expected b = 2, got %v", record.Context["b"])
	}
}

func TestLoggerCallSiteContextOverrides(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test", log.WithLoggerContext(map[string]any{"key": "logger"}))

	logger.Info("override", map[string]any{"key": "callsite"})

	record := spy.LastRecord()
	if record.Context["key"] != "callsite" {
		t.Fatalf("expected call-site context to override, got %v", record.Context["key"])
	}
}

func TestLoggerEventDispatch(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	dispatcher := newSpyDispatcher()
	logger := log.NewLogger(spy, "test", log.WithDispatcher(dispatcher))

	logger.Error("event test", map[string]any{"req_id": "abc"})

	events := dispatcher.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	msg, ok := events[0].(log.MessageLogged)
	if !ok {
		t.Fatalf("expected MessageLogged event, got %T", events[0])
	}

	if msg.Level != "error" {
		t.Fatalf("expected level 'error', got %q", msg.Level)
	}

	if msg.Message != "event test" {
		t.Fatalf("expected message 'event test', got %q", msg.Message)
	}

	if msg.Context["req_id"] != "abc" {
		t.Fatalf("expected req_id = abc, got %v", msg.Context["req_id"])
	}
}

func TestLoggerEventDispatchNilDispatcher(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test")

	// Should not panic with nil dispatcher.
	logger.Info("no dispatcher")

	if len(spy.Records()) != 1 {
		t.Fatal("expected record to be written even without dispatcher")
	}
}

func TestLoggerListenNoDispatcher(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test")

	err := logger.Listen(func(_ log.MessageLogged) {})
	if err == nil {
		t.Fatal("expected error when calling Listen without dispatcher")
	}
}

func TestLoggerWithoutContextSelectiveKeys(t *testing.T) {
	t.Parallel()

	spy := newSpyHandler(log.LevelDebug)
	logger := log.NewLogger(spy, "test", log.WithLoggerContext(map[string]any{"a": 1, "b": 2, "c": 3}))

	filtered := logger.WithoutContext("b")
	filtered.Info("filtered")

	record := spy.LastRecord()
	if record.Context["a"] != 1 {
		t.Fatalf("expected a = 1, got %v", record.Context["a"])
	}

	if _, ok := record.Context["b"]; ok {
		t.Fatal("expected key 'b' to be removed")
	}

	if record.Context["c"] != 3 {
		t.Fatalf("expected c = 3, got %v", record.Context["c"])
	}
}
