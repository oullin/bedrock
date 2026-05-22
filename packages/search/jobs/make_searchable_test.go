package jobs_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/search"
	"github.com/bedrock/packages/search/jobs"
)

// testModel implements Searchable for testing.
type testModel struct {
	id    any
	table string
}

type unsearchableTestModel struct {
	testModel
}

// fakeEngine records calls for testing.
type fakeEngine struct {
	updateCalls int
	deleteCalls int
	lastModels  []contract.Searchable
}

func (m *testModel) GetSearchKey() any                                  { return m.id }
func (m *testModel) GetSearchKeyName() string                           { return "id" }
func (m *testModel) SearchableAs() string                               { return m.table }
func (m *testModel) ToSearchableArray() map[string]any                  { return map[string]any{"id": m.id} }
func (m *testModel) ShouldBeSearchable() bool                           { return true }
func (m *testModel) SearchIndexShouldBeUpdated() bool                   { return true }
func (m *testModel) GetSearchMetadata() map[string]any                  { return nil }
func (m *testModel) WithSearchMetadata(string, any) contract.Searchable { return m }
func (m *testModel) GetTable() string                                   { return m.table }
func (m *testModel) GetKeyName() string                                 { return "id" }
func (m *testModel) GetKey() any                                        { return m.id }
func (m *testModel) GetConnectionName() string                          { return "" }
func (m *testModel) UsesSoftDelete() bool                               { return false }

func (m *unsearchableTestModel) ShouldBeSearchable() bool { return false }

func (e *fakeEngine) Update(_ context.Context, models []contract.Searchable) error {
	e.updateCalls++
	e.lastModels = models

	return nil
}
func (e *fakeEngine) Delete(_ context.Context, models []contract.Searchable) error {
	e.deleteCalls++
	e.lastModels = models

	return nil
}
func (e *fakeEngine) Search(context.Context, contract.SearchBuilder) (any, error) { return nil, nil }
func (e *fakeEngine) Paginate(context.Context, contract.SearchBuilder, int, int) (any, error) {
	return nil, nil
}
func (e *fakeEngine) MapIds(any) []any { return nil }
func (e *fakeEngine) Map(context.Context, any, contract.Searchable) ([]contract.Searchable, error) {
	return nil, nil
}
func (e *fakeEngine) LazyMap(context.Context, any, contract.Searchable) func(yield func(contract.Searchable) bool) {
	return func(yield func(contract.Searchable) bool) {}
}
func (e *fakeEngine) GetTotalCount(any) int64                                   { return 0 }
func (e *fakeEngine) Flush(context.Context, contract.Searchable) error          { return nil }
func (e *fakeEngine) CreateIndex(context.Context, string, map[string]any) error { return nil }
func (e *fakeEngine) DeleteIndex(context.Context, string) error                 { return nil }

var _ contract.Engine = (*fakeEngine)(nil)

func TestMakeSearchableHandle(t *testing.T) {
	t.Parallel()
	// MakeSearchableTest::test_handle_passes_the_collection_to_engine
	engine := &fakeEngine{}
	models := []contract.Searchable{
		&testModel{id: 1, table: "posts"},
		&testModel{id: 2, table: "posts"},
	}

	job := jobs.NewMakeSearchable(models, engine)
	err := job.Handle(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call, got %d", engine.updateCalls)
	}

	if len(engine.lastModels) != 2 {
		t.Fatalf("expected 2 models, got %d", len(engine.lastModels))
	}
}

func TestMakeSearchableHandleEmpty(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	job := jobs.NewMakeSearchable(nil, engine)
	err := job.Handle(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.updateCalls != 0 {
		t.Fatal("should not call update for empty models")
	}
}

func TestMakeSearchableFiltersUnsearchable(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	models := []contract.Searchable{
		&testModel{id: 1, table: "posts"},
		&unsearchableTestModel{testModel{id: 2, table: "posts"}},
	}

	job := jobs.NewMakeSearchable(models, engine)
	err := job.Handle(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(engine.lastModels) != 1 {
		t.Fatalf("expected 1 searchable model, got %d", len(engine.lastModels))
	}
}

func TestMakeSearchableGetModels(t *testing.T) {
	t.Parallel()
	models := []contract.Searchable{
		&testModel{id: 1, table: "posts"},
	}
	job := jobs.NewMakeSearchable(models, &fakeEngine{})

	if len(job.GetModels()) != 1 {
		t.Fatalf("expected 1 model, got %d", len(job.GetModels()))
	}
}
