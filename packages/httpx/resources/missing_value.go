package resources

// MissingValue is a sentinel indicating that a value should be omitted from the
// serialised output.
type MissingValue struct{}

// IsMissing returns true. It satisfies the PotentiallyMissing interface.

// PotentiallyMissing is implemented by values that may signal omission.
type PotentiallyMissing interface {
	IsMissing() bool
}

// MergeValue wraps a map whose entries should be merged into the parent
// resource map rather than nested under a single key.
type MergeValue struct {
	Data map[string]any
}

func (MissingValue) IsMissing() bool { return true }

// IsMissing returns false; merge values are never missing.
func (MergeValue) IsMissing() bool { return false }
