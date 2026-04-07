package view_test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/anvil/view"
)

func setupTestViews(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "hello.html.tmpl"), []byte("<h1>Hello, {{.Name}}!</h1>"), 0644)
	os.WriteFile(filepath.Join(dir, "static.html.tmpl"), []byte("<p>Static content</p>"), 0644)

	return dir
}

// Upstream: testBasicViewRendering
func TestRenderBasicView(t *testing.T) {
	t.Parallel()

	dir := setupTestViews(t)
	renderer := view.NewRenderer(dir)

	result, err := renderer.Render("hello", map[string]any{"Name": "World"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	expected := "<h1>Hello, World!</h1>"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

// Upstream: testViewWithoutData
func TestRenderStaticView(t *testing.T) {
	t.Parallel()

	dir := setupTestViews(t)
	renderer := view.NewRenderer(dir)

	result, err := renderer.Render("static", nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if result != "<p>Static content</p>" {
		t.Fatalf("expected static content, got %q", result)
	}
}

// Upstream: testViewNotFound
func TestRenderMissingViewReturnsError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	renderer := view.NewRenderer(dir)

	_, err := renderer.Render("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

// Test WriteHTML response
func TestWriteHTMLResponse(t *testing.T) {
	t.Parallel()

	dir := setupTestViews(t)
	renderer := view.NewRenderer(dir)

	rec := httptest.NewRecorder()
	err := renderer.WriteHTML(rec, 200, "hello", map[string]any{"Name": "Test"})
	if err != nil {
		t.Fatalf("WriteHTML: %v", err)
	}

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("expected text/html content-type, got %q", ct)
	}
	if rec.Body.String() != "<h1>Hello, Test!</h1>" {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

// Test WriteHTML with missing view returns error
func TestWriteHTMLMissingViewReturnsError(t *testing.T) {
	t.Parallel()

	renderer := view.NewRenderer(t.TempDir())
	rec := httptest.NewRecorder()
	err := renderer.WriteHTML(rec, 200, "missing", nil)
	if err == nil {
		t.Fatal("expected error for missing template in WriteHTML")
	}
}

// Test asset function map
func TestAssetFunctionInTemplate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "with_asset.html.tmpl"),
		[]byte(`<link href="{{asset "css/app.css"}}">`), 0644)

	renderer := view.NewRenderer(dir)
	result, err := renderer.Render("with_asset", nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	expected := `<link href="/build/css/app.css">`
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

// Test asset function with leading slash
func TestAssetFunctionWithLeadingSlash(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "abs_asset.html.tmpl"),
		[]byte(`<img src="{{asset "/images/logo.png"}}">`), 0644)

	renderer := view.NewRenderer(dir)
	result, err := renderer.Render("abs_asset", nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	// Leading slash should be preserved as-is
	expected := `<img src="/images/logo.png">`
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

// Test renderer with empty directory
func TestRendererEmptyDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	renderer := view.NewRenderer(dir)

	_, err := renderer.Render("anything", nil)
	if err == nil {
		t.Fatal("expected error for view in empty directory")
	}
}

// Upstream: testViewShareGlobalData
func TestRendererShareGlobalData(t *testing.T) {
	t.Parallel()

	renderer := view.NewRenderer(t.TempDir())

	renderer.Share("appName", "Bedrock")
	renderer.Share("version", "1.0")

	shared := renderer.Shared()
	if shared["appName"] != "Bedrock" {
		t.Fatalf("expected shared 'appName' = 'Bedrock', got %v", shared["appName"])
	}
	if shared["version"] != "1.0" {
		t.Fatalf("expected shared 'version' = '1.0', got %v", shared["version"])
	}

	// Overwrite existing shared value.
	renderer.Share("appName", "Bedrock v2")
	shared = renderer.Shared()
	if shared["appName"] != "Bedrock v2" {
		t.Fatalf("expected overwritten 'appName' = 'Bedrock v2', got %v", shared["appName"])
	}
}

func TestRendererSharedReturnsEmptyMapWhenNothingShared(t *testing.T) {
	t.Parallel()

	renderer := view.NewRenderer(t.TempDir())
	shared := renderer.Shared()

	if shared == nil {
		t.Fatal("expected non-nil map")
	}
	if len(shared) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(shared))
	}
}
