package main

import (
	"net/http"

	"github.com/bedrock/services/demo/inertia/api/auth"
	"github.com/bedrock/services/demo/inertia/api/crm"
)

func (rt *runtime) registerCRMRoutes(mux *http.ServeMux, authApp auth.App) error {
	redirectFn := func(w http.ResponseWriter, r *http.Request, url string) {
		rt.inertia.Redirect(w, r, url)
	}

	return crm.RegisterRoutes(rt.routes, mux, crm.Container{
		DB:          rt.db,
		RequireAuth: authApp.RequireAuth,
		Render:      rt.renderPage,
		Redirect:    redirectFn,

		RouteURL:    rt.routes.URL,
		SetFlash:    rt.flashStore.Set,
		CurrentUser: authApp.CurrentUser,
	})
}
