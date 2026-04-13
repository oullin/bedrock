package inception

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// --- test action ---

type testCreatesUsers struct {
	created     bool
	returnUser  cauth.Authenticatable
	returnError error
}

// --- test verifier ---

type testVerifier struct {
	sent bool
}

// --- test responder extension ---

type testRegisterResponder struct {
	testResponder
	registerCalled bool
}

// --- tests ---

// --- verifiable test user ---

type testVerifiableUser struct {
	testUser
	verified bool
}

func (a *testCreatesUsers) Create(_ context.Context, _ map[string]string) (cauth.Authenticatable, error) {
	a.created = true

	return a.returnUser, a.returnError
}

func (v *testVerifier) SendVerificationNotification(_ context.Context, _ cauth.Authenticatable) error {
	v.sent = true

	return nil
}

func (v *testVerifier) Verify(_ context.Context, _ string, _ string) error {
	return nil
}

func (r *testRegisterResponder) RegisterResponse(w http.ResponseWriter, _ *http.Request) {
	r.registerCalled = true
	w.WriteHeader(http.StatusCreated)
}

func registerRequest(email string, name string, password string) *http.Request {
	body := strings.NewReader("email=" + email + "&name=" + name + "&password=" + password + "&password_confirmation=" + password)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

func TestRegisterHandlerSuccess(t *testing.T) {
	user := &testUser{id: "1", email: "new@example.com"}
	guard := &testGuard{}
	action := &testCreatesUsers{returnUser: user}
	events := &testEvents{}
	responder := &testRegisterResponder{}

	f := buildTestInception(guard, &testProvider{}, events, &responder.testResponder)
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

	if _, ok := events.dispatched[0].(RegisteredPayload); !ok {
		t.Fatalf("expected RegisteredPayload event, got %T", events.dispatched[0])
	}
}

func TestRegisterHandlerDisabled(t *testing.T) {
	f := buildTestInception(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
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

	f := buildTestInception(&testGuard{}, &testProvider{}, &testEvents{}, &responder.testResponder)
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

	f := buildTestInception(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
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

	f := buildTestInception(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
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

func (u *testVerifiableUser) HasVerifiedEmail() bool          { return u.verified }
func (u *testVerifiableUser) MarkEmailAsVerified(_ time.Time) {}
func (u *testVerifiableUser) MarkEmailAsUnverified()          {}
func (u *testVerifiableUser) GetEmailForVerification() string { return u.email }
