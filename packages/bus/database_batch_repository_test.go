package bus_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/bus"
)

// mockDBExecutor implements bus.DBExecutor for unit testing.
type mockDBExecutor struct {
	mu          sync.Mutex
	execCalls   []mockDBExecCall
	execErr     error
	execResult  sql.Result
	queryRow    *mockRow
	queryRowErr error
	queryRows   *mockRows
	queryErr    error
}

type mockDBExecCall struct {
	Query string
	Args  []any
}

type mockResult struct {
	lastID       int64
	rowsAffected int64
}

// mockRow provides a scannable row for testing.
type mockRow struct {
	values []any
	err    error
}

// mockRows provides scannable rows for testing.
type mockRows struct {
	rows    [][]any
	cursor  int
	columns []string
}

// failed_job_ids is arg index 5.

// options is arg index 6.

// First call returns 1000 rows affected, second returns 0.

// The query should reference job_batches.

// mockDynamicResult lets RowsAffected return different values per call.
type mockDynamicResult struct {
	fn func() int64
}

func newMockDBExecutor() *mockDBExecutor {
	return &mockDBExecutor{
		execResult: &mockResult{rowsAffected: 1},
	}
}

func (db *mockDBExecutor) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	db.mu.Lock()

	defer db.mu.Unlock()

	db.execCalls = append(db.execCalls, mockDBExecCall{Query: query, Args: args})

	return db.execResult, db.execErr
}

func (db *mockDBExecutor) QueryContext(_ context.Context, query string, args ...any) (*sql.Rows, error) {
	return nil, db.queryErr
}

func (db *mockDBExecutor) QueryRowContext(_ context.Context, query string, args ...any) *sql.Row {
	return nil
}

func (r *mockResult) LastInsertId() (int64, error) { return r.lastID, nil }
func (r *mockResult) RowsAffected() (int64, error) { return r.rowsAffected, nil }

func TestDatabaseBatchRepositoryStore(t *testing.T) {
	db := newMockDBExecutor()
	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	batch := &bus.Batch{
		ID:          "batch-1",
		Name:        "test-batch",
		TotalJobs:   5,
		PendingJobs: 5,
		Options:     map[string]any{},
		CreatedAt:   time.Now(),
	}

	err := repo.Store(context.Background(), batch)

	if err != nil {
		t.Fatal(err)
	}

	db.mu.Lock()

	defer db.mu.Unlock()

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}

	call := db.execCalls[0]

	if call.Args[0] != "batch-1" {
		t.Errorf("expected batch ID 'batch-1', got %v", call.Args[0])
	}

	if call.Args[1] != "test-batch" {
		t.Errorf("expected batch name 'test-batch', got %v", call.Args[1])
	}
}

func TestDatabaseBatchRepositoryStoreSerializesJSON(t *testing.T) {
	db := newMockDBExecutor()
	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	batch := &bus.Batch{
		ID:           "batch-2",
		Name:         "json-test",
		FailedJobIDs: []string{"j1", "j2"},
		Options:      map[string]any{"allowFailures": true},
		CreatedAt:    time.Now(),
	}

	err := repo.Store(context.Background(), batch)

	if err != nil {
		t.Fatal(err)
	}

	db.mu.Lock()

	defer db.mu.Unlock()

	call := db.execCalls[0]

	failedJSON := call.Args[5].(string)

	var failedIDs []string

	if err := json.Unmarshal([]byte(failedJSON), &failedIDs); err != nil {
		t.Fatalf("failed to unmarshal failed_job_ids: %v", err)
	}

	if len(failedIDs) != 2 || failedIDs[0] != "j1" {
		t.Errorf("expected [j1, j2], got %v", failedIDs)
	}

	optsJSON := call.Args[6].(string)

	var opts map[string]any

	if err := json.Unmarshal([]byte(optsJSON), &opts); err != nil {
		t.Fatalf("failed to unmarshal options: %v", err)
	}

	if opts["allowFailures"] != true {
		t.Errorf("expected allowFailures true, got %v", opts["allowFailures"])
	}
}

func TestDatabaseBatchRepositoryStoreError(t *testing.T) {
	db := newMockDBExecutor()
	db.execErr = fmt.Errorf("db error")

	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	batch := &bus.Batch{ID: "batch-err", Options: map[string]any{}, CreatedAt: time.Now()}

	err := repo.Store(context.Background(), batch)

	if err == nil {
		t.Error("expected error from Store")
	}
}

func TestDatabaseBatchRepositoryIncrementTotalJobs(t *testing.T) {
	db := newMockDBExecutor()
	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	err := repo.IncrementTotalJobs(context.Background(), "batch-1", 3)

	if err != nil {
		t.Fatal(err)
	}

	db.mu.Lock()

	defer db.mu.Unlock()

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}

	if db.execCalls[0].Args[0] != 3 {
		t.Errorf("expected amount 3, got %v", db.execCalls[0].Args[0])
	}
}

func TestDatabaseBatchRepositoryMarkAsFinished(t *testing.T) {
	db := newMockDBExecutor()
	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	err := repo.MarkAsFinished(context.Background(), "batch-1")

	if err != nil {
		t.Fatal(err)
	}

	db.mu.Lock()

	defer db.mu.Unlock()

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}

	if db.execCalls[0].Args[1] != "batch-1" {
		t.Errorf("expected batch ID 'batch-1', got %v", db.execCalls[0].Args[1])
	}
}

func TestDatabaseBatchRepositoryCancel(t *testing.T) {
	db := newMockDBExecutor()
	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	err := repo.Cancel(context.Background(), "batch-1")

	if err != nil {
		t.Fatal(err)
	}

	db.mu.Lock()

	defer db.mu.Unlock()

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}
}

func TestDatabaseBatchRepositoryDelete(t *testing.T) {
	db := newMockDBExecutor()
	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	err := repo.Delete(context.Background(), "batch-1")

	if err != nil {
		t.Fatal(err)
	}

	db.mu.Lock()

	defer db.mu.Unlock()

	if len(db.execCalls) != 1 {
		t.Fatalf("expected 1 exec call, got %d", len(db.execCalls))
	}
}

func TestDatabaseBatchRepositoryPrune(t *testing.T) {
	db := newMockDBExecutor()

	callCount := 0
	db.execResult = &mockDynamicResult{fn: func() int64 {
		callCount++

		if callCount == 1 {
			return 1000
		}

		return 0
	}}

	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	total, err := repo.Prune(context.Background(), time.Now().Add(-24*time.Hour))

	if err != nil {
		t.Fatal(err)
	}

	if total != 1000 {
		t.Errorf("expected 1000 pruned, got %d", total)
	}
}

func TestDatabaseBatchRepositoryPruneCancelled(t *testing.T) {
	db := newMockDBExecutor()
	db.execResult = &mockResult{rowsAffected: 0}

	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	total, err := repo.PruneCancelled(context.Background(), time.Now())

	if err != nil {
		t.Fatal(err)
	}

	if total != 0 {
		t.Errorf("expected 0 pruned, got %d", total)
	}
}

func TestDatabaseBatchRepositoryPruneUnfinished(t *testing.T) {
	db := newMockDBExecutor()
	db.execResult = &mockResult{rowsAffected: 0}

	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	total, err := repo.PruneUnfinished(context.Background(), time.Now())

	if err != nil {
		t.Fatal(err)
	}

	if total != 0 {
		t.Errorf("expected 0 pruned, got %d", total)
	}
}

func TestDatabaseBatchRepositoryTransactionWithoutTransactor(t *testing.T) {
	db := newMockDBExecutor()
	repo := bus.NewDatabaseBatchRepository(db, "job_batches")

	called := false
	err := repo.Transaction(context.Background(), func(r bus.BatchRepository) error {
		called = true

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("expected fn to be called")
	}
}

func TestDatabaseBatchRepositoryDefaultTable(t *testing.T) {
	db := newMockDBExecutor()
	repo := bus.NewDatabaseBatchRepository(db, "")

	batch := &bus.Batch{ID: "b1", Options: map[string]any{}, CreatedAt: time.Now()}
	_ = repo.Store(context.Background(), batch)

	db.mu.Lock()

	defer db.mu.Unlock()

	if len(db.execCalls) == 0 {
		t.Fatal("expected at least 1 call")
	}

	query := db.execCalls[0].Query

	if len(query) < 20 {
		t.Error("query too short")
	}
}

func (r *mockDynamicResult) LastInsertId() (int64, error) { return 0, nil }
func (r *mockDynamicResult) RowsAffected() (int64, error) { return r.fn(), nil }

// Ensure DatabaseBatchRepository satisfies the interfaces.
var _ bus.BatchRepository = (*bus.DatabaseBatchRepository)(nil)

// Use a driver.Value type assertion to suppress unused import warning.
var _ driver.Value = driver.Value(nil)
