package helpers

import (
	"testing"
)

// Port of Illuminate\Tests\Support\SupportArrTest::testCollapse
func TestArrCollapse(t *testing.T) {
	t.Parallel()

	result := ArrCollapse([][]int{{1, 2}, {3, 4}, {5}})
	expected := []int{1, 2, 3, 4, 5}

	if len(result) != len(expected) {
		t.Fatalf("ArrCollapse length = %d, want %d", len(result), len(expected))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("ArrCollapse[%d] = %d, want %d", i, result[i], v)
		}
	}
}

func TestArrCollapseEmpty(t *testing.T) {
	t.Parallel()

	result := ArrCollapse([][]int{})
	if len(result) != 0 {
		t.Errorf("ArrCollapse empty = %v, want []", result)
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testFirst
func TestArrFirst(t *testing.T) {
	t.Parallel()

	val, ok := ArrFirst([]int{100, 200, 300})
	if !ok || val != 100 {
		t.Errorf("ArrFirst no predicate = (%v, %v), want (100, true)", val, ok)
	}

	val, ok = ArrFirst([]int{100, 200, 300}, func(v int, _ int) bool {
		return v >= 150
	})
	if !ok || val != 200 {
		t.Errorf("ArrFirst with predicate = (%v, %v), want (200, true)", val, ok)
	}
}

func TestArrFirstEmpty(t *testing.T) {
	t.Parallel()

	val, ok := ArrFirst([]int{})
	if ok || val != 0 {
		t.Errorf("ArrFirst empty = (%v, %v), want (0, false)", val, ok)
	}
}

func TestArrFirstNoMatch(t *testing.T) {
	t.Parallel()

	val, ok := ArrFirst([]int{1, 2, 3}, func(v int, _ int) bool {
		return v > 10
	})

	if ok || val != 0 {
		t.Errorf("ArrFirst no match = (%v, %v), want (0, false)", val, ok)
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testLast
func TestArrLast(t *testing.T) {
	t.Parallel()

	val, ok := ArrLast([]int{100, 200, 300})
	if !ok || val != 300 {
		t.Errorf("ArrLast no predicate = (%v, %v), want (300, true)", val, ok)
	}

	val, ok = ArrLast([]int{100, 200, 300}, func(v int, _ int) bool {
		return v < 250
	})
	if !ok || val != 200 {
		t.Errorf("ArrLast with predicate = (%v, %v), want (200, true)", val, ok)
	}
}

func TestArrLastEmpty(t *testing.T) {
	t.Parallel()

	val, ok := ArrLast([]int{})
	if ok || val != 0 {
		t.Errorf("ArrLast empty = (%v, %v), want (0, false)", val, ok)
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testFlatten
func TestArrFlatten(t *testing.T) {
	t.Parallel()

	items := []any{1, []any{2, 3}, []any{4, []any{5, 6}}}
	result := ArrFlatten(items)

	expected := []any{1, 2, 3, 4, 5, 6}
	if len(result) != len(expected) {
		t.Fatalf("ArrFlatten length = %d, want %d", len(result), len(expected))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("ArrFlatten[%d] = %v, want %v", i, result[i], v)
		}
	}
}

func TestArrFlattenWithDepth(t *testing.T) {
	t.Parallel()

	items := []any{1, []any{2, []any{3, 4}}}

	result := ArrFlatten(items, 1)
	if len(result) != 3 {
		t.Fatalf("ArrFlatten depth=1 length = %d, want 3", len(result))
	}

	if result[0] != 1 || result[1] != 2 {
		t.Errorf("ArrFlatten depth=1 first elements = %v", result[:2])
	}

	nested, ok := result[2].([]any)
	if !ok || len(nested) != 2 {
		t.Errorf("ArrFlatten depth=1 should preserve deeper nesting: %v", result[2])
	}
}

func TestArrFlattenZeroDepth(t *testing.T) {
	t.Parallel()

	items := []any{1, []any{2, 3}}
	result := ArrFlatten(items, 0)

	if len(result) != 2 {
		t.Errorf("ArrFlatten depth=0 should not flatten: %v", result)
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testPrepend
func TestArrPrepend(t *testing.T) {
	t.Parallel()

	result := ArrPrepend([]int{2, 3, 4}, 1)
	expected := []int{1, 2, 3, 4}

	if len(result) != len(expected) {
		t.Fatalf("ArrPrepend length = %d, want %d", len(result), len(expected))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("ArrPrepend[%d] = %d, want %d", i, result[i], v)
		}
	}
}

func TestArrPrependEmpty(t *testing.T) {
	t.Parallel()

	result := ArrPrepend([]string{}, "first")
	if len(result) != 1 || result[0] != "first" {
		t.Errorf("ArrPrepend to empty = %v", result)
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testRandom
func TestArrRandom(t *testing.T) {
	t.Parallel()

	items := []int{1, 2, 3, 4, 5}

	result, err := ArrRandom(items)
	if err != nil {
		t.Fatalf("ArrRandom error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("ArrRandom should return 1 element, got %d", len(result))
	}

	found := false
	for _, v := range items {
		if v == result[0] {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ArrRandom returned %v which is not in original slice", result[0])
	}
}

func TestArrRandomMultiple(t *testing.T) {
	t.Parallel()

	items := []int{1, 2, 3, 4, 5}
	result, err := ArrRandom(items, 3)

	if err != nil {
		t.Fatalf("ArrRandom(3) error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("ArrRandom(3) length = %d, want 3", len(result))
	}
}

func TestArrRandomEmpty(t *testing.T) {
	t.Parallel()

	_, err := ArrRandom([]int{})
	if err == nil {
		t.Error("ArrRandom on empty slice should return error")
	}
}

func TestArrRandomExceedsLength(t *testing.T) {
	t.Parallel()

	_, err := ArrRandom([]int{1, 2}, 5)
	if err == nil {
		t.Error("ArrRandom with count > length should return error")
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testSort
func TestArrSort(t *testing.T) {
	t.Parallel()

	result := ArrSort([]int{3, 1, 2})
	expected := []int{1, 2, 3}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("ArrSort[%d] = %d, want %d", i, result[i], v)
		}
	}
}

func TestArrSortStrings(t *testing.T) {
	t.Parallel()

	result := ArrSort([]string{"banana", "apple", "cherry"})
	expected := []string{"apple", "banana", "cherry"}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("ArrSort strings[%d] = %q, want %q", i, result[i], v)
		}
	}
}

func TestArrSortDoesNotMutate(t *testing.T) {
	t.Parallel()

	original := []int{3, 1, 2}
	ArrSort(original)

	if original[0] != 3 || original[1] != 1 || original[2] != 2 {
		t.Errorf("ArrSort mutated original: %v", original)
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testSortWithCallback
func TestArrSortFunc(t *testing.T) {
	t.Parallel()

	type item struct {
		Name string
		Age  int
	}

	items := []item{
		{"Taylor", 30},
		{"Abigail", 25},
		{"Dayle", 35},
	}

	result := ArrSortFunc(items, func(i item) int { return i.Age })

	if result[0].Name != "Abigail" || result[1].Name != "Taylor" || result[2].Name != "Dayle" {
		t.Errorf("ArrSortFunc = %v", result)
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testWhere
func TestArrWhere(t *testing.T) {
	t.Parallel()

	result := ArrWhere([]int{1, 2, 3, 4, 5}, func(v int, _ int) bool {
		return v%2 == 0
	})

	if len(result) != 2 || result[0] != 2 || result[1] != 4 {
		t.Errorf("ArrWhere = %v, want [2, 4]", result)
	}
}

func TestArrWhereNoMatch(t *testing.T) {
	t.Parallel()

	result := ArrWhere([]int{1, 3, 5}, func(v int, _ int) bool {
		return v%2 == 0
	})

	if len(result) != 0 {
		t.Errorf("ArrWhere no match = %v, want []", result)
	}
}

func TestArrWhereWithIndex(t *testing.T) {
	t.Parallel()

	result := ArrWhere([]string{"a", "b", "c"}, func(_ string, i int) bool {
		return i > 0
	})

	if len(result) != 2 || result[0] != "b" || result[1] != "c" {
		t.Errorf("ArrWhere with index = %v", result)
	}
}

// Port of Illuminate\Tests\Support\SupportArrTest::testWrap
func TestArrWrap(t *testing.T) {
	t.Parallel()

	result := ArrWrap[int](42)
	if len(result) != 1 || result[0] != 42 {
		t.Errorf("ArrWrap(42) = %v, want [42]", result)
	}

	result = ArrWrap[int]([]int{1, 2, 3})
	if len(result) != 3 {
		t.Errorf("ArrWrap(slice) = %v, want [1,2,3]", result)
	}

	result = ArrWrap[int](nil)
	if len(result) != 0 {
		t.Errorf("ArrWrap(nil) = %v, want []", result)
	}
}

func TestArrWrapStrings(t *testing.T) {
	t.Parallel()

	result := ArrWrap[string]("hello")
	if len(result) != 1 || result[0] != "hello" {
		t.Errorf("ArrWrap string = %v, want [hello]", result)
	}

	result = ArrWrap[string]([]string{"a", "b"})
	if len(result) != 2 {
		t.Errorf("ArrWrap string slice = %v, want [a, b]", result)
	}
}
