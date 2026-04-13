package container_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/container"
)

func TestCallInvokesCallable(t *testing.T) {
	t.Parallel()

	c := newContainer()

	v, err := c.Call(func(cc *container.Container, _ map[string]any) (any, error) {
		return "called", nil
	}, nil)

	if err != nil {
		t.Fatal(err)
	}

	if v != "called" {
		t.Fatalf("expected called, got %v", v)
	}
}

func TestCallPassesParameters(t *testing.T) {
	t.Parallel()

	c := newContainer()

	v, err := c.Call(func(_ *container.Container, params map[string]any) (any, error) {
		return params["greeting"], nil
	}, map[string]any{"greeting": "hello"})

	if err != nil {
		t.Fatal(err)
	}

	if v != "hello" {
		t.Fatalf("expected hello, got %v", v)
	}
}

func TestCallWithContainerDependencies(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	v, err := c.Call(func(cc *container.Container, _ map[string]any) (any, error) {
		return cc.Make("name")
	}, nil)

	if err != nil {
		t.Fatal(err)
	}

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestCallReturnsError(t *testing.T) {
	t.Parallel()

	c := newContainer()

	_, err := c.Call(func(_ *container.Container, _ map[string]any) (any, error) {
		return nil, errors.New("oops")
	}, nil)

	if err == nil || err.Error() != "oops" {
		t.Fatalf("expected error oops, got %v", err)
	}
}

func TestWrapReturnsDeferredClosure(t *testing.T) {
	t.Parallel()

	c := newContainer()

	count := 0
	wrapped := c.Wrap(func(_ *container.Container, _ map[string]any) (any, error) {
		count++

		return count, nil
	}, nil)

	if count != 0 {
		t.Fatal("wrap should not invoke immediately")
	}

	v, _ := wrapped()

	if v != 1 {
		t.Fatalf("expected 1, got %v", v)
	}

	v, _ = wrapped()

	if v != 2 {
		t.Fatalf("expected 2, got %v", v)
	}
}

func TestBindMethodAndCallMethodBinding(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.BindMethod("App@handle", func(_ *container.Container, params map[string]any) (any, error) {
		return params["_instance"], nil
	})

	v, err := c.CallMethodBinding("App@handle", "the-instance")

	if err != nil {
		t.Fatal(err)
	}

	if v != "the-instance" {
		t.Fatalf("expected the-instance, got %v", v)
	}
}

func TestHasMethodBinding(t *testing.T) {
	t.Parallel()

	c := newContainer()

	if c.HasMethodBinding("App@handle") {
		t.Fatal("should not have binding")
	}

	c.BindMethod("App@handle", func(_ *container.Container, _ map[string]any) (any, error) {
		return nil, nil
	})

	if !c.HasMethodBinding("App@handle") {
		t.Fatal("should have binding")
	}
}

func TestCallMethodBindingNotBound(t *testing.T) {
	t.Parallel()

	c := newContainer()

	_, err := c.CallMethodBinding("App@handle", nil)

	if !errors.Is(err, container.ErrMethodNotBound) {
		t.Fatalf("expected ErrMethodNotBound, got %v", err)
	}
}
