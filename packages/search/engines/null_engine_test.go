package engines_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/search"
	"github.com/bedrock/packages/search/engines"
)

// testModel is a minimal Searchable implementation for testing.
type testModel struct {
	id         any
	table      string
	connection string
	softDelete bool
	data       map[string]any
}

func (m *testModel) GetScoutKey() any        { return m.id }
func (m *testModel) GetScoutKeyName() string { return "id" }
func (m *testModel) SearchableAs() string    { return "test_" + m.table }
func (m *testModel) ToSearchableArray() map[string]any {
	if m.data != nil {
		return m.data
	}

	return map[string]any{"id": m.id}
}
func (m *testModel) ShouldBeSearchable() bool                          { return true }
func (m *testModel) SearchIndexShouldBeUpdated() bool                  { return true }
func (m *testModel) GetScoutMetadata() map[string]any                  { return nil }
func (m *testModel) WithScoutMetadata(string, any) contract.Searchable { return m }
func (m *testModel) GetTable() string                                  { return m.table }
func (m *testModel) GetKeyName() string                                { return "id" }
func (m *testModel) GetKey() any                                       { return m.id }
func (m *testModel) GetConnectionName() string                         { return m.connection }
func (m *testModel) UsesSoftDelete() bool                              { return m.softDelete }

func newTestModel(id any, table string) *testModel {
	return &testModel{id: id, table: table}
}

func newTestModelWithData(id any, table string, data map[string]any) *testModel {
	return &testModel{id: id, table: table, data: data}
}

func TestNullEngineUpdate(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	err := e.Update(context.Background(), []contract.Searchable{newTestModel(1, "posts")})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNullEngineDelete(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	err := e.Delete(context.Background(), []contract.Searchable{newTestModel(1, "posts")})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNullEngineSearch(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	result, err := e.Search(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestNullEnginePaginate(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	result, err := e.Paginate(context.Background(), nil, 15, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestNullEngineMapIds(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	ids := e.MapIds(nil)

	if len(ids) != 0 {
		t.Fatalf("expected empty ids, got %d", len(ids))
	}
}

func TestNullEngineMap(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	models, err := e.Map(context.Background(), nil, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(models) != 0 {
		t.Fatalf("expected empty models, got %d", len(models))
	}
}

func TestNullEngineLazyMap(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	iter := e.LazyMap(context.Background(), nil, nil)

	count := 0
	iter(func(_ contract.Searchable) bool {
		count++

		return true
	})

	if count != 0 {
		t.Fatalf("expected 0 iterations, got %d", count)
	}
}

func TestNullEngineGetTotalCount(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()

	if count := e.GetTotalCount(nil); count != 0 {
		t.Fatalf("expected 0, got %d", count)
	}
}

func TestNullEngineFlush(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	err := e.Flush(context.Background(), newTestModel(1, "posts"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNullEngineCreateIndex(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	err := e.CreateIndex(context.Background(), "test", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNullEngineDeleteIndex(t *testing.T) {
	t.Parallel()
	e := engines.NewNullEngine()
	err := e.DeleteIndex(context.Background(), "test")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
