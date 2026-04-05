package foundation

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	configpkg "github.com/gollin/packages/framework/config"
	"github.com/gollin/packages/framework/console"
	"github.com/gollin/packages/framework/foundation/configuration"
	"github.com/gollin/packages/framework/routing"
	"github.com/gollin/packages/framework/view"
)

// RoutingConfig declares the route and command registrars for an application.
type RoutingConfig struct {
	Web      func(*routing.Router)
	Commands func(*console.Kernel)
	Health   string
}

// Builder configures an application instance.
type Builder struct {
	basePath            string
	routing             RoutingConfig
	configureMiddleware func(*configuration.Middleware)
	configureExceptions func(*configuration.Exceptions)
}

// Application is a bootstrapped demo application.
type Application struct {
	basePath    string
	config      *configpkg.Repository
	router      *routing.Router
	console     *console.Kernel
	renderer    *view.Renderer
	httpHandler http.Handler
}

// Configure begins application construction.
func Configure(basePath string) *Builder {
	return &Builder{basePath: filepath.Clean(basePath)}
}

// WithRouting installs route and command registrars.
func (b *Builder) WithRouting(routingConfig RoutingConfig) *Builder {
	b.routing = routingConfig
	return b
}

// WithMiddleware configures HTTP middleware.
func (b *Builder) WithMiddleware(configure func(*configuration.Middleware)) *Builder {
	b.configureMiddleware = configure
	return b
}

// WithExceptions configures panic handling.
func (b *Builder) WithExceptions(configure func(*configuration.Exceptions)) *Builder {
	b.configureExceptions = configure
	return b
}

// Create builds the application.
func (b *Builder) Create() (*Application, error) {
	env := loadEnvironment(b.basePath)

	repo, err := configpkg.NewBuilder(filepath.Join(b.basePath, "config")).
		WithEnv(env).
		Build(context.Background())
	if err != nil {
		return nil, err
	}

	renderer := view.NewRenderer(filepath.Join(b.basePath, "resources", "views"))
	router := routing.New(renderer)
	kernel := console.New()

	if b.routing.Health != "" {
		router.Get(b.routing.Health, func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "ok")
		})
	}

	if b.routing.Web != nil {
		b.routing.Web(router)
	}

	if b.routing.Commands != nil {
		b.routing.Commands(kernel)
	}

	middleware := configuration.NewMiddleware()
	if b.configureMiddleware != nil {
		b.configureMiddleware(middleware)
	}

	exceptions := configuration.NewExceptions()
	if b.configureExceptions != nil {
		b.configureExceptions(exceptions)
	}

	publicBuild := filepath.Join(b.basePath, "public", "build")
	assets := http.StripPrefix("/build/", http.FileServer(http.Dir(publicBuild)))

	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.HasPrefix(request.URL.Path, "/build/") {
			assets.ServeHTTP(writer, request)
			return
		}

		if maintenance, ok := readMaintenancePage(b.basePath); ok {
			writer.WriteHeader(http.StatusServiceUnavailable)
			_, _ = writer.Write([]byte(maintenance))
			return
		}

		router.ServeHTTP(writer, request)
	})

	return &Application{
		basePath:    b.basePath,
		config:      repo,
		router:      router,
		console:     kernel,
		renderer:    renderer,
		httpHandler: exceptions.Wrap(middleware.Wrap(handler)),
	}, nil
}

// BasePath returns the application's base path.
func (a *Application) BasePath() string {
	return a.basePath
}

// Config returns the loaded configuration repository.
func (a *Application) Config() *configpkg.Repository {
	return a.config
}

// Console returns the application console kernel.
func (a *Application) Console() *console.Kernel {
	return a.console
}

// Renderer returns the application view renderer.
func (a *Application) Renderer() *view.Renderer {
	return a.renderer
}

// HandleCommand executes a console command and returns its exit code.
func (a *Application) HandleCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	return a.console.Handle(args, stdout, stderr)
}

// HandleRequest serves an HTTP request.
func (a *Application) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	a.httpHandler.ServeHTTP(writer, request)
}

// ServeHTTP implements http.Handler.
func (a *Application) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.HandleRequest(writer, request)
}

func loadEnvironment(basePath string) map[string]string {
	env := map[string]string{}

	for _, candidate := range []string{
		filepath.Join(basePath, ".env.example"),
		filepath.Join(basePath, ".env"),
	} {
		file, err := os.Open(candidate)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			key, value, ok := strings.Cut(line, "=")
			if ok {
				env[strings.TrimSpace(key)] = strings.TrimSpace(value)
			}
		}

		_ = file.Close()
	}

	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			env[key] = value
		}
	}

	return env
}

func readMaintenancePage(basePath string) (string, bool) {
	for _, candidate := range []string{
		filepath.Join(basePath, "storage", "framework", "maintenance.html"),
		filepath.Join(basePath, "storage", "framework", "maintenance.php"),
	} {
		payload, err := os.ReadFile(candidate)
		if err == nil {
			return string(payload), true
		}
	}

	return "", false
}
