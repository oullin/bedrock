package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// OnceTest::testResultMemoization
// OnceTest::testCallableIsCalledOnce
// OnceTest::testFlush
// OnceTest::testIsNotMemoizedWhenCallableUsesChanges
// OnceTest::testStaticMemoization
// OnceTest::testMemoizationWhenOnceIsWithinClosure
// OnceTest::testMemoizationOnGlobalFunctions
// OnceTest::testDisable
// OnceTest::testTemporaryDisable
// OnceTest::testResultIsDifferentWhenCalledFromDifferentClosures
// OnceTest::testRecursiveOnceCalls
// OnceTest::testGlobalClosures
// OnceTest::testMemoizationNullValues

func TestOnceMemoizesResultAndCallable(t *testing.T) {
	FlushOnce()

	calls := 0
	value := func() int {
		calls++

		return 10
	}

	if got := Once("value", value); got != 10 {
		t.Fatalf("Once first = %d", got)
	}

	if got := Once("value", value); got != 10 {
		t.Fatalf("Once second = %d", got)
	}

	if calls != 1 {
		t.Fatalf("expected callback to run once, got %d", calls)
	}
}

func TestOnceFlushAndKeyChanges(t *testing.T) {
	FlushOnce()

	calls := 0
	next := func() int {
		calls++

		return calls
	}

	if got := Once("flush", next); got != 1 {
		t.Fatalf("Once before flush = %d", got)
	}

	FlushOnce()

	if got := Once("flush", next); got != 2 {
		t.Fatalf("Once after flush = %d", got)
	}

	if got := Once("different-key", next); got != 3 {
		t.Fatalf("Once different key = %d", got)
	}
}

func TestOnceDisableAndTemporaryDisable(t *testing.T) {
	FlushOnce()

	calls := 0
	next := func() int {
		calls++

		return calls
	}

	DisableOnce()
	if got := Once("disabled", next); got != 1 {
		t.Fatalf("Once disabled first = %d", got)
	}
	if got := Once("disabled", next); got != 2 {
		t.Fatalf("Once disabled second = %d", got)
	}
	EnableOnce()

	if got := Once("enabled", next); got != 3 {
		t.Fatalf("Once enabled first = %d", got)
	}
	if got := Once("enabled", next); got != 3 {
		t.Fatalf("Once enabled second = %d", got)
	}

	WithoutOnce(func() {
		if got := Once("enabled", next); got != 4 {
			t.Fatalf("Once temporarily disabled = %d", got)
		}
	})

	if got := Once("enabled", next); got != 3 {
		t.Fatalf("Once restored cached value = %d", got)
	}
}

func TestOnceClosureGlobalRecursiveAndNilResults(t *testing.T) {
	FlushOnce()

	closure := func(key string) int {
		return Once(key, func() int { return len(key) })
	}

	if first, second := closure("alpha"), closure("alpha"); first != second {
		t.Fatalf("closure memoization mismatch: %d != %d", first, second)
	}

	if first, second := closure("alpha"), closure("beta"); first == second {
		t.Fatalf("different closures/keys should produce different results")
	}

	var recursive func(int) int
	recursive = func(value int) int {
		if value == 0 {
			return Once("recursive-base", func() int { return 1 })
		}

		return value + recursive(value-1)
	}

	if got := recursive(3); got != 7 {
		t.Fatalf("recursive once = %d", got)
	}

	var calls int
	nilResult := Once[*int]("nil", func() *int {
		calls++

		return nil
	})
	again := Once[*int]("nil", func() *int {
		calls++

		return nil
	})

	if nilResult != nil || again != nil || calls != 1 {
		t.Fatalf("nil memoization failed: first=%v second=%v calls=%d", nilResult, again, calls)
	}
}
