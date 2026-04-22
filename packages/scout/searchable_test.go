package scout_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/scout"
	"github.com/bedrock/packages/scout"
)

// Syncing is enabled by default.

// Disable syncing.

// Re-enable syncing.

type fakeEngine struct {
	updateCalls      int
	deleteCalls      int
	updateBatchSizes []int
	deleteBatchSizes []int
	lastModels       []contract.Searchable
}

// unsearchableModel always returns false from ShouldBeSearchable.
type unsearchableModel struct {
	testModel
}

type previousSearchableModel struct {
	testModel
	shouldBeSearchable        bool
	wasSearchableBeforeUpdate bool
	wasSearchableBeforeDelete bool
}

func (m *previousSearchableModel) ShouldBeSearchable() bool {
	return m.shouldBeSearchable
}

func (m *previousSearchableModel) WasSearchableBeforeUpdate() bool {
	return m.wasSearchableBeforeUpdate
}

func (m *previousSearchableModel) WasSearchableBeforeDelete() bool {
	return m.wasSearchableBeforeDelete
}

func TestSearchCreatesBuilder(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	b := scout.Search(model, "hello world")

	if b.GetQuery() != "hello world" {
		t.Fatalf("expected query %q, got %q", "hello world", b.GetQuery())
	}

	if b.GetModel() != model {
		t.Fatal("expected model to match")
	}
}

func TestSearchWithCallback(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")
	called := false
	cb := func(_ contract.Engine, _ string, _ map[string]any) any {
		called = true

		return nil
	}

	b := scout.Search(model, "test", cb)

	if !b.HasCallback() {
		t.Fatal("expected callback to be set")
	}

	b.GetCallback()(nil, "", nil)

	if !called {
		t.Fatal("callback was not invoked")
	}
}

func TestSearchSyncingToggle(t *testing.T) {
	model := newTestModel(1, "posts")

	if !scout.IsSearchSyncingEnabled(model) {
		t.Fatal("expected syncing to be enabled by default")
	}

	scout.DisableSearchSyncing(model)

	if scout.IsSearchSyncingEnabled(model) {
		t.Fatal("expected syncing to be disabled")
	}

	scout.EnableSearchSyncing(model)

	if !scout.IsSearchSyncingEnabled(model) {
		t.Fatal("expected syncing to be re-enabled")
	}
}

func TestWithoutSyncingToSearch(t *testing.T) {
	model := newTestModel(42, "articles")

	if !scout.IsSearchSyncingEnabled(model) {
		t.Fatal("expected syncing to be enabled before")
	}

	var syncingDuringCallback bool

	scout.WithoutSyncingToSearch(model, func() {
		syncingDuringCallback = scout.IsSearchSyncingEnabled(model)
	})

	if syncingDuringCallback {
		t.Fatal("expected syncing to be disabled during callback")
	}

	if !scout.IsSearchSyncingEnabled(model) {
		t.Fatal("expected syncing to be re-enabled after callback")
	}
}

func (e *fakeEngine) Update(_ context.Context, models []contract.Searchable) error {
	e.updateCalls++
	e.updateBatchSizes = append(e.updateBatchSizes, len(models))
	e.lastModels = models

	return nil
}
func (e *fakeEngine) Delete(_ context.Context, models []contract.Searchable) error {
	e.deleteCalls++
	e.deleteBatchSizes = append(e.deleteBatchSizes, len(models))
	e.lastModels = models

	return nil
}
func (e *fakeEngine) Search(context.Context, contract.SearchBuilder) (any, error) {
	return nil, nil
}
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

func TestMakeSearchable(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_searchable_using_update_is_called_on_collection
	// SearchableTest::test_searchable_using_update_is_called_on_collection_sync
	engine := &fakeEngine{}
	models := []contract.Searchable{
		newTestModel(1, "posts"),
		newTestModel(2, "posts"),
	}

	err := scout.MakeSearchable(context.Background(), models, engine)

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

func TestMakeSearchableEmpty(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_searchable_using_update_is_not_called_on_empty_collection
	engine := &fakeEngine{}
	err := scout.MakeSearchable(context.Background(), nil, engine)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.updateCalls != 0 {
		t.Fatal("should not call update for empty models")
	}
}

func (m *unsearchableModel) ShouldBeSearchable() bool { return false }

func TestMakeSearchableFiltersNonSearchable(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_overridden_make_searchable_is_dispatched
	engine := &fakeEngine{}
	models := []contract.Searchable{
		newTestModel(1, "posts"),
		&unsearchableModel{testModel{id: 2, table: "posts"}},
	}

	err := scout.MakeSearchable(context.Background(), models, engine)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(engine.lastModels) != 1 {
		t.Fatalf("expected 1 searchable model, got %d", len(engine.lastModels))
	}
}

func TestRemoveFromSearch(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_searchable_using_delete_is_called_on_collection
	// SearchableTest::test_searchable_using_delete_is_called_on_collection_sycn
	engine := &fakeEngine{}
	models := []contract.Searchable{
		newTestModel(1, "posts"),
	}

	err := scout.RemoveFromSearch(context.Background(), models, engine)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.deleteCalls != 1 {
		t.Fatalf("expected 1 delete call, got %d", engine.deleteCalls)
	}
}

func TestRemoveFromSearchEmpty(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_searchable_using_delete_is_not_called_on_empty_collection
	engine := &fakeEngine{}
	err := scout.RemoveFromSearch(context.Background(), nil, engine)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.deleteCalls != 0 {
		t.Fatal("should not call delete for empty models")
	}
}

func TestMakeAllSearchableChunking(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	models := make([]contract.Searchable, 12)

	for i := range models {
		models[i] = newTestModel(i+1, "posts")
	}

	err := scout.MakeAllSearchable(context.Background(), models, engine, 5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 12 models / chunk size 5 = 3 chunks (5 + 5 + 2)
	if engine.updateCalls != 3 {
		t.Fatalf("expected 3 update calls, got %d", engine.updateCalls)
	}
}

func TestMakeAllSearchableUsesDefaultChunkSize(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_searchable_using_update_is_called_on_collection
	engine := &fakeEngine{}
	models := make([]contract.Searchable, 501)

	for i := range models {
		models[i] = newTestModel(i+1, "posts")
	}

	err := scout.MakeAllSearchable(context.Background(), models, engine, 0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.updateCalls != 2 {
		t.Fatalf("expected 2 update calls for the default chunk size, got %d", engine.updateCalls)
	}

	if len(engine.updateBatchSizes) != 2 || engine.updateBatchSizes[0] != 500 || engine.updateBatchSizes[1] != 1 {
		t.Fatalf("expected default chunk sizes [500 1], got %v", engine.updateBatchSizes)
	}
}

func TestMakeAllSearchableFiltersUnsearchableModels(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_searchable_using_update_is_called_on_collection
	engine := &fakeEngine{}
	models := []contract.Searchable{
		newTestModel(1, "posts"),
		&unsearchableModel{testModel{id: 2, table: "posts"}},
		newTestModel(3, "posts"),
	}

	err := scout.MakeAllSearchable(context.Background(), models, engine, 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call after filtering, got %d", engine.updateCalls)
	}

	if len(engine.updateBatchSizes) != 1 || engine.updateBatchSizes[0] != 2 {
		t.Fatalf("expected 2 searchable models, got %v", engine.updateBatchSizes)
	}
}

func TestWasSearchableOnModelWithoutSoftDeletes(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_was_searchable_on_model_without_soft_deletes
	model := newTestModel(1, "posts")

	if !scout.WasSearchableBeforeUpdate(model) {
		t.Fatal("expected model without soft deletes to be searchable before update")
	}

	if !scout.WasSearchableBeforeDelete(model) {
		t.Fatal("expected model without soft deletes to be searchable before delete")
	}
}

func TestWasSearchableBeforeUpdateWorksFromTrueToFalse(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_was_searchable_before_update_works_from_true_to_false
	model := &previousSearchableModel{
		testModel:                 testModel{id: 1, table: "posts"},
		shouldBeSearchable:        false,
		wasSearchableBeforeUpdate: true,
	}

	if !scout.WasSearchableBeforeUpdate(model) {
		t.Fatal("expected model to report it was searchable before update")
	}

	if model.ShouldBeSearchable() {
		t.Fatal("expected model to no longer be searchable")
	}
}

func TestWasSearchableBeforeDeleteWorksWhenDeleting(t *testing.T) {
	t.Parallel()
	// SearchableTest::test_was_searchable_before_delete_works_when_deleting
	model := &previousSearchableModel{
		testModel:                 testModel{id: 1, table: "posts"},
		shouldBeSearchable:        false,
		wasSearchableBeforeDelete: true,
	}

	if !scout.WasSearchableBeforeDelete(model) {
		t.Fatal("expected model to report it was searchable before delete")
	}

	if model.ShouldBeSearchable() {
		t.Fatal("expected deleting model to no longer be searchable")
	}
}

func TestSearchableMixinDefaults(t *testing.T) {
	t.Parallel()

	var mixin scout.SearchableMixin

	if !mixin.ShouldBeSearchable() {
		t.Fatal("expected ShouldBeSearchable() to be true by default")
	}

	if !mixin.SearchIndexShouldBeUpdated() {
		t.Fatal("expected SearchIndexShouldBeUpdated() to be true by default")
	}

	if len(mixin.GetScoutMetadata()) != 0 {
		t.Fatal("expected empty metadata by default")
	}

	if mixin.UsesSoftDelete() {
		t.Fatal("expected UsesSoftDelete() to be false by default")
	}
}

func TestSearchableMixinMetadata(t *testing.T) {
	t.Parallel()

	var mixin scout.SearchableMixin

	mixin.WithScoutMetadata("__soft_deleted", 0)
	mixin.WithScoutMetadata("custom_key", "value")

	meta := mixin.GetScoutMetadata()

	if meta["__soft_deleted"] != 0 {
		t.Fatalf("expected __soft_deleted=0, got %v", meta["__soft_deleted"])
	}

	if meta["custom_key"] != "value" {
		t.Fatalf("expected custom_key=value, got %v", meta["custom_key"])
	}
}

func TestSearchableMixinPrefix(t *testing.T) {
	t.Parallel()

	var mixin scout.SearchableMixin

	if mixin.GetScoutPrefix() != "" {
		t.Fatal("expected empty prefix by default")
	}

	mixin.SetScoutPrefix("prod_")

	if mixin.GetScoutPrefix() != "prod_" {
		t.Fatalf("expected prefix prod_, got %s", mixin.GetScoutPrefix())
	}
}

func TestSearchableAsWithPrefix(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")

	if got := scout.SearchableAs(model, ""); got != "posts" {
		t.Fatalf("expected index posts, got %s", got)
	}

	if got := scout.SearchableAs(model, "prod_"); got != "prod_posts" {
		t.Fatalf("expected index prod_posts, got %s", got)
	}
}

func TestGetScoutKeyDelegatesToModel(t *testing.T) {
	t.Parallel()
	model := newTestModel(42, "posts")

	if got := scout.GetScoutKey(model); got != 42 {
		t.Fatalf("expected key 42, got %v", got)
	}
}

func TestGetScoutKeyNameDelegatesToModel(t *testing.T) {
	t.Parallel()
	model := newTestModel(1, "posts")

	if got := scout.GetScoutKeyName(model); got != "id" {
		t.Fatalf("expected key name id, got %s", got)
	}
}
