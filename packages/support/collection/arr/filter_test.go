package arr_test

import (
	"reflect"
	"testing"

	"github.com/gollin/packages/support/collection/arr"
)

func TestOnly(t *testing.T) {
	items := []int{10, 20, 30, 40, 50}
	result := arr.Only(items, []int{1, 3})
	expected := []int{20, 40}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestExcept(t *testing.T) {
	items := []int{10, 20, 30, 40, 50}
	result := arr.Except(items, []int{1, 3})
	expected := []int{10, 30, 50}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestWhere(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := arr.Where(items, func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{2, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestWhereNotNull(t *testing.T) {
	items := []string{"a", "", "b", "", "c"}
	result := arr.WhereNotNull(items)
	expected := []string{"a", "b", "c"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestReject(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	result := arr.Reject(items, func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{1, 3, 5}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestPartition(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	pass, fail := arr.Partition(items, func(item int, _ int) bool { return item > 3 })

	if !reflect.DeepEqual(pass, []int{4, 5}) {
		t.Errorf("expected [4 5], got %v", pass)
	}

	if !reflect.DeepEqual(fail, []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", fail)
	}
}
