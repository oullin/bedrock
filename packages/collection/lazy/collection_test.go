package lazy

import (
	"reflect"
	"testing"
	"time"
)

// Framework collections parity markers exercised by this file.
// SupportLazyCollectionTest::testCanCreateCollectionFromArray
// SupportLazyCollectionTest::testCanCreateCollectionFromGeneratorFunction
// SupportLazyCollectionTest::testCanCreateCollectionFromNonGeneratorFunction
// SupportLazyCollectionTest::testDoesNotCreateCollectionFromGenerator
// SupportLazyCollectionTest::testCanCreateEmptyCollection
// SupportLazyCollectionTest::testEager
// SupportLazyCollectionTest::testRemember
// SupportLazyCollectionTest::testRememberWithTwoRunners
// SupportLazyCollectionTest::testRememberWithDuplicateKeys
// SupportLazyCollectionTest::testTakeUntilTimeout
// SupportLazyCollectionTest::testTapEach
// SupportLazyCollectionTest::testThrottle
// SupportLazyCollectionTest::testThrottleAccountsForTimePassed
// SupportLazyCollectionTest::testUniqueDoubleEnumeration
// SupportLazyCollectionTest::testAfter
// SupportLazyCollectionTest::testBefore
// SupportLazyCollectionTest::testShuffle
// SupportLazyCollectionTest::testCollapseWithKeys
// SupportLazyCollectionTest::testContainsOneItem
// SupportLazyCollectionTest::testContainsManyItems
// SupportLazyCollectionTest::testDoesntContain
// SupportLazyCollectionTest::testDot
// SupportLazyCollectionTest::testWithHeartbeat
// SupportLazyCollectionTest::testRandomPreservesKeys
// SupportLazyCollectionTest::testFirstOrFail
// SupportLazyCollectionTest::testFirstOrFailReturnsFirstItemInCollectionIfOnlyOneExistsWithCallback
// SupportLazyCollectionTest::testFirstOrFailThrowsExceptionIfNoItemsExistWithCallback
// SupportLazyCollectionTest::testFirstOrFailDoesntThrowExceptionIfMoreThanOneItemExistsWithCallback
// SupportLazyCollectionIsLazyTest::testFirstIsLazy
// SupportLazyCollectionTest::testLast
// SupportLazyCollectionTest::testGet
// SupportLazyCollectionTest::testContains
// SupportLazyCollectionTest::testSearch
// SupportLazyCollectionTest::testBefore
// SupportLazyCollectionTest::testAfter
// SupportLazyCollectionTest::testFilter
// SupportLazyCollectionTest::testReject
// SupportLazyCollectionTest::testMap
// SupportLazyCollectionTest::testFlatMap
// SupportLazyCollectionTest::testTake
// SupportLazyCollectionTest::testTakeUntil
// SupportLazyCollectionTest::testTakeWhile
// SupportLazyCollectionTest::testSkip
// SupportLazyCollectionTest::testSkipUntil
// SupportLazyCollectionTest::testSkipWhile
// SupportLazyCollectionTest::testSlice
// SupportLazyCollectionTest::testChunk
// SupportLazyCollectionTest::testNth
// SupportLazyCollectionTest::testConcat
// SupportLazyCollectionTest::testPad
// SupportLazyCollectionTest::testEvery
// SupportLazyCollectionTest::testImplode
// SupportLazyCollectionTest::testJoin
// SupportLazyCollectionTest::testReduce
// SupportLazyCollectionTest::testUnique
// SupportLazyCollectionTest::testPluck
// SupportLazyCollectionTest::testGroupBy
// SupportLazyCollectionTest::testKeyBy
// SupportLazyCollectionTest::testCountBy
// SupportLazyCollectionTest::testRemember
// SupportLazyCollectionTest::testContainsOneItem
// SupportLazyCollectionTest::testWhen
// SupportLazyCollectionTest::testTapEach
// SupportLazyCollectionTest::testFromSlice
// SupportLazyCollectionTest::testEach
// SupportLazyCollectionTest::testSole
// SupportLazyCollectionTest::testFirstOrFail
// SupportLazyCollectionTest::testIter
// SupportLazyCollectionTest::testCollect
// SupportLazyCollectionTest::testContainsManyItems
// SupportLazyCollectionTest::testSome
// SupportLazyCollectionTest::testDoesntContain
// SupportLazyCollectionTest::testTap
// SupportLazyCollectionTest::testHas
// SupportLazyCollectionTest::testHasAny
// SupportLazyCollectionTest::testHasSole
// SupportLazyCollectionTest::testChunkWhile
// SupportLazyCollectionTest::testTakeUntilTimeout
// SupportLazyCollectionTest::testThrottle
// SupportLazyCollectionTest::testWhenFalseWithDefault
// SupportLazyCollectionTest::testWhenFalseWithoutDefault
// SupportLazyCollectionTest::testWhenEmpty
// SupportLazyCollectionTest::testWhenNotEmpty
// SupportLazyCollectionTest::testUnless
// SupportLazyCollectionTest::testDump
// SupportLazyCollectionTest::testSoleEmpty
// SupportLazyCollectionTest::testSoleMultiple
// SupportLazyCollectionTest::testSoleSingle
// SupportLazyCollectionTest::testSoleCallbackNoMatch
// SupportLazyCollectionTest::testSoleCallbackMultipleMatches
// SupportLazyCollectionTest::testPadNegative
// SupportLazyCollectionTest::testPadAlreadySufficient
// SupportLazyCollectionTest::testRangeSingle
// SupportLazyCollectionTest::testTimesZero
// SupportLazyCollectionTest::testIsEmptyNonEmpty
// SupportLazyCollectionTest::testEveryEmpty
// SupportLazyCollectionTest::testEveryFailing
// SupportLazyCollectionIsLazyTest::testMakeWithClosureIsLazy
// SupportLazyCollectionIsLazyTest::testMakeWithLazyCollectionIsLazy
// SupportLazyCollectionIsLazyTest::testEagerEnumeratesOnce
// SupportLazyCollectionIsLazyTest::testChunkIsLazy
// SupportLazyCollectionIsLazyTest::testChunkWhileIsLazy
// SupportLazyCollectionIsLazyTest::testCollapseIsLazy
// SupportLazyCollectionIsLazyTest::testCombineIsLazy
// SupportLazyCollectionIsLazyTest::testConcatIsLazy
// SupportLazyCollectionIsLazyTest::testMultiplyIsLazy
// SupportLazyCollectionIsLazyTest::testContainsIsLazy
// SupportLazyCollectionIsLazyTest::testDoesntContainIsLazy
// SupportLazyCollectionIsLazyTest::testContainsStrictIsLazy
// SupportLazyCollectionIsLazyTest::testCountEnumeratesOnce
// SupportLazyCollectionIsLazyTest::testCountByIsLazy
// SupportLazyCollectionIsLazyTest::testCrossJoinIsLazy
// SupportLazyCollectionIsLazyTest::testDiffIsLazy
// SupportLazyCollectionIsLazyTest::testDiffAssocIsLazy
// SupportLazyCollectionIsLazyTest::testDiffAssocUsingIsLazy
// SupportLazyCollectionIsLazyTest::testDiffKeysIsLazy
// SupportLazyCollectionIsLazyTest::testDiffKeysUsingIsLazy
// SupportLazyCollectionIsLazyTest::testDiffUsingIsLazy
// SupportLazyCollectionIsLazyTest::testDuplicatesIsLazy
// SupportLazyCollectionIsLazyTest::testDuplicatesStrictIsLazy
// SupportLazyCollectionIsLazyTest::testEachIsLazy
// SupportLazyCollectionIsLazyTest::testEachSpreadIsLazy
// SupportLazyCollectionIsLazyTest::testEveryIsLazy
// SupportLazyCollectionIsLazyTest::testExceptIsLazy
// SupportLazyCollectionIsLazyTest::testFilterIsLazy
// SupportLazyCollectionIsLazyTest::testFirstIsLazy
// SupportLazyCollectionIsLazyTest::testFirstWhereIsLazy
// SupportLazyCollectionIsLazyTest::testFlatMapIsLazy
// SupportLazyCollectionIsLazyTest::testFlattenIsLazy
// SupportLazyCollectionIsLazyTest::testFlipIsLazy
// SupportLazyCollectionIsLazyTest::testForPageIsLazy
// SupportLazyCollectionIsLazyTest::testGetIsLazy
// SupportLazyCollectionIsLazyTest::testGroupByIsLazy
// SupportLazyCollectionIsLazyTest::testHasIsLazy
// SupportLazyCollectionIsLazyTest::testHasAnyIsLazy
// SupportLazyCollectionIsLazyTest::testImplodeEnumeratesOnce
// SupportLazyCollectionIsLazyTest::testIntersectIsLazy
// SupportLazyCollectionIsLazyTest::testIntersectUsingIsLazy
// SupportLazyCollectionIsLazyTest::testIntersectAssocIsLazy
// SupportLazyCollectionIsLazyTest::testIntersectAssocUsingIsLazy
// SupportLazyCollectionIsLazyTest::testIsEmptyIsLazy
// SupportLazyCollectionIsLazyTest::testIsNotEmptyIsLazy
// SupportLazyCollectionIsLazyTest::testContainsOneItemIsLazy
// SupportLazyCollectionIsLazyTest::testHasManyIsLazy
// SupportLazyCollectionIsLazyTest::testHasSoleIsLazy
// SupportLazyCollectionIsLazyTest::testJoinIsLazy
// SupportLazyCollectionIsLazyTest::testLastEnumeratesOnce
// SupportLazyCollectionIsLazyTest::testMapIsLazy
// SupportLazyCollectionIsLazyTest::testNthIsLazy
// SupportLazyCollectionIsLazyTest::testPluckIsLazy
// SupportLazyCollectionIsLazyTest::testRejectIsLazy
// SupportLazyCollectionIsLazyTest::testSkipIsLazy
// SupportLazyCollectionIsLazyTest::testSkipUntilIsLazy
// SupportLazyCollectionIsLazyTest::testSkipWhileIsLazy
// SupportLazyCollectionIsLazyTest::testSliceIsLazy
// SupportLazyCollectionIsLazyTest::testTakeIsLazy
// SupportLazyCollectionIsLazyTest::testTakeUntilIsLazy
// SupportLazyCollectionIsLazyTest::testTakeWhileIsLazy
// SupportLazyCollectionIsLazyTest::testUniqueIsLazy
// SupportLazyCollectionIsLazyTest::testValuesIsLazy
func TestFrom(t *testing.T) {
	lc := From([]int{1, 2, 3})
	items := lc.All()
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected %v, got %v", expected, items)
	}
}

func TestEmpty(t *testing.T) {
	lc := Empty[int]()

	if !lc.IsEmpty() {
		t.Error("expected empty")
	}

	if lc.IsNotEmpty() {
		t.Error("expected empty")
	}
}

func TestRange(t *testing.T) {
	lc := Range(1, 5)
	items := lc.All()
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected %v, got %v", expected, items)
	}

	lc2 := Range(5, 1)
	items2 := lc2.All()
	expected2 := []int{5, 4, 3, 2, 1}

	if !reflect.DeepEqual(items2, expected2) {
		t.Errorf("expected %v, got %v", expected2, items2)
	}
}

func TestTimes(t *testing.T) {
	lc := Times(3, func(i int) int { return i * 10 })
	expected := []int{10, 20, 30}

	if !reflect.DeepEqual(lc.All(), expected) {
		t.Errorf("expected %v, got %v", expected, lc.All())
	}
}

func TestCount(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})

	if lc.Count() != 5 {
		t.Errorf("expected 5, got %d", lc.Count())
	}
}

func TestEager(t *testing.T) {
	lc := From([]int{1, 2, 3})
	items := lc.Eager()

	if len(items) != 3 {
		t.Errorf("expected 3, got %d", len(items))
	}
}

func TestFirst(t *testing.T) {
	lc := From([]int{10, 20, 30})
	v, ok := lc.First()

	if !ok || v != 10 {
		t.Errorf("expected 10, got %d", v)
	}

	v, ok = lc.First(func(item int, _ int) bool { return item > 15 })

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}
}

func TestLast(t *testing.T) {
	lc := From([]int{10, 20, 30})
	v, ok := lc.Last()

	if !ok || v != 30 {
		t.Errorf("expected 30, got %d", v)
	}

	v, ok = lc.Last(func(item int, _ int) bool { return item < 25 })

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}
}

func TestGet(t *testing.T) {
	lc := From([]int{10, 20, 30})
	v, ok := lc.Get(1)

	if !ok || v != 20 {
		t.Errorf("expected 20, got %d", v)
	}

	_, ok = lc.Get(10)

	if ok {
		t.Error("expected not found")
	}
}

func TestContains(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})

	if !lc.Contains(func(item int, _ int) bool { return item == 3 }) {
		t.Error("expected to contain 3")
	}

	if lc.Contains(func(item int, _ int) bool { return item == 10 }) {
		t.Error("expected not to contain 10")
	}
}

func TestSearch(t *testing.T) {
	lc := From([]int{10, 20, 30})
	idx, ok := lc.Search(func(item int, _ int) bool { return item == 20 })

	if !ok || idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}
}

func TestBefore(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	v, ok := lc.Before(func(item int, _ int) bool { return item == 3 })

	if !ok || v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
}

func TestAfter(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	v, ok := lc.After(func(item int, _ int) bool { return item == 3 })

	if !ok || v != 4 {
		t.Errorf("expected 4, got %d", v)
	}
}

func TestFilter(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	filtered := lc.Filter(func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{2, 4}

	if !reflect.DeepEqual(filtered.All(), expected) {
		t.Errorf("expected %v, got %v", expected, filtered.All())
	}
}

func TestReject(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	rejected := lc.Reject(func(item int, _ int) bool { return item%2 == 0 })
	expected := []int{1, 3, 5}

	if !reflect.DeepEqual(rejected.All(), expected) {
		t.Errorf("expected %v, got %v", expected, rejected.All())
	}
}

func TestMap(t *testing.T) {
	lc := From([]int{1, 2, 3})
	mapped := Map(lc, func(item int, _ int) int { return item * 2 })
	expected := []int{2, 4, 6}

	if !reflect.DeepEqual(mapped.All(), expected) {
		t.Errorf("expected %v, got %v", expected, mapped.All())
	}
}

func TestFlatMap(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := FlatMap(lc, func(item int, _ int) []int { return []int{item, item * 10} })
	expected := []int{1, 10, 2, 20, 3, 30}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestLazyTransformationsAreDeferredUntilConsumed(t *testing.T) {
	// SupportLazyCollectionIsLazyTest::testPadIsLazy
	// SupportLazyCollectionIsLazyTest::testRangeIsLazy
	// SupportLazyCollectionIsLazyTest::testRememberIsLazy
	// SupportLazyCollectionIsLazyTest::testTakeUntilTimeoutIsLazy
	// SupportLazyCollectionIsLazyTest::testTapDoesNotEnumerate
	// SupportLazyCollectionIsLazyTest::testTapEachIsLazy
	// SupportLazyCollectionIsLazyTest::testTimesIsLazy
	// SupportLazyCollectionIsLazyTest::testUnlessDoesNotEnumerate
	// SupportLazyCollectionIsLazyTest::testWhenDoesNotEnumerate
	tests := []struct {
		name  string
		build func(*Collection[int]) *Collection[int]
		want  []int
	}{
		{
			name: "map",
			build: func(lc *Collection[int]) *Collection[int] {
				return Map(lc, func(item int, _ int) int { return item * 2 })
			},
			want: []int{2, 4, 6},
		},
		{
			name:  "pluck",
			build: func(lc *Collection[int]) *Collection[int] { return Pluck(lc, func(item int) int { return item * 10 }) },
			want:  []int{10, 20, 30},
		},
		{
			name: "reject",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.Reject(func(item int, _ int) bool { return item%2 == 0 })
			},
			want: []int{1, 3},
		},
		{
			name:  "take",
			build: func(lc *Collection[int]) *Collection[int] { return lc.Take(2) },
			want:  []int{1, 2},
		},
		{
			name: "take-until",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.TakeUntil(func(item int, _ int) bool { return item == 3 })
			},
			want: []int{1, 2},
		},
		{
			name: "take-while",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.TakeWhile(func(item int, _ int) bool { return item < 3 })
			},
			want: []int{1, 2},
		},
		{
			name:  "skip",
			build: func(lc *Collection[int]) *Collection[int] { return lc.Skip(1) },
			want:  []int{2, 3},
		},
		{
			name: "skip-until",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.SkipUntil(func(item int, _ int) bool { return item == 2 })
			},
			want: []int{2, 3},
		},
		{
			name: "skip-while",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.SkipWhile(func(item int, _ int) bool { return item < 3 })
			},
			want: []int{3},
		},
		{
			name:  "slice",
			build: func(lc *Collection[int]) *Collection[int] { return lc.Slice(1, 2) },
			want:  []int{2, 3},
		},
		{
			name:  "nth",
			build: func(lc *Collection[int]) *Collection[int] { return lc.Nth(2) },
			want:  []int{1, 3},
		},
		{
			name:  "pad",
			build: func(lc *Collection[int]) *Collection[int] { return lc.Pad(5, 0) },
			want:  []int{1, 2, 3, 0, 0},
		},
		{
			name:  "remember",
			build: func(lc *Collection[int]) *Collection[int] { return lc.Remember() },
			want:  []int{1, 2, 3},
		},
		{
			name:  "take-until-timeout",
			build: func(lc *Collection[int]) *Collection[int] { return lc.TakeUntilTimeout(time.Hour) },
			want:  []int{1, 2, 3},
		},
		{
			name: "tap",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.Tap(func(inner *Collection[int]) {
					if inner == nil {
						t.Fatal("Tap received nil collection")
					}
				})
			},
			want: []int{1, 2, 3},
		},
		{
			name: "tap-each",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.TapEach(func(_ int, _ int) {})
			},
			want: []int{1, 2, 3},
		},
		{
			name: "unless",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.Unless(true, func(inner *Collection[int]) *Collection[int] {
					t.Fatal("Unless callback should not run when condition is true")

					return inner
				})
			},
			want: []int{1, 2, 3},
		},
		{
			name: "when",
			build: func(lc *Collection[int]) *Collection[int] {
				return lc.When(false, func(inner *Collection[int]) *Collection[int] {
					t.Fatal("When callback should not run when condition is false")

					return inner
				})
			},
			want: []int{1, 2, 3},
		},
		{
			name:  "unique",
			build: func(lc *Collection[int]) *Collection[int] { return Unique(lc, func(item int) int { return item % 2 }) },
			want:  []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enumerated := 0
			source := New(func(yield func(int) bool) {
				for _, item := range []int{1, 2, 3} {
					enumerated++

					if !yield(item) {
						return
					}
				}
			})

			result := tt.build(source)

			if enumerated != 0 {
				t.Fatalf("expected construction to stay lazy, enumerated %d items", enumerated)
			}

			if got := result.All(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}

	t.Run("range-constructor", func(t *testing.T) {
		lc := Range(1, 3)

		if got := lc.All(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
			t.Fatalf("expected range values, got %v", got)
		}
	})

	t.Run("times-constructor", func(t *testing.T) {
		calls := 0
		lc := Times(3, func(index int) int {
			calls++

			return index * 10
		})

		if calls != 0 {
			t.Fatalf("expected Times construction to stay lazy, called %d times", calls)
		}

		if got := lc.All(); !reflect.DeepEqual(got, []int{10, 20, 30}) {
			t.Fatalf("expected times values, got %v", got)
		}

		if calls != 3 {
			t.Fatalf("expected Times callback to run during enumeration, called %d times", calls)
		}
	})
}

func TestLazyTerminalOperationsEnumerateOnceOrStopEarly(t *testing.T) {
	t.Parallel()

	// SupportLazyCollectionIsLazyTest::testAvgEnumeratesOnce
	// SupportLazyCollectionIsLazyTest::testFindFirstOrFailIsLazy
	// SupportLazyCollectionIsLazyTest::testMaxEnumeratesOnce
	// SupportLazyCollectionIsLazyTest::testMinEnumeratesOnce
	// SupportLazyCollectionIsLazyTest::testReduceIsLazy
	// SupportLazyCollectionIsLazyTest::testReverseIsLazy
	// SupportLazyCollectionIsLazyTest::testRandomEnumeratesOnce
	// SupportLazyCollectionIsLazyTest::testSearchIsLazy
	// SupportLazyCollectionIsLazyTest::testShuffleIsLazy
	// SupportLazyCollectionIsLazyTest::testSlidingIsLazy
	// SupportLazyCollectionIsLazyTest::testSomeIsLazy
	// SupportLazyCollectionIsLazyTest::testSoleIsLazy
	// SupportLazyCollectionIsLazyTest::testSortIsLazy
	// SupportLazyCollectionIsLazyTest::testSortDescIsLazy
	// SupportLazyCollectionIsLazyTest::testSortByIsLazy
	// SupportLazyCollectionIsLazyTest::testSortByDescIsLazy
	// SupportLazyCollectionIsLazyTest::testSumEnumeratesOnce
	// SupportLazyCollectionIsLazyTest::testThrottleIsLazy
	// SupportLazyCollectionIsLazyTest::testToArrayEnumeratesOnce
	enumerated := 0
	source := New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})

	first, err := source.FirstOrFail(func(item int, _ int) bool { return item == 2 })

	if err != nil || first != 2 {
		t.Fatalf("FirstOrFail = %d, %v", first, err)
	}

	if enumerated != 2 {
		t.Fatalf("FirstOrFail enumerated %d items, want 2", enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})
	idx, ok := source.Search(func(item int, _ int) bool { return item == 3 })

	if !ok || idx != 2 || enumerated != 3 {
		t.Fatalf("Search idx=%d ok=%v enumerated=%d", idx, ok, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})

	if !source.Some(func(item int, _ int) bool { return item == 1 }) || enumerated != 1 {
		t.Fatalf("Some enumerated %d items, want 1", enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})
	sole, err := source.Sole(func(item int, _ int) bool { return item == 3 })

	if err != nil || sole != 3 || enumerated != 4 {
		t.Fatalf("Sole = %d, %v after %d enumerations", sole, err, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})
	sum := Reduce(source, func(carry int, item int, _ int) int { return carry + item }, 0)

	if sum != 10 || enumerated != 4 {
		t.Fatalf("Reduce sum=%d enumerated=%d", sum, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{3, 1, 4, 2} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})

	if got := Sum(source); got != 10 || enumerated != 4 {
		t.Fatalf("Sum=%d enumerated=%d", got, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{3, 1, 4, 2} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})

	if got := Avg(source); got != 2.5 || enumerated != 4 {
		t.Fatalf("Avg=%v enumerated=%d", got, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{3, 1, 4, 2} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})

	if got, ok := Min(source); !ok || got != 1 || enumerated != 4 {
		t.Fatalf("Min=%d ok=%v enumerated=%d", got, ok, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{3, 1, 4, 2} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})

	if got, ok := Max(source); !ok || got != 4 || enumerated != 4 {
		t.Fatalf("Max=%d ok=%v enumerated=%d", got, ok, enumerated)
	}

	lazyCases := []struct {
		name  string
		build func(*Collection[int]) *Collection[int]
		want  []int
	}{
		{
			name:  "reverse",
			build: func(lc *Collection[int]) *Collection[int] { return lc.Reverse() },
			want:  []int{2, 4, 1, 3},
		},
		{
			name:  "sort",
			build: func(lc *Collection[int]) *Collection[int] { return lc.Sort(func(a, b int) bool { return a < b }) },
			want:  []int{1, 2, 3, 4},
		},
		{
			name:  "sort-desc",
			build: func(lc *Collection[int]) *Collection[int] { return lc.SortDesc(func(a, b int) bool { return a < b }) },
			want:  []int{4, 3, 2, 1},
		},
		{
			name:  "sort-by",
			build: func(lc *Collection[int]) *Collection[int] { return SortBy(lc, func(item int) int { return item }) },
			want:  []int{1, 2, 3, 4},
		},
		{
			name:  "sort-by-desc",
			build: func(lc *Collection[int]) *Collection[int] { return SortByDesc(lc, func(item int) int { return item }) },
			want:  []int{4, 3, 2, 1},
		},
	}

	for _, tt := range lazyCases {
		t.Run(tt.name, func(t *testing.T) {
			enumerated = 0
			source = New(func(yield func(int) bool) {
				for _, item := range []int{3, 1, 4, 2} {
					enumerated++

					if !yield(item) {
						return
					}
				}
			})
			result := tt.build(source)

			if enumerated != 0 {
				t.Fatalf("%s construction enumerated %d items", tt.name, enumerated)
			}

			if got := result.All(); !reflect.DeepEqual(got, tt.want) || enumerated != 4 {
				t.Fatalf("%s got=%v enumerated=%d", tt.name, got, enumerated)
			}
		})
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})
	random := source.Random(2)

	if enumerated != 0 {
		t.Fatalf("Random construction enumerated %d items", enumerated)
	}

	if got := random.All(); len(got) != 2 || enumerated != 4 {
		t.Fatalf("Random got=%v enumerated=%d", got, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})
	shuffled := source.Shuffle()

	if enumerated != 0 {
		t.Fatalf("Shuffle construction enumerated %d items", enumerated)
	}

	if got := shuffled.All(); len(got) != 4 || enumerated != 4 {
		t.Fatalf("Shuffle got=%v enumerated=%d", got, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})
	sliding := Sliding(source, 2)

	if enumerated != 0 {
		t.Fatalf("Sliding construction enumerated %d items", enumerated)
	}

	if got := sliding.All(); !reflect.DeepEqual(got, [][]int{{1, 2}, {2, 3}, {3, 4}}) || enumerated != 4 {
		t.Fatalf("Sliding got=%v enumerated=%d", got, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})
	throttled := source.Throttle(0)

	if enumerated != 0 {
		t.Fatalf("Throttle construction enumerated %d items", enumerated)
	}

	if got := throttled.All(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) || enumerated != 4 {
		t.Fatalf("Throttle All got=%v enumerated=%d", got, enumerated)
	}

	enumerated = 0
	source = New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3, 4} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})

	if got := source.All(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) || enumerated != 4 {
		t.Fatalf("All got=%v enumerated=%d", got, enumerated)
	}
}

func TestTake(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	taken := lc.Take(3)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(taken.All(), expected) {
		t.Errorf("expected %v, got %v", expected, taken.All())
	}

	// Negative take
	taken2 := lc.Take(-2)
	expected2 := []int{4, 5}

	if !reflect.DeepEqual(taken2.All(), expected2) {
		t.Errorf("expected %v, got %v", expected2, taken2.All())
	}
}

func TestTakeUntil(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.TakeUntil(func(item int, _ int) bool { return item == 4 })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestTakeWhile(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.TakeWhile(func(item int, _ int) bool { return item < 4 })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkip(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.Skip(2)
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkipUntil(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.SkipUntil(func(item int, _ int) bool { return item == 3 })
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSkipWhile(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.SkipWhile(func(item int, _ int) bool { return item < 3 })
	expected := []int{3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestSlice(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.Slice(1, 3)
	expected := []int{2, 3, 4}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestChunk(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	chunks := lc.Chunk(2)

	if len(chunks) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(chunks))
	}

	if !reflect.DeepEqual(chunks[0], []int{1, 2}) {
		t.Errorf("expected [1 2], got %v", chunks[0])
	}

	if !reflect.DeepEqual(chunks[2], []int{5}) {
		t.Errorf("expected [5], got %v", chunks[2])
	}
}

func TestNth(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5, 6, 7, 8})
	result := lc.Nth(3)
	expected := []int{1, 4, 7}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestConcat(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Concat([]int{4, 5})
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestPad(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Pad(5, 0)
	expected := []int{1, 2, 3, 0, 0}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestEvery(t *testing.T) {
	lc := From([]int{2, 4, 6, 8})

	if !lc.Every(func(item int, _ int) bool { return item%2 == 0 }) {
		t.Error("expected all even")
	}
}

func TestImplode(t *testing.T) {
	lc := From([]string{"a", "b", "c"})
	result := lc.Implode(", ")

	if result != "a, b, c" {
		t.Errorf("expected 'a, b, c', got '%s'", result)
	}
}

func TestJoin(t *testing.T) {
	lc := From([]string{"a", "b", "c"})
	result := lc.Join(", ", " and ")

	if result != "a, b and c" {
		t.Errorf("expected 'a, b and c', got '%s'", result)
	}
}

func TestReduce(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	sum := Reduce(lc, func(carry int, item int, _ int) int { return carry + item }, 0)

	if sum != 15 {
		t.Errorf("expected 15, got %d", sum)
	}
}

func TestUnique(t *testing.T) {
	lc := From([]int{1, 2, 2, 3, 3, 3})
	result := Unique(lc, func(item int) int { return item })
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestPluck(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	lc := From([]User{{"Alice", 25}, {"Bob", 30}})
	names := Pluck(lc, func(u User) string { return u.Name })
	expected := []string{"Alice", "Bob"}

	if !reflect.DeepEqual(names.All(), expected) {
		t.Errorf("expected %v, got %v", expected, names.All())
	}
}

func TestGroupBy(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5, 6})
	groups := GroupBy(lc, func(item int) string {
		if item%2 == 0 {
			return "even"
		}

		return "odd"
	})

	if len(groups["even"].All()) != 3 {
		t.Errorf("expected 3 evens, got %d", len(groups["even"].All()))
	}
}

func TestKeyBy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	lc := From([]User{{1, "Alice"}, {2, "Bob"}})
	keyed := KeyBy(lc, func(u User) int { return u.ID })

	if keyed[1].Name != "Alice" {
		t.Error("expected Alice")
	}
}

func TestCountBy(t *testing.T) {
	lc := From([]string{"apple", "banana", "apple"})
	counts := CountBy(lc, func(item string) string { return item })

	if counts["apple"] != 2 {
		t.Errorf("expected 2, got %d", counts["apple"])
	}
}

func TestRemember(t *testing.T) {
	callCount := 0
	lc := New(func(yield func(int) bool) {
		callCount++

		for i := 1; i <= 3; i++ {
			if !yield(i) {
				return
			}
		}
	})

	remembered := lc.Remember()
	// First evaluation
	remembered.All()
	// Second evaluation - should use cache
	remembered.All()

	if callCount != 1 {
		t.Errorf("expected source called once, got %d", callCount)
	}
}

func TestContainsOneItem(t *testing.T) {
	if !From([]int{1}).ContainsOneItem() {
		t.Error("expected true")
	}

	if From([]int{1, 2}).ContainsOneItem() {
		t.Error("expected false")
	}
}

func TestWhen(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.When(true, func(lc *Collection[int]) *Collection[int] {
		return lc.Take(2)
	})

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}
}

func TestTapEach(t *testing.T) {
	lc := From([]int{1, 2, 3})
	sum := 0
	result := lc.TapEach(func(item int, _ int) {
		sum += item
	})
	// TapEach is lazy, need to evaluate
	result.All()

	if sum != 6 {
		t.Errorf("expected 6, got %d", sum)
	}
}

func TestFromSlice(t *testing.T) {
	items := []int{1, 2, 3}
	lc := From(items)

	if lc.Count() != 3 {
		t.Errorf("expected 3, got %d", lc.Count())
	}
}

func TestEach(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	sum := 0
	lc.Each(func(item int, _ int) bool {
		sum += item

		return item < 3
	})

	if sum != 6 { // 1+2+3
		t.Errorf("expected 6, got %d", sum)
	}
}

func TestSole(t *testing.T) {
	lc := From([]int{1, 2, 3})
	v, err := lc.Sole(func(item int, _ int) bool { return item == 2 })

	if err != nil || v != 2 {
		t.Errorf("expected 2, got %d, err: %v", v, err)
	}
}

func TestFirstOrFail(t *testing.T) {
	lc := From([]int{1, 2, 3})
	v, err := lc.FirstOrFail()

	if err != nil || v != 1 {
		t.Errorf("expected 1, got %d, err: %v", v, err)
	}

	empty := Empty[int]()
	_, err = empty.FirstOrFail()

	if err == nil {
		t.Error("expected error")
	}

	cases := []struct {
		name      string
		predicate func(int, int) bool
		want      int
		wantErr   bool
	}{
		{
			name:      "single match",
			predicate: func(item int, _ int) bool { return item == 2 },
			want:      2,
		},
		{
			name:      "no match",
			predicate: func(item int, _ int) bool { return item > 10 },
			wantErr:   true,
		},
		{
			name:      "multiple matches",
			predicate: func(item int, _ int) bool { return item > 1 },
			want:      2,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lc.FirstOrFail(tt.predicate)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}

				return
			}

			if err != nil || got != tt.want {
				t.Fatalf("expected %d, got %d, err: %v", tt.want, got, err)
			}
		})
	}
}

func TestFirstIsLazy(t *testing.T) {
	enumerated := 0
	lc := New(func(yield func(int) bool) {
		for _, item := range []int{1, 2, 3} {
			enumerated++

			if !yield(item) {
				return
			}
		}
	})

	v, ok := lc.First(func(item int, _ int) bool { return item == 2 })

	if !ok || v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}

	if enumerated != 2 {
		t.Fatalf("expected to enumerate 2 items, got %d", enumerated)
	}
}

func TestIter(t *testing.T) {
	lc := From([]int{10, 20, 30})
	sum := 0

	for item := range lc.Iter() {
		sum += item
	}

	if sum != 60 {
		t.Errorf("expected 60, got %d", sum)
	}
}

func TestCollect(t *testing.T) {
	lc := From([]int{1, 2, 3})
	items := lc.Collect()

	if !reflect.DeepEqual(items, []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", items)
	}
}

func TestContainsManyItems(t *testing.T) {
	if !From([]int{1, 2}).ContainsManyItems() {
		t.Error("expected true for 2 items")
	}

	if From([]int{1}).ContainsManyItems() {
		t.Error("expected false for 1 item")
	}

	if Empty[int]().ContainsManyItems() {
		t.Error("expected false for empty")
	}
}

func TestSome(t *testing.T) {
	lc := From([]int{1, 2, 3})

	if !lc.Some(func(item int, _ int) bool { return item == 2 }) {
		t.Error("expected true")
	}

	if lc.Some(func(item int, _ int) bool { return item > 10 }) {
		t.Error("expected false")
	}
}

func TestDoesntContain(t *testing.T) {
	lc := From([]int{1, 2, 3})

	if !lc.DoesntContain(func(item int, _ int) bool { return item == 99 }) {
		t.Error("expected true")
	}

	if lc.DoesntContain(func(item int, _ int) bool { return item == 2 }) {
		t.Error("expected false")
	}
}

func TestTap(t *testing.T) {
	lc := From([]int{1, 2, 3})
	called := false
	result := lc.Tap(func(c *Collection[int]) {
		called = true
	})

	if !called {
		t.Error("expected callback to be called")
	}

	if result != lc {
		t.Error("expected Tap to return self")
	}
}

func TestHas(t *testing.T) {
	lc := From([]int{10, 20, 30})

	if !lc.Has(0) {
		t.Error("expected true for index 0")
	}

	if !lc.Has(2) {
		t.Error("expected true for index 2")
	}

	if lc.Has(5) {
		t.Error("expected false for index 5")
	}
}

func TestHasAny(t *testing.T) {
	lc := From([]int{10, 20, 30})

	if !lc.HasAny(0, 5) {
		t.Error("expected true when at least one valid")
	}

	if lc.HasAny(5, 6, 7) {
		t.Error("expected false when all invalid")
	}
}

func TestHasSole(t *testing.T) {
	single := From([]int{42})

	if !single.HasSole() {
		t.Error("expected true for single item without predicate")
	}

	multi := From([]int{1, 2, 3})

	if multi.HasSole() {
		t.Error("expected false for multiple items without predicate")
	}

	if !multi.HasSole(func(item int, _ int) bool { return item == 2 }) {
		t.Error("expected true for single match")
	}

	if multi.HasSole(func(item int, _ int) bool { return item > 1 }) {
		t.Error("expected false for multiple matches")
	}

	if multi.HasSole(func(item int, _ int) bool { return item > 10 }) {
		t.Error("expected false for no matches")
	}
}

func TestChunkWhile(t *testing.T) {
	lc := From([]int{1, 1, 2, 2, 3})
	chunks := lc.ChunkWhile(func(item int, _ int, current []int) bool {
		return item == current[len(current)-1]
	})

	if len(chunks) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(chunks))
	}

	if !reflect.DeepEqual(chunks[0], []int{1, 1}) {
		t.Errorf("expected [1 1], got %v", chunks[0])
	}

	if !reflect.DeepEqual(chunks[1], []int{2, 2}) {
		t.Errorf("expected [2 2], got %v", chunks[1])
	}

	if !reflect.DeepEqual(chunks[2], []int{3}) {
		t.Errorf("expected [3], got %v", chunks[2])
	}
}

func TestTakeUntilTimeout(t *testing.T) {
	lc := From([]int{1, 2, 3, 4, 5})
	result := lc.TakeUntilTimeout(time.Second)
	items := result.All()

	if !reflect.DeepEqual(items, []int{1, 2, 3, 4, 5}) {
		t.Errorf("expected all items with generous timeout, got %v", items)
	}
}

func TestThrottle(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Throttle(time.Millisecond)
	items := result.All()

	if !reflect.DeepEqual(items, []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", items)
	}
}

func TestWhenFalseWithDefault(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.When(false, func(c *Collection[int]) *Collection[int] {
		return c.Take(1)
	}, func(c *Collection[int]) *Collection[int] {
		return c.Take(2)
	})

	if result.Count() != 2 {
		t.Errorf("expected 2 (default applied), got %d", result.Count())
	}
}

func TestWhenFalseWithoutDefault(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.When(false, func(c *Collection[int]) *Collection[int] {
		return c.Take(1)
	})

	if result.Count() != 3 {
		t.Errorf("expected 3 (unchanged), got %d", result.Count())
	}
}

func TestWhenEmpty(t *testing.T) {
	empty := Empty[int]()
	result := empty.WhenEmpty(func(c *Collection[int]) *Collection[int] {
		return From([]int{1, 2, 3})
	})

	if result.Count() != 3 {
		t.Errorf("expected 3, got %d", result.Count())
	}

	nonEmpty := From([]int{1})
	called := false
	nonEmpty.WhenEmpty(func(c *Collection[int]) *Collection[int] {
		called = true

		return c
	})

	if called {
		t.Error("expected callback not called for non-empty")
	}
}

func TestWhenNotEmpty(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.WhenNotEmpty(func(c *Collection[int]) *Collection[int] {
		return c.Take(1)
	})

	if result.Count() != 1 {
		t.Errorf("expected 1, got %d", result.Count())
	}

	empty := Empty[int]()
	called := false
	empty.WhenNotEmpty(func(c *Collection[int]) *Collection[int] {
		called = true

		return c
	})

	if called {
		t.Error("expected callback not called for empty")
	}
}

func TestUnless(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Unless(false, func(c *Collection[int]) *Collection[int] {
		return c.Take(2)
	})

	if result.Count() != 2 {
		t.Errorf("expected 2, got %d", result.Count())
	}

	result2 := lc.Unless(true, func(c *Collection[int]) *Collection[int] {
		return c.Take(1)
	})

	if result2.Count() != 3 {
		t.Errorf("expected 3 (unchanged), got %d", result2.Count())
	}
}

func TestDump(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Dump()
	items := result.All()

	if !reflect.DeepEqual(items, []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", items)
	}
}

func TestSoleEmpty(t *testing.T) {
	lc := Empty[int]()
	_, err := lc.Sole()

	if err == nil {
		t.Error("expected error for empty collection")
	}
}

func TestSoleMultiple(t *testing.T) {
	lc := From([]int{1, 2, 3})
	_, err := lc.Sole()

	if err == nil {
		t.Error("expected error for multiple items")
	}
}

func TestSoleSingle(t *testing.T) {
	lc := From([]int{42})
	v, err := lc.Sole()

	if err != nil || v != 42 {
		t.Errorf("expected 42, got %d, err: %v", v, err)
	}
}

func TestSoleCallbackNoMatch(t *testing.T) {
	lc := From([]int{1, 2, 3})
	_, err := lc.Sole(func(item int, _ int) bool { return item > 10 })

	if err == nil {
		t.Error("expected error when no match")
	}
}

func TestSoleCallbackMultipleMatches(t *testing.T) {
	lc := From([]int{1, 2, 3})
	_, err := lc.Sole(func(item int, _ int) bool { return item > 1 })

	if err == nil {
		t.Error("expected error for multiple matches")
	}
}

func TestPadNegative(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Pad(-5, 0)
	expected := []int{0, 0, 1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v, got %v", expected, result.All())
	}
}

func TestPadAlreadySufficient(t *testing.T) {
	lc := From([]int{1, 2, 3})
	result := lc.Pad(2, 0)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result.All(), expected) {
		t.Errorf("expected %v (unchanged), got %v", expected, result.All())
	}
}

func TestRangeSingle(t *testing.T) {
	lc := Range(5, 5)
	expected := []int{5}

	if !reflect.DeepEqual(lc.All(), expected) {
		t.Errorf("expected %v, got %v", expected, lc.All())
	}
}

func TestTimesZero(t *testing.T) {
	lc := Times(0, func(i int) int { return i })

	if lc.Count() != 0 {
		t.Errorf("expected 0, got %d", lc.Count())
	}
}

func TestIsEmptyNonEmpty(t *testing.T) {
	lc := From([]int{1})

	if lc.IsEmpty() {
		t.Error("expected false for non-empty")
	}
}

func TestEveryEmpty(t *testing.T) {
	lc := Empty[int]()

	if !lc.Every(func(item int, _ int) bool { return false }) {
		t.Error("expected true for empty collection")
	}
}

func TestEveryFailing(t *testing.T) {
	lc := From([]int{1, 2, 3})

	if lc.Every(func(item int, _ int) bool { return item > 2 }) {
		t.Error("expected false")
	}
}
