package pennant_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/bedrock/packages/pennant"
)

// ---------------------------------------------------------------------------
// Fake sql driver — serves pre-canned single-column rows.
//
// The DSN encodes the column name and pipe-separated values:
//
//	"col:value1|value2|value3"
//
// An empty DSN produces zero rows.
// ---------------------------------------------------------------------------

type fakeDriver struct{}

type fakeConn struct{ dsn string }

type fakeStmt struct{ dsn string }

// Format: "col:v1|v2|…" or just "v1|v2|…" (col defaults to "col")

type fakeRows struct {
	col    string
	values []string
	pos    int
}

// fakeDB opens a throwaway *sql.DB backed by the registered fake driver.
// It is safe to call multiple times — each call opens a fresh connection.

// openFakeRows returns a *sql.Rows containing the given string values in a
// single column named "name".

// fakeRow returns a *sql.Row that scans the given string when Scan is called.

// fakeEmptyRow returns a *sql.Row whose Scan returns sql.ErrNoRows.

// mockResult satisfies sql.Result.
type mockResult struct{}

// ---------------------------------------------------------------------------
// inMemoryDB — DBExecutor backed by an in-memory slice of rows.
// ---------------------------------------------------------------------------

type dbRow struct {
	name      string
	scope     string
	value     string // JSON-encoded
	createdAt any
	updatedAt any
}

// inMemoryDB simulates the features SQL table entirely in memory.
// It is safe for concurrent use.
type inMemoryDB struct {
	mu             sync.Mutex
	rows           []dbRow
	shouldConflict bool
	conflictCount  int
	conflictsSeen  int
	execCount      int
	queryCount     int
	rowQueryCount  int
}

func init() {
	sql.Register("pennant_fake", &fakeDriver{})
}

func (f *fakeDriver) Open(name string) (driver.Conn, error) {
	return &fakeConn{dsn: name}, nil
}

func (c *fakeConn) Prepare(_ string) (driver.Stmt, error) {
	return &fakeStmt{dsn: c.dsn}, nil
}

func (c *fakeConn) Close() error              { return nil }
func (c *fakeConn) Begin() (driver.Tx, error) { return nil, fmt.Errorf("not supported") }

func (s *fakeStmt) Close() error  { return nil }
func (s *fakeStmt) NumInput() int { return -1 }
func (s *fakeStmt) Exec(_ []driver.Value) (driver.Result, error) {
	return nil, fmt.Errorf("not supported")
}

func (s *fakeStmt) Query(_ []driver.Value) (driver.Rows, error) {
	if s.dsn == "" {
		return &fakeRows{col: "col", values: nil}, nil
	}

	col := "col"
	data := s.dsn

	if idx := strings.Index(s.dsn, ":"); idx >= 0 {
		col = s.dsn[:idx]
		data = s.dsn[idx+1:]
	}

	var values []string

	if data != "" {
		values = strings.Split(data, "|")
	}

	return &fakeRows{col: col, values: values}, nil
}

func (r *fakeRows) Columns() []string { return []string{r.col} }
func (r *fakeRows) Close() error      { return nil }

func (r *fakeRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.values) {
		return io.EOF
	}

	dest[0] = r.values[r.pos]
	r.pos++

	return nil
}

var fakeDBOnce sync.Once
var globalFakeDB *sql.DB

func openFakeDB(dsn string) *sql.DB {
	db, err := sql.Open("pennant_fake", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func openFakeRows(names []string) (*sql.Rows, error) {
	dsn := "name:" + strings.Join(names, "|")

	if len(names) == 0 {
		dsn = ""
	}

	db := openFakeDB(dsn)

	return db.QueryContext(context.Background(), "SELECT name")
}

func fakeRow(value string) *sql.Row {
	db := openFakeDB("value:" + value)

	return db.QueryRowContext(context.Background(), "SELECT value")
}

func fakeEmptyRow() *sql.Row {
	db := openFakeDB("")

	return db.QueryRowContext(context.Background(), "SELECT value")
}

func (mockResult) LastInsertId() (int64, error) { return 0, nil }
func (mockResult) RowsAffected() (int64, error) { return 1, nil }

// ExecContext dispatches INSERT, UPDATE, and DELETE statements.
func (m *inMemoryDB) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.execCount++

	upper := strings.ToUpper(strings.TrimSpace(query))

	switch {
	case strings.HasPrefix(upper, "INSERT"):
		return m.handleInsert(query, args...)
	case strings.HasPrefix(upper, "UPDATE"):
		return m.handleUpdate(args...)
	case strings.HasPrefix(upper, "DELETE"):
		return m.handleDelete(query, args...)
	default:
		return nil, fmt.Errorf("inMemoryDB: unsupported query: %s", query)
	}
}

func (m *inMemoryDB) handleInsert(query string, args ...any) (sql.Result, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("inMemoryDB INSERT: need at least 3 args")
	}

	if m.shouldConflict && m.conflictsSeen < m.conflictCount {
		m.conflictsSeen++

		return nil, fmt.Errorf("unique constraint violation")
	}

	name := fmt.Sprintf("%v", args[0])
	scope := fmt.Sprintf("%v", args[1])
	value := fmt.Sprintf("%v", args[2])

	var timestamp any

	if len(args) > 3 {
		timestamp = args[3]
	}

	isUpsert := strings.Contains(strings.ToUpper(query), "ON CONFLICT")

	for i, r := range m.rows {
		if r.name == name && r.scope == scope {
			if isUpsert {
				m.rows[i].value = value
				m.rows[i].updatedAt = timestamp

				return mockResult{}, nil
			}

			return nil, fmt.Errorf("unique constraint violation")
		}
	}

	m.rows = append(m.rows, dbRow{name: name, scope: scope, value: value, createdAt: timestamp, updatedAt: timestamp})

	return mockResult{}, nil
}

// handleUpdate: UPDATE table SET value=$1, updated_at=$2 WHERE name=$3
func (m *inMemoryDB) handleUpdate(args ...any) (sql.Result, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("inMemoryDB UPDATE: need at least 3 args")
	}

	value := fmt.Sprintf("%v", args[0])
	updatedAt := args[1]
	name := fmt.Sprintf("%v", args[2])

	for i, r := range m.rows {
		if r.name == name {
			m.rows[i].value = value
			m.rows[i].updatedAt = updatedAt
		}
	}

	return mockResult{}, nil
}

func (m *inMemoryDB) handleDelete(query string, args ...any) (sql.Result, error) {
	upper := strings.ToUpper(query)

	// Purge all — no WHERE clause.
	if !strings.Contains(upper, "WHERE") {
		m.rows = m.rows[:0]

		return mockResult{}, nil
	}

	// DELETE WHERE name=$1 AND scope=$2
	if strings.Contains(upper, "AND SCOPE") {
		if len(args) < 2 {
			return nil, fmt.Errorf("inMemoryDB DELETE by scope: need 2 args")
		}

		name := fmt.Sprintf("%v", args[0])
		scope := fmt.Sprintf("%v", args[1])

		kept := m.rows[:0]

		for _, r := range m.rows {
			if !(r.name == name && r.scope == scope) {
				kept = append(kept, r)
			}
		}

		m.rows = kept

		return mockResult{}, nil
	}

	// DELETE WHERE name IN ($1,$2,…)
	names := make(map[string]struct{}, len(args))

	for _, a := range args {
		names[fmt.Sprintf("%v", a)] = struct{}{}
	}

	kept := m.rows[:0]

	for _, r := range m.rows {
		if _, hit := names[r.name]; !hit {
			kept = append(kept, r)
		}
	}

	m.rows = kept

	return mockResult{}, nil
}

// QueryContext handles: SELECT DISTINCT name FROM table
func (m *inMemoryDB) QueryContext(_ context.Context, query string, _ ...any) (*sql.Rows, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.queryCount++

	upper := strings.ToUpper(strings.TrimSpace(query))

	if !strings.Contains(upper, "SELECT DISTINCT NAME") {
		return nil, fmt.Errorf("inMemoryDB: unsupported QueryContext: %s", query)
	}

	seen := make(map[string]struct{})

	var names []string

	for _, r := range m.rows {
		if _, ok := seen[r.name]; !ok {
			seen[r.name] = struct{}{}
			names = append(names, r.name)
		}
	}

	return openFakeRows(names)
}

// QueryRowContext handles: SELECT value FROM table WHERE name=$1 AND scope=$2
func (m *inMemoryDB) QueryRowContext(_ context.Context, _ string, args ...any) *sql.Row {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.rowQueryCount++

	if len(args) < 2 {
		return fakeEmptyRow()
	}

	name := fmt.Sprintf("%v", args[0])
	scope := fmt.Sprintf("%v", args[1])

	for _, r := range m.rows {
		if r.name == name && r.scope == scope {
			return fakeRow(r.value)
		}
	}

	return fakeEmptyRow()
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func newDB() *inMemoryDB {
	return &inMemoryDB{}
}

const testTable = "features"

func TestDatabaseDriver_InterfaceAssertions(t *testing.T) {
	t.Parallel()

	var _ pennant.Driver = (*pennant.DatabaseDriver)(nil)

	var _ pennant.StoredFeaturesLister = (*pennant.DatabaseDriver)(nil)

	var _ pennant.BulkFeatureSetter = (*pennant.DatabaseDriver)(nil)
}

func TestDatabaseDriver_Define_Get(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	drv.Define("dark-mode", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	val, err := drv.Get(ctx, "dark-mode", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != true {
		t.Fatalf("expected true, got %v", val)
	}
}

func TestDatabaseDriver_Get_UndefinedFeature(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	db := newDB()
	drv := pennant.NewDatabaseDriverWithDispatcher(db, testTable, dispatcher)
	ctx := context.Background()

	_, err := drv.Get(ctx, "unknown-feature", nil)

	if !errors.Is(err, pennant.ErrFeatureNotDefined) {
		t.Fatalf("expected ErrFeatureNotDefined, got %v", err)
	}

	if dispatcher.count("UnknownFeatureResolved") != 1 {
		t.Fatalf("expected 1 UnknownFeatureResolved event, got %d",
			dispatcher.count("UnknownFeatureResolved"))
	}
}

func TestDatabaseDriver_Get_CachesInDB(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	calls := 0

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++

		return "variant-a", nil
	})

	// First call invokes resolver and stores result.
	if _, err := drv.Get(ctx, "flag", "user:1"); err != nil {
		t.Fatalf("first Get: %v", err)
	}

	// Second call should find value in DB without invoking resolver again.
	val, err := drv.Get(ctx, "flag", "user:1")

	if err != nil {
		t.Fatalf("second Get: %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected resolver called once, got %d", calls)
	}

	if val != "variant-a" {
		t.Fatalf("expected variant-a, got %v", val)
	}
}

func TestDatabaseDriver_Set_BypassesResolver(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	if err := drv.Set(ctx, "flag", nil, false); err != nil {
		t.Fatalf("Set: %v", err)
	}

	val, err := drv.Get(ctx, "flag", nil)

	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if val != false {
		t.Fatalf("expected false (stored), got %v", val)
	}
}

func TestDatabaseDriver_SetAll(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	entries := []pennant.FeatureEntry{
		{Feature: "flag-a", Scope: "user:1", Value: true},
		{Feature: "flag-b", Scope: "user:1", Value: "variant"},
	}

	if err := drv.SetAll(ctx, entries); err != nil {
		t.Fatalf("SetAll: %v", err)
	}

	valA, err := drv.Get(ctx, "flag-a", "user:1")

	if err != nil {
		t.Fatalf("Get flag-a: %v", err)
	}

	if valA != true {
		t.Fatalf("expected true, got %v", valA)
	}

	valB, err := drv.Get(ctx, "flag-b", "user:1")

	if err != nil {
		t.Fatalf("Get flag-b: %v", err)
	}

	if valB != "variant" {
		t.Fatalf("expected variant, got %v", valB)
	}
}

func TestDatabaseDriver_SetForAllScopes(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) { return true, nil })

	// Resolve for two scopes so rows exist in DB.
	if _, err := drv.Get(ctx, "flag", "user:1"); err != nil {
		t.Fatalf("Get user:1: %v", err)
	}

	if _, err := drv.Get(ctx, "flag", "user:2"); err != nil {
		t.Fatalf("Get user:2: %v", err)
	}

	if err := drv.SetForAllScopes(ctx, "flag", false); err != nil {
		t.Fatalf("SetForAllScopes: %v", err)
	}

	for _, scope := range []any{"user:1", "user:2"} {
		val, err := drv.Get(ctx, "flag", scope)

		if err != nil {
			t.Fatalf("Get %v: %v", scope, err)
		}

		if val != false {
			t.Fatalf("scope %v: expected false, got %v", scope, val)
		}
	}
}

func TestDatabaseDriver_Delete(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	calls := 0

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++

		return true, nil
	})

	if _, err := drv.Get(ctx, "flag", nil); err != nil {
		t.Fatalf("first Get: %v", err)
	}

	if err := drv.Delete(ctx, "flag", nil); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// After deletion, Get should re-invoke the resolver.
	if _, err := drv.Get(ctx, "flag", nil); err != nil {
		t.Fatalf("second Get: %v", err)
	}

	if calls != 2 {
		t.Fatalf("expected resolver called twice (after delete), got %d", calls)
	}
}

func TestDatabaseDriver_Purge_All(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	drv.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	drv.Define("flag-b", func(_ context.Context, _ any) (any, error) { return true, nil })

	drv.Get(ctx, "flag-a", nil) //nolint:errcheck
	drv.Get(ctx, "flag-b", nil) //nolint:errcheck

	if err := drv.Purge(ctx, nil); err != nil {
		t.Fatalf("Purge: %v", err)
	}

	stored, err := drv.Stored(ctx)

	if err != nil {
		t.Fatalf("Stored: %v", err)
	}

	if len(stored) != 0 {
		t.Fatalf("expected no stored features after purge all, got %v", stored)
	}
}

func TestDatabaseDriver_Purge_Specific(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	drv.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	drv.Define("flag-b", func(_ context.Context, _ any) (any, error) { return true, nil })

	drv.Get(ctx, "flag-a", nil) //nolint:errcheck
	drv.Get(ctx, "flag-b", nil) //nolint:errcheck

	if err := drv.Purge(ctx, []string{"flag-a"}); err != nil {
		t.Fatalf("Purge: %v", err)
	}

	stored, err := drv.Stored(ctx)

	if err != nil {
		t.Fatalf("Stored: %v", err)
	}

	if len(stored) != 1 || stored[0] != "flag-b" {
		t.Fatalf("expected only flag-b stored, got %v", stored)
	}
}

func TestDatabaseDriver_Purge_EmptySlice_IsNoOp(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) { return true, nil })

	drv.Get(ctx, "flag", nil) //nolint:errcheck

	if err := drv.Purge(ctx, []string{}); err != nil {
		t.Fatalf("Purge: %v", err)
	}

	stored, err := drv.Stored(ctx)

	if err != nil {
		t.Fatalf("Stored: %v", err)
	}

	if len(stored) != 1 {
		t.Fatalf("expected 1 stored feature after empty-slice purge, got %v", stored)
	}
}

func TestDatabaseDriver_Stored(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	drv.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	drv.Define("flag-b", func(_ context.Context, _ any) (any, error) { return true, nil })

	drv.Get(ctx, "flag-a", nil) //nolint:errcheck

	stored, err := drv.Stored(ctx)

	if err != nil {
		t.Fatalf("Stored: %v", err)
	}

	if len(stored) != 1 || stored[0] != "flag-a" {
		t.Fatalf("expected [flag-a], got %v", stored)
	}
}

func TestDatabaseDriver_JSON_Bool_RoundTrip(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	if err := drv.Set(ctx, "feature", "scope1", false); err != nil {
		t.Fatalf("Set false: %v", err)
	}

	val, err := drv.Get(ctx, "feature", "scope1")

	if err != nil {
		t.Fatalf("Get false: %v", err)
	}

	b, ok := val.(bool)

	if !ok {
		t.Fatalf("expected bool, got %T (%v)", val, val)
	}

	if b != false {
		t.Fatalf("expected false, got %v", b)
	}

	if err := drv.Set(ctx, "feature2", "scope1", true); err != nil {
		t.Fatalf("Set true: %v", err)
	}

	val2, err := drv.Get(ctx, "feature2", "scope1")

	if err != nil {
		t.Fatalf("Get true: %v", err)
	}

	b2, ok := val2.(bool)

	if !ok {
		t.Fatalf("expected bool, got %T (%v)", val2, val2)
	}

	if !b2 {
		t.Fatalf("expected true, got %v", b2)
	}
}

func TestDatabaseDriver_JSON_String_RoundTrip(t *testing.T) {
	t.Parallel()

	db := newDB()
	drv := pennant.NewDatabaseDriver(db, testTable)
	ctx := context.Background()

	if err := drv.Set(ctx, "theme", "user:42", "dark"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	val, err := drv.Get(ctx, "theme", "user:42")

	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	s, ok := val.(string)

	if !ok {
		t.Fatalf("expected string, got %T (%v)", val, val)
	}

	if s != "dark" {
		t.Fatalf("expected dark, got %v", s)
	}
}
