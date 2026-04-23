package reverb_test

import (
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/reverb"
)

// ProviderTest::it_retrieves_applications_from_custom_provider
// FactoryTest::it_can_create_a_server
func TestReverbServiceProviderBindsCoreServices(t *testing.T) {
	t.Parallel()

	app := container.New()
	provider := reverb.NewReverbServiceProvider(app, reverb.Config{
		Apps: []reverb.AppConfig{
			{
				ID:             "app-1",
				Key:            "key-1",
				Secret:         "secret-1",
				AllowedOrigins: []string{"*"},
			},
		},
	})

	provider.Register()

	rawApps, err := app.Make("reverb.apps")

	if err != nil {
		t.Fatalf("resolve reverb.apps: %v", err)
	}

	apps, ok := rawApps.(*reverb.AppManager)

	if !ok {
		t.Fatalf("reverb.apps type = %T, want *reverb.AppManager", rawApps)
	}

	if len(apps.All()) != 1 {
		t.Fatalf("resolved apps count = %d, want 1", len(apps.All()))
	}

	rawDispatcher, err := app.Make("reverb.dispatcher")

	if err != nil {
		t.Fatalf("resolve reverb.dispatcher: %v", err)
	}

	if _, ok := rawDispatcher.(*reverb.SyncDispatcher); !ok {
		t.Fatalf("reverb.dispatcher type = %T, want *reverb.SyncDispatcher", rawDispatcher)
	}

	rawServer, err := app.Make("reverb.server")

	if err != nil {
		t.Fatalf("resolve reverb.server: %v", err)
	}

	if _, ok := rawServer.(*reverb.Server); !ok {
		t.Fatalf("reverb.server type = %T, want *reverb.Server", rawServer)
	}

	rawHTTP, err := app.Make("reverb.http")

	if err != nil {
		t.Fatalf("resolve reverb.http: %v", err)
	}

	if _, ok := rawHTTP.(*reverb.HTTPHandler); !ok {
		t.Fatalf("reverb.http type = %T, want *reverb.HTTPHandler", rawHTTP)
	}

	provides := provider.Provides()
	want := []string{
		"reverb.apps",
		"reverb.conns",
		"reverb.channels",
		"reverb.dispatcher",
		"reverb.server",
		"reverb.http",
	}

	if len(provides) != len(want) {
		t.Fatalf("Provides length = %d, want %d", len(provides), len(want))
	}

	for i := range want {
		if provides[i] != want[i] {
			t.Fatalf("Provides[%d] = %q, want %q", i, provides[i], want[i])
		}
	}
}
