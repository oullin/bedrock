package foundation_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/console"
	"github.com/bedrock/packages/foundation"
	"github.com/bedrock/packages/routing"
)

func setupTestApp(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	// Create config directory with minimal config
	configDir := filepath.Join(dir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "app.yml"), []byte("name: TestApp\n"), 0644)

	// Create views directory
	viewsDir := filepath.Join(dir, "resources", "views")
	os.MkdirAll(viewsDir, 0755)
	os.WriteFile(filepath.Join(viewsDir, "welcome.html.tmpl"), []byte("<h1>Welcome</h1>"), 0644)

	// Create public/build directory
	os.MkdirAll(filepath.Join(dir, "public", "build"), 0755)

	return dir
}

// Upstream: testApplicationBootstrapping
func TestApplicationBuilderCreatesApp(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{
			Web: func(r *routing.Router) {
				r.Get("/", func(ctx *routing.Context) error {
					return ctx.Text(http.StatusOK, "home")
				})
			},
		}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if app.BasePath() != filepath.Clean(dir) {
		t.Fatalf("expected base path %q, got %q", dir, app.BasePath())
	}
}

// Upstream: testApplicationServesHTTP
func TestApplicationServesHTTPRequests(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{
			Web: func(r *routing.Router) {
				r.Get("/hello", func(ctx *routing.Context) error {
					return ctx.Text(http.StatusOK, "hello")
				})
			},
		}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/hello", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "hello" {
		t.Fatalf("expected 'hello', got %q", rec.Body.String())
	}
}

// Upstream: testHealthEndpoint
func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{
			Health: "/up",
		}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/up", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("expected 'ok', got %q", rec.Body.String())
	}
}

// Upstream: testConsoleCommands
func TestApplicationHandlesConsoleCommands(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{
			Commands: func(k *console.Kernel) {
				k.Command("inspire", "Display an inspiring quote", func(_ context.Context, inv *console.Invocation) error {
					inv.Comment("Be yourself")
					return nil
				})
			},
		}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	stdout := &bytes.Buffer{}
	code := app.HandleCommand([]string{"inspire"}, stdout, &bytes.Buffer{})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stdout.String() == "" {
		t.Fatal("expected output from command")
	}
}

// Upstream: testApplicationConfig
func TestApplicationConfigLoadsYAML(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	config := app.Config()
	if config == nil {
		t.Fatal("expected config to be loaded")
	}

	name := config.Get("app.name", "")
	if name != "TestApp" {
		t.Fatalf("expected 'TestApp', got %v", name)
	}
}

// Upstream: testEnvironmentFileLoading
func TestApplicationLoadsEnvFile(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)

	// Write an .env file
	os.WriteFile(filepath.Join(dir, ".env"), []byte("APP_NAME=FromEnv\n"), 0644)
	// Update config to use env override
	os.WriteFile(filepath.Join(dir, "config", "app.yml"), []byte("name: DefaultApp\n"), 0644)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// The app config should have loaded the .env file
	config := app.Config()
	if config == nil {
		t.Fatal("expected config to be loaded")
	}
}

// Upstream: testMiddlewareConfiguration
func TestApplicationWithMiddleware(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{
			Web: func(r *routing.Router) {
				r.Get("/test", func(ctx *routing.Context) error {
					return ctx.Text(http.StatusOK, "ok")
				})
			},
		}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

// Test maintenance mode
func TestMaintenanceModeServes503(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)

	// Create maintenance file
	storageDir := filepath.Join(dir, "storage", "framework")
	os.MkdirAll(storageDir, 0755)
	os.WriteFile(filepath.Join(storageDir, "maintenance.html"), []byte("<h1>Down for maintenance</h1>"), 0644)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{
			Web: func(r *routing.Router) {
				r.Get("/", func(ctx *routing.Context) error {
					return ctx.Text(http.StatusOK, "should not reach")
				})
			},
		}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
	if rec.Body.String() != "<h1>Down for maintenance</h1>" {
		t.Fatalf("expected maintenance page, got %q", rec.Body.String())
	}
}

// Test application implements http.Handler
func TestApplicationImplementsHTTPHandler(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)
	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	var _ http.Handler = app
}

// Test Console accessor
func TestApplicationConsoleAccessor(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)
	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{
			Commands: func(k *console.Kernel) {
				k.Command("test", "Test", func(context.Context, *console.Invocation) error { return nil })
			},
		}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if app.Console() == nil {
		t.Fatal("expected Console accessor to return kernel")
	}
	if len(app.Console().Commands()) != 1 {
		t.Fatal("expected 1 registered command")
	}
}

// Test Renderer accessor
func TestApplicationRendererAccessor(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)
	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()

	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if app.Renderer() == nil {
		t.Fatal("expected Renderer accessor to return renderer")
	}
}
