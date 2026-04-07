package container

import (
	"testing"
)

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

func TestContextualWithFactory(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("db.dsn", "postgres://prod")

	c.When("analytics").Needs("db.dsn").Give(Factory(func(c *Container) (any, error) {
		return "postgres://analytics-replica", nil
	}))

	got, err := c.MakeFor("analytics", "db.dsn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "postgres://analytics-replica" {
		t.Fatalf("want analytics DSN, got %v", got)
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
