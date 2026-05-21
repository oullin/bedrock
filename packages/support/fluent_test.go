package support

import (
	"encoding/json"
	"iter"
	"testing"
)

// Ref: @bedrock/code-0376
func TestFluentAttributesSetByConstructor(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"name": "Taylor", "age": 30})

	if f.Get("name") != "Taylor" {
		t.Errorf("expected 'Taylor', got %v", f.Get("name"))
	}

	if f.Get("age") != 30 {
		t.Errorf("expected 30, got %v", f.Get("age"))
	}
}

// Ref: @bedrock/code-0376
func TestFluentAttributesSetByConstructorFromStruct(t *testing.T) {
	t.Parallel()

	type sample struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	f := NewFluent(sample{Name: "Taylor", Age: 30})

	if f.Get("name") != "Taylor" || f.Get("age") != 30 {
		t.Fatalf("constructor from struct = %v", f.All())
	}
}

// Ref: @bedrock/code-0376
func TestFluentAttributesSetByConstructorFromIterator(t *testing.T) {
	t.Parallel()

	input := iter.Seq2[string, any](func(yield func(string, any) bool) {
		_ = yield("name", "Taylor")
		_ = yield("age", 30)
	})

	f := NewFluent(input)

	if f.Get("name") != "Taylor" || f.Get("age") != 30 {
		t.Fatalf("constructor from iterator = %v", f.All())
	}
}

// Ref: @bedrock/code-0376
func TestFluentGet(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"name": "Taylor"})

	if got := f.Get("name"); got != "Taylor" {
		t.Errorf("Get existing = %v", got)
	}

	if got := f.Get("missing"); got != nil {
		t.Errorf("Get missing should be nil, got %v", got)
	}

	if got := f.Get("missing", "default"); got != "default" {
		t.Errorf("Get with default = %v", got)
	}
}

// Ref: @bedrock/code-0376
func TestFluentSet(t *testing.T) {
	t.Parallel()

	f := NewFluent()
	f.Set("key", "value")

	if f.Get("key") != "value" {
		t.Errorf("Set then Get = %v", f.Get("key"))
	}
}

// Ref: @bedrock/code-0376
func TestFluentToMap(t *testing.T) {
	t.Parallel()

	attrs := map[string]any{"a": 1, "b": 2}
	f := NewFluent(attrs)
	result := f.ToMap()

	if result["a"] != 1 || result["b"] != 2 {
		t.Errorf("ToMap = %v", result)
	}
}

// Ref: @bedrock/code-0376
func TestFluentToJSON(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"key": "value"})
	data, err := f.ToJSON()

	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}

	var result map[string]any

	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("JSON key = %v", result["key"])
	}
}

// Ref: @bedrock/code-0376
func TestFluentToPrettyJSON(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"key": "value"})
	data, err := f.ToPrettyJSON()

	if err != nil {
		t.Fatalf("ToPrettyJSON error: %v", err)
	}

	if data != "{\n  \"key\": \"value\"\n}" {
		t.Fatalf("ToPrettyJSON = %q", data)
	}
}

// Ref: @bedrock/code-0376
func TestFluentBool(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{
		"t1": "true",
		"t2": "yes",
		"t3": "on",
		"t4": "1",
		"t5": true,
		"f1": "false",
		"f2": false,
	})

	for _, key := range []string{"t1", "t2", "t3", "t4", "t5"} {
		if !f.Bool(key) {
			t.Errorf("Bool(%q) should be true", key)
		}
	}

	for _, key := range []string{"f1", "f2"} {
		if f.Bool(key) {
			t.Errorf("Bool(%q) should be false", key)
		}
	}
}

// Ref: @bedrock/code-0376
func TestFluentInt(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{
		"n1": 42,
		"n2": "100",
		"n3": float64(3.14),
	})

	if f.Int("n1") != 42 {
		t.Errorf("Int(n1) = %d", f.Int("n1"))
	}

	if f.Int("n2") != 100 {
		t.Errorf("Int(n2) = %d", f.Int("n2"))
	}

	if f.Int("n3") != 3 {
		t.Errorf("Int(n3) = %d", f.Int("n3"))
	}

	if f.Int("missing", 99) != 99 {
		t.Errorf("Int(missing, 99) = %d", f.Int("missing", 99))
	}
}

// Ref: @bedrock/code-0376
func TestFluentFloat(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{
		"f1": float64(3.14),
		"f2": "2.71",
	})

	if f.Float("f1") != 3.14 {
		t.Errorf("Float(f1) = %f", f.Float("f1"))
	}

	if f.Float("f2") != 2.71 {
		t.Errorf("Float(f2) = %f", f.Float("f2"))
	}
}

// Ports of:
// Ref: @bedrock/code-0376
func TestFluentIsEmpty(t *testing.T) {
	t.Parallel()

	empty := NewFluent()

	if !empty.IsEmpty() {
		t.Error("empty Fluent should be empty")
	}

	if empty.IsNotEmpty() {
		t.Error("empty Fluent should not be not-empty")
	}

	nonempty := NewFluent(map[string]any{"a": 1})

	if nonempty.IsEmpty() {
		t.Error("non-empty Fluent should not be empty")
	}

	if !nonempty.IsNotEmpty() {
		t.Error("non-empty Fluent should be not-empty")
	}
}

// Ref: @bedrock/code-0376
func TestFluentScope(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{
		"meta.name": "Taylor",
		"meta.role": "admin",
		"plain":     "value",
	})

	scoped := f.Scope("meta")

	if scoped.Get("name") != "Taylor" {
		t.Errorf("Scope(meta).Get(name) = %v", scoped.Get("name"))
	}

	if scoped.Get("role") != "admin" {
		t.Errorf("Scope(meta).Get(role) = %v", scoped.Get("role"))
	}

	if scoped.Get("plain") != nil {
		t.Errorf("Scope(meta) should exclude unscoped keys, got %v", scoped.Get("plain"))
	}
}

// Ref: @bedrock/code-0376
func TestFluentString(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{
		"name":  "Taylor",
		"bytes": []byte("Bedrock"),
	})

	if got := f.String("name"); got != "Taylor" {
		t.Errorf("String(name) = %q", got)
	}

	if got := f.String("bytes"); got != "Bedrock" {
		t.Errorf("String(bytes) = %q", got)
	}

	if got := f.String("missing", "default"); got != "default" {
		t.Errorf("String(missing, default) = %q", got)
	}
}

// Ref: @bedrock/code-0376
func TestFluentOnly(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"a": 1, "b": 2, "c": 3})
	only := f.Only("a", "c")

	if len(only) != 2 || only["a"] != 1 || only["c"] != 3 {
		t.Errorf("Only(a,c) = %v", only)
	}

	if _, exists := only["b"]; exists {
		t.Error("Only should not include 'b'")
	}
}

// Ref: @bedrock/code-0376
func TestFluentExcept(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"a": 1, "b": 2, "c": 3})
	except := f.Except("b")

	if _, exists := except["b"]; exists {
		t.Error("Except should exclude 'b'")
	}

	if except["a"] != 1 || except["c"] != 3 {
		t.Errorf("Except = %v", except)
	}
}

// Ref: @bedrock/code-0376
func TestFluentFill(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"a": 1})
	f.Fill(map[string]any{"b": 2, "a": 99}) // "a" should be overwritten by Fill

	if f.Get("b") != 2 {
		t.Errorf("Fill should add 'b', got %v", f.Get("b"))
	}

	if f.Get("a") != 99 {
		t.Errorf("Fill should overwrite 'a', got %v", f.Get("a"))
	}
}

// Ref: @bedrock/code-0376
func TestFluentMerge(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"a": 1})
	f.Merge(map[string]any{"b": 2, "a": 99}) // "a" should NOT be overwritten

	if f.Get("b") != 2 {
		t.Errorf("Merge should add 'b', got %v", f.Get("b"))
	}

	if f.Get("a") != 1 {
		t.Errorf("Merge should not overwrite 'a', got %v", f.Get("a"))
	}
}

// Ref: @bedrock/code-0376
func TestFluentAll(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"x": 10, "y": 20})
	all := f.All()

	if len(all) != 2 || all["x"] != 10 || all["y"] != 20 {
		t.Errorf("All() = %v", all)
	}
}

func TestFluentHasMissing(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"present": true})

	if !f.Has("present") {
		t.Error("Has should return true for existing key")
	}

	if f.Has("absent") {
		t.Error("Has should return false for missing key")
	}

	if !f.Missing("absent") {
		t.Error("Missing should return true for absent key")
	}
}

// Ref: @bedrock/code-0376
func TestFluentArray(t *testing.T) {
	t.Parallel()

	f := NewFluent(map[string]any{"a": 1})
	result := f.Array()

	if len(result) != 1 || result["a"] != 1 {
		t.Fatalf("Array() = %v", result)
	}
}
