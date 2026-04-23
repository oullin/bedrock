package scout_test

import (
	"context"
	"testing"

	contract "github.com/bedrock/packages/contracts/search"
	dbevents "github.com/bedrock/packages/database/events"
	"github.com/bedrock/packages/search"
)

// observerTestModel implements Searchable for observer tests.
type observerTestModel struct {
	id                           any
	table                        string
	shouldBeSearchable           bool
	shouldUpdateIndex            bool
	usesSoftDelete               bool
	wasSearchableBeforeUpdate    bool
	hasWasSearchableBeforeUpdate bool
	wasSearchableBeforeDelete    bool
	hasWasSearchableBeforeDelete bool
}

type dirtyObserverModel struct {
	observerTestModel
	dirty map[string]bool
}

func (m *observerTestModel) GetScoutKey() any                                  { return m.id }
func (m *observerTestModel) GetScoutKeyName() string                           { return "id" }
func (m *observerTestModel) SearchableAs() string                              { return m.table }
func (m *observerTestModel) ToSearchableArray() map[string]any                 { return map[string]any{"id": m.id} }
func (m *observerTestModel) ShouldBeSearchable() bool                          { return m.shouldBeSearchable }
func (m *observerTestModel) SearchIndexShouldBeUpdated() bool                  { return m.shouldUpdateIndex }
func (m *observerTestModel) GetScoutMetadata() map[string]any                  { return nil }
func (m *observerTestModel) WithScoutMetadata(string, any) contract.Searchable { return m }
func (m *observerTestModel) GetTable() string                                  { return m.table }
func (m *observerTestModel) GetKeyName() string                                { return "id" }
func (m *observerTestModel) GetKey() any                                       { return m.id }
func (m *observerTestModel) GetConnectionName() string                         { return "" }
func (m *observerTestModel) UsesSoftDelete() bool                              { return m.usesSoftDelete }
func (m *observerTestModel) WasSearchableBeforeUpdate() bool {
	if !m.hasWasSearchableBeforeUpdate {
		return true
	}

	return m.wasSearchableBeforeUpdate
}
func (m *observerTestModel) WasSearchableBeforeDelete() bool {
	if !m.hasWasSearchableBeforeDelete {
		return true
	}

	return m.wasSearchableBeforeDelete
}

func newObserverModel(id any, searchable, updateIndex bool) *observerTestModel {
	return &observerTestModel{
		id:                 id,
		table:              "posts",
		shouldBeSearchable: searchable,
		shouldUpdateIndex:  updateIndex,
	}
}

func (m *dirtyObserverModel) SearchIndexShouldBeUpdated() bool {
	for _, attribute := range []string{"name", "email"} {
		if m.dirty[attribute] {
			return true
		}
	}

	return false
}

func TestModelObserverSaved(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_saved_handler_makes_model_searchable
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, true)
	observer.Saved(context.Background(), dbevents.Saved{Model: model})

	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call, got %d", engine.updateCalls)
	}
}

func TestModelObserverSavedNotSearchable(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_saved_handler_makes_model_unsearchable_when_disabled_per_model_rule
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, false, true)
	observer.Saved(context.Background(), dbevents.Saved{Model: model})

	// Not searchable → should call delete (remove from search).
	if engine.deleteCalls != 1 {
		t.Fatalf("expected 1 delete call, got %d", engine.deleteCalls)
	}

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 update calls, got %d", engine.updateCalls)
	}
}

func TestModelObserverSavedAlreadyUnsearchableSkipsDelete(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, false, true)
	model.hasWasSearchableBeforeUpdate = true
	model.wasSearchableBeforeUpdate = false
	observer.Saved(context.Background(), dbevents.Saved{Model: model})

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 update calls, got %d", engine.updateCalls)
	}

	if engine.deleteCalls != 0 {
		t.Fatalf("expected 0 delete calls for an already-unsearchable model, got %d", engine.deleteCalls)
	}
}

func TestModelObserverSavedSkipsWhenIndexShouldNotUpdate(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_saved_handler_doesnt_make_model_searchable_when_search_shouldnt_update
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, false) // shouldUpdateIndex = false
	observer.Saved(context.Background(), dbevents.Saved{Model: model})

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 update calls when index should not update, got %d", engine.updateCalls)
	}
}

func TestModelObserverUpdateOnSensitiveAttributesTriggersSearch(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_update_on_sensitive_attributes_triggers_search
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := &dirtyObserverModel{
		observerTestModel: observerTestModel{
			id:                 1,
			table:              "posts",
			shouldBeSearchable: true,
		},
		dirty: map[string]bool{"password": true, "name": true},
	}

	observer.Saved(context.Background(), dbevents.Saved{Model: model})

	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call for sensitive dirty attributes, got %d", engine.updateCalls)
	}
}

func TestModelObserverUpdateOnNonSensitiveAttributesDoesNotTriggerSearch(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_update_on_non_sensitive_attributes_doesnt_trigger_search
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := &dirtyObserverModel{
		observerTestModel: observerTestModel{
			id:                 1,
			table:              "posts",
			shouldBeSearchable: true,
		},
		dirty: map[string]bool{"password": true, "remember_token": true},
	}

	observer.Saved(context.Background(), dbevents.Saved{Model: model})

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 update calls for non-sensitive dirty attributes, got %d", engine.updateCalls)
	}

	if engine.deleteCalls != 0 {
		t.Fatalf("expected 0 delete calls for non-sensitive dirty attributes, got %d", engine.deleteCalls)
	}
}

func TestModelObserverSavedWhileForcingUpdate(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_saved_handler_makes_model_searchable
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, false)
	observer.WhileForcingUpdate(func() {
		observer.Saved(context.Background(), dbevents.Saved{Model: model})
	})

	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call while forcing updates, got %d", engine.updateCalls)
	}

	if observer.IsForceUpdating() {
		t.Fatal("expected force updating to be reset after the callback")
	}
}

func TestModelObserverSavedSyncingDisabled(t *testing.T) {
	// ModelObserverTest::test_saved_handler_doesnt_make_model_searchable_when_disabled
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, true)
	search.DisableSearchSyncing(model)

	defer search.EnableSearchSyncing(model)

	observer.Saved(context.Background(), dbevents.Saved{Model: model})

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 update calls when syncing disabled, got %d", engine.updateCalls)
	}
}

func TestModelObserverDeleted(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_deleted_handler_makes_model_unsearchable
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, true)
	observer.Deleted(context.Background(), dbevents.Deleted{Model: model})

	if engine.deleteCalls != 1 {
		t.Fatalf("expected 1 delete call, got %d", engine.deleteCalls)
	}
}

func TestModelObserverDeletedSkipsAlreadyUnsearchable(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_deleted_handler_doesnt_make_model_unsearchable_when_already_unsearchable
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, true)
	model.hasWasSearchableBeforeDelete = true
	model.wasSearchableBeforeDelete = false
	observer.Deleted(context.Background(), dbevents.Deleted{Model: model})

	if engine.deleteCalls != 0 {
		t.Fatalf("expected 0 delete calls for an already-unsearchable model, got %d", engine.deleteCalls)
	}
}

func TestModelObserverDeletedSoftDeleteModelMakesUnsearchable(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_deleted_handler_on_soft_delete_model_makes_model_unsearchable
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, true)
	model.usesSoftDelete = true
	observer.Deleted(context.Background(), dbevents.Deleted{Model: model})

	if engine.deleteCalls != 1 {
		t.Fatalf("expected 1 delete call for soft-delete model without soft-delete indexing, got %d", engine.deleteCalls)
	}
}

func TestModelObserverUnsearchableIsCalledWhenDeleting(t *testing.T) {
	t.Parallel()
	// ModelObserverTest::test_unsearchable_should_be_called_when_deleting
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, false)
	observer.Deleted(context.Background(), dbevents.Deleted{Model: model})

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 update calls when deleting, got %d", engine.updateCalls)
	}

	if engine.deleteCalls != 1 {
		t.Fatalf("expected 1 delete call when deleting, got %d", engine.deleteCalls)
	}
}

func TestModelObserverDeletedWithSoftDelete(t *testing.T) {
	t.Parallel()
	// ModelObserverWithSoftDeletesTest::test_deleted_handler_makes_model_searchable_when_it_should_be_searchable
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()
	config.SoftDelete = true

	observer := search.NewModelObserver(manager, config)

	model := &observerTestModel{
		id:                 1,
		table:              "posts",
		shouldBeSearchable: true,
		shouldUpdateIndex:  true,
		usesSoftDelete:     true,
	}
	observer.Deleted(context.Background(), dbevents.Deleted{Model: model})

	// With soft deletes enabled, should update (not delete).
	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call for soft delete, got %d", engine.updateCalls)
	}

	if engine.deleteCalls != 0 {
		t.Fatalf("expected 0 delete calls for soft delete, got %d", engine.deleteCalls)
	}
}

func TestModelObserverDeletedWithSoftDeleteNotSearchable(t *testing.T) {
	t.Parallel()
	// ModelObserverWithSoftDeletesTest::test_deleted_handler_makes_model_unsearchable_when_it_should_not_be_searchable
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()
	config.SoftDelete = true

	observer := search.NewModelObserver(manager, config)

	model := &observerTestModel{
		id:                 1,
		table:              "posts",
		shouldBeSearchable: false,
		shouldUpdateIndex:  true,
		usesSoftDelete:     true,
	}
	observer.Deleted(context.Background(), dbevents.Deleted{Model: model})

	if engine.deleteCalls != 1 {
		t.Fatalf("expected 1 delete call for soft delete + unsearchable, got %d", engine.deleteCalls)
	}

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 update calls for soft delete + unsearchable, got %d", engine.updateCalls)
	}
}

func TestModelObserverForceDeleted(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, true)
	observer.ForceDeleted(context.Background(), dbevents.ForceDeleted{Model: model})

	if engine.deleteCalls != 1 {
		t.Fatalf("expected 1 delete call, got %d", engine.deleteCalls)
	}
}

func TestModelObserverRestored(t *testing.T) {
	t.Parallel()
	// ModelObserverWithSoftDeletesTest::test_restored_handler_makes_model_searchable
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, true)
	observer.Restored(context.Background(), dbevents.Restored{Model: model})

	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call on restore, got %d", engine.updateCalls)
	}
}

func TestModelObserverSavedNonSearchableModel(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	// Pass a non-Searchable model.
	observer.Saved(context.Background(), dbevents.Saved{Model: "not a searchable"})

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 calls for non-searchable, got %d", engine.updateCalls)
	}
}

func TestModelObserverWhileForcingUpdate(t *testing.T) {
	t.Parallel()
	engine := &fakeEngine{}
	manager := search.NewEngineManager()
	manager.Register("null", engine)
	manager.SetDefaultDriver("null")
	config := search.DefaultConfig()

	observer := search.NewModelObserver(manager, config)

	model := newObserverModel(1, true, false) // shouldUpdateIndex = false

	// Without forcing, should skip.
	observer.Saved(context.Background(), dbevents.Saved{Model: model})

	if engine.updateCalls != 0 {
		t.Fatalf("expected 0 update calls before force, got %d", engine.updateCalls)
	}

	// With forcing, should update even when shouldUpdateIndex is false.
	observer.WhileForcingUpdate(func() {
		observer.Saved(context.Background(), dbevents.Saved{Model: model})
	})

	if engine.updateCalls != 1 {
		t.Fatalf("expected 1 update call during force, got %d", engine.updateCalls)
	}

	if observer.IsForceUpdating() {
		t.Fatal("expected force updating to be false after callback")
	}
}
