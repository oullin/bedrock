package searcherr

import "errors"

// ErrSearchFailed is returned when a search query cannot be executed.
var ErrSearchFailed = errors.New("search: search query failed")
