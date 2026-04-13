package jsonx_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/bedrock/packages/jsonx"
)

func TestTypeToMap(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object(map[string]jsonx.SchemaType{
		"age": jsonx.Integer().Min(0).Required(),
	}).Title("User").Description("User payload").Default(map[string]any{"age": 20})

	result := schema.ToMap()

	assertEqual(t, result["type"], "object")
	assertEqual(t, result["title"], "User")
	assertEqual(t, result["description"], "User payload")
	assertMapEqual(t, result["default"], map[string]any{"age": 20})

	props, ok := result["properties"].(map[string]any)

	if !ok {
		t.Fatal("expected properties to be map[string]any")
	}

	age, ok := props["age"].(map[string]any)

	if !ok {
		t.Fatal("expected age to be map[string]any")
	}

	assertEqual(t, age["type"], "integer")
	assertEqual(t, age["minimum"], 0)

	required, ok := result["required"].([]string)

	if !ok {
		t.Fatal("expected required to be []string")
	}

	assertContains(t, required, "age")
}

func TestTypeString(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object(map[string]jsonx.SchemaType{
		"age": jsonx.Integer().Min(0).Required(),
	}).Title("User")

	str := schema.String()

	var parsed map[string]any

	if err := json.Unmarshal([]byte(str), &parsed); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	assertEqual(t, parsed["type"], "object")
	assertEqual(t, parsed["title"], "User")
}

func TestTypeStringableRepresentation(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object(map[string]jsonx.SchemaType{
		"age": jsonx.Integer().Min(0).Required(),
	}).Description("Payload")

	str := schema.String()

	var parsed map[string]any

	if err := json.Unmarshal([]byte(str), &parsed); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	assertEqual(t, parsed["type"], "object")
	assertEqual(t, parsed["description"], "Payload")
}

func TestTypesInObjectSchemaViaClosure(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object(func(f jsonx.Factory) map[string]jsonx.SchemaType {
		return map[string]jsonx.SchemaType{
			"name": f.String().Required(),
			"age":  f.Integer().Min(0),
		}
	})

	result := schema.ToMap()

	assertEqual(t, result["type"], "object")

	props, ok := result["properties"].(map[string]any)

	if !ok {
		t.Fatal("expected properties to be map[string]any")
	}

	name, ok := props["name"].(map[string]any)

	if !ok {
		t.Fatal("expected name to be map[string]any")
	}

	assertEqual(t, name["type"], "string")
}

func TestNullableString(t *testing.T) {
	t.Parallel()

	schema := jsonx.String().Nullable()
	result := schema.ToMap()

	typeVal, ok := result["type"].([]any)

	if !ok {
		t.Fatal("expected type to be []any for nullable")
	}

	if len(typeVal) != 2 || typeVal[0] != "string" || typeVal[1] != "null" {
		t.Fatalf("expected [string, null], got %v", typeVal)
	}
}

func TestNullableInteger(t *testing.T) {
	t.Parallel()

	schema := jsonx.Integer().Nullable()
	result := schema.ToMap()

	typeVal, ok := result["type"].([]any)

	if !ok {
		t.Fatal("expected type to be []any for nullable")
	}

	if len(typeVal) != 2 || typeVal[0] != "integer" || typeVal[1] != "null" {
		t.Fatalf("expected [integer, null], got %v", typeVal)
	}
}

func TestNullableNumber(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().Nullable()
	result := schema.ToMap()

	typeVal, ok := result["type"].([]any)

	if !ok {
		t.Fatal("expected type to be []any for nullable")
	}

	if len(typeVal) != 2 || typeVal[0] != "number" || typeVal[1] != "null" {
		t.Fatalf("expected [number, null], got %v", typeVal)
	}
}

func TestNullableBoolean(t *testing.T) {
	t.Parallel()

	schema := jsonx.Boolean().Nullable()
	result := schema.ToMap()

	typeVal, ok := result["type"].([]any)

	if !ok {
		t.Fatal("expected type to be []any for nullable")
	}

	if len(typeVal) != 2 || typeVal[0] != "boolean" || typeVal[1] != "null" {
		t.Fatalf("expected [boolean, null], got %v", typeVal)
	}
}

func TestNullableArray(t *testing.T) {
	t.Parallel()

	schema := jsonx.Array().Nullable()
	result := schema.ToMap()

	typeVal, ok := result["type"].([]any)

	if !ok {
		t.Fatal("expected type to be []any for nullable")
	}

	if len(typeVal) != 2 || typeVal[0] != "array" || typeVal[1] != "null" {
		t.Fatalf("expected [array, null], got %v", typeVal)
	}
}

func TestNullableFalseDoesNotMakeNullable(t *testing.T) {
	t.Parallel()

	schema := jsonx.String().Nullable(false)
	result := schema.ToMap()

	assertEqual(t, result["type"], "string")
}

func TestNestedObjectWithNullableProperty(t *testing.T) {
	t.Parallel()

	schema := jsonx.Object(map[string]jsonx.SchemaType{
		"age": jsonx.Integer().Nullable(),
	})

	result := schema.ToMap()

	props, ok := result["properties"].(map[string]any)

	if !ok {
		t.Fatal("expected properties to be map[string]any")
	}

	age, ok := props["age"].(map[string]any)

	if !ok {
		t.Fatal("expected age to be map[string]any")
	}

	typeVal, ok := age["type"].([]any)

	if !ok {
		t.Fatal("expected age type to be []any for nullable")
	}

	if len(typeVal) != 2 || typeVal[0] != "integer" || typeVal[1] != "null" {
		t.Fatalf("expected [integer, null], got %v", typeVal)
	}
}

// --- Test helpers ---

func assertEqual(t *testing.T, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v (%T), got %v (%T)", want, want, got, got)
	}
}

func assertSliceEqual(t *testing.T, got any, want []any) {
	t.Helper()

	gotSlice, ok := got.([]any)

	if !ok {
		t.Fatalf("expected []any, got %T", got)
	}

	if !reflect.DeepEqual(gotSlice, want) {
		t.Fatalf("expected %v, got %v", want, gotSlice)
	}
}

func assertMapEqual(t *testing.T, got any, want map[string]any) {
	t.Helper()

	gotMap, ok := got.(map[string]any)

	if !ok {
		t.Fatalf("expected map[string]any, got %T", got)
	}

	if !reflect.DeepEqual(gotMap, want) {
		t.Fatalf("expected %v, got %v", want, gotMap)
	}
}

func assertContains(t *testing.T, slice []string, value string) {
	t.Helper()

	for _, v := range slice {
		if v == value {
			return
		}
	}

	t.Fatalf("expected slice %v to contain %q", slice, value)
}
