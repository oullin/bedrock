package authkit

import "net/http"

// SwitchTeamHandler handles PUT /current-team requests.
type SwitchTeamHandler struct {
	js *AuthKit
}

// NewSwitchTeamHandler creates a new switch team handler.
func NewSwitchTeamHandler(js *AuthKit) *SwitchTeamHandler {
	return &SwitchTeamHandler{js: js}
}

// ServeHTTP handles the switch team request.
func (h *SwitchTeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.js, w, r)
	if err != nil {
		http.Error(w, err.Error(), statusForError(err))
		return
	}

	input := RequestInput(r, "team_id")
	teamID := input["team_id"]

	team, err := h.js.teams.FindByID(r.Context(), teamID)
	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)
		return
	}

	if !user.BelongsToTeam(team) && !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)
		return
	}

	user.SetCurrentTeam(team)

	if h.js.events != nil {
		_ = h.js.events.Dispatch(r.Context(), Event{Name: EventTeamSwitched, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}
