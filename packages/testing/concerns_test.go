package testing_test

import (
	"os"
	"strings"
	"testing"

	packagetesting "github.com/bedrock/packages/testing"
)

func TestViewConcernHelpers(t *testing.T) {
	t.Parallel()

	t.Run("TestViewsTest::testCompiledViewPathAppendsToken", func(t *testing.T) {
		got := packagetesting.CompiledViewPath("/tmp/views", "abc123")

		if got != "/tmp/views/abc123" {
			t.Fatalf("expected compiled view path to append token, got %q", got)
		}
	})

	t.Run("TestViewsTest::testCompiledViewPathTrimsTrailingSlash", func(t *testing.T) {
		got := packagetesting.CompiledViewPath("/tmp/views/", "abc123")

		if got != "/tmp/views/abc123" {
			t.Fatalf("expected trailing slash to be trimmed, got %q", got)
		}
	})

	t.Run("TestViewsTest::testCompiledViewPathReturnsNullWhenEmpty", func(t *testing.T) {
		if got := packagetesting.CompiledViewPath("", "abc123"); got != "" {
			t.Fatalf("expected empty compiled view path, got %q", got)
		}
	})

	t.Run("TestViewsTest::testSwitchToCompiledViewPathUpdatesConfig", func(t *testing.T) {
		config := map[string]any{"view": map[string]any{"compiled": "/tmp/views"}}
		updated := packagetesting.SwitchToCompiledViewPath(config, "abc123")

		view := updated["view"].(map[string]any)

		if view["compiled"] != "/tmp/views/abc123" {
			t.Fatalf("expected compiled view path to update, got %v", view["compiled"])
		}

		original := config["view"].(map[string]any)

		if original["compiled"] != "/tmp/views" {
			t.Fatalf("expected original config to remain unchanged, got %v", original["compiled"])
		}
	})

	t.Run("TestViewsTest::testSwitchToCompiledViewPathUpdatesCompilerCachePath", func(t *testing.T) {
		config := map[string]any{"view": map[string]any{"compiled": "/tmp/views"}}
		updated := packagetesting.SwitchToCompiledViewPath(config, "xyz789")

		view := updated["view"].(map[string]any)

		if view["compiled"] != "/tmp/views/xyz789" {
			t.Fatalf("expected compiler cache path to update, got %v", view["compiled"])
		}
	})

	t.Run("TestViewsTest::testCompiledViewPathWithDifferentToken", func(t *testing.T) {
		got := packagetesting.CompiledViewPath("/tmp/views", "xyz789")

		if got != "/tmp/views/xyz789" {
			t.Fatalf("expected compiled view path to use the new token, got %q", got)
		}
	})

	t.Run("TestViewsTest::testCompiledViewPath", func(t *testing.T) {
		got := packagetesting.CompiledViewPath("/tmp/views", "token")

		if got != "/tmp/views/token" {
			t.Fatalf("expected compiled view path to resolve, got %q", got)
		}
	})

	t.Run("TestViewsTest::testTearDownProcessDeletesCompiledViewDirectory", func(t *testing.T) {
		dir := t.TempDir()
		file := dir + "/compiled.php"

		if err := os.WriteFile(file, []byte("<?php"), 0o600); err != nil {
			t.Fatalf("failed to create compiled file: %v", err)
		}

		if err := packagetesting.DeleteCompiledViewDirectory(dir); err != nil {
			t.Fatalf("expected compiled view directory removal to succeed, got %v", err)
		}

		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Fatalf("expected compiled view directory to be deleted, stat err=%v", err)
		}
	})
}

func TestCacheConcernHelpers(t *testing.T) {
	t.Parallel()

	t.Run("TestCachesTest::testCachePrefixAppendsToken", func(t *testing.T) {
		if got := packagetesting.CachePrefix("cache-prefix-", "abc123"); got != "cache-prefix-abc123" {
			t.Fatalf("expected cache prefix to append token, got %q", got)
		}
	})

	t.Run("TestCachesTest::testCachePrefixPreservesOriginalPrefix", func(t *testing.T) {
		if got := packagetesting.CachePrefix("cache-prefix-", ""); got != "cache-prefix-" {
			t.Fatalf("expected original cache prefix to be preserved, got %q", got)
		}
	})

	t.Run("TestCachesTest::testSwitchToCachePrefixUpdatesConfig", func(t *testing.T) {
		config := map[string]any{"cache": map[string]any{"prefix": "cache-prefix-"}}
		updated := packagetesting.SwitchToCachePrefix(config, "abc123")

		cache := updated["cache"].(map[string]any)

		if cache["prefix"] != "cache-prefix-abc123" {
			t.Fatalf("expected cache prefix to update, got %v", cache["prefix"])
		}
	})

	t.Run("TestCachesTest::testSwitchToCachePrefixDoesNotRemoveResolvedDrivers", func(t *testing.T) {
		config := map[string]any{
			"cache": map[string]any{
				"prefix":          "cache-prefix-",
				"resolvedDrivers": map[string]any{"redis": true},
			},
		}
		updated := packagetesting.SwitchToCachePrefix(config, "abc123")

		cache := updated["cache"].(map[string]any)
		resolved := cache["resolvedDrivers"].(map[string]any)

		if resolved["redis"] != true {
			t.Fatalf("expected resolved drivers to remain intact, got %v", resolved)
		}
	})

	t.Run("TestCachesTest::testBootTestCacheRegistersSetUpTestCaseCallback", func(t *testing.T) {
		state := packagetesting.NewParallelTestingState("abc123")
		called := 0

		packagetesting.BootTestCache(state, false, func() {
			called++
		})

		state.RunCallbacks()

		if called != 1 {
			t.Fatalf("expected boot test cache callback to run once, got %d", called)
		}
	})

	t.Run("TestCachesTest::testBootTestCacheSkipsIsolationIfOptedOut", func(t *testing.T) {
		state := packagetesting.NewParallelTestingState("abc123")
		called := 0

		packagetesting.BootTestCache(state, true, func() {
			called++
		})

		state.RunCallbacks()

		if called != 0 {
			t.Fatalf("expected boot test cache callback to be skipped, got %d", called)
		}
	})
}

func TestDatabaseAndConfigHelpers(t *testing.T) {
	t.Parallel()

	t.Run("TestDatabasesTest::testSwitchToDatabaseWithoutUrl", func(t *testing.T) {
		config := map[string]any{"database": map[string]any{"url": ""}}
		updated := packagetesting.SwitchToDatabase(config, "")

		database := updated["database"].(map[string]any)

		if database["url"] != "" {
			t.Fatalf("expected empty database URL, got %v", database["url"])
		}
	})

	t.Run("TestDatabasesTest::testSwitchToDatabaseWithUrl", func(t *testing.T) {
		config := map[string]any{"database": map[string]any{"url": ""}}
		updated := packagetesting.SwitchToDatabase(config, "sqlite://:memory:")

		database := updated["database"].(map[string]any)

		if database["url"] != "sqlite://:memory:" {
			t.Fatalf("expected database URL to update, got %v", database["url"])
		}
	})

	t.Run("ConfigShowCommandTest::testDisplayConfig", func(t *testing.T) {
		config := map[string]any{
			"app":   "bedrock",
			"cache": map[string]any{"prefix": "cache-prefix-"},
		}

		got, err := packagetesting.RenderConfigShow(config, "")

		if err != nil {
			t.Fatalf("expected config render to succeed, got %v", err)
		}

		if !strings.Contains(got, "app=bedrock") || !strings.Contains(got, "cache.prefix=cache-prefix-") {
			t.Fatalf("unexpected config render: %q", got)
		}
	})

	t.Run("ConfigShowCommandTest::testDisplayNestedConfigItems", func(t *testing.T) {
		config := map[string]any{
			"app":   "bedrock",
			"cache": map[string]any{"prefix": "cache-prefix-"},
		}

		got, err := packagetesting.RenderConfigShow(config, "cache")

		if err != nil {
			t.Fatalf("expected nested config render to succeed, got %v", err)
		}

		if got != "cache.prefix=cache-prefix-" {
			t.Fatalf("expected nested config items to render deterministically, got %q", got)
		}
	})

	t.Run("ConfigShowCommandTest::testDisplaySingleValue", func(t *testing.T) {
		config := map[string]any{"app": "bedrock"}

		value, ok := packagetesting.ResolveConfigValue(config, "app")

		if !ok || value != "bedrock" {
			t.Fatalf("expected single config value, got %v (ok=%v)", value, ok)
		}
	})

	t.Run("ConfigShowCommandTest::testDisplayErrorIfConfigDoesNotExist", func(t *testing.T) {
		config := map[string]any{"app": "bedrock"}

		if _, err := packagetesting.RenderConfigShow(config, "missing"); err == nil {
			t.Fatal("expected missing config key to return an error")
		}
	})

	t.Run("Concerns/InteractsWithDatabaseTest::testCastToJsonSqlite", func(t *testing.T) {
		if got, ok := packagetesting.CastToJSONType("sqlite"); !ok || got != "text" {
			t.Fatalf("expected sqlite JSON cast type text, got %q (ok=%v)", got, ok)
		}
	})

	t.Run("Concerns/InteractsWithDatabaseTest::testCastToJsonPostgres", func(t *testing.T) {
		if got, ok := packagetesting.CastToJSONType("postgres"); !ok || got != "jsonb" {
			t.Fatalf("expected postgres JSON cast type jsonb, got %q (ok=%v)", got, ok)
		}
	})

	t.Run("Concerns/InteractsWithDatabaseTest::testCastToJsonMySql", func(t *testing.T) {
		if got, ok := packagetesting.CastToJSONType("mysql"); !ok || got != "json" {
			t.Fatalf("expected mysql JSON cast type json, got %q (ok=%v)", got, ok)
		}
	})

	t.Run("Concerns/InteractsWithDatabaseTest::testCastToJsonMariaDb", func(t *testing.T) {
		if got, ok := packagetesting.CastToJSONType("mariadb"); !ok || got != "json" {
			t.Fatalf("expected mariadb JSON cast type json, got %q (ok=%v)", got, ok)
		}
	})
}

func TestOperationalHelpers(t *testing.T) {
	t.Parallel()

	t.Run("InteractsWithDeprecationHandlingTest::testWithDeprecationHandling", func(t *testing.T) {
		reporter := packagetesting.NewDeprecationReporter(true)
		reporter.Warn("deprecated")

		if got := reporter.Messages(); len(got) != 1 || got[0] != "deprecated" {
			t.Fatalf("expected deprecation warning to be recorded, got %v", got)
		}
	})

	t.Run("InteractsWithDeprecationHandlingTest::testWithoutDeprecationHandling", func(t *testing.T) {
		reporter := packagetesting.NewDeprecationReporter(false)
		reporter.Warn("deprecated")

		if got := reporter.Messages(); len(got) != 0 {
			t.Fatalf("expected deprecation warning to be suppressed, got %v", got)
		}
	})

	t.Run("ParallelTestingTest::testToken", func(t *testing.T) {
		state := packagetesting.NewParallelTestingState("abc123")

		if got := state.Token(); got != "abc123" {
			t.Fatalf("expected parallel token, got %q", got)
		}
	})

	t.Run("ParallelTestingTest::testOptions", func(t *testing.T) {
		state := packagetesting.NewParallelTestingState("abc123")
		state.SetOption("cache", "redis")

		options := state.Options()

		if options["cache"] != "redis" {
			t.Fatalf("expected parallel testing option to be stored, got %v", options)
		}
	})

	t.Run("ParallelTestingTest::testCallbacks", func(t *testing.T) {
		state := packagetesting.NewParallelTestingState("abc123")
		called := 0
		state.RegisterCallback(func() { called++ })
		state.RunCallbacks()

		if called != 1 {
			t.Fatalf("expected parallel testing callback to run once, got %d", called)
		}
	})
}
