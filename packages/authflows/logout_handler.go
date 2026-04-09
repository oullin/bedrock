package authflows

import "net/http"

// LogoutHandler handles POST /logout requests.
type LogoutHandler struct {
	authflows *AuthFlows
}

// NewLogoutHandler creates a new logout handler.
func NewLogoutHandler(f *AuthFlows) *LogoutHandler {
	return &LogoutHandler{authflows: f}
}

// ServeHTTP handles the logout request.
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, _ := h.authflows.guard.AuthenticateRequest(ctx, w, r)

	if err := h.authflows.guard.Logout(ctx, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.authflows.events != nil && user != nil {
		_ = h.authflows.events.Dispatch(ctx, Event{
			Name:    EventLoggedOut,
			Payload: LoggedOutPayload{User: user},
		})
	}

	h.authflows.responder.LogoutResponse(w, r)
}
