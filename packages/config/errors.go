package config

import "errors"

var (
	// ErrInvalidType is returned when a config value does not match the
	// requested type.
	ErrInvalidType = errors.New("config: invalid type")
)
