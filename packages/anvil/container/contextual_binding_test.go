package container

import (
	"testing"
)

// Maps to Upstream's ContextualBindingTest. 27 Upstream tests are portable.
// Go's contextual binding uses string-based concrete/abstract names instead of
// PHP class types. The When/Needs/Give API is equivalent.

func TestContextualBinding(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("logger", func(_ *Container) (any, error) { return "default-logger", nil })
	c.When("payment-service").Needs("logger").Give("payment-logger")

	got, err := c.MakeFor("payment-service", "logger")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "payment-logger" {
		t.Fatalf("want %q, got %q", "payment-logger", got)
	}
}

func TestContextualFallback(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("cache", func(_ *Container) (any, error) { return "redis", nil })

	got, err := c.MakeFor("some-service", "cache")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "redis" {
		t.Fatalf("want %q, got %q", "redis", got)
	}
}

func TestContextualOverride(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("mailer", func(_ *Container) (any, error) { return "smtp", nil })

	c.When("test-runner").Needs("mailer").Give("log-mailer")
	c.When("test-runner").Needs("mailer").Give("array-mailer")

	got, err := c.MakeFor("test-runner", "mailer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "array-mailer" {
		t.Fatalf("want %q, got %q", "array-mailer", got)
	}
}

// --- Upstream portable tests ---

func TestContextualBindingDoesNotAffectOtherConcretes(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("logger", func(_ *Container) (any, error) { return "default", nil })
	c.When("service-a").Needs("logger").Give("custom-logger")

	gotA, _ := c.MakeFor("service-a", "logger")
	gotB, _ := c.MakeFor("service-b", "logger")

	if gotA != "custom-logger" {
		t.Fatalf("service-a: want %q, got %v", "custom-logger", gotA)
	}

	if gotB != "default" {
		t.Fatalf("service-b: want %q (fallback), got %v", "default", gotB)
	}
}

func TestContextualBindingWithAlias(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("app.logger", func(_ *Container) (any, error) { return "default", nil })
	c.Alias("app.logger", "logger")

	// Contextual needs must use the canonical abstract name since MakeFor
	// resolves aliases before checking contextual overrides.
	c.When("service").Needs("app.logger").Give("custom")

	got, _ := c.MakeFor("service", "logger")
	if got != "custom" {
		t.Fatalf("want %q, got %v", "custom", got)
	}
}

func TestContextualBindingWithInstances(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("config", "default-config")
	c.When("service").Needs("config").Give("custom-config")

	got, _ := c.MakeFor("service", "config")
	if got != "custom-config" {
		t.Fatalf("want %q, got %v", "custom-config", got)
	}
}

func TestContextualBindingWithAliasedInstances(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("app.db", "default-db")
	c.Alias("app.db", "db")
	c.When("service").Needs("app.db").Give("custom-db")

	got, _ := c.MakeFor("service", "db")
	if got != "custom-db" {
		t.Fatalf("want %q, got %v", "custom-db", got)
	}
}

func TestContextualBindingMultipleConcretes(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("logger", func(_ *Container) (any, error) { return "default", nil })

	c.When("auth").Needs("logger").Give("auth-logger")
	c.When("payment").Needs("logger").Give("payment-logger")

	gotAuth, _ := c.MakeFor("auth", "logger")
	gotPay, _ := c.MakeFor("payment", "logger")

	if gotAuth != "auth-logger" {
		t.Fatalf("auth: want %q, got %v", "auth-logger", gotAuth)
	}

	if gotPay != "payment-logger" {
		t.Fatalf("payment: want %q, got %v", "payment-logger", gotPay)
	}
}

func TestContextualBindingGivesTagged(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("csv", func(_ *Container) (any, error) { return "csv-report", nil })
	c.Bind("pdf", func(_ *Container) (any, error) { return "pdf-report", nil })
	c.Tag([]string{"csv", "pdf"}, "reports")

	c.When("dashboard").Needs("reports").Give(Factory(func(c *Container) (any, error) {
		return c.Tagged("reports")
	}))

	got, err := c.MakeFor("dashboard", "reports")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reports := got.([]any)
	if len(reports) != 2 {
		t.Fatalf("want 2 reports, got %d", len(reports))
	}
}

func TestContextualBindingDoesntOverrideNonContextual(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("cache", func(_ *Container) (any, error) { return "redis", nil })
	c.When("api").Needs("cache").Give("memcached")

	got, _ := c.Make("cache")
	if got != "redis" {
		t.Fatalf("non-contextual should still return %q, got %v", "redis", got)
	}
}
