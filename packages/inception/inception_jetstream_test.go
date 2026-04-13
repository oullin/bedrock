package inception

import (
	"context"
	"net/http"
	"strings"
	"testing"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// --- test doubles (jetstream) ---

type jsStubGuard struct{}

type jsStubTeamRepo struct{}

type jsStubInvitationRepo struct{}

type jsStubCreatesTeams struct{}

type jsStubUpdatesTeams struct{}

type jsStubDeletesTeams struct{}

type jsStubAddsMembers struct{}

type jsStubRemovesMembers struct{}

type jsStubInvitesMembers struct{}

func (g *jsStubGuard) Name() string { return "web" }
func (g *jsStubGuard) AuthenticateRequest(_ context.Context, _ http.ResponseWriter, _ *http.Request) (cauth.Authenticatable, error) {
	return nil, nil
}
func (g *jsStubGuard) Login(_ context.Context, _ http.ResponseWriter, _ cauth.Authenticatable, _ bool) error {
	return nil
}
func (g *jsStubGuard) LoginWithPendingTwoFactor(_ context.Context, _ http.ResponseWriter, _ cauth.Authenticatable) error {
	return nil
}
func (g *jsStubGuard) Logout(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (r *jsStubTeamRepo) Create(_ context.Context, _ *Team) error                   { return nil }
func (r *jsStubTeamRepo) FindByID(_ context.Context, _ string) (*Team, error)       { return nil, nil }
func (r *jsStubTeamRepo) Update(_ context.Context, _ *Team) error                   { return nil }
func (r *jsStubTeamRepo) Delete(_ context.Context, _ string) error                  { return nil }
func (r *jsStubTeamRepo) FindByOwner(_ context.Context, _ string) ([]Team, error)   { return nil, nil }
func (r *jsStubTeamRepo) Members(_ context.Context, _ string) ([]Membership, error) { return nil, nil }
func (r *jsStubTeamRepo) AddMember(_ context.Context, _ *Membership) error          { return nil }
func (r *jsStubTeamRepo) UpdateMemberRole(_ context.Context, _ string, _ string, _ string) error {
	return nil
}
func (r *jsStubTeamRepo) RemoveMember(_ context.Context, _ string, _ string) error { return nil }

func (r *jsStubInvitationRepo) Create(_ context.Context, _ *TeamInvitation) error { return nil }
func (r *jsStubInvitationRepo) FindByID(_ context.Context, _ string) (*TeamInvitation, error) {
	return nil, nil
}
func (r *jsStubInvitationRepo) FindByTeam(_ context.Context, _ string) ([]TeamInvitation, error) {
	return nil, nil
}
func (r *jsStubInvitationRepo) FindByEmail(_ context.Context, _ string, _ string) (*TeamInvitation, error) {
	return nil, nil
}
func (r *jsStubInvitationRepo) Delete(_ context.Context, _ string) error { return nil }

func (a *jsStubCreatesTeams) Create(_ context.Context, _ HasTeams, _ map[string]string) (*Team, error) {
	return &Team{}, nil
}

func (a *jsStubUpdatesTeams) Update(_ context.Context, _ HasTeams, _ *Team, _ map[string]string) error {
	return nil
}

func (a *jsStubDeletesTeams) Delete(_ context.Context, _ HasTeams, _ *Team) error { return nil }

func (a *jsStubAddsMembers) Add(_ context.Context, _ HasTeams, _ *Team, _ string, _ string) error {
	return nil
}

func (a *jsStubRemovesMembers) Remove(_ context.Context, _ HasTeams, _ *Team, _ string) error {
	return nil
}

func (a *jsStubInvitesMembers) Invite(_ context.Context, _ HasTeams, _ *Team, _ string, _ string) (*TeamInvitation, error) {
	return &TeamInvitation{}, nil
}

// --- tests ---

func TestBuildMinimalInceptionJS(t *testing.T) {
	config := DefaultConfig()
	config.Features = Features{}

	app, err := NewBuilder().
		WithConfig(config).
		WithGuard(&jsStubGuard{}).
		WithProvider(&stubProvider{}).
		WithHasher(&stubHasher{}).
		WithResponder(&stubResponder{}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.Guard().Name() != "web" {
		t.Fatal("expected guard name web")
	}
}

func TestBuildFullInceptionJS(t *testing.T) {
	config := DefaultConfig()
	config.Features = Features{
		Teams:           true,
		APITokens:       true,
		ProfilePhotos:   true,
		AccountDeletion: true,
		BrowserSessions: true,
	}

	app, err := NewBuilder().
		WithConfig(config).
		WithGuard(&jsStubGuard{}).
		WithProvider(&stubProvider{}).
		WithHasher(&stubHasher{}).
		WithResponder(&stubResponder{}).
		WithTeams(&jsStubTeamRepo{}).
		WithInvitations(&jsStubInvitationRepo{}).
		WithCreateTeam(&jsStubCreatesTeams{}).
		WithUpdateTeam(&jsStubUpdatesTeams{}).
		WithDeleteTeam(&jsStubDeletesTeams{}).
		WithAddMember(&jsStubAddsMembers{}).
		WithRemoveMember(&jsStubRemovesMembers{}).
		WithInviteMember(&jsStubInvitesMembers{}).
		WithRoles(func(rr *RoleRegistry) {
			rr.Define("admin", "Administrator", "Full access", []string{"*"})
			rr.Define("member", "Member", "Basic access", []string{"read"})
			rr.SetDefault("member")
		}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(app.Roles().All()) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(app.Roles().All()))
	}

	if app.Roles().Default() != "member" {
		t.Fatalf("expected default role member, got %s", app.Roles().Default())
	}
}

func TestBuildFailsWithoutGuardJS(t *testing.T) {
	config := DefaultConfig()
	config.Features = Features{}

	_, err := NewBuilder().
		WithConfig(config).
		WithProvider(&stubProvider{}).
		WithHasher(&stubHasher{}).
		WithResponder(&stubResponder{}).
		Build()

	if err == nil {
		t.Fatal("expected error for missing guard")
	}

	if !strings.Contains(err.Error(), "guard is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildFailsWithoutTeamDepsJS(t *testing.T) {
	config := DefaultConfig()
	config.Features = Features{Teams: true}

	_, err := NewBuilder().
		WithConfig(config).
		WithGuard(&jsStubGuard{}).
		WithProvider(&stubProvider{}).
		WithHasher(&stubHasher{}).
		WithResponder(&stubResponder{}).
		Build()

	if err == nil {
		t.Fatal("expected error for missing team dependencies")
	}

	if !strings.Contains(err.Error(), "team repository") {
		t.Fatalf("expected team repository error, got: %v", err)
	}
}

func TestDefaultFeaturesAllEnabledJS(t *testing.T) {
	f := DefaultFeatures()

	if !f.Teams || !f.APITokens || !f.ProfilePhotos || !f.AccountDeletion || !f.BrowserSessions {
		t.Fatal("expected all features enabled by default")
	}
}
