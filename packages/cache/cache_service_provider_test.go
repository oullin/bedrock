package cache_test

import (
	"testing"

	"github.com/bedrock/packages/cache"
	"github.com/bedrock/packages/container"
)

func TestCacheServiceProvider_RegistersManagerUnderCacheKey(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	app.Register(cache.NewCacheServiceProvider(app.Container, "array"))
	app.Boot()

	raw, err := app.Make("cache")

	if err != nil {
		t.Fatalf("expected cache to resolve, got error: %v", err)
	}

	mgr, ok := raw.(*cache.Manager)

	if !ok {
		t.Fatalf("expected *cache.Manager, got %T", raw)
	}

	if got := mgr.GetDefaultDriver(); got != "array" {
		t.Fatalf("expected default driver %q, got %q", "array", got)
	}
}

func TestCacheServiceProvider_Provides(t *testing.T) {
	t.Parallel()

	p := cache.NewCacheServiceProvider(container.New(), "array")

	provides := p.Provides()

	if len(provides) != 1 || provides[0] != "cache" {
		t.Fatalf("expected [\"cache\"], got %v", provides)
	}
}

func TestCacheServiceProvider_SingletonSharesInstance(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	app.Register(cache.NewCacheServiceProvider(app.Container, "file"))

	a, _ := app.Make("cache")
	b, _ := app.Make("cache")

	if a != b {
		t.Fatal("expected the cache binding to be a singleton (same pointer across resolves)")
	}
}
