package arr_test

import (
	"reflect"
	"testing"

	"github.com/gollin/packages/support/collection/arr"
)

func TestSort(t *testing.T) {
	items := []int{3, 1, 4, 1, 5}
	result := arr.Sort(items, func(a, b int) bool { return a < b })
	expected := []int{1, 1, 3, 4, 5}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSortDesc(t *testing.T) {
	items := []int{3, 1, 4, 1, 5}
	result := arr.SortDesc(items, func(a, b int) bool { return a < b })
	expected := []int{5, 4, 3, 1, 1}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSortRecursive(t *testing.T) {
	items := []int{3, 1, 2}
	result := arr.SortRecursive(items, func(a, b int) bool { return a < b })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSortRecursiveDesc(t *testing.T) {
	items := []int{1, 3, 2}
	result := arr.SortRecursiveDesc(items, func(a, b int) bool { return a < b })
	expected := []int{3, 2, 1}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestShuffle(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := arr.Shuffle(items)

	if len(result) != 5 {
		t.Errorf("expected 5 items, got %d", len(result))
	}
}

func TestRandom(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := arr.Random(items, 2)

	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
}

func TestRandomNonPositiveCount(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	for _, count := range []int{0, -1} {
		result := arr.Random(items, count)

		if len(result) != 0 {
			t.Fatalf("expected 0 items for count %d, got %d", count, len(result))
		}
	}
}
