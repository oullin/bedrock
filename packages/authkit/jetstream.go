package authkit

import (
	"errors"

	"github.com/bedrock/packages/authflows"
)

// AuthKit is the composition root for all team and profile features.
type AuthKit struct {
	features    Features
	roles       *RoleRegistry
	guard       authflows.Guard
	teams       TeamRepository
	invitations InvitationRepository
	events      EventDispatcher

	createTeam   CreatesTeams
	updateTeam   UpdatesTeamNames
	deleteTeam   DeletesTeams
	addMember    AddsTeamMembers
	removeMember RemovesTeamMembers
	inviteMember InvitesTeamMembers
}

// Features returns the AuthKit feature flags.

// Roles returns the role registry.

// Guard returns the authentication guard.

// Teams returns the team repository.

// Invitations returns the invitation repository.

// Events returns the event dispatcher.

// CreateTeam returns the team creation action.

// UpdateTeam returns the team update action.

// DeleteTeam returns the team deletion action.

// AddMember returns the add member action.

// RemoveMember returns the remove member action.

// InviteMember returns the invite member action.

// Builder constructs a AuthKit instance.
type Builder struct {
	js     *AuthKit
	errors []error
}

func (j *AuthKit) Features() Features { return j.features }

func (j *AuthKit) Roles() *RoleRegistry { return j.roles }

func (j *AuthKit) Guard() authflows.Guard { return j.guard }

func (j *AuthKit) Teams() TeamRepository { return j.teams }

func (j *AuthKit) Invitations() InvitationRepository { return j.invitations }

func (j *AuthKit) Events() EventDispatcher { return j.events }

func (j *AuthKit) CreateTeam() CreatesTeams { return j.createTeam }

func (j *AuthKit) UpdateTeam() UpdatesTeamNames { return j.updateTeam }

func (j *AuthKit) DeleteTeam() DeletesTeams { return j.deleteTeam }

func (j *AuthKit) AddMember() AddsTeamMembers { return j.addMember }

func (j *AuthKit) RemoveMember() RemovesTeamMembers { return j.removeMember }

func (j *AuthKit) InviteMember() InvitesTeamMembers { return j.inviteMember }

// NewBuilder creates a new AuthKit builder.
func NewBuilder() *Builder {
	return &Builder{
		js: &AuthKit{
			features: DefaultFeatures(),
			roles:    NewRoleRegistry(),
		},
	}
}

// WithFeatures sets the feature flags.
func (b *Builder) WithFeatures(features Features) *Builder {
	b.js.features = features

	return b
}

// WithGuard sets the authentication guard.
func (b *Builder) WithGuard(guard authflows.Guard) *Builder {
	b.js.guard = guard

	return b
}

// WithTeams sets the team repository.
func (b *Builder) WithTeams(teams TeamRepository) *Builder {
	b.js.teams = teams

	return b
}

// WithInvitations sets the invitation repository.
func (b *Builder) WithInvitations(invitations InvitationRepository) *Builder {
	b.js.invitations = invitations

	return b
}

// WithEvents sets the event dispatcher.
func (b *Builder) WithEvents(events EventDispatcher) *Builder {
	b.js.events = events

	return b
}

// WithRoles provides access to the role registry for configuration.
func (b *Builder) WithRoles(configure func(*RoleRegistry)) *Builder {
	configure(b.js.roles)

	return b
}

// WithCreateTeam sets the team creation action.
func (b *Builder) WithCreateTeam(action CreatesTeams) *Builder {
	b.js.createTeam = action

	return b
}

// WithUpdateTeam sets the team update action.
func (b *Builder) WithUpdateTeam(action UpdatesTeamNames) *Builder {
	b.js.updateTeam = action

	return b
}

// WithDeleteTeam sets the team deletion action.
func (b *Builder) WithDeleteTeam(action DeletesTeams) *Builder {
	b.js.deleteTeam = action

	return b
}

// WithAddMember sets the add member action.
func (b *Builder) WithAddMember(action AddsTeamMembers) *Builder {
	b.js.addMember = action

	return b
}

// WithRemoveMember sets the remove member action.
func (b *Builder) WithRemoveMember(action RemovesTeamMembers) *Builder {
	b.js.removeMember = action

	return b
}

// WithInviteMember sets the invite member action.
func (b *Builder) WithInviteMember(action InvitesTeamMembers) *Builder {
	b.js.inviteMember = action

	return b
}

// Build validates required dependencies and returns the AuthKit instance.
func (b *Builder) Build() (*AuthKit, error) {
	if b.js.guard == nil {
		b.errors = append(b.errors, errors.New("authkit: guard is required"))
	}

	if b.js.features.Teams {
		if b.js.teams == nil {
			b.errors = append(b.errors, errors.New("authkit: team repository is required when teams are enabled"))
		}

		if b.js.createTeam == nil {
			b.errors = append(b.errors, errors.New("authkit: CreatesTeams action is required when teams are enabled"))
		}

		if b.js.updateTeam == nil {
			b.errors = append(b.errors, errors.New("authkit: UpdatesTeamNames action is required when teams are enabled"))
		}

		if b.js.deleteTeam == nil {
			b.errors = append(b.errors, errors.New("authkit: DeletesTeams action is required when teams are enabled"))
		}

		if b.js.addMember == nil {
			b.errors = append(b.errors, errors.New("authkit: AddsTeamMembers action is required when teams are enabled"))
		}

		if b.js.removeMember == nil {
			b.errors = append(b.errors, errors.New("authkit: RemovesTeamMembers action is required when teams are enabled"))
		}

		if b.js.inviteMember == nil {
			b.errors = append(b.errors, errors.New("authkit: InvitesTeamMembers action is required when teams are enabled"))
		}

		if b.js.invitations == nil {
			b.errors = append(b.errors, errors.New("authkit: invitation repository is required when teams are enabled"))
		}
	}

	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}

	return b.js, nil
}
