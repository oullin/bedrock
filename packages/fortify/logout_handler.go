package fortify

import "net/http"

// LogoutHandler handles POST /logout requests.
type LogoutHandler struct {
	fortify *Fortify
}

// NewLogoutHandler creates a new logout handler.
func NewLogoutHandler(f *Fortify) *LogoutHandler {
	return &LogoutHandler{fortify: f}
}

// ServeHTTP handles the logout request.
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, _ := h.fortify.guard.AuthenticateRequest(ctx, w, r)

	if err := h.fortify.guard.Logout(ctx, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if h.fortify.events != nil && user != nil {
		_ = h.fortify.events.Dispatch(ctx, Event{
			Name:    EventLoggedOut,
			Payload: LoggedOutPayload{User: user},
		})
	}

	h.fortify.responder.LogoutResponse(w, r)
}
