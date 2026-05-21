package jsonx

// Factory creates JSON Schema types.
// object construction via Object(func(Factory) map[string]SchemaType{}).
type Factory struct{}

// Object creates a new object schema type with the given properties.
func (f Factory) Object(properties map[string]SchemaType) *ObjectType {
	return newObjectType(properties)
}

// Array creates a new array schema type.
func (f Factory) Array() *ArrayType {
	return newArrayType()
}

// String creates a new string schema type.
func (f Factory) String() *StringType {
	return newStringType()
}

// Integer creates a new integer schema type.
func (f Factory) Integer() *IntegerType {
	return newIntegerType()
}

// Number creates a new number schema type.
func (f Factory) Number() *NumberType {
	return newNumberType()
}

// Boolean creates a new boolean schema type.
func (f Factory) Boolean() *BooleanType {
	return newBooleanType()
}

// Object creates a new object schema type.
// It accepts either a map[string]SchemaType or a func(Factory) map[string]SchemaType.
// If no arguments are provided, an empty object type is returned.
func Object(properties ...any) *ObjectType {
	if len(properties) == 0 {
		return newObjectType(nil)
	}

	switch v := properties[0].(type) {
	case map[string]SchemaType:
		return newObjectType(v)
	case func(Factory) map[string]SchemaType:
		return newObjectType(v(Factory{}))
	default:
		panic("jsonx.Object: argument must be map[string]SchemaType or func(Factory) map[string]SchemaType")
	}
}

// Array creates a new array schema type.
func Array() *ArrayType {
	return newArrayType()
}

// String creates a new string schema type.
func String() *StringType {
	return newStringType()
}

// Integer creates a new integer schema type.
func Integer() *IntegerType {
	return newIntegerType()
}

// Number creates a new number schema type.
func Number() *NumberType {
	return newNumberType()
}

// Boolean creates a new boolean schema type.
func Boolean() *BooleanType {
	return newBooleanType()
}
