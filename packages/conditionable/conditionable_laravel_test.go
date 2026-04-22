package conditionable_test

import (
	"testing"

	"github.com/bedrock/packages/conditionable"
)

// Port of Illuminate\Tests\Conditionable\ConditionableTest::testWhen
func TestLaravelConditionableWhen(t *testing.T) {
	t.Parallel()

	target := newLogger()

	if conditionable.NewProxy(target, true) == nil {
		t.Fatal("expected true condition to create a proxy")
	}

	if conditionable.NewProxy(target, false) == nil {
		t.Fatal("expected false condition to create a proxy")
	}

	var nilCallback func(*logger, bool) *logger
	got := conditionable.When(target, false, nilCallback)

	if got != target {
		t.Fatalf("expected false condition with nil callback to return target unchanged")
	}

	got = conditionable.When(target, true, func(l *logger, _ bool) *logger {
		return l.Log("when")
	})

	if got != target {
		t.Fatalf("expected true condition callback to return target")
	}

	assertValues(t, target.Values, []any{"when"})
}

// Port of Illuminate\Tests\Conditionable\ConditionableTest::testUnless
func TestLaravelConditionableUnless(t *testing.T) {
	t.Parallel()

	target := newLogger()

	if conditionable.NewUnlessProxy(target, true) == nil {
		t.Fatal("expected true condition to create an unless proxy")
	}

	if conditionable.NewUnlessProxy(target, false) == nil {
		t.Fatal("expected false condition to create an unless proxy")
	}

	var nilCallback func(*logger, bool) *logger
	got := conditionable.Unless(target, true, nilCallback)

	if got != target {
		t.Fatalf("expected true condition with nil callback to return target unchanged")
	}

	got = conditionable.Unless(target, false, func(l *logger, _ bool) *logger {
		return l.Log("unless")
	})

	if got != target {
		t.Fatalf("expected false condition callback to return target")
	}

	assertValues(t, target.Values, []any{"unless"})
}
