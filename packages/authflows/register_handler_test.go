package authflows

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// --- test action ---

type testCreatesUsers struct {
	created     bool
	returnUser  Authenticatable
	returnError error
}

func (a *testCreatesUsers) Create(_ context.Context, _ map[string]string) (Authenticatable, error) {
	a.created = true
	return a.returnUser, a.returnError
}

// --- test verifier ---

type testVerifier struct {
	sent bool
}

func (v *testVerifier) SendVerificationNotification(_ context.Context, _ Authenticatable) error {
	v.sent = true
	return nil
}

func (v *testVerifier) Verify(_ context.Context, _ string, _ string) error {
	return nil
}

// --- test responder extension ---

type testRegisterResponder struct {
	testResponder
	registerCalled bool
}

func (r *testRegisterResponder) RegisterResponse(w http.ResponseWriter, _ *http.Request) {
	r.registerCalled = true
	w.WriteHeader(http.StatusCreated)
}

// --- helpers ---

func registerRequest(email string, name string, password string) *http.Request {
	body := strings.NewReader("email=" + email + "&name=" + name + "&password=" + password + "&password_confirmation=" + password)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

// --- tests ---

func TestRegisterHandlerSuccess(t *testing.T) {
	user := &testUser{id: "1", email: "new@example.com"}
	guard := &testGuard{}
	action := &testCreatesUsers{returnUser: user}
	events := &testEvents{}
	responder := &testRegisterResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, events, &responder.testResponder)
	f.config.Features.Registration = true
	f.createUser = action
	f.responder = responder

	handler := NewRegisterHandler(f)
	w := httptest.NewRecorder()
	r := registerRequest("new@example.com", "Test User", "password123")

	handler.ServeHTTP(w, r)

	if !action.created {
		t.Fatal("expected CreatesNewUsers.Create to be called")
	}

	if guard.loggedIn == nil {
		t.Fatal("expected user to be auto-logged in")
	}

	if guard.loggedIn.GetAuthIdentifier() != "1" {
		t.Fatalf("expected user id 1, got %s", guard.loggedIn.GetAuthIdentifier())
	}

	if !responder.registerCalled {
		t.Fatal("expected register response")
	}

	if len(events.dispatched) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events.dispatched))
	}

	if events.dispatched[0].Name != EventRegistered {
		t.Fatalf("expected Registered event, got %s", events.dispatched[0].Name)
	}
}

func TestRegisterHandlerDisabled(t *testing.T) {
	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.Registration = false

	handler := NewRegisterHandler(f)
	w := httptest.NewRecorder()
	r := registerRequest("new@example.com", "Test", "pass")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when registration disabled, got %d", w.Code)
	}
}

func TestRegisterHandlerCreateFails(t *testing.T) {
	action := &testCreatesUsers{returnError: errors.New("validation failed")}
	responder := &testRegisterResponder{}

	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.Registration = true
	f.createUser = action
	f.responder = responder

	handler := NewRegisterHandler(f)
	w := httptest.NewRecorder()
	r := registerRequest("bad@example.com", "", "")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	if responder.registerCalled {
		t.Fatal("register response should not be called on failure")
	}
}

func TestRegisterHandlerSendsVerification(t *testing.T) {
	user := &testVerifiableUser{
		testUser: testUser{id: "1", email: "new@example.com"},
		verified: false,
	}
	guard := &testGuard{}
	action := &testCreatesUsers{returnUser: user}
	verifier := &testVerifier{}
	responder := &testRegisterResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.Registration = true
	f.config.Features.EmailVerification = true
	f.createUser = action
	f.verifier = verifier
	f.responder = responder

	handler := NewRegisterHandler(f)
	w := httptest.NewRecorder()
	r := registerRequest("new@example.com", "Test", "pass")

	handler.ServeHTTP(w, r)

	if !verifier.sent {
		t.Fatal("expected verification notification to be sent")
	}
}

func TestRegisterHandlerSkipsVerificationWhenAlreadyVerified(t *testing.T) {
	user := &testVerifiableUser{
		testUser: testUser{id: "1", email: "new@example.com"},
		verified: true,
	}
	guard := &testGuard{}
	action := &testCreatesUsers{returnUser: user}
	verifier := &testVerifier{}
	responder := &testRegisterResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.Registration = true
	f.config.Features.EmailVerification = true
	f.createUser = action
	f.verifier = verifier
	f.responder = responder

	handler := NewRegisterHandler(f)
	w := httptest.NewRecorder()
	r := registerRequest("verified@example.com", "Test", "pass")

	handler.ServeHTTP(w, r)

	if verifier.sent {
		t.Fatal("should not send verification when already verified")
	}
}

// --- verifiable test user ---

type testVerifiableUser struct {
	testUser
	verified bool
}

func (u *testVerifiableUser) HasVerifiedEmail() bool          { return u.verified }
func (u *testVerifiableUser) MarkEmailAsVerified(_ time.Time) {}
func (u *testVerifiableUser) MarkEmailAsUnverified()          {}
func (u *testVerifiableUser) GetEmailForVerification() string { return u.email }
