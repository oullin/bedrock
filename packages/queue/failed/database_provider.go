package failed

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Ref: @bedrock/code-0255
// in a failed_jobs table keyed by an auto-incrementing integer id.
//
// The provider works through a small unexported store interface. The
// default constructor wires an in-memory store (which is what the
// PHPUnit sqlite tests exercise via the port); production callers can
// plug in a SQL-backed store with the same contract. Keeping the
// store contract internal lets the provider mirror the upstream observable
// semantics without embedding a SQL query engine in Go tests.
type DatabaseFailedJobProvider struct {
	store intStore
	table string
	now   func() time.Time
}

// NewDatabaseFailedJobProvider returns a provider backed by the default
// in-memory store. Production callers who need real SQL should wrap
// their drivers.DBExecer in a bespoke store implementing intStore; see
// package docs for the contract.

// Log implements FailedJobProvider.

// IDs implements FailedJobProvider.

// All implements FailedJobProvider.

// Find implements FailedJobProvider.

// Forget implements FailedJobProvider.

// Flush implements FailedJobProvider.

// Count implements Countable.

// Prune implements Prunable and returns the number of deleted rows.

// --- shared store contract ---

type record struct {
	ID         int64
	UUID       string
	Connection string
	Queue      string
	Payload    string
	Exception  string
	FailedAt   time.Time
}

// intStore is the persistence contract the integer-keyed database
// provider relies on. The default implementation is in-memory; see
// newIntMemoryStore.
type intStore interface {
	Insert(ctx context.Context, r record) (string, error)
	IDs(ctx context.Context, queueFilter string) []string
	All(ctx context.Context) []FailedJob
	Find(ctx context.Context, id string) *FailedJob
	Forget(ctx context.Context, id string) bool
	Flush(ctx context.Context, hours int, now time.Time)
	Prune(ctx context.Context, before time.Time) int64
	Count(ctx context.Context, connection, queueFilter string) int64
}

// memoryStore is the shared in-process backing for both integer and
// uuid keyed providers. The two providers wrap it with thin adapters
// that decide how ids are produced and matched.
type memoryStore struct {
	mu     sync.Mutex
	rows   []record
	nextID int64
	uuid   bool
}

func NewDatabaseFailedJobProvider(table string) *DatabaseFailedJobProvider {
	if table == "" {
		table = "failed_jobs"
	}

	return &DatabaseFailedJobProvider{
		store: newIntMemoryStore(),
		table: table,
		now:   time.Now,
	}
}

func (p *DatabaseFailedJobProvider) Log(connection, queue, payload string, exception error) (string, error) {
	return p.store.Insert(context.Background(), record{
		Connection: connection,
		Queue:      queue,
		Payload:    payload,
		Exception:  errString(exception),
		FailedAt:   p.now(),
	})
}

func (p *DatabaseFailedJobProvider) IDs(queueFilter string) ([]string, error) {
	return p.store.IDs(context.Background(), queueFilter), nil
}

func (p *DatabaseFailedJobProvider) All() ([]FailedJob, error) {
	return p.store.All(context.Background()), nil
}

func (p *DatabaseFailedJobProvider) Find(id string) (*FailedJob, error) {
	return p.store.Find(context.Background(), id), nil
}

func (p *DatabaseFailedJobProvider) Forget(id string) (bool, error) {
	return p.store.Forget(context.Background(), id), nil
}

func (p *DatabaseFailedJobProvider) Flush(hours int) error {
	p.store.Flush(context.Background(), hours, p.now())

	return nil
}

func (p *DatabaseFailedJobProvider) Count(connection, queueFilter string) (int64, error) {
	return p.store.Count(context.Background(), connection, queueFilter), nil
}

func (p *DatabaseFailedJobProvider) Prune(before time.Time) (int64, error) {
	return p.store.Prune(context.Background(), before), nil
}

func newIntMemoryStore() *memoryStore  { return &memoryStore{} }
func newUUIDMemoryStore() *memoryStore { return &memoryStore{uuid: true} }

func (s *memoryStore) Insert(_ context.Context, r record) (string, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.nextID++
	r.ID = s.nextID
	s.rows = append(s.rows, r)

	if s.uuid {
		return r.UUID, nil
	}

	return strconv.FormatInt(r.ID, 10), nil
}

func (s *memoryStore) sorted() []record {
	out := make([]record, len(s.rows))
	copy(out, s.rows)

	sort.SliceStable(out, func(i, j int) bool { return out[i].ID > out[j].ID })

	return out
}

func (s *memoryStore) IDs(_ context.Context, queueFilter string) []string {
	s.mu.Lock()

	defer s.mu.Unlock()

	rows := s.sorted()
	out := []string{}

	for _, r := range rows {
		if queueFilter != "" && r.Queue != queueFilter {
			continue
		}

		out = append(out, s.idOf(r))
	}

	return out
}

func (s *memoryStore) idOf(r record) string {
	if s.uuid {
		return r.UUID
	}

	return strconv.FormatInt(r.ID, 10)
}

func (s *memoryStore) All(_ context.Context) []FailedJob {
	s.mu.Lock()

	defer s.mu.Unlock()

	rows := s.sorted()
	out := make([]FailedJob, 0, len(rows))

	for _, r := range rows {
		out = append(out, s.toFailed(r))
	}

	return out
}

func (s *memoryStore) toFailed(r record) FailedJob {
	return FailedJob{
		ID:         s.idOf(r),
		UUID:       r.UUID,
		Connection: r.Connection,
		Queue:      r.Queue,
		Payload:    r.Payload,
		Exception:  r.Exception,
		FailedAt:   r.FailedAt,
	}
}

func (s *memoryStore) matchID(r record, id string) bool {
	return s.idOf(r) == id
}

func (s *memoryStore) Find(_ context.Context, id string) *FailedJob {
	s.mu.Lock()

	defer s.mu.Unlock()

	for _, r := range s.rows {
		if s.matchID(r, id) {
			fj := s.toFailed(r)

			return &fj
		}
	}

	return nil
}

func (s *memoryStore) Forget(_ context.Context, id string) bool {
	s.mu.Lock()

	defer s.mu.Unlock()

	for i, r := range s.rows {
		if s.matchID(r, id) {
			s.rows = append(s.rows[:i], s.rows[i+1:]...)

			return true
		}
	}

	return false
}

func (s *memoryStore) Flush(_ context.Context, hours int, now time.Time) {
	s.mu.Lock()

	defer s.mu.Unlock()

	if hours <= 0 {
		s.rows = nil

		return
	}
	// the upstream flush deletes rows where failed_at <= now - hours.
	cutoff := now.Add(-time.Duration(hours) * time.Hour)
	kept := make([]record, 0, len(s.rows))

	for _, r := range s.rows {
		if r.FailedAt.After(cutoff) {
			kept = append(kept, r)
		}
	}

	s.rows = kept
}

func (s *memoryStore) Prune(_ context.Context, before time.Time) int64 {
	s.mu.Lock()

	defer s.mu.Unlock()

	var deleted int64
	kept := make([]record, 0, len(s.rows))

	for _, r := range s.rows {
		// Upstream prune: deletes rows where failed_at < $before.
		if r.FailedAt.Before(before) {
			deleted++

			continue
		}

		kept = append(kept, r)
	}

	s.rows = kept

	return deleted
}

func (s *memoryStore) Count(_ context.Context, connection, queueFilter string) int64 {
	s.mu.Lock()

	defer s.mu.Unlock()

	var n int64

	for _, r := range s.rows {
		if connection != "" && r.Connection != connection {
			continue
		}

		if queueFilter != "" && r.Queue != queueFilter {
			continue
		}

		n++
	}

	return n
}

func errString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
