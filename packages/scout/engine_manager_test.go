package scout_test

import (
	"testing"

	contract "github.com/bedrock/packages/contracts/scout"
	"github.com/bedrock/packages/scout"
)

func TestEngineManagerRegisterAndResolve(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()
	engine := &fakeEngine{}

	m.Register("null", engine)

	got, err := m.Engine("null")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != engine {
		t.Fatal("expected same engine instance")
	}
}

func TestEngineManagerDefaultDriver(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()
	engine := &fakeEngine{}

	m.Register("database", engine)
	m.SetDefaultDriver("database")

	if m.GetDefaultDriver() != "database" {
		t.Fatalf("expected default driver database, got %s", m.GetDefaultDriver())
	}

	got, err := m.Driver()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != engine {
		t.Fatal("expected same engine instance from Driver()")
	}
}

func TestEngineManagerDefaultDriverViaEngine(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()
	engine := &fakeEngine{}

	m.Register("collection", engine)
	m.SetDefaultDriver("collection")

	// Calling Engine() with no args uses default.
	got, err := m.Engine()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != engine {
		t.Fatal("expected default engine")
	}
}

func TestEngineManagerUnregistered(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()

	_, err := m.Engine("nonexistent")

	if err == nil {
		t.Fatal("expected error for unregistered engine")
	}
}

func TestEngineManagerNoDefault(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()

	_, err := m.Engine()

	if err == nil {
		t.Fatal("expected error when no default driver is set")
	}
}

func TestEngineManagerExtendAndBuild(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()

	m.Extend("custom", func(_ map[string]any) (contract.Engine, error) {
		return &fakeEngine{}, nil
	})

	engine, err := m.Build("custom", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine == nil {
		t.Fatal("expected engine from Build")
	}
}

func TestEngineManagerBuildUnsupported(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()

	_, err := m.Build("unknown", nil)

	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}
}

func TestEngineManagerBuildAndRegister(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()

	m.Extend("custom", func(_ map[string]any) (contract.Engine, error) {
		return &fakeEngine{}, nil
	})

	engine, err := m.BuildAndRegister("my_engine", "custom", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := m.Engine("my_engine")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != engine {
		t.Fatal("expected cached engine after BuildAndRegister")
	}
}

func TestEngineManagerPurge(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()
	m.Register("null", &fakeEngine{})

	m.Purge("null")

	_, err := m.Engine("null")

	if err == nil {
		t.Fatal("expected error after purge")
	}
}

func TestEngineManagerForgetDriver(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()
	m.Register("null", &fakeEngine{})

	m.ForgetDriver("null")

	_, err := m.Engine("null")

	if err == nil {
		t.Fatal("expected error after ForgetDriver")
	}
}

func TestEngineManagerGetEngines(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()
	m.Register("a", &fakeEngine{})
	m.Register("b", &fakeEngine{})

	names := m.GetEngines()

	if len(names) != 2 {
		t.Fatalf("expected 2 engine names, got %d", len(names))
	}
}

func TestEngineManagerGetDrivers(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()
	m.Extend("x", func(_ map[string]any) (contract.Engine, error) {
		return &fakeEngine{}, nil
	})
	m.Extend("y", func(_ map[string]any) (contract.Engine, error) {
		return &fakeEngine{}, nil
	})

	names := m.GetDrivers()

	if len(names) != 2 {
		t.Fatalf("expected 2 driver names, got %d", len(names))
	}
}

func TestEngineManagerConcurrentAccess(t *testing.T) {
	t.Parallel()
	m := scout.NewEngineManager()
	m.Register("test", &fakeEngine{})
	m.SetDefaultDriver("test")

	done := make(chan struct{})

	for i := 0; i < 100; i++ {
		go func() {
			defer func() { done <- struct{}{} }()

			_, _ = m.Engine("test")
			m.Register("test", &fakeEngine{})
			_ = m.GetEngines()
			_ = m.GetDrivers()
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}
