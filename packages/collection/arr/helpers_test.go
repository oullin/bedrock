package arr_test

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/collection/arr"
)

func TestWrap(t *testing.T) {
	result := arr.Wrap(42)

	if !reflect.DeepEqual(result, []int{42}) {
		t.Errorf("expected [42], got %v", result)
	}
}

func TestWrapSlice(t *testing.T) {
	input := []int{1, 2, 3}
	result := arr.WrapSlice(input)

	if !reflect.DeepEqual(result, input) {
		t.Errorf("expected %v, got %v", input, result)
	}
}

func TestJoin(t *testing.T) {
	result := arr.Join([]string{"a", "b", "c"}, ", ")

	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}

	result = arr.Join([]string{"a", "b", "c"}, ", ", " and ")

	if result != "a, b and c" {
		t.Errorf("expected 'a, b and c', got '%s'", result)
	}

	result = arr.Join([]string{}, ", ")

	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}

	result = arr.Join([]string{"solo"}, ", ")

	if result != "solo" {
		t.Errorf("expected 'solo', got '%s'", result)
	}
}
