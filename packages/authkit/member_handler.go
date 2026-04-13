package authkit

import "net/http"

// AddTeamMemberHandler handles POST /teams/{team}/members requests.
type AddTeamMemberHandler struct {
	js *AuthKit
}

// NewAddTeamMemberHandler creates a new add team member handler.

// ServeHTTP handles the add team member request.

// UpdateTeamMemberRoleHandler handles PUT /teams/{team}/members/{user} requests.
type UpdateTeamMemberRoleHandler struct {
	js *AuthKit
}

// NewUpdateTeamMemberRoleHandler creates a new update member role handler.

// ServeHTTP handles the update team member role request.

// RemoveTeamMemberHandler handles DELETE /teams/{team}/members/{user} requests.
type RemoveTeamMemberHandler struct {
	js *AuthKit
}

func NewAddTeamMemberHandler(js *AuthKit) *AddTeamMemberHandler {
	return &AddTeamMemberHandler{js: js}
}

func (h *AddTeamMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.js, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	teamID := r.PathValue("team")

	team, err := h.js.teams.FindByID(r.Context(), teamID)

	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)

		return
	}

	if !user.HasTeamPermission(team, "addTeamMember") && !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	input := RequestInput(r, "email", "role")

	if err := h.js.addMember.Add(r.Context(), user, team, input["email"], input["role"]); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.js.events != nil {
		_, _ = h.js.events.Dispatch(r.Context(), Event{Name: EventTeamMemberAdded, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}

func NewUpdateTeamMemberRoleHandler(js *AuthKit) *UpdateTeamMemberRoleHandler {
	return &UpdateTeamMemberRoleHandler{js: js}
}

func (h *UpdateTeamMemberRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.js, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	teamID := r.PathValue("team")
	memberID := r.PathValue("user")

	team, err := h.js.teams.FindByID(r.Context(), teamID)

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

	if h.js.roles.Find(role) == nil {
		http.Error(w, "invalid role", http.StatusUnprocessableEntity)

		return
	}

	if err := h.js.teams.UpdateMemberRole(r.Context(), teamID, memberID, role); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.js.events != nil {
		_, _ = h.js.events.Dispatch(r.Context(), Event{Name: EventTeamMemberUpdated, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}

// NewRemoveTeamMemberHandler creates a new remove team member handler.
func NewRemoveTeamMemberHandler(js *AuthKit) *RemoveTeamMemberHandler {
	return &RemoveTeamMemberHandler{js: js}
}

// ServeHTTP handles the remove team member request.
func (h *RemoveTeamMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.js, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	teamID := r.PathValue("team")
	memberID := r.PathValue("user")

	team, err := h.js.teams.FindByID(r.Context(), teamID)

	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)

		return
	}

	if !user.HasTeamPermission(team, "removeTeamMember") && !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	if err := h.js.removeMember.Remove(r.Context(), user, team, memberID); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.js.events != nil {
		_, _ = h.js.events.Dispatch(r.Context(), Event{Name: EventTeamMemberRemoved, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}
