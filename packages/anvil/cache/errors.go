package cache

import "errors"

var (
	ErrNotFound     = errors.New("cache: key not found")
	ErrInvalidValue = errors.New("cache: value is not numeric")
)
