package eloquent

import "errors"

var (
	// ErrModelNotFound is returned when a find/first query returns no results.
	ErrModelNotFound = errors.New("eloquent: model not found")
	// ErrMultipleRecords is returned when Sole() finds more than one record.
	ErrMultipleRecords = errors.New("eloquent: multiple records found")
	// ErrMassAssignment is returned when a guarded attribute is mass-assigned.
	ErrMassAssignment = errors.New("eloquent: mass assignment violation")
	// ErrMissingAttribute is returned when accessing an attribute that does not exist.
	ErrMissingAttribute = errors.New("eloquent: missing attribute")
	// ErrLazyLoading is returned when lazy loading is attempted while prevented.
	ErrLazyLoading = errors.New("eloquent: lazy loading is disabled")
	// ErrInvalidCast is returned when an attribute cannot be cast to the target type.
	ErrInvalidCast = errors.New("eloquent: invalid cast")
)
