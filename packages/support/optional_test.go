package support

import "testing"

// Port of @bedrock\Tests\Support\SupportOptionalTest::testGetExistItemOnObject
func TestOptionalGetExistingValue(t *testing.T) {
	t.Parallel()

	opt := Some("hello")
	val, ok := opt.Get()

	if !ok || val != "hello" {
		t.Errorf("expected ('hello', true), got (%q, %v)", val, ok)
	}
}

// Port of @bedrock\Tests\Support\SupportOptionalTest::testGetNotExistItemOnObject
func TestOptionalGetMissingValue(t *testing.T) {
	t.Parallel()

	opt := None[string]()
	val, ok := opt.Get()

	if ok || val != "" {
		t.Errorf("expected ('', false), got (%q, %v)", val, ok)
	}
}

// Port of @bedrock\Tests\Support\SupportOptionalTest::testIssetExistItemOnObject
func TestOptionalIsPresentTrue(t *testing.T) {
	t.Parallel()

	opt := Some(42)

	if !opt.IsPresent() {
		t.Error("expected IsPresent() = true")
	}

	if opt.IsEmpty() {
		t.Error("expected IsEmpty() = false")
	}
}

// Port of @bedrock\Tests\Support\SupportOptionalTest::testIssetNotExistItemOnObject
func TestOptionalIsPresentFalse(t *testing.T) {
	t.Parallel()

	opt := None[int]()

	if opt.IsPresent() {
		t.Error("expected IsPresent() = false")
	}

	if !opt.IsEmpty() {
		t.Error("expected IsEmpty() = true")
	}
}

// Port of @bedrock\Tests\Support\SupportOptionalTest::testIssetExistItemOnNull
func TestOptionalNilPointerIsSafe(t *testing.T) {
	t.Parallel()

	var ptr *string = nil
	opt := Opt(ptr)

	if opt.IsPresent() {
		t.Error("expected IsPresent() = false for nil pointer")
	}
}

func TestOptionalOrElse(t *testing.T) {
	t.Parallel()

	opt := None[string]()

	if got := opt.OrElse("default"); got != "default" {
		t.Errorf("expected 'default', got %q", got)
	}

	opt2 := Some("value")

	if got := opt2.OrElse("default"); got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
}

func TestOptionalIfPresent(t *testing.T) {
	t.Parallel()

	called := false
	Some(10).IfPresent(func(n int) {
		called = true

		if n != 10 {
			t.Errorf("expected 10, got %d", n)
		}
	})

	if !called {
		t.Error("IfPresent callback was not called for Some")
	}

	None[int]().IfPresent(func(int) {
		t.Error("IfPresent callback should not be called for None")
	})
}

func TestOptionalMap(t *testing.T) {
	t.Parallel()

	mapped := Some(5).Map(func(n int) int { return n * 2 })
	val, ok := mapped.Get()

	if !ok || val != 10 {
		t.Errorf("expected (10, true), got (%d, %v)", val, ok)
	}

	mappedNone := None[int]().Map(func(n int) int { return n * 2 })

	if mappedNone.IsPresent() {
		t.Error("Map on None should return None")
	}
}

func TestOptionalFilter(t *testing.T) {
	t.Parallel()

	passes := Some(10).Filter(func(n int) bool { return n > 5 })

	if !passes.IsPresent() {
		t.Error("Filter should keep value when predicate passes")
	}

	fails := Some(3).Filter(func(n int) bool { return n > 5 })

	if fails.IsPresent() {
		t.Error("Filter should remove value when predicate fails")
	}
}
