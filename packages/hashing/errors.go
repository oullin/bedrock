package hashing

import "errors"

var (
	ErrPasswordTooLong   = errors.New("hashing: password exceeds maximum length")
	ErrUnsupportedDriver = errors.New("hashing: unsupported driver")
	ErrInvalidHash       = errors.New("hashing: the hash is invalid or malformed")
	ErrAlgorithmMismatch = errors.New("hashing: the password does not use the expected algorithm")
)
