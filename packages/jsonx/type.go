package jsonx

import "encoding/json"

// SchemaType is the interface satisfied by all JSON Schema type builders.
type SchemaType interface {
	ToMap() map[string]any
	String() string
}

// TypeBuilder holds JSON Schema metadata common to all types.
// It uses generics so that chainable methods return the concrete type pointer.
type TypeBuilder[T any] struct {
	self        *T
	required    *bool
	title       *string
	description *string
	defaultVal  any
	enum        []any
	nullable    *bool
}

// Required indicates that this property is required within its parent object.
func (b *TypeBuilder[T]) Required(required ...bool) *T {
	v := true

	if len(required) > 0 {
		v = required[0]
	}

	if v {
		b.required = &v
	}

	return b.self
}

// Nullable indicates that the type accepts null values.
func (b *TypeBuilder[T]) Nullable(nullable ...bool) *T {
	v := true

	if len(nullable) > 0 {
		v = nullable[0]
	}

	if v {
		b.nullable = &v
	}

	return b.self
}

// Title sets the schema title.
func (b *TypeBuilder[T]) Title(value string) *T {
	b.title = &value

	return b.self
}

// Description sets the schema description.
func (b *TypeBuilder[T]) Description(value string) *T {
	b.description = &value

	return b.self
}

// Enum restricts the value to one of the provided values.
func (b *TypeBuilder[T]) Enum(values []any) *T {
	b.enum = values

	return b.self
}

// ToMap converts the type to its map representation by delegating to Serialize.
func (b *TypeBuilder[T]) ToMap() map[string]any {
	result, err := Serialize(b.self)

	if err != nil {
		panic(err)
	}

	return result
}

// String converts the type to its JSON string representation.
func (b *TypeBuilder[T]) String() string {
	data, err := json.MarshalIndent(b.ToMap(), "", "    ")

	if err != nil {
		return ""
	}

	return string(data)
}

// isRequired reports whether the type is marked as required.
func (b *TypeBuilder[T]) isRequired() bool {
	return b.required != nil && *b.required
}

// isNullable reports whether the type is marked as nullable.
func (b *TypeBuilder[T]) isNullable() bool {
	return b.nullable != nil && *b.nullable
}
