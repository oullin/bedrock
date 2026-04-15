package boost_test

import (
	"testing"

	"github.com/bedrock/packages/boost"
	"github.com/bedrock/packages/container"
)

// TestBoostServiceProviderRegistersKey mirrors BoostServiceProviderTest.
func TestBoostServiceProviderRegistersKey(t *testing.T) {
	t.Parallel()

	app := container.New()
	p := boost.NewBoostServiceProvider(app)
	p.Register()

	// "boost" must be resolvable from the container.
	val, err := app.Make("boost")
	if err != nil {
		t.Fatalf("container.Make(\"boost\"): %v", err)
	}

	m, ok := val.(*boost.Manager)
	if !ok {
		t.Fatalf("resolved value is %T, want *boost.Manager", val)
	}

	if m == nil {
		t.Fatal("resolved Manager must not be nil")
	}
}

// TestBoostServiceProviderProvides ensures Provides() returns the "boost" key.
func TestBoostServiceProviderProvides(t *testing.T) {
	t.Parallel()

	app := container.New()
	p := boost.NewBoostServiceProvider(app)
	provides := p.Provides()

	if len(provides) != 1 || provides[0] != "boost" {
		t.Errorf("Provides() = %v, want [\"boost\"]", provides)
	}
}

// TestBoostServiceProviderSingleton verifies two resolutions return the same instance.
func TestBoostServiceProviderSingleton(t *testing.T) {
	t.Parallel()

	app := container.New()
	boost.NewBoostServiceProvider(app).Register()

	v1, _ := app.Make("boost")
	v2, _ := app.Make("boost")

	if v1 != v2 {
		t.Error("singleton: two resolutions returned different instances")
	}
}
