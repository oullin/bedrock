package main

import (
	"net/http"

	"github.com/bedrock/services/inertia-demo/api/auth"
	apierrors "github.com/bedrock/services/inertia-demo/api/errors"
)

func (rt *runtime) registerErrorRoutes(mux *http.ServeMux, authApp auth.App) error {
	return apierrors.RegisterRoutes(rt.routes, mux, apierrors.Container{
		RequireAuth: authApp.RequireAuth,
		Render:      rt.renderPage,
	})
}
