package inception

import "net/http"

// AddTeamMemberHandler handles POST /teams/{team}/members requests.
type AddTeamMemberHandler struct {
	app *Inception
}

// NewAddTeamMemberHandler creates a new add team member handler.
func NewAddTeamMemberHandler(app *Inception) *AddTeamMemberHandler {
	return &AddTeamMemberHandler{app: app}
}

// ServeHTTP handles the add team member request.
func (h *AddTeamMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if !user.HasTeamPermission(team, "addTeamMember") && !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	input := RequestInput(r, "email", "role")

	if err := h.app.addMember.Add(r.Context(), user, team, input["email"], input["role"]); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTeamMemberAdded, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateTeamMemberRoleHandler handles PUT /teams/{team}/members/{user} requests.
type UpdateTeamMemberRoleHandler struct {
	app *Inception
}

// NewUpdateTeamMemberRoleHandler creates a new update member role handler.
func NewUpdateTeamMemberRoleHandler(app *Inception) *UpdateTeamMemberRoleHandler {
	return &UpdateTeamMemberRoleHandler{app: app}
}

// ServeHTTP handles the update team member role request.
func (h *UpdateTeamMemberRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	teamID := r.PathValue("team")
	memberID := r.PathValue("user")

	team, err := h.app.teams.FindByID(r.Context(), teamID)

	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)

		return
	}

	if !user.HasTeamPermission(team, "updateTeamMember") && !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	input := RequestInput(r, "role")
	role := input["role"]

	if h.app.roles.Find(role) == nil {
		http.Error(w, "invalid role", http.StatusUnprocessableEntity)

		return
	}

	if err := h.app.teams.UpdateMemberRole(r.Context(), teamID, memberID, role); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTeamMemberUpdated, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveTeamMemberHandler handles DELETE /teams/{team}/members/{user} requests.
type RemoveTeamMemberHandler struct {
	app *Inception
}

// NewRemoveTeamMemberHandler creates a new remove team member handler.
func NewRemoveTeamMemberHandler(app *Inception) *RemoveTeamMemberHandler {
	return &RemoveTeamMemberHandler{app: app}
}

// ServeHTTP handles the remove team member request.
func (h *RemoveTeamMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	teamID := r.PathValue("team")
	memberID := r.PathValue("user")

	team, err := h.app.teams.FindByID(r.Context(), teamID)

	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)

		return
	}

	if !user.HasTeamPermission(team, "removeTeamMember") && !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	if err := h.app.removeMember.Remove(r.Context(), user, team, memberID); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTeamMemberRemoved, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}
