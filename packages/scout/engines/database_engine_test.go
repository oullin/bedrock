package engines_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	dbcontract "github.com/bedrock/packages/contracts/database"
	"github.com/bedrock/packages/scout"
	"github.com/bedrock/packages/scout/engines"
)

// mockConnection is a mock database connection for testing.
type mockConnection struct {
	driver     string
	name       string
	selectRows []map[string]any
	selectErr  error
	lastSQL    string
	lastBinds  []any
}

// Return count for aggregate queries.

// mockResolver is a mock connection resolver.
type mockResolver struct {
	conn        dbcontract.Connection
	defaultConn string
}

// Verify the SQL includes the where clause.

// Verify pagination SQL.

// Verify MySQL full-text syntax.

// Verify PostgreSQL full-text syntax.

// Verify SQLite LIKE fallback.

// errorResolver always returns an error.
type errorResolver struct {
	err error
}

func (c *mockConnection) Table(_ context.Context, table string, _ ...string) any {
	return nil
}

func (c *mockConnection) Raw(value string) dbcontract.Expression { return nil }

func (c *mockConnection) SelectOne(ctx context.Context, query string, bindings ...any) (map[string]any, error) {
	c.lastSQL = query
	c.lastBinds = bindings

	if strings.Contains(query, "count(*)") {
		return map[string]any{"aggregate": int64(len(c.selectRows))}, nil
	}

	if len(c.selectRows) > 0 {
		return c.selectRows[0], c.selectErr
	}

	return nil, c.selectErr
}

func (c *mockConnection) Select(ctx context.Context, query string, bindings ...any) ([]map[string]any, error) {
	c.lastSQL = query
	c.lastBinds = bindings

	return c.selectRows, c.selectErr
}

func (c *mockConnection) Insert(_ context.Context, _ string, _ ...any) (bool, error) {
	return true, nil
}
func (c *mockConnection) Update(_ context.Context, _ string, _ ...any) (int64, error) { return 0, nil }
func (c *mockConnection) Delete(_ context.Context, _ string, _ ...any) (int64, error) { return 0, nil }
func (c *mockConnection) Statement(_ context.Context, _ string, _ ...any) (bool, error) {
	return true, nil
}
func (c *mockConnection) AffectingStatement(_ context.Context, _ string, _ ...any) (int64, error) {
	return 0, nil
}
func (c *mockConnection) Unprepared(_ context.Context, _ string) (bool, error) { return true, nil }
func (c *mockConnection) PrepareBindings(bindings []any) []any                 { return bindings }
func (c *mockConnection) Transaction(_ context.Context, fn func(dbcontract.Connection) error, _ ...int) error {
	return fn(c)
}
func (c *mockConnection) BeginTransaction(_ context.Context) error   { return nil }
func (c *mockConnection) Commit(_ context.Context) error             { return nil }
func (c *mockConnection) Rollback(_ context.Context, _ ...int) error { return nil }
func (c *mockConnection) TransactionLevel() int                      { return 0 }
func (c *mockConnection) AfterCommit(_ func())                       {}
func (c *mockConnection) GetTablePrefix() string                     { return "" }
func (c *mockConnection) GetDatabaseName() string                    { return "test" }
func (c *mockConnection) GetDriverName() string                      { return c.driver }
func (c *mockConnection) GetName() string                            { return c.name }
func (c *mockConnection) GetConfig(_ string) any                     { return nil }

func (r *mockResolver) Connection(_ context.Context, name ...string) (dbcontract.Connection, error) {
	return r.conn, nil
}

func (r *mockResolver) GetDefaultConnection() string     { return r.defaultConn }
func (r *mockResolver) SetDefaultConnection(name string) { r.defaultConn = name }

func TestDatabaseEngineSearch(t *testing.T) {
	t.Parallel()
	// DatabaseEngineTest::test_it_can_retrieve_results

	conn := &mockConnection{
		driver: "mysql",
		selectRows: []map[string]any{
			{"id": int64(1), "title": "Hello World"},
			{"id": int64(2), "title": "Hello Go"},
		},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	b := scout.NewBuilder(model, "Hello")

	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	ids := e.MapIds(result)

	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(ids))
	}
}

func TestDatabaseEngineSearchEmptyQueryDoesNotAddSearchWhereClauses(t *testing.T) {
	t.Parallel()
	// DatabaseEngineTest::test_it_can_retrieve_results_with_empty_search
	// DatabaseEngineTest::test_it_does_not_add_search_where_clauses_with_empty_search

	conn := &mockConnection{
		driver: "mysql",
		selectRows: []map[string]any{
			{"id": int64(1), "title": "Hello World"},
			{"id": int64(2), "title": "Hello Go"},
		},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	b := scout.NewBuilder(model, "")

	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total := e.GetTotalCount(result); total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	if conn.lastSQL != "select * from posts" {
		t.Fatalf("expected SQL to omit search where clauses, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEngineSearchWithWheres(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver: "mysql",
		selectRows: []map[string]any{
			{"id": int64(1), "title": "Hello", "status": "published"},
		},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": "", "status": ""})
	b := scout.NewBuilder(model, "Hello").
		Where("status", "published")

	_, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "status = ?") {
		t.Fatalf("expected SQL to contain where clause, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEnginePaginate(t *testing.T) {
	t.Parallel()
	// DatabaseEngineTest::test_it_can_paginate_results

	conn := &mockConnection{
		driver: "sqlite",
		selectRows: []map[string]any{
			{"id": int64(1), "title": "Post 1"},
			{"id": int64(2), "title": "Post 2"},
		},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	b := scout.NewBuilder(model, "Post")

	result, err := e.Paginate(context.Background(), b, 10, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	if !strings.Contains(conn.lastSQL, "limit 10 offset 0") {
		t.Fatalf("expected SQL to contain limit/offset, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEngineLimitIsApplied(t *testing.T) {
	t.Parallel()
	// DatabaseEngineTest::test_limit_is_applied

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	b := scout.NewBuilder(model, "").Take(5)

	_, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "limit 5") {
		t.Fatalf("expected SQL to contain limit 5, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEngineOrderByIsApplied(t *testing.T) {
	t.Parallel()
	// DatabaseEngineTest::test_it_can_order_results

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	b := scout.NewBuilder(model, "").OrderBy("title", "asc")

	_, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "order by title asc") {
		t.Fatalf("expected SQL to contain order by clause, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEngineBuildFullTextMySQL(t *testing.T) {
	t.Parallel()
	// DatabaseEngineTest::test_it_adds_search_where_clauses_with_non_empty_search

	conn := &mockConnection{
		driver:     "mysql",
		selectRows: []map[string]any{},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	b := scout.NewBuilder(model, "test query")

	_, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "match") || !strings.Contains(conn.lastSQL, "against") {
		t.Fatalf("expected MySQL MATCH AGAINST syntax, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEngineBuildFullTextPostgres(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver:     "pgsql",
		selectRows: []map[string]any{},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	b := scout.NewBuilder(model, "test query")

	_, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "to_tsvector") || !strings.Contains(conn.lastSQL, "plainto_tsquery") {
		t.Fatalf("expected PostgreSQL tsvector syntax, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEngineBuildFullTextSQLite(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0, "title": ""})
	b := scout.NewBuilder(model, "test")

	_, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "like") {
		t.Fatalf("expected SQLite LIKE syntax, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEngineUpdateDeleteFlushAreNoOps(t *testing.T) {
	t.Parallel()

	resolver := &mockResolver{conn: &mockConnection{driver: "sqlite"}}
	e := engines.NewDatabaseEngine(resolver)

	if err := e.Update(context.Background(), nil); err != nil {
		t.Fatalf("Update should be no-op: %v", err)
	}

	if err := e.Delete(context.Background(), nil); err != nil {
		t.Fatalf("Delete should be no-op: %v", err)
	}

	if err := e.Flush(context.Background(), nil); err != nil {
		t.Fatalf("Flush should be no-op: %v", err)
	}

	if err := e.CreateIndex(context.Background(), "test", nil); err != nil {
		t.Fatalf("CreateIndex should be no-op: %v", err)
	}

	if err := e.DeleteIndex(context.Background(), "test"); err != nil {
		t.Fatalf("DeleteIndex should be no-op: %v", err)
	}
}

func TestDatabaseEngineMapIdsInvalidResult(t *testing.T) {
	t.Parallel()

	resolver := &mockResolver{conn: &mockConnection{driver: "sqlite"}}
	e := engines.NewDatabaseEngine(resolver)

	ids := e.MapIds("invalid")

	if len(ids) != 0 {
		t.Fatalf("expected empty ids for invalid result, got %d", len(ids))
	}
}

func TestDatabaseEngineTotalCountInvalidResult(t *testing.T) {
	t.Parallel()

	resolver := &mockResolver{conn: &mockConnection{driver: "sqlite"}}
	e := engines.NewDatabaseEngine(resolver)

	count := e.GetTotalCount("invalid")

	if count != 0 {
		t.Fatalf("expected 0 for invalid result, got %d", count)
	}
}

func TestDatabaseEngineConnectionError(t *testing.T) {
	t.Parallel()

	resolver := &errorResolver{err: fmt.Errorf("connection failed")}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test")

	_, err := e.Search(context.Background(), b)

	if err == nil {
		t.Fatal("expected error from failed connection")
	}

	if !strings.Contains(err.Error(), "connection failed") {
		t.Fatalf("expected connection failed error, got: %v", err)
	}
}

func TestDatabaseEngineSearchErrorWrapsSentinels(t *testing.T) {
	t.Parallel()

	selectErr := errors.New("select failed")
	conn := &mockConnection{
		driver:    "sqlite",
		selectErr: selectErr,
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(1, "posts", map[string]any{"id": 1})
	b := scout.NewBuilder(model, "test")

	_, err := e.Search(context.Background(), b)

	if !errors.Is(err, scout.ErrSearchFailed) {
		t.Fatalf("Search() error = %v, want %v", err, scout.ErrSearchFailed)
	}

	if !errors.Is(err, selectErr) {
		t.Fatalf("Search() error = %v, want %v", err, selectErr)
	}
}

func (r *errorResolver) Connection(_ context.Context, _ ...string) (dbcontract.Connection, error) {
	return nil, r.err
}
func (r *errorResolver) GetDefaultConnection() string  { return "" }
func (r *errorResolver) SetDefaultConnection(_ string) {}

func TestDatabaseEngineSoftDeleteFilter(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver, true) // soft delete enabled

	model := &testModel{
		id:         1,
		table:      "posts",
		softDelete: true,
		data:       map[string]any{"id": 1, "title": "Test"},
	}
	b := scout.NewBuilder(model, "Test")

	_, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify soft delete filter is applied.
	if !strings.Contains(conn.lastSQL, "deleted_at is null") {
		t.Fatalf("expected soft delete filter, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEnginePaginateUsingDatabase(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{{"id": int64(1)}},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0})
	b := scout.NewBuilder(model, "")

	_, err := e.PaginateUsingDatabase(context.Background(), b, 10, 2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "limit 10 offset 10") {
		t.Fatalf("expected limit 10 offset 10, got: %s", conn.lastSQL)
	}
}

func TestDatabaseEngineWhereInAndNotIn(t *testing.T) {
	t.Parallel()

	conn := &mockConnection{
		driver:     "sqlite",
		selectRows: []map[string]any{},
	}
	resolver := &mockResolver{conn: conn}
	e := engines.NewDatabaseEngine(resolver)

	model := newTestModelWithData(0, "posts", map[string]any{"id": 0})
	b := scout.NewBuilder(model, "").
		WhereIn("status", []any{"published", "draft"}).
		WhereNotIn("category", []any{"archived"})

	_, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(conn.lastSQL, "in (?, ?)") {
		t.Fatalf("expected whereIn clause, got: %s", conn.lastSQL)
	}

	if !strings.Contains(conn.lastSQL, "not in (?)") {
		t.Fatalf("expected whereNotIn clause, got: %s", conn.lastSQL)
	}
}
