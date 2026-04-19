package helpers

import (
	"testing"
)

func TestDotGet(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "Taylor",
		"user": map[string]any{
			"name": "Taylor",
			"age":  30,
			"address": map[string]any{
				"city": "Little Rock",
			},
		},
	}

	val, ok := dotGet(m, "name")

	if !ok || val != "Taylor" {
		t.Errorf("dotGet(name) = (%v, %v), want (Taylor, true)", val, ok)
	}

	val, ok = dotGet(m, "user.name")

	if !ok || val != "Taylor" {
		t.Errorf("dotGet(user.name) = (%v, %v), want (Taylor, true)", val, ok)
	}

	val, ok = dotGet(m, "user.age")

	if !ok || val != 30 {
		t.Errorf("dotGet(user.age) = (%v, %v), want (30, true)", val, ok)
	}

	val, ok = dotGet(m, "user.address.city")

	if !ok || val != "Little Rock" {
		t.Errorf("dotGet(user.address.city) = (%v, %v), want (Little Rock, true)", val, ok)
	}
}

func TestDotGetMissingKey(t *testing.T) {
	t.Parallel()

	m := map[string]any{"user": map[string]any{"name": "Taylor"}}

	val, ok := dotGet(m, "missing")

	if ok || val != nil {
		t.Errorf("dotGet(missing) = (%v, %v), want (nil, false)", val, ok)
	}

	val, ok = dotGet(m, "user.missing")

	if ok || val != nil {
		t.Errorf("dotGet(user.missing) = (%v, %v), want (nil, false)", val, ok)
	}

	val, ok = dotGet(m, "user.name.deep")

	if ok || val != nil {
		t.Errorf("dotGet(user.name.deep) = (%v, %v), want (nil, false)", val, ok)
	}
}

func TestDotGetLiteralDotKey(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"user.name": "literal",
		"user": map[string]any{
			"name": "nested",
		},
	}

	val, ok := dotGet(m, "user.name")

	if !ok || val != "literal" {
		t.Errorf("dotGet(user.name) should prefer literal key, got (%v, %v)", val, ok)
	}
}

func TestDotSet(t *testing.T) {
	t.Parallel()

	m := make(map[string]any)
	dotSet(m, "name", "Taylor")

	if m["name"] != "Taylor" {
		t.Errorf("dotSet(name) = %v, want Taylor", m["name"])
	}

	dotSet(m, "user.name", "Taylor")
	user, ok := m["user"].(map[string]any)

	if !ok || user["name"] != "Taylor" {
		t.Errorf("dotSet(user.name) failed: %v", m)
	}

	dotSet(m, "user.address.city", "Little Rock")
	addr, ok := user["address"].(map[string]any)

	if !ok || addr["city"] != "Little Rock" {
		t.Errorf("dotSet(user.address.city) failed: %v", m)
	}
}

func TestDotSetOverwritesNonMap(t *testing.T) {
	t.Parallel()

	m := map[string]any{"user": "string_value"}
	dotSet(m, "user.name", "Taylor")

	user, ok := m["user"].(map[string]any)

	if !ok || user["name"] != "Taylor" {
		t.Errorf("dotSet should overwrite non-map intermediate: %v", m)
	}
}

func TestDotHas(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "Taylor",
		"user": map[string]any{
			"age": 30,
		},
	}

	if !dotHas(m, "name") {
		t.Error("dotHas(name) should be true")
	}

	if !dotHas(m, "user.age") {
		t.Error("dotHas(user.age) should be true")
	}

	if dotHas(m, "missing") {
		t.Error("dotHas(missing) should be false")
	}

	if dotHas(m, "user.missing") {
		t.Error("dotHas(user.missing) should be false")
	}
}

func TestDotForget(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "Taylor",
		"user": map[string]any{
			"name": "Taylor",
			"age":  30,
		},
	}

	dotForget(m, "name")

	if _, ok := m["name"]; ok {
		t.Error("dotForget(name) should remove top-level key")
	}

	dotForget(m, "user.age")
	user := m["user"].(map[string]any)

	if _, ok := user["age"]; ok {
		t.Error("dotForget(user.age) should remove nested key")
	}

	if user["name"] != "Taylor" {
		t.Error("dotForget should not affect other keys")
	}
}

func TestDotForgetNonExistent(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor"}
	dotForget(m, "missing")
	dotForget(m, "missing.deep")

	if m["name"] != "Taylor" {
		t.Error("dotForget on missing key should not affect existing keys")
	}
}

func TestDotFlatten(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "Taylor",
		"user": map[string]any{
			"name": "Taylor",
			"address": map[string]any{
				"city": "Little Rock",
			},
		},
	}

	result := make(map[string]any)
	dotFlatten(m, "", result)

	expected := map[string]any{
		"name":              "Taylor",
		"user.name":         "Taylor",
		"user.address.city": "Little Rock",
	}

	if len(result) != len(expected) {
		t.Errorf("dotFlatten length = %d, want %d", len(result), len(expected))
	}

	for k, want := range expected {
		if result[k] != want {
			t.Errorf("dotFlatten[%s] = %v, want %v", k, result[k], want)
		}
	}
}

func TestDotFlattenWithPrepend(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "Taylor",
	}

	result := make(map[string]any)
	dotFlatten(m, "prefix", result)

	if result["prefix.name"] != "Taylor" {
		t.Errorf("dotFlatten with prepend: got %v", result)
	}
}
