package arr_test

import (
	"reflect"
	"testing"

	"github.com/gollin/packages/support/collection/arr"
)

func TestPrepend(t *testing.T) {
	items := []int{2, 3}
	result := arr.Prepend(items, 1)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestPush(t *testing.T) {
	items := []int{1, 2}
	result := arr.Push(items, 3, 4)
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSet(t *testing.T) {
	items := []int{10, 20, 30}

	result := arr.Set(items, 1, 99)
	expected := []int{10, 99, 30}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	result = arr.Set(items, 10, 99)

	if !reflect.DeepEqual(result, items) {
		t.Errorf("expected unchanged %v, got %v", items, result)
	}

	result = arr.Set(items, -1, 99)

	if !reflect.DeepEqual(result, items) {
		t.Errorf("expected unchanged %v for negative index, got %v", items, result)
	}
}

func TestForget(t *testing.T) {
	items := []int{10, 20, 30}

	result := arr.Forget(items, 1)
	expected := []int{10, 30}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	result = arr.Forget(items, 10)

	if !reflect.DeepEqual(result, items) {
		t.Errorf("expected unchanged %v, got %v", items, result)
	}

	result = arr.Forget(items, -1)

	if !reflect.DeepEqual(result, items) {
		t.Errorf("expected unchanged %v for negative index, got %v", items, result)
	}
}

func TestPull(t *testing.T) {
	items := []int{10, 20, 30}

	v, result := arr.Pull(items, 1)

	if v != 20 {
		t.Errorf("expected pulled value 20, got %d", v)
	}

	expected := []int{10, 30}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	v, result = arr.Pull(items, 10)

	if v != 0 {
		t.Errorf("expected zero value, got %d", v)
	}

	if !reflect.DeepEqual(result, items) {
		t.Errorf("expected unchanged %v, got %v", items, result)
	}

	v, result = arr.Pull(items, -1)

	if v != 0 {
		t.Errorf("expected zero value for negative index, got %d", v)
	}

	if !reflect.DeepEqual(result, items) {
		t.Errorf("expected unchanged %v, got %v", items, result)
	}
}
