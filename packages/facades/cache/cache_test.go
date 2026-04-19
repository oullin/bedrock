package cache_test

import (
	"testing"

	"github.com/bedrock/packages/bootstrap"
	cachepkg "github.com/bedrock/packages/cache"
	"github.com/bedrock/packages/container"
	cachefacade "github.com/bedrock/packages/facades/cache"
)

func TestFacade_ResolvesManagerFromGlobalApp(t *testing.T) {
	t.Cleanup(func() {
		bootstrap.SetApp(nil)
		cachefacade.Reset()
	})

	application := container.NewApplication()
	application.Register(cachepkg.NewCacheServiceProvider(application.Container, "array"))
	application.Boot()
	bootstrap.SetApp(application)
	cachefacade.Reset() // forget any cached manager from previous tests

	mgr := cachefacade.Manager()

	if mgr == nil {
		t.Fatal("facade returned nil manager")
	}

	if got := mgr.GetDefaultDriver(); got != "array" {
		t.Fatalf("expected default driver %q, got %q", "array", got)
	}
}

func TestFacade_CachesManagerAcrossCalls(t *testing.T) {
	t.Cleanup(func() {
		bootstrap.SetApp(nil)
		cachefacade.Reset()
	})

	application := container.NewApplication()
	application.Register(cachepkg.NewCacheServiceProvider(application.Container, "file"))
	bootstrap.SetApp(application)
	cachefacade.Reset()

	a := cachefacade.Manager()
	b := cachefacade.Manager()

	if a != b {
		t.Fatal("expected facade to cache the resolved manager")
	}
}

func TestFacade_ResetForcesReResolve(t *testing.T) {
	t.Cleanup(func() {
		bootstrap.SetApp(nil)
		cachefacade.Reset()
	})

	app1 := container.NewApplication()
	app1.Register(cachepkg.NewCacheServiceProvider(app1.Container, "array"))
	bootstrap.SetApp(app1)
	cachefacade.Reset()

	first := cachefacade.Manager()

	// Install a different application and reset the facade cache.
	app2 := container.NewApplication()
	app2.Register(cachepkg.NewCacheServiceProvider(app2.Container, "file"))
	bootstrap.SetApp(app2)
	cachefacade.Reset()

	second := cachefacade.Manager()

	if first == second {
		t.Fatal("expected Reset to drop the cached manager so a fresh one is resolved")
	}

	if got := second.GetDefaultDriver(); got != "file" {
		t.Fatalf("expected new manager to have driver %q, got %q", "file", got)
	}
}
