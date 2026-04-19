package scout_test

import (
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/search"
)

func TestScoutServiceProviderRegister(t *testing.T) {
	t.Parallel()
	app := container.New()
	config := search.DefaultConfig()
	config.Driver = "null"

	provider := search.NewScoutServiceProvider(app, config)
	provider.Register()

	result, err := app.Make("search")

	if err != nil {
		t.Fatalf("unexpected error resolving search: %v", err)
	}

	manager, ok := result.(*search.EngineManager)

	if !ok {
		t.Fatal("expected *search.EngineManager from container")
	}

	if manager.GetDefaultDriver() != "null" {
		t.Fatalf("expected default driver null, got %s", manager.GetDefaultDriver())
	}
}

func TestScoutServiceProviderRegistersBuiltInEngines(t *testing.T) {
	t.Parallel()
	app := container.New()
	config := search.DefaultConfig()
	config.Driver = "null"

	provider := search.NewScoutServiceProvider(app, config)
	provider.Register()

	result, _ := app.Make("search")
	manager := result.(*search.EngineManager)

	// Check that null engine is available.
	_, err := manager.Engine("null")

	if err != nil {
		t.Fatalf("expected null engine to be registered: %v", err)
	}

	// Check that collection engine is available.
	_, err = manager.Engine("collection")

	if err != nil {
		t.Fatalf("expected collection engine to be registered: %v", err)
	}
}

func TestScoutServiceProviderProvides(t *testing.T) {
	t.Parallel()
	app := container.New()
	config := search.DefaultConfig()
	provider := search.NewScoutServiceProvider(app, config)

	provides := provider.Provides()

	if len(provides) != 1 || provides[0] != "search" {
		t.Fatalf("expected [search], got %v", provides)
	}
}

func TestScoutServiceProviderSingleton(t *testing.T) {
	t.Parallel()
	app := container.New()
	config := search.DefaultConfig()
	config.Driver = "null"

	provider := search.NewScoutServiceProvider(app, config)
	provider.Register()

	result1, _ := app.Make("search")
	result2, _ := app.Make("search")

	if result1 != result2 {
		t.Fatal("expected singleton: same instance on repeated resolution")
	}
}
