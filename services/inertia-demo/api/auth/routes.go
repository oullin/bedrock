package auth

import "net/http"

// RegisterRoutes mounts the auth HTTP routes onto the provided mux.
func (a App) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/login", a.GuestOnly(http.HandlerFunc(a.loginHandler)))
	mux.Handle("/logout", a.RequireAuth(http.HandlerFunc(a.logoutHandler)))
}
