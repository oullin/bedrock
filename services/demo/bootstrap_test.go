package demo_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/services/demo/api"
)

func standardOptions(t *testing.T) api.Options {
	t.Helper()

	return api.Options{
		BasePath:      t.TempDir(),
		StoragePath:   filepath.Join(t.TempDir(), "storage"),
		DatabaseURL:   "sqlite:///:memory:",
		RunMigrations: ptr(true),
		Seed:          ptr(true),
	}
}

func ptr(value bool) *bool {
	return &value
}

func assertStandardBindings(t *testing.T, application *container.Application) {
	t.Helper()

	standardKeys := []string{
		"events", "hash", "files", "cookie", "validator", "concurrency",
		"cache", "session", "queue", "log", "auth",
		"bus", "notifications", "db", "router",
	}

	for _, key := range standardKeys {
		v, err := application.Make(key)

		if err != nil {
			t.Errorf("Make(%q) failed: %v", key, err)

			continue
		}

		if v == nil {
			t.Errorf("Make(%q) returned nil", key)
		}
	}
}

func TestNewApplication_RegistersAndBootsAllStandardProviders(t *testing.T) {
	t.Parallel()

	application := api.NewApplication(standardOptions(t))

	if !application.Booted() {
		t.Fatal("expected app to be booted")
	}

	assertStandardBindings(t, application)
}

func TestStandardProviders_ManualCompositionBootsAllStandardProviders(t *testing.T) {
	t.Parallel()

	application := container.NewApplication()
	application.RegisterMany(api.StandardProviders(application))
	application.Boot()

	if !application.Booted() {
		t.Fatal("expected app to be booted")
	}

	assertStandardBindings(t, application)
}

func TestNewApplication_EncryptionSkippedWithoutKey(t *testing.T) {
	t.Parallel()

	application := api.NewApplication(standardOptions(t))

	_, err := application.Make("encrypter")

	if !errors.Is(err, container.ErrNotBound) {
		t.Fatalf("expected ErrNotBound for encrypter without key, got %v", err)
	}
}

func TestNewApplication_EncryptionRegisteredWhenKeyProvided(t *testing.T) {
	t.Parallel()

	opts := standardOptions(t)
	opts.EncryptionKey = make([]byte, 32)

	application := api.NewApplication(opts)

	v, err := application.Make("encrypter")

	if err != nil {
		t.Fatalf("expected encrypter to resolve when key provided, got %v", err)
	}

	if v == nil {
		t.Fatal("encrypter is nil")
	}
}

func TestNewApplication_OptionsOverrideDefaults(t *testing.T) {
	t.Parallel()

	opts := standardOptions(t)
	opts.CacheDefaultDriver = "redis"

	application := api.NewApplication(opts)

	type defaultDriverGetter interface{ GetDefaultDriver() string }

	raw, err := application.Make("cache")

	if err != nil {
		t.Fatal(err)
	}

	mgr, ok := raw.(defaultDriverGetter)

	if !ok {
		t.Fatalf("cache binding does not expose GetDefaultDriver: %T", raw)
	}

	if got := mgr.GetDefaultDriver(); got != "redis" {
		t.Fatalf("expected cache default driver %q, got %q", "redis", got)
	}
}
