package collectible

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gollin/packages/support/collection/support"
)

func TestNew(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})

	if m.Count() != 3 {
		t.Errorf("expected 3, got %d", m.Count())
	}
}

func TestFromPairs(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	v, ok := m.Get("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestMapGet(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	v, ok := m.Get("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	_, ok = m.Get("z")

	if ok {
		t.Error("expected not found")
	}

	v, _ = m.Get("z", 99)

	if v != 99 {
		t.Errorf("expected default 99, got %d", v)
	}
}

func TestMapGetOrPut(t *testing.T) {
	m := New(map[string]int{"a": 1})

	v := m.GetOrPut("a", 99)

	if v != 1 {
		t.Errorf("expected existing 1, got %d", v)
	}

	v = m.GetOrPut("b", 99)

	if v != 99 {
		t.Errorf("expected default 99, got %d", v)
	}

	if m.Count() != 2 {
		t.Errorf("expected 2 items, got %d", m.Count())
	}
}

func TestMapHas(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.Has("a") {
		t.Error("expected to have 'a'")
	}

	if m.Has("z") {
		t.Error("expected not to have 'z'")
	}
}

func TestMapHasAny(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.HasAny("z", "a") {
		t.Error("expected HasAny to return true")
	}

	if m.HasAny("x", "y", "z") {
		t.Error("expected HasAny to return false")
	}
}

func TestMapPut(t *testing.T) {
	m := New(map[string]int{})
	m.Put("a", 1)
	v, ok := m.Get("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestMapPull(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	v, ok := m.Pull("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	if m.Has("a") {
		t.Error("expected 'a' to be removed")
	}
}

func TestMapForget(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	m.Forget("a", "c")

	if m.Count() != 1 {
		t.Errorf("expected 1, got %d", m.Count())
	}
}

func TestMapOnly(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.Only("a", "c")

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}

	if !result.Has("a") || !result.Has("c") {
		t.Error("expected to have 'a' and 'c'")
	}
}

func TestMapExcept(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.Except("b")

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}

	if result.Has("b") {
		t.Error("expected 'b' to be excluded")
	}
}

func TestMapKeys(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	keys := m.Keys()

	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestMapValues(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	values := m.Values()

	if len(values) != 2 {
		t.Errorf("expected 2 values, got %d", len(values))
	}
}

func TestMapContains(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.Contains(func(v int, k string) bool { return v == 2 }) {
		t.Error("expected to contain value 2")
	}

	if m.Contains(func(v int, k string) bool { return v == 99 }) {
		t.Error("expected not to contain 99")
	}
}

func TestMapFirst(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	v, ok := m.First()

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestMapLast(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	v, ok := m.Last()

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
}

func TestMapFilter(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	filtered := m.Filter(func(v int, k string) bool { return v > 1 })

	if filtered.Count() != 2 {
		t.Errorf("expected 2, got %d", filtered.Count())
	}
}

func TestMapReject(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	rejected := m.Reject(func(v int, k string) bool { return v > 1 })

	if rejected.Count() != 1 {
		t.Errorf("expected 1, got %d", rejected.Count())
	}
}

func TestMapValues_Transform(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	result := MapValues(m, func(v int, k string) int { return v * 10 })
	v, _ := result.Get("a")

	if v != 10 {
		t.Errorf("expected 10, got %d", v)
	}
}

func TestMapEvery(t *testing.T) {
	m := New(map[string]int{"a": 2, "b": 4})

	if !m.Every(func(v int, k string) bool { return v%2 == 0 }) {
		t.Error("expected all even")
	}
}

func TestMapPartition(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	pass, fail := m.Partition(func(v int, k string) bool { return v > 1 })

	if pass.Count() != 2 {
		t.Errorf("expected 2 passing, got %d", pass.Count())
	}

	if fail.Count() != 1 {
		t.Errorf("expected 1 failing, got %d", fail.Count())
	}
}

func TestMapMerge(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	result := m.Merge(map[string]int{"b": 20, "c": 3})
	v, _ := result.Get("b")

	if v != 20 {
		t.Errorf("expected 20 (overwritten), got %d", v)
	}

	if result.Count() != 3 {
		t.Errorf("expected 3, got %d", result.Count())
	}
}

func TestMapUnion(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	result := m.Union(map[string]int{"b": 20, "c": 3})
	v, _ := result.Get("b")

	if v != 2 {
		t.Errorf("expected 2 (not overwritten), got %d", v)
	}

	if result.Count() != 3 {
		t.Errorf("expected 3, got %d", result.Count())
	}
}

func TestMapDiffKeys(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.DiffKeys(map[string]int{"b": 99})

	if result.Has("b") {
		t.Error("expected 'b' to be excluded")
	}

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestMapIntersectByKeys(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.IntersectByKeys(map[string]int{"a": 99, "c": 99})

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestFlip(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	flipped := Flip(m)
	v, ok := flipped.Get(1)

	if !ok || v != "a" {
		t.Errorf("expected 'a', got '%s'", v)
	}
}

func TestSortKeys(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "c", Value: 3},
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	sorted := SortKeys(m)
	keys := sorted.Keys()
	expected := []string{"a", "b", "c"}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestSortKeysDesc(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "c", Value: 3},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	sorted := SortKeysDesc(m)
	keys := sorted.Keys()
	expected := []string{"c", "b", "a"}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestMapImplode(t *testing.T) {
	m := FromPairs(
		support.Pair[string, string]{Key: "a", Value: "hello"},
		support.Pair[string, string]{Key: "b", Value: "world"},
	)
	result := m.Implode(", ")

	if result != "hello, world" {
		t.Errorf("expected 'hello, world', got '%s'", result)
	}
}

func TestMapToJSON(t *testing.T) {
	m := New(map[string]int{"a": 1})
	b, err := m.ToJSON()

	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `{"a":1}` {
		t.Errorf("expected {\"a\":1}, got %s", string(b))
	}
}

func TestMapCopy(t *testing.T) {
	m := New(map[string]int{"a": 1})
	m2 := m.Copy()
	m.Put("b", 2)

	if m2.Has("b") {
		t.Error("copy should not be affected")
	}
}

func TestMapToPairs(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	pairs := m.ToPairs()

	if len(pairs) != 2 {
		t.Errorf("expected 2 pairs, got %d", len(pairs))
	}
}

func TestMapWhen(t *testing.T) {
	m := New(map[string]int{"a": 1})
	result := m.When(true, func(mc *Collection[string, int]) *Collection[string, int] {
		mc.Put("b", 2)

		return mc
	})

	if !result.Has("b") {
		t.Error("expected 'b' to be added")
	}
}

func TestMapSearch(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	key, ok := m.Search(func(v int, k string) bool { return v == 2 })

	if !ok || key != "b" {
		t.Errorf("expected 'b', got '%s'", key)
	}
}

func TestMapIsEmpty(t *testing.T) {
	m := New(map[string]int{})

	if !m.IsEmpty() {
		t.Error("expected empty")
	}

	if m.IsNotEmpty() {
		t.Error("expected empty")
	}
}

func TestMapContainsOneItem(t *testing.T) {
	m := New(map[string]int{"a": 1})

	if !m.ContainsOneItem() {
		t.Error("expected true")
	}
}

func TestMapContainsManyItems(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.ContainsManyItems() {
		t.Error("expected true")
	}
}

func TestMapEach(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	sum := 0
	m.Each(func(v int, k string) bool {
		sum += v

		return true
	})

	if sum != 3 {
		t.Errorf("expected 3, got %d", sum)
	}
}

func TestMapDoesntContain(t *testing.T) {
	m := New(map[string]int{"a": 1})

	if !m.DoesntContain(func(v int, k string) bool { return v == 99 }) {
		t.Error("expected true")
	}
}

func TestMapSortKeysUsing(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "banana", Value: 1},
		support.Pair[string, int]{Key: "apple", Value: 2},
		support.Pair[string, int]{Key: "cherry", Value: 3},
	)
	sorted := m.SortKeysUsing(func(a, b string) bool { return a < b })
	keys := sorted.Keys()
	expected := []string{"apple", "banana", "cherry"}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestMapReplace(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})
	result := m.Replace(map[string]int{"b": 20, "c": 3})
	v, _ := result.Get("b")

	if v != 20 {
		t.Errorf("expected 20, got %d", v)
	}
}

func TestMapSome(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2})

	if !m.Some(func(v int, k string) bool { return v == 2 }) {
		t.Error("expected Some to return true")
	}

	if m.Some(func(v int, k string) bool { return v == 99 }) {
		t.Error("expected Some to return false")
	}
}

func TestMapHasSole(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})

	if !m.HasSole(func(v int, k string) bool { return v == 2 }) {
		t.Error("expected true for single match")
	}

	if m.HasSole(func(v int, k string) bool { return v > 1 }) {
		t.Error("expected false for multiple matches")
	}

	if m.HasSole(func(v int, k string) bool { return v > 10 }) {
		t.Error("expected false for no matches")
	}
}

func TestMapFirstWithPredicate(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
		support.Pair[string, int]{Key: "c", Value: 3},
	)

	v, ok := m.First(func(v int, k string) bool { return v > 1 })

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}

	_, ok = m.First(func(v int, k string) bool { return v > 10 })

	if ok {
		t.Error("expected false when no match")
	}

	empty := New(map[string]int{})

	_, ok = empty.First()

	if ok {
		t.Error("expected false for empty collection")
	}
}

func TestMapLastWithPredicate(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
		support.Pair[string, int]{Key: "c", Value: 3},
	)

	v, ok := m.Last(func(v int, k string) bool { return v < 3 })

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}

	_, ok = m.Last(func(v int, k string) bool { return v > 10 })

	if ok {
		t.Error("expected false when no match")
	}

	empty := New(map[string]int{})

	_, ok = empty.Last()

	if ok {
		t.Error("expected false for empty collection")
	}
}

func TestMapSearchNotFound(t *testing.T) {
	m := New(map[string]int{"a": 1})

	_, ok := m.Search(func(v int, k string) bool { return v == 99 })

	if ok {
		t.Error("expected false when no match")
	}
}

func TestMapEachEarlyReturn(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
		support.Pair[string, int]{Key: "c", Value: 3},
	)
	count := 0
	m.Each(func(v int, k string) bool {
		count++

		return count < 2
	})

	if count != 2 {
		t.Errorf("expected 2 iterations, got %d", count)
	}
}

func TestMapEveryEmpty(t *testing.T) {
	m := New(map[string]int{})

	if !m.Every(func(v int, k string) bool { return v > 0 }) {
		t.Error("expected true for empty collection")
	}
}

func TestMergeRecursive(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
	)
	result := MergeRecursive(m, map[string]int{"b": 2})

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestMapDiffKeysUsing(t *testing.T) {
	m := New(map[string]int{"a": 1, "b": 2, "c": 3})
	result := m.DiffKeysUsing(map[string]int{"A": 99}, func(k1, k2 string) bool {
		return strings.EqualFold(k1, k2)
	})

	if result.Has("a") {
		t.Error("expected 'a' to be excluded via case-insensitive match")
	}

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestMapDiffAssoc(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
		support.Pair[string, int]{Key: "c", Value: 3},
	)
	result := DiffAssoc(m, map[string]int{"a": 1, "b": 99})

	if result.Has("a") {
		t.Error("expected 'a' excluded (same key and value)")
	}

	if !result.Has("b") {
		t.Error("expected 'b' included (value differs)")
	}

	if !result.Has("c") {
		t.Error("expected 'c' included (key missing from items)")
	}
}

func TestMapIntersectAssoc(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
		support.Pair[string, int]{Key: "c", Value: 3},
	)
	result := IntersectAssoc(m, map[string]int{"a": 1, "b": 99})

	if !result.Has("a") {
		t.Error("expected 'a' included (key and value match)")
	}

	if result.Has("b") {
		t.Error("expected 'b' excluded (value differs)")
	}

	if result.Has("c") {
		t.Error("expected 'c' excluded (key missing)")
	}

	if result.Count() != 1 {
		t.Errorf("expected 1, got %d", result.Count())
	}
}

func TestMapJoin(t *testing.T) {
	m := FromPairs(
		support.Pair[string, string]{Key: "a", Value: "hello"},
		support.Pair[string, string]{Key: "b", Value: "beautiful"},
		support.Pair[string, string]{Key: "c", Value: "world"},
	)
	result := m.Join(", ", " and ")

	if result != "hello, beautiful and world" {
		t.Errorf("expected 'hello, beautiful and world', got '%s'", result)
	}

	result = m.Join(", ")

	if result != "hello, beautiful, world" {
		t.Errorf("expected 'hello, beautiful, world', got '%s'", result)
	}
}

func TestMapTap(t *testing.T) {
	m := New(map[string]int{"a": 1})
	called := false
	result := m.Tap(func(mc *Collection[string, int]) {
		called = true
	})

	if !called {
		t.Error("expected callback to be called")
	}

	if result != m {
		t.Error("expected Tap to return self")
	}
}

func TestMapUnless(t *testing.T) {
	m := New(map[string]int{"a": 1})
	result := m.Unless(false, func(mc *Collection[string, int]) *Collection[string, int] {
		mc.Put("b", 2)

		return mc
	})

	if !result.Has("b") {
		t.Error("expected callback to be applied when condition is false")
	}

	m2 := New(map[string]int{"a": 1})
	result2 := m2.Unless(true, func(mc *Collection[string, int]) *Collection[string, int] {
		mc.Put("b", 2)

		return mc
	})

	if result2.Has("b") {
		t.Error("expected callback not applied when condition is true")
	}
}

func TestMapWhenFalseWithDefault(t *testing.T) {
	m := New(map[string]int{"a": 1})
	result := m.When(false, func(mc *Collection[string, int]) *Collection[string, int] {
		mc.Put("b", 2)

		return mc
	}, func(mc *Collection[string, int]) *Collection[string, int] {
		mc.Put("c", 3)

		return mc
	})

	if result.Has("b") {
		t.Error("expected primary callback not called")
	}

	if !result.Has("c") {
		t.Error("expected default callback applied")
	}
}

func TestMapWhenFalseWithoutDefault(t *testing.T) {
	m := New(map[string]int{"a": 1})
	result := m.When(false, func(mc *Collection[string, int]) *Collection[string, int] {
		mc.Put("b", 2)

		return mc
	})

	if result.Has("b") {
		t.Error("expected callback not applied")
	}

	if result.Count() != 1 {
		t.Errorf("expected 1, got %d", result.Count())
	}
}

func TestMapToPrettyJSON(t *testing.T) {
	m := New(map[string]int{"a": 1})
	b, err := m.ToPrettyJSON()

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(b), "\n") {
		t.Error("expected pretty JSON with newlines")
	}
}

func TestMapString(t *testing.T) {
	m := New(map[string]int{"a": 1})
	s := m.String()

	if s != `{"a":1}` {
		t.Errorf("expected {\"a\":1}, got %s", s)
	}
}

func TestMapMarshalJSON(t *testing.T) {
	m := New(map[string]int{"a": 1})
	b, err := m.MarshalJSON()

	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `{"a":1}` {
		t.Errorf("expected {\"a\":1}, got %s", string(b))
	}
}

func TestMapUnmarshalJSON(t *testing.T) {
	m := New(map[string]int{})
	err := m.UnmarshalJSON([]byte(`{"a":1,"b":2}`))

	if err != nil {
		t.Fatal(err)
	}

	v, ok := m.Get("a")

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestMapDump(t *testing.T) {
	m := New(map[string]int{"a": 1})
	result := m.Dump()

	if result != m {
		t.Error("expected Dump to return self")
	}
}

func TestMapIter(t *testing.T) {
	m := FromPairs(
		support.Pair[string, int]{Key: "a", Value: 1},
		support.Pair[string, int]{Key: "b", Value: 2},
	)
	sum := 0

	for _, v := range m.Iter() {
		sum += v
	}

	if sum != 3 {
		t.Errorf("expected 3, got %d", sum)
	}
}

func TestMapNewNil(t *testing.T) {
	m := New[string, int](nil)

	if !m.IsEmpty() {
		t.Error("expected empty collection from nil map")
	}
}
