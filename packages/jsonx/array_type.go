package jsonx

// ArrayType represents a JSON Schema array type.
type ArrayType struct {
	TypeBuilder[ArrayType]
	minItems    *int
	maxItems    *int
	items       SchemaType
	uniqueItems *bool
}

// newArrayType creates a new ArrayType instance.
func newArrayType() *ArrayType {
	t := &ArrayType{}
	t.self = t

	return t
}

// Min sets the minimum number of items (inclusive).
func (t *ArrayType) Min(value int) *ArrayType {
	t.minItems = &value

	return t
}

// Max sets the maximum number of items (inclusive).
func (t *ArrayType) Max(value int) *ArrayType {
	t.maxItems = &value

	return t
}

// Items sets the schema for array items.
func (t *ArrayType) Items(schema SchemaType) *ArrayType {
	t.items = schema

	return t
}

// Unique indicates that the array items must be unique.
func (t *ArrayType) Unique(unique ...bool) *ArrayType {
	v := true

	if len(unique) > 0 {
		v = unique[0]
	}

	if v {
		t.uniqueItems = &v
	}

	return t
}

// Default sets the default value for the array type.
func (t *ArrayType) Default(value []any) *ArrayType {
	t.defaultVal = value

	return t
}
