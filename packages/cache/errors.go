package cache

import "errors"

var (
	// ErrNotFound is returned when a key does not exist or has expired.
	ErrNotFound = errors.New("cache: key not found")
	// ErrInvalidValue is returned when a non-numeric value is incremented.
	ErrInvalidValue = errors.New("cache: value is not numeric")
	// ErrLockTimeout is returned when a lock cannot be acquired within the
	// configured timeout.
	ErrLockTimeout = errors.New("cache: lock timeout")
)
