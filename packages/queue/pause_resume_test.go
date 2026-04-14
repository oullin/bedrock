package queue_test

import (
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/events"
)

// Ports of Framework\Tests\Queue\QueuePauseResumeTest.
//
// Upstream's test constructs a QueueManager against a Carbon test clock
// and an ArrayStore-backed cache. The Go equivalent exercises the
// underlying PauseResumer + InMemoryPauseStore directly — when Step 6
// wires up the full Manager, the same contract must still hold and
// these tests continue to pass because Manager just delegates.
//
// Event assertions use the clockRecorder/pauseEventRecorder helpers below,
// which match Mockery's role in the PHP test (capture the dispatched
// event and assert on its fields).

// --- helpers ----------------------------------------------------------

// mockClock is a deterministic time source we can freeze and advance,
// mirroring Upstream's Carbon::setTestNow().
type mockClock struct {
	mu  sync.Mutex
	now time.Time
}

// pauseEventRecorder collects every event the PauseResumer emits so tests
// can assert the dispatched payloads.
type pauseEventRecorder struct {
	mu     sync.Mutex
	events []any
}

func newMockClock(now time.Time) *mockClock { return &mockClock{now: now} }

func (c *mockClock) Now() time.Time {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.now
}

func (c *mockClock) Advance(d time.Duration) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.now = c.now.Add(d)
}

func (r *pauseEventRecorder) Emit(event any) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.events = append(r.events, event)
}

func (r *pauseEventRecorder) last() any {
	r.mu.Lock()

	defer r.mu.Unlock()

	if len(r.events) == 0 {
		return nil
	}

	return r.events[len(r.events)-1]
}

// newPauseResumer wires up a PauseResumer with its own in-memory store,
// recorder, and mock clock. The returned clock is bound into the store
// so PauseFor tests can advance it deterministically.
func newPauseResumer() (*queue.PauseResumer, *pauseEventRecorder, *mockClock) {
	store := queue.NewInMemoryPauseStore()
	clock := newMockClock(time.Date(2026, 4, 14, 12, 0, 0, 0, time.UTC))
	store.SetClock(clock.Now)
	rec := &pauseEventRecorder{}

	return queue.NewPauseResumer(store, rec), rec, clock
}

// --- ports ------------------------------------------------------------

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testPauseQueueWithConnection
func TestPauseQueueWithConnection(t *testing.T) {
	t.Parallel()

	pr, _, _ := newPauseResumer()

	if err := pr.Pause("redis", "default"); err != nil {
		t.Fatalf("Pause: %v", err)
	}

	if !pr.IsPaused("redis", "default") {
		t.Error("expected redis:default to be paused")
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testPauseQueueWithTTL
func TestPauseQueueWithTTL(t *testing.T) {
	t.Parallel()

	pr, _, clock := newPauseResumer()

	if err := pr.PauseFor("redis", "default", 30*time.Second); err != nil {
		t.Fatalf("PauseFor: %v", err)
	}

	if !pr.IsPaused("redis", "default") {
		t.Error("expected redis:default to be paused immediately after PauseFor")
	}

	clock.Advance(time.Minute)

	if pr.IsPaused("redis", "default") {
		t.Error("expected redis:default to have expired after advancing 1 minute")
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testPauseQueueIndefinitely
func TestPauseQueueIndefinitely(t *testing.T) {
	t.Parallel()

	pr, _, clock := newPauseResumer()

	if err := pr.Pause("redis", "default"); err != nil {
		t.Fatalf("Pause: %v", err)
	}

	if !pr.IsPaused("redis", "default") {
		t.Error("expected indefinite pause to be active immediately")
	}

	clock.Advance(365 * 24 * time.Hour)

	if !pr.IsPaused("redis", "default") {
		t.Error("expected indefinite pause to persist after 1 year")
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testResumeQueue
func TestResumeQueue(t *testing.T) {
	t.Parallel()

	pr, _, _ := newPauseResumer()

	_ = pr.Pause("redis", "default")

	if !pr.IsPaused("redis", "default") {
		t.Fatal("pre-condition: expected paused")
	}

	if err := pr.Resume("redis", "default"); err != nil {
		t.Fatalf("Resume: %v", err)
	}

	if pr.IsPaused("redis", "default") {
		t.Error("expected redis:default to be resumed")
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testPausingQueueOnOneConnectionDoesNotAffectAnother
func TestPausingQueueOnOneConnectionDoesNotAffectAnother(t *testing.T) {
	t.Parallel()

	pr, _, _ := newPauseResumer()

	_ = pr.Pause("redis", "default")

	if !pr.IsPaused("redis", "default") {
		t.Error("expected redis:default paused")
	}

	if pr.IsPaused("database", "default") {
		t.Error("expected database:default NOT paused")
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testPausingDifferentQueuesOnSameConnection
func TestPausingDifferentQueuesOnSameConnection(t *testing.T) {
	t.Parallel()

	pr, _, _ := newPauseResumer()

	_ = pr.Pause("redis", "emails")
	_ = pr.Pause("redis", "notifications")

	if !pr.IsPaused("redis", "emails") {
		t.Error("expected redis:emails paused")
	}

	if !pr.IsPaused("redis", "notifications") {
		t.Error("expected redis:notifications paused")
	}

	if pr.IsPaused("redis", "default") {
		t.Error("expected redis:default NOT paused")
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testResumingOnlyAffectsSpecificQueue
func TestResumingOnlyAffectsSpecificQueue(t *testing.T) {
	t.Parallel()

	pr, _, _ := newPauseResumer()

	_ = pr.Pause("redis", "emails")
	_ = pr.Pause("redis", "notifications")

	_ = pr.Resume("redis", "emails")

	if pr.IsPaused("redis", "emails") {
		t.Error("expected redis:emails resumed")
	}

	if !pr.IsPaused("redis", "notifications") {
		t.Error("expected redis:notifications still paused")
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testPauseDispatchesQueuePausedEvent
func TestPauseDispatchesQueuePausedEvent(t *testing.T) {
	t.Parallel()

	pr, rec, _ := newPauseResumer()

	_ = pr.Pause("redis", "default")

	ev, ok := rec.last().(events.QueuePaused)

	if !ok {
		t.Fatalf("expected QueuePaused event, got %T (%v)", rec.last(), rec.last())
	}

	if ev.ConnectionName != "redis" {
		t.Errorf("ConnectionName: got %q, want redis", ev.ConnectionName)
	}

	if ev.Queue != "default" {
		t.Errorf("Queue: got %q, want default", ev.Queue)
	}

	if ev.TTL != nil {
		t.Errorf("TTL: got %v, want nil", *ev.TTL)
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testPauseForDispatchesQueuePausedEventWithTTL
func TestPauseForDispatchesQueuePausedEventWithTTL(t *testing.T) {
	t.Parallel()

	pr, rec, _ := newPauseResumer()

	_ = pr.PauseFor("redis", "emails", 60*time.Second)

	ev, ok := rec.last().(events.QueuePaused)

	if !ok {
		t.Fatalf("expected QueuePaused event, got %T", rec.last())
	}

	if ev.ConnectionName != "redis" || ev.Queue != "emails" {
		t.Errorf("ConnectionName/Queue: got %q/%q", ev.ConnectionName, ev.Queue)
	}

	if ev.TTL == nil {
		t.Fatal("TTL: nil, want non-nil 60s")
	}

	if *ev.TTL != 60*time.Second {
		t.Errorf("TTL: got %s, want 60s", *ev.TTL)
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testResumeDispatchesQueueResumedEvent
func TestResumeDispatchesQueueResumedEvent(t *testing.T) {
	t.Parallel()

	pr, rec, _ := newPauseResumer()

	_ = pr.Resume("database", "notifications")

	ev, ok := rec.last().(events.QueueResumed)

	if !ok {
		t.Fatalf("expected QueueResumed event, got %T", rec.last())
	}

	if ev.ConnectionName != "database" {
		t.Errorf("ConnectionName: got %q, want database", ev.ConnectionName)
	}

	if ev.Queue != "notifications" {
		t.Errorf("Queue: got %q, want notifications", ev.Queue)
	}
}

// Port of Framework\Tests\Queue\QueuePauseResumeTest::testParsingQueueString
func TestParsingQueueString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		raw            string
		wantConnection string
		wantQueue      string
	}{
		{"", "redis", "default"},
		{"emails", "redis", "emails"},
		{"database:notifications", "database", "notifications"},
		{"redis:foo:bar", "redis", "foo:bar"},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.raw, func(t *testing.T) {
			t.Parallel()

			gotConn, gotQ := queue.ParseQueue(tc.raw, "redis")

			if gotConn != tc.wantConnection || gotQ != tc.wantQueue {
				t.Errorf("ParseQueue(%q): got (%q, %q), want (%q, %q)",
					tc.raw, gotConn, gotQ, tc.wantConnection, tc.wantQueue)
			}
		})
	}
}
