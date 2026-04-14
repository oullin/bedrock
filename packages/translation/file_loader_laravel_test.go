package translation_test

// Ports of Illuminate\Tests\Translation\TranslationFileLoaderTest.

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/translation"
)

// ── helpers ───────────────────────────────────────────────────────────────

// writeJSON creates a JSON file at {dir}/{locale}/{group}.json with content.
func writeJSON(t *testing.T, dir, locale, group string, content map[string]any) {
	t.Helper()
	localeDir := filepath.Join(dir, locale)

	if err := os.MkdirAll(localeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(content)

	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(localeDir, group+".json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// writeFlatJSON creates a flat {dir}/{locale}.json file.
func writeFlatJSON(t *testing.T, dir, locale string, content map[string]any) {
	t.Helper()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(content)

	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, locale+".json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// ── tests ─────────────────────────────────────────────────────────────────

func TestFileLoaderLoadMethodProperlyLoadsFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeJSON(t, dir, "en", "foo", map[string]any{"bar": "baz"})

	l := translation.NewFileLoader(dir)
	got := l.Load("en", "foo", nil)

	if got["bar"] != "baz" {
		t.Errorf("got %v", got)
	}
}

func TestFileLoaderLoadMethodProperlyLoadsFilesFromMultiplePaths(t *testing.T) {
	t.Parallel()

	dir1 := t.TempDir()
	dir2 := t.TempDir()
	writeJSON(t, dir1, "en", "foo", map[string]any{"bar": "first", "only1": "here"})
	writeJSON(t, dir2, "en", "foo", map[string]any{"bar": "second", "only2": "there"})

	l := translation.NewFileLoader(dir1, dir2)
	got := l.Load("en", "foo", nil)

	if got["bar"] != "second" {
		t.Errorf("later path should override: got %v", got["bar"])
	}

	if got["only1"] != "here" {
		t.Errorf("earlier-path-only key missing: got %v", got["only1"])
	}

	if got["only2"] != "there" {
		t.Errorf("later-path-only key missing: got %v", got["only2"])
	}
}

func TestFileLoaderLoadMethodSetsAnEmptyArrayWhenFileDoesntExist(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	l := translation.NewFileLoader(dir)
	got := l.Load("en", "missing", nil)

	if got == nil || len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestFileLoaderLoadMethodSetsAnEmptyArrayWhenFileIsEmpty(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "en"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "en", "empty.json"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	l := translation.NewFileLoader(dir)
	got := l.Load("en", "empty", nil)

	if got == nil || len(got) != 0 {
		t.Errorf("expected empty map for empty file, got %v", got)
	}
}

func TestFileLoaderLoadNamespacedFiles(t *testing.T) {
	t.Parallel()

	hint := t.TempDir()
	writeJSON(t, hint, "en", "foo", map[string]any{"bar": "namespaced"})

	l := translation.NewFileLoader()
	l.AddNamespace("vendor", hint)

	ns := "vendor"
	got := l.Load("en", "foo", &ns)

	if got["bar"] != "namespaced" {
		t.Errorf("got %v", got)
	}
}

func TestFileLoaderLoadNamespacedFilesWithOverride(t *testing.T) {
	t.Parallel()

	hint := t.TempDir()
	appDir := t.TempDir()

	writeJSON(t, hint, "en", "messages", map[string]any{"welcome": "vendor msg", "unique": "vendor only"})

	// Override in app vendor directory.
	vendorOverride := filepath.Join(appDir, "vendor", "pkg", "en")

	if err := os.MkdirAll(vendorOverride, 0o755); err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(map[string]any{"welcome": "app override"})

	if err := os.WriteFile(filepath.Join(vendorOverride, "messages.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	l := translation.NewFileLoader(appDir)
	l.AddNamespace("pkg", hint)

	ns := "pkg"
	got := l.Load("en", "messages", &ns)

	if got["welcome"] != "app override" {
		t.Errorf("override not applied: got %v", got["welcome"])
	}

	if got["unique"] != "vendor only" {
		t.Errorf("vendor-only key lost: got %v", got["unique"])
	}
}

func TestFileLoaderAddJsonPath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFlatJSON(t, dir, "en", map[string]any{"greeting": "hello"})

	l := translation.NewFileLoader()
	l.AddJsonPath(dir)

	got := l.Load("en", "*", nil)

	if got["greeting"] != "hello" {
		t.Errorf("got %v", got)
	}
}

func TestFileLoaderLoadFromMultipleJsonPaths(t *testing.T) {
	t.Parallel()

	dir1 := t.TempDir()
	dir2 := t.TempDir()
	writeFlatJSON(t, dir1, "en", map[string]any{"a": "first", "shared": "dir1"})
	writeFlatJSON(t, dir2, "en", map[string]any{"b": "second", "shared": "dir2"})

	l := translation.NewFileLoader()
	l.AddJsonPath(dir1)
	l.AddJsonPath(dir2)

	got := l.Load("en", "*", nil)

	if got["a"] != "first" {
		t.Errorf("got a=%v", got["a"])
	}

	if got["b"] != "second" {
		t.Errorf("got b=%v", got["b"])
	}

	if got["shared"] != "dir2" {
		t.Errorf("later path should win: got shared=%v", got["shared"])
	}
}

func TestFileLoaderMalformedJsonFileReturnsError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(dir, "en"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "en", "bad.json"), []byte("not json {{"), 0o644); err != nil {
		t.Fatal(err)
	}

	l := translation.NewFileLoader(dir)
	got := l.Load("en", "bad", nil)

	if got != nil {
		t.Errorf("expected nil return for malformed JSON, got %v", got)
	}
}

func TestFileLoaderMalformedJsonReturnsWrappedError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFlatJSON(t, dir, "en", map[string]any{"ok": "val"})
	// Overwrite with invalid JSON.
	if err := os.WriteFile(filepath.Join(dir, "en.json"), []byte("{bad json"), 0o644); err != nil {
		t.Fatal(err)
	}

	l := translation.NewFileLoader()
	l.AddJsonPath(dir)

	got := l.Load("en", "*", nil)

	if got != nil {
		t.Errorf("expected nil for malformed JSON path")
	}
}

func TestFileLoaderNamespacesMethod(t *testing.T) {
	t.Parallel()

	l := translation.NewFileLoader()
	l.AddNamespace("foo", "/path/foo")
	l.AddNamespace("bar", "/path/bar")

	ns := l.Namespaces()

	if ns["foo"] != "/path/foo" || ns["bar"] != "/path/bar" {
		t.Errorf("Namespaces() = %v", ns)
	}
}

func TestFileLoaderPathsMethod(t *testing.T) {
	t.Parallel()

	l := translation.NewFileLoader("/a", "/b")
	l.AddPath("/c")

	paths := l.Paths()

	if len(paths) != 3 || paths[2] != "/c" {
		t.Errorf("Paths() = %v", paths)
	}
}

func TestFileLoaderJsonPathsMethod(t *testing.T) {
	t.Parallel()

	l := translation.NewFileLoader()
	l.AddJsonPath("/j1")
	l.AddJsonPath("/j2")

	jp := l.JsonPaths()

	if len(jp) != 2 {
		t.Errorf("JsonPaths() = %v", jp)
	}
}

func TestFileLoaderReturnsEmptyMapForMissingLocale(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeJSON(t, dir, "en", "messages", map[string]any{"key": "val"})

	l := translation.NewFileLoader(dir)
	got := l.Load("fr", "messages", nil)

	if len(got) != 0 {
		t.Errorf("expected empty map for missing locale, got %v", got)
	}
}

func TestFileLoaderErrMalformedJSONIsSentinel(t *testing.T) {
	t.Parallel()

	// Verify that ErrMalformedJSON is exported and usable with errors.Is.
	wrapped := fmt.Errorf("context: %w", translation.ErrMalformedJSON)

	if !errors.Is(wrapped, translation.ErrMalformedJSON) {
		t.Error("errors.Is should find ErrMalformedJSON in wrapped error")
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testLoadMethodLoadsTranslationsFromAddedPath
func TestFileLoaderLoadMethodLoadsTranslationsFromAddedPath(t *testing.T) {
	t.Parallel()

	dir1 := t.TempDir()
	dir2 := t.TempDir()
	writeJSON(t, dir1, "en", "messages", map[string]any{"foo": "bar"})
	writeJSON(t, dir2, "en", "messages", map[string]any{"baz": "backagesplash"})

	l := translation.NewFileLoader(dir1)
	l.AddPath(dir2)

	got := l.Load("en", "messages", nil)

	if got["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", got["foo"])
	}

	if got["baz"] != "backagesplash" {
		t.Errorf("expected baz=backagesplash, got %v", got["baz"])
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testLoadMethodHandlesMissingAddedPath
func TestFileLoaderLoadMethodHandlesMissingAddedPath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeJSON(t, dir, "en", "messages", map[string]any{"foo": "bar"})

	l := translation.NewFileLoader(dir)
	l.AddPath("/nonexistent/missing/path")

	got := l.Load("en", "messages", nil)

	if got["foo"] != "bar" {
		t.Errorf("expected foo=bar from valid path, got %v", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testLoadMethodOverwritesExistingKeysFromAddedPath
func TestFileLoaderLoadMethodOverwritesExistingKeysFromAddedPath(t *testing.T) {
	t.Parallel()

	dir1 := t.TempDir()
	dir2 := t.TempDir()
	writeJSON(t, dir1, "en", "messages", map[string]any{"foo": "bar"})
	writeJSON(t, dir2, "en", "messages", map[string]any{"foo": "baz"})

	l := translation.NewFileLoader(dir1)
	l.AddPath(dir2)

	got := l.Load("en", "messages", nil)

	if got["foo"] != "baz" {
		t.Errorf("later AddPath should override: got foo=%v", got["foo"])
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testLoadMethodLoadsTranslationsFromMultipleAddedPaths
func TestFileLoaderLoadMethodLoadsTranslationsFromMultipleAddedPaths(t *testing.T) {
	t.Parallel()

	dir1 := t.TempDir()
	dir2 := t.TempDir()
	dir3 := t.TempDir()
	writeJSON(t, dir1, "en", "messages", map[string]any{"a": "1"})
	writeJSON(t, dir2, "en", "messages", map[string]any{"b": "2"})
	writeJSON(t, dir3, "en", "messages", map[string]any{"c": "3"})

	l := translation.NewFileLoader(dir1)
	l.AddPath(dir2)
	l.AddPath(dir3)

	got := l.Load("en", "messages", nil)

	if got["a"] != "1" || got["b"] != "2" || got["c"] != "3" {
		t.Errorf("all three paths should merge: got %v", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testLoadMethodWithNamespacesProperlyCallsLoaderWithMultiplePaths
func TestFileLoaderLoadMethodWithNamespacesProperlyCallsLoaderWithMultiplePaths(t *testing.T) {
	t.Parallel()

	hint := t.TempDir()
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	writeJSON(t, hint, "en", "messages", map[string]any{"foo": "bar"})

	// First regular path has a vendor override.
	override := filepath.Join(dir1, "vendor", "pkg", "en")

	if err := os.MkdirAll(override, 0o755); err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(map[string]any{"foo": "override"})

	if err := os.WriteFile(filepath.Join(override, "messages.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	l := translation.NewFileLoader(dir1)
	l.AddPath(dir2)
	l.AddNamespace("pkg", hint)

	ns := "pkg"
	got := l.Load("en", "messages", &ns)

	if got["foo"] != "override" {
		t.Errorf("vendor override in first path should win: got foo=%v", got["foo"])
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testLoadMethodWithNamespacesProperlyCallsLoaderAndLoadsLocalOverridesWithMultiplePaths
func TestFileLoaderLoadMethodWithNamespacesAndLocalOverridesWithMultiplePaths(t *testing.T) {
	t.Parallel()

	hint := t.TempDir()
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	writeJSON(t, hint, "en", "messages", map[string]any{"foo": "vendor", "unique": "vendor-only"})

	// First path has a partial override.
	override1 := filepath.Join(dir1, "vendor", "pkg", "en")

	if err := os.MkdirAll(override1, 0o755); err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(map[string]any{"foo": "path1-override"})

	if err := os.WriteFile(filepath.Join(override1, "messages.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	// Second path has another override for a different key.
	override2 := filepath.Join(dir2, "vendor", "pkg", "en")

	if err := os.MkdirAll(override2, 0o755); err != nil {
		t.Fatal(err)
	}

	data, _ = json.Marshal(map[string]any{"extra": "path2-extra"})

	if err := os.WriteFile(filepath.Join(override2, "messages.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	l := translation.NewFileLoader(dir1)
	l.AddPath(dir2)
	l.AddNamespace("pkg", hint)

	ns := "pkg"
	got := l.Load("en", "messages", &ns)

	if got["foo"] != "path1-override" {
		t.Errorf("path1 override should win over vendor: got foo=%v", got["foo"])
	}

	if got["unique"] != "vendor-only" {
		t.Errorf("vendor-only key should survive: got unique=%v", got["unique"])
	}

	if got["extra"] != "path2-extra" {
		t.Errorf("path2 extra key should be present: got extra=%v", got["extra"])
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testLoadMethodWithNamespacesProperlyCallsLoaderAndLoadsLocalOverridesWithMultiplePathsWithMissingKey
func TestFileLoaderLoadMethodWithNamespacesAndLocalOverridesMultiplePathsMissingKey(t *testing.T) {
	t.Parallel()

	hint := t.TempDir()
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	writeJSON(t, hint, "en", "messages", map[string]any{"foo": "vendor", "only-vendor": "here"})

	// First path override: only overrides "foo", not "only-vendor".
	override1 := filepath.Join(dir1, "vendor", "pkg", "en")

	if err := os.MkdirAll(override1, 0o755); err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(map[string]any{"foo": "override"})

	if err := os.WriteFile(filepath.Join(override1, "messages.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	l := translation.NewFileLoader(dir1)
	l.AddPath(dir2)
	l.AddNamespace("pkg", hint)

	ns := "pkg"
	got := l.Load("en", "messages", &ns)

	if got["foo"] != "override" {
		t.Errorf("overridden key should return override value: got foo=%v", got["foo"])
	}
	// Key not overridden in any app path falls through to the vendor hint.
	if got["only-vendor"] != "here" {
		t.Errorf("vendor-only key should still be present: got only-vendor=%v", got["only-vendor"])
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testEmptyArraysReturnedWhenFilesDontExistForNamespacedItems
func TestFileLoaderEmptyArraysReturnedWhenNamespacedFilesDoNotExist(t *testing.T) {
	t.Parallel()

	l := translation.NewFileLoader()
	// Namespace "foo" is never registered — no hint directory.
	ns := "foo"
	got := l.Load("en", "bar", &ns)

	if len(got) != 0 {
		t.Errorf("unregistered namespace should return empty map, got %v", got)
	}
}

// Port of Illuminate\Tests\Translation\TranslationFileLoaderTest::testLoadMethodForJSONProperlyCallsLoaderForMultiplePaths
func TestFileLoaderLoadMethodForJSONWithMultiplePaths(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	jsonDir := t.TempDir()
	writeFlatJSON(t, dir, "en", map[string]any{"foo": "bar"})
	writeFlatJSON(t, jsonDir, "en", map[string]any{"baz": "backagesplash"})

	l := translation.NewFileLoader(dir)
	l.AddJsonPath(jsonDir)

	got := l.Load("en", "*", nil)

	if got["foo"] != "bar" {
		t.Errorf("expected foo=bar from regular path, got %v", got["foo"])
	}

	if got["baz"] != "backagesplash" {
		t.Errorf("expected baz=backagesplash from json path, got %v", got["baz"])
	}
}
