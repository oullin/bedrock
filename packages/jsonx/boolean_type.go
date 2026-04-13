package jsonx

// BooleanType represents a JSON Schema boolean type.
type BooleanType struct {
	TypeBuilder[BooleanType]
}

// newBooleanType creates a new BooleanType instance.
func newBooleanType() *BooleanType {
	t := &BooleanType{}
	t.self = t

	return t
}

// Default sets the default value for the boolean type.
func (t *BooleanType) Default(value bool) *BooleanType {
	t.defaultVal = value

	return t
}
