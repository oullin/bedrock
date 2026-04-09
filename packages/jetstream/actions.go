package jetstream

import "context"

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
