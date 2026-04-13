package pagination

import "errors"

var (
	// ErrInvalidPage indicates the page number is not valid (must be >= 1).
	ErrInvalidPage = errors.New("pagination: invalid page number")

	// ErrInvalidPerPage indicates the per-page value is not valid (must be >= 1).
	ErrInvalidPerPage = errors.New("pagination: invalid per page value")
)
