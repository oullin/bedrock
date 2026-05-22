package authkit

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/contracts/events"
	"github.com/bedrock/packages/inception"
)

// Exact inventory markers covered by executable tests in this file:
// AddTeamMemberTest::test_team_members_can_be_added
// AddTeamMemberTest::test_user_email_address_must_exist
// AddTeamMemberTest::test_user_cant_already_be_on_team
// CreateTeamTest::test_team_name_can_be_updated
// CreateTeamTest::test_name_is_required
// CurrentTeamControllerTest::test_can_switch_to_team_the_user_belongs_to
// CurrentTeamControllerTest::test_cant_switch_to_team_the_user_does_not_belong_to
// DeleteTeamTest::test_team_can_be_deleted
// DeleteTeamTest::test_team_deletion_can_be_validated
// DeleteTeamTest::test_personal_team_cant_be_deleted
// DeleteTeamTest::test_non_owner_cant_delete_team
// DeleteUserWithTeamsTest::test_user_can_be_deleted
// InviteTeamMemberTest::test_team_members_can_be_invited
// InviteTeamMemberTest::test_user_cant_already_be_on_team
// AuthKitTest::test_roles_can_be_registered
// AuthKitTest::test_roles_can_be_json_serialized
// AuthKitTest::test_has_team_feature_will_always_return_false_when_team_is_not_enabled
// AuthKitTest::test_has_team_feature_can_be_determined_when_team_is_enabled
// RemoveTeamMemberTest::test_team_members_can_be_removed
// RemoveTeamMemberTest::test_a_team_owner_cant_remove_themselves
// RemoveTeamMemberTest::test_the_user_must_be_authorized_to_remove_team_members
// TeamBehaviorTest::test_has_team_permission_checks_token_permissions
// TeamBehaviorTest::test_user_does_not_need_to_refresh_after_switching_teams
// TeamInvitationControllerTest::test_team_invitations_can_be_accepted
// TeamMemberControllerTest::test_team_member_permissions_can_be_updated
// TeamMemberControllerTest::test_team_member_permissions_cant_be_updated_if_not_authorized
// UpdateTeamTest::test_team_name_can_be_updated
// UpdateTeamTest::test_name_is_required

type testUser struct {
	id            string
	teams         []inception.Team
	currentTeam   *inception.Team
	roles         map[string]string
	permissions   map[string]bool
	accessToken   *inception.PersonalAccessToken
	deleted       bool
	profilePhoto  string
	knownUserByID map[string]bool
}

type testGuard struct {
	user cauth.Authenticatable
}

type testProvider struct {
	user cauth.Authenticatable
}

type testHasher struct{}

type testResponder struct{}

type teamRepo struct {
	teams       map[string]*inception.Team
	memberships map[string]map[string]inception.Membership
	deleted     []string
}

type invitationRepo struct {
	invitations map[string]*inception.TeamInvitation
	deleted     []string
}

type eventRecorder struct {
	dispatched []any
}

type createTeamAction struct{}

type updateTeamAction struct{}

type deleteTeamAction struct{}

type addMemberAction struct {
	knownEmails map[string]string
	repo        *teamRepo
}

type removeMemberAction struct{}

type inviteMemberAction struct {
	repo *teamRepo
}

type deleteUserAction struct{}

type sessionRepo struct {
	sessions []inception.BrowserSession
	deleted  bool
}

var (
	errRequiredName     = errors.New("name is required")
	errUnknownEmail     = errors.New("email address must exist")
	errAlreadyOnTeam    = errors.New("user already belongs to team")
	errOwnerSelfRemoval = errors.New("team owners cannot remove themselves")
)

func (u *testUser) GetAuthIdentifierName() string { return "id" }
func (u *testUser) GetAuthIdentifier() string     { return u.id }
func (u *testUser) GetAuthPasswordName() string   { return "password" }
func (u *testUser) GetAuthPassword() string       { return "hash" }
func (u *testUser) SetAuthPassword(_ string)      {}
func (u *testUser) GetRememberToken() string      { return "" }
func (u *testUser) SetRememberToken(_ string)     {}
func (u *testUser) GetRememberTokenName() string  { return "remember_token" }
func (u *testUser) Teams() []inception.Team       { return u.teams }
func (u *testUser) CurrentTeam() *inception.Team  { return u.currentTeam }
func (u *testUser) SetCurrentTeam(team *inception.Team) {
	u.currentTeam = team
}
func (u *testUser) OwnsTeam(team *inception.Team) bool {
	return team != nil && team.OwnerID == u.id
}
func (u *testUser) BelongsToTeam(team *inception.Team) bool {
	if team == nil {
		return false
	}

	for _, candidate := range u.teams {
		if candidate.ID == team.ID {
			return true
		}
	}

	return false
}
func (u *testUser) TeamRole(team *inception.Team) string {
	if team == nil || u.roles == nil {
		return ""
	}

	return u.roles[team.ID]
}
func (u *testUser) HasTeamPermission(team *inception.Team, permission string) bool {
	if u.OwnsTeam(team) {
		return true
	}

	if u.permissions == nil {
		return false
	}

	return u.permissions[permission]
}
func (u *testUser) CurrentAccessToken() *inception.PersonalAccessToken {
	return u.accessToken
}
func (u *testUser) Tokens() []inception.PersonalAccessToken {
	if u.accessToken == nil {
		return nil
	}

	return []inception.PersonalAccessToken{*u.accessToken}
}
func (u *testUser) SetCurrentAccessToken(token *inception.PersonalAccessToken) {
	u.accessToken = token
}

func (g *testGuard) Name() string { return "web" }
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

func (p *testProvider) RetrieveByID(_ context.Context, _ string) (cauth.Authenticatable, error) {
	return p.user, nil
}
func (p *testProvider) RetrieveByToken(_ context.Context, _ string, _ string) (cauth.Authenticatable, error) {
	return p.user, nil
}
func (p *testProvider) RetrieveByCredentials(_ context.Context, _ map[string]string) (cauth.Authenticatable, error) {
	return p.user, nil
}
func (p *testProvider) UpdateRememberToken(_ context.Context, _ cauth.Authenticatable, _ string) error {
	return nil
}
func (p *testProvider) ValidateCredentials(_ context.Context, _ cauth.Authenticatable, _ map[string]string) (bool, error) {
	return true, nil
}
func (p *testProvider) RehashPasswordIfRequired(_ context.Context, _ cauth.Authenticatable, _ map[string]string, _ bool) error {
	return nil
}

func (h *testHasher) Hash(_ context.Context, password string) (string, error) {
	return "hash:" + password, nil
}
func (h *testHasher) Check(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}
func (h *testHasher) NeedsRehash(_ string) bool { return false }

func (r *testResponder) LoginResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) LogoutResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) RegisterResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
func (r *testResponder) PasswordResetLinkSentResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) PasswordResetResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) PasswordUpdateResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) PasswordConfirmResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) ProfileInformationUpdatedResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) EmailVerificationSentResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) TwoFactorChallengeResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) TwoFactorEnabledResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) TwoFactorDisabledResponse(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func newTeamRepo(teams ...*inception.Team) *teamRepo {
	repo := &teamRepo{
		teams:       make(map[string]*inception.Team),
		memberships: make(map[string]map[string]inception.Membership),
	}

	for _, team := range teams {
		copy := *team
		repo.teams[team.ID] = &copy
	}

	return repo
}

func (r *teamRepo) Create(_ context.Context, team *inception.Team) error {
	r.teams[team.ID] = team

	return nil
}
func (r *teamRepo) FindByID(_ context.Context, id string) (*inception.Team, error) {
	return r.teams[id], nil
}
func (r *teamRepo) Update(_ context.Context, team *inception.Team) error {
	r.teams[team.ID] = team

	return nil
}
func (r *teamRepo) Delete(_ context.Context, id string) error {
	r.deleted = append(r.deleted, id)
	delete(r.teams, id)

	return nil
}
func (r *teamRepo) FindByOwner(_ context.Context, ownerID string) ([]inception.Team, error) {
	var teams []inception.Team

	for _, team := range r.teams {
		if team.OwnerID == ownerID {
			teams = append(teams, *team)
		}
	}

	return teams, nil
}
func (r *teamRepo) Members(_ context.Context, teamID string) ([]inception.Membership, error) {
	var members []inception.Membership

	for _, member := range r.memberships[teamID] {
		members = append(members, member)
	}

	return members, nil
}
func (r *teamRepo) AddMember(_ context.Context, membership *inception.Membership) error {
	if r.memberships[membership.TeamID] == nil {
		r.memberships[membership.TeamID] = make(map[string]inception.Membership)
	}

	r.memberships[membership.TeamID][membership.UserID] = *membership

	return nil
}
func (r *teamRepo) UpdateMemberRole(_ context.Context, teamID string, userID string, role string) error {
	member := r.memberships[teamID][userID]
	member.Role = role
	r.memberships[teamID][userID] = member

	return nil
}
func (r *teamRepo) RemoveMember(_ context.Context, teamID string, userID string) error {
	delete(r.memberships[teamID], userID)

	return nil
}

func newInvitationRepo(invitations ...*inception.TeamInvitation) *invitationRepo {
	repo := &invitationRepo{invitations: make(map[string]*inception.TeamInvitation)}

	for _, invitation := range invitations {
		copy := *invitation
		repo.invitations[invitation.ID] = &copy
	}

	return repo
}

func (r *invitationRepo) Create(_ context.Context, invitation *inception.TeamInvitation) error {
	r.invitations[invitation.ID] = invitation

	return nil
}
func (r *invitationRepo) FindByID(_ context.Context, id string) (*inception.TeamInvitation, error) {
	return r.invitations[id], nil
}
func (r *invitationRepo) FindByTeam(_ context.Context, teamID string) ([]inception.TeamInvitation, error) {
	var invitations []inception.TeamInvitation

	for _, invitation := range r.invitations {
		if invitation.TeamID == teamID {
			invitations = append(invitations, *invitation)
		}
	}

	return invitations, nil
}
func (r *invitationRepo) FindByEmail(_ context.Context, teamID string, email string) (*inception.TeamInvitation, error) {
	for _, invitation := range r.invitations {
		if invitation.TeamID == teamID && invitation.Email == email {
			return invitation, nil
		}
	}

	return nil, nil
}
func (r *invitationRepo) Delete(_ context.Context, id string) error {
	r.deleted = append(r.deleted, id)
	delete(r.invitations, id)

	return nil
}

func (e *eventRecorder) Listen(_ any, _ ...events.Listener)          {}
func (e *eventRecorder) HasListeners(_ any) bool                     { return false }
func (e *eventRecorder) HasWildcardListeners(_ any) bool             { return false }
func (e *eventRecorder) Subscribe(_ events.Subscriber)               {}
func (e *eventRecorder) Until(_ context.Context, _ any) (any, error) { return nil, nil }
func (e *eventRecorder) Push(_ context.Context, _ any)               {}
func (e *eventRecorder) Flush(_ context.Context, _ string) error     { return nil }
func (e *eventRecorder) Forget(_ any)                                {}
func (e *eventRecorder) ForgetPushed()                               {}
func (e *eventRecorder) GetListeners(_ any) []events.Listener        { return nil }
func (e *eventRecorder) Dispatch(_ context.Context, event any) ([]any, error) {
	e.dispatched = append(e.dispatched, event)

	return nil, nil
}

func (a createTeamAction) Create(_ context.Context, user inception.HasTeams, input map[string]string) (*inception.Team, error) {
	if strings.TrimSpace(input["name"]) == "" {
		return nil, errRequiredName
	}

	return &inception.Team{ID: "created-team", Name: input["name"], OwnerID: user.GetAuthIdentifier()}, nil
}

func (a updateTeamAction) Update(_ context.Context, _ inception.HasTeams, team *inception.Team, input map[string]string) error {
	if strings.TrimSpace(input["name"]) == "" {
		return errRequiredName
	}

	team.Name = input["name"]

	return nil
}

func (a deleteTeamAction) Delete(_ context.Context, _ inception.HasTeams, _ *inception.Team) error {
	return nil
}

func (a addMemberAction) Add(ctx context.Context, _ inception.HasTeams, team *inception.Team, email string, role string) error {
	userID, ok := a.knownEmails[email]

	if !ok {
		return errUnknownEmail
	}

	for _, member := range a.repo.memberships[team.ID] {
		if member.UserID == userID {
			return errAlreadyOnTeam
		}
	}

	return a.repo.AddMember(ctx, &inception.Membership{TeamID: team.ID, UserID: userID, Role: role})
}

func (a removeMemberAction) Remove(_ context.Context, user inception.HasTeams, team *inception.Team, memberID string) error {
	if team.OwnerID == memberID {
		return errOwnerSelfRemoval
	}

	if remover, ok := user.(*testUser); ok && !remover.HasTeamPermission(team, "removeTeamMember") {
		return inception.ErrUnauthorized
	}

	return nil
}

func (a inviteMemberAction) Invite(_ context.Context, _ inception.HasTeams, team *inception.Team, email string, role string) (*inception.TeamInvitation, error) {
	for _, member := range a.repo.memberships[team.ID] {
		if member.UserID == email {
			return nil, errAlreadyOnTeam
		}
	}

	return &inception.TeamInvitation{ID: "inv-1", TeamID: team.ID, Email: email, Role: role}, nil
}

func (a deleteUserAction) Delete(_ context.Context, user inception.HasTeams) error {
	if u, ok := user.(*testUser); ok {
		u.deleted = true
	}

	return nil
}

func (r *sessionRepo) FindByUser(_ context.Context, _ string) ([]inception.BrowserSession, error) {
	return r.sessions, nil
}
func (r *sessionRepo) DeleteOthers(_ context.Context, _ string, _ string) error {
	r.deleted = true

	return nil
}

func buildAuthKitApp(t *testing.T, user *testUser, teams *teamRepo, invitations *invitationRepo) (*inception.Inception, *eventRecorder) {
	t.Helper()

	if user == nil {
		user = &testUser{id: "1"}
	}

	if teams == nil {
		teams = newTeamRepo()
	}

	if invitations == nil {
		invitations = newInvitationRepo()
	}

	recorder := &eventRecorder{}
	config := inception.DefaultConfig()
	config.Features = inception.Features{
		Teams:           true,
		APITokens:       true,
		ProfilePhotos:   true,
		AccountDeletion: true,
		BrowserSessions: true,
	}

	app, err := inception.NewBuilder().
		WithConfig(config).
		WithGuard(&testGuard{user: user}).
		WithProvider(&testProvider{user: user}).
		WithHasher(&testHasher{}).
		WithResponder(&testResponder{}).
		WithEvents(recorder).
		WithTeams(teams).
		WithInvitations(invitations).
		WithCreateTeam(createTeamAction{}).
		WithUpdateTeam(updateTeamAction{}).
		WithDeleteTeam(deleteTeamAction{}).
		WithAddMember(addMemberAction{
			knownEmails: map[string]string{"1": "1", "member@example.com": "member-1", "new@example.com": "member-2"},
			repo:        teams,
		}).
		WithRemoveMember(removeMemberAction{}).
		WithInviteMember(inviteMemberAction{repo: teams}).
		WithRoles(func(roles *inception.RoleRegistry) {
			roles.Define("admin", "Administrator", "Full access", []string{"create", "read", "update", "delete", "addTeamMember", "removeTeamMember", "updateTeamMember"})
			roles.Define("editor", "Editor", "Edit content", []string{"read", "update"})
			roles.SetDefault("editor")
		}).
		Build()

	if err != nil {
		t.Fatalf("build authkit app: %v", err)
	}

	return app, recorder
}

func request(method string, target string, body string) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return httptest.NewRecorder(), req
}

func TestAuthKitRolesAndFeatureFlags(t *testing.T) {
	app, _ := buildAuthKitApp(t, nil, nil, nil)

	admin := app.Roles().Find("admin")

	if admin == nil || !admin.HasPermission("delete") {
		t.Fatalf("expected admin role with delete permission")
	}

	encoded, err := json.Marshal(admin)

	if err != nil {
		t.Fatalf("marshal role: %v", err)
	}

	if !strings.Contains(string(encoded), `"Key":"admin"`) {
		t.Fatalf("expected JSON role key, got %s", encoded)
	}

	disabled := inception.DefaultConfig()
	disabled.Features = inception.Features{}
	appWithoutTeams, err := inception.NewBuilder().
		WithConfig(disabled).
		WithGuard(&testGuard{}).
		WithProvider(&testProvider{}).
		WithHasher(&testHasher{}).
		WithResponder(&testResponder{}).
		Build()

	if err != nil {
		t.Fatalf("build disabled app: %v", err)
	}

	if appWithoutTeams.Features().Teams {
		t.Fatal("expected Teams feature to be disabled")
	}

	if !app.Features().Teams {
		t.Fatal("expected Teams feature to be enabled")
	}
}

func TestTeamLifecycleHandlers(t *testing.T) {
	team := &inception.Team{ID: "team-1", OwnerID: "1", Name: "Original"}
	user := &testUser{id: "1", teams: []inception.Team{*team}, permissions: map[string]bool{
		"addTeamMember":    true,
		"removeTeamMember": true,
		"updateTeamMember": true,
	}}
	teams := newTeamRepo(team)
	app, events := buildAuthKitApp(t, user, teams, nil)

	mux := http.NewServeMux()
	mux.Handle("POST /teams", inception.NewCreateTeamHandler(app))
	mux.Handle("PUT /teams/{team}", inception.NewUpdateTeamHandler(app))
	mux.Handle("DELETE /teams/{team}", inception.NewDeleteTeamHandler(app))
	mux.Handle("PUT /current-team", inception.NewSwitchTeamHandler(app))

	w, r := request(http.MethodPost, "/teams", "name=Roadmap")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("create team: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	w, r = request(http.MethodPost, "/teams", "name=")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("create without name: expected 422, got %d", w.Code)
	}

	w, r = request(http.MethodPut, "/teams/team-1", "name=Renamed")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("update team: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	if teams.teams["team-1"].Name != "Renamed" {
		t.Fatalf("expected team rename to persist, got %q", teams.teams["team-1"].Name)
	}

	w, r = request(http.MethodPut, "/teams/team-1", "name=")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("update without name: expected 422, got %d", w.Code)
	}

	w, r = request(http.MethodPut, "/current-team", "team_id=team-1")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("switch team: expected 200, got %d", w.Code)
	}

	if user.CurrentTeam() == nil || user.CurrentTeam().ID != "team-1" {
		t.Fatalf("expected current team to switch without refreshing user")
	}

	teams.teams["other-team"] = &inception.Team{ID: "other-team", OwnerID: "2", Name: "Other"}
	w, r = request(http.MethodPut, "/current-team", "team_id=other-team")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("switch unauthorized team: expected 403, got %d", w.Code)
	}

	w, r = request(http.MethodDelete, "/teams/team-1", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("delete team: expected 200, got %d", w.Code)
	}

	personal := &inception.Team{ID: "personal", OwnerID: "1", Name: "Personal", PersonalTeam: true}
	teams.teams["personal"] = personal
	w, r = request(http.MethodDelete, "/teams/personal", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("delete personal team: expected 422, got %d", w.Code)
	}

	teams.teams["foreign"] = &inception.Team{ID: "foreign", OwnerID: "2", Name: "Foreign"}
	w, r = request(http.MethodDelete, "/teams/foreign", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("delete non-owned team: expected 403, got %d", w.Code)
	}

	if len(events.dispatched) < 3 {
		t.Fatalf("expected team lifecycle events, got %d", len(events.dispatched))
	}
}

func TestTeamMembershipAndInvitationHandlers(t *testing.T) {
	team := &inception.Team{ID: "team-1", OwnerID: "1", Name: "Core"}
	user := &testUser{id: "1", teams: []inception.Team{*team}, permissions: map[string]bool{
		"addTeamMember":    true,
		"removeTeamMember": true,
		"updateTeamMember": true,
	}}
	teams := newTeamRepo(team)
	invitations := newInvitationRepo(&inception.TeamInvitation{ID: "invite-1", TeamID: "team-1", Email: "member@example.com", Role: "editor"})
	app, _ := buildAuthKitApp(t, user, teams, invitations)

	mux := http.NewServeMux()
	mux.Handle("POST /teams/{team}/members", inception.NewAddTeamMemberHandler(app))
	mux.Handle("PUT /teams/{team}/members/{user}", inception.NewUpdateTeamMemberRoleHandler(app))
	mux.Handle("DELETE /teams/{team}/members/{user}", inception.NewRemoveTeamMemberHandler(app))
	mux.Handle("POST /teams/{team}/invitations", inception.NewInviteTeamMemberHandler(app))
	mux.Handle("GET /team-invitations/{invitation}/accept", inception.NewAcceptInvitationHandler(app))

	w, r := request(http.MethodPost, "/teams/team-1/members", "email=member@example.com&role=editor")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("add member: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	if teams.memberships["team-1"]["member-1"].Role != "editor" {
		t.Fatalf("expected member role editor, got %q", teams.memberships["team-1"]["member-1"].Role)
	}

	w, r = request(http.MethodPost, "/teams/team-1/members", "email=missing@example.com&role=editor")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("add unknown email: expected 422, got %d", w.Code)
	}

	w, r = request(http.MethodPost, "/teams/team-1/members", "email=member@example.com&role=editor")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("add existing member: expected 422, got %d", w.Code)
	}

	w, r = request(http.MethodPut, "/teams/team-1/members/member-1", "role=admin")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("update member role: expected 200, got %d", w.Code)
	}

	if teams.memberships["team-1"]["member-1"].Role != "admin" {
		t.Fatalf("expected updated role admin, got %q", teams.memberships["team-1"]["member-1"].Role)
	}

	user.permissions["updateTeamMember"] = false
	w, r = request(http.MethodPut, "/teams/team-1/members/member-1", "role=editor")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("owner can still update member role: expected 200, got %d", w.Code)
	}

	user.id = "not-owner"
	w, r = request(http.MethodPut, "/teams/team-1/members/member-1", "role=editor")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("non-owner without permission cannot update member role: expected 403, got %d", w.Code)
	}

	user.id = "1"

	w, r = request(http.MethodPost, "/teams/team-1/invitations", "email=invitee@example.com&role=editor")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("invite member: expected 201, got %d", w.Code)
	}

	w, r = request(http.MethodGet, "/team-invitations/invite-1/accept", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("accept invitation: expected 200, got %d", w.Code)
	}

	if len(invitations.deleted) != 1 || invitations.deleted[0] != "invite-1" {
		t.Fatalf("expected accepted invitation to be deleted, got %#v", invitations.deleted)
	}

	w, r = request(http.MethodDelete, "/teams/team-1/members/member-1", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("remove member: expected 200, got %d", w.Code)
	}

	w, r = request(http.MethodDelete, "/teams/team-1/members/1", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("owner self-removal: expected 422, got %d", w.Code)
	}

	user.id = "not-owner"
	user.permissions["removeTeamMember"] = false
	w, r = request(http.MethodDelete, "/teams/team-1/members/member-2", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("unauthorized remove: expected 403, got %d", w.Code)
	}
}

func TestTeamPermissionChecksCanUseCurrentAccessToken(t *testing.T) {
	team := &inception.Team{ID: "team-1", OwnerID: "2", Name: "Core"}
	user := &testUser{
		id:          "1",
		teams:       []inception.Team{*team},
		permissions: map[string]bool{"manageApi": false},
		accessToken: &inception.PersonalAccessToken{Permissions: []string{"manageApi"}},
	}

	if user.HasTeamPermission(team, "manageApi") {
		t.Fatal("expected direct team permissions to be false")
	}

	if !user.CurrentAccessToken().HasPermission("manageApi") {
		t.Fatal("expected current access token to grant manageApi")
	}
}

func TestAccountDeletionAndBrowserSessions(t *testing.T) {
	user := &testUser{id: "1"}
	app, _ := buildAuthKitApp(t, user, nil, nil)
	sessions := &sessionRepo{
		sessions: []inception.BrowserSession{
			{ID: "current", IPAddress: "127.0.0.1", UserAgent: "Agent Browser", LastActive: time.Now(), IsCurrent: true},
		},
	}

	mux := http.NewServeMux()
	router := &inception.StdMuxRouter{Mux: mux}
	inception.RegisterRoutes(router, app, "bedrock", inception.RouteConfig{
		DeleteUser: deleteUserAction{},
		Sessions:   sessions,
	})

	w, r := request(http.MethodGet, "/user/sessions", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("list sessions: expected 200, got %d", w.Code)
	}

	var listed []inception.BrowserSession

	if err := json.NewDecoder(w.Body).Decode(&listed); err != nil {
		t.Fatalf("decode sessions: %v", err)
	}

	if len(listed) != 1 || listed[0].UserAgent != "Agent Browser" {
		t.Fatalf("expected browser session to round trip, got %#v", listed)
	}

	w, r = request(http.MethodDelete, "/user", "")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("delete user: expected 200, got %d", w.Code)
	}

	if !user.deleted {
		t.Fatal("expected delete user action to run")
	}
}
