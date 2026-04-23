// Package api exposes the HTTP surface of the Bedrock Horizon dashboard. It
// ports Laravel Horizon's Vue-facing JSON controllers (BatchesController,
// DashboardStatsController, MasterSupervisorController, MonitoringController,
// metrics and job retrieval endpoints) on top of the Go primitives shipped in
// packages/horizon.
package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/bedrock/packages/horizon"
)

// MasterSupervisor describes one supervisor master process reported by the
// dashboard. Laravel Horizon exposes the same shape via MasterSupervisor::all().
type MasterSupervisor struct {
	Name        string   `json:"name"`
	PID         int      `json:"pid"`
	Status      string   `json:"status"`
	Supervisors []string `json:"supervisors"`
}

// Batch is the minimal batch record returned by /api/batches. It matches the
// JSON shape consumed by Horizon's Vue dashboard.
type Batch struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	TotalJobs     int        `json:"totalJobs"`
	PendingJobs   int        `json:"pendingJobs"`
	FailedJobs    int        `json:"failedJobs"`
	ProcessedJobs int        `json:"processedJobs"`
	CreatedAt     time.Time  `json:"createdAt"`
	FinishedAt    *time.Time `json:"finishedAt,omitempty"`
	Cancelled     bool       `json:"cancelled"`
}

// BatchStore retains Horizon job batches. Callers plug in their own
// persistence; the dashboard reads through this interface.
type BatchStore interface {
	Search(name string, afterID string, limit int) []Batch
}

// MasterSupervisorSource returns the master supervisors that should appear in
// the dashboard supervisor list.
type MasterSupervisorSource interface {
	All() []MasterSupervisor
}

// Options bundles the collaborators required to mount the dashboard handler.
// Any zero-valued field falls back to a deterministic empty source so partial
// deployments still render a working dashboard.
type Options struct {
	Auth         AuthCallback
	Repository   horizon.Repository
	Jobs         *horizon.JobRepository
	Monitoring   *horizon.MonitoringRepository
	Metrics      *horizon.MetricsRepository
	Supervisors  MasterSupervisorSource
	Batches      BatchStore
	Now          func() time.Time
	SilencedJobs []horizon.JobRecord
}

// NewHandler builds the dashboard HTTP handler with all JSON routes mounted
// under /api and guarded by the auth middleware.

type emptySupervisors struct{}

type emptyBatches struct{}

// SliceSupervisors adapts a slice of MasterSupervisor values to the
// MasterSupervisorSource interface.
type SliceSupervisors struct {
	mu    sync.RWMutex
	items []MasterSupervisor
}

// NewSliceSupervisors returns a supervisor source backed by the supplied slice.

// All returns a defensive copy of the stored supervisors.

// Replace swaps the stored supervisors — useful in tests.

// InMemoryBatches is a trivial BatchStore used by tests and local runs.
type InMemoryBatches struct {
	mu    sync.RWMutex
	items []Batch
}

func NewHandler(opts Options) http.Handler {
	opts = resolveOptions(opts)

	mux := http.NewServeMux()
	registerRoutes(mux, opts)

	return RequireAuth(opts.Auth)(mux)
}

func resolveOptions(opts Options) Options {
	if opts.Auth == nil {
		opts.Auth = AllowAll
	}

	if opts.Repository == nil {
		opts.Repository = horizon.NewInMemoryRepository()
	}

	if opts.Jobs == nil {
		opts.Jobs = horizon.NewJobRepository(opts.Now)
	}

	if opts.Monitoring == nil {
		opts.Monitoring = horizon.NewMonitoringRepository()
	}

	if opts.Metrics == nil {
		opts.Metrics = horizon.NewMetricsRepository()
	}

	if opts.Supervisors == nil {
		opts.Supervisors = emptySupervisors{}
	}

	if opts.Batches == nil {
		opts.Batches = emptyBatches{}
	}

	if opts.Now == nil {
		opts.Now = time.Now
	}

	return opts
}

func (emptySupervisors) All() []MasterSupervisor { return nil }

func (emptyBatches) Search(string, string, int) []Batch { return nil }

func NewSliceSupervisors(items []MasterSupervisor) *SliceSupervisors {
	return &SliceSupervisors{items: append([]MasterSupervisor(nil), items...)}
}

func (s *SliceSupervisors) All() []MasterSupervisor {
	s.mu.RLock()

	defer s.mu.RUnlock()

	out := make([]MasterSupervisor, len(s.items))
	copy(out, s.items)

	return out
}

func (s *SliceSupervisors) Replace(items []MasterSupervisor) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.items = append([]MasterSupervisor(nil), items...)
}

// NewInMemoryBatches returns an empty in-memory batch store.
func NewInMemoryBatches() *InMemoryBatches { return &InMemoryBatches{} }

// Store appends a batch to the store.
func (b *InMemoryBatches) Store(batch Batch) {
	b.mu.Lock()

	defer b.mu.Unlock()

	b.items = append(b.items, batch)
}

// Search ports Laravel Horizon's BatchRepository::getRecentUnfinished +
// searchBy helpers: matches are case-insensitive substring matches on the
// batch name, cursor pagination is ID-based, and wildcard characters in the
// query are treated as literals (LIKE wildcard escaping).
func (b *InMemoryBatches) Search(name string, afterID string, limit int) []Batch {
	b.mu.RLock()

	defer b.mu.RUnlock()

	matches := make([]Batch, 0, len(b.items))
	needle := normaliseBatchQuery(name)

	for _, batch := range b.items {
		if needle != "" && !containsFold(batch.Name, needle) {
			continue
		}

		matches = append(matches, batch)
	}

	if afterID != "" {
		cursor := -1

		for i, batch := range matches {
			if batch.ID == afterID {
				cursor = i

				break
			}
		}

		if cursor >= 0 && cursor+1 < len(matches) {
			matches = matches[cursor+1:]
		} else if cursor >= 0 {
			matches = nil
		}
	}

	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}

	return append([]Batch(nil), matches...)
}
