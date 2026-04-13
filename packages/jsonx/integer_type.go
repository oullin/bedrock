package jsonx

// IntegerType represents a JSON Schema integer type.
type IntegerType struct {
	TypeBuilder[IntegerType]
	minimum    *int
	maximum    *int
	multipleOf *int
}

// newIntegerType creates a new IntegerType instance.
func newIntegerType() *IntegerType {
	t := &IntegerType{}
	t.self = t

	return t
}

// Min sets the minimum value (inclusive).
func (t *IntegerType) Min(value int) *IntegerType {
	t.minimum = &value

	return t
}

// Max sets the maximum value (inclusive).
func (t *IntegerType) Max(value int) *IntegerType {
	t.maximum = &value

	return t
}

// MultipleOf sets the number the value must be a multiple of.
func (t *IntegerType) MultipleOf(value int) *IntegerType {
	t.multipleOf = &value

	return t
}

// Default sets the default value for the integer type.
func (t *IntegerType) Default(value int) *IntegerType {
	t.defaultVal = value

	return t
}
