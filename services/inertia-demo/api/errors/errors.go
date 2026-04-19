package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/routegen"
)

// Container contains the host application integrations required by the errors package.
type Container struct {
	RequireAuth func(http.Handler) http.Handler
	Render      func(http.ResponseWriter, *http.Request, string, protocol.Props)
}

// Validate checks that all required dependencies are set.
func (c Container) Validate() error {
	var errs []error

	if c.RequireAuth == nil {
		errs = append(errs, stderrors.New("errors: RequireAuth must not be nil"))
	}

	if c.Render == nil {
		errs = append(errs, stderrors.New("errors: Render must not be nil"))
	}

	return stderrors.Join(errs...)
}

// DefineRoutes registers error showcase route metadata (name, method, pattern)
// on the given registry without mounting handlers.
func DefineRoutes(routes *routegen.Registry) {
	routes.Group("features.errors", "/features/errors", func(g *routegen.Group) {
		g.Add("http-error", "GET", "/http-error")
		g.Add("network-errors", "GET", "/network-errors")
	})
}

// RegisterRoutes mounts the error showcase HTTP routes onto the provided mux.
func RegisterRoutes(routes *routegen.Registry, mux *http.ServeMux, container Container) error {
	if err := container.Validate(); err != nil {
		return fmt.Errorf("errors: %w", err)
	}

	auth := func(h http.HandlerFunc) http.Handler {
		return container.RequireAuth(h)
	}

	routes.Handle("features.errors.http-error", auth(httpErrorHandler(container)), mux)

	mux.Handle("/features/errors/http-error/{code}", auth(httpErrorTriggerHandler()))

	routes.Handle("features.errors.network-errors", auth(networkErrorsHandler(container)), mux)

	return nil
}

func httpErrorHandler(container Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		container.Render(w, r, "Features/Errors/HttpError", protocol.Props{})
	}
}

func httpErrorTriggerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		switch code {
		case "403":
			http.Error(w, "Forbidden", http.StatusForbidden)
		case "404":
			http.Error(w, "Not Found", http.StatusNotFound)
		case "500":
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		case "unhandled":
			http.Error(w, "I'm a teapot", http.StatusTeapot)
		default:
			http.NotFound(w, r)
		}
	}
}

func networkErrorsHandler(container Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		container.Render(w, r, "Features/Errors/NetworkErrors", protocol.Props{})
	}
}
