package arr_test

import (
	"testing"

	"github.com/gollin/packages/support/collection/arr"
)

func TestAccessible(t *testing.T) {
	if !arr.Accessible([]int{1, 2, 3}) {
		t.Error("expected true for non-nil slice")
	}

	if arr.Accessible(nil) {
		t.Error("expected false for nil")
	}
}

func TestIsList(t *testing.T) {
	if !arr.IsList([]int{1, 2, 3}) {
		t.Error("expected true")
	}
}

func TestExists(t *testing.T) {
	items := []int{10, 20, 30}

	if !arr.Exists(items, 1) {
		t.Error("expected index 1 to exist")
	}

	if arr.Exists(items, 10) {
		t.Error("expected index 10 to not exist")
	}
}

func TestHas(t *testing.T) {
	items := []int{10, 20, 30}

	if !arr.Has(items, 0, 1, 2) {
		t.Error("expected all indices to exist")
	}

	if arr.Has(items, 0, 5) {
		t.Error("expected to fail for index 5")
	}
}

func TestHasAny(t *testing.T) {
	items := []int{10, 20, 30}

	if !arr.HasAny(items, 0, 99) {
		t.Error("expected at least one valid index")
	}

	if arr.HasAny(items, 5, 10) {
		t.Error("expected no valid indices")
	}
}
