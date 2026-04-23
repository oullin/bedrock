package kv_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bedrock/packages/collection/kv"
)

// Framework collections parity markers exercised by this file.
// SupportArrTest::testGet
// SupportArrTest::testSet
// SupportArrTest::testAdd
// SupportArrTest::testHas
// SupportArrTest::testHasAllMethod
// SupportArrTest::testHasAnyMethod
// SupportArrTest::testPull
// SupportArrTest::testForget
// SupportArrTest::testDot
// SupportArrTest::testDotWithDepth
// SupportArrTest::testUndot
// SupportArrTest::testOnly
// SupportArrTest::testExcept
// SupportArrTest::testIsAssoc
// SupportArrTest::testQuery
// SupportArrTest::testToCssClasses
// SupportArrTest::testToCssStyles
// SupportArrTest::testSort
// SupportArrTest::testSortRecursive
// SupportArrTest::testMap
// SupportArrTest::testWhere
// SupportArrTest::testPrependKeysWith
// SupportArrTest::testReplaceArray
func TestGet(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
	}
	v := kv.Get(data, "user.name")

	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}

	v = kv.Get(data, "user.email", "default@test.com")

	if v != "default@test.com" {
		t.Errorf("expected default, got %v", v)
	}

	v = kv.Get(data, "missing.key")

	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}

	v = kv.Get(data, "")

	if v == nil {
		t.Error("expected non-nil for empty key")
	}
}

func TestSet(t *testing.T) {
	data := make(map[string]any)
	kv.Set(data, "user.name", "Alice")

	v := kv.Get(data, "user.name")

	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}
}

func TestSetNoOverwrite(t *testing.T) {
	data := map[string]any{"name": "Alice"}
	kv.Set(data, "name", "Bob", false)

	if data["name"] != "Alice" {
		t.Errorf("expected Alice (not overwritten), got %v", data["name"])
	}
}

func TestHas(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
	}

	if !kv.Has(data, "user.name") {
		t.Error("expected true")
	}

	if kv.Has(data, "user.email") {
		t.Error("expected false")
	}

	if kv.Has(data, "") {
		t.Error("expected false for empty key")
	}
}

func TestFill(t *testing.T) {
	data := map[string]any{
		"name": "Alice",
	}

	kv.Fill(data, "name", "Bob") // Should not overwrite

	if data["name"] != "Alice" {
		t.Errorf("expected Alice (not overwritten), got %v", data["name"])
	}

	kv.Fill(data, "email", "alice@test.com") // Should fill

	if data["email"] != "alice@test.com" {
		t.Errorf("expected email to be filled, got %v", data["email"])
	}
}

func TestAdd(t *testing.T) {
	data := map[string]any{"name": "Alice"}
	kv.Add(data, "name", "Bob")

	if data["name"] != "Alice" {
		t.Errorf("expected Alice, got %v", data["name"])
	}

	kv.Add(data, "email", "alice@test.com")

	if data["email"] != "alice@test.com" {
		t.Errorf("expected email set, got %v", data["email"])
	}
}

func TestForget(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name":  "Alice",
			"email": "alice@test.com",
		},
	}

	kv.Forget(data, "user.email")

	if kv.Has(data, "user.email") {
		t.Error("expected email to be removed")
	}

	if !kv.Has(data, "user.name") {
		t.Error("expected name to still exist")
	}
}

func TestForgetMany(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name":  "Alice",
			"email": "alice@test.com",
		},
	}
	kv.ForgetMany(data, "user.email", "user.name")

	if kv.Has(data, "user.email") || kv.Has(data, "user.name") {
		t.Error("expected both keys removed")
	}
}

func TestPull(t *testing.T) {
	data := map[string]any{"name": "Alice", "age": 25}
	v := kv.Pull(data, "name")

	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}

	if _, ok := data["name"]; ok {
		t.Error("expected name to be removed")
	}
}

func TestHasAll(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
	}

	if !kv.HasAll(data, "user.name") {
		t.Error("expected true")
	}

	if kv.HasAll(data, "user.name", "user.email") {
		t.Error("expected false when one key missing")
	}
}

func TestHasAny(t *testing.T) {
	data := map[string]any{
		"name": "Alice",
	}

	if !kv.HasAny(data, "email", "name") {
		t.Error("expected true")
	}

	if kv.HasAny(data, "email", "phone") {
		t.Error("expected false")
	}
}

func TestDot(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name":  "Alice",
			"email": "alice@test.com",
		},
		"status": "active",
	}
	dotted := kv.Dot(data)

	if dotted["user.name"] != "Alice" {
		t.Errorf("expected Alice, got %v", dotted["user.name"])
	}

	if dotted["status"] != "active" {
		t.Errorf("expected active, got %v", dotted["status"])
	}
}

func TestDotWithPrefix(t *testing.T) {
	data := map[string]any{
		"name": "Alice",
	}
	result := kv.Dot(data, "user")

	if result["user.name"] != "Alice" {
		t.Errorf("expected Alice at 'user.name', got %v", result["user.name"])
	}
}

func TestUndot(t *testing.T) {
	dotted := map[string]any{
		"user.name":  "Alice",
		"user.email": "alice@test.com",
	}
	undotted := kv.Undot(dotted)
	v := kv.Get(undotted, "user.name")

	if v != "Alice" {
		t.Errorf("expected Alice, got %v", v)
	}
}

func TestOnly(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2, "c": 3}
	result := kv.Only(data, "a", "c")

	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}

	if result["a"] != 1 {
		t.Error("expected a=1")
	}
}

func TestExcept(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2, "c": 3}
	result := kv.Except(data, "b")

	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}

	if _, ok := result["b"]; ok {
		t.Error("expected b to be excluded")
	}
}

func TestIsAssoc(t *testing.T) {
	if !kv.IsAssoc(map[string]any{"a": 1}) {
		t.Error("expected true")
	}

	if kv.IsAssoc(map[string]any{}) {
		t.Error("expected false for empty map")
	}
}

func TestQuery(t *testing.T) {
	data := map[string]any{"name": "Alice", "age": 25}
	result := kv.Query(data)

	if !strings.Contains(result, "name=Alice") {
		t.Errorf("expected name=Alice in query, got %s", result)
	}
}

func TestToCssClasses(t *testing.T) {
	classes := map[string]bool{
		"active":   true,
		"disabled": false,
		"visible":  true,
	}
	result := kv.ToCssClasses(classes)

	if !strings.Contains(result, "active") {
		t.Error("expected 'active' in result")
	}

	if strings.Contains(result, "disabled") {
		t.Error("expected 'disabled' to be excluded")
	}
}

func TestToCssStyles(t *testing.T) {
	styles := map[string]bool{
		"color: red":        true,
		"display: none":     false,
		"font-weight: bold": true,
	}
	result := kv.ToCssStyles(styles)

	if !strings.Contains(result, "color: red;") {
		t.Errorf("expected 'color: red;' in result, got '%s'", result)
	}
}

func TestSort(t *testing.T) {
	data := map[string]any{"c": 3, "a": 1, "b": 2}
	result := kv.Sort(data)

	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}
}

func TestSortRecursive(t *testing.T) {
	data := map[string]any{
		"b": 2,
		"a": map[string]any{
			"z": 26,
			"y": 25,
		},
	}
	result := kv.SortRecursive(data)

	if result == nil {
		t.Error("expected non-nil result")
	}

	nested, ok := result["a"].(map[string]any)

	if !ok {
		t.Error("expected nested map")
	}

	_ = nested
}

func TestMap(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2}
	result := kv.Map(data, func(v int, k string) int { return v * 10 })

	if result["a"] != 10 {
		t.Errorf("expected 10, got %d", result["a"])
	}
}

func TestWhere(t *testing.T) {
	data := map[string]int{"a": 1, "b": 2, "c": 3}
	result := kv.Where(data, func(v int, k string) bool { return v > 1 })

	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestPrependKeysWith(t *testing.T) {
	data := map[string]int{"name": 1, "age": 2}
	result := kv.PrependKeysWith(data, "user_")

	if _, ok := result["user_name"]; !ok {
		t.Error("expected 'user_name' key")
	}
}

func TestReplace(t *testing.T) {
	data := map[string]any{"a": 1, "b": 2}
	result := kv.Replace(data, map[string]any{"b": 20, "c": 3})

	if result["b"] != 20 {
		t.Errorf("expected 20, got %v", result["b"])
	}

	if result["c"] != 3 {
		t.Errorf("expected 3, got %v", result["c"])
	}
}

func TestMapTransformTypes(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	data := map[string]User{
		"a": {1, "Alice"},
		"b": {2, "Bob"},
	}
	result := kv.Map(data, func(u User, _ string) string { return u.Name })
	expected := map[string]string{"a": "Alice", "b": "Bob"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
