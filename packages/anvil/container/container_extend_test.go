package container

import (
	"testing"
)

// Maps to Upstream's ContainerExtendTest. All 11 tests are portable.

func TestExtendedBindings(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "original", nil })

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	got, err := c.Make("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "original-extended" {
		t.Fatalf("want %q, got %q", "original-extended", got)
	}
}

func TestExtendedSingletonBindings(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) { return "singleton", nil })

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	got, _ := c.Make("svc")
	if got != "singleton-extended" {
		t.Fatalf("want %q, got %q", "singleton-extended", got)
	}
}

func TestExtendInstancesArePreserved(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("svc", "original")

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-ext1"
	})

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-ext2"
	})

	got, _ := c.Make("svc")
	if got != "original-ext1-ext2" {
		t.Fatalf("want %q, got %q", "original-ext1-ext2", got)
	}
}

func TestExtendIsLazyInitialized(t *testing.T) {
	t.Parallel()

	c := New()
	extended := false

	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })
	c.Extend("svc", func(value any, _ *Container) any {
		extended = true
		return value
	})

	if extended {
		t.Fatal("extend callback should not fire until Make is called")
	}

	_, _ = c.Make("svc")

	if !extended {
		t.Fatal("extend callback should fire on Make")
	}
}

func TestExtendCanBeCalledBeforeBind(t *testing.T) {
	t.Parallel()

	c := New()

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	c.Bind("svc", func(_ *Container) (any, error) { return "late-bind", nil })

	got, _ := c.Make("svc")
	if got != "late-bind-extended" {
		t.Fatalf("want %q, got %q", "late-bind-extended", got)
	}
}

func TestExtendInstanceRebindingCallback(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("svc", "original")

	var captured any
	c.Rebinding("svc", func(val any) { captured = val })

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	if captured != "original-extended" {
		t.Fatalf("want rebound with %q, got %v", "original-extended", captured)
	}
}

func TestExtendBindRebindingCallback(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })
	_, _ = c.Make("svc")

	var captured any
	c.Rebinding("svc", func(val any) { captured = val })

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	if captured != "value-extended" {
		t.Fatalf("want rebound with %q, got %v", "value-extended", captured)
	}
}

func TestExtensionWorksOnAliasedBindings(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc.full", func(_ *Container) (any, error) { return "value", nil })
	c.Alias("svc.full", "svc")

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	got, _ := c.Make("svc")
	if got != "value-extended" {
		t.Fatalf("want %q, got %q", "value-extended", got)
	}
}

func TestMultipleExtends(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "base", nil })

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-first"
	})

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-second"
	})

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-third"
	})

	got, _ := c.Make("svc")
	if got != "base-first-second-third" {
		t.Fatalf("want %q, got %q", "base-first-second-third", got)
	}
}

func TestUnsetExtend(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	c.ForgetExtenders("svc")

	got, _ := c.Make("svc")
	if got != "value" {
		t.Fatalf("want %q after ForgetExtenders, got %q", "value", got)
	}
}

func TestExtendContextualBinding(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "default", nil })
	c.When("consumer").Needs("svc").Give("contextual")

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	got, _ := c.Make("svc")
	if got != "default-extended" {
		t.Fatalf("want %q, got %q", "default-extended", got)
	}
}

func TestExtendContextualBindingAfterResolution(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "default", nil })
	_, _ = c.Make("svc")

	c.Extend("svc", func(value any, _ *Container) any {
		return value.(string) + "-extended"
	})

	got, _ := c.Make("svc")
	if got != "default-extended" {
		t.Fatalf("want %q, got %q", "default-extended", got)
	}
}

// --- Rebound listeners (moved from container_test.go) ---

func TestReboundListeners(t *testing.T) {
	t.Parallel()

	c := New()

	var captured any
	c.Bind("svc", func(_ *Container) (any, error) { return "v1", nil })
	_, _ = c.Make("svc")

	c.Rebinding("svc", func(val any) { captured = val })

	c.Bind("svc", func(_ *Container) (any, error) { return "v2", nil })

	if captured != "v2" {
		t.Fatalf("want rebound callback with %q, got %v", "v2", captured)
	}
}

func TestReboundListenersOnInstances(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("svc", "v1")
	_, _ = c.Make("svc")

	var captured any
	c.Rebinding("svc", func(val any) { captured = val })

	c.Instance("svc", "v2")

	if captured != "v2" {
		t.Fatalf("want rebound callback with %q, got %v", "v2", captured)
	}
}

func TestReboundListenersOnInstancesOnlyFiresIfWasAlreadyBound(t *testing.T) {
	t.Parallel()

	c := New()

	var fired bool
	c.Rebinding("svc", func(any) { fired = true })

	c.Instance("svc", "v1")

	if fired {
		t.Fatal("rebound should not fire for first Instance")
	}
}

// --- Binding override / precedence (moved from container_test.go) ---

func TestBindingsCanBeOverridden(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("val", func(_ *Container) (any, error) { return "first", nil })
	c.Bind("val", func(_ *Container) (any, error) { return "second", nil })

	got, err := c.Make("val")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "second" {
		t.Fatalf("want %q, got %q", "second", got)
	}
}

func TestSingletonBindingsNotRespectedWithNewBind(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) { return "singleton-v1", nil })
	_, _ = c.Make("svc")

	c.Bind("svc", func(_ *Container) (any, error) { return "transient", nil })

	got, _ := c.Make("svc")
	if got != "transient" {
		t.Fatalf("want %q, got %q", "transient", got)
	}
}

func TestInstanceClearsBinding(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "from-factory", nil })
	c.Instance("svc", "from-instance")

	got, _ := c.Make("svc")
	if got != "from-instance" {
		t.Fatalf("want %q, got %q", "from-instance", got)
	}
}

func TestBindClearsInstance(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("svc", "from-instance")
	c.Bind("svc", func(_ *Container) (any, error) { return "from-factory", nil })

	got, _ := c.Make("svc")
	if got != "from-factory" {
		t.Fatalf("want %q, got %q", "from-factory", got)
	}
}

func TestScopedSingletonWithBind(t *testing.T) {
	t.Parallel()

	c := New()
	c.Scoped("svc", func(_ *Container) (any, error) { return "scoped", nil })

	got, _ := c.Make("svc")
	if got != "scoped" {
		t.Fatalf("want %q, got %q", "scoped", got)
	}

	c.Bind("svc", func(_ *Container) (any, error) { return "transient", nil })

	got, _ = c.Make("svc")
	if got != "transient" {
		t.Fatalf("want %q, got %q", "transient", got)
	}
}

func TestSingletonWithBind(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) { return "singleton", nil })

	got, _ := c.Make("svc")
	if got != "singleton" {
		t.Fatalf("want %q, got %q", "singleton", got)
	}

	c.Bind("svc", func(_ *Container) (any, error) { return "transient", nil })

	got, _ = c.Make("svc")
	if got != "transient" {
		t.Fatalf("want %q, got %q", "transient", got)
	}
}
