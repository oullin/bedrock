package scout_test

import (
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/scout"
)

func TestScoutServiceProviderRegister(t *testing.T) {
	t.Parallel()
	app := container.New()
	config := scout.DefaultConfig()
	config.Driver = "null"

	provider := scout.NewScoutServiceProvider(app, config)
	provider.Register()

	result, err := app.Make("scout")

	if err != nil {
		t.Fatalf("unexpected error resolving scout: %v", err)
	}

	manager, ok := result.(*scout.EngineManager)

	if !ok {
		t.Fatal("expected *scout.EngineManager from container")
	}

	if manager.GetDefaultDriver() != "null" {
		t.Fatalf("expected default driver null, got %s", manager.GetDefaultDriver())
	}
}

func TestScoutServiceProviderRegistersBuiltInEngines(t *testing.T) {
	t.Parallel()
	app := container.New()
	config := scout.DefaultConfig()
	config.Driver = "null"

	provider := scout.NewScoutServiceProvider(app, config)
	provider.Register()

	result, _ := app.Make("scout")
	manager := result.(*scout.EngineManager)

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
	config := scout.DefaultConfig()
	provider := scout.NewScoutServiceProvider(app, config)

	provides := provider.Provides()

	if len(provides) != 1 || provides[0] != "scout" {
		t.Fatalf("expected [scout], got %v", provides)
	}
}

func TestScoutServiceProviderSingleton(t *testing.T) {
	t.Parallel()
	app := container.New()
	config := scout.DefaultConfig()
	config.Driver = "null"

	provider := scout.NewScoutServiceProvider(app, config)
	provider.Register()

	result1, _ := app.Make("scout")
	result2, _ := app.Make("scout")

	if result1 != result2 {
		t.Fatal("expected singleton: same instance on repeated resolution")
	}
}
