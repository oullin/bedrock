package jetstream

import (
	"errors"

	"github.com/bedrock/packages/fortify"
)

// Jetstream is the composition root for all team and profile features.
type Jetstream struct {
	features    Features
	roles       *RoleRegistry
	guard       fortify.Guard
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

// Features returns the Jetstream feature flags.
func (j *Jetstream) Features() Features { return j.features }

// Roles returns the role registry.
func (j *Jetstream) Roles() *RoleRegistry { return j.roles }

// Guard returns the authentication guard.
func (j *Jetstream) Guard() fortify.Guard { return j.guard }

// Teams returns the team repository.
func (j *Jetstream) Teams() TeamRepository { return j.teams }

// Invitations returns the invitation repository.
func (j *Jetstream) Invitations() InvitationRepository { return j.invitations }

// Events returns the event dispatcher.
func (j *Jetstream) Events() EventDispatcher { return j.events }

// CreateTeam returns the team creation action.
func (j *Jetstream) CreateTeam() CreatesTeams { return j.createTeam }

// UpdateTeam returns the team update action.
func (j *Jetstream) UpdateTeam() UpdatesTeamNames { return j.updateTeam }

// DeleteTeam returns the team deletion action.
func (j *Jetstream) DeleteTeam() DeletesTeams { return j.deleteTeam }

// AddMember returns the add member action.
func (j *Jetstream) AddMember() AddsTeamMembers { return j.addMember }

// RemoveMember returns the remove member action.
func (j *Jetstream) RemoveMember() RemovesTeamMembers { return j.removeMember }

// InviteMember returns the invite member action.
func (j *Jetstream) InviteMember() InvitesTeamMembers { return j.inviteMember }

// Builder constructs a Jetstream instance.
type Builder struct {
	js     *Jetstream
	errors []error
}

// NewBuilder creates a new Jetstream builder.
func NewBuilder() *Builder {
	return &Builder{
		js: &Jetstream{
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
func (b *Builder) WithGuard(guard fortify.Guard) *Builder {
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

// Build validates required dependencies and returns the Jetstream instance.
func (b *Builder) Build() (*Jetstream, error) {
	if b.js.guard == nil {
		b.errors = append(b.errors, errors.New("jetstream: guard is required"))
	}

	if b.js.features.Teams {
		if b.js.teams == nil {
			b.errors = append(b.errors, errors.New("jetstream: team repository is required when teams are enabled"))
		}

		if b.js.createTeam == nil {
			b.errors = append(b.errors, errors.New("jetstream: CreatesTeams action is required when teams are enabled"))
		}

		if b.js.updateTeam == nil {
			b.errors = append(b.errors, errors.New("jetstream: UpdatesTeamNames action is required when teams are enabled"))
		}

		if b.js.deleteTeam == nil {
			b.errors = append(b.errors, errors.New("jetstream: DeletesTeams action is required when teams are enabled"))
		}

		if b.js.addMember == nil {
			b.errors = append(b.errors, errors.New("jetstream: AddsTeamMembers action is required when teams are enabled"))
		}

		if b.js.removeMember == nil {
			b.errors = append(b.errors, errors.New("jetstream: RemovesTeamMembers action is required when teams are enabled"))
		}

		if b.js.inviteMember == nil {
			b.errors = append(b.errors, errors.New("jetstream: InvitesTeamMembers action is required when teams are enabled"))
		}

		if b.js.invitations == nil {
			b.errors = append(b.errors, errors.New("jetstream: invitation repository is required when teams are enabled"))
		}
	}

	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}

	return b.js, nil
}
