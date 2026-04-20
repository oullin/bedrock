package arr_test

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/collection/arr"
)

func TestMap(t *testing.T) {
	items := []int{1, 2, 3}
	result := arr.Map(items, func(item int, _ int) int { return item * 2 })
	expected := []int{2, 4, 6}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMapWithKeys(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	users := []User{{1, "Alice"}, {2, "Bob"}}
	result := arr.MapWithKeys(users, func(u User) (int, string) { return u.ID, u.Name })

	if result[1] != "Alice" {
		t.Error("expected Alice at key 1")
	}
}

func TestMapSpread(t *testing.T) {
	items := []int{1, 2, 3}
	result := arr.MapSpread(items, func(item int, _ int) int { return item * 10 })
	expected := []int{10, 20, 30}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestKeyBy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	users := []User{{1, "Alice"}, {2, "Bob"}}
	result := arr.KeyBy(users, func(u User) int { return u.ID })

	if result[1].Name != "Alice" {
		t.Error("expected Alice")
	}
}

func TestPluck(t *testing.T) {
	type User struct {
		Name string
	}

	users := []User{{"Alice"}, {"Bob"}}
	result := arr.Pluck(users, func(u User) string { return u.Name })
	expected := []string{"Alice", "Bob"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestDivide(t *testing.T) {
	items := []string{"a", "b", "c"}
	keys, values := arr.Divide(items)

	if !reflect.DeepEqual(keys, []int{0, 1, 2}) {
		t.Errorf("expected [0 1 2], got %v", keys)
	}

	if !reflect.DeepEqual(values, items) {
		t.Errorf("expected %v, got %v", items, values)
	}
}
