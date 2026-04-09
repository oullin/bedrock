package authflows

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- test action doubles ---

type testUpdatesPasswords struct {
	called      bool
	returnError error
}

func (a *testUpdatesPasswords) Update(_ context.Context, _ Authenticatable, _ map[string]string) error {
	a.called = true
	return a.returnError
}

type testConfirmsPasswords struct {
	called      bool
	returnError error
}

func (a *testConfirmsPasswords) Confirm(_ context.Context, _ Authenticatable, _ string) error {
	a.called = true
	return a.returnError
}

// --- test responder extension ---

type testPasswordUpdateResponder struct {
	testResponder
	updateCalled  bool
	confirmCalled bool
}

func (r *testPasswordUpdateResponder) PasswordUpdateResponse(w http.ResponseWriter, _ *http.Request) {
	r.updateCalled = true
	w.WriteHeader(http.StatusOK)
}

func (r *testPasswordUpdateResponder) PasswordConfirmResponse(w http.ResponseWriter, _ *http.Request) {
	r.confirmCalled = true
	w.WriteHeader(http.StatusOK)
}

// --- update password tests ---

func TestUpdatePasswordSuccess(t *testing.T) {
	user := &testUser{id: "1"}
	guard := &testGuard{authenticatedUser: user}
	action := &testUpdatesPasswords{}
	events := &testEvents{}
	responder := &testPasswordUpdateResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, events, &responder.testResponder)
	f.config.Features.UpdatePasswords = true
	f.updatePass = action
	f.responder = responder

	handler := NewUpdatePasswordHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/password", "current_password=old&password=new&password_confirmation=new")
	r.Method = http.MethodPut

	handler.ServeHTTP(w, r)

	if !action.called {
		t.Fatal("expected UpdatesUserPasswords.Update to be called")
	}

	if !responder.updateCalled {
		t.Fatal("expected password update response")
	}

	if len(events.dispatched) != 1 || events.dispatched[0].Name != EventPasswordUpdated {
		t.Fatal("expected PasswordUpdated event")
	}
}

func TestUpdatePasswordDisabled(t *testing.T) {
	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.UpdatePasswords = false

	handler := NewUpdatePasswordHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/user/password", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdatePasswordUnauthenticated(t *testing.T) {
	guard := &testGuard{authenticatedUser: nil}
	responder := &testPasswordUpdateResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.UpdatePasswords = true
	f.updatePass = &testUpdatesPasswords{}
	f.responder = responder

	handler := NewUpdatePasswordHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/user/password", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestUpdatePasswordActionFails(t *testing.T) {
	user := &testUser{id: "1"}
	guard := &testGuard{authenticatedUser: user}
	action := &testUpdatesPasswords{returnError: errors.New("current password is incorrect")}
	responder := &testPasswordUpdateResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.UpdatePasswords = true
	f.updatePass = action
	f.responder = responder

	handler := NewUpdatePasswordHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/password", "current_password=wrong&password=new&password_confirmation=new")
	r.Method = http.MethodPut

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

// --- confirm password tests ---

func TestConfirmPasswordSuccess(t *testing.T) {
	user := &testUser{id: "1"}
	guard := &testGuard{authenticatedUser: user}
	action := &testConfirmsPasswords{}
	responder := &testPasswordUpdateResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.ConfirmPassword = true
	f.confirmPass = action
	f.responder = responder

	handler := NewConfirmPasswordHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/confirm-password", "password=mysecret")

	handler.ServeHTTP(w, r)

	if !action.called {
		t.Fatal("expected ConfirmsPasswords.Confirm to be called")
	}

	if !responder.confirmCalled {
		t.Fatal("expected password confirm response")
	}
}

func TestConfirmPasswordDisabled(t *testing.T) {
	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.ConfirmPassword = false

	handler := NewConfirmPasswordHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/confirm-password", "password=test")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestConfirmPasswordWrongPassword(t *testing.T) {
	user := &testUser{id: "1"}
	guard := &testGuard{authenticatedUser: user}
	action := &testConfirmsPasswords{returnError: errors.New("password does not match")}
	responder := &testPasswordUpdateResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.ConfirmPassword = true
	f.confirmPass = action
	f.responder = responder

	handler := NewConfirmPasswordHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/confirm-password", "password=wrong")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}
