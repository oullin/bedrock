package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportReflectorTest::testGetClassName
// SupportReflectorTest::testEmptyClassName
// SupportReflectorTest::testStringTypeName
// SupportReflectorTest::testParameterSubclassOfInterface
// SupportReflectorTest::testIsCallable

type reflectorSample struct{}

type reflectorContract interface {
	ReflectorMethod()
}

type reflectorImplementation struct{}

func (reflectorImplementation) ReflectorMethod() {}

func TestReflectorTypeName(t *testing.T) {
	t.Parallel()

	if got := TypeName(reflectorSample{}); got != "reflectorSample" {
		t.Fatalf("TypeName(struct) = %q", got)
	}

	if got := TypeName(nil); got != "" {
		t.Fatalf("TypeName(nil) = %q", got)
	}

	if got := TypeName("value"); got != "string" {
		t.Fatalf("TypeName(string) = %q", got)
	}
}

func TestReflectorCallableAndInterface(t *testing.T) {
	t.Parallel()

	if !IsCallable(func() {}) {
		t.Fatal("expected function to be callable")
	}

	if IsCallable("not callable") {
		t.Fatal("string should not be callable")
	}

	if !Implements(reflectorImplementation{}, (*reflectorContract)(nil)) {
		t.Fatal("expected implementation to satisfy contract")
	}
}
