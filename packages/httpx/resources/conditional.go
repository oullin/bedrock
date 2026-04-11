package resources

// ConditionalValue wraps a value that is included in the resource output only
// when its condition is true.
type ConditionalValue struct {
	Condition bool
	Value     any
}

// IsMissing returns true when the condition is false, signalling omission.
func (c ConditionalValue) IsMissing() bool { return !c.Condition }

// When returns a ConditionalValue that is included when condition is true.
func When(condition bool, value any) ConditionalValue {
	return ConditionalValue{Condition: condition, Value: value}
}

// Unless returns a ConditionalValue that is included when condition is false.
func Unless(condition bool, value any) ConditionalValue {
	return ConditionalValue{Condition: !condition, Value: value}
}

// MergeWhen returns a MergeValue that merges its entries into the parent map
// only when condition is true. When false, a MissingValue is returned.
func MergeWhen(condition bool, data map[string]any) any {
	if condition {
		return MergeValue{Data: data}
	}

	return MissingValue{}
}
