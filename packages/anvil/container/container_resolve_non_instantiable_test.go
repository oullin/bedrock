package container

import (
	"errors"
	"testing"
)

// Upstream's ContainerResolveNonInstantiableTest covers resolving
// non-instantiable PHP classes (interfaces, abstracts) with default/variadic
// parameters. All 3 Upstream tests rely on PHP reflection and are skipped.
// This file tests the Go equivalents: error paths for unbound and broken
// factories.
//
// INTENTIONAL-SKIP: testResolvingNonInstantiableWithDefaultRemovesWiths
// INTENTIONAL-SKIP: testResolvingNonInstantiableWithVariadicRemovesWiths
// INTENTIONAL-SKIP: testResolveVariadicPrimitive

func TestUnknownEntryThrowsException(t *testing.T) {
	t.Parallel()

	c := New()
	_, err := c.Make("missing")

	if err == nil {
		t.Fatal("expected error for unbound abstract")
	}

	if !errors.Is(err, ErrNotBound) {
		t.Fatalf("expected ErrNotBound, got: %v", err)
	}
}

func TestBindingResolutionExceptionMessage(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("broken", func(_ *Container) (any, error) {
		return nil, errors.New("factory failed")
	})

	_, err := c.Make("broken")
	if err == nil {
		t.Fatal("expected error from factory")
	}

	if !errors.Is(err, ErrResolve) {
		t.Fatalf("expected ErrResolve, got: %v", err)
	}
}

func TestSingletonFactoryErrorCached(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("broken", func(_ *Container) (any, error) {
		return nil, errors.New("init failed")
	})

	_, err1 := c.Make("broken")
	_, err2 := c.Make("broken")

	if !errors.Is(err1, ErrResolve) {
		t.Fatalf("first call: expected ErrResolve, got %v", err1)
	}

	if !errors.Is(err2, ErrResolve) {
		t.Fatalf("second call: expected same error, got %v", err2)
	}
}

func TestMakeForUnboundReturnsError(t *testing.T) {
	t.Parallel()

	c := New()
	_, err := c.MakeFor("consumer", "unbound")

	if err == nil {
		t.Fatal("expected error for unbound abstract with no contextual binding")
	}

	if !errors.Is(err, ErrNotBound) {
		t.Fatalf("expected ErrNotBound, got: %v", err)
	}
}

func TestFactoryReturningNilNoError(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("nullable", func(_ *Container) (any, error) {
		return nil, nil
	})

	got, err := c.Make("nullable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != nil {
		t.Fatalf("want nil, got %v", got)
	}
}
