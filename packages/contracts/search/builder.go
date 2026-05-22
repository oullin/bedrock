package search

// SearchBuilder exposes the read surface of a Search builder to engines.
type SearchBuilder interface {
	// GetQuery returns the search query string.
	GetQuery() string
	// GetModel returns the model being searched.
	GetModel() Searchable
	// GetIndex returns the custom index name, or empty string for default.
	GetIndex() string
	// GetLimit returns the result limit, or 0 for no limit.
	GetLimit() int
	// GetWheres returns the where constraints as key-value pairs.
	GetWheres() map[string]any
	// GetWhereIns returns the whereIn constraints.
	GetWhereIns() map[string][]any
	// GetWhereNotIns returns the whereNotIn constraints.
	GetWhereNotIns() map[string][]any
	// GetOrders returns the order-by directives.
	GetOrders() []Order
	// GetOptions returns engine-specific options.
	GetOptions() map[string]any
	// GetCallback returns the engine-specific query callback, or nil.
	GetCallback() func(Engine, string, map[string]any) any
	// HasCallback reports whether a callback was set.
	HasCallback() bool
}

// Order represents a sort directive.
type Order struct {
	Column    string
	Direction string // "asc" or "desc"
}
