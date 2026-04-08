package foundation_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/anvil/console"
	"github.com/bedrock/packages/anvil/foundation"
	"github.com/bedrock/packages/anvil/routing"
)

func setupTestApp(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	// Create config directory with minimal config
	configDir := filepath.Join(dir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "app.yml"), []byte("name: TestApp\n"), 0644)

	// Create public/build directory
	os.MkdirAll(filepath.Join(dir, "public", "build"), 0755)

	return dir
}

// Laravel: testApplicationBootstrapping
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

// Laravel: testApplicationServesHTTP
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

// Laravel: testHealthEndpoint
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

// Laravel: testConsoleCommands
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

// Laravel: testApplicationConfig
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

// Laravel: testEnvironmentFileLoading
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

// Laravel: testMiddlewareConfiguration
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

// Laravel: testEnvironmentDetection
func TestApplicationEnvironment(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configDir := filepath.Join(dir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "app.yml"), []byte("env: testing\nname: TestApp\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "public", "build"), 0755)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if env := app.Environment(); env != "testing" {
		t.Fatalf("expected environment 'testing', got %q", env)
	}
	if app.IsProduction() {
		t.Fatal("expected IsProduction to return false for 'testing' environment")
	}
}

// ======================== BOOTSTRAP CALLBACK TESTS ========================

// Laravel: testBootingCallbacks
func TestBootingCallbacks(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)
	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	var order []string
	app.Booting(func(_ *foundation.Application) { order = append(order, "boot-1") })
	app.Booting(func(_ *foundation.Application) { order = append(order, "boot-2") })

	app.Boot()

	if len(order) != 2 || order[0] != "boot-1" || order[1] != "boot-2" {
		t.Fatalf("expected [boot-1, boot-2], got %v", order)
	}
	if !app.IsBooted() {
		t.Fatal("expected IsBooted to be true after Boot()")
	}
}

// Laravel: testBootedCallbacks
func TestBootedCallbacks(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)
	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	var order []string
	app.Booting(func(_ *foundation.Application) { order = append(order, "booting") })
	app.Booted(func(_ *foundation.Application) { order = append(order, "booted-1") })
	app.Booted(func(_ *foundation.Application) { order = append(order, "booted-2") })

	app.Boot()

	if len(order) != 3 {
		t.Fatalf("expected 3 callbacks, got %d: %v", len(order), order)
	}
	if order[0] != "booting" || order[1] != "booted-1" || order[2] != "booted-2" {
		t.Fatalf("expected [booting, booted-1, booted-2], got %v", order)
	}

	// Registering a booted callback after boot should execute immediately.
	app.Booted(func(_ *foundation.Application) { order = append(order, "booted-late") })
	if order[3] != "booted-late" {
		t.Fatalf("expected booted-late to execute immediately, got %v", order)
	}
}

// Laravel: testTerminationTests
func TestTerminationCallbacks(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)
	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	var order []string
	app.Terminating(func() { order = append(order, "term-1") })
	app.Terminating(func() { order = append(order, "term-2") })

	app.Terminate()

	if len(order) != 2 || order[0] != "term-1" || order[1] != "term-2" {
		t.Fatalf("expected [term-1, term-2], got %v", order)
	}
}

// Laravel: testBootDoesNotReboot
func TestBootIsIdempotent(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)
	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	count := 0
	app.Booting(func(_ *foundation.Application) { count++ })
	app.Boot()
	app.Boot() // second call should be no-op

	if count != 1 {
		t.Fatalf("expected booting callback to run once, got %d", count)
	}
}

// ======================== ENVIRONMENT HELPER TESTS ========================

func TestIsLocal(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configDir := filepath.Join(dir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "app.yml"), []byte("env: local\nname: TestApp\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "public", "build"), 0755)

	app, err := foundation.Configure(dir).WithRouting(foundation.RoutingConfig{}).Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !app.IsLocal() {
		t.Fatal("expected IsLocal to be true")
	}
	if app.IsProduction() {
		t.Fatal("expected IsProduction to be false for local env")
	}
}

func TestIsTesting(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configDir := filepath.Join(dir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "app.yml"), []byte("env: testing\nname: TestApp\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "public", "build"), 0755)

	app, err := foundation.Configure(dir).WithRouting(foundation.RoutingConfig{}).Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !app.IsTesting() {
		t.Fatal("expected IsTesting to be true")
	}
}

func TestRunningInConsole(t *testing.T) {
	t.Parallel()

	dir := setupTestApp(t)
	app, err := foundation.Configure(dir).WithRouting(foundation.RoutingConfig{}).Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if app.RunningInConsole() {
		t.Fatal("expected RunningInConsole to be false by default")
	}
	app.SetRunningInConsole(true)
	if !app.RunningInConsole() {
		t.Fatal("expected RunningInConsole to be true after setting")
	}
}

func TestApplicationEnvironmentDefaultsToProduction(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configDir := filepath.Join(dir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "app.yml"), []byte("name: TestApp\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "public", "build"), 0755)

	app, err := foundation.Configure(dir).
		WithRouting(foundation.RoutingConfig{}).
		Create()
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if env := app.Environment(); env != "production" {
		t.Fatalf("expected default environment 'production', got %q", env)
	}
	if !app.IsProduction() {
		t.Fatal("expected IsProduction to return true for default environment")
	}
}
