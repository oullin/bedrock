package schema

import "errors"

var (
	// ErrTableNotFound is returned when the requested table does not exist.
	ErrTableNotFound = errors.New("schema: table not found")
	// ErrColumnNotFound is returned when the requested column does not exist.
	ErrColumnNotFound = errors.New("schema: column not found")
)
