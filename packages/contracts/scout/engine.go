package search

import "context"

// Engine is the contract for search engine backends.
type Engine interface {
	// Update indexes or updates the given models in the search engine.
	Update(ctx context.Context, models []Searchable) error
	// Delete removes the given models from the search engine.
	Delete(ctx context.Context, models []Searchable) error
	// Search performs a search query and returns raw engine results.
	Search(ctx context.Context, builder SearchBuilder) (any, error)
	// Paginate performs a paginated search query.
	Paginate(ctx context.Context, builder SearchBuilder, perPage, page int) (any, error)
	// MapIds extracts model IDs from raw search results.
	MapIds(results any) []any
	// Map hydrates models from raw search results, using the given model as prototype.
	Map(ctx context.Context, results any, model Searchable) ([]Searchable, error)
	// LazyMap returns a lazy iterator over hydrated models from search results.
	LazyMap(ctx context.Context, results any, model Searchable) func(yield func(Searchable) bool)
	// GetTotalCount returns the total count of results from a raw engine response.
	GetTotalCount(results any) int64
	// Flush removes all records for the given model type from the index.
	Flush(ctx context.Context, model Searchable) error
	// CreateIndex creates a search index with optional settings.
	CreateIndex(ctx context.Context, name string, options map[string]any) error
	// DeleteIndex deletes a search index.
	DeleteIndex(ctx context.Context, name string) error
}

// UpdatesIndexSettings is implemented by engines that support
// configuring index settings (e.g. Algolia, Meilisearch).
type UpdatesIndexSettings interface {
	// UpdateIndexSettings updates the settings for the given index.
	UpdateIndexSettings(ctx context.Context, name string, settings map[string]any) error
	// GetIndexSettings returns the current settings for the given index.
	GetIndexSettings(ctx context.Context, name string) (map[string]any, error)
}

// PaginatesUsingDatabase is implemented by engines that paginate
// results using the database directly (e.g. DatabaseEngine).
type PaginatesUsingDatabase interface {
	// PaginateUsingDatabase performs paginated search using database queries.
	PaginateUsingDatabase(ctx context.Context, builder SearchBuilder, perPage, page int) (any, error)
	// SimplePaginateUsingDatabase performs simple pagination using database queries.
	SimplePaginateUsingDatabase(ctx context.Context, builder SearchBuilder, perPage, page int) (any, error)
}
