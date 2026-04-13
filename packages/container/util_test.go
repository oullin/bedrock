package container_test

import (
	"testing"

	"github.com/bedrock/packages/container"
)

func TestNewReturnsInitializedContainer(t *testing.T) {
	t.Parallel()

	c := container.New()

	if c == nil {
		t.Fatal("expected non-nil container")
	}

	// Should not panic when calling methods on fresh container.
	if c.Bound("anything") {
		t.Fatal("fresh container should have no bindings")
	}
}

func TestContainerCanBeUsedAsDependency(t *testing.T) {
	t.Parallel()

	c := container.New()
	c.Instance("self", c)

	v, err := c.Make("self")

	if err != nil {
		t.Fatal(err)
	}

	if v != c {
		t.Fatal("expected container itself")
	}
}

func TestFlushClearsCallbacks(t *testing.T) {
	t.Parallel()

	c := container.New()

	var called bool

	c.ResolvingAny(func(_ any, _ *container.Container) {
		called = true
	})

	c.Flush()

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Make("name") //nolint:errcheck

	if called {
		t.Fatal("callbacks should be cleared after flush")
	}
}

func TestFlushClearsMethodBindings(t *testing.T) {
	t.Parallel()

	c := container.New()

	c.BindMethod("App@handle", func(_ *container.Container, _ map[string]any) (any, error) {
		return nil, nil
	})

	c.Flush()

	if c.HasMethodBinding("App@handle") {
		t.Fatal("method bindings should be cleared after flush")
	}
}

func TestFlushClearsTags(t *testing.T) {
	t.Parallel()

	c := container.New()

	c.Bind("service", func(_ *container.Container) (any, error) { return "s", nil }, false)
	c.Tag([]string{"service"}, "tag")

	c.Flush()

	results := c.Tagged("tag")

	if len(results) != 0 {
		t.Fatal("tags should be cleared after flush")
	}
}
