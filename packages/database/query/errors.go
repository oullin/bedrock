package query

import "errors"

var (
	// ErrInvalidOperator is returned when an unsupported comparison operator is used.
	ErrInvalidOperator = errors.New("query: invalid operator")
	// ErrInvalidBinding is returned when a binding type is unsupported.
	ErrInvalidBinding = errors.New("query: invalid binding type")
	// ErrEmptyColumns is returned when a select or insert has no columns.
	ErrEmptyColumns = errors.New("query: empty columns")
)
