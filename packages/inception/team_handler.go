package inception

import (
	"encoding/json"
	"net/http"
)

// CreateTeamHandler handles POST /teams requests.
type CreateTeamHandler struct {
	app *Inception
}

// NewCreateTeamHandler creates a new create team handler.
func NewCreateTeamHandler(app *Inception) *CreateTeamHandler {
	return &CreateTeamHandler{app: app}
}

// ServeHTTP handles the create team request.
func (h *CreateTeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	input := RequestInput(r, "name")

	team, err := h.app.createTeam.Create(r.Context(), user, input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTeamCreated, Payload: team})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(team)
}

// UpdateTeamHandler handles PUT /teams/{team} requests.
type UpdateTeamHandler struct {
	app *Inception
}

// NewUpdateTeamHandler creates a new update team handler.
func NewUpdateTeamHandler(app *Inception) *UpdateTeamHandler {
	return &UpdateTeamHandler{app: app}
}

// ServeHTTP handles the update team request.
func (h *UpdateTeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	teamID := r.PathValue("team")

	team, err := h.app.teams.FindByID(r.Context(), teamID)

	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)

		return
	}

	if !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	input := RequestInput(r, "name")

	if err := h.app.updateTeam.Update(r.Context(), user, team, input); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTeamUpdated, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteTeamHandler handles DELETE /teams/{team} requests.
type DeleteTeamHandler struct {
	app *Inception
}

// NewDeleteTeamHandler creates a new delete team handler.
func NewDeleteTeamHandler(app *Inception) *DeleteTeamHandler {
	return &DeleteTeamHandler{app: app}
}

// ServeHTTP handles the delete team request.
func (h *DeleteTeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	teamID := r.PathValue("team")

	team, err := h.app.teams.FindByID(r.Context(), teamID)

	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)

		return
	}

	if !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	if team.PersonalTeam {
		http.Error(w, "cannot delete personal team", http.StatusUnprocessableEntity)

		return
	}

	if err := h.app.deleteTeam.Delete(r.Context(), user, team); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTeamDeleted, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}

func statusForError(err error) int {
	switch err {
	case ErrUnauthenticated:
		return http.StatusUnauthorized
	case ErrNotTeamUser:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
