package cookie

import "errors"

// ErrEmptyName is returned when a cookie with an empty name is queued.
var ErrEmptyName = errors.New("cookie: name must not be empty")
