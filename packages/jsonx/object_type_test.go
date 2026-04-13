package jsonx_test

import (
	"testing"

	"github.com/bedrock/packages/jsonx"
)

func TestObjectNoProperties(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object().Title("Payload")
	result := schema.ToMap()

	assertEqual(t, result["type"], "object")
	assertEqual(t, result["title"], "Payload")

	if _, ok := result["properties"]; ok {
		t.Fatal("expected no properties key for empty object")
	}
}

func TestObjectWithClosure(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object(func(f jsonx.Factory) map[string]jsonx.SchemaType {
		return map[string]jsonx.SchemaType{}
	}).Title("Payload")

	result := schema.ToMap()

	assertEqual(t, result["type"], "object")
	assertEqual(t, result["title"], "Payload")

	if _, ok := result["properties"]; ok {
		t.Fatal("expected no properties key for empty closure result")
	}
}

func TestObjectWithProperties(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object(map[string]jsonx.SchemaType{
		"age-a": jsonx.Integer().Min(0).Required(),
		"age-b": jsonx.Integer().Default(30).Max(45),
	}).Description("Root object")

	result := schema.ToMap()

	assertEqual(t, result["type"], "object")
	assertEqual(t, result["description"], "Root object")

	props, ok := result["properties"].(map[string]any)

	if !ok {
		t.Fatal("expected properties to be map[string]any")
	}

	ageA, ok := props["age-a"].(map[string]any)

	if !ok {
		t.Fatal("expected age-a to be map[string]any")
	}

	assertEqual(t, ageA["type"], "integer")
	assertEqual(t, ageA["minimum"], 0)

	ageB, ok := props["age-b"].(map[string]any)

	if !ok {
		t.Fatal("expected age-b to be map[string]any")
	}

	assertEqual(t, ageB["type"], "integer")
	assertEqual(t, ageB["default"], 30)
	assertEqual(t, ageB["maximum"], 45)

	required, ok := result["required"].([]string)

	if !ok {
		t.Fatal("expected required to be []string")
	}

	assertContains(t, required, "age-a")
}

func TestObjectWithClosureAndProperties(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object(func(f jsonx.Factory) map[string]jsonx.SchemaType {
		return map[string]jsonx.SchemaType{
			"age-a": f.Integer().Min(0).Required(),
			"age-b": f.Integer().Default(30).Max(45),
		}
	}).Description("Root object")

	result := schema.ToMap()

	assertEqual(t, result["type"], "object")
	assertEqual(t, result["description"], "Root object")

	props, ok := result["properties"].(map[string]any)

	if !ok {
		t.Fatal("expected properties to be map[string]any")
	}

	ageA, ok := props["age-a"].(map[string]any)

	if !ok {
		t.Fatal("expected age-a to be map[string]any")
	}

	assertEqual(t, ageA["type"], "integer")
	assertEqual(t, ageA["minimum"], 0)

	required, ok := result["required"].([]string)

	if !ok {
		t.Fatal("expected required to be []string")
	}

	assertContains(t, required, "age-a")
}

func TestObjectWithoutAdditionalProperties(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object().Default(map[string]any{"age": 1}).WithoutAdditionalProperties()
	result := schema.ToMap()

	assertEqual(t, result["type"], "object")
	assertMapEqual(t, result["default"], map[string]any{"age": 1})
	assertEqual(t, result["additionalProperties"], false)
}

func TestObjectEnum(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object().Enum([]any{
		map[string]any{"a": 1},
		map[string]any{"a": 2},
	})

	result := schema.ToMap()

	assertEqual(t, result["type"], "object")

	enum, ok := result["enum"].([]any)

	if !ok {
		t.Fatal("expected enum to be []any")
	}

	if len(enum) != 2 {
		t.Fatalf("expected 2 enum values, got %d", len(enum))
	}
}
