package jsonx

// NumberType represents a JSON Schema number type.
type NumberType struct {
	TypeBuilder[NumberType]
	minimum    *float64
	maximum    *float64
	multipleOf *float64
}

// newNumberType creates a new NumberType instance.
func newNumberType() *NumberType {
	t := &NumberType{}
	t.self = t

	return t
}

// Min sets the minimum value (inclusive).
func (t *NumberType) Min(value float64) *NumberType {
	t.minimum = &value

	return t
}

// Max sets the maximum value (inclusive).
func (t *NumberType) Max(value float64) *NumberType {
	t.maximum = &value

	return t
}

// MultipleOf sets the number the value must be a multiple of.
func (t *NumberType) MultipleOf(value float64) *NumberType {
	t.multipleOf = &value

	return t
}

// Default sets the default value for the number type.
func (t *NumberType) Default(value float64) *NumberType {
	t.defaultVal = value

	return t
}
