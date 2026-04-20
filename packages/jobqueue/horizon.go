package jobqueue

import (
	"context"
	"sort"
	"sync"
	"time"
)

// QueueStatus is a point-in-time view of one queue.
type QueueStatus struct {
	Name       string
	Pending    int
	Reserved   int
	Processing int
	Failed     int
	Throughput int
	Wait       time.Duration
	Paused     bool
}

// Snapshot groups queue statuses captured at the same time.
type Snapshot struct {
	GeneratedAt time.Time
	Queues      []QueueStatus
}

// QueueSource captures queue state from a backend.
type QueueSource interface {
	Snapshot(ctx context.Context) ([]QueueStatus, error)
}

// QueueSourceFunc adapts a function into a QueueSource.
type QueueSourceFunc func(ctx context.Context) ([]QueueStatus, error)

// Snapshot captures queue state by calling f.

// Repository stores JobQueue snapshots.
type Repository interface {
	Record(ctx context.Context, snapshot Snapshot) error
	Latest(ctx context.Context) (Snapshot, error)
}

// Monitor captures queue snapshots and records them in a Repository.
type Monitor struct {
	repository Repository
	sources    []QueueSource
	now        func() time.Time
}

// NewMonitor creates a queue monitor.

// Capture gathers queue state from all sources, stores it, and returns it.

// InMemoryRepository stores snapshots in memory.
type InMemoryRepository struct {
	mu        sync.RWMutex
	snapshots []Snapshot
}

func (f QueueSourceFunc) Snapshot(ctx context.Context) ([]QueueStatus, error) {
	return f(ctx)
}

func NewMonitor(repository Repository, sources ...QueueSource) *Monitor {
	return &Monitor{
		repository: repository,
		sources:    sources,
		now:        time.Now,
	}
}

func (m *Monitor) Capture(ctx context.Context) (Snapshot, error) {
	queues := make([]QueueStatus, 0)

	for _, source := range m.sources {
		statuses, err := source.Snapshot(ctx)

		if err != nil {
			return Snapshot{}, err
		}

		queues = append(queues, statuses...)
	}

	sort.SliceStable(queues, func(i, j int) bool {
		return queues[i].Name < queues[j].Name
	})

	snapshot := Snapshot{
		GeneratedAt: m.now().UTC(),
		Queues:      queues,
	}

	if m.repository != nil {
		if err := m.repository.Record(ctx, snapshot); err != nil {
			return Snapshot{}, err
		}
	}

	return snapshot, nil
}

// NewInMemoryRepository creates an empty repository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

// Record appends a snapshot.
func (r *InMemoryRepository) Record(ctx context.Context, snapshot Snapshot) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	r.snapshots = append(r.snapshots, cloneSnapshot(snapshot))

	return nil
}

// Latest returns the newest snapshot.
func (r *InMemoryRepository) Latest(ctx context.Context) (Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return Snapshot{}, err
	}

	r.mu.RLock()

	defer r.mu.RUnlock()

	if len(r.snapshots) == 0 {
		return Snapshot{}, ErrNoSnapshots
	}

	return cloneSnapshot(r.snapshots[len(r.snapshots)-1]), nil
}

func cloneSnapshot(snapshot Snapshot) Snapshot {
	queues := make([]QueueStatus, len(snapshot.Queues))
	copy(queues, snapshot.Queues)
	snapshot.Queues = queues

	return snapshot
}
