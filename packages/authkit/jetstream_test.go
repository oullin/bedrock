package authkit

import (
	"context"
	"net/http"
	"strings"
	"testing"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// --- test doubles ---

type stubGuard struct{}

type stubTeamRepo struct{}

type stubInvitationRepo struct{}

type stubCreatesTeams struct{}

type stubUpdatesTeams struct{}

type stubDeletesTeams struct{}

type stubAddsMembers struct{}

type stubRemovesMembers struct{}

type stubInvitesMembers struct{}

func (g *stubGuard) Name() string { return "web" }
func (g *stubGuard) AuthenticateRequest(_ context.Context, _ http.ResponseWriter, _ *http.Request) (cauth.Authenticatable, error) {
	return nil, nil
}
func (g *stubGuard) Login(_ context.Context, _ http.ResponseWriter, _ cauth.Authenticatable, _ bool) error {
	return nil
}
func (g *stubGuard) LoginWithPendingTwoFactor(_ context.Context, _ http.ResponseWriter, _ cauth.Authenticatable) error {
	return nil
}
func (g *stubGuard) Logout(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (r *stubTeamRepo) Create(_ context.Context, _ *Team) error                   { return nil }
func (r *stubTeamRepo) FindByID(_ context.Context, _ string) (*Team, error)       { return nil, nil }
func (r *stubTeamRepo) Update(_ context.Context, _ *Team) error                   { return nil }
func (r *stubTeamRepo) Delete(_ context.Context, _ string) error                  { return nil }
func (r *stubTeamRepo) FindByOwner(_ context.Context, _ string) ([]Team, error)   { return nil, nil }
func (r *stubTeamRepo) Members(_ context.Context, _ string) ([]Membership, error) { return nil, nil }
func (r *stubTeamRepo) AddMember(_ context.Context, _ *Membership) error          { return nil }
func (r *stubTeamRepo) UpdateMemberRole(_ context.Context, _ string, _ string, _ string) error {
	return nil
}
func (r *stubTeamRepo) RemoveMember(_ context.Context, _ string, _ string) error { return nil }

func (r *stubInvitationRepo) Create(_ context.Context, _ *TeamInvitation) error { return nil }
func (r *stubInvitationRepo) FindByID(_ context.Context, _ string) (*TeamInvitation, error) {
	return nil, nil
}
func (r *stubInvitationRepo) FindByTeam(_ context.Context, _ string) ([]TeamInvitation, error) {
	return nil, nil
}
func (r *stubInvitationRepo) FindByEmail(_ context.Context, _ string, _ string) (*TeamInvitation, error) {
	return nil, nil
}
func (r *stubInvitationRepo) Delete(_ context.Context, _ string) error { return nil }

func (a *stubCreatesTeams) Create(_ context.Context, _ HasTeams, _ map[string]string) (*Team, error) {
	return &Team{}, nil
}

func (a *stubUpdatesTeams) Update(_ context.Context, _ HasTeams, _ *Team, _ map[string]string) error {
	return nil
}

func (a *stubDeletesTeams) Delete(_ context.Context, _ HasTeams, _ *Team) error { return nil }

func (a *stubAddsMembers) Add(_ context.Context, _ HasTeams, _ *Team, _ string, _ string) error {
	return nil
}

func (a *stubRemovesMembers) Remove(_ context.Context, _ HasTeams, _ *Team, _ string) error {
	return nil
}

func (a *stubInvitesMembers) Invite(_ context.Context, _ HasTeams, _ *Team, _ string, _ string) (*TeamInvitation, error) {
	return &TeamInvitation{}, nil
}

// --- tests ---

func TestBuildMinimalAuthKit(t *testing.T) {
	features := Features{}

	js, err := NewBuilder().
		WithFeatures(features).
		WithGuard(&stubGuard{}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if js.Guard().Name() != "web" {
		t.Fatal("expected guard name web")
	}
}

func TestBuildFullAuthKit(t *testing.T) {
	js, err := NewBuilder().
		WithFeatures(DefaultFeatures()).
		WithGuard(&stubGuard{}).
		WithTeams(&stubTeamRepo{}).
		WithInvitations(&stubInvitationRepo{}).
		WithCreateTeam(&stubCreatesTeams{}).
		WithUpdateTeam(&stubUpdatesTeams{}).
		WithDeleteTeam(&stubDeletesTeams{}).
		WithAddMember(&stubAddsMembers{}).
		WithRemoveMember(&stubRemovesMembers{}).
		WithInviteMember(&stubInvitesMembers{}).
		WithRoles(func(rr *RoleRegistry) {
			rr.Define("admin", "Administrator", "Full access", []string{"*"})
			rr.Define("member", "Member", "Basic access", []string{"read"})
			rr.SetDefault("member")
		}).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(js.Roles().All()) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(js.Roles().All()))
	}

	if js.Roles().Default() != "member" {
		t.Fatalf("expected default role member, got %s", js.Roles().Default())
	}
}

func TestBuildFailsWithoutGuard(t *testing.T) {
	_, err := NewBuilder().
		WithFeatures(Features{}).
		Build()

	if err == nil {
		t.Fatal("expected error for missing guard")
	}

	if !strings.Contains(err.Error(), "guard is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildFailsWithoutTeamDeps(t *testing.T) {
	_, err := NewBuilder().
		WithFeatures(Features{Teams: true}).
		WithGuard(&stubGuard{}).
		Build()

	if err == nil {
		t.Fatal("expected error for missing team dependencies")
	}

	if !strings.Contains(err.Error(), "team repository") {
		t.Fatalf("expected team repository error, got: %v", err)
	}
}

func TestDefaultFeaturesAllEnabled(t *testing.T) {
	f := DefaultFeatures()

	if !f.Teams || !f.APITokens || !f.ProfilePhotos || !f.AccountDeletion || !f.BrowserSessions {
		t.Fatal("expected all features enabled by default")
	}
}
