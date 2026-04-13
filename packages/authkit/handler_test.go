package authkit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/contracts/events"
)

// --- test team user ---

type testTeamUser struct {
	id          string
	teams       []Team
	currentTeam *Team
	ownedTeamID string
	permissions map[string]bool
}

// --- test guard ---

type testGuard struct {
	user cauth.Authenticatable
}

// --- test repos ---

type testTeamRepo struct {
	teams   map[string]*Team
	created bool
	deleted bool
}

type testInvitationRepo struct {
	invitations map[string]*TeamInvitation
	deleted     bool
}

// --- test events ---

type testEvents struct {
	dispatched []any
}

// --- test actions ---

type testCreatesTeams struct {
	created bool
}

type testUpdatesTeams struct{ called bool }

type testDeletesTeams struct{ called bool }

type testAddsMembers struct{ called bool }

type testRemovesMembers struct{ called bool }

type testInvitesMembers struct{ called bool }

func (u *testTeamUser) GetAuthIdentifierName() string { return "id" }
func (u *testTeamUser) GetAuthIdentifier() string     { return u.id }
func (u *testTeamUser) GetAuthPasswordName() string   { return "password" }
func (u *testTeamUser) GetAuthPassword() string       { return "" }
func (u *testTeamUser) SetAuthPassword(_ string)      {}
func (u *testTeamUser) GetRememberToken() string      { return "" }
func (u *testTeamUser) SetRememberToken(_ string)     {}
func (u *testTeamUser) GetRememberTokenName() string  { return "remember_token" }
func (u *testTeamUser) Teams() []Team                 { return u.teams }
func (u *testTeamUser) CurrentTeam() *Team            { return u.currentTeam }
func (u *testTeamUser) SetCurrentTeam(t *Team)        { u.currentTeam = t }
func (u *testTeamUser) OwnsTeam(t *Team) bool         { return t.OwnerID == u.id }
func (u *testTeamUser) BelongsToTeam(t *Team) bool {
	for _, team := range u.teams {
		if team.ID == t.ID {
			return true
		}
	}

	return false
}
func (u *testTeamUser) TeamRole(_ *Team) string { return "admin" }
func (u *testTeamUser) HasTeamPermission(_ *Team, perm string) bool {
	if u.permissions == nil {
		return true
	}

	return u.permissions[perm]
}

func (g *testGuard) Name() string { return "test" }
func (g *testGuard) AuthenticateRequest(_ context.Context, _ http.ResponseWriter, _ *http.Request) (cauth.Authenticatable, error) {
	return g.user, nil
}
func (g *testGuard) Login(_ context.Context, _ http.ResponseWriter, _ cauth.Authenticatable, _ bool) error {
	return nil
}
func (g *testGuard) LoginWithPendingTwoFactor(_ context.Context, _ http.ResponseWriter, _ cauth.Authenticatable) error {
	return nil
}
func (g *testGuard) Logout(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newTestTeamRepo(teams ...*Team) *testTeamRepo {
	r := &testTeamRepo{teams: make(map[string]*Team)}

	for _, t := range teams {
		r.teams[t.ID] = t
	}

	return r
}

func (r *testTeamRepo) Create(_ context.Context, t *Team) error {
	r.created = true
	r.teams[t.ID] = t

	return nil
}
func (r *testTeamRepo) FindByID(_ context.Context, id string) (*Team, error) {
	return r.teams[id], nil
}
func (r *testTeamRepo) Update(_ context.Context, _ *Team) error { return nil }
func (r *testTeamRepo) Delete(_ context.Context, id string) error {
	r.deleted = true
	delete(r.teams, id)

	return nil
}
func (r *testTeamRepo) FindByOwner(_ context.Context, _ string) ([]Team, error)   { return nil, nil }
func (r *testTeamRepo) Members(_ context.Context, _ string) ([]Membership, error) { return nil, nil }
func (r *testTeamRepo) AddMember(_ context.Context, _ *Membership) error          { return nil }
func (r *testTeamRepo) UpdateMemberRole(_ context.Context, _ string, _ string, _ string) error {
	return nil
}
func (r *testTeamRepo) RemoveMember(_ context.Context, _ string, _ string) error { return nil }

func newTestInvitationRepo(invs ...*TeamInvitation) *testInvitationRepo {
	r := &testInvitationRepo{invitations: make(map[string]*TeamInvitation)}

	for _, inv := range invs {
		r.invitations[inv.ID] = inv
	}

	return r
}

func (r *testInvitationRepo) Create(_ context.Context, inv *TeamInvitation) error {
	r.invitations[inv.ID] = inv

	return nil
}
func (r *testInvitationRepo) FindByID(_ context.Context, id string) (*TeamInvitation, error) {
	return r.invitations[id], nil
}
func (r *testInvitationRepo) FindByTeam(_ context.Context, _ string) ([]TeamInvitation, error) {
	return nil, nil
}
func (r *testInvitationRepo) FindByEmail(_ context.Context, _ string, _ string) (*TeamInvitation, error) {
	return nil, nil
}
func (r *testInvitationRepo) Delete(_ context.Context, id string) error {
	r.deleted = true
	delete(r.invitations, id)

	return nil
}

func (e *testEvents) Listen(_ any, _ ...events.Listener)          {}
func (e *testEvents) HasListeners(_ any) bool                     { return false }
func (e *testEvents) HasWildcardListeners(_ any) bool             { return false }
func (e *testEvents) Subscribe(_ events.Subscriber)               {}
func (e *testEvents) Until(_ context.Context, _ any) (any, error) { return nil, nil }
func (e *testEvents) Push(_ context.Context, _ any)               {}
func (e *testEvents) Flush(_ context.Context, _ string) error     { return nil }
func (e *testEvents) Forget(_ any)                                {}
func (e *testEvents) ForgetPushed()                               {}
func (e *testEvents) GetListeners(_ any) []events.Listener        { return nil }

func (e *testEvents) Dispatch(_ context.Context, event any) ([]any, error) {
	e.dispatched = append(e.dispatched, event)

	return nil, nil
}

func (a *testCreatesTeams) Create(_ context.Context, _ HasTeams, input map[string]string) (*Team, error) {
	a.created = true

	return &Team{ID: "new-team", Name: input["name"], OwnerID: "1"}, nil
}

func (a *testUpdatesTeams) Update(_ context.Context, _ HasTeams, _ *Team, _ map[string]string) error {
	a.called = true

	return nil
}

func (a *testDeletesTeams) Delete(_ context.Context, _ HasTeams, _ *Team) error {
	a.called = true

	return nil
}

func (a *testAddsMembers) Add(_ context.Context, _ HasTeams, _ *Team, _ string, _ string) error {
	a.called = true

	return nil
}

func (a *testRemovesMembers) Remove(_ context.Context, _ HasTeams, _ *Team, _ string) error {
	a.called = true

	return nil
}

func (a *testInvitesMembers) Invite(_ context.Context, _ HasTeams, _ *Team, email string, role string) (*TeamInvitation, error) {
	a.called = true

	return &TeamInvitation{ID: "inv-1", Email: email, Role: role}, nil
}

// --- helpers ---

func buildTestAuthKit(user *testTeamUser, teams *testTeamRepo, invitations *testInvitationRepo) (*AuthKit, *testEvents) {
	events := &testEvents{}
	roles := NewRoleRegistry()
	roles.Define("admin", "Admin", "", []string{"*"})
	roles.Define("member", "Member", "", []string{"read"})
	roles.SetDefault("member")

	return &AuthKit{
		features:     DefaultFeatures(),
		roles:        roles,
		guard:        &testGuard{user: user},
		teams:        teams,
		invitations:  invitations,
		events:       events,
		createTeam:   &testCreatesTeams{},
		updateTeam:   &testUpdatesTeams{},
		deleteTeam:   &testDeletesTeams{},
		addMember:    &testAddsMembers{},
		removeMember: &testRemovesMembers{},
		inviteMember: &testInvitesMembers{},
	}, events
}

func postForm(path string, body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

func putForm(path string, body string) *http.Request {
	req := httptest.NewRequest(http.MethodPut, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

// --- create team tests ---

func TestCreateTeamSuccess(t *testing.T) {
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo()
	js, events := buildTestAuthKit(user, teams, newTestInvitationRepo())

	handler := NewCreateTeamHandler(js)
	w := httptest.NewRecorder()
	r := postForm("/teams", "name=My+Team")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var team Team

	_ = json.NewDecoder(w.Body).Decode(&team)

	if team.Name != "My Team" {
		t.Fatalf("expected team name 'My Team', got %s", team.Name)
	}

	if len(events.dispatched) != 1 || events.dispatched[0].(Event).Name != EventTeamCreated {
		t.Fatal("expected TeamCreated event")
	}
}

func TestCreateTeamUnauthenticated(t *testing.T) {
	js, _ := buildTestAuthKit(nil, newTestTeamRepo(), newTestInvitationRepo())
	js.guard = &testGuard{user: nil}

	handler := NewCreateTeamHandler(js)
	w := httptest.NewRecorder()
	r := postForm("/teams", "name=Test")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// --- update team tests ---

func TestUpdateTeamSuccess(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1", Name: "Old Name"}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	js, events := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("PUT /teams/{team}", NewUpdateTeamHandler(js))

	w := httptest.NewRecorder()
	r := putForm("/teams/t1", "name=New+Name")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(events.dispatched) != 1 || events.dispatched[0].(Event).Name != EventTeamUpdated {
		t.Fatal("expected TeamUpdated event")
	}
}

func TestUpdateTeamForbiddenForNonOwner(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "other", Name: "Team"}
	user := &testTeamUser{id: "1", permissions: map[string]bool{}}
	teams := newTestTeamRepo(team)
	js, _ := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("PUT /teams/{team}", NewUpdateTeamHandler(js))

	w := httptest.NewRecorder()
	r := putForm("/teams/t1", "name=Hacked")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// --- delete team tests ---

func TestDeleteTeamSuccess(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1", Name: "Team", PersonalTeam: false}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	js, events := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("DELETE /teams/{team}", NewDeleteTeamHandler(js))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/teams/t1", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(events.dispatched) != 1 || events.dispatched[0].(Event).Name != EventTeamDeleted {
		t.Fatal("expected TeamDeleted event")
	}
}

func TestDeletePersonalTeamFails(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1", PersonalTeam: true}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	js, _ := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("DELETE /teams/{team}", NewDeleteTeamHandler(js))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/teams/t1", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

// --- add member tests ---

func TestAddTeamMemberSuccess(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1"}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	js, events := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("POST /teams/{team}/members", NewAddTeamMemberHandler(js))

	w := httptest.NewRecorder()
	r := postForm("/teams/t1/members", "email=new@example.com&role=member")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(events.dispatched) != 1 || events.dispatched[0].(Event).Name != EventTeamMemberAdded {
		t.Fatal("expected TeamMemberAdded event")
	}
}

// --- update member role tests ---

func TestUpdateMemberRoleSuccess(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1"}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	js, _ := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("PUT /teams/{team}/members/{user}", NewUpdateTeamMemberRoleHandler(js))

	w := httptest.NewRecorder()
	r := putForm("/teams/t1/members/u2", "role=admin")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateMemberRoleInvalidRole(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1"}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	js, _ := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("PUT /teams/{team}/members/{user}", NewUpdateTeamMemberRoleHandler(js))

	w := httptest.NewRecorder()
	r := putForm("/teams/t1/members/u2", "role=nonexistent")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

// --- remove member tests ---

func TestRemoveTeamMemberSuccess(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1"}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	js, events := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("DELETE /teams/{team}/members/{user}", NewRemoveTeamMemberHandler(js))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/teams/t1/members/u2", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(events.dispatched) != 1 || events.dispatched[0].(Event).Name != EventTeamMemberRemoved {
		t.Fatal("expected TeamMemberRemoved event")
	}
}

// --- invitation tests ---

func TestInviteTeamMemberSuccess(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1"}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	js, events := buildTestAuthKit(user, teams, newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("POST /teams/{team}/invitations", NewInviteTeamMemberHandler(js))

	w := httptest.NewRecorder()
	r := postForm("/teams/t1/invitations", "email=invited@example.com&role=member")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	if len(events.dispatched) != 1 || events.dispatched[0].(Event).Name != EventTeamMemberInvited {
		t.Fatal("expected TeamMemberInvited event")
	}
}

func TestCancelInvitationSuccess(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "1"}
	invitation := &TeamInvitation{ID: "inv-1", TeamID: "t1", Email: "test@example.com"}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	invitations := newTestInvitationRepo(invitation)
	js, _ := buildTestAuthKit(user, teams, invitations)

	mux := http.NewServeMux()
	mux.Handle("DELETE /team-invitations/{invitation}", NewCancelInvitationHandler(js))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/team-invitations/inv-1", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !invitations.deleted {
		t.Fatal("expected invitation to be deleted")
	}
}

func TestAcceptInvitationSuccess(t *testing.T) {
	team := &Team{ID: "t1", OwnerID: "owner"}
	invitation := &TeamInvitation{ID: "inv-1", TeamID: "t1", Email: "user@example.com", Role: "member"}
	user := &testTeamUser{id: "1"}
	teams := newTestTeamRepo(team)
	invitations := newTestInvitationRepo(invitation)
	js, events := buildTestAuthKit(user, teams, invitations)

	mux := http.NewServeMux()
	mux.Handle("GET /team-invitations/{invitation}/accept", NewAcceptInvitationHandler(js))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/team-invitations/inv-1/accept", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !invitations.deleted {
		t.Fatal("expected invitation to be deleted after acceptance")
	}

	if len(events.dispatched) != 1 || events.dispatched[0].(Event).Name != EventTeamMemberAdded {
		t.Fatal("expected TeamMemberAdded event")
	}
}

// --- switch team tests ---

func TestSwitchTeamSuccess(t *testing.T) {
	team := &Team{ID: "t2", OwnerID: "1"}
	user := &testTeamUser{id: "1", teams: []Team{*team}}
	teams := newTestTeamRepo(team)
	js, events := buildTestAuthKit(user, teams, newTestInvitationRepo())

	handler := NewSwitchTeamHandler(js)
	w := httptest.NewRecorder()
	r := putForm("/current-team", "team_id=t2")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if user.currentTeam == nil || user.currentTeam.ID != "t2" {
		t.Fatal("expected current team to be switched to t2")
	}

	if len(events.dispatched) != 1 || events.dispatched[0].(Event).Name != EventTeamSwitched {
		t.Fatal("expected TeamSwitched event")
	}
}

func TestSwitchTeamForbidden(t *testing.T) {
	team := &Team{ID: "t2", OwnerID: "other"}
	user := &testTeamUser{id: "1", teams: []Team{}}
	teams := newTestTeamRepo(team)
	js, _ := buildTestAuthKit(user, teams, newTestInvitationRepo())

	handler := NewSwitchTeamHandler(js)
	w := httptest.NewRecorder()
	r := putForm("/current-team", "team_id=t2")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}
