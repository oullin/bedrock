package bootstrap

import (
	"path/filepath"
	"runtime"

	"github.com/gollin/demo/routes"
	"github.com/gollin/packages/framework/foundation"
	foundationconfiguration "github.com/gollin/packages/framework/foundation/configuration"
)

// New bootstraps the root demo application.
func New() (*foundation.Application, error) {
	_, file, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(file))

	return foundation.Configure(basePath).
		WithRouting(foundation.RoutingConfig{
			Web:      routes.RegisterWeb,
			Commands: routes.RegisterConsole,
			Health:   "/up",
		}).
		WithMiddleware(func(*foundationconfiguration.Middleware) {}).
		WithExceptions(func(*foundationconfiguration.Exceptions) {}).
		Create()
}
