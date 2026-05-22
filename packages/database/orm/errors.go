package orm

import "errors"

var (
	// ErrModelNotFound is returned when a find/first query returns no results.
	ErrModelNotFound = errors.New("orm: model not found")
	// ErrMultipleRecords is returned when Sole() finds more than one record.
	ErrMultipleRecords = errors.New("orm: multiple records found")
	// ErrMassAssignment is returned when a guarded attribute is mass-assigned.
	ErrMassAssignment = errors.New("orm: mass assignment violation")
	// ErrMissingAttribute is returned when accessing an attribute that does not exist.
	ErrMissingAttribute = errors.New("orm: missing attribute")
	// ErrLazyLoading is returned when lazy loading is attempted while prevented.
	ErrLazyLoading = errors.New("orm: lazy loading is disabled")
	// ErrInvalidCast is returned when an attribute cannot be cast to the target type.
	ErrInvalidCast = errors.New("orm: invalid cast")
)
