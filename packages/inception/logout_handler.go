package inception

import "net/http"

// LogoutHandler handles POST /logout requests.
type LogoutHandler struct {
	app *Inception
}

// NewLogoutHandler creates a new logout handler.
func NewLogoutHandler(app *Inception) *LogoutHandler {
	return &LogoutHandler{app: app}
}

// ServeHTTP handles the logout request.
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, _ := h.app.guard.AuthenticateRequest(ctx, w, r)

	if err := h.app.guard.Logout(ctx, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if h.app.events != nil && user != nil {
		_, _ = h.app.events.Dispatch(ctx, LoggedOutPayload{User: user})
	}

	h.app.responder.LogoutResponse(w, r)
}
