package scout_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/scout"
	"github.com/bedrock/packages/scout"
)

// testModel is a minimal Searchable implementation for testing.
type testModel struct {
	id    any
	table string
}

// Default column

// Custom column

// Without Within, uses model default.

// With custom index.

type paginateResult struct {
	models []contract.Searchable
	total  int64
}

type fakePaginationEngine struct {
	paginateCalls     int
	mapCalls          int
	getTotalCalls     int
	callbackInvoked   bool
	lastCallbackQuery string
}

func (m *testModel) GetScoutKey() any                                  { return m.id }
func (m *testModel) GetScoutKeyName() string                           { return "id" }
func (m *testModel) SearchableAs() string                              { return "test_" + m.table }
func (m *testModel) ToSearchableArray() map[string]any                 { return map[string]any{"id": m.id} }
func (m *testModel) ShouldBeSearchable() bool                          { return true }
func (m *testModel) SearchIndexShouldBeUpdated() bool                  { return true }
func (m *testModel) GetScoutMetadata() map[string]any                  { return nil }
func (m *testModel) WithScoutMetadata(string, any) contract.Searchable { return m }
func (m *testModel) GetTable() string                                  { return m.table }
func (m *testModel) GetKeyName() string                                { return "id" }
func (m *testModel) GetKey() any                                       { return m.id }
func (m *testModel) GetConnectionName() string                         { return "" }
func (m *testModel) UsesSoftDelete() bool                              { return false }

func newTestModel(id any, table string) *testModel {
	return &testModel{id: id, table: table}
}

func TestNewBuilder(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "hello world")

	if b.GetQuery() != "hello world" {
		t.Fatalf("expected query %q, got %q", "hello world", b.GetQuery())
	}

	if b.GetModel() != model {
		t.Fatal("expected model to match")
	}

	if b.GetLimit() != 0 {
		t.Fatalf("expected no limit, got %d", b.GetLimit())
	}

	if b.HasCallback() {
		t.Fatal("expected no callback")
	}
}

func TestBuilderWhere(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test").
		Where("status", "published").
		Where("category", "tech")

	wheres := b.GetWheres()

	if wheres["status"] != "published" {
		t.Fatalf("expected status=published, got %v", wheres["status"])
	}

	if wheres["category"] != "tech" {
		t.Fatalf("expected category=tech, got %v", wheres["category"])
	}
}

func TestBuilderWhereIn(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test").
		WhereIn("id", []any{1, 2, 3}).
		WhereNotIn("status", []any{"draft", "archived"})

	whereIns := b.GetWhereIns()

	if len(whereIns["id"]) != 3 {
		t.Fatalf("expected 3 whereIn values, got %d", len(whereIns["id"]))
	}

	whereNotIns := b.GetWhereNotIns()

	if len(whereNotIns["status"]) != 2 {
		t.Fatalf("expected 2 whereNotIn values, got %d", len(whereNotIns["status"]))
	}
}

func TestBuilderOrderBy(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test").
		OrderBy("title", "asc").
		OrderBy("created_at", "desc")

	orders := b.GetOrders()

	if len(orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(orders))
	}

	if orders[0].Column != "title" || orders[0].Direction != "asc" {
		t.Fatalf("unexpected first order: %+v", orders[0])
	}

	if orders[1].Column != "created_at" || orders[1].Direction != "desc" {
		t.Fatalf("unexpected second order: %+v", orders[1])
	}
}

func TestBuilderLatest(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")

	b := scout.NewBuilder(model, "test").Latest()
	orders := b.GetOrders()

	if len(orders) != 1 || orders[0].Column != "created_at" || orders[0].Direction != "desc" {
		t.Fatalf("Latest() default: unexpected %+v", orders)
	}

	b2 := scout.NewBuilder(model, "test").Latest("updated_at")
	orders2 := b2.GetOrders()

	if orders2[0].Column != "updated_at" {
		t.Fatalf("Latest(updated_at): expected updated_at, got %s", orders2[0].Column)
	}
}

func TestBuilderOldest(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")

	b := scout.NewBuilder(model, "test").Oldest()
	orders := b.GetOrders()

	if len(orders) != 1 || orders[0].Column != "created_at" || orders[0].Direction != "asc" {
		t.Fatalf("Oldest() default: unexpected %+v", orders)
	}
}

func TestBuilderTake(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test").Take(25)

	if b.GetLimit() != 25 {
		t.Fatalf("expected limit 25, got %d", b.GetLimit())
	}
}

func TestBuilderWithin(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")

	b := scout.NewBuilder(model, "test")

	if b.GetIndex() != "test_posts" {
		t.Fatalf("expected index test_posts, got %s", b.GetIndex())
	}

	b2 := scout.NewBuilder(model, "test").Within("custom_index")

	if b2.GetIndex() != "custom_index" {
		t.Fatalf("expected index custom_index, got %s", b2.GetIndex())
	}
}

func TestBuilderWithOptions(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test").
		WithOptions(map[string]any{"hitsPerPage": 20}).
		WithOptions(map[string]any{"typoTolerance": true})

	opts := b.GetOptions()

	if opts["hitsPerPage"] != 20 {
		t.Fatalf("expected hitsPerPage=20, got %v", opts["hitsPerPage"])
	}

	if opts["typoTolerance"] != true {
		t.Fatalf("expected typoTolerance=true, got %v", opts["typoTolerance"])
	}
}

func TestBuilderCallback(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")

	called := false
	cb := func(_ contract.Engine, _ string, _ map[string]any) any {
		called = true

		return nil
	}

	b := scout.NewBuilder(model, "test", cb)

	if !b.HasCallback() {
		t.Fatal("expected callback to be set")
	}

	b.GetCallback()(nil, "", nil)

	if !called {
		t.Fatal("expected callback to be invoked")
	}
}

func (e *fakePaginationEngine) Update(context.Context, []contract.Searchable) error { return nil }
func (e *fakePaginationEngine) Delete(context.Context, []contract.Searchable) error { return nil }
func (e *fakePaginationEngine) Search(context.Context, contract.SearchBuilder) (any, error) {
	return nil, nil
}
func (e *fakePaginationEngine) Paginate(_ context.Context, builder contract.SearchBuilder, perPage, page int) (any, error) {
	e.paginateCalls++

	if builder.HasCallback() {
		cb := builder.GetCallback()

		if cb != nil {
			e.callbackInvoked = true
			e.lastCallbackQuery = builder.GetQuery()
			_ = cb(e, builder.GetQuery(), builder.GetOptions())
		}
	}

	start := (page - 1) * perPage
	end := start + perPage

	if start < 0 {
		start = 0
	}

	if end < start {
		end = start
	}

	models := []contract.Searchable{
		newTestModel(1, "posts"),
		newTestModel(2, "posts"),
		newTestModel(3, "posts"),
	}

	if start > len(models) {
		start = len(models)
	}

	if end > len(models) {
		end = len(models)
	}

	return &paginateResult{models: models[start:end], total: int64(len(models))}, nil
}
func (e *fakePaginationEngine) MapIds(any) []any { return nil }
func (e *fakePaginationEngine) Map(_ context.Context, results any, _ contract.Searchable) ([]contract.Searchable, error) {
	e.mapCalls++
	pr, ok := results.(*paginateResult)

	if !ok {
		return nil, nil
	}

	return pr.models, nil
}
func (e *fakePaginationEngine) LazyMap(context.Context, any, contract.Searchable) func(yield func(contract.Searchable) bool) {
	return func(func(contract.Searchable) bool) {}
}
func (e *fakePaginationEngine) GetTotalCount(results any) int64 {
	e.getTotalCalls++
	pr, ok := results.(*paginateResult)

	if !ok {
		return 0
	}

	return pr.total
}
func (e *fakePaginationEngine) Flush(context.Context, contract.Searchable) error          { return nil }
func (e *fakePaginationEngine) CreateIndex(context.Context, string, map[string]any) error { return nil }
func (e *fakePaginationEngine) DeleteIndex(context.Context, string) error                 { return nil }

var _ contract.Engine = (*fakePaginationEngine)(nil)

func TestBuilderPaginateWithoutCustomQueryCallback(t *testing.T) {
	t.Parallel()
	// BuilderTest::test_it_can_paginate_without_custom_query_callback
	engine := &fakePaginationEngine{}
	model := newTestModel(1, "posts")
	builder := scout.NewBuilder(model, "hello").SetEngine(engine)

	result, err := builder.Paginate(context.Background(), 2, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.paginateCalls != 1 {
		t.Fatalf("expected 1 paginate call, got %d", engine.paginateCalls)
	}

	if engine.mapCalls != 1 {
		t.Fatalf("expected 1 map call, got %d", engine.mapCalls)
	}

	if engine.callbackInvoked {
		t.Fatal("expected callback not to be invoked")
	}

	if result.Total != 3 {
		t.Fatalf("expected total 3, got %d", result.Total)
	}

	if len(result.Models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(result.Models))
	}
}

func TestBuilderPaginateWithCustomQueryCallback(t *testing.T) {
	t.Parallel()
	// BuilderTest::test_it_can_paginate_with_custom_query_callback
	engine := &fakePaginationEngine{}
	model := newTestModel(1, "posts")

	builder := scout.NewBuilder(model, "hello", func(_ contract.Engine, query string, _ map[string]any) any {
		if query != "hello" {
			t.Fatalf("expected callback query %q, got %q", "hello", query)
		}

		return nil
	}).SetEngine(engine)

	_, err := builder.Paginate(context.Background(), 2, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !engine.callbackInvoked {
		t.Fatal("expected callback to be invoked")
	}

	if engine.lastCallbackQuery != "hello" {
		t.Fatalf("expected last callback query %q, got %q", "hello", engine.lastCallbackQuery)
	}
}

func TestBuilderPaginateRawWithoutCustomQueryCallback(t *testing.T) {
	t.Parallel()
	// BuilderTest::test_it_can_paginate_raw_without_custom_query_callback
	engine := &fakePaginationEngine{}
	model := newTestModel(1, "posts")
	builder := scout.NewBuilder(model, "hello").SetEngine(engine)

	result, err := builder.PaginateRaw(context.Background(), 2, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.paginateCalls != 1 {
		t.Fatalf("expected 1 paginate call, got %d", engine.paginateCalls)
	}

	if engine.mapCalls != 0 {
		t.Fatalf("expected no map calls for PaginateRaw, got %d", engine.mapCalls)
	}

	if result.Total != 3 {
		t.Fatalf("expected total 3, got %d", result.Total)
	}

	if result.Results == nil {
		t.Fatal("expected raw results to be non-nil")
	}
}

func TestBuilderSimplePaginateRawWithoutCustomQueryCallback(t *testing.T) {
	t.Parallel()

	engine := &fakePaginationEngine{}
	model := newTestModel(1, "posts")
	builder := scout.NewBuilder(model, "hello").SetEngine(engine)

	result, err := builder.SimplePaginateRaw(context.Background(), 2, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 3 {
		t.Fatalf("expected total 3, got %d", result.Total)
	}

	if result.Results == nil {
		t.Fatal("expected raw results to be non-nil")
	}

	if !result.HasMorePages() {
		t.Fatal("expected simple raw pagination to report more pages")
	}
}

func TestBuilderFluentChaining(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "go programming").
		Where("status", "published").
		WhereIn("tag", []any{"go", "programming"}).
		WhereNotIn("category", []any{"archived"}).
		OrderBy("relevance", "desc").
		Latest().
		Take(10).
		Within("posts_v2").
		WithOptions(map[string]any{"highlight": true})

	if b.GetQuery() != "go programming" {
		t.Fatal("query mismatch")
	}

	if len(b.GetWheres()) != 1 {
		t.Fatal("wheres mismatch")
	}

	if len(b.GetWhereIns()) != 1 {
		t.Fatal("whereIns mismatch")
	}

	if len(b.GetWhereNotIns()) != 1 {
		t.Fatal("whereNotIns mismatch")
	}

	if len(b.GetOrders()) != 2 {
		t.Fatal("orders mismatch")
	}

	if b.GetLimit() != 10 {
		t.Fatal("limit mismatch")
	}

	if b.GetIndex() != "posts_v2" {
		t.Fatal("index mismatch")
	}

	if b.GetOptions()["highlight"] != true {
		t.Fatal("options mismatch")
	}
}

func TestBuilderGetWithoutEngine(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test")

	_, err := b.Get(context.Background())

	if err != scout.ErrEngineNotConfigured {
		t.Fatalf("expected ErrEngineNotConfigured, got %v", err)
	}
}

func TestBuilderRawWithoutEngine(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test")

	_, err := b.Raw(context.Background())

	if err != scout.ErrEngineNotConfigured {
		t.Fatalf("expected ErrEngineNotConfigured, got %v", err)
	}
}

func TestBuilderKeysWithoutEngine(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test")

	_, err := b.Keys(context.Background())

	if err != scout.ErrEngineNotConfigured {
		t.Fatalf("expected ErrEngineNotConfigured, got %v", err)
	}
}

func TestBuilderCursorWithoutEngine(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test")

	_, err := b.Cursor(context.Background())

	if err != scout.ErrEngineNotConfigured {
		t.Fatalf("expected ErrEngineNotConfigured, got %v", err)
	}
}

func TestBuilderPaginateWithoutEngine(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.NewBuilder(model, "test")

	_, err := b.Paginate(context.Background(), 15, 1)

	if err != scout.ErrEngineNotConfigured {
		t.Fatalf("expected ErrEngineNotConfigured, got %v", err)
	}
}

func TestPaginatedResultHasMorePages(t *testing.T) {
	t.Parallel()
	p := &scout.PaginatedResult{Total: 50, PerPage: 15, CurrentPage: 1}

	if !p.HasMorePages() {
		t.Fatal("expected more pages")
	}

	p2 := &scout.PaginatedResult{Total: 50, PerPage: 15, CurrentPage: 4}

	if p2.HasMorePages() {
		t.Fatal("expected no more pages")
	}
}

func TestPaginatedResultLastPage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		total   int64
		perPage int
		want    int
	}{
		{50, 15, 4},
		{30, 15, 2},
		{0, 15, 1},
		{1, 15, 1},
		{15, 15, 1},
		{16, 15, 2},
	}

	for _, tt := range tests {
		p := &scout.PaginatedResult{Total: tt.total, PerPage: tt.perPage}

		if got := p.LastPage(); got != tt.want {
			t.Errorf("LastPage(total=%d, perPage=%d) = %d, want %d", tt.total, tt.perPage, got, tt.want)
		}
	}
}
