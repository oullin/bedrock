package queue_test

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// fakeInspectableQueue is a minimal Queue that also satisfies the two
// optional inspection contracts (QueueNamer + JobInspector). It serves
// the manager fan-out tests without dragging in a real driver.
type fakeInspectableQueue struct {
	connection string
	names      []string
	pending    map[string][]queue.InspectedJob
	delayed    map[string][]queue.InspectedJob
	reserved   map[string][]queue.InspectedJob
	pendingErr map[string]error
}

// inertQueue implements queue.Queue but neither QueueNamer nor
// JobInspector. The manager's fan-out should surface ErrNotSupported
// for it.
type inertQueue struct{ connection string }

func (f *fakeInspectableQueue) Push(context.Context, string, []byte) (string, error) {
	return "", nil
}

func (f *fakeInspectableQueue) PushDelayed(context.Context, string, []byte, time.Duration) (string, error) {
	return "", nil
}

func (f *fakeInspectableQueue) PushMultiple(_ context.Context, _ string, p [][]byte) ([]string, error) {
	return make([]string, len(p)), nil
}

func (f *fakeInspectableQueue) Pop(context.Context, string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (f *fakeInspectableQueue) Size(context.Context, string) (int64, error) { return 0, nil }

func (f *fakeInspectableQueue) PendingSize(context.Context, string) (int64, error) {
	return 0, nil
}

func (f *fakeInspectableQueue) DelayedSize(context.Context, string) (int64, error) {
	return 0, nil
}

func (f *fakeInspectableQueue) ReservedSize(context.Context, string) (int64, error) {
	return 0, nil
}

func (f *fakeInspectableQueue) ConnectionName() string { return f.connection }

func (f *fakeInspectableQueue) QueueNames(_ context.Context) ([]string, error) {
	return append([]string(nil), f.names...), nil
}

func (f *fakeInspectableQueue) PendingJobs(_ context.Context, name string) ([]queue.InspectedJob, error) {
	if err, ok := f.pendingErr[name]; ok {
		return nil, err
	}

	return f.pending[name], nil
}

func (f *fakeInspectableQueue) DelayedJobs(_ context.Context, name string) ([]queue.InspectedJob, error) {
	return f.delayed[name], nil
}

func (f *fakeInspectableQueue) ReservedJobs(_ context.Context, name string) ([]queue.InspectedJob, error) {
	return f.reserved[name], nil
}

func (q *inertQueue) Push(context.Context, string, []byte) (string, error) { return "", nil }

func (q *inertQueue) PushDelayed(context.Context, string, []byte, time.Duration) (string, error) {
	return "", nil
}

func (q *inertQueue) PushMultiple(_ context.Context, _ string, p [][]byte) ([]string, error) {
	return make([]string, len(p)), nil
}

func (q *inertQueue) Pop(context.Context, string) (queue.Job, error)      { return nil, queue.ErrNoJob }
func (q *inertQueue) Size(context.Context, string) (int64, error)         { return 0, nil }
func (q *inertQueue) PendingSize(context.Context, string) (int64, error)  { return 0, nil }
func (q *inertQueue) DelayedSize(context.Context, string) (int64, error)  { return 0, nil }
func (q *inertQueue) ReservedSize(context.Context, string) (int64, error) { return 0, nil }
func (q *inertQueue) ConnectionName() string                              { return q.connection }

func newManagerWithQueue(t *testing.T, name string, q queue.Queue) *queue.Manager {
	t.Helper()

	m := queue.NewManager()
	m.SetConfig(name, map[string]any{"driver": "fake"})
	m.Register("fake", func(_ map[string]any) (queue.Queue, error) { return q, nil })

	return m
}

func TestManagerAllPendingJobsConcatenatesAcrossQueues(t *testing.T) {
	t.Parallel()

	fq := &fakeInspectableQueue{
		connection: "primary",
		names:      []string{"high", "low"},
		pending: map[string][]queue.InspectedJob{
			"high": {{ID: 1, Queue: "high", Connection: "primary", UUID: "u-high"}},
			"low":  {{ID: 2, Queue: "low", Connection: "primary", UUID: "u-low-a"}, {ID: 3, Queue: "low", Connection: "primary", UUID: "u-low-b"}},
		},
	}

	m := newManagerWithQueue(t, "primary", fq)

	got, err := m.AllPendingJobs(context.Background(), "primary")

	if err != nil {
		t.Fatalf("AllPendingJobs: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("got %d jobs, want 3", len(got))
	}

	uuids := []string{got[0].UUID, got[1].UUID, got[2].UUID}
	want := []string{"u-high", "u-low-a", "u-low-b"}

	if !reflect.DeepEqual(uuids, want) {
		t.Errorf("UUIDs: got %v, want %v", uuids, want)
	}
}

func TestManagerAllDelayedJobsDeduplicatesByQueueIteration(t *testing.T) {
	t.Parallel()

	fq := &fakeInspectableQueue{
		connection: "primary",
		names:      []string{"only"},
		delayed: map[string][]queue.InspectedJob{
			"only": {{ID: 1, Queue: "only", UUID: "u-1"}, {ID: 2, Queue: "only", UUID: "u-2"}},
		},
	}

	m := newManagerWithQueue(t, "primary", fq)

	got, err := m.AllDelayedJobs(context.Background(), "primary")

	if err != nil {
		t.Fatalf("AllDelayedJobs: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("got %d, want 2", len(got))
	}
}

func TestManagerAllReservedJobsDriverWithoutInspector(t *testing.T) {
	t.Parallel()

	m := newManagerWithQueue(t, "shy", &inertQueue{connection: "shy"})

	_, err := m.AllReservedJobs(context.Background(), "shy")

	if !errors.Is(err, queue.ErrNotSupported) {
		t.Fatalf("error: got %v, want ErrNotSupported", err)
	}
}

func TestManagerAllPendingJobsSkipsPerQueueNotSupported(t *testing.T) {
	t.Parallel()

	fq := &fakeInspectableQueue{
		connection: "mixed",
		names:      []string{"yes", "no"},
		pending: map[string][]queue.InspectedJob{
			"yes": {{ID: 1, Queue: "yes", UUID: "u-yes"}},
		},
		pendingErr: map[string]error{
			"no": queue.ErrNotSupported,
		},
	}

	m := newManagerWithQueue(t, "mixed", fq)

	got, err := m.AllPendingJobs(context.Background(), "mixed")

	if err != nil {
		t.Fatalf("AllPendingJobs: %v", err)
	}

	if len(got) != 1 || got[0].UUID != "u-yes" {
		t.Errorf("got %+v, want single u-yes snapshot", got)
	}
}

func TestManagerAllPendingJobsAggregatesRealErrors(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("backend down")

	fq := &fakeInspectableQueue{
		connection: "broken",
		names:      []string{"q1", "q2"},
		pending: map[string][]queue.InspectedJob{
			"q1": {{ID: 1, UUID: "u-1"}},
		},
		pendingErr: map[string]error{"q2": sentinel},
	}

	m := newManagerWithQueue(t, "broken", fq)

	got, err := m.AllPendingJobs(context.Background(), "broken")

	if !errors.Is(err, sentinel) {
		t.Fatalf("error: got %v, want sentinel", err)
	}

	if len(got) != 1 {
		t.Errorf("snapshots: got %d, want 1 (q1 should still surface)", len(got))
	}
}

// TestQueueNamerOrderingIsDeterministic guards against the manager
// silently re-sorting queue names; ordering should match the driver's
// reported order so dashboards stay stable.
func TestQueueNamerOrderingIsDeterministic(t *testing.T) {
	t.Parallel()

	fq := &fakeInspectableQueue{
		connection: "stable",
		names:      []string{"c", "a", "b"},
		pending: map[string][]queue.InspectedJob{
			"a": {{UUID: "a"}},
			"b": {{UUID: "b"}},
			"c": {{UUID: "c"}},
		},
	}

	m := newManagerWithQueue(t, "stable", fq)

	got, err := m.AllPendingJobs(context.Background(), "stable")

	if err != nil {
		t.Fatalf("AllPendingJobs: %v", err)
	}

	uuids := []string{got[0].UUID, got[1].UUID, got[2].UUID}

	if !reflect.DeepEqual(uuids, []string{"c", "a", "b"}) {
		t.Errorf("ordering not preserved: got %v", uuids)
	}

	// sanity: the test would catch a re-sort
	sorted := append([]string(nil), uuids...)

	sort.Strings(sorted)

	if reflect.DeepEqual(uuids, sorted) {
		t.Errorf("ordering matches sorted output — manager may be sorting")
	}
}
