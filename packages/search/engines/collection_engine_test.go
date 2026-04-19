package engines_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/search"
	"github.com/bedrock/packages/search"
	"github.com/bedrock/packages/search/engines"
)

func makeBuilder(model contract.Searchable, query string, models []contract.Searchable) *search.Builder {
	b := search.NewBuilder(model, query).
		WithOptions(map[string]any{"__models": models})

	return b
}

func TestCollectionEngineSearch(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := []contract.Searchable{
		newTestModelWithData(1, "posts", map[string]any{"id": 1, "title": "Hello World", "status": "published"}),
		newTestModelWithData(2, "posts", map[string]any{"id": 2, "title": "Go Programming", "status": "published"}),
		newTestModelWithData(3, "posts", map[string]any{"id": 3, "title": "Hello Go", "status": "draft"}),
	}

	b := makeBuilder(newTestModel(0, "posts"), "Hello", models)

	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 2 {
		t.Fatalf("expected 2 results, got %d", total)
	}

	ids := e.MapIds(result)

	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(ids))
	}
}

func TestCollectionEngineSearchCaseInsensitive(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := []contract.Searchable{
		newTestModelWithData(1, "posts", map[string]any{"id": 1, "title": "HELLO WORLD"}),
		newTestModelWithData(2, "posts", map[string]any{"id": 2, "title": "hello world"}),
	}

	b := makeBuilder(newTestModel(0, "posts"), "hello", models)
	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 2 {
		t.Fatalf("expected 2 results for case-insensitive search, got %d", total)
	}
}

func TestCollectionEngineWhereFilter(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := []contract.Searchable{
		newTestModelWithData(1, "posts", map[string]any{"id": 1, "title": "Hello", "status": "published"}),
		newTestModelWithData(2, "posts", map[string]any{"id": 2, "title": "World", "status": "draft"}),
	}

	b := makeBuilder(newTestModel(0, "posts"), "", models).
		Where("status", "published")

	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 1 {
		t.Fatalf("expected 1 result, got %d", total)
	}
}

func TestCollectionEngineWhereInFilter(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := []contract.Searchable{
		newTestModelWithData(1, "posts", map[string]any{"id": 1, "status": "published"}),
		newTestModelWithData(2, "posts", map[string]any{"id": 2, "status": "draft"}),
		newTestModelWithData(3, "posts", map[string]any{"id": 3, "status": "archived"}),
	}

	b := makeBuilder(newTestModel(0, "posts"), "", models).
		WhereIn("status", []any{"published", "draft"})

	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 2 {
		t.Fatalf("expected 2 results, got %d", total)
	}
}

func TestCollectionEngineWhereNotInFilter(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := []contract.Searchable{
		newTestModelWithData(1, "posts", map[string]any{"id": 1, "status": "published"}),
		newTestModelWithData(2, "posts", map[string]any{"id": 2, "status": "draft"}),
		newTestModelWithData(3, "posts", map[string]any{"id": 3, "status": "archived"}),
	}

	b := makeBuilder(newTestModel(0, "posts"), "", models).
		WhereNotIn("status", []any{"archived"})

	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 2 {
		t.Fatalf("expected 2 results, got %d", total)
	}
}

func TestCollectionEngineOrdering(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := []contract.Searchable{
		newTestModelWithData(1, "posts", map[string]any{"id": 1, "title": "Bravo"}),
		newTestModelWithData(2, "posts", map[string]any{"id": 2, "title": "Alpha"}),
		newTestModelWithData(3, "posts", map[string]any{"id": 3, "title": "Charlie"}),
	}

	b := makeBuilder(newTestModel(0, "posts"), "", models).
		OrderBy("title", "asc")

	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mapped, _ := e.Map(context.Background(), result, nil)

	if len(mapped) != 3 {
		t.Fatalf("expected 3 results, got %d", len(mapped))
	}

	if mapped[0].GetScoutKey() != 2 {
		t.Fatalf("expected Alpha (id=2) first, got id=%v", mapped[0].GetScoutKey())
	}

	if mapped[1].GetScoutKey() != 1 {
		t.Fatalf("expected Bravo (id=1) second, got id=%v", mapped[1].GetScoutKey())
	}
}

func TestCollectionEngineLimit(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := make([]contract.Searchable, 10)

	for i := range models {
		models[i] = newTestModelWithData(i+1, "posts", map[string]any{"id": i + 1, "title": "Post"})
	}

	b := makeBuilder(newTestModel(0, "posts"), "", models).Take(3)

	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mapped, _ := e.Map(context.Background(), result, nil)

	if len(mapped) != 3 {
		t.Fatalf("expected 3 results, got %d", len(mapped))
	}
}

func TestCollectionEnginePaginate(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := make([]contract.Searchable, 10)

	for i := range models {
		models[i] = newTestModelWithData(i+1, "posts", map[string]any{"id": i + 1, "title": "Post"})
	}

	b := makeBuilder(newTestModel(0, "posts"), "", models)

	// Page 1
	result, err := e.Paginate(context.Background(), b, 3, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 10 {
		t.Fatalf("expected total 10, got %d", total)
	}

	mapped, _ := e.Map(context.Background(), result, nil)

	if len(mapped) != 3 {
		t.Fatalf("expected 3 on page 1, got %d", len(mapped))
	}

	// Page 4 (last partial page)
	result2, err := e.Paginate(context.Background(), b, 3, 4)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mapped2, _ := e.Map(context.Background(), result2, nil)

	if len(mapped2) != 1 {
		t.Fatalf("expected 1 on page 4, got %d", len(mapped2))
	}
}

func TestCollectionEngineLazyMap(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	models := []contract.Searchable{
		newTestModelWithData(1, "posts", map[string]any{"id": 1, "title": "A"}),
		newTestModelWithData(2, "posts", map[string]any{"id": 2, "title": "B"}),
	}

	b := makeBuilder(newTestModel(0, "posts"), "", models)
	result, _ := e.Search(context.Background(), b)

	iter := e.LazyMap(context.Background(), result, nil)
	count := 0
	iter(func(m contract.Searchable) bool {
		count++

		return true
	})

	if count != 2 {
		t.Fatalf("expected 2 iterations, got %d", count)
	}
}

func TestCollectionEngineEmptySearch(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

	b := makeBuilder(newTestModel(0, "posts"), "test", nil)
	result, err := e.Search(context.Background(), b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := e.GetTotalCount(result)

	if total != 0 {
		t.Fatalf("expected 0 results, got %d", total)
	}
}

func TestCollectionEngineUpdateAndDeleteAreNoOps(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()

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

func TestCollectionEngineMapIdsInvalidResult(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()
	ids := e.MapIds("invalid")

	if len(ids) != 0 {
		t.Fatalf("expected empty ids for invalid result, got %d", len(ids))
	}
}

func TestCollectionEngineTotalCountInvalidResult(t *testing.T) {
	t.Parallel()
	e := engines.NewCollectionEngine()
	count := e.GetTotalCount("invalid")

	if count != 0 {
		t.Fatalf("expected 0 for invalid result, got %d", count)
	}
}
