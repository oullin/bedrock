package search

import (
	"context"
	"reflect"
	"sync"

	contract "github.com/bedrock/packages/contracts/search"
)

// syncingDisabled tracks model types for which search syncing is disabled.
// Keys are reflect.Type values, values are struct{}.

// SearchableMixin provides default implementations of the Searchable
// interface methods. User models embed this struct alongside orm.Model
// to gain Search search capabilities.
//
// In Upstream this is the Searchable trait. In Go, we use composition:
//
//	type Post struct {
//	    orm.Model
//	    search.SearchableMixin
//	}
type SearchableMixin struct {
	scoutPrefix   string
	scoutKeyName  string
	scoutMetadata map[string]any
}

var syncingDisabled sync.Map

// SetScoutPrefix sets the index name prefix.
func (s *SearchableMixin) SetScoutPrefix(prefix string) {
	s.scoutPrefix = prefix
}

// GetScoutPrefix returns the index name prefix.
func (s *SearchableMixin) GetScoutPrefix() string {
	return s.scoutPrefix
}

// SetScoutKeyName overrides the key name used for the search index.
func (s *SearchableMixin) SetScoutKeyName(name string) {
	s.scoutKeyName = name
}

// ShouldBeSearchable reports whether this model should be indexed.
// Override this on your model to add conditions.
func (s *SearchableMixin) ShouldBeSearchable() bool {
	return true
}

// SearchIndexShouldBeUpdated reports whether a save should trigger re-indexing.
func (s *SearchableMixin) SearchIndexShouldBeUpdated() bool {
	return true
}

// GetScoutMetadata returns engine-specific metadata for the model.
func (s *SearchableMixin) GetScoutMetadata() map[string]any {
	if s.scoutMetadata == nil {
		return map[string]any{}
	}

	return s.scoutMetadata
}

// WithScoutMetadata sets a metadata key-value pair.
func (s *SearchableMixin) WithScoutMetadata(key string, value any) {
	if s.scoutMetadata == nil {
		s.scoutMetadata = make(map[string]any)
	}

	s.scoutMetadata[key] = value
}

// UsesSoftDelete reports whether the model uses soft deletes.
// This default returns false; models with SoftDeletes should override.
func (s *SearchableMixin) UsesSoftDelete() bool {
	return false
}

// Search creates a new Builder for the given model and query.
func Search(model contract.Searchable, query string, callback ...func(contract.Engine, string, map[string]any) any) *Builder {
	return NewBuilder(model, query, callback...)
}

// MakeAllSearchable indexes all models of the given type in chunks.
func MakeAllSearchable(ctx context.Context, models []contract.Searchable, engine contract.Engine, chunkSize int) error {
	if len(models) == 0 {
		return nil
	}

	if chunkSize <= 0 {
		chunkSize = 500
	}

	for i := 0; i < len(models); i += chunkSize {
		end := i + chunkSize

		if end > len(models) {
			end = len(models)
		}

		chunk := models[i:end]

		// Filter to only searchable models.
		var searchable []contract.Searchable

		for _, m := range chunk {
			if m.ShouldBeSearchable() {
				searchable = append(searchable, m)
			}
		}

		if len(searchable) == 0 {
			continue
		}

		if err := engine.Update(ctx, searchable); err != nil {
			return err
		}
	}

	return nil
}

// MakeSearchable indexes the given models in the search engine.
func MakeSearchable(ctx context.Context, models []contract.Searchable, engine contract.Engine) error {
	if len(models) == 0 {
		return nil
	}

	var searchable []contract.Searchable

	for _, m := range models {
		if m.ShouldBeSearchable() {
			searchable = append(searchable, m)
		}
	}

	if len(searchable) == 0 {
		return nil
	}

	return engine.Update(ctx, searchable)
}

// RemoveFromSearch removes the given models from the search engine.
func RemoveFromSearch(ctx context.Context, models []contract.Searchable, engine contract.Engine) error {
	if len(models) == 0 {
		return nil
	}

	return engine.Delete(ctx, models)
}

// EnableSearchSyncing enables automatic search syncing for the given model type.
func EnableSearchSyncing(model contract.Searchable) {
	syncingDisabled.Delete(reflect.TypeOf(model))
}

// DisableSearchSyncing disables automatic search syncing for the given model type.
func DisableSearchSyncing(model contract.Searchable) {
	syncingDisabled.Store(reflect.TypeOf(model), struct{}{})
}

// IsSearchSyncingEnabled reports whether syncing is enabled for the model type.
func IsSearchSyncingEnabled(model contract.Searchable) bool {
	_, disabled := syncingDisabled.Load(reflect.TypeOf(model))

	return !disabled
}

// WithoutSyncingToSearch temporarily disables search syncing for the given
// model type while fn executes, then re-enables it.
func WithoutSyncingToSearch(model contract.Searchable, fn func()) {
	DisableSearchSyncing(model)

	defer EnableSearchSyncing(model)

	fn()
}

// SearchableAs returns the index name for the given model, respecting
// the prefix from the SearchableMixin if available.
func SearchableAs(model contract.Searchable, prefix string) string {
	if prefix != "" {
		return prefix + model.GetTable()
	}

	return model.GetTable()
}

// GetScoutKey returns the search key for the model. By default, this is
// the model's primary key value.
func GetScoutKey(model contract.Searchable) any {
	return model.GetKey()
}

// GetScoutKeyName returns the search key name for the model. By default,
// this is the model's primary key column name.
func GetScoutKeyName(model contract.Searchable) string {
	return model.GetKeyName()
}

// ToSearchableArray converts a searchable model to a map for indexing.
// By default, this returns all model attributes (from GetKey/GetTable/etc).
func ToSearchableArray(model contract.Searchable) map[string]any {
	return model.ToSearchableArray()
}
