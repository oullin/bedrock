package container_test

import (
	"testing"

	"github.com/bedrock/packages/container"
)

func TestExtendModifiesResolvedInstance(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	v, _ := c.Make("name")

	if v != "Taylor Otwell" {
		t.Fatalf("expected Taylor Otwell, got %v", v)
	}
}

func TestExtendAppliedToExistingInstance(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	v, _ := c.Make("name")

	if v != "Taylor Otwell" {
		t.Fatalf("expected Taylor Otwell, got %v", v)
	}
}

func TestExtendIsLazyInitialized(t *testing.T) {
	t.Parallel()

	c := newContainer()

	called := false

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		called = true

		return instance, nil
	})

	if called {
		t.Fatal("extender should not be called before resolution")
	}

	c.Make("name") //nolint:errcheck

	if !called {
		t.Fatal("extender should be called on resolution")
	}
}

func TestExtendCanBeCalledBeforeBind(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	v, _ := c.Make("name")

	if v != "Taylor Otwell" {
		t.Fatalf("expected Taylor Otwell, got %v", v)
	}
}

func TestMultipleExtendersAppliedInOrder(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "!", nil
	})

	v, _ := c.Make("name")

	if v != "Taylor Otwell!" {
		t.Fatalf("expected Taylor Otwell!, got %v", v)
	}
}

func TestForgetExtenders(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	c.ForgetExtenders("name")

	v, _ := c.Make("name")

	if v != "Taylor" {
		t.Fatalf("expected Taylor after forget, got %v", v)
	}
}

func TestExtendOnSingleton(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Singleton("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	})

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	a, _ := c.Make("name")
	b, _ := c.Make("name")

	if a != "Taylor Otwell" {
		t.Fatalf("expected Taylor Otwell, got %v", a)
	}

	if a != b {
		t.Fatal("singleton extend should return same instance")
	}
}

func TestExtendOnAlias(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Alias("name", "shortName")

	c.Extend("shortName", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	v, _ := c.Make("shortName")

	if v != "Taylor Otwell" {
		t.Fatalf("expected Taylor Otwell, got %v", v)
	}
}

func TestExtendInstanceRebindingCallback(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	var rebound bool

	c.Rebinding("name", func(_ any, _ *container.Container) { //nolint:errcheck
		rebound = true
	})

	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	if !rebound {
		t.Fatal("extend on instance should trigger rebound callback")
	}
}
