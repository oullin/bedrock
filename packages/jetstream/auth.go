package jetstream

import (
	"errors"
	"net/http"
)

// ErrUnauthenticated is returned when a request has no authenticated user.
var ErrUnauthenticated = errors.New("jetstream: unauthenticated")

// ErrNotTeamUser is returned when the authenticated user does not implement HasTeams.
var ErrNotTeamUser = errors.New("jetstream: user does not support teams")

// ErrUnauthorized is returned when the user lacks permission for an action.
var ErrUnauthorized = errors.New("jetstream: unauthorized")

// authenticateTeamUser resolves the authenticated user as a HasTeams from the request.
func authenticateTeamUser(j *Jetstream, w http.ResponseWriter, r *http.Request) (HasTeams, error) {
	user, err := j.guard.AuthenticateRequest(r.Context(), w, r)
	if err != nil || user == nil {
		return nil, ErrUnauthenticated
	}

	teamUser, ok := user.(HasTeams)
	if !ok {
		return nil, ErrNotTeamUser
	}

	return teamUser, nil
}
