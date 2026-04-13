package inception

import (
	"context"
	"encoding/json"
	"net/http"
)

// currentSessionKey is the context key used by the consuming app to
// inject the current session ID into the request context.
type currentSessionKey struct{}

// WithCurrentSessionID returns a context with the current session ID.
func WithCurrentSessionID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, currentSessionKey{}, id)
}

// CurrentSessionID retrieves the current session ID from context.
func CurrentSessionID(ctx context.Context) string {
	id, _ := ctx.Value(currentSessionKey{}).(string)

	return id
}

// ListSessionsHandler handles GET /user/sessions requests.
type ListSessionsHandler struct {
	app      *Inception
	sessions SessionRepository
}

// NewListSessionsHandler creates a new list sessions handler.
func NewListSessionsHandler(app *Inception, sessions SessionRepository) *ListSessionsHandler {
	return &ListSessionsHandler{app: app, sessions: sessions}
}

// ServeHTTP handles the list sessions request.
func (h *ListSessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.BrowserSessions {
		http.Error(w, "browser sessions are disabled", http.StatusNotFound)

		return
	}

	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	sessions, err := h.sessions.FindByUser(r.Context(), user.GetAuthIdentifier())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sessions)
}

// DeleteOtherSessionsHandler handles DELETE /user/other-sessions requests.
type DeleteOtherSessionsHandler struct {
	app      *Inception
	sessions SessionRepository
}

// NewDeleteOtherSessionsHandler creates a new delete other sessions handler.
func NewDeleteOtherSessionsHandler(app *Inception, sessions SessionRepository) *DeleteOtherSessionsHandler {
	return &DeleteOtherSessionsHandler{app: app, sessions: sessions}
}

// ServeHTTP handles the delete other sessions request.
func (h *DeleteOtherSessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.app.config.Features.BrowserSessions {
		http.Error(w, "browser sessions are disabled", http.StatusNotFound)

		return
	}

	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	currentID := CurrentSessionID(r.Context())

	if currentID == "" {
		http.Error(w, "current session ID not found", http.StatusBadRequest)

		return
	}

	if err := h.sessions.DeleteOthers(r.Context(), user.GetAuthIdentifier(), currentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
