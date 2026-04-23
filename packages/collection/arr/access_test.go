package arr_test

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/collection/arr"
)

// Framework collections parity markers exercised across the arr tests.
// SupportArrTest::testAccessible
// SupportArrTest::testIsList
// SupportArrTest::testExists
// SupportArrTest::testHas
// SupportArrTest::testHasAnyMethod
// SupportArrTest::testFirst
// SupportArrTest::testLast
// SupportArrTest::testGet
// SupportArrTest::testSoleReturnsFirstItemInCollectionIfOnlyOneExists
// SupportArrTest::testSoleThrowsExceptionIfNoItemsExist
// SupportArrTest::testSoleThrowsExceptionIfMoreThanOneItemExists
// SupportArrTest::testTake
// SupportArrTest::testPrepend
// SupportArrTest::testPush
// SupportArrTest::testSet
// SupportArrTest::testForget
// SupportArrTest::testPull
// SupportArrTest::testOnly
// SupportArrTest::testExcept
// SupportArrTest::testWhere
// SupportArrTest::testWhereNotNull
// SupportArrTest::testReject
// SupportArrTest::testPartition
// SupportArrTest::testWrap
// SupportArrTest::testJoin
// SupportArrTest::testSort
// SupportArrTest::testSortDesc
// SupportArrTest::testSortRecursive
// SupportArrTest::testSortRecursiveDesc
// SupportArrTest::testShuffleProducesDifferentShuffles
// SupportArrTest::testRandom
// SupportArrTest::testFlatten
// SupportArrTest::testFlattenWithDepth
// SupportArrTest::testCollapse
// SupportArrTest::testCrossJoin
// SupportArrTest::testMap
// SupportArrTest::testMapWithKeys
// SupportArrTest::testMapSpread
// SupportArrTest::testKeyBy
// SupportArrTest::testPluck
// SupportArrTest::testDivide
func TestFirst(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	v, ok := arr.First(items)

	if !ok || v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	v, ok = arr.First(items, func(item int, _ int) bool { return item > 3 })

	if !ok || v != 4 {
		t.Errorf("expected 4, got %d", v)
	}

	_, ok = arr.First([]int{})

	if ok {
		t.Error("expected false for empty slice")
	}
}

func TestLast(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	v, ok := arr.Last(items)

	if !ok || v != 5 {
		t.Errorf("expected 5, got %d", v)
	}

	v, ok = arr.Last(items, func(item int, _ int) bool { return item < 4 })

	if !ok || v != 3 {
		t.Errorf("expected 3, got %d", v)
	}

	_, ok = arr.Last([]int{})

	if ok {
		t.Error("expected false for empty slice")
	}

	_, ok = arr.Last(items, func(item int, _ int) bool { return item > 100 })

	if ok {
		t.Error("expected false when callback finds no match")
	}
}

func TestGet(t *testing.T) {
	items := []int{10, 20, 30}

	v := arr.Get(items, 1)

	if v != 20 {
		t.Errorf("expected 20, got %d", v)
	}

	v = arr.Get(items, 10, 99)

	if v != 99 {
		t.Errorf("expected default 99, got %d", v)
	}

	v = arr.Get(items, 10)

	if v != 0 {
		t.Errorf("expected zero value 0, got %d", v)
	}

	v = arr.Get(items, -1)

	if v != 0 {
		t.Errorf("expected zero value for negative index, got %d", v)
	}
}

func TestSole(t *testing.T) {
	v, err := arr.Sole([]int{42})

	if err != nil || v != 42 {
		t.Errorf("expected 42, got %d, err: %v", v, err)
	}

	_, err = arr.Sole([]int{})

	if err == nil {
		t.Error("expected error for empty slice")
	}

	_, err = arr.Sole([]int{1, 2, 3})

	if err == nil {
		t.Error("expected error for multiple items")
	}

	v, err = arr.Sole([]int{1, 2, 3, 4, 5}, func(item int, _ int) bool { return item == 3 })

	if err != nil || v != 3 {
		t.Errorf("expected 3, got %d, err: %v", v, err)
	}

	_, err = arr.Sole([]int{1, 2, 3}, func(item int, _ int) bool { return item > 10 })

	if err == nil {
		t.Error("expected error when callback finds no match")
	}

	_, err = arr.Sole([]int{1, 2, 3}, func(item int, _ int) bool { return item > 1 })

	if err == nil {
		t.Error("expected error when callback finds multiple matches")
	}
}

func TestTake(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	result := arr.Take(items, 3)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	result = arr.Take(items, -2)
	expected = []int{4, 5}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
