package crm

import (
	"fmt"
	"net/http"

	"github.com/bedrock/packages/inertia/wayfinder"
)

type app struct {
	container Container
	repo      *databaseRepository
}

// DefineRoutes registers CRM route metadata (name, method, pattern) on the
// given registry without mounting handlers.
func DefineRoutes(routes *wayfinder.Registry) {
	routes.Add("dashboard", "GET", "/dashboard")

	routes.Group("contacts", "/contacts", func(g *wayfinder.Group) {
		g.Add("index", "GET", "")
		g.Add("create", "GET", "/create")
		g.Add("store", "POST", "")
		g.Add("show", "GET", "/{contact}")
		g.Add("edit", "GET", "/{contact}/edit")
		g.Add("update", "POST", "/{contact}")
		g.Add("destroy", "DELETE", "/{contact}")
		g.Add("favorite", "POST", "/{contact}/favorite")

		g.Group("notes", "", func(ng *wayfinder.Group) {
			ng.Add("store", "POST", "/{contact}/notes")
		})
	})

	routes.Group("organizations", "/organizations", func(g *wayfinder.Group) {
		g.Add("index", "GET", "")
		g.Add("show", "GET", "/{organization}")
		g.Add("update", "POST", "/{organization}")
	})
}

// RegisterRoutes mounts the CRM HTTP routes onto the provided mux.
func RegisterRoutes(routes *wayfinder.Registry, mux *http.ServeMux, container Container) error {
	if err := container.Validate(); err != nil {
		return fmt.Errorf("crm: %w", err)
	}

	app, err := newApp(container)

	if err != nil {
		return fmt.Errorf("crm: %w", err)
	}

	auth := func(h http.HandlerFunc) http.Handler {
		return container.RequireAuth(h)
	}

	routes.Handle("dashboard", auth(app.dashboardHandler), mux)
	routes.Handle("contacts.index", auth(app.contactsHandler), mux)
	routes.Handle("contacts.create", auth(app.contactsCreateHandler), mux)

	mux.Handle("/contacts/", auth(app.contactByIDHandler))

	routes.Handle("organizations.index", auth(app.organizationsHandler), mux)

	mux.Handle("/organizations/", auth(app.organizationByIDHandler))

	return nil
}

func newApp(container Container) (app, error) {
	repo, err := newRepository(container.DB)

	if err != nil {
		return app{}, err
	}

	return app{
		container: container,
		repo:      repo,
	}, nil
}
