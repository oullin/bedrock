package arr_test

import (
	"testing"

	"github.com/gollin/packages/support/collection/arr"
)

func TestEvery(t *testing.T) {
	items := []int{2, 4, 6}

	if !arr.Every(items, func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected all even")
	}

	if arr.Every([]int{2, 3, 4}, func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected not all even")
	}
}

func TestSome(t *testing.T) {
	items := []int{1, 2, 3}

	if !arr.Some(items, func(item int, _ int) bool { return item == 2 }) {
		t.Error("expected to find 2")
	}

	if arr.Some(items, func(item int, _ int) bool { return item == 10 }) {
		t.Error("expected not to find 10")
	}
}
