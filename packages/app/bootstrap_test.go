package app_test

import (
	"errors"
	"testing"

	bedrockapp "github.com/bedrock/app"
	"github.com/bedrock/packages/container"
)

func TestDefault_RegistersAndBootsAllStandardProviders(t *testing.T) {
	t.Parallel()

	application := bedrockapp.Default()

	if !application.Booted() {
		t.Fatal("expected app to be booted")
	}

	// Every standard binding should resolve.
	standardKeys := []string{
		"events", "hash", "files", "cookie", "validator", "concurrency",
		"cache", "session", "queue", "log", "auth",
		"bus", "notifications", "router",
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

func TestDefault_EncryptionSkippedWithoutKey(t *testing.T) {
	t.Parallel()

	application := bedrockapp.Default()

	_, err := application.Make("encrypter")

	if !errors.Is(err, container.ErrNotBound) {
		t.Fatalf("expected ErrNotBound for encrypter without key, got %v", err)
	}
}

func TestDefault_EncryptionRegisteredWhenKeyProvided(t *testing.T) {
	t.Parallel()

	application := bedrockapp.Default(bedrockapp.Options{
		EncryptionKey: make([]byte, 32),
	})

	v, err := application.Make("encrypter")

	if err != nil {
		t.Fatalf("expected encrypter to resolve when key provided, got %v", err)
	}

	if v == nil {
		t.Fatal("encrypter is nil")
	}
}

func TestDefault_OptionsOverrideDefaults(t *testing.T) {
	t.Parallel()

	application := bedrockapp.Default(bedrockapp.Options{
		CacheDefaultDriver: "redis",
	})

	// Resolve the cache and check the default driver was applied.
	// We can't import cache here without circular imports through bootstrap's
	// own deps, so resolve as any and rely on a method check.
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
