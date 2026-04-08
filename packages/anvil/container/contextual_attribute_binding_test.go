package container

import (
	"errors"
	"testing"
)

// Upstream's ContextualAttributeBindingTest relies on PHP attributes
// (#[Attribute]) for dependency injection. All 20 Upstream tests are skipped.
// In Go, the closest analog to attribute-based contextual binding is passing
// a Factory to Give(), which provides dynamic resolution.
//
// INTENTIONAL-SKIP: testDependencyCanBeResolvedFromAttributeBinding
// INTENTIONAL-SKIP: testSimpleDependencyCanBeResolvedCorrectlyFromGiveAttributeBinding
// INTENTIONAL-SKIP: testComplexDependencyCanBeResolvedCorrectlyFromGiveAttributeBinding
// INTENTIONAL-SKIP: testScalarDependencyCanBeResolvedFromAttributeBinding
// INTENTIONAL-SKIP: testScalarDependencyCanBeResolvedFromAttributeResolveMethod
// INTENTIONAL-SKIP: testDependencyWithAfterCallbackAttributeCanBeResolved
// INTENTIONAL-SKIP: testAuthedAttribute
// INTENTIONAL-SKIP: testCacheAttribute
// INTENTIONAL-SKIP: testConfigAttribute
// INTENTIONAL-SKIP: testDatabaseAttribute
// INTENTIONAL-SKIP: testAuthAttribute
// INTENTIONAL-SKIP: testLogAttribute
// INTENTIONAL-SKIP: testRouteParameterAttribute
// INTENTIONAL-SKIP: testContextAttribute
// INTENTIONAL-SKIP: testContextAttributeInteractingWithHidden
// INTENTIONAL-SKIP: testStorageAttribute
// INTENTIONAL-SKIP: testInjectionWithAttributeOnAppCall
// INTENTIONAL-SKIP: testAttributeOnAppCall
// INTENTIONAL-SKIP: testNestedAttributeOnAppCall
// INTENTIONAL-SKIP: testTagAttribute

// --- Go equivalent: Factory-based contextual bindings ---

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

func TestContextualFactoryResolvingDeps(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("config.region", "us-east-1")

	c.When("storage").Needs("endpoint").Give(Factory(func(c *Container) (any, error) {
		region, err := c.Make("config.region")
		if err != nil {
			return nil, err
		}
		return "https://s3." + region.(string) + ".amazonaws.com", nil
	}))

	got, err := c.MakeFor("storage", "endpoint")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "https://s3.us-east-1.amazonaws.com" {
		t.Fatalf("want S3 endpoint, got %v", got)
	}
}

func TestContextualFactoryError(t *testing.T) {
	t.Parallel()

	c := New()
	c.When("consumer").Needs("broken").Give(Factory(func(_ *Container) (any, error) {
		return nil, errors.New("contextual factory failed")
	}))

	_, err := c.MakeFor("consumer", "broken")
	if err == nil {
		t.Fatal("expected error from contextual factory")
	}

	if !errors.Is(err, ErrResolve) {
		t.Fatalf("expected ErrResolve, got: %v", err)
	}
}

func TestContextualMultipleConcretesSameAbstract(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("cache", func(_ *Container) (any, error) { return "default-cache", nil })

	c.When("api").Needs("cache").Give(Factory(func(_ *Container) (any, error) {
		return "redis-cache", nil
	}))

	c.When("worker").Needs("cache").Give(Factory(func(_ *Container) (any, error) {
		return "sqs-cache", nil
	}))

	apiGot, _ := c.MakeFor("api", "cache")
	workerGot, _ := c.MakeFor("worker", "cache")
	defaultGot, _ := c.Make("cache")

	if apiGot != "redis-cache" {
		t.Fatalf("api: want %q, got %v", "redis-cache", apiGot)
	}

	if workerGot != "sqs-cache" {
		t.Fatalf("worker: want %q, got %v", "sqs-cache", workerGot)
	}

	if defaultGot != "default-cache" {
		t.Fatalf("default: want %q, got %v", "default-cache", defaultGot)
	}
}
