package inception

import (
	"context"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// CreatesNewUsers creates a new user from registration input.
type CreatesNewUsers interface {
	Create(ctx context.Context, input map[string]string) (cauth.Authenticatable, error)
}

// AuthenticatesUsers customizes the authentication attempt.
// Implement this to override the default credential-based login.
type AuthenticatesUsers interface {
	Authenticate(ctx context.Context, input map[string]string) (cauth.Authenticatable, error)
}

// UpdatesUserProfileInformation updates a user's profile.
type UpdatesUserProfileInformation interface {
	Update(ctx context.Context, user cauth.Authenticatable, input map[string]string) error
}

// UpdatesUserPasswords changes an authenticated user's password.
type UpdatesUserPasswords interface {
	Update(ctx context.Context, user cauth.Authenticatable, input map[string]string) error
}

// ResetsUserPasswords sets a new password during the reset flow.
type ResetsUserPasswords interface {
	Reset(ctx context.Context, user cauth.Authenticatable, password string) error
}

// ConfirmsPasswords validates the user's current password.
type ConfirmsPasswords interface {
	Confirm(ctx context.Context, user cauth.Authenticatable, password string) error
}

// CreatesTeams creates a new team.
type CreatesTeams interface {
	Create(ctx context.Context, user HasTeams, input map[string]string) (*Team, error)
}

// UpdatesTeamNames updates a team's name.
type UpdatesTeamNames interface {
	Update(ctx context.Context, user HasTeams, team *Team, input map[string]string) error
}

// DeletesTeams deletes a team.
type DeletesTeams interface {
	Delete(ctx context.Context, user HasTeams, team *Team) error
}

// AddsTeamMembers adds a member to a team.
type AddsTeamMembers interface {
	Add(ctx context.Context, user HasTeams, team *Team, email string, role string) error
}

// RemovesTeamMembers removes a member from a team.
type RemovesTeamMembers interface {
	Remove(ctx context.Context, user HasTeams, team *Team, memberID string) error
}

// InvitesTeamMembers sends an invitation to join a team.
type InvitesTeamMembers interface {
	Invite(ctx context.Context, user HasTeams, team *Team, email string, role string) (*TeamInvitation, error)
}
