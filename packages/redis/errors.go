package redis

import "errors"

// Sentinel errors.
var (
	// ErrNil is returned when a command finds no value (null reply).
	ErrNil = errors.New("redis: nil")

	// ErrConnectionNotFound is returned when a connection name is not registered.
	ErrConnectionNotFound = errors.New("redis: connection not registered")

	// ErrDriverNotFound is returned when a driver name is not registered.
	ErrDriverNotFound = errors.New("redis: driver not registered")

	// ErrUnexpectedReply is returned when a reply has a type the helper
	// cannot coerce.
	ErrUnexpectedReply = errors.New("redis: unexpected reply type")

	// ErrLimiterTimeout is returned when a limiter fails to acquire within
	// the block duration.
	ErrLimiterTimeout = errors.New("redis: limiter timeout")

	// ErrClosed is returned when an operation is attempted on a closed
	// connection.
	ErrClosed = errors.New("redis: connection closed")
)
