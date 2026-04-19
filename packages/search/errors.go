package search

import "errors"

var (
	// ErrEngineNotConfigured is returned when a named engine has not been registered.
	ErrEngineNotConfigured = errors.New("search: engine is not configured")
	// ErrDriverNotSupported is returned when a driver name is unrecognised.
	ErrDriverNotSupported = errors.New("search: driver is not supported")
	// ErrModelNotSearchable is returned when a model does not implement Searchable.
	ErrModelNotSearchable = errors.New("search: model does not implement Searchable")
	// ErrIndexNotFound is returned when a search index cannot be located.
	ErrIndexNotFound = errors.New("search: index not found")
	// ErrSearchFailed is returned when a search query cannot be executed.
	ErrSearchFailed = errors.New("search: search query failed")
	// ErrIndexingFailed is returned when an indexing operation fails.
	ErrIndexingFailed = errors.New("search: indexing operation failed")
	// ErrFlushFailed is returned when a flush operation fails.
	ErrFlushFailed = errors.New("search: flush operation failed")
)
