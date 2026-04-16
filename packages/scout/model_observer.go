package scout

import (
	"context"

	contract "github.com/bedrock/packages/contracts/scout"
	dbevents "github.com/bedrock/packages/database/events"
)

// ModelObserver observes Eloquent model lifecycle events and synchronises
// the search index automatically. It mirrors Laravel Scout's ModelObserver.
type ModelObserver struct {
	manager     *EngineManager
	config      Config
	afterCommit bool

	// forceUpdating is true when WhileForcingUpdate is active.
	forceUpdating bool
}

// NewModelObserver creates a new observer.
func NewModelObserver(manager *EngineManager, config Config) *ModelObserver {
	return &ModelObserver{
		manager:     manager,
		config:      config,
		afterCommit: config.AfterCommit,
	}
}

// Saved handles the model saved event. It indexes the model if it should be
// searchable, or removes it if not.
func (o *ModelObserver) Saved(ctx context.Context, event dbevents.Saved) {
	model, ok := event.Model.(contract.Searchable)

	if !ok {
		return
	}

	if !IsSearchSyncingEnabled(model) {
		return
	}

	if !o.forceUpdating && !model.SearchIndexShouldBeUpdated() {
		return
	}

	if model.ShouldBeSearchable() {
		o.makeSearchable(ctx, model)
	} else {
		o.removeFromSearch(ctx, model)
	}
}

// Deleted handles the model deleted event. It removes the model from the
// search index, unless soft deletes are enabled and the model uses them.
func (o *ModelObserver) Deleted(ctx context.Context, event dbevents.Deleted) {
	model, ok := event.Model.(contract.Searchable)

	if !ok {
		return
	}

	if !IsSearchSyncingEnabled(model) {
		return
	}

	if o.config.SoftDelete && model.UsesSoftDelete() {
		// When soft deletes are enabled, update the model in the index
		// (so it includes the soft-delete metadata) instead of removing it.
		o.makeSearchable(ctx, model)

		return
	}

	o.removeFromSearch(ctx, model)
}

// ForceDeleted handles the model force-deleted event. Always removes from index.
func (o *ModelObserver) ForceDeleted(ctx context.Context, event dbevents.ForceDeleted) {
	model, ok := event.Model.(contract.Searchable)

	if !ok {
		return
	}

	if !IsSearchSyncingEnabled(model) {
		return
	}

	o.removeFromSearch(ctx, model)
}

// Restored handles the model restored event. Re-indexes the model.
func (o *ModelObserver) Restored(ctx context.Context, event dbevents.Restored) {
	model, ok := event.Model.(contract.Searchable)

	if !ok {
		return
	}

	if !IsSearchSyncingEnabled(model) {
		return
	}

	o.makeSearchable(ctx, model)
}

// WhileForcingUpdate runs fn with force-updating enabled, ensuring that
// all saves during fn trigger re-indexing regardless of SearchIndexShouldBeUpdated.
func (o *ModelObserver) WhileForcingUpdate(fn func()) {
	o.forceUpdating = true

	defer func() { o.forceUpdating = false }()

	fn()
}

// IsForceUpdating reports whether the observer is currently force-updating.
func (o *ModelObserver) IsForceUpdating() bool {
	return o.forceUpdating
}

func (o *ModelObserver) makeSearchable(ctx context.Context, model contract.Searchable) {
	engine, err := o.manager.Engine()

	if err != nil {
		return
	}

	_ = MakeSearchable(ctx, []contract.Searchable{model}, engine)
}

func (o *ModelObserver) removeFromSearch(ctx context.Context, model contract.Searchable) {
	engine, err := o.manager.Engine()

	if err != nil {
		return
	}

	_ = RemoveFromSearch(ctx, []contract.Searchable{model}, engine)
}
