package authflows

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// --- test profile user ---

type testProfileUser struct {
	testUser
	email      string
	verified   bool
	unverified bool
}

// --- test profile action ---

type testUpdatesProfile struct {
	called   bool
	newEmail string
	err      error
}

// --- test responder extension ---

type testProfileResponder struct {
	testResponder
	profileUpdatedCalled bool
}

func (u *testProfileUser) HasVerifiedEmail() bool          { return u.verified }
func (u *testProfileUser) MarkEmailAsVerified(_ time.Time) {}
func (u *testProfileUser) MarkEmailAsUnverified()          { u.unverified = true }
func (u *testProfileUser) GetEmailForVerification() string { return u.email }

func (a *testUpdatesProfile) Update(_ context.Context, user cauth.Authenticatable, input map[string]string) error {
	a.called = true

	if a.err != nil {
		return a.err
	}

	if newEmail, ok := input["email"]; ok && newEmail != "" {
		if pu, ok := user.(*testProfileUser); ok {
			a.newEmail = newEmail
			pu.email = newEmail
		}
	}

	return nil
}

func (r *testProfileResponder) ProfileInformationUpdatedResponse(w http.ResponseWriter, _ *http.Request) {
	r.profileUpdatedCalled = true
	w.WriteHeader(http.StatusOK)
}

// --- tests ---

func TestUpdateProfileSuccess(t *testing.T) {
	user := &testProfileUser{
		testUser: testUser{id: "1"},
		email:    "old@example.com",
		verified: true,
	}

	guard := &testGuard{authenticatedUser: user}
	action := &testUpdatesProfile{}
	events := &testEvents{}
	responder := &testProfileResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, events, &responder.testResponder)
	f.config.Features.UpdateProfileInformation = true
	f.updateProfile = action
	f.responder = responder

	handler := NewUpdateProfileHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/profile-information", "name=Updated&email=old@example.com")
	r.Method = http.MethodPut

	handler.ServeHTTP(w, r)

	if !action.called {
		t.Fatal("expected UpdatesUserProfileInformation.Update to be called")
	}

	if !responder.profileUpdatedCalled {
		t.Fatal("expected profile updated response")
	}

	if len(events.dispatched) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events.dispatched))
	}
}

func TestUpdateProfileDisabled(t *testing.T) {
	f := buildTestAuthFlows(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.UpdateProfileInformation = false

	handler := NewUpdateProfileHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/user/profile-information", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateProfileUnauthenticated(t *testing.T) {
	responder := &testProfileResponder{}

	f := buildTestAuthFlows(&testGuard{authenticatedUser: nil}, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.UpdateProfileInformation = true
	f.updateProfile = &testUpdatesProfile{}
	f.responder = responder

	handler := NewUpdateProfileHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/user/profile-information", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestUpdateProfileActionFails(t *testing.T) {
	user := &testProfileUser{testUser: testUser{id: "1"}, email: "user@example.com"}
	guard := &testGuard{authenticatedUser: user}
	action := &testUpdatesProfile{err: errors.New("name is required")}
	responder := &testProfileResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.UpdateProfileInformation = true
	f.updateProfile = action
	f.responder = responder

	handler := NewUpdateProfileHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/profile-information", "name=&email=user@example.com")
	r.Method = http.MethodPut

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestUpdateProfileReVerifiesOnEmailChange(t *testing.T) {
	user := &testProfileUser{
		testUser: testUser{id: "1"},
		email:    "old@example.com",
		verified: true,
	}

	guard := &testGuard{authenticatedUser: user}
	action := &testUpdatesProfile{}
	verifier := &testVerifier{}
	responder := &testProfileResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.UpdateProfileInformation = true
	f.config.Features.EmailVerification = true
	f.updateProfile = action
	f.verifier = verifier
	f.responder = responder

	handler := NewUpdateProfileHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/profile-information", "name=Test&email=new@example.com")
	r.Method = http.MethodPut

	handler.ServeHTTP(w, r)

	if !user.unverified {
		t.Fatal("expected email to be marked as unverified after change")
	}

	if !verifier.sent {
		t.Fatal("expected verification notification to be sent after email change")
	}
}

func TestUpdateProfileSkipsVerificationWhenEmailUnchanged(t *testing.T) {
	user := &testProfileUser{
		testUser: testUser{id: "1"},
		email:    "same@example.com",
		verified: true,
	}

	guard := &testGuard{authenticatedUser: user}
	action := &testUpdatesProfile{}
	verifier := &testVerifier{}
	responder := &testProfileResponder{}

	f := buildTestAuthFlows(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.UpdateProfileInformation = true
	f.config.Features.EmailVerification = true
	f.updateProfile = action
	f.verifier = verifier
	f.responder = responder

	handler := NewUpdateProfileHandler(f)
	w := httptest.NewRecorder()
	r := postWithBody("/user/profile-information", "name=Test&email=same@example.com")
	r.Method = http.MethodPut

	handler.ServeHTTP(w, r)

	if user.unverified {
		t.Fatal("should not mark as unverified when email unchanged")
	}

	if verifier.sent {
		t.Fatal("should not send verification when email unchanged")
	}
}
