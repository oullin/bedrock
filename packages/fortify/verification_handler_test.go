package fortify

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// --- test responder extension ---

type testVerificationResponder struct {
	testResponder
	verificationSentCalled bool
}

// --- test verifier with error support ---

type testVerifierWithError struct {
	testVerifier
	verifyError error
}

// --- tests: send verification ---

// --- tests: verify email ---

// --- failing verifier ---

type failingVerifier struct {
	err error
}

func (r *testVerificationResponder) EmailVerificationSentResponse(w http.ResponseWriter, _ *http.Request) {
	r.verificationSentCalled = true
	w.WriteHeader(http.StatusOK)
}

func (v *testVerifierWithError) Verify(_ context.Context, _ string, _ string) error {
	return v.verifyError
}

func TestSendVerificationSuccess(t *testing.T) {
	user := &testUser{id: "1", email: "user@example.com"}
	guard := &testGuard{authenticatedUser: user}
	verifier := &testVerifier{}
	responder := &testVerificationResponder{}

	f := buildTestFortify(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.EmailVerification = true
	f.verifier = verifier
	f.responder = responder

	handler := NewSendVerificationHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/email/verification-notification", nil)

	handler.ServeHTTP(w, r)

	if !verifier.sent {
		t.Fatal("expected verification notification to be sent")
	}

	if !responder.verificationSentCalled {
		t.Fatal("expected verification sent response")
	}
}

func TestSendVerificationDisabled(t *testing.T) {
	f := buildTestFortify(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.EmailVerification = false

	handler := NewSendVerificationHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/email/verification-notification", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestSendVerificationUnauthenticated(t *testing.T) {
	guard := &testGuard{authenticatedUser: nil}
	responder := &testVerificationResponder{}

	f := buildTestFortify(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.EmailVerification = true
	f.verifier = &testVerifier{}
	f.responder = responder

	handler := NewSendVerificationHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/email/verification-notification", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestSendVerificationAlreadyVerified(t *testing.T) {
	user := &testVerifiableUser{
		testUser: testUser{id: "1", email: "user@example.com"},
		verified: true,
	}
	guard := &testGuard{authenticatedUser: user}
	verifier := &testVerifier{}
	responder := &testVerificationResponder{}

	f := buildTestFortify(guard, &testProvider{}, &testEvents{}, &responder.testResponder)
	f.config.Features.EmailVerification = true
	f.verifier = verifier
	f.responder = responder

	handler := NewSendVerificationHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/email/verification-notification", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	if verifier.sent {
		t.Fatal("should not send verification for already-verified user")
	}
}

func TestVerifyEmailSuccess(t *testing.T) {
	verifier := &testVerifier{}
	events := &testEvents{}

	f := buildTestFortify(&testGuard{}, &testProvider{}, events, &testResponder{})
	f.config.Features.EmailVerification = true
	f.verifier = verifier

	handler := NewVerifyEmailHandler(f)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.Handle("GET /verify-email/{id}/{hash}", handler)

	r := httptest.NewRequest(http.MethodGet, "/verify-email/123/abc123hash", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(events.dispatched) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events.dispatched))
	}

	if s, ok := events.dispatched[0].(string); !ok || s != EventVerified {
		t.Fatalf("expected Verified event string, got %T", events.dispatched[0])
	}
}

func TestVerifyEmailDisabled(t *testing.T) {
	f := buildTestFortify(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.EmailVerification = false

	handler := NewVerifyEmailHandler(f)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.Handle("GET /verify-email/{id}/{hash}", handler)

	r := httptest.NewRequest(http.MethodGet, "/verify-email/123/abc", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestVerifyEmailInvalidLink(t *testing.T) {
	f := buildTestFortify(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.EmailVerification = true
	f.verifier = &testVerifier{}

	handler := NewVerifyEmailHandler(f)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/verify-email/123/abc", nil)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing path values, got %d", w.Code)
	}
}

func TestVerifyEmailVerifierFails(t *testing.T) {
	verifier := &failingVerifier{err: errors.New("invalid hash")}

	f := buildTestFortify(&testGuard{}, &testProvider{}, &testEvents{}, &testResponder{})
	f.config.Features.EmailVerification = true
	f.verifier = verifier

	handler := NewVerifyEmailHandler(f)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.Handle("GET /verify-email/{id}/{hash}", handler)

	r := httptest.NewRequest(http.MethodGet, "/verify-email/123/badhash", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func (v *failingVerifier) SendVerificationNotification(_ context.Context, _ cauth.Authenticatable) error {
	return nil
}

func (v *failingVerifier) Verify(_ context.Context, _ string, _ string) error {
	return v.err
}
