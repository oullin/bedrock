package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportConditionableTest::testWhenConditionCallback
// SupportConditionableTest::testWhenDefaultCallback
// SupportConditionableTest::testUnlessConditionCallback
// SupportConditionableTest::testUnlessDefaultCallback
// SupportTappableTest::testTappableClassWithCallback
// SupportTappableTest::testTappableClassWithoutCallback

func TestConditionableWhenUnless(t *testing.T) {
	t.Parallel()

	if got := When("base", true, func(value string) string { return value + "-when" }); got != "base-when" {
		t.Fatalf("When condition callback = %q", got)
	}

	if got := When("base", false, func(value string) string { return "unused" }, func(value string) string { return value + "-default" }); got != "base-default" {
		t.Fatalf("When default callback = %q", got)
	}

	if got := Unless("base", false, func(value string) string { return value + "-unless" }); got != "base-unless" {
		t.Fatalf("Unless condition callback = %q", got)
	}

	if got := Unless("base", true, func(value string) string { return "unused" }, func(value string) string { return value + "-default" }); got != "base-default" {
		t.Fatalf("Unless default callback = %q", got)
	}
}

func TestTapValue(t *testing.T) {
	t.Parallel()

	called := false
	value := TapValue("base", func(value string) {
		called = value == "base"
	})

	if value != "base" || !called {
		t.Fatalf("TapValue with callback = %q called=%v", value, called)
	}

	if got := TapValue("base"); got != "base" {
		t.Fatalf("TapValue without callback = %q", got)
	}
}
