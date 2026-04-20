package arr_test

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/collection/arr"
)

func TestFlatten(t *testing.T) {
	items := [][]int{{1, 2}, {3, 4}, {5}}
	result := arr.Flatten(items)
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestCollapse(t *testing.T) {
	items := [][]int{{1, 2}, {3, 4}}
	result := arr.Collapse(items)

	if len(result) != 4 {
		t.Errorf("expected 4 items, got %d", len(result))
	}
}

func TestCrossJoin(t *testing.T) {
	result := arr.CrossJoin([]int{1, 2}, []int{10, 20})

	if len(result) != 4 {
		t.Errorf("expected 4 combinations, got %d", len(result))
	}
}
