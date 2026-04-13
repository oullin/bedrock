package jsonx

// ObjectType represents a JSON Schema object type.
type ObjectType struct {
	TypeBuilder[ObjectType]
	additionalProperties *bool
	properties           map[string]SchemaType
}

// newObjectType creates a new ObjectType instance with the given properties.
func newObjectType(properties map[string]SchemaType) *ObjectType {
	t := &ObjectType{
		properties: properties,
	}
	t.self = t

	return t
}

// WithoutAdditionalProperties disallows additional properties beyond those defined.
func (t *ObjectType) WithoutAdditionalProperties() *ObjectType {
	v := false
	t.additionalProperties = &v

	return t
}

// Default sets the default value for the object type.
func (t *ObjectType) Default(value map[string]any) *ObjectType {
	t.defaultVal = value

	return t
}
