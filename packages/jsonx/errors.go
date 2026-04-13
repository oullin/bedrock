package jsonx

import "errors"

// ErrUnknownType is returned when attempting to serialize an unrecognized schema type.
var ErrUnknownType = errors.New("jsonx: unknown schema type")
