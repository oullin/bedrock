package authkit

import "net/http"

// Router is a minimal interface for registering routes.
type Router interface {
	Post(pattern string, handler http.Handler)
	Put(pattern string, handler http.Handler)
	Get(pattern string, handler http.Handler)
	Delete(pattern string, handler http.Handler)
}

// StdMuxRouter adapts net/http.ServeMux to the Router interface.
type StdMuxRouter struct {
	Mux *http.ServeMux
}

// Post registers a POST route.

// Put registers a PUT route.

// Get registers a GET route.

// Delete registers a DELETE route.

// RouteConfig holds optional dependencies for route registration.
// Only the dependencies relevant to enabled features need to be set.
type RouteConfig struct {
	Tokens       TokenRepository
	Sessions     SessionRepository
	Photos       UpdatesProfilePhotos
	DeletePhotos DeletesProfilePhotos
	DeleteUser   DeletesUsers
}

func (r *StdMuxRouter) Post(pattern string, handler http.Handler) {
	r.Mux.Handle("POST "+pattern, handler)
}

func (r *StdMuxRouter) Put(pattern string, handler http.Handler) {
	r.Mux.Handle("PUT "+pattern, handler)
}

func (r *StdMuxRouter) Get(pattern string, handler http.Handler) {
	r.Mux.Handle("GET "+pattern, handler)
}

func (r *StdMuxRouter) Delete(pattern string, handler http.Handler) {
	r.Mux.Handle("DELETE "+pattern, handler)
}

// RegisterRoutes registers all enabled AuthKit routes on the given router.
func RegisterRoutes(router Router, js *AuthKit, config RouteConfig) {
	// Teams
	if js.features.Teams {
		router.Post("/teams", NewCreateTeamHandler(js))
		router.Put("/teams/{team}", NewUpdateTeamHandler(js))
		router.Delete("/teams/{team}", NewDeleteTeamHandler(js))

		router.Post("/teams/{team}/members", NewAddTeamMemberHandler(js))
		router.Put("/teams/{team}/members/{user}", NewUpdateTeamMemberRoleHandler(js))
		router.Delete("/teams/{team}/members/{user}", NewRemoveTeamMemberHandler(js))

		router.Post("/teams/{team}/invitations", NewInviteTeamMemberHandler(js))
		router.Delete("/team-invitations/{invitation}", NewCancelInvitationHandler(js))
		router.Get("/team-invitations/{invitation}/accept", NewAcceptInvitationHandler(js))

		router.Put("/current-team", NewSwitchTeamHandler(js))
	}

	// API Tokens
	if js.features.APITokens && config.Tokens != nil {
		router.Post("/user/api-tokens", NewCreateTokenHandler(js, config.Tokens))
		router.Put("/user/api-tokens/{token}", NewUpdateTokenHandler(js, config.Tokens))
		router.Delete("/user/api-tokens/{token}", NewDeleteTokenHandler(js, config.Tokens))
	}

	// Profile Photos
	if js.features.ProfilePhotos && config.Photos != nil {
		router.Put("/user/profile-photo", NewUpdateProfilePhotoHandler(js, config.Photos))
	}

	if js.features.ProfilePhotos && config.DeletePhotos != nil {
		router.Delete("/user/profile-photo", NewDeleteProfilePhotoHandler(js, config.DeletePhotos))
	}

	// Account Deletion
	if js.features.AccountDeletion && config.DeleteUser != nil {
		router.Delete("/user", NewDeleteAccountHandler(js, config.DeleteUser))
	}

	// Browser Sessions
	if js.features.BrowserSessions && config.Sessions != nil {
		router.Get("/user/sessions", NewListSessionsHandler(js, config.Sessions))
		router.Delete("/user/other-sessions", NewDeleteOtherSessionsHandler(js, config.Sessions))
	}
}
