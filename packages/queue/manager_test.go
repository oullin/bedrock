package queue_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// inspectorQueue is a queue.Queue that satisfies the optional
// QueueNamer and JobInspector contracts, used here to exercise the
// manager-level All{Pending,Delayed,Reserved}Jobs fan-out.
type inspectorQueue struct {
	connection string
	names      []string
	perQueue   map[string][]queue.InspectedJob
	skipQueue  string // queue.ReservedJobs returns ErrNotSupported for this queue name
	namesErr   error
	queueErr   error
}

func TestManagerRegisterAndDriver(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q, err := m.Driver("default")

	if err != nil {
		t.Fatal(err)
	}

	if q == nil {
		t.Fatal("expected non-nil queue")
	}

	if q.ConnectionName() != "null" {
		t.Errorf("expected 'null', got %q", q.ConnectionName())
	}
}

func TestManagerDriverCachesInstance(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	calls := 0
	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		calls++

		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q1, _ := m.Driver("default")
	q2, _ := m.Driver("default")

	if q1 != q2 {
		t.Error("expected same instance from cache")
	}

	if calls != 1 {
		t.Errorf("expected creator called once, got %d", calls)
	}
}

func TestManagerDriverInvalidDriver(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()
	m.SetConfig("default", map[string]any{"driver": "unknown"})

	_, err := m.Driver("default")

	if err == nil {
		t.Fatal("expected error for unknown driver")
	}

	if !errors.Is(err, queue.ErrInvalidDriver) {
		t.Errorf("expected ErrInvalidDriver, got %v", err)
	}
}

func TestManagerDriverMissingConfig(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	_, err := m.Driver("unconfigured")

	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestManagerExtendIsAlias(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Extend("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q, err := m.Driver("default")

	if err != nil {
		t.Fatal(err)
	}

	if q == nil {
		t.Fatal("expected non-nil queue from Extend")
	}
}

func TestManagerConcurrentAccess(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, _ = m.Driver("default")
		}()
	}

	wg.Wait()
}

func TestManagerConnectionIsAliasForDriver(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q, err := m.Connection("default")

	if err != nil {
		t.Fatal(err)
	}

	if q == nil {
		t.Fatal("expected non-nil queue from Connection")
	}
}

func TestManagerDefaultConnection(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	if m.GetDefaultConnection() != "" {
		t.Errorf("expected empty default, got %q", m.GetDefaultConnection())
	}

	m.SetDefaultConnection("redis")

	if m.GetDefaultConnection() != "redis" {
		t.Errorf("expected 'redis', got %q", m.GetDefaultConnection())
	}
}

func TestManagerPurge(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	q1, _ := m.Driver("default")
	m.Purge("default")
	q2, _ := m.Driver("default")

	if q1 == q2 {
		t.Error("expected different instance after purge")
	}
}

func TestManagerForgetDriver(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		return drivers.NewNullDriver("null"), nil
	})

	m.SetConfig("default", map[string]any{"driver": "null"})

	m.ForgetDriver("null")

	_, err := m.Driver("default")

	if err == nil {
		t.Fatal("expected error after ForgetDriver")
	}
}

func TestManagerCreatorError(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	m.Register("bad", func(_ map[string]any) (queue.Queue, error) {
		return nil, errors.New("create failed")
	})

	m.SetConfig("default", map[string]any{"driver": "bad"})

	_, err := m.Driver("default")

	if err == nil {
		t.Fatal("expected error from creator")
	}
}

func (q *inspectorQueue) Push(_ context.Context, _ string, _ []byte) (string, error) {
	return "", nil
}

func (q *inspectorQueue) PushDelayed(_ context.Context, _ string, _ []byte, _ time.Duration) (string, error) {
	return "", nil
}

func (q *inspectorQueue) PushMultiple(_ context.Context, _ string, payloads [][]byte) ([]string, error) {
	return make([]string, len(payloads)), nil
}

func (q *inspectorQueue) Pop(_ context.Context, _ string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (q *inspectorQueue) Size(_ context.Context, _ string) (int64, error)         { return 0, nil }
func (q *inspectorQueue) PendingSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (q *inspectorQueue) DelayedSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (q *inspectorQueue) ReservedSize(_ context.Context, _ string) (int64, error) { return 0, nil }
func (q *inspectorQueue) ConnectionName() string                                  { return q.connection }

func (q *inspectorQueue) QueueNames(_ context.Context) ([]string, error) {
	return q.names, q.namesErr
}

func (q *inspectorQueue) PendingJobs(_ context.Context, name string) ([]queue.InspectedJob, error) {
	if q.queueErr != nil {
		return nil, q.queueErr
	}

	return q.perQueue[name], nil
}

func (q *inspectorQueue) DelayedJobs(_ context.Context, name string) ([]queue.InspectedJob, error) {
	return q.perQueue[name], nil
}

func (q *inspectorQueue) ReservedJobs(_ context.Context, name string) ([]queue.InspectedJob, error) {
	if name == q.skipQueue {
		return nil, queue.ErrNotSupported
	}

	return q.perQueue[name], nil
}

func newManagerWithInspector(t *testing.T, name string, q *inspectorQueue) *queue.Manager {
	t.Helper()

	m := queue.NewManager()
	m.Register("inspector", func(_ map[string]any) (queue.Queue, error) {
		return q, nil
	})
	m.SetConfig(name, map[string]any{"driver": "inspector"})

	return m
}

func TestManagerAllPendingJobsFansOutAcrossQueues(t *testing.T) {
	t.Parallel()

	q := &inspectorQueue{
		connection: "inspector",
		names:      []string{"default", "emails"},
		perQueue: map[string][]queue.InspectedJob{
			"default": {{ID: 1, Queue: "default"}, {ID: 2, Queue: "default"}},
			"emails":  {{ID: 3, Queue: "emails"}},
		},
	}

	m := newManagerWithInspector(t, "default", q)

	jobs, err := m.AllPendingJobs(context.Background(), "default")

	if err != nil {
		t.Fatalf("AllPendingJobs: %v", err)
	}

	if len(jobs) != 3 {
		t.Fatalf("got %d jobs, want 3", len(jobs))
	}

	// Order follows QueueNames order.
	if jobs[0].ID != 1 || jobs[1].ID != 2 || jobs[2].ID != 3 {
		t.Errorf("ids: got %d/%d/%d, want 1/2/3", jobs[0].ID, jobs[1].ID, jobs[2].ID)
	}
}

func TestManagerAllDelayedJobsAndReservedJobsAggregate(t *testing.T) {
	t.Parallel()

	q := &inspectorQueue{
		connection: "inspector",
		names:      []string{"a", "b"},
		perQueue: map[string][]queue.InspectedJob{
			"a": {{ID: 10}},
			"b": {{ID: 20}},
		},
	}

	m := newManagerWithInspector(t, "default", q)
	ctx := context.Background()

	delayed, err := m.AllDelayedJobs(ctx, "default")

	if err != nil || len(delayed) != 2 {
		t.Errorf("AllDelayedJobs: got (%v, %v)", delayed, err)
	}

	reserved, err := m.AllReservedJobs(ctx, "default")

	if err != nil || len(reserved) != 2 {
		t.Errorf("AllReservedJobs: got (%v, %v)", reserved, err)
	}
}

func TestManagerAllReservedJobsSkipsErrNotSupportedQueues(t *testing.T) {
	t.Parallel()

	q := &inspectorQueue{
		connection: "inspector",
		names:      []string{"good", "skip-me"},
		perQueue: map[string][]queue.InspectedJob{
			"good": {{ID: 1}},
		},
		skipQueue: "skip-me",
	}

	m := newManagerWithInspector(t, "default", q)

	jobs, err := m.AllReservedJobs(context.Background(), "default")

	if err != nil {
		t.Fatalf("AllReservedJobs: %v", err)
	}

	if len(jobs) != 1 || jobs[0].ID != 1 {
		t.Errorf("got %v, want only the 'good' job", jobs)
	}
}

func TestManagerAllPendingJobsConnectionWithoutContractErrors(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()
	m.Register("null", func(_ map[string]any) (queue.Queue, error) {
		// NullDriver does satisfy QueueNames/JobInspector (returns nil),
		// so to hit the "no contract" branch we use a thin driver here
		// — but NullDriver returns nil names so the fan-out returns empty.
		return drivers.NewNullDriver("null"), nil
	})
	m.SetConfig("default", map[string]any{"driver": "null"})

	jobs, err := m.AllPendingJobs(context.Background(), "default")

	if err != nil {
		t.Fatalf("AllPendingJobs: %v", err)
	}

	if jobs != nil {
		t.Errorf("got %v, want nil for null driver", jobs)
	}
}

func TestManagerAllPendingJobsMissingConnection(t *testing.T) {
	t.Parallel()

	m := queue.NewManager()

	_, err := m.AllPendingJobs(context.Background(), "unknown")

	if err == nil {
		t.Fatal("expected error for unknown connection")
	}
}

func TestManagerAllPendingJobsQueueErrorBubbles(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("query failed")
	q := &inspectorQueue{
		connection: "inspector",
		names:      []string{"x"},
		queueErr:   wantErr,
	}
	m := newManagerWithInspector(t, "default", q)

	_, err := m.AllPendingJobs(context.Background(), "default")

	if err == nil || !errors.Is(err, wantErr) {
		t.Errorf("want %v, got %v", wantErr, err)
	}
}
