package inception

import (
	"encoding/json"
	"net/http"
)

// InviteTeamMemberHandler handles POST /teams/{team}/invitations requests.
type InviteTeamMemberHandler struct {
	app *Inception
}

// NewInviteTeamMemberHandler creates a new invite team member handler.

// ServeHTTP handles the invite team member request.

// CancelInvitationHandler handles DELETE /team-invitations/{invitation} requests.
type CancelInvitationHandler struct {
	app *Inception
}

// NewCancelInvitationHandler creates a new cancel invitation handler.

// ServeHTTP handles the cancel invitation request.

// AcceptInvitationHandler handles GET /team-invitations/{invitation}/accept requests.
type AcceptInvitationHandler struct {
	app *Inception
}

func NewInviteTeamMemberHandler(app *Inception) *InviteTeamMemberHandler {
	return &InviteTeamMemberHandler{app: app}
}

func (h *InviteTeamMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	invitation, err := h.app.inviteMember.Invite(r.Context(), user, team, input["email"], input["role"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTeamMemberInvited, Payload: invitation})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(invitation)
}

func NewCancelInvitationHandler(app *Inception) *CancelInvitationHandler {
	return &CancelInvitationHandler{app: app}
}

func (h *CancelInvitationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	invitationID := r.PathValue("invitation")

	invitation, err := h.app.invitations.FindByID(r.Context(), invitationID)

	if err != nil || invitation == nil {
		http.Error(w, "invitation not found", http.StatusNotFound)

		return
	}

	team, err := h.app.teams.FindByID(r.Context(), invitation.TeamID)

	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)

		return
	}

	if !user.OwnsTeam(team) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	if err := h.app.invitations.Delete(r.Context(), invitationID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// NewAcceptInvitationHandler creates a new accept invitation handler.
func NewAcceptInvitationHandler(app *Inception) *AcceptInvitationHandler {
	return &AcceptInvitationHandler{app: app}
}

// ServeHTTP handles the accept invitation request.
func (h *AcceptInvitationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateTeamUser(h.app, w, r)

	if err != nil {
		http.Error(w, err.Error(), statusForError(err))

		return
	}

	invitationID := r.PathValue("invitation")

	invitation, err := h.app.invitations.FindByID(r.Context(), invitationID)

	if err != nil || invitation == nil {
		http.Error(w, "invitation not found", http.StatusNotFound)

		return
	}

	team, err := h.app.teams.FindByID(r.Context(), invitation.TeamID)

	if err != nil || team == nil {
		http.Error(w, "team not found", http.StatusNotFound)

		return
	}

	if err := h.app.addMember.Add(r.Context(), user, team, user.GetAuthIdentifier(), invitation.Role); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	if err := h.app.invitations.Delete(r.Context(), invitationID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if h.app.events != nil {
		_, _ = h.app.events.Dispatch(r.Context(), Event{Name: EventTeamMemberAdded, Payload: team})
	}

	w.WriteHeader(http.StatusOK)
}
