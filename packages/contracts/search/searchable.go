package search

// Searchable is the contract that models must satisfy to participate in
// Search search. In Go, this replaces the upstream Searchable trait.
type Searchable interface {
	// GetSearchKey returns the value used as the document ID in the search index.
	GetSearchKey() any
	// GetSearchKeyName returns the attribute name used as the document ID.
	GetSearchKeyName() string
	// SearchableAs returns the index name for this model.
	SearchableAs() string
	// ToSearchableArray converts the model to a map for indexing.
	ToSearchableArray() map[string]any
	// ShouldBeSearchable reports whether this model instance should be indexed.
	ShouldBeSearchable() bool
	// SearchIndexShouldBeUpdated reports whether a save should trigger re-indexing.
	SearchIndexShouldBeUpdated() bool
	// GetSearchMetadata returns engine-specific metadata for the model.
	GetSearchMetadata() map[string]any
	// WithSearchMetadata sets a metadata key-value pair on the model.
	WithSearchMetadata(key string, value any) Searchable

	// GetTable returns the table name. Inherited from orm.Model.
	GetTable() string
	// GetKeyName returns the primary key column name.
	GetKeyName() string
	// GetKey returns the value of the primary key.
	GetKey() any
	// GetConnectionName returns the database connection name for this model.
	GetConnectionName() string
	// UsesSoftDelete reports whether the model uses soft deletes.
	UsesSoftDelete() bool
}
