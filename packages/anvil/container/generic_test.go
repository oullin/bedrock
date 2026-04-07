package container

import (
	"errors"
	"testing"
)

type testService struct {
	Name string
}

type testLogger interface {
	Log(msg string)
}

type testConsoleLogger struct{}

func (testConsoleLogger) Log(string) {}

func TestMakeGeneric(t *testing.T) {
	t.Parallel()

	c := New()
	svc := &testService{Name: "auth"}
	c.Instance("svc", svc)

	got, err := Make[*testService](c, "svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Name != "auth" {
		t.Fatalf("want Name=%q, got %q", "auth", got.Name)
	}

	if got != svc {
		t.Fatal("expected same pointer")
	}
}

func TestMakeGenericInterface(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("logger", testConsoleLogger{})

	got, err := Make[testLogger](c, "logger")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestMakeGenericMismatch(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("val", "a string")

	_, err := Make[int](c, "val")
	if err == nil {
		t.Fatal("expected type mismatch error")
	}

	if !errors.Is(err, ErrTypeMismatch) {
		t.Fatalf("expected ErrTypeMismatch, got: %v", err)
	}
}

func TestMustMakeGenericPanic(t *testing.T) {
	t.Parallel()

	c := New()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic from MustMake")
		}
	}()

	MustMake[string](c, "missing")
}
