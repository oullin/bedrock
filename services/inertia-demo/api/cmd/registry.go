package main

import (
	"errors"
	"net/http"

	"github.com/bedrock/packages/inertia"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/routegen"
	"github.com/bedrock/services/inertia-demo/api/auth"
	"github.com/bedrock/services/inertia-demo/api/crm"
	demoerrors "github.com/bedrock/services/inertia-demo/api/errors"
	"github.com/bedrock/services/inertia-demo/api/features"
)

func initRoutes() *routegen.Registry {
	routes := routegen.New()

	routes.Add("login", "GET", "/login")

	routes.Add("logout", "POST", "/logout")

	crm.DefineRoutes(routes)

	features.DefineRoutes(routes)

	demoerrors.DefineRoutes(routes)

	return routes
}

func (rt *runtime) withDemoProps(authApp auth.App, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := authApp.CurrentUser(r)

		sidebarOpen := true

		if cookie, err := r.Cookie("sidebar_open"); err == nil {
			sidebarOpen = cookie.Value != "false"
		}

		ctx := r.Context()
		ctx = inertia.SetProps(ctx, protocol.Props{
			"sidebarOpen": sidebarOpen,
			"app": map[string]any{
				"name":        "Inertia.js Kitchen Sink",
				"productLine": "Go Demo Port",
				"environment": "Demo",
			},
			"auth": map[string]any{
				"user": authApp.PublicUser(user),
			},
			"workspace": map[string]any{
				"name": "Inertia Go",
				"plan": "Porting",
			},
			"routes": rt.routes.ManifestProps(),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rt *runtime) renderPage(w http.ResponseWriter, r *http.Request, component string, pageProps protocol.Props) {
	ctx := r.Context()

	if err := rt.inertia.Render(w, r.WithContext(ctx), component, pageProps); err != nil {
		switch {
		case errors.Is(err, protocol.ErrNotFound):
			http.Error(w, "page not found", http.StatusNotFound)
		default:
			http.Error(w, "demo internal error", http.StatusInternalServerError)
		}
	}
}
