package container_test

import (
	"testing"

	"github.com/bedrock/packages/container"
)

func TestContextualBindingResolvesCorrectImplementation(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("logger", func(_ *container.Container) (any, error) {
		return "default-logger", nil
	}, false)

	c.When("service-a").Needs("logger").Give(container.Factory(func(_ *container.Container) (any, error) {
		return "special-logger", nil
	}))

	c.Bind("service-a", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)

	c.Bind("service-b", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)

	a, _ := c.Make("service-a")
	b, _ := c.Make("service-b")

	if a != "special-logger" {
		t.Fatalf("expected special-logger for service-a, got %v", a)
	}

	if b != "default-logger" {
		t.Fatalf("expected default-logger for service-b, got %v", b)
	}
}

func TestContextualBindingWithRawValue(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.When("service").Needs("$apiKey").Give("secret-key-123")

	c.Bind("service", func(cc *container.Container) (any, error) {
		return cc.Make("$apiKey")
	}, false)

	v, _ := c.Make("service")

	if v != "secret-key-123" {
		t.Fatalf("expected secret-key-123, got %v", v)
	}
}

func TestContextualBindingForMultipleConcretes(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("logger", func(_ *container.Container) (any, error) {
		return "default-logger", nil
	}, false)

	c.When("service-a", "service-b").Needs("logger").Give(container.Factory(func(_ *container.Container) (any, error) {
		return "shared-logger", nil
	}))

	c.Bind("service-a", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)

	c.Bind("service-b", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)

	c.Bind("service-c", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)

	a, _ := c.Make("service-a")
	b, _ := c.Make("service-b")
	d, _ := c.Make("service-c")

	if a != "shared-logger" {
		t.Fatalf("expected shared-logger for service-a, got %v", a)
	}

	if b != "shared-logger" {
		t.Fatalf("expected shared-logger for service-b, got %v", b)
	}

	if d != "default-logger" {
		t.Fatalf("expected default-logger for service-c, got %v", d)
	}
}

func TestContextualBindingDoesNotOverrideNonContextualResolution(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("logger", func(_ *container.Container) (any, error) {
		return "default-logger", nil
	}, false)

	c.When("service-a").Needs("logger").Give(container.Factory(func(_ *container.Container) (any, error) {
		return "special-logger", nil
	}))

	// Direct resolution (not through service-a) should get the default.
	v, _ := c.Make("logger")

	if v != "default-logger" {
		t.Fatalf("expected default-logger, got %v", v)
	}
}

func TestContextualBindingGiveTagged(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("report-a", func(_ *container.Container) (any, error) {
		return "report-a", nil
	}, false)

	c.Bind("report-b", func(_ *container.Container) (any, error) {
		return "report-b", nil
	}, false)

	c.Tag([]string{"report-a", "report-b"}, "reports")

	c.When("aggregator").Needs("reports").GiveTagged("reports")

	c.Bind("aggregator", func(cc *container.Container) (any, error) {
		return cc.Make("reports")
	}, false)

	v, err := c.Make("aggregator")

	if err != nil {
		t.Fatal(err)
	}

	reports, ok := v.([]any)

	if !ok {
		t.Fatalf("expected []any, got %T", v)
	}

	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}

	if reports[0] != "report-a" || reports[1] != "report-b" {
		t.Fatalf("unexpected report values: %v", reports)
	}
}

func TestContextualBindingGiveConfig(t *testing.T) {
	t.Parallel()

	c := newContainer()

	// Register a mock config.
	type mockConfig struct{}
	c.Instance("config", &mockConfig{})

	// GiveConfig needs a Get method.
	// Since mockConfig doesn't implement configGetter, it should use the fallback.
	c.When("service").Needs("$timeout").GiveConfig("app.timeout", 30)

	c.Bind("service", func(cc *container.Container) (any, error) {
		return cc.Make("$timeout")
	}, false)

	v, _ := c.Make("service")

	if v != 30 {
		t.Fatalf("expected 30, got %v", v)
	}
}

func TestContextualBindingGiveConfigWithGetter(t *testing.T) {
	t.Parallel()

	c := newContainer()

	type mockConfig struct{}
	mc := &struct {
		mockConfig
	}{}

	// Create a config that implements Get.
	cfg := &configStub{data: map[string]any{
		"app.timeout": 60,
	}}
	c.Instance("config", cfg)

	_ = mc // unused but shows intent

	c.When("service").Needs("$timeout").GiveConfig("app.timeout", 30)

	c.Bind("service", func(cc *container.Container) (any, error) {
		return cc.Make("$timeout")
	}, false)

	v, _ := c.Make("service")

	if v != 60 {
		t.Fatalf("expected 60, got %v", v)
	}
}

type configStub struct {
	data map[string]any
}

func (cs *configStub) Get(key string, fallback ...any) any {
	if v, ok := cs.data[key]; ok {
		return v
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return nil
}

func TestContextualBindingWorksForExistingInstancedBindings(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("logger", "instanced-logger")

	c.When("service").Needs("logger").Give(container.Factory(func(_ *container.Container) (any, error) {
		return "contextual-logger", nil
	}))

	c.Bind("service", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)

	v, _ := c.Make("service")

	if v != "contextual-logger" {
		t.Fatalf("expected contextual-logger, got %v", v)
	}
}

func TestContextualBindingWorksWithAliasedTargets(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("logger", func(_ *container.Container) (any, error) {
		return "default", nil
	}, false)

	c.Alias("logger", "log")

	c.When("service").Needs("log").Give(container.Factory(func(_ *container.Container) (any, error) {
		return "contextual", nil
	}))

	c.Bind("service", func(cc *container.Container) (any, error) {
		return cc.Make("log")
	}, false)

	v, _ := c.Make("service")

	if v != "contextual" {
		t.Fatalf("expected contextual, got %v", v)
	}
}

func TestContextualBindingNotRecreatedUnnecessarily(t *testing.T) {
	t.Parallel()

	c := newContainer()

	count := 0

	c.When("service").Needs("dep").Give(container.Factory(func(_ *container.Container) (any, error) {
		count++

		return count, nil
	}))

	c.Bind("service", func(cc *container.Container) (any, error) {
		return cc.Make("dep")
	}, false)

	a, _ := c.Make("service")
	b, _ := c.Make("service")

	// Each resolution should call the factory.
	if a == b {
		t.Fatal("expected different instances for non-shared contextual binding")
	}
}

func TestAddContextualBindingDirectly(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.AddContextualBinding("service", "dep", container.Factory(func(_ *container.Container) (any, error) {
		return "direct-contextual", nil
	}))

	c.Bind("service", func(cc *container.Container) (any, error) {
		return cc.Make("dep")
	}, false)

	v, _ := c.Make("service")

	if v != "direct-contextual" {
		t.Fatalf("expected direct-contextual, got %v", v)
	}
}
