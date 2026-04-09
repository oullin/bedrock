package routes

import (
	"net/http"

	"github.com/bedrock/packages/anvil/routing"
)

// RegisterWeb installs the demo HTTP routes.
func RegisterWeb(router *routing.Router) {
	router.Get("/", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "Bedrock")
	})

	// Auth page routes (Inertia renders)
	router.Get("/login", inertiaPage("Auth/Login"))
	router.Get("/register", inertiaPage("Auth/Register"))
	router.Get("/forgot-password", inertiaPage("Auth/ForgotPassword"))
	router.Get("/reset-password/{token}", inertiaPage("Auth/ResetPassword"))
	router.Get("/verify-email", inertiaPage("Auth/VerifyEmail"))
	router.Get("/two-factor-challenge", inertiaPage("Auth/TwoFactorChallenge"))

	// Authenticated page routes
	router.Get("/dashboard", inertiaPage("Dashboard"))
	router.Get("/user/profile", inertiaPage("Profile/Show"))
	router.Get("/user/api-tokens", inertiaPage("API/Index"))
	router.Get("/teams/create", inertiaPage("Teams/Create"))
	router.Get("/teams/{team}", inertiaPage("Teams/Show"))
}

// inertiaPage returns a handler that renders an Inertia page component.
// In a full setup this would use the inertia-go adapter; here it returns
// a JSON stub matching the Inertia protocol for the testing app.
func inertiaPage(component string) routing.HandlerFunc {
	return func(ctx *routing.Context) error {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.WriteHeader(http.StatusOK)
		_, _ = ctx.Writer.Write([]byte(`{"component":"` + component + `","props":{},"url":"` + ctx.Request.URL.Path + `","version":"1"}`))

		return nil
	}
}
