package authflows

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutHandler(t *testing.T) {
	user := &testUser{id: "1", email: "user@example.com"}
	guard := &testGuard{authenticatedUser: user}
	events := &testEvents{}
	responder := &testResponder{}
	f := buildTestAuthFlows(guard, &testProvider{}, events, responder)

	handler := NewLogoutHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/logout", nil)

	handler.ServeHTTP(w, r)

	if !guard.loggedOut {
		t.Fatal("expected guard.Logout to be called")
	}

	if !responder.logoutCalled {
		t.Fatal("expected logout response")
	}

	if len(events.dispatched) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events.dispatched))
	}

	if events.dispatched[0].Name != EventLoggedOut {
		t.Fatalf("expected LoggedOut event, got %s", events.dispatched[0].Name)
	}
}

func TestLogoutHandlerWithoutAuthenticatedUser(t *testing.T) {
	guard := &testGuard{authenticatedUser: nil}
	events := &testEvents{}
	responder := &testResponder{}
	f := buildTestAuthFlows(guard, &testProvider{}, events, responder)

	handler := NewLogoutHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/logout", nil)

	handler.ServeHTTP(w, r)

	if !guard.loggedOut {
		t.Fatal("expected guard.Logout to be called")
	}

	if !responder.logoutCalled {
		t.Fatal("expected logout response")
	}

	if len(events.dispatched) != 0 {
		t.Fatalf("expected 0 events when no user, got %d", len(events.dispatched))
	}
}
