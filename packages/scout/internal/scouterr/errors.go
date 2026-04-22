package scouterr

import "errors"

// ErrSearchFailed is returned when a search query cannot be executed.
var ErrSearchFailed = errors.New("scout: search query failed")
