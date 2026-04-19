package scout

import "errors"

var (
	// ErrEngineNotConfigured is returned when a named engine has not been registered.
	ErrEngineNotConfigured = errors.New("scout: engine is not configured")
	// ErrDriverNotSupported is returned when a driver name is unrecognised.
	ErrDriverNotSupported = errors.New("scout: driver is not supported")
	// ErrModelNotSearchable is returned when a model does not implement Searchable.
	ErrModelNotSearchable = errors.New("scout: model does not implement Searchable")
	// ErrIndexNotFound is returned when a search index cannot be located.
	ErrIndexNotFound = errors.New("scout: index not found")
	// ErrSearchFailed is returned when a search query cannot be executed.
	ErrSearchFailed = errors.New("scout: search query failed")
	// ErrIndexingFailed is returned when an indexing operation fails.
	ErrIndexingFailed = errors.New("scout: indexing operation failed")
	// ErrFlushFailed is returned when a flush operation fails.
	ErrFlushFailed = errors.New("scout: flush operation failed")
)
