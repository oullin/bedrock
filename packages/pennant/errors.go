package pennant

import "errors"

var (
	// ErrFeatureNotDefined is returned when Get is called for a feature with no
	// registered resolver and no stored state.
	ErrFeatureNotDefined = errors.New("pennant: feature not defined")

	// ErrUnserializableScope is returned by SerializeScope when the scope value
	// cannot be converted to a stable string key.
	ErrUnserializableScope = errors.New("pennant: scope cannot be serialized")

	// ErrDriverNotFound is returned by Manager when a named driver has not been
	// registered.
	ErrDriverNotFound = errors.New("pennant: driver not registered")

	// ErrStorageConflict is returned by DatabaseDriver when an upsert conflict
	// cannot be resolved within the maximum number of retries.
	ErrStorageConflict = errors.New("pennant: storage conflict after max retries")
)
