package jetstream

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func buildFullJetstream() (*Jetstream, *testTeamUser, RouteConfig) {
	team := &Team{ID: "t1", OwnerID: "1", Name: "Test Team"}
	invitation := &TeamInvitation{ID: "inv-1", TeamID: "t1", Email: "invite@example.com", Role: "member"}
	invitation2 := &TeamInvitation{ID: "inv-2", TeamID: "t1", Email: "invite2@example.com", Role: "member"}

	user := &testTeamUser{
		id:    "1",
		teams: []Team{*team},
	}

	teams := newTestTeamRepo(team)
	invitations := newTestInvitationRepo(invitation, invitation2)
	tokens := newTestTokenRepo(&PersonalAccessToken{ID: "tok-1", UserID: "1", TokenHash: "hash", Permissions: []string{"read"}})
	sessions := &testSessionRepo{
		sessions: []BrowserSession{
			{ID: "s1", IPAddress: "127.0.0.1", UserAgent: "Test", LastActive: time.Now(), IsCurrent: true},
		},
	}

	js, _ := buildTestJetstream(user, teams, invitations)

	config := RouteConfig{
		Tokens:       tokens,
		Sessions:     sessions,
		Photos:       &testUpdatesPhotos{url: "https://cdn.example.com/photo.jpg"},
		DeletePhotos: &testDeletesPhotos{},
		DeleteUser:   &testDeletesUsers{},
	}

	return js, user, config
}

func TestRegisterRoutesAllFeatures(t *testing.T) {
	js, _, config := buildFullJetstream()

	mux := http.NewServeMux()
	router := &StdMuxRouter{Mux: mux}
	RegisterRoutes(router, js, config)

	routes := []struct {
		method string
		path   string
	}{
		// Teams
		{"POST", "/teams"},
		{"PUT", "/teams/t1"},
		{"DELETE", "/teams/t1"},
		{"POST", "/teams/t1/members"},
		{"PUT", "/teams/t1/members/u2"},
		{"DELETE", "/teams/t1/members/u2"},
		{"POST", "/teams/t1/invitations"},
		{"GET", "/team-invitations/inv-1/accept"},
		{"DELETE", "/team-invitations/inv-2"},
		{"PUT", "/current-team"},
		// API Tokens
		{"POST", "/user/api-tokens"},
		{"PUT", "/user/api-tokens/tok-1"},
		{"DELETE", "/user/api-tokens/tok-1"},
		// Profile Photos
		{"PUT", "/user/profile-photo"},
		{"DELETE", "/user/profile-photo"},
		// Account Deletion
		{"DELETE", "/user"},
		// Browser Sessions
		{"GET", "/user/sessions"},
		{"DELETE", "/user/other-sessions"},
	}

	for _, tc := range routes {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			var body *strings.Reader
			switch tc.method {
			case "POST":
				body = strings.NewReader("name=Test&email=test@example.com&role=member&permissions=read")
			case "PUT":
				body = strings.NewReader("name=Updated&role=admin&team_id=t1&permissions=read")
			default:
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(tc.method, tc.path, body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code == http.StatusMethodNotAllowed || w.Code == 404 {
				t.Fatalf("route %s %s not registered, got %d", tc.method, tc.path, w.Code)
			}
		})
	}
}

func TestRegisterRoutesMinimalFeatures(t *testing.T) {
	js, _, config := buildFullJetstream()
	js.features = Features{}

	mux := http.NewServeMux()
	router := &StdMuxRouter{Mux: mux}
	RegisterRoutes(router, js, config)

	disabled := []struct {
		method string
		path   string
	}{
		{"POST", "/teams"},
		{"POST", "/user/api-tokens"},
		{"PUT", "/user/profile-photo"},
		{"DELETE", "/user"},
		{"GET", "/user/sessions"},
	}

	for _, tc := range disabled {
		t.Run(tc.method+" "+tc.path+" disabled", func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != 404 && w.Code != http.StatusMethodNotAllowed {
				t.Fatalf("route %s %s should not be registered, got %d", tc.method, tc.path, w.Code)
			}
		})
	}
}

func TestFullTeamLifecycle(t *testing.T) {
	js, user, config := buildFullJetstream()

	mux := http.NewServeMux()
	router := &StdMuxRouter{Mux: mux}
	RegisterRoutes(router, js, config)

	// Step 1: Create team
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/teams", strings.NewReader("name=New+Team"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("step 1 create: expected 201, got %d", w.Code)
	}

	var created Team
	_ = json.NewDecoder(w.Body).Decode(&created)
	if created.Name != "New Team" {
		t.Fatalf("step 1: expected 'New Team', got %s", created.Name)
	}

	// Add the created team to the repo so FindByID works
	teamRepo := js.teams.(*testTeamRepo)
	teamRepo.teams[created.ID] = &created

	// Step 2: Switch to new team (user belongs via ownership)
	user.teams = append(user.teams, created)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("PUT", "/current-team", strings.NewReader("team_id="+created.ID))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("step 2 switch: expected 200, got %d", w.Code)
	}

	if user.currentTeam == nil || user.currentTeam.ID != created.ID {
		t.Fatal("step 2: expected current team to be switched")
	}

	// Step 3: List sessions
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/user/sessions", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("step 3 sessions: expected 200, got %d", w.Code)
	}

	var sessions []BrowserSession
	_ = json.NewDecoder(w.Body).Decode(&sessions)
	if len(sessions) != 1 {
		t.Fatalf("step 3: expected 1 session, got %d", len(sessions))
	}
}
