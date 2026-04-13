package conditionable_test

import (
	"slices"
	"testing"

	"github.com/bedrock/packages/conditionable"
)

// logger mirrors Laravel's ConditionableLogger test helper.
type logger struct {
	Values []any
	Toggle bool
}

func newLogger() *logger {
	return &logger{}
}

func (l *logger) Log(values ...any) *logger {
	l.Values = append(l.Values, values...)
	return l
}

func (l *logger) Has(value any) bool {
	return slices.Contains(l.Values, value)
}

func (l *logger) DoToggle() *logger {
	l.Toggle = !l.Toggle
	return l
}

// --- When with static condition ---

func TestWhenConditionCallback(t *testing.T) {
	t.Parallel()

	l := conditionable.When(
		newLogger(),
		2,
		func(l *logger, v any) *logger { return l.Log("when", v) },
		func(l *logger, v any) *logger { return l.Log("default", v) },
	)

	assertValues(t, l.Values, []any{"when", 2})
}

// --- When with closure condition ---

func TestWhenFuncConditionCallback(t *testing.T) {
	t.Parallel()

	l := conditionable.WhenFunc(
		newLogger().Log("init"),
		func(l *logger) any { return l.Has("init") },
		func(l *logger, v any) *logger { return l.Log("when", v) },
		func(l *logger, v any) *logger { return l.Log("default", v) },
	)

	assertValues(t, l.Values, []any{"init", "when", true})
}

// --- When default with static nil condition ---

func TestWhenDefaultCallback(t *testing.T) {
	t.Parallel()

	l := conditionable.When(
		newLogger(),
		nil,
		func(l *logger, v any) *logger { return l.Log("when", v) },
		func(l *logger, v any) *logger { return l.Log("default", v) },
	)

	assertValues(t, l.Values, []any{"default", nil})
}

// --- When default with closure condition ---

func TestWhenFuncDefaultCallback(t *testing.T) {
	t.Parallel()

	l := conditionable.WhenFunc(
		newLogger(),
		func(l *logger) any { return l.Has("missing") },
		func(l *logger, v any) *logger { return l.Log("when", v) },
		func(l *logger, v any) *logger { return l.Log("default", v) },
	)

	assertValues(t, l.Values, []any{"default", false})
}

// --- Unless with static nil condition ---

func TestUnlessConditionCallback(t *testing.T) {
	t.Parallel()

	l := conditionable.Unless(
		newLogger(),
		nil,
		func(l *logger, v any) *logger { return l.Log("unless", v) },
		func(l *logger, v any) *logger { return l.Log("default", v) },
	)

	assertValues(t, l.Values, []any{"unless", nil})
}

// --- Unless with closure condition ---

func TestUnlessFuncConditionCallback(t *testing.T) {
	t.Parallel()

	l := conditionable.UnlessFunc(
		newLogger(),
		func(l *logger) any { return l.Has("missing") },
		func(l *logger, v any) *logger { return l.Log("unless", v) },
		func(l *logger, v any) *logger { return l.Log("default", v) },
	)

	assertValues(t, l.Values, []any{"unless", false})
}

// --- Unless default with static truthy condition ---

func TestUnlessDefaultCallback(t *testing.T) {
	t.Parallel()

	l := conditionable.Unless(
		newLogger(),
		2,
		func(l *logger, v any) *logger { return l.Log("unless", v) },
		func(l *logger, v any) *logger { return l.Log("default", v) },
	)

	assertValues(t, l.Values, []any{"default", 2})
}

// --- Unless default with closure condition ---

func TestUnlessFuncDefaultCallback(t *testing.T) {
	t.Parallel()

	l := conditionable.UnlessFunc(
		newLogger().Log("init"),
		func(l *logger) any { return l.Has("init") },
		func(l *logger, v any) *logger { return l.Log("unless", v) },
		func(l *logger, v any) *logger { return l.Log("default", v) },
	)

	assertValues(t, l.Values, []any{"init", "default", true})
}

// --- When without default (no-op on falsy) ---

func TestWhenWithoutDefault(t *testing.T) {
	t.Parallel()

	l := conditionable.When(
		newLogger().Log("init"),
		false,
		func(l *logger, v any) *logger { return l.Log("when") },
	)

	assertValues(t, l.Values, []any{"init"})
}

// --- Unless without default (no-op on truthy) ---

func TestUnlessWithoutDefault(t *testing.T) {
	t.Parallel()

	l := conditionable.Unless(
		newLogger().Log("init"),
		true,
		func(l *logger, v any) *logger { return l.Log("unless") },
	)

	assertValues(t, l.Values, []any{"init"})
}

// --- When proxy ---

func TestWhenProxy(t *testing.T) {
	t.Parallel()

	// when(true) executes, when(false) skips
	l := newLogger()
	l = conditionable.NewProxy(l, true).Then(func(l *logger) *logger {
		return l.Log("one")
	})
	l = conditionable.NewProxy(l, false).Then(func(l *logger) *logger {
		return l.Log("two")
	})

	assertValues(t, l.Values, []any{"one"})
}

func TestWhenProxyWithFunc(t *testing.T) {
	t.Parallel()

	l := newLogger().Log("init")

	// when(func) resolves to true → executes
	l = conditionable.NewProxy(l, l.Has("init")).Then(func(l *logger) *logger {
		return l.Log("one")
	})

	// when(func) resolves to false → skips
	l = conditionable.NewProxy(l, l.Has("missing")).Then(func(l *logger) *logger {
		return l.Log("two")
	})

	// toggle starts false, so proxy skips
	l = conditionable.NewProxy(l, l.Toggle).Then(func(l *logger) *logger {
		return l.Log("three")
	})

	l = l.DoToggle()

	// toggle now true, proxy executes
	l = conditionable.NewProxy(l, l.Toggle).Then(func(l *logger) *logger {
		return l.Log("four")
	})

	assertValues(t, l.Values, []any{"init", "one", "four"})
}

// --- Unless proxy ---

func TestUnlessProxy(t *testing.T) {
	t.Parallel()

	// unless(true) skips, unless(false) executes
	l := newLogger()
	l = conditionable.NewUnlessProxy(l, true).Then(func(l *logger) *logger {
		return l.Log("one")
	})
	l = conditionable.NewUnlessProxy(l, false).Then(func(l *logger) *logger {
		return l.Log("two")
	})

	assertValues(t, l.Values, []any{"two"})
}

func TestUnlessProxyWithFunc(t *testing.T) {
	t.Parallel()

	l := newLogger().Log("init")

	// unless(has("init")) → condition true → skip
	l = conditionable.NewUnlessProxy(l, l.Has("init")).Then(func(l *logger) *logger {
		return l.Log("one")
	})

	// unless(has("missing")) → condition false → execute
	l = conditionable.NewUnlessProxy(l, l.Has("missing")).Then(func(l *logger) *logger {
		return l.Log("two")
	})

	// toggle is false, unless(false) → execute
	l = conditionable.NewUnlessProxy(l, l.Toggle).Then(func(l *logger) *logger {
		return l.Log("three")
	})

	l = l.DoToggle()

	// toggle is true, unless(true) → skip
	l = conditionable.NewUnlessProxy(l, l.Toggle).Then(func(l *logger) *logger {
		return l.Log("four")
	})

	assertValues(t, l.Values, []any{"init", "two", "three"})
}

// assertValues compares two any slices for deep equality.
func assertValues(t *testing.T, got, want []any) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("got %v (len %d), want %v (len %d)", got, len(got), want, len(want))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %v (%T), want %v (%T)", i, got[i], got[i], want[i], want[i])
		}
	}
}
