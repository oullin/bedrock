package scout

// Searchable is the contract that models must satisfy to participate in
// Scout search. In Go, this replaces the Laravel Searchable trait.
type Searchable interface {
	// GetScoutKey returns the value used as the document ID in the search index.
	GetScoutKey() any
	// GetScoutKeyName returns the attribute name used as the document ID.
	GetScoutKeyName() string
	// SearchableAs returns the index name for this model.
	SearchableAs() string
	// ToSearchableArray converts the model to a map for indexing.
	ToSearchableArray() map[string]any
	// ShouldBeSearchable reports whether this model instance should be indexed.
	ShouldBeSearchable() bool
	// SearchIndexShouldBeUpdated reports whether a save should trigger re-indexing.
	SearchIndexShouldBeUpdated() bool
	// GetScoutMetadata returns engine-specific metadata for the model.
	GetScoutMetadata() map[string]any
	// WithScoutMetadata sets a metadata key-value pair on the model.
	WithScoutMetadata(key string, value any) Searchable

	// GetTable returns the table name. Inherited from eloquent.Model.
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
