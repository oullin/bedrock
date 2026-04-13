package inception

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// --- test session repo ---

type testSessionRepo struct {
	sessions      []BrowserSession
	deletedOthers bool
}

func (r *testSessionRepo) FindByUser(_ context.Context, _ string) ([]BrowserSession, error) {
	return r.sessions, nil
}

func (r *testSessionRepo) DeleteOthers(_ context.Context, _ string, _ string) error {
	r.deletedOthers = true

	return nil
}

// --- list sessions tests ---

func TestListSessionsSuccess(t *testing.T) {
	user := &testTeamUser{id: "1"}
	sessions := &testSessionRepo{
		sessions: []BrowserSession{
			{ID: "s1", IPAddress: "127.0.0.1", UserAgent: "Chrome", LastActive: time.Now(), IsCurrent: true},
			{ID: "s2", IPAddress: "10.0.0.1", UserAgent: "Firefox", LastActive: time.Now()},
		},
	}

	app, _ := buildTestInceptionJS(user, newTestTeamRepo(), newTestInvitationRepo())
	app.config.Features.BrowserSessions = true

	handler := NewListSessionsHandler(app, sessions)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/user/sessions", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result []BrowserSession

	_ = json.NewDecoder(w.Body).Decode(&result)

	if len(result) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(result))
	}
}

func TestListSessionsDisabled(t *testing.T) {
	app, _ := buildTestInceptionJS(&testTeamUser{id: "1"}, newTestTeamRepo(), newTestInvitationRepo())
	app.config.Features.BrowserSessions = false

	handler := NewListSessionsHandler(app, &testSessionRepo{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/user/sessions", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- delete other sessions tests ---

func TestDeleteOtherSessionsSuccess(t *testing.T) {
	user := &testTeamUser{id: "1"}
	sessions := &testSessionRepo{}

	app, _ := buildTestInceptionJS(user, newTestTeamRepo(), newTestInvitationRepo())
	app.config.Features.BrowserSessions = true

	handler := NewDeleteOtherSessionsHandler(app, sessions)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/user/other-sessions", nil)
	r = r.WithContext(WithCurrentSessionID(r.Context(), "current-session-id"))

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !sessions.deletedOthers {
		t.Fatal("expected other sessions to be deleted")
	}
}

func TestDeleteOtherSessionsMissingSessionID(t *testing.T) {
	user := &testTeamUser{id: "1"}
	sessions := &testSessionRepo{}

	app, _ := buildTestInceptionJS(user, newTestTeamRepo(), newTestInvitationRepo())
	app.config.Features.BrowserSessions = true

	handler := NewDeleteOtherSessionsHandler(app, sessions)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/user/other-sessions", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteOtherSessionsDisabled(t *testing.T) {
	app, _ := buildTestInceptionJS(&testTeamUser{id: "1"}, newTestTeamRepo(), newTestInvitationRepo())
	app.config.Features.BrowserSessions = false

	handler := NewDeleteOtherSessionsHandler(app, &testSessionRepo{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/user/other-sessions", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
