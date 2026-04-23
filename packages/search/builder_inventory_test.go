package scout_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/search"
	"github.com/bedrock/packages/search"
)

type terminalEngine struct {
	models       []contract.Searchable
	ids          []any
	total        int64
	searchCalls  int
	paginateCall int
	lastBuilder  contract.SearchBuilder
}

func (e *terminalEngine) Update(context.Context, []contract.Searchable) error { return nil }
func (e *terminalEngine) Delete(context.Context, []contract.Searchable) error { return nil }
func (e *terminalEngine) Search(_ context.Context, builder contract.SearchBuilder) (any, error) {
	e.searchCalls++
	e.lastBuilder = builder

	return e.models, nil
}
func (e *terminalEngine) Paginate(_ context.Context, builder contract.SearchBuilder, _ int, _ int) (any, error) {
	e.paginateCall++
	e.lastBuilder = builder

	return e.models, nil
}
func (e *terminalEngine) MapIds(any) []any { return e.ids }
func (e *terminalEngine) Map(_ context.Context, _ any, _ contract.Searchable) ([]contract.Searchable, error) {
	return e.models, nil
}
func (e *terminalEngine) LazyMap(_ context.Context, _ any, _ contract.Searchable) func(func(contract.Searchable) bool) {
	return func(yield func(contract.Searchable) bool) {
		for _, model := range e.models {
			if !yield(model) {
				return
			}
		}
	}
}
func (e *terminalEngine) GetTotalCount(any) int64 { return e.total }
func (e *terminalEngine) Flush(context.Context, contract.Searchable) error {
	return nil
}
func (e *terminalEngine) CreateIndex(context.Context, string, map[string]any) error {
	return nil
}
func (e *terminalEngine) DeleteIndex(context.Context, string) error { return nil }

var _ contract.Engine = (*terminalEngine)(nil)

// BuilderTest::test_pagination_correctly_handles_paginated_results
func TestBuilderInventoryPaginateWithoutCustomQueryCallback(t *testing.T) {
	t.Parallel()

	models := []contract.Searchable{newTestModel(1, "posts"), newTestModel(2, "posts")}
	engine := &terminalEngine{models: models, total: 5}
	result, err := search.NewBuilder(newTestModel(0, "posts"), "go").
		SetEngine(engine).
		Paginate(context.Background(), 2, 2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.paginateCall != 1 {
		t.Fatalf("expected one paginate call, got %d", engine.paginateCall)
	}

	if len(result.Models) != 2 || result.Total != 5 || !result.HasMorePages() || result.LastPage() != 3 {
		t.Fatalf("unexpected pagination result: %+v", result)
	}
}

func TestBuilderInventoryPaginateWithCustomQueryCallback(t *testing.T) {
	t.Parallel()

	engine := &terminalEngine{models: []contract.Searchable{newTestModel(1, "posts")}, total: 1}
	result, err := search.NewBuilder(newTestModel(0, "posts"), "go").
		SetEngine(engine).
		WithQueryCallback(func(builder *search.Builder) {
			builder.Where("status", "published").Take(10)
		}).
		Paginate(context.Background(), 15, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}

	if engine.lastBuilder.GetWheres()["status"] != "published" || engine.lastBuilder.GetLimit() != 10 {
		t.Fatalf("query callback was not applied: wheres=%v limit=%d", engine.lastBuilder.GetWheres(), engine.lastBuilder.GetLimit())
	}
}

// BuilderTest::test_simple_pagination_correctly_handles_paginated_results
func TestBuilderInventoryPaginateRawAndSimplePaginate(t *testing.T) {
	t.Parallel()

	engine := &terminalEngine{models: []contract.Searchable{newTestModel(1, "posts")}, total: 3}
	builder := search.NewBuilder(newTestModel(0, "posts"), "go").SetEngine(engine)
	raw, err := builder.PaginateRaw(context.Background(), 2, 1)

	if err != nil {
		t.Fatalf("unexpected raw pagination error: %v", err)
	}

	if raw.Total != 3 || !raw.HasMorePages() {
		t.Fatalf("unexpected raw pagination result: %+v", raw)
	}

	simple, err := builder.SimplePaginate(context.Background(), 2, 1)

	if err != nil {
		t.Fatalf("unexpected simple pagination error: %v", err)
	}

	if simple.Total != 3 || !simple.HasMorePages() {
		t.Fatalf("unexpected simple pagination result: %+v", simple)
	}
}

// BuilderTest::test_simple_pagination_correctly_handles_paginated_results_without_more_pages
func TestBuilderInventorySimplePaginationWithoutMorePages(t *testing.T) {
	t.Parallel()

	engine := &terminalEngine{models: []contract.Searchable{newTestModel(1, "posts")}, total: 2}
	result, err := search.NewBuilder(newTestModel(0, "posts"), "go").
		SetEngine(engine).
		SimplePaginate(context.Background(), 2, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HasMorePages() {
		t.Fatalf("expected no more pages: %+v", result)
	}
}
