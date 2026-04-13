package inception

import (
	"errors"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/contracts/events"
)

// Inception is the composition root for all authentication, team,
// and profile features. It holds references to contracts and action
// implementations provided by the consuming application.
type Inception struct {
	config    Config
	guard     cauth.HTTPGuard
	provider  cauth.UserProvider
	hasher    cauth.PasswordHasher
	broker    PasswordBroker
	verifier  EmailVerifier
	events    events.Dispatcher
	limiter   RateLimiter
	responder Responder

	createUser    CreatesNewUsers
	authenticator AuthenticatesUsers
	updateProfile UpdatesUserProfileInformation
	updatePass    UpdatesUserPasswords
	resetPass     ResetsUserPasswords
	confirmPass   ConfirmsPasswords

	roles       *RoleRegistry
	teams       TeamRepository
	invitations InvitationRepository

	createTeam   CreatesTeams
	updateTeam   UpdatesTeamNames
	deleteTeam   DeletesTeams
	addMember    AddsTeamMembers
	removeMember RemovesTeamMembers
	inviteMember InvitesTeamMembers
}

// Builder constructs an Inception instance with required dependencies.
type Builder struct {
	app    *Inception
	errors []error
}

func (i *Inception) Config() Config                               { return i.config }
func (i *Inception) Guard() cauth.HTTPGuard                       { return i.guard }
func (i *Inception) Provider() cauth.UserProvider                 { return i.provider }
func (i *Inception) Hasher() cauth.PasswordHasher                 { return i.hasher }
func (i *Inception) Broker() PasswordBroker                       { return i.broker }
func (i *Inception) Verifier() EmailVerifier                      { return i.verifier }
func (i *Inception) Events() events.Dispatcher                    { return i.events }
func (i *Inception) Limiter() RateLimiter                         { return i.limiter }
func (i *Inception) Responder() Responder                         { return i.responder }
func (i *Inception) CreateUser() CreatesNewUsers                  { return i.createUser }
func (i *Inception) Authenticator() AuthenticatesUsers            { return i.authenticator }
func (i *Inception) UpdateProfile() UpdatesUserProfileInformation { return i.updateProfile }
func (i *Inception) UpdatePassword() UpdatesUserPasswords         { return i.updatePass }
func (i *Inception) ResetPassword() ResetsUserPasswords           { return i.resetPass }
func (i *Inception) ConfirmPassword() ConfirmsPasswords           { return i.confirmPass }
func (i *Inception) Roles() *RoleRegistry                         { return i.roles }
func (i *Inception) Teams() TeamRepository                        { return i.teams }
func (i *Inception) Invitations() InvitationRepository            { return i.invitations }
func (i *Inception) CreateTeam() CreatesTeams                     { return i.createTeam }
func (i *Inception) UpdateTeam() UpdatesTeamNames                 { return i.updateTeam }
func (i *Inception) DeleteTeam() DeletesTeams                     { return i.deleteTeam }
func (i *Inception) AddMember() AddsTeamMembers                   { return i.addMember }
func (i *Inception) RemoveMember() RemovesTeamMembers             { return i.removeMember }
func (i *Inception) InviteMember() InvitesTeamMembers             { return i.inviteMember }
func (i *Inception) Features() Features                           { return i.config.Features }

// NewBuilder creates a new Inception builder.
func NewBuilder() *Builder {
	return &Builder{
		app: &Inception{
			config: DefaultConfig(),
			roles:  NewRoleRegistry(),
		},
	}
}

func (b *Builder) WithConfig(config Config) *Builder {
	b.app.config = config

	return b
}

func (b *Builder) WithGuard(guard cauth.HTTPGuard) *Builder {
	b.app.guard = guard

	return b
}

func (b *Builder) WithProvider(provider cauth.UserProvider) *Builder {
	b.app.provider = provider

	return b
}

func (b *Builder) WithHasher(hasher cauth.PasswordHasher) *Builder {
	b.app.hasher = hasher

	return b
}

func (b *Builder) WithBroker(broker PasswordBroker) *Builder {
	b.app.broker = broker

	return b
}

func (b *Builder) WithVerifier(verifier EmailVerifier) *Builder {
	b.app.verifier = verifier

	return b
}

func (b *Builder) WithEvents(dispatcher events.Dispatcher) *Builder {
	b.app.events = dispatcher

	return b
}

func (b *Builder) WithLimiter(limiter RateLimiter) *Builder {
	b.app.limiter = limiter

	return b
}

func (b *Builder) WithResponder(responder Responder) *Builder {
	b.app.responder = responder

	return b
}

func (b *Builder) WithCreateUser(action CreatesNewUsers) *Builder {
	b.app.createUser = action

	return b
}

func (b *Builder) WithAuthenticator(action AuthenticatesUsers) *Builder {
	b.app.authenticator = action

	return b
}

func (b *Builder) WithUpdateProfile(action UpdatesUserProfileInformation) *Builder {
	b.app.updateProfile = action

	return b
}

func (b *Builder) WithUpdatePassword(action UpdatesUserPasswords) *Builder {
	b.app.updatePass = action

	return b
}

func (b *Builder) WithResetPassword(action ResetsUserPasswords) *Builder {
	b.app.resetPass = action

	return b
}

func (b *Builder) WithConfirmPassword(action ConfirmsPasswords) *Builder {
	b.app.confirmPass = action

	return b
}

func (b *Builder) WithTeams(teams TeamRepository) *Builder {
	b.app.teams = teams

	return b
}

func (b *Builder) WithInvitations(invitations InvitationRepository) *Builder {
	b.app.invitations = invitations

	return b
}

func (b *Builder) WithRoles(configure func(*RoleRegistry)) *Builder {
	configure(b.app.roles)

	return b
}

func (b *Builder) WithCreateTeam(action CreatesTeams) *Builder {
	b.app.createTeam = action

	return b
}

func (b *Builder) WithUpdateTeam(action UpdatesTeamNames) *Builder {
	b.app.updateTeam = action

	return b
}

func (b *Builder) WithDeleteTeam(action DeletesTeams) *Builder {
	b.app.deleteTeam = action

	return b
}

func (b *Builder) WithAddMember(action AddsTeamMembers) *Builder {
	b.app.addMember = action

	return b
}

func (b *Builder) WithRemoveMember(action RemovesTeamMembers) *Builder {
	b.app.removeMember = action

	return b
}

func (b *Builder) WithInviteMember(action InvitesTeamMembers) *Builder {
	b.app.inviteMember = action

	return b
}

// Build validates required dependencies and returns the Inception instance.
func (b *Builder) Build() (*Inception, error) {
	if b.app.guard == nil {
		b.errors = append(b.errors, errors.New("inception: guard is required"))
	}

	if b.app.provider == nil {
		b.errors = append(b.errors, errors.New("inception: user provider is required"))
	}

	if b.app.hasher == nil {
		b.errors = append(b.errors, errors.New("inception: password hasher is required"))
	}

	if b.app.responder == nil {
		b.errors = append(b.errors, errors.New("inception: responder is required"))
	}

	if b.app.config.Features.Registration && b.app.createUser == nil {
		b.errors = append(b.errors, errors.New("inception: CreatesNewUsers action is required when registration is enabled"))
	}

	if b.app.config.Features.UpdateProfileInformation && b.app.updateProfile == nil {
		b.errors = append(b.errors, errors.New("inception: UpdatesUserProfileInformation action is required when profile updates are enabled"))
	}

	if b.app.config.Features.UpdatePasswords && b.app.updatePass == nil {
		b.errors = append(b.errors, errors.New("inception: UpdatesUserPasswords action is required when password updates are enabled"))
	}

	if b.app.config.Features.ResetPasswords && b.app.broker == nil {
		b.errors = append(b.errors, errors.New("inception: password broker is required when password resets are enabled"))
	}

	if b.app.config.Features.ResetPasswords && b.app.resetPass == nil {
		b.errors = append(b.errors, errors.New("inception: ResetsUserPasswords action is required when password resets are enabled"))
	}

	if b.app.config.Features.EmailVerification && b.app.verifier == nil {
		b.errors = append(b.errors, errors.New("inception: email verifier is required when email verification is enabled"))
	}

	if b.app.config.Features.ConfirmPassword && b.app.confirmPass == nil {
		b.errors = append(b.errors, errors.New("inception: ConfirmsPasswords action is required when password confirmation is enabled"))
	}

	if b.app.config.Features.Teams {
		if b.app.teams == nil {
			b.errors = append(b.errors, errors.New("inception: team repository is required when teams are enabled"))
		}

		if b.app.createTeam == nil {
			b.errors = append(b.errors, errors.New("inception: CreatesTeams action is required when teams are enabled"))
		}

		if b.app.updateTeam == nil {
			b.errors = append(b.errors, errors.New("inception: UpdatesTeamNames action is required when teams are enabled"))
		}

		if b.app.deleteTeam == nil {
			b.errors = append(b.errors, errors.New("inception: DeletesTeams action is required when teams are enabled"))
		}

		if b.app.addMember == nil {
			b.errors = append(b.errors, errors.New("inception: AddsTeamMembers action is required when teams are enabled"))
		}

		if b.app.removeMember == nil {
			b.errors = append(b.errors, errors.New("inception: RemovesTeamMembers action is required when teams are enabled"))
		}

		if b.app.inviteMember == nil {
			b.errors = append(b.errors, errors.New("inception: InvitesTeamMembers action is required when teams are enabled"))
		}

		if b.app.invitations == nil {
			b.errors = append(b.errors, errors.New("inception: invitation repository is required when teams are enabled"))
		}
	}

	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}

	return b.app, nil
}
