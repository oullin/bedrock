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
	id                 any
	table              string
	shouldBeSearchable bool
	shouldUpdateIndex  bool
	usesSoftDelete     bool
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

func newObserverModel(id any, searchable, updateIndex bool) *observerTestModel {
	return &observerTestModel{
		id:                 id,
		table:              "posts",
		shouldBeSearchable: searchable,
		shouldUpdateIndex:  updateIndex,
	}
}

func TestModelObserverSaved(t *testing.T) {
	t.Parallel()
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

func TestModelObserverSavedSkipsWhenIndexShouldNotUpdate(t *testing.T) {
	t.Parallel()
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

func TestModelObserverSavedSyncingDisabled(t *testing.T) {
	t.Parallel()
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

func TestModelObserverDeletedWithSoftDelete(t *testing.T) {
	t.Parallel()
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
