package filesystem

import "errors"

var (
	ErrNotFound      = errors.New("filesystem: file not found")
	ErrNotDirectory  = errors.New("filesystem: not a directory")
	ErrNotFile       = errors.New("filesystem: not a file")
	ErrHashAlgorithm = errors.New("filesystem: unsupported hash algorithm")
	ErrLockFailed    = errors.New("filesystem: failed to acquire lock")
)
