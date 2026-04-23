package reverb_test

import (
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/websockets"
)

// ProviderTest::it_retrieves_applications_from_custom_provider
// FactoryTest::it_can_create_a_server
func TestWebSocketsServiceProviderBindsCoreServices(t *testing.T) {
	t.Parallel()

	app := container.New()
	provider := websockets.NewWebSocketsServiceProvider(app, websockets.Config{
		Apps: []websockets.AppConfig{
			{
				ID:             "app-1",
				Key:            "key-1",
				Secret:         "secret-1",
				AllowedOrigins: []string{"*"},
			},
		},
	})

	provider.Register()

	rawApps, err := app.Make("websockets.apps")

	if err != nil {
		t.Fatalf("resolve websockets.apps: %v", err)
	}

	apps, ok := rawApps.(*websockets.AppManager)

	if !ok {
		t.Fatalf("websockets.apps type = %T, want *websockets.AppManager", rawApps)
	}

	if len(apps.All()) != 1 {
		t.Fatalf("resolved apps count = %d, want 1", len(apps.All()))
	}

	rawDispatcher, err := app.Make("websockets.dispatcher")

	if err != nil {
		t.Fatalf("resolve websockets.dispatcher: %v", err)
	}

	if _, ok := rawDispatcher.(*websockets.SyncDispatcher); !ok {
		t.Fatalf("websockets.dispatcher type = %T, want *websockets.SyncDispatcher", rawDispatcher)
	}

	rawServer, err := app.Make("websockets.server")

	if err != nil {
		t.Fatalf("resolve websockets.server: %v", err)
	}

	if _, ok := rawServer.(*websockets.Server); !ok {
		t.Fatalf("websockets.server type = %T, want *websockets.Server", rawServer)
	}

	rawHTTP, err := app.Make("websockets.http")

	if err != nil {
		t.Fatalf("resolve websockets.http: %v", err)
	}

	if _, ok := rawHTTP.(*websockets.HTTPHandler); !ok {
		t.Fatalf("websockets.http type = %T, want *websockets.HTTPHandler", rawHTTP)
	}

	provides := provider.Provides()
	want := []string{
		"websockets.apps",
		"websockets.conns",
		"websockets.channels",
		"websockets.dispatcher",
		"websockets.server",
		"websockets.http",
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
