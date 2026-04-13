package jsonx

import "fmt"

// Serialize converts a SchemaType into its JSON Schema map representation.
func Serialize(t any) (map[string]any, error) {
	var (
		typeName string
		attrs    = make(map[string]any)
	)

	switch v := t.(type) {
	case *ArrayType:
		typeName = "array"
		serializeBase(&v.TypeBuilder, attrs)
		serializeArrayFields(v, attrs)
	case *BooleanType:
		typeName = "boolean"
		serializeBase(&v.TypeBuilder, attrs)
	case *IntegerType:
		typeName = "integer"
		serializeBase(&v.TypeBuilder, attrs)
		serializeIntegerFields(v, attrs)
	case *NumberType:
		typeName = "number"
		serializeBase(&v.TypeBuilder, attrs)
		serializeNumberFields(v, attrs)
	case *ObjectType:
		typeName = "object"
		serializeBase(&v.TypeBuilder, attrs)
		serializeObjectFields(v, attrs)
	case *StringType:
		typeName = "string"
		serializeBase(&v.TypeBuilder, attrs)
		serializeStringFields(v, attrs)
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnknownType, t)
	}

	if isNullable(t) {
		attrs["type"] = []any{typeName, "null"}
	} else {
		attrs["type"] = typeName
	}

	return attrs, nil
}

func serializeBase[T any](b *TypeBuilder[T], attrs map[string]any) {
	if b.title != nil {
		attrs["title"] = *b.title
	}

	if b.description != nil {
		attrs["description"] = *b.description
	}

	if b.defaultVal != nil {
		attrs["default"] = b.defaultVal
	}

	if b.enum != nil {
		attrs["enum"] = b.enum
	}
}

func serializeStringFields(t *StringType, attrs map[string]any) {
	if t.minLength != nil {
		attrs["minLength"] = *t.minLength
	}

	if t.maxLength != nil {
		attrs["maxLength"] = *t.maxLength
	}

	if t.pattern != nil {
		attrs["pattern"] = *t.pattern
	}

	if t.format != nil {
		attrs["format"] = *t.format
	}
}

func serializeIntegerFields(t *IntegerType, attrs map[string]any) {
	if t.minimum != nil {
		attrs["minimum"] = *t.minimum
	}

	if t.maximum != nil {
		attrs["maximum"] = *t.maximum
	}

	if t.multipleOf != nil {
		attrs["multipleOf"] = *t.multipleOf
	}
}

func serializeNumberFields(t *NumberType, attrs map[string]any) {
	if t.minimum != nil {
		attrs["minimum"] = *t.minimum
	}

	if t.maximum != nil {
		attrs["maximum"] = *t.maximum
	}

	if t.multipleOf != nil {
		attrs["multipleOf"] = *t.multipleOf
	}
}

func serializeArrayFields(t *ArrayType, attrs map[string]any) {
	if t.minItems != nil {
		attrs["minItems"] = *t.minItems
	}

	if t.maxItems != nil {
		attrs["maxItems"] = *t.maxItems
	}

	if t.items != nil {
		items, err := Serialize(t.items)

		if err != nil {
			panic(err)
		}

		attrs["items"] = items
	}

	if t.uniqueItems != nil {
		attrs["uniqueItems"] = *t.uniqueItems
	}
}

func serializeObjectFields(t *ObjectType, attrs map[string]any) {
	if t.additionalProperties != nil {
		attrs["additionalProperties"] = *t.additionalProperties
	}

	if len(t.properties) == 0 {
		return
	}

	properties := make(map[string]any, len(t.properties))

	var required []string

	for key, prop := range t.properties {
		serialized, err := Serialize(prop)

		if err != nil {
			panic(err)
		}

		properties[key] = serialized

		if isRequiredProp(prop) {
			required = append(required, key)
		}
	}

	attrs["properties"] = properties

	if len(required) > 0 {
		attrs["required"] = required
	}
}

func isNullable(t any) bool {
	switch v := t.(type) {
	case *ArrayType:
		return v.TypeBuilder.isNullable()
	case *BooleanType:
		return v.TypeBuilder.isNullable()
	case *IntegerType:
		return v.TypeBuilder.isNullable()
	case *NumberType:
		return v.TypeBuilder.isNullable()
	case *ObjectType:
		return v.TypeBuilder.isNullable()
	case *StringType:
		return v.TypeBuilder.isNullable()
	default:
		return false
	}
}

func isRequiredProp(t SchemaType) bool {
	switch v := t.(type) {
	case *ArrayType:
		return v.TypeBuilder.isRequired()
	case *BooleanType:
		return v.TypeBuilder.isRequired()
	case *IntegerType:
		return v.TypeBuilder.isRequired()
	case *NumberType:
		return v.TypeBuilder.isRequired()
	case *ObjectType:
		return v.TypeBuilder.isRequired()
	case *StringType:
		return v.TypeBuilder.isRequired()
	default:
		return false
	}
}
